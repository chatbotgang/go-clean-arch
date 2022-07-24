package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/chatbotgang/go-clean-architecture-template/internal/app"
	"github.com/chatbotgang/go-clean-architecture-template/internal/app/service/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

func PostGood(app *app.Application) gin.HandlerFunc {
	type Body struct {
		Name string `json:"name" binding:"required"`
	}

	type Response struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		OwnerID   int       `json:"owner_id"`
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

		// Get current trader
		trader, err := GetCurrentTrader(c)
		if err != nil {
			respondWithError(c, err)
			return
		}

		// Invoke service
		good, err := app.BarterService.PostGood(ctx, barter.PostGoodParam{
			Trader:   *trader,
			GoodName: body.Name,
		})
		if err != nil {
			respondWithError(c, err)
			return
		}

		resp := Response{
			ID:        good.ID,
			Name:      good.Name,
			OwnerID:   good.OwnerID,
			CreatedAt: trader.CreatedAt,
		}
		respondWithJSON(c, http.StatusCreated, resp)
	}
}

func ListMyGoods(app *app.Application) gin.HandlerFunc {

	type Good struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		OwnerID   int       `json:"owner_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	type Response struct {
		Goods []Good `json:"goods"`
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get current trader
		trader, err := GetCurrentTrader(c)
		if err != nil {
			respondWithError(c, err)
			return
		}

		// Invoke service
		goods, err := app.BarterService.ListMyGoods(ctx, barter.ListMyGoodsParam{
			Trader: *trader,
		})
		if err != nil {
			respondWithError(c, err)
			return
		}

		resp := Response{
			Goods: []Good{},
		}
		for i := range goods {
			g := goods[i]
			resp.Goods = append(resp.Goods, Good(g))
		}

		respondWithJSON(c, http.StatusOK, resp)
	}
}

func ListOthersGoods(app *app.Application) gin.HandlerFunc {

	type Good struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		OwnerID   int       `json:"owner_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	type Response struct {
		Goods []Good `json:"goods"`
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get current trader
		trader, err := GetCurrentTrader(c)
		if err != nil {
			respondWithError(c, err)
			return
		}

		// Invoke service
		goods, err := app.BarterService.ListOthersGoods(ctx, barter.ListOthersGoodsParam{
			Trader: *trader,
		})
		if err != nil {
			respondWithError(c, err)
			return
		}

		resp := Response{
			Goods: []Good{},
		}
		for i := range goods {
			g := goods[i]
			resp.Goods = append(resp.Goods, Good(g))
		}

		respondWithJSON(c, http.StatusOK, resp)
	}
}

func RemoveMyGood(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Validate parameters
		goodID, err := GetParamInt(c, "good_id")
		if err != nil {
			respondWithError(c, err)
			return
		}

		// Get current trader
		trader, err := GetCurrentTrader(c)
		if err != nil {
			respondWithError(c, err)
			return
		}

		// Invoke service
		err = app.BarterService.RemoveMyGood(ctx, barter.RemoveGoodParam{
			Trader: *trader,
			GoodID: goodID,
		})
		if err != nil {
			respondWithError(c, err)
			return
		}

		respondWithoutBody(c, http.StatusNoContent)
	}
}
