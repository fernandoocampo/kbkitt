package updates

import (
	"context"
	"errors"
	"fmt"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/spf13/cobra"
)

type updateKB struct {
	service     *kbs.Service
	kb          *kbs.KB
	id          string
	interactive bool
}

const (
	doYouWantToUpdateLabel = "Are you sure you want to update this knowledge base? [y/n]: "
	kbToUpdateLabel        = "...KB to update..."
	updateQuestionLabel    = "> do you want to update it? [y/n]: "
)

// field labels
const (
	keyLabel       = "key%s: "
	valueLabel     = "value%s: "
	notesLabel     = "notes%s: "
	kindLabel      = "class%s: "
	tagLabel       = "tag: "
	referenceLabel = "reference%s: "
	tagsLabel      = "tags (space separated values): "
	showValueLabel = "(%s)"
)

var updateKBData *updateKB

func MakeUpdateCommand(service *kbs.Service) *cobra.Command {
	newCmd := cobra.Command{
		Use:   "update",
		Short: "update an existing kb",
		Long:  "update a kb with a given id in classic or interactive mode",
		Run:   makeRunUpdateCommand(),
	}

	updateKBData = &updateKB{
		service: service,
	}

	newCmd.PersistentFlags().StringVarP(&updateKBData.id, "id", "i", "", "knowledge base id")
	newCmd.PersistentFlags().BoolVarP(&updateKBData.interactive, "ux", "u", false, "show result in interactive mode")

	return &newCmd
}

func makeRunUpdateCommand() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fillMissingUpdateFields()

		ctx := context.Background()

		showKBToUpdate(ctx)
	}
}

func showKBToUpdate(ctx context.Context) error {
	kbToUpdate, err := updateKBData.service.GetByID(ctx, updateKBData.id)
	if err != nil {
		return fmt.Errorf("unable to show kb to update: %w", err)
	}

	if kbToUpdate == nil {
		return errors.New("kb with given id does not exist")
	}

	updateKBData.kb = kbToUpdate

	fmt.Println()
	fmt.Println(updateKBData.kb)
	fmt.Println()

	if !cmds.AreYouSure(doYouWantToUpdateLabel) {
		fmt.Println("bye")
		return nil
	}

	fmt.Println()

	err = runInteractive()
	if err != nil {
		return fmt.Errorf("unable to show form: %w", err)
	}

	if !confirmKBData() {
		fmt.Println("bye")
		return nil
	}

	err = updateKBData.service.Update(ctx, getKBToUpdate())
	if err != nil {
		return fmt.Errorf("unable to update kb: %w", err)
	}

	return nil
}

func fillMissingUpdateFields() {
	if !kbs.IsStringEmpty(updateKBData.id) {
		return
	}

	updateKBData.id = cmds.RequestStringValue(cmds.GetKBIDLabel)

	fillMissingUpdateFields()
}

func confirmKBData() bool {
	fmt.Println(kbToUpdateLabel)
	fmt.Println()
	fmt.Println(updateKBData.kb)
	fmt.Println()
	if cmds.AreYouSure(updateQuestionLabel) {
		fmt.Println()
		return true
	}
	fmt.Println()
	return false
}

func getKBToUpdate() kbs.KB {
	return *updateKBData.kb
}
