package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisDriver implements Driver using Redis lists.
type RedisDriver struct {
	client *redis.Client
	prefix string
}

// NewRedisDriver creates a Redis-backed queue driver.
func NewRedisDriver(client *redis.Client) *RedisDriver {
	return &RedisDriver{
		client: client,
		prefix: "rapidgo:queue:",
	}
}

func (d *RedisDriver) key(queue string) string {
	return d.prefix + queue
}

func (d *RedisDriver) failedKey(queue string) string {
	return d.prefix + "failed:" + queue
}

func (d *RedisDriver) Push(ctx context.Context, job *Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("queue: redis marshal failed: %w", err)
	}
	return d.client.LPush(ctx, d.key(job.Queue), data).Err()
}

func (d *RedisDriver) Pop(ctx context.Context, queue string) (*Job, error) {
	data, err := d.client.RPop(ctx, d.key(queue)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("queue: redis pop failed: %w", err)
	}

	var job Job
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, fmt.Errorf("queue: redis unmarshal failed: %w", err)
	}

	now := time.Now()
	job.ReservedAt = &now
	return &job, nil
}

func (d *RedisDriver) Delete(_ context.Context, _ *Job) error {
	// Already removed by RPop — nothing to do.
	return nil
}

func (d *RedisDriver) Release(ctx context.Context, job *Job, delay time.Duration) error {
	job.Attempts++
	job.ReservedAt = nil
	job.AvailableAt = time.Now().Add(delay)

	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("queue: redis marshal failed: %w", err)
	}
	return d.client.LPush(ctx, d.key(job.Queue), data).Err()
}

func (d *RedisDriver) Fail(ctx context.Context, job *Job, _ error) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("queue: redis marshal failed: %w", err)
	}
	return d.client.LPush(ctx, d.failedKey(job.Queue), data).Err()
}

func (d *RedisDriver) Size(ctx context.Context, queue string) (int64, error) {
	return d.client.LLen(ctx, d.key(queue)).Result()
}
