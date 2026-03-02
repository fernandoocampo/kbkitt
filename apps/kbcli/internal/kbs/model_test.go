package kbs_test

import (
	"context"
	"errors"
	"iter"
	"slices"
	"testing"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKBQueryFilter(t *testing.T) {
	filter := kbs.NewKBQueryFilter("mykey", "keyword")

	assert.Equal(t, "mykey", filter.Key)
	assert.Equal(t, "keyword", filter.Keyword)
	assert.Equal(t, uint32(5), filter.Limit)
	assert.Equal(t, uint32(0), filter.Offset)
}

func TestSearchResultKeysEarlyExit(t *testing.T) {
	result := &kbs.SearchResult{
		Items: []kbs.KBItem{{Key: "alpha"}, {Key: "beta"}},
	}

	seq := result.Keys()
	var keys []string
	seq(func(s string) bool {
		keys = append(keys, s)
		return false // stop after first item
	})

	assert.Len(t, keys, 1)
}

func TestSearchResultCategoriesEarlyExit(t *testing.T) {
	result := &kbs.SearchResult{
		Items: []kbs.KBItem{{Category: "cat1"}, {Category: "cat2"}},
	}

	seq := result.Categories()
	var cats []string
	seq(func(s string) bool {
		cats = append(cats, s)
		return false
	})

	assert.Len(t, cats, 1)
}

func TestSearchResultTagsEarlyExit(t *testing.T) {
	result := &kbs.SearchResult{
		Items: []kbs.KBItem{{Tags: []string{"a"}}, {Tags: []string{"b"}}},
	}

	seq := result.Tags()
	var tags []string
	seq(func(s string) bool {
		tags = append(tags, s)
		return false
	})

	assert.Len(t, tags, 1)
}

func TestSearchResultNamespacesEarlyExit(t *testing.T) {
	result := &kbs.SearchResult{
		Items: []kbs.KBItem{{Namespace: "ns1"}, {Namespace: "ns2"}},
	}

	seq := result.Namespaces()
	var ns []string
	seq(func(s string) bool {
		ns = append(ns, s)
		return false
	})

	assert.Len(t, ns, 1)
}

func TestSearchResultKeys(t *testing.T) {
	result := &kbs.SearchResult{
		Items: []kbs.KBItem{
			{ID: "1", Key: "alpha"},
			{ID: "2", Key: "beta"},
		},
	}

	keys := slices.Collect(result.Keys())

	assert.Equal(t, []string{"alpha", "beta"}, keys)
}

func TestSearchResultCategories(t *testing.T) {
	result := &kbs.SearchResult{
		Items: []kbs.KBItem{
			{ID: "1", Category: "cat1"},
			{ID: "2", Category: "cat2"},
		},
	}

	cats := slices.Collect(result.Categories())

	assert.Equal(t, []string{"cat1", "cat2"}, cats)
}

func TestSearchResultTags(t *testing.T) {
	result := &kbs.SearchResult{
		Items: []kbs.KBItem{
			{ID: "1", Tags: []string{"a", "b"}},
			{ID: "2", Tags: []string{"c"}},
		},
	}

	tags := slices.Collect(result.Tags())

	assert.Equal(t, []string{"a b", "c"}, tags)
}

func TestSearchResultNamespaces(t *testing.T) {
	result := &kbs.SearchResult{
		Items: []kbs.KBItem{
			{ID: "1", Namespace: "ns1"},
			{ID: "2", Namespace: "ns2"},
		},
	}

	ns := slices.Collect(result.Namespaces())

	assert.Equal(t, []string{"ns1", "ns2"}, ns)
}

func TestSearchResultTotalPages(t *testing.T) {
	result := &kbs.SearchResult{Total: 25, Limit: 10}
	assert.Equal(t, 3, result.TotalPages())
}

func TestSearchResultTotalPagesExact(t *testing.T) {
	result := &kbs.SearchResult{Total: 20, Limit: 10}
	assert.Equal(t, 2, result.TotalPages())
}

func TestKBItemToArray(t *testing.T) {
	item := kbs.KBItem{
		Key:       "mykey",
		Category:  "mycat",
		Namespace: "myns",
		Tags:      []string{"a", "b"},
	}

	arr := item.ToArray()

	assert.Equal(t, []string{"mykey", "mycat", "myns", "a b"}, arr)
}

func TestKBToYAML(t *testing.T) {
	kb := kbs.KB{
		Key:       "halving",
		Value:     "Bitcoin value",
		Notes:     "Notes",
		Category:  "bitcoin",
		Namespace: "cryptos",
		Tags:      []string{"bitcoin"},
	}

	data, err := kb.ToYAML()

	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "halving")
}

func TestImportResultOk(t *testing.T) {
	result := &kbs.ImportResult{
		NewIDs:     map[string]string{"key": "id"},
		FailedKeys: map[string]string{},
	}
	assert.True(t, result.Ok())
}

func TestImportResultNotOkWhenFailed(t *testing.T) {
	result := &kbs.ImportResult{
		NewIDs:     map[string]string{},
		FailedKeys: map[string]string{"key": "err"},
	}
	assert.False(t, result.Ok())
}

func TestImportResultNotOkWhenEmpty(t *testing.T) {
	result := &kbs.ImportResult{
		NewIDs:     map[string]string{},
		FailedKeys: map[string]string{},
	}
	assert.False(t, result.Ok())
}

func TestSyncResultOk(t *testing.T) {
	result := &kbs.SyncResult{
		NewIDs:     map[string]string{"key": "id"},
		FailedKeys: map[string]string{},
	}
	assert.True(t, result.Ok())
}

func TestSyncResultNotOk(t *testing.T) {
	result := &kbs.SyncResult{
		NewIDs:     map[string]string{},
		FailedKeys: map[string]string{"key": "err"},
	}
	assert.False(t, result.Ok())
}

func TestSyncResultEmpty(t *testing.T) {
	result := &kbs.SyncResult{
		NewIDs:     map[string]string{},
		FailedKeys: map[string]string{},
	}
	assert.True(t, result.Empty())
}

func TestSyncResultNotEmpty(t *testing.T) {
	result := &kbs.SyncResult{
		NewIDs:     map[string]string{"key": "id"},
		FailedKeys: map[string]string{},
	}
	assert.False(t, result.Empty())
}

func TestNewServerError(t *testing.T) {
	err := kbs.NewServerError("server failed")
	assert.Equal(t, "server failed", err.Error())
}

func TestNewServerErrorWithWrapper(t *testing.T) {
	inner := errors.New("inner error")
	err := kbs.NewServerErrorWithWrapper("server failed", inner)
	assert.Equal(t, "server failed: inner error", err.Error())
	assert.Equal(t, inner, err.Unwrap())
}

func TestNewClientError(t *testing.T) {
	err := kbs.NewClientError("bad request")
	assert.Equal(t, "bad request", err.Error())
}

func TestNewClientErrorWithWrapper(t *testing.T) {
	inner := errors.New("inner error")
	err := kbs.NewClientErrorWithWrapper("bad request", inner)
	assert.Equal(t, "bad request: inner error", err.Error())
	assert.Equal(t, inner, err.Unwrap())
}

func TestIsStringEmpty(t *testing.T) {
	assert.True(t, kbs.IsStringEmpty(""))
	assert.True(t, kbs.IsStringEmpty("   "))
	assert.False(t, kbs.IsStringEmpty("value"))
}

func TestGetLongerText(t *testing.T) {
	title := "short"
	seq := func(yield func(string) bool) {
		yield("longer text")
		yield("x")
	}

	length := kbs.GetLongerText(title, iter.Seq[string](seq))

	assert.Equal(t, len("longer text"), length)
}

func TestGetLongerTextTitleWins(t *testing.T) {
	title := "very long title indeed"
	seq := func(yield func(string) bool) {
		yield("short")
	}

	length := kbs.GetLongerText(title, iter.Seq[string](seq))

	assert.Equal(t, len(title), length)
}

func TestKBString(t *testing.T) {
	kb := kbs.KB{
		ID:        "uuid",
		Key:       "mykey",
		Value:     "myvalue",
		Notes:     "mynotes",
		Category:  "mycat",
		Namespace: "myns",
		Tags:      []string{"a", "b"},
	}

	s := kb.String()

	assert.Contains(t, s, "mykey")
	assert.Contains(t, s, "myvalue")
}

func TestNewKBStringNoMediaType(t *testing.T) {
	newKB := kbs.NewKB{
		Key:       "mykey",
		Value:     "myvalue",
		Notes:     "Notes",
		Category:  "mycat",
		Namespace: "myns",
		Tags:      []string{"tag"},
	}

	s := newKB.String()

	assert.Contains(t, s, "mykey")
	assert.NotContains(t, s, "Media Type")
}

func TestNewKBStringWithMediaType(t *testing.T) {
	newKB := kbs.NewKB{
		Key:       "mykey",
		Value:     "https://example.com/image.jpg",
		Notes:     "Notes",
		Category:  "media",
		Namespace: "general",
		MediaType: "jpeg",
		Tags:      []string{"image"},
	}

	s := newKB.String()

	assert.Contains(t, s, "mykey")
	assert.Contains(t, s, "jpeg")
}

func TestGetAllKBsLimitOverMax(t *testing.T) {
	filter := kbs.KBQueryFilter{Limit: 101}

	ctx := context.TODO()
	svc := kbs.NewService(kbs.ServiceSetup{})

	_, err := svc.GetAllKBs(ctx, filter)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid limit number")
}
