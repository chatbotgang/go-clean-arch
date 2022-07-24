package barter

import (
	"context"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

type ExchangeGoodsParam struct {
	Trader        barter.Trader
	RequestGoodID int
	TargetGoodID  int
}

func (s *BarterService) ExchangeGoods(ctx context.Context, param ExchangeGoodsParam) common.Error {
	// 1. Check ownership of request Good
	requestGood, err := s.goodRepo.GetGoodByID(ctx, param.RequestGoodID)
	if err != nil {
		s.logger(ctx).Error().Err(err).Msg("failed to get request good")
		return err
	}
	if !requestGood.IsMyGood(param.Trader) {
		s.logger(ctx).Error().Msg("not the owner of request good")
		return common.NewError(common.ErrorCodeAuthPermissionDenied, nil)
	}

	// 2. Check the target good exist or not
	targetGood, err := s.goodRepo.GetGoodByID(ctx, param.TargetGoodID)
	if err != nil {
		s.logger(ctx).Error().Err(err).Msg("failed to get target good")
		return err
	}

	// 3. Exchange ownerships of two goods
	_, err = s.goodRepo.UpdateGoods(ctx, barter.ExchangeGoods(*requestGood, *targetGood))
	if err != nil {
		s.logger(ctx).Error().Err(err).Msg("failed to exchange goods")
		return err
	}

	return nil
}
