package barter

import (
	"context"
	"errors"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

type RegisterTraderParam struct {
	Email    string
	Name     string
	Password string
}

func (s *BarterService) RegisterTrader(ctx context.Context, param RegisterTraderParam) (*barter.Trader, common.Error) {
	// Check the given trader email exist or not
	_, err := s.traderRepo.GetTraderByEmail(ctx, param.Email)
	if err == nil {
		msg := "email exists"
		s.logger(ctx).Error().Msg(msg)
		return nil, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg))
	}

	// If not existed:
	// 1. Register a new account to Crescendo's auth server.
	// 2. Create a trader in the application.
	uid, err := s.authServer.RegisterAccount(ctx, param.Email, param.Password)
	if err != nil {
		s.logger(ctx).Error().Err(err).Msg("failed to register account in Crescendo")
		return nil, err
	}

	trader, err := s.traderRepo.CreateTrader(ctx, barter.NewTrader(uid, param.Email, param.Name))
	if err != nil {
		s.logger(ctx).Error().Err(err).Msg("failed to register trader")
		return nil, err
	}

	return trader, nil
}

type LoginTraderParam struct {
	Email    string
	Password string
}

func (s *BarterService) LoginTrader(ctx context.Context, param LoginTraderParam) (*barter.Trader, common.Error) {
	// Check the given trader email exist or not
	trader, err := s.traderRepo.GetTraderByEmail(ctx, param.Email)
	if err != nil {
		s.logger(ctx).Error().Err(err).Msg("failed to get trader")
		msg := "email does not exist"
		return nil, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg))
	}

	// Authenticate the account
	err = s.authServer.AuthenticateAccount(ctx, param.Email, param.Password)
	if err != nil {
		s.logger(ctx).Error().Err(err).Msg("failed to authenticate account")
		msg := "invalid password"
		s.logger(ctx).Error().Msg(msg)
		return nil, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg))
	}

	return trader, nil
}
