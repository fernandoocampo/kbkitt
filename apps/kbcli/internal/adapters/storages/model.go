package storages

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
)

type kb struct {
	KeyID       string
	InternalID  string
	Key         string
	Value       string
	Notes       string
	Category    string
	Reference   string
	Tags        string
	DateCreated time.Time
}

type kbItem struct {
	ID   string
	Key  string
	Kind string
	Tags string
}

type filterBuilder struct {
	query          string
	countStatement string
	filters        []string
	queryArgs      []interface{}
	countArgs      []interface{}
}

const (
	aSpace = " "

	equalsOperator = "="
	whereOperator  = "WHERE"
	andOperator    = "AND"
	likeOperator   = "LIKE"
	matchOperator  = "MATCH"
	limitOperator  = "LIMIT"
	offsetOperator = "OFFSET"
)

// user columns.
const (
	keyColumn              = "k.KB_KEY"
	tagValuesVirtualColumn = "t.tag_values"
	rowIDVirtualColum      = "t.rowid"
	internalIDColumn       = "k.INTERNAL_ID"
)

// errors
var (
	errUnableToSearchKBS = errors.New("unable to search kbs")
)

func (k kb) toKB() *kbs.KB {
	newKB := kbs.KB{
		ID:        k.KeyID,
		Key:       k.Key,
		Value:     k.Value,
		Notes:     k.Notes,
		Kind:      k.Category,
		Reference: k.Reference,
		Tags:      strings.Split(k.Tags, aSpace),
	}

	return &newKB
}

func (f *filterBuilder) addCondition(field, operator string, value interface{}) *filterBuilder {
	isHint := false
	condition := whereOperator

	if len(f.filters) > 0 {
		condition = " " + andOperator
	}

	newStatement := fmt.Sprintf("%s %s %s", condition, field, operator)

	return f.addFilter(newStatement, value, isHint)
}

func (f *filterBuilder) addFilter(statement string, value interface{}, isHint bool) *filterBuilder {
	index := len(f.filters) + 1

	statement = fmt.Sprintf("%s $%d", statement, index)

	f.filters = append(f.filters, statement)

	if !isHint {
		f.countArgs = append(f.countArgs, value)
	}

	f.queryArgs = append(f.queryArgs, value)

	return f
}

func toDBKB(akb *kbs.KB) kb {
	return kb{
		KeyID:       akb.ID,
		Key:         akb.Key,
		Value:       akb.Value,
		Notes:       akb.Notes,
		Category:    akb.Kind,
		Reference:   akb.Reference,
		Tags:        strings.Join(akb.Tags, aSpace),
		DateCreated: time.Now().UTC(),
	}
}

func (k kbItem) toKBItem() kbs.KBItem {
	return kbs.KBItem{
		ID:   k.ID,
		Key:  k.Key,
		Kind: k.Kind,
		Tags: strings.Split(k.Tags, aSpace),
	}
}

func toKBItems(dbKBItems []kbItem) []kbs.KBItem {
	result := make([]kbs.KBItem, 0, len(dbKBItems))

	for _, v := range dbKBItems {
		result = append(result, v.toKBItem())
	}

	return result
}
