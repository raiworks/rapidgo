package pubsub

import (
	"context"
	"log/slog"
	"time"
)

// Handler processes a received pub/sub message.
type Handler func(ctx context.Context, channel string, payload string)

// Publisher sends messages to a Redis channel.
type Publisher interface {
	Publish(ctx context.Context, channel string, payload string) error
}

// Subscriber listens on one or more Redis channels.
// Subscribe blocks until ctx is cancelled. Reconnects automatically on disconnect.
type Subscriber interface {
	Subscribe(ctx context.Context, channels []string, h Handler) error
}

// SubscriberOptions configures reconnect and health-check behavior.
type SubscriberOptions struct {
	MinBackoff time.Duration // Initial reconnect delay. Default: 500ms.
	MaxBackoff time.Duration // Maximum reconnect delay (caps exponential growth). Default: 30s.
	PingEvery  time.Duration // Health-check ping interval. Must be < Redis server timeout. Default: 30s.
	Logger     *slog.Logger  // Structured logger for reconnect/error events. Default: slog.Default().
}

// applyDefaults fills zero-value fields with sensible defaults.
func applyDefaults(opts SubscriberOptions) SubscriberOptions {
	if opts.MinBackoff <= 0 {
		opts.MinBackoff = 500 * time.Millisecond
	}
	if opts.MaxBackoff <= 0 {
		opts.MaxBackoff = 30 * time.Second
	}
	if opts.PingEvery <= 0 {
		opts.PingEvery = 30 * time.Second
	}
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}
	return opts
}
