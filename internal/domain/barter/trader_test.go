package barter

import (
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestTrader_Token(t *testing.T) {
	t.Parallel()
	// Args
	type Args struct {
		Trader         Trader
		ExpiryDuration time.Duration
		SigningKey     []byte
		Issuer         string
	}
	var args Args
	_ = faker.FakeData(&args)

	// Test cases
	testCases := []struct {
		Name                string
		SetupArgs           func(t *testing.T) Args
		ExpectGenerateError bool
		ExpectParseError    bool
	}{
		{
			Name: "success",
			SetupArgs: func(t *testing.T) Args {
				a := args
				a.ExpiryDuration = time.Hour

				return a
			},
			ExpectGenerateError: false,
			ExpectParseError:    false,
		},
		{
			Name: "expired token",
			SetupArgs: func(t *testing.T) Args {
				a := args
				a.ExpiryDuration = -1 * time.Hour

				return a
			},
			ExpectGenerateError: false,
			ExpectParseError:    true,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			a := c.SetupArgs(t)
			token, err := GenerateTraderToken(a.Trader, a.SigningKey, a.ExpiryDuration, a.Issuer)

			if c.ExpectGenerateError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				trader, err := ParseTraderFromToken(token, a.SigningKey)
				if c.ExpectParseError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, a.Trader.ID, trader.ID)
				}
			}
		})
	}

}
