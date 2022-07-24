package barter

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"

	"github.com/chatbotgang/go-clean-architecture-template/internal/app/service/barter/automock"
)

type serviceMock struct {
	GoodRepo *automock.MockGoodRepository
}

func buildServiceMock(ctrl *gomock.Controller) serviceMock {
	return serviceMock{
		GoodRepo: automock.NewMockGoodRepository(ctrl),
	}
}
func buildService(mock serviceMock) *BarterService {
	param := BarterServiceParam{
		GoodRepo: mock.GoodRepo,
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
