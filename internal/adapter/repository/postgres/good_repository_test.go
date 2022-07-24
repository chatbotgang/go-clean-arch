package postgres

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"github.com/chatbotgang/go-clean-architecture-template/testdata"
)

func assertGood(t *testing.T, expected *barter.Good, actual *barter.Good) {
	require.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.OwnerID, actual.OwnerID)
}

func TestPostgresRepository_CreateGood(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db, testdata.Path(testdata.TestDataTrader))

	// Args
	type Args struct {
		barter.Good
	}
	var args Args
	_ = faker.FakeData(&args)
	traderID := 1
	args.OwnerID = traderID

	good, err := repo.CreateGood(context.Background(), args.Good)
	require.NoError(t, err)
	assertGood(t, &args.Good, good)
}

func TestPostgresRepository_GetGoodByID(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataTrader),
		testdata.Path(testdata.TestDataGood),
	)
	goodID := 1

	_, err := repo.GetGoodByID(context.Background(), goodID)
	require.NoError(t, err)
}

func TestPostgresRepository_ListGoods(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataTrader),
		testdata.Path(testdata.TestDataGood),
	)

	goods, err := repo.ListGoods(context.Background())
	require.NoError(t, err)
	assert.Len(t, goods, 2)
}

func TestPostgresRepository_ListGoodsByOwner(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataTrader),
		testdata.Path(testdata.TestDataGood),
	)
	traderID := 500

	goods, err := repo.ListGoodsByOwner(context.Background(), traderID)
	require.NoError(t, err)
	assert.Len(t, goods, 0)
}

func TestPostgresRepository_UpdateGood(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataTrader),
		testdata.Path(testdata.TestDataGood),
	)
	goodID := 1

	good, err := repo.GetGoodByID(context.Background(), goodID)
	require.NoError(t, err)

	newName := "good 123"
	good.Name = newName

	updatedGood, err := repo.UpdateGood(context.Background(), *good)
	require.NoError(t, err)
	assertGood(t, good, updatedGood)
}

func TestPostgresRepository_DeleteGoodByID(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataTrader),
		testdata.Path(testdata.TestDataGood),
	)
	goodID := 1

	_, err := repo.GetGoodByID(context.Background(), goodID)
	require.NoError(t, err)

	err = repo.DeleteGoodByID(context.Background(), goodID)
	require.NoError(t, err)

	_, err = repo.GetGoodByID(context.Background(), goodID)
	require.Error(t, err)
}
