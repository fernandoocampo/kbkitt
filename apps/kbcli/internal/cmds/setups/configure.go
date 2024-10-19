package setups

import (
	"context"
	"fmt"
	"os"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/settings"
	"github.com/spf13/cobra"
)

const (
	startConfigurationMessage = "> do you want to setup kbkitt? [y/n]: "
	hostLabel                 = "kbkitt host (http(s)://): "
	filePathForSyncLabel      = "file path to save kbs for synchronization: "
	dirForMediaLabel          = "dir path to save kb media files: "
)

const (
	apiVersion      = "0.1.0"
	filePathForSync = "0.1.0"
)

func MakeConfigureCommand() *cobra.Command {
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

		ctx := context.Background()

		err = saveConfiguration(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to save configuration: %s", err)
			fmt.Println()
			os.Exit(1)
		}

		fmt.Println("done")
	}
}

func saveConfiguration(ctx context.Context) error {
	newConfiguration := newKBKitt()
	err := settings.Save(newConfiguration)
	if err != nil {
		return fmt.Errorf("unable to save file with given settings")
	}

	storage, err := cmds.NewStorage(newConfiguration)
	if err != nil {
		return fmt.Errorf("unable to initialize internal database: %w", err)
	}

	defer storage.Close()

	err = settings.CreateDatabaseIfNotExist(ctx, newConfiguration, storage)
	if err != nil {
		return fmt.Errorf("unable to initialize internal database: %w", err)
	}

	return nil
}

func startConfiguration() bool {
	var yesOrNot string
	fmt.Print(startConfigurationMessage)
	fmt.Scan(&yesOrNot)

	return cmds.Yes(yesOrNot)
}

func newKBKitt() *settings.Configuration {
	var newConfiguration settings.Configuration

	newConfiguration.Version = apiVersion
	newConfiguration.Server = &settings.Server{}

	newConfiguration.Server.URL = cmds.RequestStringValue(hostLabel)
	newConfiguration.FileForSyncPath = cmds.RequestStringValue(filePathForSyncLabel)
	newConfiguration.DirForMediaPath = cmds.RequestStringValue(dirForMediaLabel)

	return &newConfiguration
}
