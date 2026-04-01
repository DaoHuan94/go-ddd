package register

import (
	"context"
	"strings"
	"time"

	"go-ddd/domain/auth"
	"go-ddd/infra/security"
)

type Usecase interface {
	Execute(ctx context.Context, arg RegisterArg) (TokensResult, error)
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

func (u *UsecaseImpl) Execute(ctx context.Context, arg RegisterArg) (TokensResult, error) {
	// Pre-check to avoid relying on infra-specific duplicate errors.
	if arg.Email != "" {
		if _, err := u.repo.GetUserByEmail(ctx, arg.Email); err == nil {
			return TokensResult{}, ErrEmailAlreadyExists
		}
	}

	passwordHash, err := security.HashPassword(arg.Password)
	if err != nil {
		return TokensResult{}, err
	}

	user, err := u.repo.CreateUser(ctx, arg.Email, passwordHash, arg.Name, arg.AvatarURL)
	if err != nil {
		if isEmailAlreadyExistsErr(err) {
			return TokensResult{}, ErrEmailAlreadyExists
		}
		return TokensResult{}, err
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

func isEmailAlreadyExistsErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "email already exists") ||
		(strings.Contains(msg, "duplicate key value") && strings.Contains(msg, "email"))
}
