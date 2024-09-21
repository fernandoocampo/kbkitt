package adds

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
)

// ui model
type model struct {
	inputs  []textinput.Model
	focused int
	err     error
}

// ui form fields
const (
	key = iota
	kind
	value
	notes
	reference
	tags
)

// ui constants
const (
	hotGreen = lipgloss.Color("#3aeb34")
	darkGray = lipgloss.Color("#767676")
)

// ui style
var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotGreen)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

func runInteractive() error {
	p := tea.NewProgram(initialModel())

	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("unable to run interactive mode: %w", err)
	}

	return nil
}

func initialModel() model {
	var inputs []textinput.Model = make([]textinput.Model, 6)
	inputs[key] = textinput.New()
	inputs[key].Placeholder = "key"
	inputs[key].Focus()
	inputs[key].CharLimit = 64
	inputs[key].Width = 70
	inputs[key].Prompt = ""
	inputs[key].SetValue(addKBData.key)

	inputs[kind] = textinput.New()
	inputs[kind].Placeholder = "category"
	inputs[kind].CharLimit = 64
	inputs[kind].Width = 70
	inputs[kind].Prompt = ""
	inputs[kind].SetValue(addKBData.kind)

	inputs[value] = textinput.New()
	inputs[value].Placeholder = "values"
	inputs[value].Width = 100
	inputs[value].Prompt = ""
	inputs[value].SetValue(addKBData.value)

	inputs[notes] = textinput.New()
	inputs[notes].Placeholder = ""
	inputs[notes].Width = 100
	inputs[notes].Prompt = ""
	inputs[notes].SetValue(addKBData.notes)

	inputs[reference] = textinput.New()
	inputs[reference].Placeholder = ""
	inputs[reference].CharLimit = 64
	inputs[reference].Width = 70
	inputs[reference].Prompt = ""
	inputs[reference].SetValue(addKBData.reference)

	inputs[tags] = textinput.New()
	inputs[tags].Placeholder = "keyword1 keyword2 keyword3 keywordN"
	inputs[tags].CharLimit = 100
	inputs[tags].Prompt = ""
	inputs[tags].SetValue(addKBData.reference)

	return model{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				m.toAddKBParams()
				return m, tea.Quit
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			os.Exit(0)
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf(
		` Adding a new KB:

 %s
 %s

 %s
 %s

 %s
 %s

 %s
 %s

 %s
 %s

 %s
 %s

 %s
`,
		inputStyle.Width(30).Render("Key"),
		m.inputs[key].View(),
		inputStyle.Width(8).Render("Category"),
		m.inputs[kind].View(),
		inputStyle.Width(6).Render("Value"),
		m.inputs[value].View(),
		inputStyle.Width(6).Render("Notes"),
		m.inputs[notes].View(),
		inputStyle.Width(9).Render("Reference"),
		m.inputs[reference].View(),
		inputStyle.Width(9).Render("Tags"),
		m.inputs[tags].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

// nextInput focuses the next input field
func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func (m *model) toAddKBParams() {
	addKBData.key = m.inputs[key].Value()
	addKBData.kind = m.inputs[kind].Value()
	addKBData.value = m.inputs[value].Value()
	addKBData.notes = m.inputs[notes].Value()
	addKBData.reference = m.inputs[reference].Value()
	addKBData.tags = m.convertTagsToArray()
}

func (m *model) convertTagsToArray() []string {
	if kbs.IsStringEmpty(m.inputs[tags].Value()) {
		return nil
	}

	return strings.Split(m.inputs[tags].Value(), " ")
}
