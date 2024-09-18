package cmds

import (
	"context"
	"fmt"
	"os"
	"strings"

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
			fmt.Println()
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
	keyLength := len(keyCol)
	for key := range kbs.Keys() {
		if len(key) > keyLength {
			keyLength = len(key)
		}
	}
	kindLength := len(keyCol)
	for kind := range kbs.Kinds() {
		if len(kind) > kindLength {
			kindLength = len(kind)
		}
	}
	fmt.Println()
	fmt.Println(totalLabel, kbs.Total)
	fmt.Println(limitLabel, kbs.Total)
	fmt.Println(offsetLabel, kbs.Total)
	fmt.Println()
	fmt.Println(fmt.Sprintf("%-36s", idCol), fmt.Sprintf("%s%*s", keyCol, keyLength-len(keyCol), ""), fmt.Sprintf("%s%*s", kindCol, kindLength-len(kindCol), ""), tagCol)
	fmt.Println(fmt.Sprintf("%-36s", idColSeparator), fmt.Sprintf("%s%*s", keyColSeparator, keyLength-len(keyCol), ""), fmt.Sprintf("%s%*s", kindColSeparator, kindLength-len(kindCol), ""), tagColSeparator)
	for _, kb := range kbs.Items {
		fmt.Println(kb.ID, fmt.Sprintf("%s%*s", kb.Key, keyLength-len(kb.Key), ""), fmt.Sprintf("%s%*s", kb.Kind, kindLength-len(kb.Kind), ""), strings.Join(kb.Tags, ","))
	}
}
