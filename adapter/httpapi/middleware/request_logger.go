package middleware

import (
	"time"

	"github.com/labstack/echo/v4"

	"go-ddd/infra/logger"
)

// RequestLogger logs basic request info and duration.
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Now()

			// status code may not be set if the handler panicked; keep it best-effort.
			status := c.Response().Status
			logger.Log.Info().
				Str("method", c.Request().Method).
				Str("path", c.Path()).
				Int("status", status).
				Dur("duration", stop.Sub(start)).
				Msg("http request")

			return err
		}
	}
}

