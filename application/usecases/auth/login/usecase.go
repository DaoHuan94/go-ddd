package login

import (
	"context"
	"time"

	"go-ddd/domain/auth"
	"go-ddd/infra/security"
)

type Usecase interface {
	Execute(ctx context.Context, arg LoginArg) (TokensResult, error)
}

type UsecaseImpl struct {
	repo       auth.AuthRepository
	accessCfg  security.AccessTokenConfig
	refreshTTL time.Duration
}

func NewUsecase(
	repo auth.AuthRepository,
	accessCfg security.AccessTokenConfig,
	refreshTTL time.Duration,
) Usecase {
	return &UsecaseImpl{
		repo:       repo,
		accessCfg:  accessCfg,
		refreshTTL: refreshTTL,
	}
}

func (u *UsecaseImpl) Execute(ctx context.Context, arg LoginArg) (TokensResult, error) {
	user, err := u.repo.GetUserByEmail(ctx, arg.Email)
	if err != nil {
		return TokensResult{}, ErrInvalidCredentials
	}

	if err := security.VerifyPassword(user.PasswordHash, arg.Password); err != nil {
		return TokensResult{}, ErrInvalidCredentials
	}

	accessToken, err := security.GenerateAccessToken(u.accessCfg, user.ID, user.Email)
	if err != nil {
		return TokensResult{}, err
	}

	refreshToken, refreshHash, err := security.GenerateRefreshToken()
	if err != nil {
		return TokensResult{}, err
	}

	expiresAt := time.Now().Add(u.refreshTTL)
	if err := u.repo.CreateRefreshToken(ctx, user.ID, refreshHash, expiresAt); err != nil {
		return TokensResult{}, err
	}

	return TokensResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

