package login

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	authusecase "go-ddd/application/usecases/auth"
	"go-ddd/adapter/httpapi/v1/auth/common"
)

type Controller struct {
	usecase authusecase.AuthUsecase
}

func NewController(usecase authusecase.AuthUsecase) *Controller {
	return &Controller{usecase: usecase}
}

func (ctl *Controller) Handle(c echo.Context) error {
	if ctl.usecase == nil {
		return c.JSON(http.StatusNotImplemented, common.ErrorResponse{Message: "not implemented"})
	}

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "invalid request"})
	}

	tokens, err := ctl.usecase.Login(context.Background(), authusecase.LoginArg{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return common.HandleAuthError(c, err)
	}

	return c.JSON(http.StatusOK, common.SuccessResponse[common.TokensResponseData]{
		Message: "success",
		Data: common.TokensResponseData{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	})
}

