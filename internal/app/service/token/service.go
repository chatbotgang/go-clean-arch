package token

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

type TokenService struct {
	signingKey     []byte
	expiryDuration time.Duration
	issuer         string
}

type TokenServiceParam struct {
	SigningKey     []byte
	ExpiryDuration time.Duration
	Issuer         string
}

type CustomClaims struct {
	jwt.StandardClaims
	barter.Trader
}

func NewTokenService(_ context.Context, param TokenServiceParam) *TokenService {
	return &TokenService{
		signingKey:     param.SigningKey,
		expiryDuration: param.ExpiryDuration,
		issuer:         param.Issuer,
	}
}

func (s *TokenService) GenerateTraderToken(_ context.Context, trader barter.Trader) (string, common.Error) {
	signedToken, err := barter.GenerateTraderToken(trader, s.signingKey, s.expiryDuration, s.issuer)
	if err != nil {
		return "", common.NewError(common.ErrorCodeParameterInvalid, err, common.WithMsg(err.ClientMsg()))
	}
	return signedToken, nil
}

func (s *TokenService) ValidateTraderToken(_ context.Context, signedToken string) (*barter.Trader, common.Error) {
	trader, err := barter.ParseTraderFromToken(signedToken, s.signingKey)
	if err != nil {
		return nil, common.NewError(common.ErrorCodeAuthNotAuthenticated, err, common.WithMsg(err.ClientMsg()))
	}
	return trader, nil
}
