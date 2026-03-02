package storages_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/storages"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) *storages.SQLite {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	t.Cleanup(func() { db.Close() })

	setup := storages.SQLiteSetup{DB: db}
	storage := storages.NewSQLite(&setup)

	err = storage.InitializeDB(context.Background())
	require.NoError(t, err)

	return storage
}

func makeTestKB() kbs.KB {
	return kbs.KB{
		ID:        "test-uuid-1234-5678-9012-3456",
		Key:       "bitcoin-halving",
		Value:     "The number of bitcoins generated per block is decreased 50% every four years",
		Notes:     "Bitcoins have a finite supply",
		Category:  "bitcoin",
		Namespace: "cryptos",
		Reference: "https://bitcoin.org",
		Tags:      []string{"bitcoin", "halving"},
	}
}

// ---- Create ----

func TestCreateKB(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()
	kb := makeTestKB()

	id, err := storage.Create(ctx, kb)

	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestCreateKBDuplicateKey(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()
	kb := makeTestKB()

	_, err := storage.Create(ctx, kb)
	require.NoError(t, err)

	_, err = storage.Create(ctx, kb)

	assert.Error(t, err)
}

// ---- GetByID ----

func TestGetByIDExisting(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()
	kb := makeTestKB()

	_, err := storage.Create(ctx, kb)
	require.NoError(t, err)

	result, err := storage.GetByID(ctx, kb.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, kb.ID, result.ID)
	assert.Equal(t, kb.Key, result.Key)
	assert.Equal(t, kb.Value, result.Value)
	assert.Equal(t, kb.Category, result.Category)
	assert.Equal(t, kb.Namespace, result.Namespace)
}

func TestGetByIDNotFound(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()

	result, err := storage.GetByID(ctx, "nonexistent-id")

	assert.NoError(t, err)
	assert.Nil(t, result)
}

// ---- GetByKey ----

func TestGetByKeyExisting(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()
	kb := makeTestKB()

	_, err := storage.Create(ctx, kb)
	require.NoError(t, err)

	result, err := storage.GetByKey(ctx, kb.Key)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, kb.ID, result.ID)
	assert.Equal(t, kb.Key, result.Key)
}

func TestGetByKeyNotFound(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()

	result, err := storage.GetByKey(ctx, "nonexistent-key")

	assert.NoError(t, err)
	assert.Nil(t, result)
}

// ---- Update ----

func TestUpdateKB(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()
	kb := makeTestKB()

	_, err := storage.Create(ctx, kb)
	require.NoError(t, err)

	kb.Value = "Updated value"
	kb.Notes = "Updated notes"
	kb.Tags = []string{"bitcoin", "halving", "updated"}

	err = storage.Update(ctx, &kb)
	require.NoError(t, err)

	result, err := storage.GetByKey(ctx, kb.Key)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated value", result.Value)
	assert.Equal(t, "Updated notes", result.Notes)
}

// ---- GetAll ----

func TestGetAllEmpty(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()
	filter := kbs.KBQueryFilter{Limit: 10, Offset: 0}

	result, err := storage.GetAll(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, result.Total)
	assert.Empty(t, result.KBs)
}

func TestGetAllWithRecords(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()

	kb1 := makeTestKB()
	_, err := storage.Create(ctx, kb1)
	require.NoError(t, err)

	kb2 := kbs.KB{
		ID:        "test-uuid-different-0000",
		Key:       "ethereum",
		Value:     "Ethereum is a decentralized platform",
		Notes:     "Notes about ethereum",
		Category:  "crypto",
		Namespace: "blockchain",
		Tags:      []string{"ethereum", "blockchain"},
	}
	_, err = storage.Create(ctx, kb2)
	require.NoError(t, err)

	filter := kbs.KBQueryFilter{Limit: 10, Offset: 0}

	result, err := storage.GetAll(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.KBs, 2)
}

func TestGetAllWithCategoryFilter(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()

	kb1 := makeTestKB()
	_, err := storage.Create(ctx, kb1)
	require.NoError(t, err)

	kb2 := kbs.KB{
		ID:        "other-uuid-1234",
		Key:       "ethereum",
		Value:     "Ethereum is a platform",
		Notes:     "Notes",
		Category:  "other",
		Namespace: "blockchain",
		Tags:      []string{"ethereum"},
	}
	_, err = storage.Create(ctx, kb2)
	require.NoError(t, err)

	filter := kbs.KBQueryFilter{Category: "bitcoin", Limit: 10}

	result, err := storage.GetAll(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.KBs, 1)
}

func TestGetAllPagination(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		kb := kbs.KB{
			ID:        fmt.Sprintf("test-uuid-page-%04d", i),
			Key:       fmt.Sprintf("key-%d", i),
			Value:     fmt.Sprintf("Value %d", i),
			Notes:     "Notes",
			Category:  "test",
			Namespace: "testns",
			Tags:      []string{"tag"},
		}
		_, err := storage.Create(ctx, kb)
		require.NoError(t, err)
	}

	filter := kbs.KBQueryFilter{Limit: 2, Offset: 0}

	result, err := storage.GetAll(ctx, filter)

	require.NoError(t, err)
	assert.Equal(t, 5, result.Total)
	assert.Len(t, result.KBs, 2)
}

// ---- Search ----

func TestSearchByKey(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()
	kb := makeTestKB()

	_, err := storage.Create(ctx, kb)
	require.NoError(t, err)

	filter := kbs.KBQueryFilter{Key: "halving", Limit: 10}

	result, err := storage.Search(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, kb.Key, result.Items[0].Key)
}

func TestSearchByCategory(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()
	kb := makeTestKB()

	_, err := storage.Create(ctx, kb)
	require.NoError(t, err)

	filter := kbs.KBQueryFilter{Category: "bitcoin", Limit: 10}

	result, err := storage.Search(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
}

func TestSearchByNamespace(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()
	kb := makeTestKB()

	_, err := storage.Create(ctx, kb)
	require.NoError(t, err)

	filter := kbs.KBQueryFilter{Namespace: "cryptos", Limit: 10}

	result, err := storage.Search(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
}

func TestSearchByKeyword(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()
	kb := makeTestKB()

	_, err := storage.Create(ctx, kb)
	require.NoError(t, err)

	filter := kbs.KBQueryFilter{Keyword: "bitcoin", Limit: 10}

	result, err := storage.Search(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Items, 1)
}

func TestSearchByKeyAndCategory(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()
	kb := makeTestKB()

	_, err := storage.Create(ctx, kb)
	require.NoError(t, err)

	// Multiple filter conditions exercises the AND operator in filterBuilder
	filter := kbs.KBQueryFilter{Key: "halving", Category: "bitcoin", Limit: 10}

	result, err := storage.Search(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
}

func TestSearchByKeywordAndNamespace(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()
	kb := makeTestKB()

	_, err := storage.Create(ctx, kb)
	require.NoError(t, err)

	filter := kbs.KBQueryFilter{Keyword: "bitcoin", Namespace: "cryptos", Limit: 10}

	result, err := storage.Search(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
}

func TestSearchNoResults(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()

	filter := kbs.KBQueryFilter{Key: "nonexistent", Limit: 10}

	result, err := storage.Search(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, result.Total)
	assert.Empty(t, result.Items)
}

// ---- CountByCategory ----

func TestCountByCategoryEmpty(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()

	count, err := storage.CountByCategory(ctx, "bitcoin")

	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestCountByCategory(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()

	kb := makeTestKB()
	_, err := storage.Create(ctx, kb)
	require.NoError(t, err)

	count, err := storage.CountByCategory(ctx, "bitcoin")

	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

// ---- Version ----

func TestVersion(t *testing.T) {
	storage := newTestDB(t)

	version, err := storage.Version()

	require.NoError(t, err)
	assert.NotEmpty(t, version)
}

// ---- Close ----

func TestClose(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	setup := storages.SQLiteSetup{DB: db}
	storage := storages.NewSQLite(&setup)

	// Close should not panic
	assert.NotPanics(t, func() { storage.Close() })
}

func TestCloseNilDB(t *testing.T) {
	setup := storages.SQLiteSetup{DB: nil}
	storage := storages.NewSQLite(&setup)

	// Close with nil DB should not panic
	assert.NotPanics(t, func() { storage.Close() })
}

// ---- Error paths via closed DB ----

func newClosedDB(t *testing.T) *storages.SQLite {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	setup := storages.SQLiteSetup{DB: db}
	storage := storages.NewSQLite(&setup)

	err = storage.InitializeDB(context.Background())
	require.NoError(t, err)

	db.Close() // close to force errors on subsequent calls

	return storage
}

func TestSearchErrorOnClosedDB(t *testing.T) {
	storage := newClosedDB(t)
	ctx := context.Background()

	filter := kbs.KBQueryFilter{Key: "test", Limit: 10}

	_, err := storage.Search(ctx, filter)

	assert.Error(t, err)
}

func TestSearchKeywordErrorOnClosedDB(t *testing.T) {
	storage := newClosedDB(t)
	ctx := context.Background()

	filter := kbs.KBQueryFilter{Keyword: "test", Limit: 10}

	_, err := storage.Search(ctx, filter)

	assert.Error(t, err)
}

func TestGetAllErrorOnClosedDB(t *testing.T) {
	storage := newClosedDB(t)
	ctx := context.Background()

	filter := kbs.KBQueryFilter{Limit: 10}

	_, err := storage.GetAll(ctx, filter)

	assert.Error(t, err)
}

func TestCountByCategoryErrorOnClosedDB(t *testing.T) {
	storage := newClosedDB(t)
	ctx := context.Background()

	_, err := storage.CountByCategory(ctx, "bitcoin")

	assert.Error(t, err)
}

// ---- CreateSQLiteConnection ----

func TestCreateSQLiteConnection(t *testing.T) {
	dbPath := fmt.Sprintf("%s/test.db", t.TempDir())

	db, err := storages.CreateSQLiteConnection(dbPath)

	require.NoError(t, err)
	require.NotNil(t, db)
	db.Close()
}

func TestCountByCategoryMultiple(t *testing.T) {
	storage := newTestDB(t)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		kb := kbs.KB{
			ID:        fmt.Sprintf("test-uuid-cat-%04d", i),
			Key:       fmt.Sprintf("quote-%d", i),
			Value:     fmt.Sprintf("Quote value %d", i),
			Notes:     "Notes",
			Category:  "quote",
			Namespace: "general",
			Tags:      []string{"inspirational"},
		}
		_, err := storage.Create(ctx, kb)
		require.NoError(t, err)
	}

	count, err := storage.CountByCategory(ctx, "quote")

	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}
