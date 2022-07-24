package auth

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"

	"github.com/chatbotgang/go-clean-architecture-template/internal/app/service/auth/automock"
)

type serviceMock struct {
	AuthServer *automock.MockAuthServer
	TraderRepo *automock.MockTraderRepository
}

func buildServiceMock(ctrl *gomock.Controller) serviceMock {
	return serviceMock{
		AuthServer: automock.NewMockAuthServer(ctrl),
		TraderRepo: automock.NewMockTraderRepository(ctrl),
	}
}
func buildService(mock serviceMock) *AuthService {
	param := AuthServiceParam{
		AuthServer: mock.AuthServer,
		TraderRepo: mock.TraderRepo,
	}
	return NewAuthService(context.Background(), param)
}

// nolint
func TestMain(m *testing.M) {
	// To avoid getting an empty object slice
	_ = faker.SetRandomMapAndSliceMinSize(2)

	// To avoid getting a zero random number
	_ = faker.SetRandomNumberBoundaries(1, 100)

	m.Run()
}
