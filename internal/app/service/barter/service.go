package barter

import (
	"context"

	"github.com/rs/zerolog"
)

type BarterService struct {
	authServer AuthServer
	traderRepo TraderRepository
}

type BarterServiceParam struct {
	AuthServer AuthServer
	TraderRepo TraderRepository
}

func NewBarterService(_ context.Context, param BarterServiceParam) *BarterService {
	return &BarterService{
		authServer: param.AuthServer,
		traderRepo: param.TraderRepo,
	}
}

// logger wrap the execution context with component info
func (s *BarterService) logger(ctx context.Context) *zerolog.Logger {
	l := zerolog.Ctx(ctx).With().Str("component", "barter-service").Logger()
	return &l
}
