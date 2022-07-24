package router

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/chatbotgang/go-clean-architecture-template/internal/app"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

type AuthMiddlewareBearer struct {
	app *app.Application
}

func NewAuthMiddlewareBearer(app *app.Application) *AuthMiddlewareBearer {
	return &AuthMiddlewareBearer{
		app: app,
	}
}

func (m *AuthMiddlewareBearer) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get Bearer token
		token, err := GetAuthorizationToken(c)
		if err != nil {
			respondWithError(c, err)
			return
		}
		tokens := strings.Split(token, "Bearer ")
		if len(tokens) != 2 {
			msg := "bearer token is needed"
			respondWithError(c, common.NewError(common.ErrorCodeAuthNotAuthenticated, errors.New(msg), common.WithMsg(msg)))
			return
		}

		// Validate token
		trader, err := m.app.AuthService.ValidateTraderToken(ctx, tokens[1])
		if err != nil {
			respondWithError(c, common.NewError(common.ErrorCodeAuthNotAuthenticated, errors.New(err.Error()), common.WithMsg(err.ClientMsg())))
			return
		}

		// Set trader to context
		if err = SetTrader(c, *trader); err != nil {
			respondWithError(c, err)
			return
		}
		c.Next()
	}
}
