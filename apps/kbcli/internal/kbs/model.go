package kbs

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type KB struct {
	ID    string   `json:"id"`
	Key   string   `json:"key"`
	Value string   `json:"value"`
	Notes string   `json:"notes"`
	Kind  string   `json:"kind"`
	Tags  []string `json:"tags"`
}

type NewKB struct {
	Key   string   `json:"key"`
	Value string   `json:"value"`
	Notes string   `json:"notes"`
	Kind  string   `json:"kind"`
	Tags  []string `json:"tags"`
}

type SearchResult struct {
	Items []KBItem `json:"items"`
	Total int64    `json:"total"`
	// determines the number of rows.
	Limit uint16 `json:"limit"`
	// skips the offset rows before beginning to return the rows.
	Offset uint16 `json:"offset"`
}

type KBItem struct {
	ID   string   `json:"id"`
	Key  string   `json:"key"`
	Kind string   `json:"kind"`
	Tags []string `json:"tags"`
}

type KBQueryFilter struct {
	Keyword string `json:"keyword"`
	Key     string `json:"key"`
	// determines the number of rows.
	Limit uint16 `json:"limit"`
	// skips the offset rows before beginning to return the rows.
	Offset uint16 `json:"offset"`
}

var (
	errEmptyKBKey   = errors.New("kb key is empty")
	errEmptyKBValue = errors.New("kb value is empty")
	errEmptyKBKind  = errors.New("kb kind is empty")
	errEmptyKBTags  = errors.New("kb tags is empty")
)

func NewKBQueryFilter(key, keyword string) KBQueryFilter {
	return KBQueryFilter{
		Key:     key,
		Keyword: keyword,
		Limit:   5,
		Offset:  0,
	}
}

func (n NewKB) validate() error {
	var err error

	if n.Key == "" {
		err = errors.Join(err, errEmptyKBKey)
	}

	if n.Value == "" {
		err = errors.Join(err, errEmptyKBValue)
	}

	if n.Kind == "" {
		err = errors.Join(err, errEmptyKBKind)
	}

	if len(n.Tags) == 0 {
		err = errors.Join(err, errEmptyKBTags)
	}

	return err
}

func (k KBQueryFilter) validate() error {
	if IsStringEmpty(k.Key) && IsStringEmpty(k.Keyword) {
		return fmt.Errorf("invalid data to search kbs")
	}

	return nil
}

func (n NewKB) toKB(id string) KB {
	return KB{
		ID:    id,
		Key:   n.Key,
		Value: n.Value,
		Notes: n.Notes,
		Kind:  n.Kind,
		Tags:  slices.Clone(n.Tags),
	}
}

func (k KB) String() string {
	return fmt.Sprintf(`ID: %s
Key: %s
Value: %s
Notes: %s
Kind: %s
Tags: %+v
`, k.ID, k.Key, k.Value, k.Notes, k.Kind, k.Tags)
}

func (n NewKB) String() string {
	return fmt.Sprintf(`Key: %s
Value: %s
Notes: %s
Kind: %s
Tags: %+v
`, n.Key, n.Value, n.Notes, n.Kind, n.Tags)
}

func IsStringEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}
