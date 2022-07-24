package barter

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"gotest.tools/assert"
)

func TestGood_IsMyGood(t *testing.T) {
	t.Parallel()
	// Args
	type Args struct {
		Trader Trader
		Good   Good
	}
	var args Args
	_ = faker.FakeData(&args)

	// Test cases
	testCases := []struct {
		Name         string
		SetupArgs    func(t *testing.T) Args
		ExpectResult bool
	}{
		{
			Name: "my good",
			SetupArgs: func(t *testing.T) Args {
				a := args
				a.Good.OwnerID = a.Trader.ID
				return a
			},
			ExpectResult: true,
		},
		{
			Name: "others' good",
			SetupArgs: func(t *testing.T) Args {
				a := args
				a.Good.OwnerID = a.Trader.ID + 1
				return a
			},
			ExpectResult: false,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			a := c.SetupArgs(t)
			ok := a.Good.IsMyGood(a.Trader)

			assert.Equal(t, c.ExpectResult, ok)
		})
	}

}
