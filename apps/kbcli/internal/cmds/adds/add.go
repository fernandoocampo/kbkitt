package adds

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/spf13/cobra"
)

// AddKBParams contains parameters required by add command to add a new KB.
type addKBParams struct {
	key         string
	value       string
	notes       string
	category    string
	namespace   string
	reference   string
	mediaType   string
	rawTags     string
	interactive bool
	tags        []string
}

// add messages
const (
	byeMessage             = "Bye!"
	kbToSaveLabel          = "...KB to save..."
	saveForLaterLabel      = "> do you want to save this KB to sync later? [y/n]: "
	saveQuestionLabel      = "> do you want to save it? [y/n]: "
	retryQuestionLabel     = "> do you want to retry? [y/n]: "
	kbAddedSuccessfully    = "kb added successfully"
	mediaSavedSuccessfully = "kb media added successfully"
	kbSavedForSyncSuccess  = "kb successfully saved for later sync"
	missingMediaType       = "it seems that the type of media cannot be determined..."
	saveMediaQuestionLabel = "> do you want to save this media kb locally on your computer? [y/n]: "
)

// field labels
const (
	keyLabel       = "key%s: "
	valueLabel     = "value%s: "
	notesLabel     = "notes%s: "
	categoryLabel  = "category%s: "
	namespaceLabel = "namespace%s: "
	tagLabel       = "tags%s: "
	referenceLabel = "reference%s: "
	mediaTypeLabel = "media type%s: "
	tagsLabel      = "tags (space separated values): "
	showValueLabel = "(%s)"
)

var addKBData addKBParams
var exitGUI bool

func MakeAddCommand(service *kbs.Service) *cobra.Command {
	newCmd := cobra.Command{
		Use:   "add",
		Short: "add a new knowledge base",
		Long:  `add a new knowledge base such as: concepts, commands, prompts, etc.`,
		Run:   makeRunAddKBCommand(service),
	}

	newCmd.PersistentFlags().StringVarP(&addKBData.key, "key", "k", "", "knowledge base key")
	newCmd.PersistentFlags().StringVarP(&addKBData.value, "value", "v", "", "knowledge base value")
	newCmd.PersistentFlags().StringVarP(&addKBData.notes, "notes", "o", "", "knowledge base notes")
	newCmd.PersistentFlags().StringVarP(&addKBData.category, "category", "c", "", "category of knowledge base")
	newCmd.PersistentFlags().StringVarP(&addKBData.namespace, "namespace", "n", "default", "namespace of knowledge base")
	newCmd.PersistentFlags().StringVarP(&addKBData.reference, "reference", "r", "", "author or refence of this kb")
	newCmd.PersistentFlags().StringSliceVarP(&addKBData.tags, "tags", "t", []string{}, "comma separated tags for this kb")
	newCmd.PersistentFlags().BoolVarP(&addKBData.interactive, "ux", "u", false, "add KB in interactive mode")

	return &newCmd
}

func makeRunAddKBCommand(service *kbs.Service) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := collectData()
		if err != nil {
			fmt.Fprintln(os.Stderr, "collecting data", err)
			os.Exit(1)
		}

		if exitGUI {
			os.Exit(0)
		}

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
		if saveMedia() {
			fmt.Println()
			err := service.SaveMedia(ctx, addKBData.toNewKB())
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to save media locally:", err)
				fmt.Println()
				os.Exit(1)
			}
			fmt.Println(mediaSavedSuccessfully)
		}
		fmt.Println()
	}
}

func collectData() error {
	if addKBData.interactive {
		err := runInteractive()
		if err != nil {
			return fmt.Errorf("unable to collect parameters: %w", err)
		}

		checkMediaType()

		return nil
	}

	fillMissingAddFields()

	return nil
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
	if cmds.AreYouSure(saveQuestionLabel) {
		fmt.Println()
		return true
	}
	fmt.Println()
	return false
}

func saveForLater() bool {
	if cmds.AreYouSure(saveForLaterLabel) {
		fmt.Println()
		return true
	}
	fmt.Println()
	return false
}

func saveMedia() bool {
	if !isMediaType() {
		return false
	}

	fmt.Println()
	return cmds.AreYouSure(saveMediaQuestionLabel)
}

func wantToRetry() bool {
	if cmds.AreYouSure(retryQuestionLabel) {
		fmt.Println()
		return true
	}
	fmt.Println()
	return false
}

func fillMissingAddFields() {
	if kbs.IsStringEmpty(addKBData.key) {
		addKBData.key = cmds.RequestStringValue(getLabel(keyLabel, addKBData.key))
	}
	if kbs.IsStringEmpty(addKBData.value) {
		addKBData.value = cmds.RequestStringValue(getLabel(valueLabel, addKBData.value))
	}
	if kbs.IsStringEmpty(addKBData.notes) {
		addKBData.notes = cmds.RequestStringValue(getLabel(notesLabel, addKBData.notes))
	}
	if kbs.IsStringEmpty(addKBData.category) {
		addKBData.category = cmds.RequestStringValue(getLabel(categoryLabel, addKBData.category))
	}
	if kbs.IsStringEmpty(addKBData.namespace) {
		addKBData.namespace = strings.ToLower(cmds.RequestStringValue(getLabel(namespaceLabel, addKBData.namespace)))
	}
	if kbs.IsStringEmpty(addKBData.reference) {
		addKBData.reference = cmds.RequestStringValue(getLabel(referenceLabel, addKBData.reference))
	}
	if kbs.IsStringEmpty(addKBData.rawTags) && len(addKBData.tags) == 0 {
		addKBData.rawTags = cmds.RequestStringValue(getLabel(tagLabel, addKBData.rawTags))
		addKBData.buildTags()
	}

	checkMediaType()
}

func checkMediaType() {
	if isMediaTypeEmpty() {
		addKBData.mediaType = requestMediaTypeValue()
	}
}

func isMediaTypeEmpty() bool {
	if !isMediaType() {
		return false
	}

	if kbs.IsStringEmpty(addKBData.mediaType) {
		return true
	}

	return false
}

func isMediaType() bool {
	return strings.EqualFold(addKBData.category, kbs.MediaCategory)
}

func requestMediaTypeValue() string {
	index := strings.LastIndex(addKBData.value, ".")
	if index > 0 &&
		len(addKBData.value[index:]) > 1 &&
		slices.Contains(kbs.MediaExtensions, addKBData.value[index+1:]) {
		return addKBData.value[index+1:]
	}

	fmt.Println()
	fmt.Println(missingMediaType)
	fmt.Println()

	return cmds.RequestStringValue(getLabel(mediaTypeLabel, addKBData.mediaType))
}

func fillExistingFields() {
	addKBData.key = readStringValue(keyLabel, addKBData.key)
	addKBData.value = readStringValue(valueLabel, addKBData.value)
	addKBData.notes = readStringValue(notesLabel, addKBData.notes)
	addKBData.category = readStringValue(categoryLabel, addKBData.category)
	addKBData.namespace = readStringValue(namespaceLabel, addKBData.namespace)
	addKBData.reference = readStringValue(referenceLabel, addKBData.reference)
	addKBData.rawTags = readStringValue(tagLabel, addKBData.rawTags)
}

func readStringValue(label, currentValue string) string {
	value := cmds.RequestStringValue(getLabel(label, currentValue))
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
		Category:  a.category,
		Namespace: a.namespace,
		Notes:     a.notes,
		Reference: a.reference,
		MediaType: a.mediaType,
		Tags:      make([]string, len(a.tags)),
	}

	copy(newKB.Tags, a.tags)

	return newKB
}

func (a *addKBParams) buildTags() {
	a.tags = strings.Split(a.rawTags, " ")
}
