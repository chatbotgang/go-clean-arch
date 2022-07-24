package app

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/chatbotgang/go-clean-architecture-template/internal/adapter/repository/postgres"
	"github.com/chatbotgang/go-clean-architecture-template/internal/adapter/server"
	"github.com/chatbotgang/go-clean-architecture-template/internal/app/service/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/app/service/token"
)

type Application struct {
	Params        ApplicationParams
	BarterService *barter.BarterService
	TokenService  *token.TokenService
}

type ApplicationParams struct {
	// General configuration
	Env string

	// Database parameters
	DatabaseDSN string

	// Token parameter
	TokenSigningKey     []byte
	TokenExpiryDuration time.Duration
	TokenIssuer         string
}

func MustNewApplication(ctx context.Context, wg *sync.WaitGroup, params ApplicationParams) *Application {
	app, err := NewApplication(ctx, wg, params)
	if err != nil {
		log.Panicf("fail to new application, err: %s", err.Error())
	}
	return app
}

func NewApplication(ctx context.Context, wg *sync.WaitGroup, params ApplicationParams) (*Application, error) {
	// Create repositories
	db := sqlx.MustOpen("postgres", params.DatabaseDSN)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	pgRepo := postgres.NewPostgresRepository(ctx, db)

	// Create servers
	authServer := server.NewAuthServer(ctx, server.AuthServerParam{})

	// Create application
	app := &Application{
		Params: params,
		BarterService: barter.NewBarterService(ctx, barter.BarterServiceParam{
			AuthServer: authServer,
			TraderRepo: pgRepo,
		}),
		TokenService: token.NewTokenService(ctx, token.TokenServiceParam{
			SigningKey:     params.TokenSigningKey,
			ExpiryDuration: params.TokenExpiryDuration,
			Issuer:         params.TokenIssuer,
		}),
	}

	return app, nil
}
