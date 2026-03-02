package kbs_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ---- Add ----

func TestAddKBInvalidValues(t *testing.T) {
	// Given
	newKB := kbs.NewKB{
		Key:   "halving",
		Notes: "Bitcoins have a finite supply, which makes them a scarce digital commodity",
		Tags:  []string{"bitcoin", "halving"},
	}
	expectedError := errors.New("the given values are not valid: kb value is empty\nkb category is empty\nkb namespace is empty")

	ctx := context.TODO()
	settings := kbs.ServiceSetup{}
	kbService := kbs.NewService(settings)
	// When
	_, err := kbService.Add(ctx, newKB)
	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError.Error(), err.Error())
}

func TestAddKB(t *testing.T) {
	// Given
	newKB := kbs.NewKB{
		Key:       "halving",
		Value:     "The number of bitcoins generated per block is decreased 50% every four years",
		Notes:     "Bitcoins have a finite supply, which makes them a scarce digital commodity",
		Category:  "bitcoin",
		Namespace: "cryptos",
		Tags:      []string{"bitcoin", "halving", "bitcoin"},
	}
	expectedKB := &kbs.KB{
		ID:        "88ac1fa1-2cdd-4f64-a4a3-13c6d162f504",
		Key:       "halving",
		Value:     "The number of bitcoins generated per block is decreased 50% every four years",
		Notes:     "Bitcoins have a finite supply, which makes them a scarce digital commodity",
		Category:  "bitcoin",
		Namespace: "cryptos",
		Tags:      []string{"bitcoin", "halving"},
	}
	ctx := context.TODO()
	storageMock := newStorageMock()
	storageMock.On("Create", ctx, mock.AnythingOfType("kbs.KB")).Return("1", nil)

	settings := kbs.ServiceSetup{
		KBStorage: storageMock,
	}
	kbService := kbs.NewService(settings)
	// When
	gotKB, err := kbService.Add(ctx, newKB)
	// Then
	assert.NoError(t, err)
	expectedKB.ID = gotKB.ID
	assert.Equal(t, expectedKB, gotKB)
}

func TestAddKBInvalidTagChars(t *testing.T) {
	newKB := kbs.NewKB{
		Key:       "halving",
		Value:     "value",
		Notes:     "notes",
		Category:  "bitcoin",
		Namespace: "cryptos",
		Tags:      []string{"valid", "invalid tag!"},
	}

	ctx := context.TODO()
	svc := kbs.NewService(kbs.ServiceSetup{})

	_, err := svc.Add(ctx, newKB)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "alphabetic")
}

// ---- Update ----

func TestUpdateKBInvalidValues(t *testing.T) {
	kb := kbs.KB{
		Key: "halving",
		// Missing: Value, Category, Namespace, Tags
	}

	ctx := context.TODO()
	settings := kbs.ServiceSetup{}
	kbService := kbs.NewService(settings)

	err := kbService.Update(ctx, kb)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "the given values are not valid")
}

func TestUpdateKB(t *testing.T) {
	kb := kbs.KB{
		ID:        "some-uuid",
		Key:       "halving",
		Value:     "The number of bitcoins generated per block is decreased 50% every four years",
		Notes:     "Bitcoins have a finite supply",
		Category:  "bitcoin",
		Namespace: "cryptos",
		Tags:      []string{"bitcoin", "halving"},
	}

	ctx := context.TODO()
	storageMock := newStorageMock()
	storageMock.On("Update", ctx, mock.AnythingOfType("*kbs.KB")).Return(nil)

	settings := kbs.ServiceSetup{KBStorage: storageMock}
	kbService := kbs.NewService(settings)

	err := kbService.Update(ctx, kb)

	assert.NoError(t, err)
	storageMock.AssertExpectations(t)
}

// ---- Import ----

func TestImportKBsInvalidValues(t *testing.T) {
	newKBs := []kbs.NewKB{
		{Key: "halving"}, // missing Value, Category, Namespace, Tags
	}

	ctx := context.TODO()
	settings := kbs.ServiceSetup{}
	kbService := kbs.NewService(settings)

	result, err := kbService.Import(ctx, newKBs)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestImportKBs(t *testing.T) {
	newKBs := []kbs.NewKB{
		{
			Key:       "halving",
			Value:     "The number of bitcoins generated per block is decreased 50% every four years",
			Notes:     "Bitcoins have a finite supply",
			Category:  "bitcoin",
			Namespace: "cryptos",
			Tags:      []string{"bitcoin", "halving"},
		},
	}

	ctx := context.TODO()
	storageMock := newStorageMock()
	storageMock.On("Create", ctx, mock.AnythingOfType("kbs.KB")).Return("1", nil)

	settings := kbs.ServiceSetup{KBStorage: storageMock}
	kbService := kbs.NewService(settings)

	result, err := kbService.Import(ctx, newKBs)

	require.NoError(t, err)
	assert.Len(t, result.NewIDs, 1)
	assert.Empty(t, result.FailedKeys)
	storageMock.AssertExpectations(t)
}

func TestImportKBsWithStorageError(t *testing.T) {
	newKBs := []kbs.NewKB{
		{
			Key:       "halving",
			Value:     "The number of bitcoins generated per block is decreased 50% every four years",
			Notes:     "Bitcoins have a finite supply",
			Category:  "bitcoin",
			Namespace: "cryptos",
			Tags:      []string{"bitcoin", "halving"},
		},
	}

	ctx := context.TODO()
	storageMock := newStorageMock()
	storageMock.On("Create", ctx, mock.AnythingOfType("kbs.KB")).Return("", errors.New("duplicate key"))

	settings := kbs.ServiceSetup{KBStorage: storageMock}
	kbService := kbs.NewService(settings)

	result, err := kbService.Import(ctx, newKBs)

	require.NoError(t, err)
	assert.Empty(t, result.NewIDs)
	assert.Len(t, result.FailedKeys, 1)
}

// ---- GetAllKBs ----

func TestGetAllKBsInvalidFilter(t *testing.T) {
	filter := kbs.KBQueryFilter{Limit: 0}

	ctx := context.TODO()
	settings := kbs.ServiceSetup{}
	kbService := kbs.NewService(settings)

	result, err := kbService.GetAllKBs(ctx, filter)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetAllKBs(t *testing.T) {
	filter := kbs.KBQueryFilter{Limit: 10, Offset: 0}
	expectedResult := &kbs.GetAllResult{
		KBs:   []kbs.KB{},
		Total: 0,
		Limit: 10,
	}

	ctx := context.TODO()
	storageMock := newStorageMock()
	storageMock.On("GetAll", ctx, filter).Return(expectedResult, nil)

	settings := kbs.ServiceSetup{KBStorage: storageMock}
	kbService := kbs.NewService(settings)

	result, err := kbService.GetAllKBs(ctx, filter)

	require.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	storageMock.AssertExpectations(t)
}

// ---- Search ----

func TestSearchKBsNothingToLookFor(t *testing.T) {
	filter := kbs.KBQueryFilter{}

	ctx := context.TODO()
	settings := kbs.ServiceSetup{}
	kbService := kbs.NewService(settings)

	result, err := kbService.Search(ctx, filter)

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestSearchKBs(t *testing.T) {
	filter := kbs.KBQueryFilter{Key: "halving", Limit: 5}
	expectedResult := &kbs.SearchResult{
		Items: []kbs.KBItem{},
		Total: 0,
		Limit: 5,
	}

	ctx := context.TODO()
	storageMock := newStorageMock()
	storageMock.On("Search", ctx, filter).Return(expectedResult, nil)

	settings := kbs.ServiceSetup{KBStorage: storageMock}
	kbService := kbs.NewService(settings)

	result, err := kbService.Search(ctx, filter)

	require.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	storageMock.AssertExpectations(t)
}

// ---- GetByID ----

func TestGetByIDEmptyID(t *testing.T) {
	ctx := context.TODO()
	settings := kbs.ServiceSetup{}
	kbService := kbs.NewService(settings)

	result, err := kbService.GetByID(ctx, "  ")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetByID(t *testing.T) {
	id := "some-uuid"
	expectedKB := &kbs.KB{
		ID:        id,
		Key:       "halving",
		Value:     "The number of bitcoins generated per block is decreased 50% every four years",
		Category:  "bitcoin",
		Namespace: "cryptos",
		Tags:      []string{"bitcoin", "halving"},
	}

	ctx := context.TODO()
	storageMock := newStorageMock()
	storageMock.On("GetByID", ctx, id).Return(expectedKB, nil)

	settings := kbs.ServiceSetup{KBStorage: storageMock}
	kbService := kbs.NewService(settings)

	result, err := kbService.GetByID(ctx, id)

	require.NoError(t, err)
	assert.Equal(t, expectedKB, result)
	storageMock.AssertExpectations(t)
}

// ---- GetByKey ----

func TestGetByKeyEmptyKey(t *testing.T) {
	ctx := context.TODO()
	settings := kbs.ServiceSetup{}
	kbService := kbs.NewService(settings)

	result, err := kbService.GetByKey(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetByKey(t *testing.T) {
	key := "halving"
	expectedKB := &kbs.KB{
		ID:        "some-uuid",
		Key:       key,
		Value:     "The number of bitcoins generated per block is decreased 50% every four years",
		Category:  "bitcoin",
		Namespace: "cryptos",
		Tags:      []string{"bitcoin", "halving"},
	}

	ctx := context.TODO()
	storageMock := newStorageMock()
	storageMock.On("GetByKey", ctx, key).Return(expectedKB, nil)

	settings := kbs.ServiceSetup{KBStorage: storageMock}
	kbService := kbs.NewService(settings)

	result, err := kbService.GetByKey(ctx, key)

	require.NoError(t, err)
	assert.Equal(t, expectedKB, result)
	storageMock.AssertExpectations(t)
}

// ---- GetRandomQuote ----

func TestGetRandomQuoteNoQuotes(t *testing.T) {
	ctx := context.TODO()
	storageMock := newStorageMock()
	storageMock.On("CountByCategory", ctx, kbs.QuoteCategory).Return(int64(0), nil)

	settings := kbs.ServiceSetup{KBStorage: storageMock}
	kbService := kbs.NewService(settings)

	result, err := kbService.GetRandomQuote(ctx)

	assert.NoError(t, err)
	assert.Nil(t, result)
	storageMock.AssertExpectations(t)
}

func TestGetRandomQuoteOneQuote(t *testing.T) {
	expectedKB := kbs.KB{
		ID:        "some-uuid",
		Key:       "my-quote",
		Value:     "A quote value",
		Category:  kbs.QuoteCategory,
		Namespace: "general",
		Tags:      []string{"inspirational"},
	}
	expectedFilter := kbs.KBQueryFilter{
		Category: kbs.QuoteCategory,
		Limit:    1,
		Offset:   0,
	}
	expectedResult := &kbs.GetAllResult{
		KBs:   []kbs.KB{expectedKB},
		Total: 1,
		Limit: 1,
	}

	ctx := context.TODO()
	storageMock := newStorageMock()
	storageMock.On("CountByCategory", ctx, kbs.QuoteCategory).Return(int64(1), nil)
	storageMock.On("GetAll", ctx, expectedFilter).Return(expectedResult, nil)

	settings := kbs.ServiceSetup{KBStorage: storageMock}
	kbService := kbs.NewService(settings)

	result, err := kbService.GetRandomQuote(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedKB.Key, result.Key)
	storageMock.AssertExpectations(t)
}

func TestGetRandomQuoteManyQuotes(t *testing.T) {
	expectedKB := kbs.KB{
		ID:        "some-uuid",
		Key:       "my-quote",
		Value:     "A quote value",
		Category:  kbs.QuoteCategory,
		Namespace: "general",
		Tags:      []string{"inspirational"},
	}
	expectedResult := &kbs.GetAllResult{
		KBs:   []kbs.KB{expectedKB},
		Total: 5,
		Limit: 1,
	}

	ctx := context.TODO()
	storageMock := newStorageMock()
	storageMock.On("CountByCategory", ctx, kbs.QuoteCategory).Return(int64(5), nil)
	storageMock.On("GetAll", ctx, mock.MatchedBy(func(f kbs.KBQueryFilter) bool {
		return f.Category == kbs.QuoteCategory && f.Limit == 1
	})).Return(expectedResult, nil)

	settings := kbs.ServiceSetup{KBStorage: storageMock}
	kbService := kbs.NewService(settings)

	result, err := kbService.GetRandomQuote(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedKB.Key, result.Key)
	storageMock.AssertExpectations(t)
}

// ---- SaveForSync ----

func TestSaveForSync(t *testing.T) {
	syncFilePath := filepath.Join(t.TempDir(), "sync.yaml")
	newKB := kbs.NewKB{
		Key:       "halving",
		Value:     "The number of bitcoins generated per block is decreased 50% every four years",
		Notes:     "Bitcoins have a finite supply",
		Category:  "bitcoin",
		Namespace: "cryptos",
		Tags:      []string{"bitcoin", "halving"},
	}

	ctx := context.TODO()
	settings := kbs.ServiceSetup{FileForSyncPath: syncFilePath}
	kbService := kbs.NewService(settings)

	err := kbService.SaveForSync(ctx, newKB)

	require.NoError(t, err)
	content, readErr := os.ReadFile(syncFilePath)
	require.NoError(t, readErr)
	assert.NotEmpty(t, content)
	assert.Contains(t, string(content), "halving")
}

// ---- Sync ----

func TestSyncNoFile(t *testing.T) {
	syncFilePath := filepath.Join(t.TempDir(), "sync.yaml")
	// file does not exist

	ctx := context.TODO()
	settings := kbs.ServiceSetup{FileForSyncPath: syncFilePath}
	kbService := kbs.NewService(settings)

	result, err := kbService.Sync(ctx)

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestSyncWithItems(t *testing.T) {
	syncFilePath := filepath.Join(t.TempDir(), "sync.yaml")
	syncContent := `Key: halving
Value: The number of bitcoins generated per block is decreased 50% every four years
Notes: Bitcoins have a finite supply
Category: bitcoin
Namespace: cryptos
Tags:
    - bitcoin
    - halving
`
	err := os.WriteFile(syncFilePath, []byte(syncContent), 0644)
	require.NoError(t, err)

	ctx := context.TODO()
	kbClientMock := newKBClientMock()
	kbClientMock.On("Create", ctx, mock.AnythingOfType("kbs.NewKB")).Return("new-id", nil)

	settings := kbs.ServiceSetup{
		FileForSyncPath: syncFilePath,
		KBClient:        kbClientMock,
	}
	kbService := kbs.NewService(settings)

	result, err := kbService.Sync(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.NewIDs, 1)
	assert.Empty(t, result.FailedKeys)
	kbClientMock.AssertExpectations(t)
}

func TestSyncWithPartialFailure(t *testing.T) {
	syncFilePath := filepath.Join(t.TempDir(), "sync.yaml")
	syncContent := `Key: halving
Value: The number of bitcoins generated per block is decreased 50% every four years
Notes: Bitcoins have a finite supply
Category: bitcoin
Namespace: cryptos
Tags:
    - bitcoin
    - halving
`
	err := os.WriteFile(syncFilePath, []byte(syncContent), 0644)
	require.NoError(t, err)

	ctx := context.TODO()
	kbClientMock := newKBClientMock()
	kbClientMock.On("Create", ctx, mock.AnythingOfType("kbs.NewKB")).Return("", errors.New("server error"))

	settings := kbs.ServiceSetup{
		FileForSyncPath: syncFilePath,
		KBClient:        kbClientMock,
	}
	kbService := kbs.NewService(settings)

	result, err := kbService.Sync(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, result.NewIDs)
	assert.Len(t, result.FailedKeys, 1)
}

// ---- SaveMedia ----

func TestSaveMediaNotMediaFile(t *testing.T) {
	newKB := kbs.NewKB{
		Key:   "test",
		Value: "not-a-url-or-existing-file",
	}

	ctx := context.TODO()
	settings := kbs.ServiceSetup{}
	kbService := kbs.NewService(settings)

	err := kbService.SaveMedia(ctx, newKB)

	assert.ErrorIs(t, err, kbs.ErrIsNotMediaFile)
}


// ---- Mocks ----

type storageDummy struct {
	mock.Mock
}

func newStorageMock() *storageDummy {
	return &storageDummy{}
}

func (k *storageDummy) Create(ctx context.Context, newKB kbs.KB) (string, error) {
	args := k.Called(ctx, newKB)

	return args.String(0), args.Error(1)
}

func (k *storageDummy) Search(ctx context.Context, filter kbs.KBQueryFilter) (*kbs.SearchResult, error) {
	args := k.Called(ctx, filter)

	return args.Get(0).(*kbs.SearchResult), args.Error(1)
}

func (k *storageDummy) GetByID(ctx context.Context, id string) (*kbs.KB, error) {
	args := k.Called(ctx, id)

	return args.Get(0).(*kbs.KB), args.Error(1)
}

func (k *storageDummy) GetByKey(ctx context.Context, id string) (*kbs.KB, error) {
	args := k.Called(ctx, id)

	return args.Get(0).(*kbs.KB), args.Error(1)
}

func (k *storageDummy) Update(ctx context.Context, kb *kbs.KB) error {
	args := k.Called(ctx, kb)

	return args.Error(0)
}

func (k *storageDummy) GetAll(ctx context.Context, filter kbs.KBQueryFilter) (*kbs.GetAllResult, error) {
	args := k.Called(ctx, filter)

	return args.Get(0).(*kbs.GetAllResult), args.Error(1)
}

func (k *storageDummy) CountByCategory(ctx context.Context, category string) (int64, error) {
	args := k.Called(ctx, category)

	return args.Get(0).(int64), args.Error(1)
}

type kbClientDummy struct {
	mock.Mock
}

func newKBClientMock() *kbClientDummy {
	return &kbClientDummy{}
}

func (k *kbClientDummy) Create(ctx context.Context, newKB kbs.NewKB) (string, error) {
	args := k.Called(ctx, newKB)

	return args.String(0), args.Error(1)
}

func (k *kbClientDummy) Update(ctx context.Context, kb *kbs.KB) error {
	args := k.Called(ctx, kb)

	return args.Error(0)
}

func (k *kbClientDummy) Search(ctx context.Context, filter kbs.KBQueryFilter) (*kbs.SearchResult, error) {
	args := k.Called(ctx, filter)

	return args.Get(0).(*kbs.SearchResult), args.Error(1)
}

func (k *kbClientDummy) Get(ctx context.Context, id string) (*kbs.KB, error) {
	args := k.Called(ctx, id)

	return args.Get(0).(*kbs.KB), args.Error(1)
}
