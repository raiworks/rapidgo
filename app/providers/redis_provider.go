package providers

import (
	"github.com/RAiWorks/RapidGo/core/config"
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/redis/go-redis/v9"
)

// RedisProvider registers a shared *redis.Client singleton in the container.
type RedisProvider struct{}

// Register binds a *redis.Client singleton. Connection is lazy.
func (p *RedisProvider) Register(c *container.Container) {
	c.Singleton("redis", func(c *container.Container) interface{} {
		return redis.NewClient(&redis.Options{
			Addr:     config.Env("REDIS_HOST", "localhost") + ":" + config.Env("REDIS_PORT", "6379"),
			Password: config.Env("REDIS_PASSWORD", ""),
		})
	})
}

// Boot is a no-op.
func (p *RedisProvider) Boot(c *container.Container) {}
