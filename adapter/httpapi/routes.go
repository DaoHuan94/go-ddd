package httpapi

import (
	"github.com/labstack/echo/v4"

	authv1 "go-ddd/adapter/httpapi/v1/auth"
	loginUsecase "go-ddd/application/usecases/auth/login"
	logoutUsecase "go-ddd/application/usecases/auth/logout"
	refreshUsecase "go-ddd/application/usecases/auth/refresh"
	registerUsecase "go-ddd/application/usecases/auth/register"
)

// RegisterRoutes wires all HTTP routes (including versioned groups).
func RegisterRoutes(
	e *echo.Echo,
	loginUsecase loginUsecase.Usecase,
	logoutUsecase logoutUsecase.Usecase,
	refreshUsecase refreshUsecase.Usecase,
	registerUsecase registerUsecase.Usecase,
) {
	authv1.RegisterAuthRoutesV1(e.Group("/api/v1"), loginUsecase, logoutUsecase, refreshUsecase, registerUsecase)
}
