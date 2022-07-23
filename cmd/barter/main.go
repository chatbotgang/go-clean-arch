package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/chatbotgang/go-clean-architecture-template/internal/app"
)

var (
	AppName    = "crescendo-barter"
	AppVersion = "unknown_version"
	AppBuild   = "unknown_build"
)

const (
	defaultEnv                     = "staging"
	defaultLogLevel                = "info"
	defaultPort                    = "9000"
	defaultTokenSigningKey         = "cb-signing-key" // nolint
	defaultTokenExpiryDurationHour = "8"
	defaultTokenTokenIssuer        = "crescendo-barter"
)

type AppConfig struct {
	// General configuration
	Env      *string
	LogLevel *string

	// Database configuration
	DatabaseDSN *string

	// HTTP configuration
	Port *int

	// Token configuration
	TokenSigningKey         *string
	TokenExpiryDurationHour *int
	TokenIssuer             *string
}

func initAppConfig() AppConfig {
	// Setup basic application information
	app := kingpin.New(AppName, "The HTTP server").
		Version(fmt.Sprintf("version: %s, build: %s", AppVersion, AppBuild))

	var config AppConfig

	config.Env = app.
		Flag("env", "The running environment").
		Envar("CB_ENV").Default(defaultEnv).Enum("staging", "production")

	config.LogLevel = app.
		Flag("log_level", "Log filtering level").
		Envar("CB_LOG_LEVEL").Default(defaultLogLevel).Enum("error", "warn", "info", "debug", "disabled")

	config.Port = app.
		Flag("port", "The HTTP server port").
		Envar("CB_PORT").Default(defaultPort).Int()

	config.DatabaseDSN = app.
		Flag("database_dsn", "The database DSN").
		Envar("CB_DATABASE_DSN").Required().String()

	config.TokenSigningKey = app.
		Flag("token_signing_key", "Token signing key").
		Envar("CB_TOKEN_SIGNING_KEY").Default(defaultTokenSigningKey).String()
	config.TokenExpiryDurationHour = app.
		Flag("token_expiry_duration_hour", "Token expiry time").
		Envar("CB_TOKEN_EXPIRY_DURATION_HOUR").Default(defaultTokenExpiryDurationHour).Int()
	config.TokenIssuer = app.
		Flag("token_issuer", "Token issuer").
		Envar("CB_TOKEN_ISSUER").Default(defaultTokenTokenIssuer).String()

	kingpin.MustParse(app.Parse(os.Args[1:]))

	return config
}

func initRootLogger(levelStr, env string) zerolog.Logger {
	// Set global log level
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Set logger time format
	const rfc3339Micro = "2006-01-02T15:04:05.000000Z07:00"
	zerolog.TimeFieldFormat = rfc3339Micro

	serviceName := fmt.Sprintf("%s-%s", AppName, env)
	rootLogger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", serviceName).
		Logger()

	return rootLogger
}

func main() {
	// Setup app configuration
	cfg := initAppConfig()

	// Create root logger
	rootLogger := initRootLogger(*cfg.LogLevel, *cfg.Env)

	// Create root context
	rootCtx, rootCtxCancelFunc := context.WithCancel(context.Background())
	rootCtx = rootLogger.WithContext(rootCtx)

	rootLogger.Info().
		Str("version", AppVersion).
		Str("build", AppBuild).
		Msgf("Launching %s", AppName)

	wg := sync.WaitGroup{}
	// Create application
	app := app.MustNewApplication(rootCtx, &wg, app.ApplicationParams{
		Env:                 *cfg.Env,
		DatabaseDSN:         *cfg.DatabaseDSN,
		TokenSigningKey:     []byte(*cfg.TokenSigningKey),
		TokenExpiryDuration: time.Duration(*cfg.TokenExpiryDurationHour) * time.Hour,
		TokenIssuer:         *cfg.TokenIssuer,
	})

	// Run server
	wg.Add(1)
	runHTTPServer(rootCtx, &wg, *cfg.Port, app)

	// Listen to SIGTERM/SIGINT to close
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	<-gracefulStop
	rootCtxCancelFunc()

	// Wait for all services to close with a specific timeout
	var waitUntilDone = make(chan struct{})
	go func() {
		wg.Wait()
		close(waitUntilDone)
	}()
	select {
	case <-waitUntilDone:
		rootLogger.Info().Msg("success to close all services")
	case <-time.After(10 * time.Second):
		rootLogger.Err(context.DeadlineExceeded).Msg("fail to close all services")
	}
}
