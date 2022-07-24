package auth

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type AuthService struct {
	authServer AuthServer
	traderRepo TraderRepository

	signingKey     []byte
	expiryDuration time.Duration
	issuer         string
}

type AuthServiceParam struct {
	AuthServer AuthServer
	TraderRepo TraderRepository

	SigningKey     []byte
	ExpiryDuration time.Duration
	Issuer         string
}

func NewAuthService(_ context.Context, param AuthServiceParam) *AuthService {
	return &AuthService{
		authServer: param.AuthServer,
		traderRepo: param.TraderRepo,

		signingKey:     param.SigningKey,
		expiryDuration: param.ExpiryDuration,
		issuer:         param.Issuer,
	}
}

// logger wrap the execution context with component info
func (s *AuthService) logger(ctx context.Context) *zerolog.Logger {
	l := zerolog.Ctx(ctx).With().Str("component", "auth-service").Logger()
	return &l
}
