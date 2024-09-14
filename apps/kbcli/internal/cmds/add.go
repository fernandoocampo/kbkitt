package cmds

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/spf13/cobra"
)

// AddKBParams contains parameters required by add command to add a new KB.
type addKBParams struct {
	key       string
	value     string
	notes     string
	kind      string
	reference string
	tags      []string
}

// add messages
const (
	byeMessage            = "Bye!"
	areYouDoneLabel       = "> are you done? [y/n]: "
	kbToSaveLabel         = "...KB to save..."
	saveForLaterLabel     = "> do you want to save this KB to sync later? [y/n]: "
	saveQuestionLabel     = "> do you want to save it? [y/n]: "
	retryQuestionLabel    = "> do you want to retry? [y/n]: "
	kbAddedSuccessfully   = "kb added successfully"
	kbSavedForSyncSuccess = "kb successfully saved for later sync"
	missingTags           = "it seems tag values are missing, they will be useful to find this kb entry, Please indicate some of them..."
)

// field labels
const (
	keyLabel       = "key%s: "
	valueLabel     = "value%s: "
	notesLabel     = "notes%s: "
	kindLabel      = "class%s: "
	tagLabel       = "tag: "
	referenceLabel = "reference%s: "
	tagsLabel      = "tags (comma separated values): "
	showValueLabel = "(%s)"
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
	newCmd.PersistentFlags().StringVarP(&addKBData.reference, "reference", "r", "", "author or refence of this kb")
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

		fillMissingAddFields()

		ctx := context.Background()

		for {
			newKBToSave := addKBData.toNewKB()
			if !confirmKBData(&newKBToSave) {
				fmt.Println(byeMessage)
				os.Exit(0)
			}

			newKB, err := service.Add(ctx, newKBToSave)
			if errors.As(err, &kbs.DataError{}) {
				printAddingKBError(newKBToSave, err)
				if retry() {
					continue
				}

				fmt.Println(byeMessage)
				break
			}

			if err != nil {
				printAddingKBError(newKBToSave, err)
				if !saveForLater() {
					fmt.Println(byeMessage)
					break
				}

				errSaveForLater := service.SaveForSync(ctx, newKBToSave)
				if errSaveForLater != nil {
					fmt.Fprintln(os.Stderr, "unable to save new kb for sync:", errSaveForLater)
					fmt.Println()
					os.Exit(1)
				}

				fmt.Println(kbSavedForSyncSuccess)
				break
			}

			fmt.Println(kbAddedSuccessfully)
			fmt.Println(newKB)
			break
		}
		fmt.Println()
	}
}

func printAddingKBError(newKBToSave kbs.NewKB, err error) {
	fmt.Fprintln(os.Stderr, "unable to add new kb:", err)
	fmt.Println()
	fmt.Println(newKBToSave)
	fmt.Println()
}

func retry() bool {
	if wantToRetry() {
		fillExistingFields()
		return true
	}
	return false
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

func saveForLater() bool {
	if areYouSure(saveForLaterLabel) {
		fmt.Println()
		return true
	}
	fmt.Println()
	return false
}

func wantToRetry() bool {
	if areYouSure(retryQuestionLabel) {
		fmt.Println()
		return true
	}
	fmt.Println()
	return false
}

func fillMissingAddFields() {
	if kbs.IsStringEmpty(addKBData.key) {
		addKBData.key = requestStringValue(getLabel(keyLabel, addKBData.key))
	}
	if kbs.IsStringEmpty(addKBData.value) {
		addKBData.value = requestStringValue(getLabel(valueLabel, addKBData.value))
	}
	if kbs.IsStringEmpty(addKBData.notes) {
		addKBData.notes = requestStringValue(getLabel(notesLabel, addKBData.notes))
	}
	if kbs.IsStringEmpty(addKBData.kind) {
		addKBData.kind = requestStringValue(getLabel(kindLabel, addKBData.kind))
	}
	if kbs.IsStringEmpty(addKBData.reference) {
		addKBData.reference = requestStringValue(getLabel(referenceLabel, addKBData.reference))
	}
	if len(addKBData.tags) == 0 {
		fmt.Println()
		fmt.Println(missingTags)
		fmt.Println()
		addKBData.tags = readCSVFromStdin(tagLabel)
	}
}

func fillExistingFields() {
	addKBData.key = readStringValue(keyLabel, addKBData.key)
	addKBData.value = readStringValue(valueLabel, addKBData.value)
	addKBData.notes = readStringValue(notesLabel, addKBData.notes)
	addKBData.kind = readStringValue(kindLabel, addKBData.kind)
	addKBData.reference = readStringValue(referenceLabel, addKBData.reference)
	if len(addKBData.tags) == 0 {
		fmt.Println()
		fmt.Println(missingTags)
		fmt.Println()
		addKBData.tags = readCSVFromStdin(tagLabel)
	}
}

func readStringValue(label, currentValue string) string {
	value := requestStringValue(getLabel(label, currentValue))
	if kbs.IsStringEmpty(value) {
		return currentValue
	}

	return value
}

func getLabel(labelWithPattern string, value string) string {
	if kbs.IsStringEmpty(value) {
		return fmt.Sprintf(labelWithPattern, "")
	}

	return fmt.Sprintf(labelWithPattern, fmt.Sprintf(showValueLabel, value))
}

func (a addKBParams) toNewKB() kbs.NewKB {
	newKB := kbs.NewKB{
		Key:       a.key,
		Value:     a.value,
		Kind:      a.kind,
		Notes:     a.notes,
		Reference: a.reference,
		Tags:      make([]string, len(a.tags)),
	}

	copy(newKB.Tags, a.tags)

	return newKB
}
