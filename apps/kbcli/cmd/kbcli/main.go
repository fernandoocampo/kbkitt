package main

import (
	"fmt"
	"os"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/apps"
)

func main() {
	err := apps.NewApplication().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to start app:", err)
		os.Exit(1)
	}
}
