package apps

import (
	"errors"
	"fmt"
	"os"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/kbkitt"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/storages"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/adds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/exports"
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
	// isItSet indicates if application is configured to run.
	isItSet bool
	// rootCommand is the main command of this CLI app.
	rootCommand *cobra.Command
	// configuration contains the runtime settings for this app.
	configuration *settings.Configuration
	// storage reference the storage mechanism for this app.
	storage *storages.SQLite
	// service is the object in charge of handling business logic.
	service *kbs.Service
	// kbkitClient provides logic related to the central kb server.
	kbkitClient *kbkitt.Client
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

	a.initializeRootCommand()

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

	if err != nil && errors.Is(err, cmds.ErrNoConfiguration) {
		a.isItSet = false
		return nil
	}

	a.isItSet = true

	a.configuration = configuration

	return nil
}

func (a *Application) initializeStorage() error {
	if !a.itIsSet() {
		return nil
	}

	storage, err := cmds.NewStorage(a.configuration)
	if err != nil {
		return fmt.Errorf("unable to load service: %w", err)
	}

	a.storage = storage

	return nil
}

func (a *Application) initializeKBKittClient() {
	if !a.itIsSet() {
		return
	}

	kbkittSetup := kbkitt.Setup{
		URL: a.configuration.Server.URL,
	}

	a.kbkitClient = kbkitt.NewClient(kbkittSetup)
}

func (a *Application) initializeService() error {
	if !a.itIsSet() {
		return nil
	}

	serviceSetup := kbs.ServiceSetup{
		KBClient:        a.kbkitClient,
		KBStorage:       a.storage,
		FileForSyncPath: a.configuration.FileForSyncPath,
		DirForMediaPath: a.configuration.DirForMediaPath,
	}

	a.service = kbs.NewService(serviceSetup)

	return nil
}

func (a *Application) initializeRootCommand() {
	if !a.itIsSet() {
		return
	}

	a.rootCommand.AddCommand(adds.MakeAddCommand(a.service))
	a.rootCommand.AddCommand(imports.MakeImportCommand(a.service))
	a.rootCommand.AddCommand(exports.MakeExportCommand(a.service))
	a.rootCommand.AddCommand(gets.MakeGetCommand(a.service))
	a.rootCommand.AddCommand(syncs.MakeSyncCommand(a.service))
	a.rootCommand.AddCommand(updates.MakeUpdateCommand(a.service))
}

func (a *Application) itIsSet() bool {
	return a.isItSet
}
