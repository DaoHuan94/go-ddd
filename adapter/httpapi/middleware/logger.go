package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func LoggerMiddleware(base zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			reqIDBytes := make([]byte, 16)
			_, _ = rand.Read(reqIDBytes)
			reqId := hex.EncodeToString(reqIDBytes)
			log := base.With().Str("request_id", reqId).
				Str("method", c.Request().Method).
				Str("path", c.Path()).
				Logger()

			c.Set("logger", log)

			log.Info().Msg("request started")

			err := next(c)
			if err != nil {
				log.Error().Err(err).Msg("request failed")
			}

			return err
		}
	}
}
