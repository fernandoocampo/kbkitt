package storages_test

import (
	"context"
	"flag"
	"testing"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/storages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var E2ETest = flag.Bool("e2e-test", false, "this flags indicates if this is an e2e test")

func TestUpdateUser(t *testing.T) {
	if !*E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx := context.TODO()

	db, err := storages.CreateSQLiteConnection("/Users/Fernando_Ocampo/.kbkitt/kbkitt.db")
	if err != nil {
		require.NoError(t, err)
	}
	storageSetup := storages.SQLiteSetup{
		DB: db,
	}
	storage := storages.NewSQLite(&storageSetup)
	defer storage.Close()

	existingKB, err := storage.GetByID(ctx, "e442e8a9-7417-4a25-9299-e25a1c2c11f3")
	if err != nil {
		require.NoError(t, err)
	}

	// When
	err = storage.Update(ctx, existingKB)

	// Then
	assert.NoError(t, err)
}
