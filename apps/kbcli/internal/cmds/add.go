package cmds

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/spf13/cobra"
)

// AddKBParams contains parameters required by add command to add a new KB.
type addKBParams struct {
	key   string
	value string
	notes string
	kind  string
	tags  []string
}

// add messages
const (
	areYouDoneLabel = "are you done? [y/n]: "
	missingTags     = "it seems tag values are missing, they will be useful to find this kb entry, Please indicate some..."
)

// field labels
const (
	keyLabel   = "key: "
	valueLabel = "value: "
	notesLabel = "notes: "
	kindLabel  = "class: "
	tagLabel   = "tag: "
	tagsLabel  = "tags (comma separated values): "
)

// values
const (
	yes      = "yes"
	yesShort = "y"
)

var addKBData addKBParams

func makeAddCommand() *cobra.Command {
	newCmd := cobra.Command{
		Use:   "add",
		Short: "add a new knowledge base",
		Long:  `add a new knowledge base such as: concepts, commands, prompts, etc.`,
		Run:   makeRunAddKBCommand(),
	}

	newCmd.PersistentFlags().StringVarP(&addKBData.key, "key", "k", "", "knowledge base key")
	newCmd.PersistentFlags().StringVarP(&addKBData.value, "value", "v", "", "knowledge base value")
	newCmd.PersistentFlags().StringVarP(&addKBData.notes, "notes", "n", "", "knowledge base notes")
	newCmd.PersistentFlags().StringVarP(&addKBData.kind, "class", "c", "", "kind of knowledge base")
	newCmd.PersistentFlags().StringSliceVarP(&addKBData.tags, "tags", "t", []string{}, "comma separated tags for this kb")

	return &newCmd
}

func makeRunAddKBCommand() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fillMissingFields()

		newKB := addKBData.toNewKB()

		fmt.Printf("%+v", newKB)
	}
}

func fillMissingFields() {
	if isStringEmpty(addKBData.key) {
		addKBData.key = requestStringValue(keyLabel)
	}
	if isStringEmpty(addKBData.value) {
		addKBData.value = requestStringValue(valueLabel)
	}
	if isStringEmpty(addKBData.notes) {
		addKBData.notes = requestStringValue(notesLabel)
	}
	if isStringEmpty(addKBData.kind) {
		addKBData.kind = requestStringValue(kindLabel)
	}
	if len(addKBData.tags) == 0 {
		fmt.Println(missingTags)
		fmt.Println()
		addKBData.tags = readCSVFromStdin(tagLabel)
	}
}

func requestStringValue(label string) string {
	var output string
	fmt.Print(label)

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		output = scanner.Text()
	}

	return output
}

func readCSVFromStdin(label string) []string {
	var result []string
	for {
		var value string
		fmt.Print(label)
		fmt.Scan(&value)

		result = append(result, value)

		var done string
		fmt.Print(areYouDoneLabel)
		fmt.Scan(&done)

		if strings.EqualFold(done, yes) || strings.EqualFold(done, yesShort) {
			fmt.Println()
			break
		}
	}
	return result
}

func isStringEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func (a addKBParams) toNewKB() kbs.NewKB {
	newKB := kbs.NewKB{
		Key:   a.key,
		Value: a.value,
		Kind:  a.kind,
		Notes: a.notes,
		Tags:  make([]string, len(a.tags)),
	}

	copy(newKB.Tags, a.tags)

	return newKB
}
