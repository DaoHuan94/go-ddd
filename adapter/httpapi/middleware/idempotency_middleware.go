package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type CachedResponse struct {
	Status int             `json:"status"`
	Body   json.RawMessage `json:"body"`
}

var idemScript = redis.NewScript(`
-- KEYS[1] = idempotency key
-- KEYS[2] = lock key

-- ARGV[1] = PROCESSING
-- ARGV[2] = lock value
-- ARGV[3] = TTL

local val = redis.call("GET", KEYS[1])

-- ✅ Already finished
if val and val ~= ARGV[1] then
    return { "DONE", val }
end

-- ⏳ Still processing
if val == ARGV[1] then
    return { "PROCESSING", nil }
end

-- 🔒 Try lock
local lock = redis.call("SETNX", KEYS[2], ARGV[2])

if lock == 1 then
    redis.call("EXPIRE", KEYS[2], ARGV[3])
    redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[3])
    return { "LOCKED", nil }
end

return { "RETRY", nil }

`)

func IdempotencyMiddleware(rdb *redis.Client) echo.MiddlewareFunc {
	const (
		processingMarker = "PROCESSING"
		cacheTTL         = 5 * time.Minute
	)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			userID := c.Request().Header.Get("X-User-ID")
			idemKey := c.Request().Header.Get("Idempotency-Key")

			if idemKey == "" || userID == "" {
				return next(c) // skip if not provided
			}

			key := "idem:" + userID + ":" + idemKey
			lockKey := "lock:" + key
			lockValue := strconv.FormatInt(time.Now().UnixNano(), 10)

			// Atomically decide whether to replay, wait/retry, or process.
			res, err := idemScript.Run(
				ctx,
				rdb,
				[]string{key, lockKey},
				processingMarker,
				lockValue,
				int(cacheTTL.Seconds()),
			).Result()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "idempotency error")
			}

			fmt.Println("res-----:", res)

			result, ok := res.([]interface{})
			if !ok || len(result) == 0 {
				return echo.NewHTTPError(http.StatusInternalServerError, "invalid idempotency state")
			}

			fmt.Println("result-----:", result)

			state, _ := result[0].(string)
			fmt.Println("state-----:", state)
			switch state {
			case "DONE":
				if len(result) < 2 {
					return echo.NewHTTPError(http.StatusInternalServerError, "invalid idempotency cache")
				}
				raw, _ := result[1].(string)
				var cached CachedResponse
				if err := json.Unmarshal([]byte(raw), &cached); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "invalid cached response")
				}
				return c.Blob(cached.Status, echo.MIMEApplicationJSON, cached.Body)
			case "PROCESSING", "RETRY":
				return echo.NewHTTPError(http.StatusConflict, "request already in progress")
			case "LOCKED":
				// Current request owns the lock and can process.
			default:
				return echo.NewHTTPError(http.StatusInternalServerError, "unknown idempotency state")
			}

			// Capture response for caching once handler completes.
			rec := &responseRecorder{ResponseWriter: c.Response().Writer}
			c.Response().Writer = rec

			err = next(c)
			if err != nil {
				// Release keys so client can retry failed execution.
				_ = rdb.Del(ctx, key, lockKey).Err()
				return err
			}

			status := rec.status
			if status == 0 {
				status = http.StatusOK
			}

			cache := CachedResponse{
				Status: status,
				Body:   append(json.RawMessage(nil), rec.body.Bytes()...),
			}

			data, err := json.Marshal(cache)
			if err != nil {
				_ = rdb.Del(ctx, key, lockKey).Err()
				return echo.NewHTTPError(http.StatusInternalServerError, "idempotency serialization error")
			}

			if err := rdb.Set(ctx, key, data, cacheTTL).Err(); err != nil {
				_ = rdb.Del(ctx, key, lockKey).Err()
				return echo.NewHTTPError(http.StatusInternalServerError, "idempotency store error")
			}

			// Release lock; cached value is the replay source.
			_ = rdb.Del(ctx, lockKey).Err()

			return nil
		}
	}
}

type responseRecorder struct {
	http.ResponseWriter
	body   bytes.Buffer
	status int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
