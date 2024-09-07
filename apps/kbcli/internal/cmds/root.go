package cmds

import (
	"fmt"
	"os"

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

func Execute() {
	rootCommand := makeRootCommand()
	rootCommand.AddCommand(makeVersionCommand())
	rootCommand.AddCommand(makeAddCommand())
	rootCommand.AddCommand(makeImportCommand())
	rootCommand.AddCommand(makeGetCommand())
	rootCommand.AddCommand(makeConfigureCommand())

	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
