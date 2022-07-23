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

func assertTrader(t *testing.T, expected *barter.Trader, actual *barter.Trader) {
	require.NotNil(t, actual)
	assert.Equal(t, expected.UID, actual.UID)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Name, actual.Name)
}

func TestPostgresRepository_CreateTrader(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db)

	// Args
	type Args struct {
		Trader barter.Trader
	}
	var args Args
	_ = faker.FakeData(&args)

	trader, err := repo.CreateTrader(context.Background(), args.Trader)

	require.NoError(t, err)
	assertTrader(t, &args.Trader, trader)

	// No duplicate
	_, err = repo.CreateTrader(context.Background(), args.Trader)
	require.Error(t, err)
}

func TestPostgresRepository_GetTraderByEmail(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db, testdata.Path(testdata.TestDataTrader))

	email := "trader1@cresclab.com"

	_, err := repo.GetTraderByEmail(context.Background(), email)
	require.NoError(t, err)
}
