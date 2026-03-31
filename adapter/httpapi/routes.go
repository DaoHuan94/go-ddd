package httpapi

import (
	"github.com/labstack/echo/v4"

	authv1 "go-ddd/adapter/httpapi/v1/auth"
	authusecase "go-ddd/application/usecases/auth"
)

// RegisterRoutes wires all HTTP routes (including versioned groups).
func RegisterRoutes(
	e *echo.Echo,
	authCtrl authusecase.AuthUsecase,
) {
	authv1.RegisterAuthRoutesV1(e.Group("/api/v1"), authCtrl)
}
