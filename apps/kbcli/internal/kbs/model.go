package kbs

import (
	"bytes"
	"errors"
	"fmt"
	"iter"
	"math"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/filesystems"
	"github.com/google/uuid"
	yaml "gopkg.in/yaml.v3"
)

type KB struct {
	ID        string   `json:"id" yaml:"-"`
	Key       string   `json:"key" yaml:"Key"`
	Value     string   `json:"value" yaml:"Value"`
	Notes     string   `json:"notes" yaml:"Notes"`
	Category  string   `json:"category" yaml:"Category"`
	Reference string   `json:"reference,omitempty" yaml:"Reference"`
	Namespace string   `json:"namespace,omitempty" yaml:"Namespace"`
	Tags      []string `json:"tags" yaml:"Tags"`
}

type NewKB struct {
	Key       string   `json:"key" yaml:"Key"`
	Value     string   `json:"value" yaml:"Value"`
	Notes     string   `json:"notes" yaml:"Notes"`
	Category  string   `json:"category" yaml:"Category"`
	Reference string   `json:"reference,omitempty" yaml:"Reference,omitempty"`
	MediaType string   `json:"media_type,omitempty" yaml:"MediaType,omitempty"`
	Namespace string   `json:"namespace,omitempty" yaml:"Namespace"`
	Tags      []string `json:"tags" yaml:"Tags"`
}

type SearchResult struct {
	Items []KBItem `json:"items"`
	Total int      `json:"total"`
	// determines the number of rows.
	Limit uint32 `json:"limit"`
	// skips the offset rows before beginning to return the rows.
	Offset uint32 `json:"offset"`
}

type GetAllResult struct {
	KBs   []KB `json:"kbs"`
	Total int  `json:"total"`
	// determines the number of rows.
	Limit uint32 `json:"limit"`
	// skips the offset rows before beginning to return the rows.
	Offset uint32 `json:"offset"`
}

type ImportResult struct {
	// new kb keys and ids generated
	NewIDs map[string]string `json:"ids"`
	// failed kb keys with its respective error
	FailedKeys map[string]string `json:"failed_keys"`
}

type SyncResult struct {
	// new kb keys and ids generated
	NewIDs map[string]string `json:"ids"`
	// failed kb keys with its respective error
	FailedKeys map[string]string `json:"failed_keys"`
}

type KBItem struct {
	ID        string   `json:"id"`
	Key       string   `json:"key"`
	Category  string   `json:"category"`
	Namespace string   `json:"namespace,omitempty"`
	Tags      []string `json:"tags"`
}

type KBQueryFilter struct {
	Keyword   string `json:"keyword"`
	Key       string `json:"key"`
	Category  string `json:"category"`
	Namespace string `json:"namespace"`
	// determines the number of rows.
	Limit uint32 `json:"limit"`
	// skips the offset rows before beginning to return the rows.
	Offset uint32 `json:"offset"`
}

type MediaFileData struct {
	IsMediaFile bool
	Path        string
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

const (
	IDLabel        = "ID"
	KeyLabel       = "Key"
	ValueLabel     = "Value"
	NotesLabel     = "Notes"
	CategoryLabel  = "Category"
	NamespaceLabel = "Namespace"
	ReferenceLabel = "Reference"
	MediaTypeLabel = "Media Type"
	TagsLabel      = "Tags"
)

const (
	MediaCategory = "media"
	mediaFolder   = "media"
)

// magic values
const (
	maxAllowedGetAllKBLimit = 100
)

// media file type
var (
	MediaExtensions = []string{
		"apng", "avif", "csv", "gif", "jpeg",
		"mp4", "pdf", "png", "svg", "tar.gz",
		"txt", "webp", "yaml", "zip",
	}
)

var (
	errEmptyKBKey       = errors.New("kb key is empty")
	errEmptyKBValue     = errors.New("kb value is empty")
	errEmptyKBCategory  = errors.New("kb category is empty")
	errEmptyKBNamespace = errors.New("kb namespace is empty")
	errEmptyKBTags      = errors.New("kb tags is empty")
	errKBTagValues      = errors.New("kb tag must contain only alphabetic characters")
	errMinGetAllKBLimit = errors.New("minimum number of records to retrieve is no valid")
)

var IsLetter = regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString

func (s *SearchResult) Keys() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range s.Items {
			if !yield(v.Key) {
				return
			}
		}
	}
}

func (s *SearchResult) Categories() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range s.Items {
			if !yield(v.Category) {
				return
			}
		}
	}
}

func (s *SearchResult) Tags() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range s.Items {
			if !yield(strings.Join(v.Tags, " ")) {
				return
			}
		}
	}
}

func (s *SearchResult) Namespaces() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range s.Items {
			if !yield(v.Namespace) {
				return
			}
		}
	}
}

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

	if n.Category == "" {
		err = errors.Join(err, errEmptyKBCategory)
	}

	if n.Namespace == "" {
		err = errors.Join(err, errEmptyKBNamespace)
	}

	if len(n.Tags) == 0 {
		err = errors.Join(err, errEmptyKBTags)
	}

	for _, tag := range n.Tags {
		if !IsLetter(tag) {
			err = errors.Join(fmt.Errorf("%q", tag), errKBTagValues)
		}
	}

	return err
}

func (k KB) validate() error {
	var err error

	if k.Key == "" {
		err = errors.Join(err, errEmptyKBKey)
	}

	if k.Value == "" {
		err = errors.Join(err, errEmptyKBValue)
	}

	if k.Category == "" {
		err = errors.Join(err, errEmptyKBCategory)
	}

	if k.Namespace == "" {
		err = errors.Join(err, errEmptyKBNamespace)
	}

	if len(k.Tags) == 0 {
		err = errors.Join(err, errEmptyKBTags)
	}

	for _, tag := range k.Tags {
		if !IsLetter(tag) {
			err = errors.Join(fmt.Errorf("%q", tag), errKBTagValues)
		}
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

func (k KB) ToYAML() ([]byte, error) {
	kbData, err := yaml.Marshal(k)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal kb to yaml: %w", err)
	}

	return kbData, nil
}

func (k KBQueryFilter) validate() error {
	if IsStringEmpty(k.Key) && IsStringEmpty(k.Keyword) {
		return fmt.Errorf("invalid data to search kbs")
	}

	return nil
}

func (n NewKB) toKB() KB {
	return KB{
		ID:        uuid.New().String(),
		Key:       strings.ToLower(n.Key),
		Value:     n.Value,
		Notes:     n.Notes,
		Category:  strings.ToLower(n.Category),
		Namespace: strings.ToLower(n.Namespace),
		Reference: n.Reference,
		Tags:      slices.Clone(n.Tags),
	}
}

func (k KB) String() string {
	return fmt.Sprintf(`%s: %s
%s: %s
%s: %s
%s: %s
%s: %s
%s: %s
%s: %s
%s: %+v
`,
		IDLabel, k.ID,
		KeyLabel, k.Key,
		ValueLabel, k.Value,
		NotesLabel, k.Notes,
		CategoryLabel, k.Category,
		ReferenceLabel, k.Reference,
		NamespaceLabel, k.Namespace,
		TagsLabel, k.Tags)
}

func (n NewKB) String() string {
	if n.MediaType == "" {
		return fmt.Sprintf(`Key: %s
Value: %s
Notes: %s
Category: %s
Reference: %s
Namespace: %s
Tags: %+v
`, n.Key, n.Value, n.Notes, n.Category, n.Reference, n.Namespace, n.Tags)
	}
	return fmt.Sprintf(`Key: %s
Value: %s
Notes: %s
Category: %s
Reference: %s
Namespace: %s
Media Type: %s
Tags: %+v
`, n.Key, n.Value, n.Notes, n.Category, n.Reference, n.Namespace, n.MediaType, n.Tags)
}

func (i *ImportResult) Ok() bool {
	return len(i.FailedKeys) == 0 && len(i.NewIDs) > 0
}

func (i *ImportResult) anyError() bool {
	return len(i.FailedKeys) > 0
}

func (s *SyncResult) Ok() bool {
	return len(s.FailedKeys) == 0 && len(s.NewIDs) > 0
}

func (s *SyncResult) anyError() bool {
	return len(s.FailedKeys) > 0
}

func (s *SyncResult) Empty() bool {
	return len(s.FailedKeys) == 0 && len(s.NewIDs) == 0
}

func IsStringEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func loadSyncFile(syncFile string) ([]NewKB, error) {
	file, err := filesystems.ReadFile(syncFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read file for synchronization (%q): %w", syncFile, err)
	}

	dec := yaml.NewDecoder(bytes.NewReader(file))

	var kbItems []NewKB
	var kbItem NewKB
	for dec.Decode(&kbItem) == nil {
		kbItems = append(kbItems, kbItem)
		kbItem = NewKB{}
	}

	return kbItems, nil
}

func GetLongerText(title string, iterator iter.Seq[string]) int {
	textLength := len(title)
	for value := range iterator {
		if len(value) > textLength {
			textLength = len(value)
		}
	}
	return textLength
}

func (s *SearchResult) TotalPages() int {
	return int(math.Ceil(float64(s.Total) / float64(s.Limit)))
}

func (k KBItem) ToArray() []string {
	return []string{k.Key, k.Category, k.Namespace, strings.Join(k.Tags, " ")}
}

func (k KBQueryFilter) nothingToLookFor() bool {
	return IsStringEmpty(k.Key) &&
		IsStringEmpty(k.Keyword) &&
		IsStringEmpty(k.Category) &&
		IsStringEmpty(k.Namespace)
}

func (g KBQueryFilter) valid() error {
	var err error

	if g.Limit < 1 {
		err = errors.Join(err, errMinGetAllKBLimit)
	}

	if g.Limit > maxAllowedGetAllKBLimit {
		err = errors.Join(err, fmt.Errorf("invalid limit number, maximum should be %d", maxAllowedGetAllKBLimit))
	}

	return err
}

func isWebURL(anURL string) bool {
	result, err := url.ParseRequestURI(anURL)
	if err != nil {
		return false
	}

	if result == nil {
		return false
	}

	return true
}
