package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chatbotgang/go-clean-architecture-template/internal/app"
	"github.com/chatbotgang/go-clean-architecture-template/internal/app/service/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

func ExchangeGoods(app *app.Application) gin.HandlerFunc {
	type Body struct {
		RequestGoodID int `json:"request_good_id" binding:"required"`
		TargetGoodID  int `json:"target_good_id" binding:"required"`
	}

	type Response struct {
		UUID string `json:"uuid"`
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

		uid := c.GetHeader("X-Idempotency-Key")
		if uid == "" {
			msg := "no idempotency key"
			respondWithError(c, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg)))
			return
		}

		// Get current trader
		trader, err := GetCurrentTrader(c)
		if err != nil {
			respondWithError(c, err)
			return
		}

		// Invoke service
		err = app.BarterService.ExchangeGoods(ctx, barter.ExchangeGoodsParam{
			Trader:        *trader,
			RequestGoodID: body.RequestGoodID,
			TargetGoodID:  body.TargetGoodID,
		})
		if err != nil {
			respondWithError(c, err)
			return
		}

		resp := Response{
			UUID: uid,
		}
		respondWithJSON(c, http.StatusOK, resp)
	}
}
