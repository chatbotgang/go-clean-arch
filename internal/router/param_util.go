package router

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

const KeyAuth = "Authorization"
const KeyTrader = "Trader"

func GetAuthorizationToken(c *gin.Context) (string, common.Error) {
	token := c.GetHeader(KeyAuth)
	if token == "" {
		msg := "no Authorization"
		return "", common.NewError(common.ErrorCodeAuthNotAuthenticated, errors.New(msg), common.WithMsg(msg))
	}
	return token, nil
}

func GetCurrentTrader(c *gin.Context) (*barter.Trader, common.Error) {
	data, ok := c.Get(KeyTrader)
	if !ok {
		return nil, common.NewError(common.ErrorCodeAuthNotAuthenticated, errors.New("no credential"))
	}

	cred, ok := data.(barter.Trader)
	if !ok {
		return nil, common.NewError(common.ErrorCodeAuthNotAuthenticated, errors.New("failed to assert credential"))
	}
	return &cred, nil
}

func SetTrader(c *gin.Context, trader barter.Trader) common.Error {
	c.Set(KeyTrader, trader)
	return nil
}

// GetParamInt gets a key's value from Gin's URL param and transform it to int.
func GetParamInt(c *gin.Context, key string) (int, common.Error) {
	s := c.Param(key)
	if s == "" {
		msg := fmt.Sprintf("no %s", key)
		return 0, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg))
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		msg := fmt.Sprintf("invalid %s", key)
		return 0, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg))
	}
	return i, nil
}

// GetQueryInt gets a key's value from Gin's URL query and transform it to int.
func GetQueryInt(c *gin.Context, key string) (int, common.Error) {
	s := c.Query(key)
	if s == "" {
		msg := fmt.Sprintf("no %s", key)
		return 0, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg))
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		msg := fmt.Sprintf("invalid %s", key)
		return 0, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg))
	}
	return i, nil
}

// GetQueryBool gets a key's value from Gin's URL query and transform it to bool.
func GetQueryBool(c *gin.Context, key string) (bool, common.Error) {
	s := c.Query(key)
	if s == "" {
		msg := fmt.Sprintf("no %s", key)
		return false, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg))
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		msg := fmt.Sprintf("invalid %s", key)
		return false, common.NewError(common.ErrorCodeParameterInvalid, errors.New(msg), common.WithMsg(msg))
	}
	return b, nil
}
