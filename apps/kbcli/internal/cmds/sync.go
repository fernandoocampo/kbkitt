package cmds

import (
	"context"
	"fmt"
	"os"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/spf13/cobra"
)

// syncKBParams contains parameters required by sync command to sync locally saved kbs.
type syncKBParams struct {
	showFailedKBs bool
	showAddedKBs  bool
}

// field labels
const (
	syncedLabel             = "Synced KBs"
	notSyncedLabel          = "Not Synced KBs"
	notSyncedErrorLabel     = "ERROR"
	notSyncedErrorSeparator = "-----"
	totalSyncedLabel        = "Total:"
	totalNotSyncedLabel     = "Total:"
)

var syncKBData syncKBParams

func makeSyncCommand() *cobra.Command {
	newCmd := cobra.Command{
		Use:   "sync",
		Short: "sync locally saved kbs with the server",
		Long:  "sync locally saved kbs with the server, these are kbs that could not be saved due to some server errors",
		Run:   makeRunSyncCommand(),
	}

	newCmd.PersistentFlags().BoolVarP(&syncKBData.showAddedKBs, "show-added-kbs", "", false, "knowledge base key")
	newCmd.PersistentFlags().BoolVarP(&syncKBData.showFailedKBs, "show-failed-kbs", "", false, "knowledge base key")

	return &newCmd
}

func makeRunSyncCommand() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		service, err := newService()
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to load service: %s", err)
			fmt.Println()
			os.Exit(1)
		}

		ctx := context.Background()

		result, err := service.Sync(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to process synchronization:", err)
			os.Exit(1)
		}

		if result == nil || result.Empty() {
			fmt.Println("nothing was processed")
		}

		printSyncedReport(result)
	}
}

func printSyncedReport(kbs *kbs.SyncResult) {
	if kbs == nil {
		return
	}

	if syncKBData.dontPrint() {
		fmt.Println("report output was not requested")
		return
	}

	if syncKBData.showAddedKBs && len(kbs.NewIDs) > 0 {
		printSyncedKBs(kbs)
	}

	if syncKBData.showFailedKBs && len(kbs.FailedKeys) > 0 {
		printNotSyncedKBs(kbs)
	}
}

func printNotSyncedKBs(kbs *kbs.SyncResult) {
	length := len(keyCol)
	for key := range kbs.FailedKeys {
		if len(key) > length {
			length = len(key)
		}
	}

	fmt.Println()
	fmt.Println(notSyncedLabel)
	fmt.Println(titleSeparator)
	fmt.Println(totalNotSyncedLabel, len(kbs.FailedKeys))
	fmt.Println()
	fmt.Println(fmt.Sprintf("%s%*s", keyCol, length-len(keyCol), ""), notSyncedErrorLabel)
	fmt.Println(fmt.Sprintf("%s%*s", keyColSeparator, length-len(keyCol), ""), notSyncedErrorSeparator)
	for key, errorMessage := range kbs.FailedKeys {
		fmt.Println(fmt.Sprintf("%s%*s", key, length-len(key), ""), errorMessage)
	}
}

func printSyncedKBs(kbs *kbs.SyncResult) {
	fmt.Println()
	fmt.Println(syncedLabel)
	fmt.Println(titleSeparator)
	fmt.Println(totalSyncedLabel, len(kbs.NewIDs))
	fmt.Println()
	fmt.Println(fmt.Sprintf("%-36s", idCol), keyCol)
	fmt.Println(fmt.Sprintf("%-36s", idColSeparator), keyColSeparator)
	for key, id := range kbs.NewIDs {
		fmt.Println(id, key)
	}
}

func (i syncKBParams) dontPrint() bool {
	return !syncKBData.showAddedKBs && !syncKBData.showFailedKBs
}
