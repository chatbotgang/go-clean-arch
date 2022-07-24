package barter

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

func TestBarterService_RemoveMyGood(t *testing.T) {
	t.Parallel()
	// Args
	type Args struct {
		Trader barter.Trader
		Good   barter.Good
	}
	var args Args
	_ = faker.FakeData(&args)

	// Init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test cases
	testCases := []struct {
		Name         string
		SetupService func(t *testing.T) *BarterService
		ExpectError  bool
	}{
		{
			Name: "remove my good",
			SetupService: func(t *testing.T) *BarterService {
				mock := buildServiceMock(ctrl)

				trader := args.Trader
				good := args.Good
				good.OwnerID = trader.ID
				mock.GoodRepo.EXPECT().GetGoodByID(gomock.Any(), good.ID).Return(&good, nil)
				mock.GoodRepo.EXPECT().DeleteGoodByID(gomock.Any(), good.ID).Return(nil)

				service := buildService(mock)
				return service
			},
			ExpectError: false,
		},
		{
			Name: "remove others' good",
			SetupService: func(t *testing.T) *BarterService {
				mock := buildServiceMock(ctrl)

				trader := args.Trader
				good := args.Good
				good.OwnerID = trader.ID + 1
				mock.GoodRepo.EXPECT().GetGoodByID(gomock.Any(), good.ID).Return(&good, nil)

				service := buildService(mock)
				return service
			},
			ExpectError: true,
		},
		{
			Name: "good not found",
			SetupService: func(t *testing.T) *BarterService {
				mock := buildServiceMock(ctrl)

				good := args.Good
				mock.GoodRepo.EXPECT().GetGoodByID(gomock.Any(), good.ID).Return(nil, common.BaseError{})

				service := buildService(mock)
				return service
			},
			ExpectError: true,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			service := c.SetupService(t)
			param := RemoveGoodParam{
				Trader: args.Trader,
				GoodID: args.Good.ID,
			}

			err := service.RemoveMyGood(context.Background(), param)

			if c.ExpectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
