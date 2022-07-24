package barter

import (
	"context"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

//go:generate mockgen -destination automock/good_repository.go -package=automock . GoodRepository
type GoodRepository interface {
	CreateGood(ctx context.Context, param barter.Good) (*barter.Good, common.Error)
	GetGoodByID(ctx context.Context, id int) (*barter.Good, common.Error)
	ListGoods(ctx context.Context) ([]barter.Good, common.Error)
	ListGoodsByOwner(ctx context.Context, ownerID int) ([]barter.Good, common.Error)
	UpdateGood(ctx context.Context, good barter.Good) (*barter.Good, common.Error)
	DeleteGoodByID(ctx context.Context, id int) common.Error
}
