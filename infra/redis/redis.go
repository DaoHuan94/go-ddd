package redis

import (
	"context"
	"fmt"
	"go-ddd/infra/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Rdb *redis.Client
	cfg config.Config
}

func NewRedisClient(ctx context.Context, cfg config.Config) *RedisClient {
	return &RedisClient{
		Rdb: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}),
		cfg: cfg,
	}
}

func (c *RedisClient) Close() error {
	return c.Rdb.Close()
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.Rdb.Get(ctx, key).Result()
}

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Rdb.Set(ctx, key, value, expiration).Err()
}

func (c *RedisClient) Del(ctx context.Context, key string) error {
	return c.Rdb.Del(ctx, key).Err()
}

func (c *RedisClient) Exists(ctx context.Context, key string) (int64, error) {
	return c.Rdb.Exists(ctx, key).Result()
}

func (c *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.Rdb.Expire(ctx, key, expiration).Err()
}
