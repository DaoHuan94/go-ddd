package refresh

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrRefreshTokenInvalid = errors.New("refresh token invalid")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
