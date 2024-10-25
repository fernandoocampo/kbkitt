package exports

import (
	"context"
	"fmt"
	"os"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/spf13/cobra"
)

// exportKBParams contains parameters required by export command.
type exportKBParams struct {
	namespace string
	category  string
}

// field labels
const (
	namespaceLabel     = "namespace: "
	categoryLabel      = "category: "
	totalExportedLabel = "Total:"
)

var exportKBData exportKBParams

func MakeExportCommand(service *kbs.Service) *cobra.Command {
	newCmd := cobra.Command{
		Use:   "export",
		Short: "get knowledge bases in yaml format",
		Long:  `get all knowledge bases in yaml format based on a given criteria`,
		Run:   makeRunExportedKBCommand(service),
	}

	newCmd.PersistentFlags().StringVarP(&exportKBData.namespace, "namespace", "n", "", "get all kbs with this namespace")
	newCmd.PersistentFlags().StringVarP(&exportKBData.category, "category", "c", "", "get all kbs with this category")

	return &newCmd
}

func makeRunExportedKBCommand(service *kbs.Service) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := exportData(service)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to process import:", err)
			os.Exit(1)
		}
	}
}

func exportData(service *kbs.Service) error {
	ctx := context.Background()

	total := 1 // hypotetical number
	filter := exportKBData.toGetAllKBFilter()

	for int(filter.Offset) <= total {
		result, err := service.GetAllKBs(ctx, filter)
		if err != nil {
			return fmt.Errorf("unable to export kbs: %w", err)
		}

		if result == nil {
			fmt.Println("no records were found")
			return nil
		}

		kbData, err := toYAMLDocuments(result.KBs, filter.Offset == 0)
		if err != nil {
			return fmt.Errorf("unable to export kbs: %w", err)
		}

		fmt.Println(string(kbData))

		total = result.Total
		filter.Offset += filter.Limit
	}

	printExportedKBs(total)

	return nil
}

func (e exportKBParams) toGetAllKBFilter() kbs.KBQueryFilter {
	return kbs.KBQueryFilter{
		Namespace: e.namespace,
		Category:  e.category,
		Limit:     20,
		Offset:    0,
	}
}

func printExportedKBs(total int) {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, totalExportedLabel, total)
	fmt.Fprintln(os.Stderr)
}

func toYAMLDocuments(kbs []kbs.KB, firstBlock bool) ([]byte, error) {
	var content []byte

	for _, kb := range kbs {
		newKBYAML, err := kb.ToYAML()
		if err != nil {
			return nil, fmt.Errorf("unable to save new kb for later sync: %w", err)
		}

		if firstBlock {
			firstBlock = false
		} else {
			content = append(content, []byte("---\n")...)
		}
		content = append(content, newKBYAML...)
	}

	return content, nil
}
