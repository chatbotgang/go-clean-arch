package barter

import (
	"context"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

type PostGoodParam struct {
	Trader   barter.Trader
	GoodName string
}

func (s *BarterService) PostGood(ctx context.Context, param PostGoodParam) (*barter.Good, common.Error) {
	return s.goodRepo.CreateGood(ctx, barter.NewGood(param.Trader, param.GoodName))
}

type ListMyGoodsParam struct {
	Trader barter.Trader
}

func (s *BarterService) ListMyGoods(ctx context.Context, param ListMyGoodsParam) ([]barter.Good, common.Error) {
	goods, err := s.goodRepo.ListGoodsByOwner(ctx, param.Trader.ID)
	if err != nil {
		return nil, err
	}
	return goods, nil
}

type ListOthersGoodsParam struct {
	Trader barter.Trader
}

func (s *BarterService) ListOthersGoods(ctx context.Context, param ListOthersGoodsParam) ([]barter.Good, common.Error) {
	goods, err := s.goodRepo.ListGoods(ctx)
	if err != nil {
		return nil, err
	}

	// Filter out goods of mine
	var filteredGoods []barter.Good
	for i := range goods {
		g := goods[i]
		if !g.MyGood(param.Trader) {
			filteredGoods = append(filteredGoods, g)
		}
	}
	return filteredGoods, nil
}

type RemoveGoodParam struct {
	Trader barter.Trader
	GoodID int
}

func (s *BarterService) RemoveMyGood(ctx context.Context, param RemoveGoodParam) common.Error {
	// Check the good exist otr not
	good, err := s.goodRepo.GetGoodByID(ctx, param.GoodID)
	if err != nil {
		s.logger(ctx).Error().Err(err).Msg("failed to get good")
		return err
	}

	// Check the ownership
	if !good.MyGood(param.Trader) {
		s.logger(ctx).Error().
			Int("traderID", param.Trader.ID).
			Int("goodOwnerID", good.OwnerID).
			Msg("cannot remove others' good")
		return common.NewError(common.ErrorCodeAuthPermissionDenied, nil)
	}

	// Remove the good
	return s.goodRepo.DeleteGoodByID(ctx, good.ID)
}
