package apps

import (
	"fmt"
	"os"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/adds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/gets"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/imports"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/setups"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/syncs"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/versions"
	"github.com/spf13/cobra"
)

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

func Execute() error {
	service, err := cmds.NewService()
	if err != nil {
		return fmt.Errorf("unable to load service: %w", err)
	}

	rootCommand := makeRootCommand()
	rootCommand.AddCommand(versions.MakeVersionCommand())
	rootCommand.AddCommand(adds.MakeAddCommand(service))
	rootCommand.AddCommand(imports.MakeImportCommand(service))
	rootCommand.AddCommand(gets.MakeGetCommand(service))
	rootCommand.AddCommand(setups.MakeConfigureCommand())
	rootCommand.AddCommand(syncs.MakeSyncCommand(service))

	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return fmt.Errorf("unable to execute app")
	}

	return nil
}
