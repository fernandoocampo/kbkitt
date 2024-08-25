package kbs

import (
	"errors"
	"slices"
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

type KBItem struct {
	ID   string   `json:"id"`
	Key  string   `json:"key"`
	Kind string   `json:"kind"`
	Tags []string `json:"tags"`
}

type KBQueryFilter struct {
	Keyword string
	Key     string
	Limit   uint16
	Offset  uint16
}

var (
	errEmptyKBKey   = errors.New("kb key is empty")
	errEmptyKBValue = errors.New("kb value is empty")
	errEmptyKBKind  = errors.New("kb kind is empty")
	errEmptyKBTags  = errors.New("kb tags is empty")
)

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
