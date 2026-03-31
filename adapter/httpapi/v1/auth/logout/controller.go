package logout

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"go-ddd/adapter/httpapi/v1/auth/common"
	authUsecase "go-ddd/application/usecases/auth"
)

type Controller struct {
	usecase authUsecase.AuthUsecase
}

func NewController(usecase authUsecase.AuthUsecase) *Controller {
	return &Controller{usecase: usecase}
}

func (ctl *Controller) Handle(c echo.Context) error {
	if ctl.usecase == nil {
		return c.JSON(http.StatusNotImplemented, common.ErrorResponse{Message: "not implemented"})
	}

	var req LogoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "invalid request"})
	}

	if err := ctl.usecase.Logout(context.Background(), authUsecase.LogoutArg{
		RefreshToken: req.RefreshToken,
	}); err != nil {
		return common.HandleAuthError(c, err)
	}

	return c.JSON(http.StatusOK, common.ErrorResponse{Message: "success"})
}
