package httpapi

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"go-ddd/adapter/httpapi/middleware"
	loginUsecase "go-ddd/application/usecases/auth/login"
	logoutUsecase "go-ddd/application/usecases/auth/logout"
	refreshUsecase "go-ddd/application/usecases/auth/refresh"
	registerUsecase "go-ddd/application/usecases/auth/register"
)

type Server struct {
	echo            *echo.Echo
	loginUsecase    loginUsecase.Usecase
	logoutUsecase   logoutUsecase.Usecase
	refreshUsecase  refreshUsecase.Usecase
	registerUsecase registerUsecase.Usecase
}

func NewServer(
	loginUsecase loginUsecase.Usecase,
	logoutUsecase logoutUsecase.Usecase,
	refreshUsecase refreshUsecase.Usecase,
	registerUsecase registerUsecase.Usecase,
) *Server {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())

	return &Server{
		echo:            e,
		loginUsecase:    loginUsecase,
		logoutUsecase:   logoutUsecase,
		refreshUsecase:  refreshUsecase,
		registerUsecase: registerUsecase,
	}
}

func (s *Server) Register() {
	RegisterRoutes(s.echo, s.loginUsecase, s.logoutUsecase, s.refreshUsecase, s.registerUsecase)
}

func (s *Server) Start(ctx context.Context) error {
	s.Register()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	// Stop server on context cancellation.
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.echo.Shutdown(shutdownCtx)
	}()

	return s.echo.Start(addr)
}
