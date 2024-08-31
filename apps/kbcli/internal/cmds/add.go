package cmds

import (
	"context"
	"fmt"
	"os"

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
	areYouDoneLabel     = "are you done? [y/n]: "
	kbToSaveLabel       = "...KB to save..."
	saveQuestionLabel   = "do you want to save it? [y/n]: "
	kbAddedSuccessfully = "kb added successfully"
	missingTags         = "it seems tag values are missing, they will be useful to find this kb entry, Please indicate some of them..."
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
		service, err := newService()
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to load service: %s", err)
			fmt.Println()
			os.Exit(1)
		}

		fillMissingFields()

		ctx := context.Background()
		newKBToSave := addKBData.toNewKB()

		if !confirmKBData(&newKBToSave) {
			fmt.Println("bye")
			os.Exit(0)
		}

		newKB, err := service.Add(ctx, newKBToSave)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to add new kb:", err)
			fmt.Println()
			os.Exit(1)
		}

		fmt.Println(kbAddedSuccessfully)
		fmt.Println(newKB)
		fmt.Println()
	}
}

func confirmKBData(newKB *kbs.NewKB) bool {
	fmt.Println(kbToSaveLabel)
	fmt.Println()
	fmt.Println(newKB)
	fmt.Println()
	if areYouSure(saveQuestionLabel) {
		fmt.Println()
		return true
	}
	fmt.Println()
	return false
}

func fillMissingFields() {
	if kbs.IsStringEmpty(addKBData.key) {
		addKBData.key = requestStringValue(keyLabel)
	}
	if kbs.IsStringEmpty(addKBData.value) {
		addKBData.value = requestStringValue(valueLabel)
	}
	if kbs.IsStringEmpty(addKBData.notes) {
		addKBData.notes = requestStringValue(notesLabel)
	}
	if kbs.IsStringEmpty(addKBData.kind) {
		addKBData.kind = requestStringValue(kindLabel)
	}
	if len(addKBData.tags) == 0 {
		fmt.Println()
		fmt.Println(missingTags)
		fmt.Println()
		addKBData.tags = readCSVFromStdin(tagLabel)
	}
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
