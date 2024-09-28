package adds

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
)

type errMsg error

// ui model
type model struct {
	inputs  []cmds.InputComponent
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
	// var inputs []textinput.Model = make([]textinput.Model, 6)
	var inputs []cmds.InputComponent = make([]cmds.InputComponent, 6)
	keyInput := textinput.New()
	keyInput.Placeholder = "key"
	keyInput.Focus()
	keyInput.CharLimit = 64
	keyInput.Width = 70
	keyInput.Prompt = ""
	keyInput.SetValue(addKBData.key)
	inputs[key].TextInput = &keyInput

	kindInput := textinput.New()
	kindInput.Placeholder = "category"
	kindInput.CharLimit = 64
	kindInput.Width = 70
	kindInput.Prompt = ""
	kindInput.SetValue(addKBData.kind)
	inputs[kind].TextInput = &kindInput

	valueInput := textarea.New()
	valueInput.Placeholder = "..."
	valueInput.Prompt = ""
	valueInput.CharLimit = 700
	valueInput.ShowLineNumbers = false
	valueInput.SetHeight(4)
	valueInput.SetWidth(80)
	valueInput.SetValue(addKBData.value)
	inputs[value].TextArea = &valueInput

	notesInput := textarea.New()
	notesInput.Placeholder = ""
	notesInput.Prompt = ""
	notesInput.CharLimit = 700
	notesInput.ShowLineNumbers = false
	notesInput.SetHeight(4)
	notesInput.SetWidth(80)
	notesInput.SetValue(addKBData.notes)
	inputs[notes].TextArea = &notesInput

	refInput := textinput.New()
	refInput.Placeholder = ""
	refInput.CharLimit = 64
	refInput.Width = 70
	refInput.Prompt = ""
	refInput.SetValue(addKBData.reference)
	inputs[reference].TextInput = &refInput

	tagsInput := textinput.New()
	tagsInput.Placeholder = "keyword1 keyword2 keyword3 keywordN"
	tagsInput.CharLimit = 100
	tagsInput.Prompt = ""
	tagsInput.SetValue(addKBData.reference)
	inputs[tags].TextInput = &tagsInput

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
		case tea.KeyCtrlC, tea.KeyEsc:
			os.Exit(0)
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			if m.focused == len(m.inputs)-1 {
				m.toAddKBParams()
				return m, tea.Quit
			}
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
		if m.inputs[i].TextInput != nil {
			textInputModel, textInputCmd := m.inputs[i].TextInput.Update(msg)
			m.inputs[i].TextInput, cmds[i] = &textInputModel, textInputCmd
			continue
		}
		textInputModel, textInputCmd := m.inputs[i].TextArea.Update(msg)
		m.inputs[i].TextArea, cmds[i] = &textInputModel, textInputCmd
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

• tab fields • shift+tab fields • ctrl+c: quit

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
