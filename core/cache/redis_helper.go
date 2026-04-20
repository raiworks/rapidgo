package cache

import (
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient builds a go-redis client from environment variables.
// If dbOverride is non-nil, it takes precedence over the REDIS_DB env var.
//
// Environment variables read:
//
//	REDIS_HOST          (default "localhost")
//	REDIS_PORT          (default "6379")
//	REDIS_PASSWORD      (default "")
//	REDIS_DB            (default 0)
//	REDIS_POOL_SIZE     (default 10)
//	REDIS_DIAL_TIMEOUT  (default "5s")
//	REDIS_READ_TIMEOUT  (default "3s")
//	REDIS_WRITE_TIMEOUT (default "3s")
func NewRedisClient(dbOverride *int) *redis.Client {
	db := envInt("REDIS_DB", 0)
	if dbOverride != nil {
		db = *dbOverride
	}
	return redis.NewClient(&redis.Options{
		Addr:         envStr("REDIS_HOST", "localhost") + ":" + envStr("REDIS_PORT", "6379"),
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           db,
		PoolSize:     envInt("REDIS_POOL_SIZE", 10),
		DialTimeout:  envDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:  envDuration("REDIS_READ_TIMEOUT", 3*time.Second),
		WriteTimeout: envDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
	})
}

// envStr returns the value of the named environment variable or fallback.
func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// envInt returns the integer value of the named environment variable or fallback.
// Returns fallback if the value is empty, non-numeric, or negative.
func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return fallback
	}
	return n
}

// envDuration returns the duration value of the named environment variable or fallback.
// The value must be a valid Go duration string (e.g. "5s", "100ms", "2m").
func envDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
