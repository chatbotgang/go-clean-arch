package chat

import (
	"context"

	"github.com/rs/zerolog"
)

type Service struct {
}

type ServiceParam struct {
}

func NewService(_ context.Context, param ServiceParam) *Service {
	return &Service{}
}

// logger wrap the execution context with component info
func (s *Service) logger(ctx context.Context) *zerolog.Logger {
	l := zerolog.Ctx(ctx).With().Str("component", "service").Logger()
	return &l
}
