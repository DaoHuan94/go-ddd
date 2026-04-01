package common

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	loginUsecase "go-ddd/application/usecases/auth/login"
	refreshUsecase "go-ddd/application/usecases/auth/refresh"
	registerUsecase "go-ddd/application/usecases/auth/register"
)

type TokensResponseData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SuccessResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func HandleAuthError(c echo.Context, err error) error {
	// Map usecase errors to HTTP status codes.
	switch {
	case errors.Is(err, registerUsecase.ErrEmailAlreadyExists):
		return c.JSON(http.StatusConflict, ErrorResponse{Message: err.Error()})
	case errors.Is(err, loginUsecase.ErrInvalidCredentials):
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
	case errors.Is(err, refreshUsecase.ErrRefreshTokenExpired):
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
	case errors.Is(err, refreshUsecase.ErrRefreshTokenInvalid):
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
	default:
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "internal server error"})
	}
}
