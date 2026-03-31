package middleware

import (
	"runtime/debug"

	"github.com/labstack/echo/v4"

	"go-ddd/infra/logger"
)

// Recover prevents panics from crashing the server.
func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Log.Error().
						Interface("panic", r).
						Bytes("stack", debug.Stack()).
						Msg("panic recovered")

					_ = c.JSON(500, map[string]string{
						"message": "internal server error",
					})
				}
			}()

			return next(c)
		}
	}
}

