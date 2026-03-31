package httpapi

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"go-ddd/adapter/httpapi/middleware"
	authusecase "go-ddd/application/usecases/auth"
)

type Server struct {
	echo     *echo.Echo
	authCtrl authusecase.AuthUsecase
}

func NewServer(
	authCtrl authusecase.AuthUsecase,
) *Server {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())

	return &Server{
		echo:     e,
		authCtrl: authCtrl,
	}
}

func (s *Server) Register() {
	RegisterRoutes(s.echo, s.authCtrl)
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
