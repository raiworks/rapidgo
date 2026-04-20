package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisSubscriber struct {
	client *redis.Client
	opts   SubscriberOptions
}

// NewRedisSubscriber returns a Subscriber that listens on Redis channels with
// automatic reconnect and exponential backoff.
func NewRedisSubscriber(client *redis.Client, opts SubscriberOptions) Subscriber {
	return &redisSubscriber{
		client: client,
		opts:   applyDefaults(opts),
	}
}

// Subscribe blocks until ctx is cancelled. On disconnect it reconnects with
// exponential backoff. Returns nil on clean shutdown (context cancelled).
func (s *redisSubscriber) Subscribe(ctx context.Context, channels []string, h Handler) error {
	if len(channels) == 0 {
		return fmt.Errorf("pubsub: no channels specified")
	}

	delay := s.opts.MinBackoff

	for {
		ps := s.client.Subscribe(ctx, channels...)

		// Wait for confirmation that subscription is active.
		_, err := ps.Receive(ctx)
		if err != nil {
			ps.Close()
			if ctx.Err() != nil {
				return nil
			}
			s.opts.Logger.Warn("pubsub: subscribe failed", "error", err, "backoff", delay)
			if !s.sleep(ctx, delay) {
				return nil
			}
			delay = s.nextDelay(delay)
			continue
		}

		// Start pinging to detect dead connections.
		pingCtx, pingCancel := context.WithCancel(ctx)
		go s.ping(pingCtx, ps)

		// Block receiving messages until error or context cancellation.
		err = s.listen(ctx, ps, h)

		pingCancel()
		ps.Close()

		if ctx.Err() != nil {
			return nil
		}

		// Reset delay on successful connection (was receiving messages).
		delay = s.opts.MinBackoff
		s.opts.Logger.Warn("pubsub: disconnected, reconnecting", "error", err, "backoff", delay)
		if !s.sleep(ctx, delay) {
			return nil
		}
		delay = s.nextDelay(delay)
	}
}

// listen receives messages in a loop and dispatches them to the handler.
// Returns on error (triggers reconnect) or context cancellation.
func (s *redisSubscriber) listen(ctx context.Context, ps *redis.PubSub, h Handler) error {
	ch := ps.Channel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-ch:
			if !ok {
				return fmt.Errorf("pubsub: channel closed")
			}
			s.safeHandle(ctx, h, msg.Channel, msg.Payload)
		}
	}
}

// safeHandle calls the handler with panic recovery.
func (s *redisSubscriber) safeHandle(ctx context.Context, h Handler, channel, payload string) {
	defer func() {
		if r := recover(); r != nil {
			s.opts.Logger.Error("pubsub: handler panic", "channel", channel, "panic", r)
		}
	}()
	h(ctx, channel, payload)
}

// ping sends periodic PING commands on the PubSub connection to detect dead connections.
func (s *redisSubscriber) ping(ctx context.Context, ps *redis.PubSub) {
	ticker := time.NewTicker(s.opts.PingEvery)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := ps.Ping(ctx); err != nil {
				s.opts.Logger.Warn("pubsub: ping failed", "error", err)
				return
			}
		}
	}
}

// sleep waits for the given duration or until ctx is cancelled.
// Returns false if ctx was cancelled (caller should exit).
func (s *redisSubscriber) sleep(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

// nextDelay doubles the delay, capped at MaxBackoff.
func (s *redisSubscriber) nextDelay(current time.Duration) time.Duration {
	next := current * 2
	if next > s.opts.MaxBackoff {
		next = s.opts.MaxBackoff
	}
	return next
}
