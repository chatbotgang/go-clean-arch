package barter

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
)

type barterServiceMock struct {
}

func buildBarterServiceMock(ctrl *gomock.Controller) barterServiceMock {
	return barterServiceMock{}
}
func buildService(mock barterServiceMock) *BarterService {
	param := BarterServiceParam{}
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
