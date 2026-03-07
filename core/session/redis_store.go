package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore stores sessions in Redis.
type RedisStore struct {
	Client *redis.Client
	Prefix string
}

// NewRedisStore creates a Redis-backed session store from the given client.
func NewRedisStore(client *redis.Client, prefix string) *RedisStore {
	if prefix == "" {
		prefix = "session:"
	}
	return &RedisStore{Client: client, Prefix: prefix}
}

func (s *RedisStore) key(id string) string {
	return s.Prefix + id
}

func (s *RedisStore) Read(id string) (map[string]interface{}, error) {
	ctx := context.Background()
	raw, err := s.Client.Get(ctx, s.key(id)).Result()
	if err == redis.Nil {
		return make(map[string]interface{}), nil
	}
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *RedisStore) Write(id string, data map[string]interface{}, lifetime time.Duration) error {
	ctx := context.Background()
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, s.key(id), string(raw), lifetime).Err()
}

func (s *RedisStore) Destroy(id string) error {
	ctx := context.Background()
	return s.Client.Del(ctx, s.key(id)).Err()
}

func (s *RedisStore) GC(maxLifetime time.Duration) error {
	// Redis handles expiry automatically via TTL — no manual GC needed.
	return nil
}
