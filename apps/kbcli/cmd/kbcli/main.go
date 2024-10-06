package main

import (
	"fmt"
	"os"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/apps"
)

func main() {
	// /Users/Fernando_Ocampo/Workspaces/gomodws/epicgames/uas-replacement-pocs/libauth/internal/apps/cmd.go
	err := apps.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to start app:", err)
		os.Exit(1)
	}
}
