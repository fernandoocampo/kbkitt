package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

func makeVersionCommand() *cobra.Command {
	newCmd := cobra.Command{
		Use:   "version",
		Short: "print the version number of kb-kitt",
		Long:  "print the version number of kb-kitt",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}

	return &newCmd
}
