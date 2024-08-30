package cmds

import (
	"fmt"
	"os"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/settings"
	"github.com/spf13/cobra"
)

const (
	startConfigurationMessage = "do you want to setup kbkitt? [y/n]: "
	hostLabel                 = "kbkitt host (http(s)://): "
)

const (
	apiVersion = "0.1.0"
)

func makeConfigureCommand() *cobra.Command {
	newCmd := cobra.Command{
		Use:   "configure",
		Short: "configure kb-kitt",
		Long:  "configure kb-kitt",
		Run:   makeRunConfigureCommand(),
	}

	return &newCmd
}

func makeRunConfigureCommand() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := settings.CheckAndCreateKBKittFolder()
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to verify kbkitt folder: %s", err)
			fmt.Println()
			os.Exit(1)
		}

		configuration, err := settings.LoadConfiguration()
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to load configuration: %s", err)
			fmt.Println()
			os.Exit(1)
		}

		if configuration != nil && !startConfiguration() {
			fmt.Println("ok")
			os.Exit(0)
		}

		err = settings.Save(newKBKitt())
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to save configuration: %s", err)
			fmt.Println()
			os.Exit(1)
		}

		fmt.Println("done")
	}
}

func startConfiguration() bool {
	var yesOrNot string
	fmt.Print(startConfigurationMessage)
	fmt.Scan(&yesOrNot)

	if yes(yesOrNot) {
		return true
	}

	return false
}

func newKBKitt() *settings.Configuration {
	var newConfiguration settings.Configuration

	newConfiguration.Version = apiVersion
	newConfiguration.Server = &settings.Server{}

	newConfiguration.Server.URL = requestStringValue(hostLabel)

	return &newConfiguration
}
