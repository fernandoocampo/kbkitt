package cmds

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/storages"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/settings"
)

type InputType interface {
	Blur()
	Focus() tea.Cmd
	View() string
	Value() string
}

type InputComponent struct {
	TextInput *textinput.Model
	TextArea  *textarea.Model
}

// values
const (
	yesValue      = "yes"
	yesShortValue = "y"
)

// common labels
const (
	TitleSeparator       = "-------------"
	IDCol                = "ID"
	IDColSeparator       = "--"
	KeyCol               = "KEY"
	KeyColSeparator      = "---"
	CategoryCol          = "CATEGORY"
	CategoryColSeparator = "----"
	TagCol               = "TAGS"
	TagColSeparator      = "----"
	GetKBIDLabel         = "id: "
)

var ErrNoConfiguration = errors.New("no configuration has been created yet")

func GetConfiguration() (*settings.Configuration, error) {
	configuration, err := settings.LoadConfiguration()
	if err != nil {
		return nil, fmt.Errorf("unable get configuration: %w", err)
	}

	if configuration == nil {
		return nil, ErrNoConfiguration
	}

	if configuration.Invalid() {
		return nil, errors.New("kbkitt settings are not good, please verify")
	}

	return configuration, nil
}

func NewStorage(configuration *settings.Configuration) (*storages.SQLite, error) {
	sqlConn, err := storages.CreateSQLiteConnection(configuration.GetDBPath())
	if err != nil {
		return nil, fmt.Errorf("unable to create db connection: %w", err)
	}

	setup := storages.SQLiteSetup{
		DB: sqlConn,
	}

	return storages.NewSQLite(&setup), nil
}

func Yes(answer string) bool {
	return strings.EqualFold(answer, yesValue) || strings.EqualFold(answer, yesShortValue)
}

func No(answer string) bool {
	return !Yes(answer)
}

func RequestStringValue(label string) string {
	var output string
	fmt.Print(label)

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		output = scanner.Text()
	}

	return output
}

func AreYouSure(label string) bool {
	var done string
	fmt.Print(label)
	fmt.Scan(&done)

	return Yes(done)
}

func PrintKB(k *kbs.KB) string {
	return fmt.Sprintf(`%s:
%s
%s:
%s
%s:
%s
%s:
%s
%s:
%s
%s:
%s
%s:
%+v
`,
		kbs.IDLabel, k.ID,
		kbs.KeyLabel, k.Key,
		kbs.ValueLabel, k.Value,
		kbs.NotesLabel, k.Notes,
		kbs.CategoryLabel, k.Category,
		kbs.ReferenceLabel, k.Reference,
		kbs.TagsLabel, k.Tags)
}

func (i *InputComponent) Blur() {
	if i.TextArea != nil {
		i.TextArea.Blur()
		return
	}

	i.TextInput.Blur()
}

func (i *InputComponent) Focus() tea.Cmd {
	if i.TextArea != nil {
		return i.TextArea.Focus()
	}

	return i.TextInput.Focus()
}

func (i *InputComponent) View() string {
	if i.TextArea != nil {
		return i.TextArea.View()
	}

	return i.TextInput.View()
}

func (i *InputComponent) Value() string {
	if i.TextArea != nil {
		return i.TextArea.Value()
	}

	return i.TextInput.Value()
}
