package repository_impl

import (
	"context"
	"errors"
	"strings"
	"time"

	"go-ddd/domain/auth"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type userRow struct {
	ID           int64  `gorm:"column:id;primaryKey"`
	Email        string `gorm:"column:email;uniqueIndex"`
	PasswordHash string `gorm:"column:password_hash"`
	Name         string `gorm:"column:name"`
	AvatarURL    string `gorm:"column:avatar_url"`
}

type refreshTokenRow struct {
	ID        int64      `gorm:"column:id;primaryKey"`
	UserID    int64      `gorm:"column:user_id"`
	TokenHash string     `gorm:"column:token_hash;uniqueIndex"`
	ExpiresAt time.Time  `gorm:"column:expires_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
}

// PostgresAuthRepository stores auth users and refresh token state.
// For passwords, hashing is handled by application/usecases/auth.
type PostgresAuthRepository struct {
	db *gorm.DB
}

func NewPostgresAuthRepository(db *gorm.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}

func (r *PostgresAuthRepository) CreateUser(
	ctx context.Context,
	email string,
	passwordHash string,
	name string,
	avatarURL string,
) (auth.User, error) {
	row := userRow{
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		AvatarURL:    avatarURL,
	}

	if err := r.db.WithContext(ctx).
		Table("users").
		Create(&row).Error; err != nil {
		// Postgres unique violation message contains "duplicate key value" and the constraint/index.
		msg := err.Error()
		if strings.Contains(msg, "duplicate key value") && strings.Contains(strings.ToLower(msg), "email") {
			return auth.User{}, ErrEmailAlreadyExists
		}
		return auth.User{}, err
	}

	return auth.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Name:         row.Name,
		AvatarURL:    row.AvatarURL,
	}, nil
}

func (r *PostgresAuthRepository) GetUserByEmail(ctx context.Context, email string) (auth.User, error) {
	var row userRow
	err := r.db.WithContext(ctx).
		Table("users").
		First(&row, "email = ?", email).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return auth.User{}, ErrUserNotFound
		}
		return auth.User{}, err
	}

	return auth.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Name:         row.Name,
		AvatarURL:    row.AvatarURL,
	}, nil
}

func (r *PostgresAuthRepository) GetUserByID(ctx context.Context, userID int64) (auth.User, error) {
	var row userRow
	err := r.db.WithContext(ctx).
		Table("users").
		First(&row, "id = ?", userID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return auth.User{}, ErrUserNotFound
		}
		return auth.User{}, err
	}

	return auth.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Name:         row.Name,
		AvatarURL:    row.AvatarURL,
	}, nil
}

func (r *PostgresAuthRepository) CreateRefreshToken(
	ctx context.Context,
	userID int64,
	tokenHash string,
	expiresAt time.Time,
) error {
	row := refreshTokenRow{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	return r.db.WithContext(ctx).
		Table("refresh_tokens").
		Create(&row).Error
}

func (r *PostgresAuthRepository) GetRefreshToken(ctx context.Context, tokenHash string) (auth.RefreshToken, error) {
	var row refreshTokenRow
	err := r.db.WithContext(ctx).
		Table("refresh_tokens").
		First(&row, "token_hash = ?", tokenHash).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return auth.RefreshToken{}, ErrRefreshTokenNotFound
		}
		return auth.RefreshToken{}, err
	}

	return auth.RefreshToken{
		UserID:    row.UserID,
		ExpiresAt: row.ExpiresAt,
		RevokedAt: row.RevokedAt,
	}, nil
}

func (r *PostgresAuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string, revokedAt time.Time) error {
	res := r.db.WithContext(ctx).
		Table("refresh_tokens").
		Where("token_hash = ? AND revoked_at IS NULL", tokenHash).
		Updates(map[string]any{
			"revoked_at": revokedAt,
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}

