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
	id          string
	key         string
	category    string
	namespace   string
	keyword     string
	limit       uint32
	offset      uint32
	interactive bool
}

// field labels
const (
	totalLabel          = "Total:"
	limitLabel          = "Limit:"
	offsetLabel         = "Offset:"
	getKBIDLabel        = "id (hit <enter> if want to keep it empty): "
	getKBKeyLabel       = "key (hit <enter> if want to keep it empty): "
	getKBCategoryLabel  = "category (hit <enter> if want to keep it empty): "
	getKBNamespaceLabel = "namespace (hit <enter> if want to keep it empty): "
	getKBKeywordLabel   = "keyword (hit <enter> if want to keep it empty): "
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
	newCmd.PersistentFlags().StringVarP(&getKBData.category, "category", "c", "", "knowledge base category. e.g bookmark, quote, etc")
	newCmd.PersistentFlags().StringVarP(&getKBData.namespace, "namespace", "n", "", "knowledge base namespace")
	newCmd.PersistentFlags().StringVarP(&getKBData.keyword, "keyword", "w", "", "knowledge base keyword to search based on tags")
	newCmd.PersistentFlags().Uint32VarP(&getKBData.limit, "limit", "l", 5, "number of rows you want to retrieve")
	newCmd.PersistentFlags().Uint32VarP(&getKBData.offset, "offset", "o", 0, "number of rows to skip before starting to return result rows")
	newCmd.PersistentFlags().BoolVarP(&getKBData.interactive, "ux", "u", false, "show result in interactive mode")

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

		if getKBData.nothingToLookFor() {
			os.Exit(0)
		}

		err := search(ctx, service)
		if err != nil {
			fmt.Fprintln(os.Stderr, "searching:", err)
			fmt.Println()
			os.Exit(1)
		}
		fmt.Println()
	}
}

func search(ctx context.Context, service *kbs.Service) error {
	if getKBData.interactive {
		err := runInteractive(ctx, service)
		if err != nil {
			return fmt.Errorf("unable to run interactive search: %w", err)
		}
		return nil
	}

	err := searchBasic(ctx, service)
	if err != nil {
		return fmt.Errorf("unable to run basic search: %w", err)
	}

	return nil
}

func searchBasic(ctx context.Context, service *kbs.Service) error {
	result, err := service.Search(ctx, getKBData.toKBQueryFilter())
	if err != nil {
		return fmt.Errorf("unable to search: %w", err)
	}

	printSimpleKBReport(result)

	return nil
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

	if kbs.IsStringEmpty(getKBData.category) {
		getKBData.category = cmds.RequestStringValue(getKBCategoryLabel)
	}

	if kbs.IsStringEmpty(getKBData.namespace) {
		getKBData.namespace = cmds.RequestStringValue(getKBNamespaceLabel)
	}

	if kbs.IsStringEmpty(getKBData.keyword) {
		getKBData.keyword = cmds.RequestStringValue(getKBKeywordLabel)
	}

	if getKBData.somethingToLookFor() {
		return
	}

	fillFilterFields()
}

func printSimpleKBReport(result *kbs.SearchResult) {
	keyLength := kbs.GetLongerText(cmds.KeyCol, result.Keys())
	categoryLength := kbs.GetLongerText(cmds.CategoryCol, result.Categories())
	namespaceLength := kbs.GetLongerText(cmds.NamespaceCol, result.Namespaces())
	fmt.Println()
	fmt.Println(totalLabel, result.Total)
	fmt.Println(limitLabel, result.Limit)
	fmt.Println(offsetLabel, result.Offset)
	fmt.Println()
	printReportTitles(keyLength, categoryLength, namespaceLength)
	for _, kb := range result.Items {
		fmt.Println(kb.ID, fmt.Sprintf("%s%*s", kb.Key, keyLength-len(kb.Key), ""), fmt.Sprintf("%s%*s", kb.Category, categoryLength-len(kb.Category), ""), strings.Join(kb.Tags, ","))
	}
}

func printReportTitles(keyLength, categoryLength, namespaceLength int) {
	fmt.Println(
		fmt.Sprintf("%-36s", cmds.IDCol),
		fmt.Sprintf("%s%*s", cmds.KeyCol, keyLength-len(cmds.KeyCol), ""),
		fmt.Sprintf("%s%*s", cmds.CategoryCol, categoryLength-len(cmds.CategoryCol), ""),
		fmt.Sprintf("%s%*s", cmds.NamespaceCol, namespaceLength-len(cmds.NamespaceCol), ""),
		cmds.TagCol,
	)

	fmt.Println(
		fmt.Sprintf("%-36s", cmds.IDColSeparator),
		fmt.Sprintf("%s%*s", cmds.KeyColSeparator, keyLength-len(cmds.KeyCol), ""),
		fmt.Sprintf("%s%*s", cmds.CategoryColSeparator, categoryLength-len(cmds.CategoryCol), ""),
		fmt.Sprintf("%s%*s", cmds.NamespaceColSeparator, namespaceLength-len(cmds.NamespaceCol), ""),
		cmds.TagColSeparator)
}

func (g *getKBParams) toKBQueryFilter() kbs.KBQueryFilter {
	return kbs.KBQueryFilter{
		Keyword:   getKBData.keyword,
		Key:       getKBData.key,
		Category:  getKBData.category,
		Namespace: getKBData.namespace,
		Limit:     getKBData.limit,
		Offset:    getKBData.offset,
	}
}

func (g *getKBParams) nothingToLookFor() bool {
	return kbs.IsStringEmpty(getKBData.key) &&
		kbs.IsStringEmpty(getKBData.keyword) &&
		kbs.IsStringEmpty(getKBData.category) &&
		kbs.IsStringEmpty(getKBData.namespace)
}

func (g *getKBParams) somethingToLookFor() bool {
	return !kbs.IsStringEmpty(getKBData.keyword) ||
		!kbs.IsStringEmpty(getKBData.category) ||
		!kbs.IsStringEmpty(getKBData.namespace)
}
