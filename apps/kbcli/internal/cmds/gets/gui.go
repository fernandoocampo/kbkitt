package gets

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
)

// mode defines get mode
type mode int

type model struct {
	mode         mode
	result       *kbs.SearchResult
	table        table.Model
	paginator    *paginator.Model
	service      *kbs.Service
	ctx          context.Context
	selectedItem *kbs.KB
	itemViewport *viewport.Model
}

const (
	searchMode mode = iota
	itemMode
)

const hotGreen = lipgloss.Color("#3aeb34")

var (
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
	inputStyle = lipgloss.NewStyle().Foreground(hotGreen)
)

func runInteractive(ctx context.Context, service *kbs.Service) error {
	model := newModel(ctx, service)

	err := model.searchKBItems()
	if err != nil {
		return fmt.Errorf("unable to run interactive: %w", err)
	}

	if model.empty() {
		fmt.Println()
		fmt.Println("zero occurrences with that filter")
		return nil
	}

	itemViewPort, err := newItemViewport()
	if err != nil {
		return fmt.Errorf("unable to run interactive: %w", err)
	}

	model.itemViewport = itemViewPort
	model.paginator = newPaginator(model.result.Limit, model.result.Total)
	model.updateTable()

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
		service: service,
		ctx:     ctx,
		mode:    searchMode,
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

func newPaginator(limit uint16, total int) *paginator.Model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = int(limit)
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(int(total))

	return &p
}

func (m *model) updateTable() {
	keyLength := m.result.GetLongerKey()
	kindLength := m.result.GetLongerKind()
	tagLength := m.result.GetLongerTags()

	columns := []table.Column{
		{Title: cmds.IDCol, Width: 36},
		{Title: cmds.KeyCol, Width: keyLength},
		{Title: cmds.KindCol, Width: kindLength},
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
		table.WithRows(toTableRow(m.result.Items)),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	t.SetStyles(s)

	m.table = t
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case tea.KeyLeft.String():
			if (int(getKBData.offset) - int(getKBData.limit)) < 0 {
				return m, cmd
			}
			getKBData.offset = getKBData.offset - getKBData.limit

			err := m.searchKBItems()
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to search: %w", err)
				os.Exit(1)
			}
		case tea.KeyRight.String():
			if (uint32(getKBData.offset) + getKBData.limit) >= uint32(m.result.Total) {
				return m, cmd
			}
			getKBData.offset = getKBData.offset + getKBData.limit
			err := m.searchKBItems()
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to search: %w", err)
				os.Exit(1)
			}
		case tea.KeyDown.String():
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		case tea.KeyUp.String():
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		case tea.KeyCtrlR.String():
			m.mode = searchMode
			m.selectedItem = nil
		case tea.KeyEnter.String():
			selectedRowID := m.table.SelectedRow()[0]
			err := m.loadKBItem(selectedRowID)
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to search: %w", err)
				os.Exit(1)
			}
			m.mode = itemMode
			newItemViewport, cmd := m.itemViewport.Update(msg)
			newItemViewport.SetContent(m.content())
			m.itemViewport = &newItemViewport
			return m, cmd
		}
	default:
		return m, nil
	}
	m.updateTable()
	m.table.UpdateViewport()
	newpaginator, cmd := m.paginator.Update(msg)
	m.paginator = &newpaginator

	return m, cmd
}

func (m *model) View() string {
	switch m.mode {
	case itemMode:
		return m.drawKBViewer()
	case searchMode:
		fallthrough
	default:
		return m.drawTable()
	}
}

func (m *model) drawKBViewer() string {
	return m.itemViewport.View() + m.helpView()
}

func (m *model) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate • Ctrl+R: Return • q: Quit\n")
}

func (m *model) drawTable() string {
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
	var b strings.Builder
	b.WriteString("\n  Knowledge Base Found \n\n")
	b.WriteString(baseStyle.Render(m.table.View()))
	b.WriteString("\n\n")
	b.WriteString("  " + m.paginator.View())
	b.WriteString("\n\n  ←/→ page • q: quit\n")
	return b.String()
}

func (m *model) searchKBItems() error {
	result, err := m.service.Search(m.ctx, getKBData.toKBQueryFilter())
	if err != nil {
		return fmt.Errorf("unable to search: %w", err)
	}

	m.result = result

	return nil
}

func (m *model) loadKBItem(kbID string) error {
	kb, err := m.service.GetByID(m.ctx, kbID)
	if err != nil {
		return fmt.Errorf("unable to get kb: %w", err)
	}

	m.selectedItem = kb

	return nil
}

func (m *model) empty() bool {
	return len(m.result.Items) == 0
}

func (m *model) formatItems() []string {
	var items []string
	keyLength := m.result.GetLongerKey()
	kindLength := m.result.GetLongerKind()

	for _, kb := range m.result.Items {
		row := fmt.Sprintln(kb.ID, fmt.Sprintf("%s%*s", kb.Key, keyLength-len(kb.Key), ""), fmt.Sprintf("%s%*s", kb.Kind, kindLength-len(kb.Kind), ""), strings.Join(kb.Tags, ","))
		items = append(items, row)
	}

	return items
}

func toTableRow(items []kbs.KBItem) []table.Row {
	result := make([]table.Row, 0, len(items))

	for _, v := range items {
		result = append(result, v.ToArray())
	}
	return result
}

func (m *model) content() string {
	return renderKBItem(m.selectedItem)
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
%+v
`,
		inputStyle.Width(30).Render("ID"), k.ID,
		inputStyle.Width(30).Render("Key"), k.Key,
		inputStyle.Width(30).Render("Value"), k.Value,
		inputStyle.Width(30).Render("Notes"), k.Notes,
		inputStyle.Width(30).Render("Kind"), k.Kind,
		inputStyle.Width(30).Render("Reference"), k.Reference,
		inputStyle.Width(30).Render("Tags"), k.Tags)
}
