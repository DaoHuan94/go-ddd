package common

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	authUsecase "go-ddd/application/usecases/auth"
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
	case errors.Is(err, authUsecase.ErrEmailAlreadyExists):
		return c.JSON(http.StatusConflict, ErrorResponse{Message: err.Error()})
	case errors.Is(err, authUsecase.ErrInvalidCredentials):
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
	case errors.Is(err, authUsecase.ErrRefreshTokenExpired):
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
	case errors.Is(err, authUsecase.ErrRefreshTokenInvalid):
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
	default:
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "internal server error"})
	}
}
