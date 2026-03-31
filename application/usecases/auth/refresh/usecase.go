package refresh

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"go-ddd/domain/auth"
	"go-ddd/infra/security"
)

type Usecase interface {
	Execute(ctx context.Context, arg RefreshArg) (TokensResult, error)
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

func (u *UsecaseImpl) Execute(ctx context.Context, arg RefreshArg) (TokensResult, error) {
	if arg.RefreshToken == "" {
		return TokensResult{}, ErrRefreshTokenInvalid
	}

	tokenHash := refreshTokenHash(arg.RefreshToken)
	rt, err := u.repo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return TokensResult{}, ErrRefreshTokenInvalid
	}

	now := time.Now()
	if rt.RevokedAt != nil {
		return TokensResult{}, ErrRefreshTokenInvalid
	}
	if rt.ExpiresAt.Before(now) {
		return TokensResult{}, ErrRefreshTokenExpired
	}

	// Rotate: revoke old refresh token
	_ = u.repo.RevokeRefreshToken(ctx, tokenHash, now)

	user, err := u.repo.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return TokensResult{}, ErrRefreshTokenInvalid
	}

	// Issue new tokens
	return issueTokens(ctx, u.repo, u.accessCfg, u.refreshTTL, user.ID, user.Email)
}

func issueTokens(
	ctx context.Context,
	repo auth.AuthRepository,
	accessCfg security.AccessTokenConfig,
	refreshTTL time.Duration,
	userID int64,
	email string,
) (TokensResult, error) {
	accessToken, err := security.GenerateAccessToken(accessCfg, userID, email)
	if err != nil {
		return TokensResult{}, err
	}

	refreshToken, refreshHash, err := security.GenerateRefreshToken()
	if err != nil {
		return TokensResult{}, err
	}

	expiresAt := time.Now().Add(refreshTTL)
	if err := repo.CreateRefreshToken(ctx, userID, refreshHash, expiresAt); err != nil {
		return TokensResult{}, err
	}

	return TokensResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func refreshTokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

