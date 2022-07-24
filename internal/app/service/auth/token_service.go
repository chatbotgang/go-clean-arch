package auth

import (
	"context"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

func (s *AuthService) GenerateTraderToken(_ context.Context, trader barter.Trader) (string, common.Error) {
	signedToken, err := barter.GenerateTraderToken(trader, s.signingKey, s.expiryDuration, s.issuer)
	if err != nil {
		return "", common.NewError(common.ErrorCodeParameterInvalid, err, common.WithMsg(err.ClientMsg()))
	}
	return signedToken, nil
}

func (s *AuthService) ValidateTraderToken(_ context.Context, signedToken string) (*barter.Trader, common.Error) {
	trader, err := barter.ParseTraderFromToken(signedToken, s.signingKey)
	if err != nil {
		return nil, common.NewError(common.ErrorCodeAuthNotAuthenticated, err, common.WithMsg(err.ClientMsg()))
	}
	return trader, nil
}
