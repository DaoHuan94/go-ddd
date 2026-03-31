package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenConfig struct {
	Secret    string
	TTL       time.Duration
	Issuer    string
	Audience  string
}

func GenerateAccessToken(cfg AccessTokenConfig, userID int64, email string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"iat":   now.Unix(),
		"exp":   now.Add(cfg.TTL).Unix(),
	}
	if cfg.Issuer != "" {
		claims["iss"] = cfg.Issuer
	}
	if cfg.Audience != "" {
		claims["aud"] = cfg.Audience
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(cfg.Secret))
}

// GenerateRefreshToken creates an opaque refresh token + its SHA256 hash
// (the hash is what should be stored in the DB).
func GenerateRefreshToken() (token string, tokenHash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	token = base64.RawURLEncoding.EncodeToString(b)

	sum := sha256.Sum256([]byte(token))
	tokenHash = hex.EncodeToString(sum[:])
	return token, tokenHash, nil
}

func VerifyAccessToken(cfg AccessTokenConfig, tokenString string) (userID int64, email string, ok bool) {
	parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil || parsed == nil || !parsed.Valid {
		return 0, "", false
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", false
	}

	sub, _ := claims["sub"].(float64) // jwt.MapClaims decodes numbers as float64
	email, _ = claims["email"].(string)
	return int64(sub), email, true
}

