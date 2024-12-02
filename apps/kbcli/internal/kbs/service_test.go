package kbs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
