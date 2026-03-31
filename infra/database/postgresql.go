package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultConnectTimeout = 10 * time.Second

// DBFromEnv creates a PostgreSQL connection using GORM.
// Expected format: postgres://user:pass@host:port/db?sslmode=disable
func ConnectDB(ctx context.Context, dbURL string) (*gorm.DB, error) {
	if dbURL == "" {
		return nil, fmt.Errorf("dbURL is not set")
	}

	// Ensure DB connect attempts do not block indefinitely.
	dbURL = withConnectTimeout(dbURL, defaultConnectTimeout)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, defaultConnectTimeout)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return db, nil
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func withConnectTimeout(dbURL string, timeout time.Duration) string {
	u, err := url.Parse(dbURL)
	if err != nil {
		// If URL parsing fails, keep the original DSN.
		return dbURL
	}

	q := u.Query()
	if q.Get("connect_timeout") == "" {
		q.Set("connect_timeout", fmt.Sprintf("%d", int(timeout.Seconds())))
		u.RawQuery = q.Encode()
	}

	return u.String()
}
