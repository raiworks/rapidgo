package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache is a cache backend backed by Redis.
type RedisCache struct {
	Client *redis.Client
}

// NewRedisCache returns a Redis-backed cache store.
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{Client: client}
}

func (c *RedisCache) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := c.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (c *RedisCache) Set(key, value string, ttl time.Duration) error {
	ctx := context.Background()
	return c.Client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Delete(key string) error {
	ctx := context.Background()
	return c.Client.Del(ctx, key).Err()
}

func (c *RedisCache) Flush() error {
	ctx := context.Background()
	return c.Client.FlushDB(ctx).Err()
}
