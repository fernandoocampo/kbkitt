package gets

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"golang.design/x/clipboard"
)

// mode defines get mode
type mode int

type filterView struct {
	inputs  [4]cmds.InputComponent
	focused int
}

type itemView struct {
	selectedItem *kbs.KB
	itemViewport *viewport.Model
}

type searchView struct {
	result    *kbs.SearchResult
	table     table.Model
	paginator *paginator.Model
}

type model struct {
	mode       mode
	searchView *searchView
	itemView   *itemView
	filterView *filterView
	service    *kbs.Service
	ctx        context.Context
	message    string
}

const (
	searchMode mode = iota
	filterMode
	itemMode
)

// ui form fields for filter view
const (
	category = iota
	namespace
	key
	keyword
)

const (
	hotGreen = lipgloss.Color("#3aeb34")
	darkGray = lipgloss.Color("#767676")
)

var (
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
	inputStyle    = lipgloss.NewStyle().Foreground(hotGreen)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

func runInteractive(ctx context.Context, service *kbs.Service) error {
	model := newModel(ctx, service)

	model.filterView = newFilterViewModel()

	itemViewPort, err := newItemViewport()
	if err != nil {
		return fmt.Errorf("unable to run interactive: %w", err)
	}

	model.itemView.itemViewport = itemViewPort

	p := tea.NewProgram(
		model,
	)

	_, err = p.Run()
	if err != nil {
		return fmt.Errorf("unable to run interactive mode: %w", err)
	}

	return nil
}

func newModel(ctx context.Context, service *kbs.Service) *model {
	newModel := model{
		service:    service,
		searchView: &searchView{},
		itemView:   &itemView{},
		ctx:        ctx,
		mode:       filterMode,
	}

	return &newModel
}

func newItemViewport() (*viewport.Model, error) {
	const width = 98

	vp := viewport.New(width, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil, err
	}

	str, err := renderer.Render("")
	if err != nil {
		return nil, err
	}

	vp.SetContent(str)

	return &vp, nil
}

func newPaginator(limit uint32, total int) *paginator.Model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = int(limit)
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(int(total))

	return &p
}

func newFilterViewModel() *filterView {
	var inputs [4]cmds.InputComponent

	categoryInput := textinput.New()
	categoryInput.Placeholder = "category"
	categoryInput.CharLimit = 64
	categoryInput.Width = 70
	categoryInput.Prompt = ""
	categoryInput.Focus()
	categoryInput.SetValue(getKBData.category)
	inputs[category].TextInput = &categoryInput

	namespaceInput := textinput.New()
	namespaceInput.Placeholder = "namespace"
	namespaceInput.CharLimit = 64
	namespaceInput.Width = 70
	namespaceInput.Prompt = ""
	namespaceInput.SetValue(getKBData.namespace)
	inputs[namespace].TextInput = &namespaceInput

	keyInput := textinput.New()
	keyInput.Placeholder = "key"
	keyInput.CharLimit = 64
	keyInput.Width = 70
	keyInput.Prompt = ""
	keyInput.SetValue(getKBData.key)
	inputs[key].TextInput = &keyInput

	keywordInput := textinput.New()
	keywordInput.Placeholder = "key"
	keywordInput.CharLimit = 64
	keywordInput.Width = 70
	keywordInput.Prompt = ""
	keywordInput.SetValue(getKBData.keyword)
	inputs[keyword].TextInput = &keywordInput

	filterView := filterView{
		inputs:  inputs,
		focused: 0,
	}

	return &filterView
}

func (m *model) updateTable() {
	keyLength := kbs.GetLongerText(cmds.KeyCol, m.searchView.result.Keys())
	categoryLength := kbs.GetLongerText(cmds.CategoryCol, m.searchView.result.Categories())
	namespaceLength := kbs.GetLongerText(cmds.NamespaceCol, m.searchView.result.Namespaces())
	tagLength := kbs.GetLongerText(cmds.TagCol, m.searchView.result.Tags())

	columns := []table.Column{
		{Title: cmds.KeyCol, Width: keyLength},
		{Title: cmds.CategoryCol, Width: categoryLength},
		{Title: cmds.NamespaceCol, Width: namespaceLength},
		{Title: cmds.TagCol, Width: tagLength},
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(toTableRow(m.searchView.result.Items)),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	t.SetStyles(s)

	m.searchView.table = t
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.filterView.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlQ:
			return m, tea.Quit
		case tea.KeyCtrlC:
			m.copyToClipboard()
			return m, cmd
		case tea.KeyCtrlF:
			switch m.mode {
			case filterMode:
				m.mode = searchMode
				getKBData.offset = 0
				err := m.search()
				if err != nil {
					fmt.Fprintln(os.Stderr, "searching kbs: %w", err)
					return m, tea.Quit
				}
				m.message = fmt.Sprintf("%d", m.searchView.result.Total)
			case searchMode:
				m.mode = filterMode
			}
		case tea.KeyCtrlO:
			m.openBrowser()
			return m, cmd
		case tea.KeyLeft:
			if (int(getKBData.offset) - int(getKBData.limit)) < 0 {
				return m, cmd
			}
			getKBData.offset = getKBData.offset - getKBData.limit

			err := m.searchKBItems()
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to search: %w", err)
				return m, tea.Quit
			}
		case tea.KeyRight:
			if (uint32(getKBData.offset) + getKBData.limit) >= uint32(m.searchView.result.Total) {
				return m, cmd
			}
			getKBData.offset = getKBData.offset + getKBData.limit
			err := m.searchKBItems()
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to search: %w", err)
				return m, tea.Quit
			}
		case tea.KeyDown:
			m.searchView.table, cmd = m.searchView.table.Update(msg)
			return m, cmd
		case tea.KeyUp:
			m.searchView.table, cmd = m.searchView.table.Update(msg)
			return m, cmd
		case tea.KeyCtrlR:
			m.mode = searchMode
			m.itemView.selectedItem = nil
		case tea.KeyShiftTab, tea.KeyCtrlP:
			if m.mode != filterMode {
				return m, cmd
			}
			m.mode = filterMode
			m.filterView.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			if m.mode != filterMode {
				return m, cmd
			}
			if m.filterView.focused == len(m.filterView.inputs)-1 {
				m.mode = searchMode
				getKBData.offset = 0
				err := m.search()
				if err != nil {
					fmt.Fprintln(os.Stderr, "searching kbs: %w", err)
					return m, tea.Quit
				}
				m.message = fmt.Sprintf("%d", m.searchView.result.Total)
			} else {
				m.filterView.nextInput()
			}
		case tea.KeyEnter:
			if m.mode != searchMode {
				return m, cmd
			}
			selectedRowKey := m.searchView.table.SelectedRow()[0]
			err := m.loadKBItem(selectedRowKey)
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to search: %w", err)
				return m, tea.Quit
			}
			m.mode = itemMode
			newItemViewport, cmd := m.itemView.itemViewport.Update(msg)
			newItemViewport.SetContent(m.content())
			m.itemView.itemViewport = &newItemViewport
			return m, cmd
		}
	default:
		return m, nil
	}

	switch m.mode {
	case searchMode:
		m.updateTable()
		newpaginator, cmd := m.searchView.paginator.Update(msg)
		m.searchView.paginator = &newpaginator
		m.searchView.table.UpdateViewport()

		return m, cmd
	case filterMode:
		for i := range m.filterView.inputs {
			m.filterView.inputs[i].Blur()
		}
		m.filterView.inputs[m.filterView.focused].Focus()

		for i := range m.filterView.inputs {
			if m.filterView.inputs[i].TextInput != nil {
				textInputModel, textInputCmd := m.filterView.inputs[i].TextInput.Update(msg)
				m.filterView.inputs[i].TextInput, cmds[i] = &textInputModel, textInputCmd
			}
		}
		return m, tea.Batch(cmds...)
	default:
		return m, cmd
	}
}

func (m *model) View() string {
	switch m.mode {
	case itemMode:
		return m.drawKBViewer()
	case filterMode:
		return m.renderFilters()
	case searchMode:
		fallthrough
	default:
		return m.drawTable()
	}
}

func (m *model) drawKBViewer() string {
	return m.itemView.itemViewport.View() + m.helpView()
}

func (m *model) helpView() string {
	if m.itemView.selectedItem == nil {
		return helpStyle("\n  • Ctrl+R: Back • q: Quit\n")
	}

	if m.itemView.selectedItem.Category == kbs.BookmarkCategory {
		return helpStyle("\n  ↑/↓: Navigate • Ctrl+R: Back • Ctrl+c: Copy • Ctrl+o: Open • Esc: Quit\n")
	}

	return helpStyle("\n  ↑/↓: Navigate • Ctrl+R: Back • Ctrl+c: Copy • Esc: Quit\n")
}

func (m *model) drawTable() string {
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
	var b strings.Builder
	b.WriteString("\n  Knowledge Base Results: ")
	b.WriteString(m.message)
	b.WriteString("\n\n")
	b.WriteString(baseStyle.Render(m.searchView.table.View()))
	b.WriteString("\n\n")
	b.WriteString("  " + m.searchView.paginator.View())
	b.WriteString("\n\n  ←/→ page • Ctrl+F: filters • Esc: quit\n")
	return b.String()
}

func (m *model) search() error {
	m.toGetKBParams()
	m.itemView.selectedItem = nil

	err := m.searchKBItems()
	if err != nil {
		return fmt.Errorf("unable to search: %w", err)
	}

	m.searchView.paginator = newPaginator(m.searchView.result.Limit, m.searchView.result.Total)

	return nil
}

func (m *model) searchKBItems() error {
	result, err := m.service.Search(m.ctx, getKBData.toKBQueryFilter())
	if err != nil {
		return fmt.Errorf("unable to search kb items: %w", err)
	}

	if result == nil {
		m.searchView.result = &kbs.SearchResult{}
		return nil
	}

	m.searchView.result = result

	return nil
}

func (m *model) loadKBItem(kbKey string) error {
	kb, err := m.service.GetByKey(m.ctx, kbKey)
	if err != nil {
		return fmt.Errorf("unable to get kb: %w", err)
	}

	m.itemView.selectedItem = kb

	return nil
}

func toTableRow(items []kbs.KBItem) []table.Row {
	result := make([]table.Row, 0, len(items))

	for _, v := range items {
		result = append(result, v.ToArray())
	}
	return result
}

func (m *model) content() string {
	return renderKBItem(m.itemView.selectedItem)
}

func renderKBItem(k *kbs.KB) string {
	return fmt.Sprintf(`%s
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
%s
%+v
`,
		inputStyle.Width(30).Render("ID"), k.ID,
		inputStyle.Width(30).Render("Key"), k.Key,
		inputStyle.Width(30).Render("Category"), k.Category,
		inputStyle.Width(30).Render("Namespace"), k.Namespace,
		inputStyle.Width(30).Render("Value"), k.Value,
		inputStyle.Width(30).Render("Notes"), k.Notes,
		inputStyle.Width(30).Render("Reference"), k.Reference,
		inputStyle.Width(30).Render("Tags"), k.Tags)
}

func (m model) renderFilters() string {
	return fmt.Sprintf(
		` Criteria to search kbs:

%s
%s

%s
%s

%s
%s

%s
%s

%s

• tab: next • shift+tab: previous • Ctrl+F: find • Esc: quit

`,

		inputStyle.Width(8).Render(kbs.CategoryLabel),
		m.filterView.inputs[category].View(),
		inputStyle.Width(9).Render(kbs.NamespaceLabel),
		m.filterView.inputs[namespace].View(),
		inputStyle.Width(6).Render(kbs.KeyLabel),
		m.filterView.inputs[key].View(),
		inputStyle.Width(9).Render(kbs.KeywordLabel),
		m.filterView.inputs[keyword].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func (m *model) copyToClipboard() {
	if m.itemView.selectedItem == nil {
		return
	}

	err := clipboard.Init()
	if err != nil {
		// let's ignore error
		return
	}

	var value string

	switch m.itemView.selectedItem.Category {
	case kbs.QuoteCategory:
		value = fmt.Sprintf("%q ~ %s", m.itemView.selectedItem.Value,
			m.itemView.selectedItem.Reference)
	default:
		value = m.itemView.selectedItem.Value
	}

	_ = clipboard.Write(clipboard.FmtText, []byte(value))
}

func (m *model) openBrowser() {
	if m.itemView.selectedItem == nil || m.itemView.selectedItem.Category != kbs.BookmarkCategory {
		return
	}

	url := m.itemView.selectedItem.Value

	switch runtime.GOOS {
	case "linux":
		_ = exec.Command("xdg-open", url).Start()
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		_ = exec.Command("open", url).Start()
	default:
		return
	}
}

// nextInput focuses the next input field
func (f *filterView) nextInput() {
	f.focused = (f.focused + 1) % len(f.inputs)
}

// prevInput focuses the previous input field
func (f *filterView) prevInput() {
	f.focused--
	// Wrap around
	if f.focused < 0 {
		f.focused = 1
	}
}

func (m *model) toGetKBParams() {
	getKBData.category = strings.ToLower(m.filterView.inputs[category].Value())
	getKBData.namespace = strings.ToLower(m.filterView.inputs[namespace].Value())
	getKBData.key = strings.ToLower(m.filterView.inputs[key].Value())
	getKBData.keyword = strings.ToLower(m.filterView.inputs[keyword].Value())
}
