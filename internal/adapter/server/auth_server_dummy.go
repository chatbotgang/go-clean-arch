package server

import (
	"context"

	"github.com/google/uuid"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

// We don't really implement the auth server, which communicates with the Crescendo's account management service,
// because this application is only used for introduction to the clean architecture.

type AuthServerParam struct{}

type AuthServer struct{}

func NewAuthServer(_ context.Context, param AuthServerParam) *AuthServer {
	return &AuthServer{}
}

func (s *AuthServer) RegisterAccount(ctx context.Context, email string, password string) (string, common.Error) {
	return uuid.NewString(), nil
}
