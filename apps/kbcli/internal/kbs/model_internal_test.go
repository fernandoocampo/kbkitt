package kbs

// Internal test file to access package-private functions.

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomNumberZero(t *testing.T) {
	assert.Equal(t, int64(0), newRandomNumber(0))
}

func TestNewRandomNumberPositive(t *testing.T) {
	// Just verify it returns a value in [0, n)
	n := int64(100)
	result := newRandomNumber(n)
	assert.GreaterOrEqual(t, result, int64(0))
	assert.Less(t, result, n)
}

func TestIsWebURLTrue(t *testing.T) {
	assert.True(t, isWebURL("https://example.com/image.png"))
	assert.True(t, isWebURL("http://example.com"))
}

func TestIsWebURLFalse(t *testing.T) {
	assert.False(t, isWebURL("not-a-url"))
	assert.False(t, isWebURL(""))
	assert.False(t, isWebURL("relative/path"))
}

func TestKBQueryFilterValidate(t *testing.T) {
	valid := KBQueryFilter{Keyword: "test"}
	assert.NoError(t, valid.validate())

	validKey := KBQueryFilter{Key: "mykey"}
	assert.NoError(t, validKey.validate())

	invalid := KBQueryFilter{}
	assert.Error(t, invalid.validate())
}

func TestImportResultAnyError(t *testing.T) {
	r := &ImportResult{FailedKeys: map[string]string{"key": "err"}}
	assert.True(t, r.anyError())

	r2 := &ImportResult{FailedKeys: map[string]string{}}
	assert.False(t, r2.anyError())
}

func TestSyncResultAnyError(t *testing.T) {
	r := &SyncResult{FailedKeys: map[string]string{"key": "err"}}
	assert.True(t, r.anyError())

	r2 := &SyncResult{FailedKeys: map[string]string{}}
	assert.False(t, r2.anyError())
}
