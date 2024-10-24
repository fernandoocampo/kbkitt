package apps

import (
	"errors"
	"fmt"
	"os"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/kbkitt"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/storages"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/adds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/gets"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/imports"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/setups"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/syncs"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/updates"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/versions"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/settings"
	"github.com/spf13/cobra"
)

type Application struct {
	rootCommand   *cobra.Command
	configuration *settings.Configuration
	storage       *storages.SQLite
	service       *kbs.Service
	kbkitClient   *kbkitt.Client
}

func NewApplication() *Application {
	newApp := Application{}

	return &newApp
}

func makeRootCommand() *cobra.Command {
	newCmd := cobra.Command{
		Use:   "kb",
		Short: "kb is a knowledge base manager",
		Long:  `A knowledge base manager to manage the concepts you use every day.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				fmt.Println(err)
			}
		},
	}

	return &newCmd
}

func (a *Application) Execute() error {
	a.rootCommand = makeRootCommand()
	a.rootCommand.AddCommand(versions.MakeVersionCommand())
	a.rootCommand.AddCommand(setups.MakeConfigureCommand())

	err := a.initializeConfiguration()
	if err != nil && !errors.Is(err, cmds.ErrNoConfiguration) {
		return fmt.Errorf("unable to start kbkitt: %w", err)
	}

	err = a.initializeStorage()
	if err != nil {
		return fmt.Errorf("unable to start kbkitt: %w", err)
	}

	defer a.storage.Close()

	a.initializeKBKittClient()

	err = a.initializeService()
	if err != nil {
		return fmt.Errorf("unable to start kbkitt: %w", err)
	}

	err = a.initializeRootCommand()
	if err != nil {
		return fmt.Errorf("unable to initialize service commands: %w", err)
	}

	if err := a.rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return fmt.Errorf("unable to execute app")
	}

	return nil
}

func (a *Application) initializeConfiguration() error {
	configuration, err := cmds.GetConfiguration()
	if err != nil && !errors.Is(err, cmds.ErrNoConfiguration) {
		return fmt.Errorf("unable to load configuration: %w", err)
	}

	a.configuration = configuration

	return nil
}

func (a *Application) initializeStorage() error {
	storage, err := cmds.NewStorage(a.configuration)
	if err != nil {
		return fmt.Errorf("unable to load service: %w", err)
	}

	a.storage = storage

	return nil
}

func (a *Application) initializeKBKittClient() {
	kbkittSetup := kbkitt.Setup{
		URL: a.configuration.Server.URL,
	}

	a.kbkitClient = kbkitt.NewClient(kbkittSetup)
}

func (a *Application) initializeService() error {
	serviceSetup := kbs.ServiceSetup{
		KBClient:        a.kbkitClient,
		KBStorage:       a.storage,
		FileForSyncPath: a.configuration.FileForSyncPath,
		DirForMediaPath: a.configuration.DirForMediaPath,
	}

	a.service = kbs.NewService(serviceSetup)

	return nil
}

func (a *Application) initializeRootCommand() error {
	a.rootCommand.AddCommand(adds.MakeAddCommand(a.service))
	a.rootCommand.AddCommand(imports.MakeImportCommand(a.service))
	a.rootCommand.AddCommand(gets.MakeGetCommand(a.service))
	a.rootCommand.AddCommand(syncs.MakeSyncCommand(a.service))
	a.rootCommand.AddCommand(updates.MakeUpdateCommand(a.service))

	return nil
}
