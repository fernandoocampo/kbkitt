package gets

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds"
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

func MakeGetCommand(service *kbs.Service) *cobra.Command {
	newCmd := cobra.Command{
		Use:   "get",
		Short: "get knowledge base content",
		Long:  `get a kb with id or key or other filter criteria based on tags`,
		Run:   makeGetKBCommand(service),
	}

	newCmd.PersistentFlags().StringVarP(&getKBData.id, "id", "i", "", "knowledge base id")
	newCmd.PersistentFlags().StringVarP(&getKBData.key, "key", "k", "", "knowledge base key")
	newCmd.PersistentFlags().StringVarP(&getKBData.keyword, "keyword", "w", "", "knowledge base keyword to search based on tags")

	return &newCmd
}

func makeGetKBCommand(service *kbs.Service) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fillFilterFields()

		ctx := context.Background()

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
		getKBData.id = cmds.RequestStringValue(getKBIDLabel)
	}

	if !kbs.IsStringEmpty(getKBData.id) {
		return
	}

	if kbs.IsStringEmpty(getKBData.key) {
		getKBData.key = cmds.RequestStringValue(getKBKeyLabel)
	}

	if !kbs.IsStringEmpty(getKBData.key) {
		return
	}

	if kbs.IsStringEmpty(getKBData.keyword) {
		getKBData.keyword = cmds.RequestStringValue(getKBKeywordLabel)
	}

	if !kbs.IsStringEmpty(getKBData.keyword) {
		return
	}

	fillFilterFields()
}

func printKBReport(kbs *kbs.SearchResult) {
	keyLength := len(cmds.KeyCol)
	for key := range kbs.Keys() {
		if len(key) > keyLength {
			keyLength = len(key)
		}
	}
	kindLength := len(cmds.KeyCol)
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
	fmt.Println(fmt.Sprintf("%-36s", cmds.IDCol), fmt.Sprintf("%s%*s", cmds.KeyCol, keyLength-len(cmds.KeyCol), ""), fmt.Sprintf("%s%*s", cmds.KindCol, kindLength-len(cmds.KindCol), ""), cmds.TagCol)
	fmt.Println(fmt.Sprintf("%-36s", cmds.IDColSeparator), fmt.Sprintf("%s%*s", cmds.KeyColSeparator, keyLength-len(cmds.KeyCol), ""), fmt.Sprintf("%s%*s", cmds.KindColSeparator, kindLength-len(cmds.KindCol), ""), cmds.TagColSeparator)
	for _, kb := range kbs.Items {
		fmt.Println(kb.ID, fmt.Sprintf("%s%*s", kb.Key, keyLength-len(kb.Key), ""), fmt.Sprintf("%s%*s", kb.Kind, kindLength-len(kb.Kind), ""), strings.Join(kb.Tags, ","))
	}
}
