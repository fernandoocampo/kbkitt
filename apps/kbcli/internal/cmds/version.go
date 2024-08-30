package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version    string
	BuildDate  string
	CommitHash string
)

// versionFormat x.y.z (commit_hash commit_date)
const versionFormat = "kbcli %s (%s %s)"

func makeVersionCommand() *cobra.Command {
	newCmd := cobra.Command{
		Use:   "version",
		Short: "print the version number of kb-kitt",
		Long:  "print the version number of kb-kitt",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(versionFormat, Version, CommitHash, BuildDate)
		},
	}

	return &newCmd
}
