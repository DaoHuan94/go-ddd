package main

import (
	"context"
	"errors"
	"fmt"
	"go-ddd/adapter/httpapi"
	"go-ddd/adapter/httpapi/middleware"
	loginUsecase "go-ddd/application/usecases/auth/login"
	logoutUsecase "go-ddd/application/usecases/auth/logout"
	refreshUsecase "go-ddd/application/usecases/auth/refresh"
	registerUsecase "go-ddd/application/usecases/auth/register"
	"go-ddd/infra/config"
	"go-ddd/infra/database"
	"go-ddd/infra/database/repository_impl"
	"go-ddd/infra/logger"
	"go-ddd/infra/redis"
	"go-ddd/infra/security"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const shutdownTimeout = 10 * time.Second

func main() {
	logger.Init()

	cfg := mustLoadConfig()
	db := mustConnectDB(cfg)
	redisClient := redis.NewRedisClient(context.Background(), *cfg)

	e := echo.New()
	e.Use(middleware.LoggerMiddleware(logger.Log))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.IdempotencyMiddleware(redisClient.Rdb))
	// e.Use(middleware.LockMiddleware(redisClient.Rdb))

	// Wire usecases + routes.
	accessSecret := getenvDefault("JWT_ACCESS_SECRET", "dev-access-secret-change-me")
	accessTTLSeconds := getenvIntDefault("JWT_ACCESS_TTL_SECONDS", 900)
	refreshTTLSeconds := getenvIntDefault("JWT_REFRESH_TTL_SECONDS", 2592000)

	authAccessCfg := security.AccessTokenConfig{
		Secret: accessSecret,
		TTL:    time.Duration(accessTTLSeconds) * time.Second,
		Issuer: cfg.App.Name,
	}

	authRepo := repository_impl.NewPostgresAuthRepository(db)
	loginUsecase := loginUsecase.NewUsecase(
		authRepo,
		authAccessCfg,
		time.Duration(refreshTTLSeconds)*time.Second,
	)
	logoutUsecase := logoutUsecase.NewUsecase(
		authRepo,
		authAccessCfg,
		time.Duration(refreshTTLSeconds)*time.Second,
	)
	refreshUsecase := refreshUsecase.NewUsecase(
		authRepo,
		authAccessCfg,
		time.Duration(refreshTTLSeconds)*time.Second,
	)
	registerUsecase := registerUsecase.NewUsecase(
		authRepo,
		authAccessCfg,
		time.Duration(refreshTTLSeconds)*time.Second,
	)

	httpapi.RegisterRoutes(e, loginUsecase, logoutUsecase, refreshUsecase, registerUsecase)

	serverErr := startHTTPServer(e, cfg.App.Port)
	if !waitForSignalOrServerError(serverErr) {
		return
	}
	shutdown(e, db)
}

func mustLoadConfig() *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	return cfg
}

func mustConnectDB(cfg *config.Config) *gorm.DB {
	if cfg.Database.DbURL == "" {
		return nil
	}
	db, err := database.ConnectDB(context.Background(), cfg.Database.DbURL)
	if err != nil {
		log.Fatal("Failed to connect DB: ", err)
	}
	return db
}

func getenvDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getenvIntDefault(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func startHTTPServer(e *echo.Echo, port int) chan error {
	serverErr := make(chan error, 1)
	serverAddr := fmt.Sprintf(":%d", port)
	go func() {
		serverErr <- e.Start(serverAddr)
	}()
	return serverErr
}

func waitForSignalOrServerError(serverErr <-chan error) bool {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(quit)

	select {
	case sig := <-quit:
		logger.Log.Info().Str("signal", sig.String()).Msg("shutdown signal received")
		return true
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal().Err(err).Msg("server start failed")
		}
		return false
	}
}

func shutdown(e *echo.Echo, db *gorm.DB) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("failed to shutdown http server gracefully")
	}

	if db != nil {
		if err := database.CloseDB(db); err != nil {
			logger.Log.Error().Err(err).Msg("failed to close db connection")
		}
	}

	logger.Log.Info().Msg("graceful shutdown completed")
}
