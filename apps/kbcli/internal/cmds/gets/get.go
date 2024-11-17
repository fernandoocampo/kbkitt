package gets

import (
	"context"
	"fmt"
	"os"

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
		ctx := context.Background()

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
	err := runInteractive(ctx, service)
	if err != nil {
		return fmt.Errorf("unable to run interactive search: %w", err)
	}
	return nil
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
