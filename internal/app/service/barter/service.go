package barter

import (
	"context"

	"github.com/rs/zerolog"
)

type BarterService struct {
	goodRepo GoodRepository
}

type BarterServiceParam struct {
	GoodRepo GoodRepository
}

func NewBarterService(_ context.Context, param BarterServiceParam) *BarterService {
	return &BarterService{
		goodRepo: param.GoodRepo,
	}
}

// logger wrap the execution context with component info
func (s *BarterService) logger(ctx context.Context) *zerolog.Logger {
	l := zerolog.Ctx(ctx).With().Str("component", "barter-service").Logger()
	return &l
}
