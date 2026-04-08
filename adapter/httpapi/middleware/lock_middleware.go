package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func LockMiddleware(rdb *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.Background()

			userID := c.Request().Header.Get("X-User-ID")
			action := c.Path() // or custom action name

			if userID == "" {
				return next(c)
			}

			lockKey := "lock:" + userID + ":" + action

			// Try acquire lock
			res, err := rdb.SetArgs(ctx, lockKey, "1", redis.SetArgs{
				Mode: "NX",
				TTL:  10 * time.Second,
			}).Result()
				if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "lock error")
			}
			ok := res == "OK"

			if !ok {
				return echo.NewHTTPError(http.StatusConflict, "Request already in progress")
			}

			defer rdb.Del(ctx, lockKey)

			return next(c)
		}
	}
}
