package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/chatbotgang/go-clean-architecture-template/internal/app"
	"github.com/chatbotgang/go-clean-architecture-template/internal/app/service/auth"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

func RegisterTrader(app *app.Application) gin.HandlerFunc {
	type Body struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	type Response struct {
		ID        int       `json:"id"`
		UID       string    `json:"uid"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Validate parameters
		var body Body
		err := c.ShouldBind(&body)
		if err != nil {
			respondWithError(c, common.NewError(common.ErrorCodeParameterInvalid, err, common.WithMsg("invalid parameter")))
			return
		}

		// Invoke service
		trader, err := app.AuthService.RegisterTrader(ctx, auth.RegisterTraderParam{
			Email:    body.Email,
			Name:     body.Name,
			Password: body.Password,
		})
		if err != nil {
			respondWithError(c, err)
			return
		}

		resp := Response{
			ID:        trader.ID,
			UID:       trader.UID,
			Email:     trader.Email,
			Name:      trader.Name,
			CreatedAt: trader.CreatedAt,
		}
		respondWithJSON(c, http.StatusCreated, resp)
	}
}

func LoginTrader(app *app.Application) gin.HandlerFunc {
	type Body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	type Response struct {
		TraderID int    `json:"trader_id"`
		Token    string `json:"token"`
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Validate parameters
		var body Body
		err := c.ShouldBind(&body)
		if err != nil {
			respondWithError(c, common.NewError(common.ErrorCodeParameterInvalid, err, common.WithMsg("invalid parameter")))
			return
		}

		// Invoke service
		trader, err := app.AuthService.LoginTrader(ctx, auth.LoginTraderParam{
			Email:    body.Email,
			Password: body.Password,
		})
		if err != nil {
			respondWithError(c, err)
			return
		}
		token, err := app.AuthService.GenerateTraderToken(ctx, *trader)
		if err != nil {
			respondWithError(c, err)
			return
		}

		resp := Response{
			TraderID: trader.ID,
			Token:    token,
		}
		respondWithJSON(c, http.StatusOK, resp)
	}
}
