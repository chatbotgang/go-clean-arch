package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlerHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}
