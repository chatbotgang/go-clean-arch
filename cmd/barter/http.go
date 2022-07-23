package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/chatbotgang/go-clean-architecture-template/internal/app"
	"github.com/chatbotgang/go-clean-architecture-template/internal/router"
)

func runHTTPServer(rootCtx context.Context, wg *sync.WaitGroup, port int, app *app.Application) {
	// Set to release mode to disable Gin logger
	gin.SetMode(gin.ReleaseMode)

	// Create gin router
	ginRouter := gin.New()

	// Set general middleware
	router.SetGeneralMiddlewares(rootCtx, ginRouter)

	// Register all handlers
	router.RegisterHandlers(ginRouter, app)

	// Build HTTP server
	httpAddr := fmt.Sprintf("0.0.0.0:%d", port)
	server := &http.Server{
		Addr:    httpAddr,
		Handler: ginRouter,
	}

	// Run the server in a goroutine
	go func() {
		zerolog.Ctx(rootCtx).Info().Msgf("HTTP server is on http://%s", httpAddr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			zerolog.Ctx(rootCtx).Panic().Err(err).Str("addr", httpAddr).Msg("fail to start HTTP server")
		}
	}()

	// Wait for rootCtx done
	go func() {
		<-rootCtx.Done()

		// Graceful shutdown http server with a timeout
		zerolog.Ctx(rootCtx).Info().Msgf("HTTP server is closing")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("fail to shutdown HTTP server")
		}

		// Notify when server is closed
		zerolog.Ctx(rootCtx).Info().Msgf("HTTP server is closed")
		wg.Done()
	}()
}
