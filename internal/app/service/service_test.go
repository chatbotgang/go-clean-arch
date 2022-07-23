package chat

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
)

type serviceMock struct {
}

func buildServiceMock(ctrl *gomock.Controller) serviceMock {
	return serviceMock{}
}
func buildService(mock serviceMock) *Service {
	param := ServiceParam{}
	return NewService(context.Background(), param)
}

// nolint
func TestMain(m *testing.M) {
	// To avoid getting an empty object slice
	_ = faker.SetRandomMapAndSliceMinSize(2)

	// To avoid getting a zero random number
	_ = faker.SetRandomNumberBoundaries(1, 100)

	m.Run()
}
