package auth

import (
	"github.com/labstack/echo/v4"

	"go-ddd/adapter/httpapi/v1/auth/login"
	"go-ddd/adapter/httpapi/v1/auth/logout"
	"go-ddd/adapter/httpapi/v1/auth/refresh"
	"go-ddd/adapter/httpapi/v1/auth/register"
	loginUsecase "go-ddd/application/usecases/auth/login"
	logoutUsecase "go-ddd/application/usecases/auth/logout"
	refreshUsecase "go-ddd/application/usecases/auth/refresh"
	registerUsecase "go-ddd/application/usecases/auth/register"
)

// RegisterRoutesV1 wires v1 auth endpoints under /api/v1/auth.
func RegisterAuthRoutesV1(
	g *echo.Group,
	loginUsecase loginUsecase.Usecase,
	logoutUsecase logoutUsecase.Usecase,
	refreshUsecase refreshUsecase.Usecase,
	registerUsecase registerUsecase.Usecase,
) {
	grp := g.Group("/auth")

	grp.POST("/register", register.NewController(registerUsecase).Handle)
	grp.POST("/login", login.NewController(loginUsecase).Handle)
	grp.POST("/refresh", refresh.NewController(refreshUsecase).Handle)
	grp.POST("/logout", logout.NewController(logoutUsecase).Handle)
}
