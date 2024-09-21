package cmds

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/kbkitt"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/settings"
)

// values
const (
	yesValue      = "yes"
	yesShortValue = "y"
)

// common labels
const (
	AreYouDoneLabel  = "> are you done? [y/n]: "
	TitleSeparator   = "-------------"
	IDCol            = "ID"
	IDColSeparator   = "--"
	KeyCol           = "KEY"
	KeyColSeparator  = "---"
	KindCol          = "KIND"
	KindColSeparator = "----"
	TagCol           = "TAGS"
	TagColSeparator  = "----"
)

func getConfiguration() (*settings.Configuration, error) {
	configuration, err := settings.LoadConfiguration()
	if err != nil {
		return nil, fmt.Errorf("unable to load kbkitt settings: %w", err)
	}

	if configuration.Invalid() {
		return nil, errors.New("kbkitt settings are not good, please verify")
	}

	return configuration, nil
}

func NewService() (*kbs.Service, error) {
	configuration, err := getConfiguration()
	if err != nil {
		return nil, fmt.Errorf("unable to load configuration: %w", err)
	}

	service, err := getKBKittService(configuration)
	if err != nil {
		return nil, fmt.Errorf("unable to create service: %w", err)
	}

	return service, nil
}

func getKBKittService(conf *settings.Configuration) (*kbs.Service, error) {
	serviceSetup := kbs.ServiceSetup{
		KBClient:        newKBKittClient(conf),
		FileForSyncPath: conf.FilepathForSyncPath,
	}

	newService := kbs.NewService(serviceSetup)

	return newService, nil
}

func newKBKittClient(conf *settings.Configuration) *kbkitt.Client {
	kbkittSetup := kbkitt.Setup{
		URL: conf.Server.URL,
	}

	return kbkitt.NewClient(kbkittSetup)
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

func ReadCSVFromStdin(label string) []string {
	var result []string
	for {
		var value string
		fmt.Print(label)
		fmt.Scan(&value)

		result = append(result, value)

		if AreYouSure(AreYouDoneLabel) {
			fmt.Println()
			break
		}
	}
	return result
}

func AreYouSure(label string) bool {
	var done string
	fmt.Print(label)
	fmt.Scan(&done)

	return Yes(done)
}
