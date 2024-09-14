package kbs

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type KB struct {
	ID        string   `json:"id"`
	Key       string   `json:"key"`
	Value     string   `json:"value"`
	Notes     string   `json:"notes"`
	Kind      string   `json:"kind"`
	Reference string   `json:"reference,omitempty"`
	Tags      []string `json:"tags"`
}

type NewKB struct {
	Key       string   `json:"key" yaml:"Key"`
	Value     string   `json:"value" yaml:"Value"`
	Notes     string   `json:"notes" yaml:"Notes"`
	Kind      string   `json:"kind" yaml:"Kind"`
	Reference string   `json:"reference,omitempty" yaml:"Reference,omitempty"`
	Tags      []string `json:"tags" yaml:"Tags"`
}

type SearchResult struct {
	Items []KBItem `json:"items"`
	Total int64    `json:"total"`
	// determines the number of rows.
	Limit uint16 `json:"limit"`
	// skips the offset rows before beginning to return the rows.
	Offset uint16 `json:"offset"`
}

type ImportResult struct {
	// new kb keys and ids generated
	NewIDs map[string]string `json:"ids"`
	// failed kb keys with its respective error
	FailedKeys map[string]string `json:"failed_keys"`
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

// DataError defines an error to indicate that provided data was not valid.
type DataError struct {
	message string
}

// ServerError defines an error that was propagated by the server
type ServerError struct {
	inner   error
	message string
}

// ClientError defines an error that was propagated by the server
// due to some client request error
type ClientError struct {
	inner   error
	message string
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

func NewDataError(message string) DataError {
	return DataError{
		message: message,
	}
}

func NewServerError(message string) *ServerError {
	return &ServerError{
		message: message,
	}
}

func NewServerErrorWithWrapper(message string, err error) *ServerError {
	return &ServerError{
		inner:   err,
		message: message,
	}
}

func NewClientError(message string) *ClientError {
	return &ClientError{
		message: message,
	}
}

func NewClientErrorWithWrapper(message string, err error) *ClientError {
	return &ClientError{
		inner:   err,
		message: message,
	}
}

func (e DataError) Error() string {
	return e.message
}

func (e *ServerError) Error() string {
	if e.inner != nil {
		return fmt.Sprintf("%s: %s", e.message, e.inner.Error())
	}

	return e.message
}

func (e *ServerError) Unwrap() error {
	return e.inner
}

func (e ClientError) Error() string {
	if e.inner != nil {
		return fmt.Sprintf("%s: %s", e.message, e.inner.Error())
	}

	return e.message
}

func (e *ClientError) Unwrap() error {
	return e.inner
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

func (n NewKB) toYAML() ([]byte, error) {
	kbData, err := yaml.Marshal(n)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal new kb to yaml: %w", err)
	}

	return kbData, nil
}

func (k KBQueryFilter) validate() error {
	if IsStringEmpty(k.Key) && IsStringEmpty(k.Keyword) {
		return fmt.Errorf("invalid data to search kbs")
	}

	return nil
}

func (n NewKB) toKB(id string) KB {
	return KB{
		ID:        id,
		Key:       n.Key,
		Value:     n.Value,
		Notes:     n.Notes,
		Kind:      n.Kind,
		Reference: n.Reference,
		Tags:      slices.Clone(n.Tags),
	}
}

func (k KB) String() string {
	return fmt.Sprintf(`ID: %s
Key: %s
Value: %s
Notes: %s
Kind: %s
Reference: %s
Tags: %+v
`, k.ID, k.Key, k.Value, k.Notes, k.Kind, k.Reference, k.Tags)
}

func (n NewKB) String() string {
	return fmt.Sprintf(`Key: %s
Value: %s
Notes: %s
Kind: %s
Reference: %s
Tags: %+v
`, n.Key, n.Value, n.Notes, n.Kind, n.Reference, n.Tags)
}

func (i *ImportResult) Ok() bool {
	return len(i.FailedKeys) == 0 && len(i.NewIDs) > 0
}

func (i *ImportResult) anyError() bool {
	return len(i.FailedKeys) > 0
}

func IsStringEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}
