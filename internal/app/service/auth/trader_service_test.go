package auth

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

func TestBarterService_RegisterTrader(t *testing.T) {
	t.Parallel()
	// Args
	type Args struct {
		Trader   barter.Trader
		Password string
	}
	var args Args
	_ = faker.FakeData(&args)

	// Init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test cases
	testCases := []struct {
		Name         string
		SetupService func(t *testing.T) *AuthService
		ExpectError  bool
	}{
		{
			Name: "trader does not exist",
			SetupService: func(t *testing.T) *AuthService {
				mock := buildServiceMock(ctrl)

				mock.TraderRepo.EXPECT().GetTraderByEmail(gomock.Any(), args.Trader.Email).Return(nil, common.DomainError{})
				mock.AuthServer.EXPECT().RegisterAccount(gomock.Any(), args.Trader.Email, args.Password).Return(args.Trader.UID, nil)
				mock.TraderRepo.EXPECT().CreateTrader(gomock.Any(), gomock.Any()).Return(&args.Trader, nil)

				service := buildService(mock)
				return service
			},
			ExpectError: false,
		},
		{
			Name: "failed to register trader",
			SetupService: func(t *testing.T) *AuthService {
				mock := buildServiceMock(ctrl)

				mock.TraderRepo.EXPECT().GetTraderByEmail(gomock.Any(), args.Trader.Email).Return(nil, common.DomainError{})
				mock.AuthServer.EXPECT().RegisterAccount(gomock.Any(), args.Trader.Email, args.Password).Return(args.Trader.UID, nil)
				mock.TraderRepo.EXPECT().CreateTrader(gomock.Any(), gomock.Any()).Return(nil, common.DomainError{})

				service := buildService(mock)
				return service
			},
			ExpectError: true,
		},
		{
			Name: "failed to register account",
			SetupService: func(t *testing.T) *AuthService {
				mock := buildServiceMock(ctrl)

				mock.TraderRepo.EXPECT().GetTraderByEmail(gomock.Any(), args.Trader.Email).Return(nil, common.DomainError{})
				mock.AuthServer.EXPECT().RegisterAccount(gomock.Any(), args.Trader.Email, args.Password).Return("", common.DomainError{})

				service := buildService(mock)
				return service
			},
			ExpectError: true,
		},
		{
			Name: "trader exist",
			SetupService: func(t *testing.T) *AuthService {
				mock := buildServiceMock(ctrl)

				mock.TraderRepo.EXPECT().GetTraderByEmail(gomock.Any(), args.Trader.Email).Return(&args.Trader, nil)

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
			param := RegisterTraderParam{
				Email:    args.Trader.Email,
				Name:     args.Trader.Name,
				Password: args.Password,
			}

			_, err := service.RegisterTrader(context.Background(), param)

			if c.ExpectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBarterService_LoginTrader(t *testing.T) {
	t.Parallel()
	// Args
	type Args struct {
		Trader   barter.Trader
		Password string
	}
	var args Args
	_ = faker.FakeData(&args)

	// Init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test cases
	testCases := []struct {
		Name         string
		SetupService func(t *testing.T) *AuthService
		ExpectError  bool
	}{
		{
			Name: "success",
			SetupService: func(t *testing.T) *AuthService {
				mock := buildServiceMock(ctrl)

				mock.TraderRepo.EXPECT().GetTraderByEmail(gomock.Any(), args.Trader.Email).Return(&args.Trader, nil)
				mock.AuthServer.EXPECT().AuthenticateAccount(gomock.Any(), args.Trader.Email, args.Password).Return(nil)

				service := buildService(mock)
				return service
			},
			ExpectError: false,
		},
		{
			Name: "invalid args.Password",
			SetupService: func(t *testing.T) *AuthService {
				mock := buildServiceMock(ctrl)

				mock.TraderRepo.EXPECT().GetTraderByEmail(gomock.Any(), args.Trader.Email).Return(&args.Trader, nil)
				mock.AuthServer.EXPECT().AuthenticateAccount(gomock.Any(), args.Trader.Email, args.Password).Return(common.DomainError{})

				service := buildService(mock)
				return service
			},
			ExpectError: true,
		},
		{
			Name: "email does not exist",
			SetupService: func(t *testing.T) *AuthService {
				mock := buildServiceMock(ctrl)

				mock.TraderRepo.EXPECT().GetTraderByEmail(gomock.Any(), args.Trader.Email).Return(nil, common.DomainError{})

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
			param := LoginTraderParam{
				Email:    args.Trader.Email,
				Password: args.Password,
			}

			_, err := service.LoginTrader(context.Background(), param)

			if c.ExpectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
