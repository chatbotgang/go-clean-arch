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

func TestBarterService_ExchangeGoods(t *testing.T) {
	t.Parallel()
	// Args
	type Args struct {
		Trader      barter.Trader
		RequestGood barter.Good
		TargetGood  barter.Good
	}
	var args Args
	_ = faker.FakeData(&args)
	args.RequestGood.OwnerID = args.Trader.ID

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
			Name: "success",
			SetupService: func(t *testing.T) *BarterService {
				mock := buildServiceMock(ctrl)

				requestGood := args.RequestGood
				targetGood := args.TargetGood
				mock.GoodRepo.EXPECT().GetGoodByID(gomock.Any(), requestGood.ID).Return(&requestGood, nil)
				mock.GoodRepo.EXPECT().GetGoodByID(gomock.Any(), targetGood.ID).Return(&targetGood, nil)
				mock.GoodRepo.EXPECT().UpdateGoods(gomock.Any(), gomock.Any()).Return([]barter.Good{requestGood, targetGood}, nil)

				service := buildService(mock)
				return service
			},
			ExpectError: false,
		},
		{
			Name: "failed to exchange goods",
			SetupService: func(t *testing.T) *BarterService {
				mock := buildServiceMock(ctrl)

				requestGood := args.RequestGood
				targetGood := args.TargetGood
				mock.GoodRepo.EXPECT().GetGoodByID(gomock.Any(), requestGood.ID).Return(&requestGood, nil)
				mock.GoodRepo.EXPECT().GetGoodByID(gomock.Any(), targetGood.ID).Return(&targetGood, nil)
				mock.GoodRepo.EXPECT().UpdateGoods(gomock.Any(), gomock.Any()).Return(nil, common.DomainError{})

				service := buildService(mock)
				return service
			},
			ExpectError: true,
		},
		{
			Name: "failed to get target good",
			SetupService: func(t *testing.T) *BarterService {
				mock := buildServiceMock(ctrl)

				requestGood := args.RequestGood
				targetGood := args.TargetGood
				mock.GoodRepo.EXPECT().GetGoodByID(gomock.Any(), requestGood.ID).Return(&requestGood, nil)
				mock.GoodRepo.EXPECT().GetGoodByID(gomock.Any(), targetGood.ID).Return(nil, common.DomainError{})

				service := buildService(mock)
				return service
			},
			ExpectError: true,
		},
		{
			Name: "no ownership of request good",
			SetupService: func(t *testing.T) *BarterService {
				mock := buildServiceMock(ctrl)

				requestGood := args.RequestGood
				requestGood.OwnerID = requestGood.OwnerID + 1
				mock.GoodRepo.EXPECT().GetGoodByID(gomock.Any(), requestGood.ID).Return(&requestGood, nil)

				service := buildService(mock)
				return service
			},
			ExpectError: true,
		},
		{
			Name: "failed to get request good",
			SetupService: func(t *testing.T) *BarterService {
				mock := buildServiceMock(ctrl)

				requestGood := args.RequestGood
				mock.GoodRepo.EXPECT().GetGoodByID(gomock.Any(), requestGood.ID).Return(nil, common.DomainError{})

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
			param := ExchangeGoodsParam{
				Trader:        args.Trader,
				RequestGoodID: args.RequestGood.ID,
				TargetGoodID:  args.TargetGood.ID,
			}

			err := service.ExchangeGoods(context.Background(), param)

			if c.ExpectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
