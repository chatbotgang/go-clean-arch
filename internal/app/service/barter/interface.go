package barter

import (
	"context"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

//go:generate mockgen -destination automock/auth_server.go -package=automock . AuthServer
type AuthServer interface {
	RegisterAccount(ctx context.Context, email string, password string) (string, common.Error)
	AuthenticateAccount(ctx context.Context, email string, password string) common.Error
}

//go:generate mockgen -destination automock/trader_repository.go -package=automock . TraderRepository
type TraderRepository interface {
	GetTraderByEmail(ctx context.Context, email string) (*barter.Trader, common.Error)
	CreateTrader(ctx context.Context, trader barter.Trader) (*barter.Trader, common.Error)
}
