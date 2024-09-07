package cmds

import (
	"context"
	"fmt"
	"os"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/spf13/cobra"
)

// getKBParams contains parameters required by get command
type getKBParams struct {
	id      string
	key     string
	keyword string
}

// field labels
const (
	totalLabel        = "Total:"
	limitLabel        = "Limit:"
	offsetLabel       = "Offset:"
	getKBIDLabel      = "id (hit <enter> if want to keep it empty): "
	getKBKeyLabel     = "key (hit <enter> if want to keep it empty): "
	getKBKeywordLabel = "keyword (hit <enter> if want to keep it empty): "
)

var getKBData getKBParams

func makeGetCommand() *cobra.Command {
	newCmd := cobra.Command{
		Use:   "get",
		Short: "get knowledge base content",
		Long:  `get a kb with id or key or other filter criteria based on tags`,
		Run:   makeGetKBCommand(),
	}

	newCmd.PersistentFlags().StringVarP(&getKBData.id, "id", "i", "", "knowledge base id")
	newCmd.PersistentFlags().StringVarP(&getKBData.key, "key", "k", "", "knowledge base key")
	newCmd.PersistentFlags().StringVarP(&getKBData.keyword, "keyword", "w", "", "knowledge base keyword to search based on tags")

	return &newCmd
}

func makeGetKBCommand() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fillFilterFields()

		ctx := context.Background()

		service, err := newService()
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to load service: %s", err)
			fmt.Println()
			os.Exit(1)
		}

		if !kbs.IsStringEmpty(getKBData.id) {
			kb, err := service.GetByID(ctx, getKBData.id)
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to get kb:", err)
				fmt.Println()
				os.Exit(1)
			}
			fmt.Println(kb)
			fmt.Println()
			return
		}

		if !kbs.IsStringEmpty(getKBData.key) || !kbs.IsStringEmpty(getKBData.keyword) {
			kbs, err := service.Search(ctx, kbs.NewKBQueryFilter(getKBData.key, getKBData.keyword))
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to get kb with given key:", err)
				fmt.Println()
				os.Exit(1)
			}
			printKBReport(kbs)
			fmt.Println()
		}
	}
}

func fillFilterFields() {
	if !kbs.IsStringEmpty(getKBData.id) {
		getKBData.key = ""
		getKBData.keyword = ""
		fmt.Println("using id to get kb")
		return
	}

	if kbs.IsStringEmpty(getKBData.id) {
		getKBData.id = requestStringValue(getKBIDLabel)
	}

	if !kbs.IsStringEmpty(getKBData.id) {
		return
	}

	if kbs.IsStringEmpty(getKBData.key) {
		getKBData.key = requestStringValue(getKBKeyLabel)
	}

	if !kbs.IsStringEmpty(getKBData.key) {
		return
	}

	if kbs.IsStringEmpty(getKBData.keyword) {
		getKBData.keyword = requestStringValue(getKBKeywordLabel)
	}

	if !kbs.IsStringEmpty(getKBData.keyword) {
		return
	}

	fillFilterFields()
}

func printKBReport(kbs *kbs.SearchResult) {
	length := len(keyCol)
	for _, v := range kbs.Items {
		if len(v.Key) > length {
			length = len(v.Key)
		}
	}
	fmt.Println()
	fmt.Println(totalLabel, kbs.Total)
	fmt.Println(limitLabel, kbs.Total)
	fmt.Println(offsetLabel, kbs.Total)
	fmt.Println()
	fmt.Println(fmt.Sprintf("%-36s", "ID"), fmt.Sprintf("%s%*s", "KEY", length-len(keyCol), ""), "KIND")
	for _, kb := range kbs.Items {
		fmt.Println(kb.ID, fmt.Sprintf("%s%*s", kb.Key, length-len(kb.Key), ""), kb.Kind)
	}
}
