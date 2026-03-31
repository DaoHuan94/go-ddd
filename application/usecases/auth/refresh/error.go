package refresh

import authtypes "go-ddd/application/usecases/auth/types"

var (
	ErrInvalidCredentials  = authtypes.ErrInvalidCredentials
	ErrRefreshTokenInvalid = authtypes.ErrRefreshTokenInvalid
	ErrRefreshTokenExpired = authtypes.ErrRefreshTokenExpired
)

