package router

import (
	"context"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// SetGeneralMiddlewares add general-purpose middlewares
func SetGeneralMiddlewares(ctx context.Context, ginRouter *gin.Engine) {
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(CORSMiddleware())
	ginRouter.Use(requestid.New())
	ginRouter.Use(LoggerMiddleware(ctx))

}
