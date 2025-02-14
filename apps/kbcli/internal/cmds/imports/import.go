package imports

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/filesystems"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

// importKBParams contains parameters required by add command to add a new KB.
type importKBParams struct {
	file          string
	showFailedKBs bool
	showAddedKBs  bool
}

// field labels
const (
	fileLabel                = "file path: "
	importedLabel            = "Imported KBs"
	unImportedLabel          = "Uimported KBs"
	unImportedErrorLabel     = "ERROR"
	unImportedErrorSeparator = "-----"
	totalImportedLabel       = "Total:"
	totalUnimportedLabel     = "Total:"
)

var importKBData importKBParams

func MakeImportCommand(service *kbs.Service) *cobra.Command {
	newCmd := cobra.Command{
		Use:   "import",
		Short: "import knowledge bases",
		Long:  `import knowledge bases from a file or other sources to your own kb repository`,
		Run:   makeRunImportKBCommand(service),
	}

	newCmd.PersistentFlags().StringVarP(&importKBData.file, "file", "f", "", "knowledge base key")
	newCmd.PersistentFlags().BoolVarP(&importKBData.showAddedKBs, "show-added-kbs", "", false, "knowledge base key")
	newCmd.PersistentFlags().BoolVarP(&importKBData.showFailedKBs, "show-failed-kbs", "", false, "knowledge base key")

	return &newCmd
}

func makeRunImportKBCommand(service *kbs.Service) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fillMissingImportFields()

		result, err := importFile(service)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to process import:", err)
			os.Exit(1)
		}

		printImportedReport(result)
	}
}

func importFile(service *kbs.Service) (*kbs.ImportResult, error) {
	newKBs, err := loadImportFile()
	if err != nil {
		return nil, fmt.Errorf("unable to import file: %w", err)
	}

	ctx := context.Background()
	result, err := service.Import(ctx, newKBs)
	if err != nil {
		return nil, fmt.Errorf("unable to import kbs: %w", err)
	}

	return result, nil
}

func loadImportFile() ([]kbs.NewKB, error) {
	file, err := filesystems.ReadFile(importKBData.file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file to import (%q): %w", importKBData.file, err)
	}

	dec := yaml.NewDecoder(bytes.NewReader(file))

	var kbItems []kbs.NewKB
	var kbItem kbs.NewKB
	for dec.Decode(&kbItem) == nil {
		kbItems = append(kbItems, kbItem)
		kbItem = kbs.NewKB{}
	}

	return kbItems, nil
}

func fillMissingImportFields() {
	if kbs.IsStringEmpty(importKBData.file) {
		importKBData.file = cmds.RequestStringValue(fileLabel)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func printImportedReport(kbs *kbs.ImportResult) {
	if kbs == nil {
		return
	}

	if importKBData.dontPrint() {
		return
	}

	if importKBData.showAddedKBs && len(kbs.NewIDs) > 0 {
		printImportedKBs(kbs)
	}

	if importKBData.showFailedKBs && len(kbs.FailedKeys) > 0 {
		printUnimportedKBs(kbs)
	}
}

func printUnimportedKBs(kbs *kbs.ImportResult) {
	length := len(cmds.KeyCol)
	for key := range kbs.FailedKeys {
		if len(key) > length {
			length = len(key)
		}
	}

	fmt.Println()
	fmt.Println(unImportedLabel)
	fmt.Println(cmds.TitleSeparator)
	fmt.Println(totalUnimportedLabel, len(kbs.FailedKeys))
	fmt.Println()
	fmt.Println(fmt.Sprintf("%s%*s", cmds.KeyCol, length-len(cmds.KeyCol), ""), unImportedErrorLabel)
	fmt.Println(fmt.Sprintf("%s%*s", cmds.KeyColSeparator, length-len(cmds.KeyCol), ""), unImportedErrorSeparator)
	for key, errorMessage := range kbs.FailedKeys {
		fmt.Println(fmt.Sprintf("%s%*s", key, length-len(key), ""), errorMessage)
	}
}

func printImportedKBs(kbs *kbs.ImportResult) {
	fmt.Println()
	fmt.Println(importedLabel)
	fmt.Println(cmds.TitleSeparator)
	fmt.Println(totalImportedLabel, len(kbs.NewIDs))
	fmt.Println()
	fmt.Println(fmt.Sprintf("%-36s", cmds.IDCol), cmds.KeyCol)
	fmt.Println(fmt.Sprintf("%-36s", cmds.IDColSeparator), cmds.KeyColSeparator)
	for key, id := range kbs.NewIDs {
		fmt.Println(id, key)
	}
}

func (i importKBParams) dontPrint() bool {
	return !importKBData.showAddedKBs && !importKBData.showFailedKBs
}
