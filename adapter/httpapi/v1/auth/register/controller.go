package register

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"go-ddd/adapter/httpapi/v1/auth/common"
	registerUsecase "go-ddd/application/usecases/auth/register"
)

type Controller struct {
	usecase registerUsecase.Usecase
}

func NewController(
	usecase registerUsecase.Usecase) *Controller {
	return &Controller{usecase: usecase}
}

func (ctl *Controller) Handle(c echo.Context) error {
	if ctl.usecase == nil {
		return c.JSON(http.StatusNotImplemented, common.ErrorResponse{Message: "not implemented"})
	}

	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "invalid request"})
	}

	tokens, err := ctl.usecase.Execute(context.Background(), registerUsecase.RegisterArg{
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		AvatarURL: req.AvatarURL,
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
