package router

import (
	"github.com/gin-gonic/gin"

	"github.com/chatbotgang/go-clean-architecture-template/internal/app"
)

func RegisterHandlers(router *gin.Engine, app *app.Application) {
	registerAPIHandlers(router, app)
}

func registerAPIHandlers(router *gin.Engine, app *app.Application) {
	// Build middlewares
	BearerToken := NewAuthMiddlewareBearer(app)

	// We mount all handlers under /api path
	r := router.Group("/api")
	v1 := r.Group("/v1")

	// Add health-check
	v1.GET("/health", handlerHealthCheck())

	// Add auth namespace
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/traders", RegisterTrader(app))
		authGroup.POST("/traders/login", LoginTrader(app))
	}

	// Add barter namespace
	barterGroup := v1.Group("/barter", BearerToken.Required())
	{
		barterGroup.POST("/goods", PostGood(app))
		barterGroup.GET("/goods", ListMyGoods(app))
		barterGroup.GET("/goods/traders", ListOthersGoods(app))
		barterGroup.DELETE("/goods/:good_id", RemoveMyGood(app))
	}
}
