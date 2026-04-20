package pubsub

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type redisPublisher struct {
	client *redis.Client
}

// NewRedisPublisher returns a Publisher that sends messages via Redis PUBLISH.
func NewRedisPublisher(client *redis.Client) Publisher {
	return &redisPublisher{client: client}
}

func (p *redisPublisher) Publish(ctx context.Context, channel string, payload string) error {
	return p.client.Publish(ctx, channel, payload).Err()
}
