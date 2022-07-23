package barter

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"

	"github.com/chatbotgang/go-clean-architecture-template/internal/app/service/barter/automock"
)

type barterServiceMock struct {
	AuthServer *automock.MockAuthServer
	TraderRepo *automock.MockTraderRepository
}

func buildBarterServiceMock(ctrl *gomock.Controller) barterServiceMock {
	return barterServiceMock{
		AuthServer: automock.NewMockAuthServer(ctrl),
		TraderRepo: automock.NewMockTraderRepository(ctrl),
	}
}
func buildBarterService(mock barterServiceMock) *BarterService {
	param := BarterServiceParam{
		AuthServer: mock.AuthServer,
		TraderRepo: mock.TraderRepo,
	}
	return NewBarterService(context.Background(), param)
}

// nolint
func TestMain(m *testing.M) {
	// To avoid getting an empty object slice
	_ = faker.SetRandomMapAndSliceMinSize(2)

	// To avoid getting a zero random number
	_ = faker.SetRandomNumberBoundaries(1, 100)

	m.Run()
}
