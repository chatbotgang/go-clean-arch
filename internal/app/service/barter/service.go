package barter

import (
	"context"

	"github.com/rs/zerolog"
)

type BarterService struct {
}

type BarterServiceParam struct {
}

func NewBarterService(_ context.Context, param BarterServiceParam) *BarterService {
	return &BarterService{}
}

// logger wrap the execution context with component info
func (s *BarterService) logger(ctx context.Context) *zerolog.Logger {
	l := zerolog.Ctx(ctx).With().Str("component", "barter-service").Logger()
	return &l
}
