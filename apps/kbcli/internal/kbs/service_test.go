package kbs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
	"github.com/stretchr/testify/assert"
)

func TestAddKBInvalidValues(t *testing.T) {
	// Given
	newKB := kbs.NewKB{
		Key:   "halving",
		Notes: "Bitcoins have a finite supply, which makes them a scarce digital commodity",
		Tags:  []string{"bitcoin", "halving"},
	}
	expectedError := errors.New("the given values are not valid: kb value is empty\nkb kind is empty")

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
	withKBIDDummy := "88ac1fa1-2cdd-4f64-a4a3-13c6d162f504"
	newKB := kbs.NewKB{
		Key:   "halving",
		Value: "The number of bitcoins generated per block is decreased 50% every four years",
		Notes: "Bitcoins have a finite supply, which makes them a scarce digital commodity",
		Kind:  "bitcoin",
		Tags:  []string{"bitcoin", "halving"},
	}
	expectedKB := &kbs.KB{
		ID:    "88ac1fa1-2cdd-4f64-a4a3-13c6d162f504",
		Key:   "halving",
		Value: "The number of bitcoins generated per block is decreased 50% every four years",
		Notes: "Bitcoins have a finite supply, which makes them a scarce digital commodity",
		Kind:  "bitcoin",
		Tags:  []string{"bitcoin", "halving"},
	}
	ctx := context.TODO()
	settings := kbs.ServiceSetup{
		KBClient: newKBServiceDummyWithCreate(withKBIDDummy, nil),
	}
	kbService := kbs.NewService(settings)
	// When
	gotKB, err := kbService.Add(ctx, newKB)
	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedKB, gotKB)
}

type kbServiceDummy struct {
	createKBID    string
	newKBInput    *kbs.NewKB
	createKBError error
}

func newKBServiceDummyWithCreate(kbid string, err error) *kbServiceDummy {
	return &kbServiceDummy{
		createKBID:    kbid,
		createKBError: err,
	}
}

func (k *kbServiceDummy) Create(ctx context.Context, newKB kbs.NewKB) (string, error) {
	k.newKBInput = &newKB

	if k.createKBError != nil {
		return "", k.createKBError
	}

	return k.createKBID, nil
}

func (k *kbServiceDummy) Search(ctx context.Context, filter kbs.KBQueryFilter) (*kbs.SearchResult, error) {
	return nil, nil
}

func (k *kbServiceDummy) Get(ctx context.Context, id string) (*kbs.KB, error) {
	return nil, nil
}
