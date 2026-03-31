package logout

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"go-ddd/domain/auth"
	"go-ddd/infra/security"
)

type Usecase interface {
	Execute(ctx context.Context, arg LogoutArg) error
}

type UsecaseImpl struct {
	repo       auth.AuthRepository
	refreshTTL time.Duration
	accessCfg  security.AccessTokenConfig
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

func (u *UsecaseImpl) Execute(ctx context.Context, arg LogoutArg) error {
	if arg.RefreshToken == "" {
		return nil
	}

	now := time.Now()
	tokenHash := refreshTokenHash(arg.RefreshToken)
	_ = u.repo.RevokeRefreshToken(ctx, tokenHash, now)
	return nil
}

func refreshTokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

