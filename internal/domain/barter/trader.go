package barter

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

type Trader struct {
	ID        int
	UID       string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTrader(uid string, email string, name string) Trader {
	return Trader{
		UID:   uid,
		Email: email,
		Name:  name,
	}
}

type traderClaim struct {
	jwt.StandardClaims
	Trader
}

func GenerateTraderToken(trader Trader, signingKey []byte, expiryDuration time.Duration, issuer string) (string, common.Error) {
	claim := &traderClaim{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiryDuration).Unix(),
			Issuer:    issuer,
			IssuedAt:  time.Now().Unix(),
		},
		trader,
	}

	// Generate Signed JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", common.NewError(common.ErrorCodeInternalProcess, err, common.WithMsg("failed to generate token"))
	}
	return signedToken, nil
}

func ParseTraderFromToken(signedToken string, signingKey []byte) (*Trader, common.Error) {
	token, err := jwt.ParseWithClaims(signedToken, &traderClaim{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok && e.Errors == jwt.ValidationErrorExpired {
			msg := "token is expired"
			return nil, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg))
		} else {
			return nil, common.NewError(common.ErrorCodeParameterInvalid, err, common.WithMsg("failed to parse token"))
		}
	}

	if !token.Valid {
		msg := "invalid token"
		return nil, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg))
	}

	claim, ok := token.Claims.(*traderClaim)
	if !ok {
		msg := "failed to parse claim"
		return nil, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg))
	}

	return &claim.Trader, nil
}
