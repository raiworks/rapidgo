---
title: "Pub/Sub"
version: "0.1.0"
status: "Final"
date: "2026-04-20"
last_updated: "2026-04-20"
authors:
  - "raiworks"
supersedes: ""
---

# Pub/Sub

## Abstract

This document covers the `core/pubsub` package — a cross-process
publish/subscribe layer built on Redis. It provides `Publisher` and
`Subscriber` interfaces with automatic reconnect and exponential
backoff. For in-process events, see [`core/events`](../../features/).

## Table of Contents

1. [When to Use What](#1-when-to-use-what)
2. [API Reference](#2-api-reference)
3. [Usage Examples](#3-usage-examples)
4. [Reconnect Semantics](#4-reconnect-semantics)
5. [Provider Registration](#5-provider-registration)
6. [Limitations](#6-limitations)
7. [Security Considerations](#7-security-considerations)
8. [References](#8-references)

## 1. When to Use What

| Scenario | Package | Why |
|---|---|---|
| Single-process event (e.g., "user created" triggers email) | `core/events` | In-process, zero deps, sync or async goroutines |
| Cross-process notification (e.g., cache invalidation across instances) | `core/pubsub` | Network-bound, Redis-backed, auto-reconnect |
| Durable job processing (e.g., send email in background) | `core/queue` | Persistent (Redis lists), retry, backoff, worker pool |

**Rule of thumb**: If the receiver might be a different OS process,
use `core/pubsub`. If it's always the same process, use `core/events`.
If it needs to survive a crash, use `core/queue`.

## 2. API Reference

### Handler

```go
type Handler func(ctx context.Context, channel string, payload string)
```

Processes a received message. Called once per message. If a handler
panics, the panic is recovered and logged — the subscriber continues
operating.

### Publisher

```go
type Publisher interface {
    Publish(ctx context.Context, channel string, payload string) error
}

func NewRedisPublisher(client *redis.Client) Publisher
```

Thin wrapper around `redis.Client.Publish()`. Returns an error if
the Redis command fails.

### Subscriber

```go
type Subscriber interface {
    Subscribe(ctx context.Context, channels []string, h Handler) error
}

func NewRedisSubscriber(client *redis.Client, opts SubscriberOptions) Subscriber
```

`Subscribe` **blocks** until `ctx` is cancelled. On disconnect it
reconnects with exponential backoff. Returns `nil` on clean shutdown
(context cancelled). Returns an error only for unrecoverable problems
(e.g., empty channels list).

### SubscriberOptions

```go
type SubscriberOptions struct {
    MinBackoff time.Duration // Default: 500ms
    MaxBackoff time.Duration // Default: 30s
    PingEvery  time.Duration // Default: 30s
    Logger     *slog.Logger  // Default: slog.Default()
}
```

| Field | Default | Description |
|---|---|---|
| `MinBackoff` | 500ms | Initial delay before first reconnect attempt |
| `MaxBackoff` | 30s | Maximum delay (caps exponential doubling) |
| `PingEvery` | 30s | Health-check ping interval. Must be shorter than the Redis server's `timeout` setting |
| `Logger` | `slog.Default()` | Structured logger for reconnect and error events |

Zero values are replaced with defaults automatically.

## 3. Usage Examples

### Publish a Message

```go
client := cache.NewRedisClient(nil)
pub := pubsub.NewRedisPublisher(client)

err := pub.Publish(ctx, "myapp:events:user.created", `{"id":42}`)
if err != nil {
    slog.Error("publish failed", "error", err)
}
```

### Subscribe to Channels

```go
client := cache.NewRedisClient(nil)
sub := pubsub.NewRedisSubscriber(client, pubsub.SubscriberOptions{})

// Subscribe blocks — run in a goroutine.
go func() {
    err := sub.Subscribe(ctx, []string{"myapp:events:user.created"}, func(ctx context.Context, ch, payload string) {
        slog.Info("received", "channel", ch, "payload", payload)
    })
    if err != nil {
        slog.Error("subscriber error", "error", err)
    }
}()
```

### Graceful Shutdown

```go
ctx, cancel := context.WithCancel(context.Background())

go sub.Subscribe(ctx, channels, handler)

// On shutdown signal:
cancel()
// Subscribe() returns nil — goroutine exits cleanly.
```

### Multiple Named Clients (Different DBs)

```go
// In a provider:
c.Singleton("redis.pubsub", func(_ *container.Container) interface{} {
    db := 5
    return cache.NewRedisClient(&db)
})

// In Boot() or app startup:
client := container.MustMake[*redis.Client](c, "redis.pubsub")
pub := pubsub.NewRedisPublisher(client)
sub := pubsub.NewRedisSubscriber(client, pubsub.SubscriberOptions{})
```

## 4. Reconnect Semantics

When the Redis connection drops, the subscriber:

1. Closes the old PubSub connection
2. Logs a warning with the current backoff delay
3. Sleeps for the backoff duration
4. Doubles the delay (capped at `MaxBackoff`)
5. Attempts to re-subscribe

On a successful reconnection, the backoff resets to `MinBackoff`.

**Ping health check**: The subscriber sends a `PING` on the PubSub
connection every `PingEvery` interval. This detects dead connections
that would otherwise hang silently (common behind NAT or load
balancers). Set `PingEvery` shorter than the Redis server's `timeout`
configuration.

## 5. Provider Registration

```go
type PubSubProvider struct{}

func (p *PubSubProvider) Register(c *container.Container) {
    // Dedicated Redis client for pub/sub (optional — can share with cache)
    c.Singleton("redis.pubsub", func(_ *container.Container) interface{} {
        db := 5
        return cache.NewRedisClient(&db)
    })
}

func (p *PubSubProvider) Boot(c *container.Container) {
    client := container.MustMake[*redis.Client](c, "redis.pubsub")

    // Register publisher
    c.Singleton("pubsub.publisher", func(_ *container.Container) interface{} {
        return pubsub.NewRedisPublisher(client)
    })

    // Start subscriber in background
    sub := pubsub.NewRedisSubscriber(client, pubsub.SubscriberOptions{})
    go sub.Subscribe(context.Background(), []string{"myapp:events"}, func(ctx context.Context, ch, payload string) {
        // Handle messages...
    })
}
```

## 6. Limitations

These are fundamental to Redis pub/sub and **cannot** be worked around
at the library level:

1. **At-most-once delivery** — If a subscriber is disconnected when a
   message is published, that message is lost permanently. Redis pub/sub
   has no message persistence or replay.

2. **No ordering guarantees** — Messages from a single publisher to a
   single channel arrive in order. Cross-publisher or cross-channel
   ordering is not guaranteed.

3. **No consumer groups** — Every subscriber receives every message.
   There is no load-balancing across subscribers (unlike Redis Streams
   or Kafka).

4. **Not suitable for durable work** — If you need at-least-once
   delivery or retry on failure, use `core/queue` instead.

5. **No pattern subscriptions** — `PSUBSCRIBE` (wildcard channels) is
   not supported in v1. Can be added in a future release.

## 7. Security Considerations

- Redis authentication is handled by the `*redis.Client` passed in
  (via `REDIS_PASSWORD`). The pubsub package does not manage auth.
- Channel names **SHOULD** use well-defined prefixes (e.g.,
  `myapp:events:`) to avoid collisions with other apps on the same
  Redis instance.
- Payloads are raw strings. The pubsub package does **not** validate
  or deserialize them. Handlers **MUST** validate payloads before use.

## 8. References

- [Redis PUBLISH documentation](https://redis.io/commands/publish)
- [Redis SUBSCRIBE documentation](https://redis.io/commands/subscribe)
- [go-redis PubSub](https://pkg.go.dev/github.com/redis/go-redis/v9#PubSub)
- [Caching documentation](caching.md) — Redis client factory (`cache.NewRedisClient`)
