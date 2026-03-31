package auth

import (
	"context"
	"time"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Name         string
	AvatarURL    string
}

type RefreshToken struct {
	UserID    int64
	ExpiresAt time.Time
	RevokedAt *time.Time
}

type AuthRepository interface {
	// User operations
	CreateUser(ctx context.Context, email string, passwordHash string, name string, avatarURL string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, userID int64) (User, error)

	// Refresh token operations
	CreateRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, tokenHash string) (RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string, revokedAt time.Time) error
}

