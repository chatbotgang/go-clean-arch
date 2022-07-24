package barter

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExchangeGoods(t *testing.T) {
	t.Parallel()
	// Args
	type Args struct {
		Trader1 Trader
		Good1   Good
		Trader2 Trader
		Good2   Good
	}
	var args Args
	_ = faker.FakeData(&args)
	args.Good1.OwnerID = args.Trader1.ID
	args.Good2.OwnerID = args.Trader2.ID

	goods := ExchangeGoods(args.Good1, args.Good2)

	require.Len(t, goods, 2)
	assert.True(t, goods[0].IsMyGood(args.Trader2))
	assert.True(t, goods[1].IsMyGood(args.Trader1))
}
