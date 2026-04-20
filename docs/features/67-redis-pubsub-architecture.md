# 🏗️ Architecture: Redis Pub/Sub Package

> **Feature**: `67` — Redis Pub/Sub Package
> **Discussion**: [`67-redis-pubsub-discussion.md`](67-redis-pubsub-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-04-20

---

## Overview

New `core/pubsub` package providing cross-process publish/subscribe messaging over Redis. Exposes `Publisher` and `Subscriber` interfaces backed by go-redis v9 PUBLISH/SUBSCRIBE. The subscriber auto-reconnects with exponential backoff and health-pings the connection on a configurable interval. Both constructors accept a pre-built `*redis.Client` (same pattern as `core/queue.RedisDriver`).

## File Structure

```
core/pubsub/
├── pubsub.go              # Interfaces (Publisher, Subscriber, Handler), SubscriberOptions, defaults
├── publisher.go           # redisPublisher — wraps *redis.Client.Publish()
├── subscriber.go          # redisSubscriber — Subscribe loop with reconnect + backoff + ping
└── subscriber_test.go     # Tests using miniredis: publish/receive, reconnect, context cancel, multi-channel
```

**Modified files:**

```
docs/framework/infrastructure/
├── pubsub.md              # NEW — full documentation
└── caching.md             # MODIFY — add cross-ref to pubsub.md in §7.1

rapidgo-starter/
└── app/providers/redis_provider.go   # MODIFY — uncomment redis.pubsub example, add Boot() pubsub usage
```

## Data Model

N/A — this feature does not introduce data models. Redis pub/sub is fire-and-forget with no persistence.

## Component Design

### Interfaces — `pubsub.go`

**Responsibility**: Define contracts and configuration types
**Location**: `core/pubsub/pubsub.go`

```go
// Handler processes a received message.
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
    MinBackoff time.Duration // Default: 500ms
    MaxBackoff time.Duration // Default: 30s
    PingEvery  time.Duration // Default: 30s — must be < Redis server timeout
    Logger     *slog.Logger  // Default: slog.Default()
}
```

### redisPublisher — `publisher.go`

**Responsibility**: Thin wrapper around `*redis.Client.Publish()`
**Location**: `core/pubsub/publisher.go`

```
Key Methods:
└── Publish(ctx, channel, payload) → error    # Delegates to client.Publish().Err()
```

```go
type redisPublisher struct {
    client *redis.Client
}

func NewRedisPublisher(client *redis.Client) Publisher
```

### redisSubscriber — `subscriber.go`

**Responsibility**: Subscribe to channels with reconnect, backoff, and ping
**Location**: `core/pubsub/subscriber.go`

```
Key Methods:
├── Subscribe(ctx, channels, handler) → error   # Outer loop: reconnect with backoff
├── listen(ctx, ps, handler) → error             # Inner loop: ReceiveMessage + dispatch
└── applyDefaults(opts) → SubscriberOptions      # Fill zero-value fields with defaults
```

```go
type redisSubscriber struct {
    client *redis.Client
    opts   SubscriberOptions
}

func NewRedisSubscriber(client *redis.Client, opts SubscriberOptions) Subscriber
```

**Subscribe flow:**

```
Subscribe(ctx, channels, handler)
│
├── outer loop (reconnect):
│   ├── client.Subscribe(ctx, channels...)  → *redis.PubSub
│   ├── start ping ticker (opts.PingEvery)
│   ├── listen(ctx, ps, handler)            → blocks until error or ctx done
│   ├── ps.Close()                          → cleanup
│   ├── if ctx.Err() != nil → return nil    → clean shutdown
│   ├── log.Warn("disconnected, reconnecting", "backoff", delay)
│   ├── sleep(delay) with jitter
│   └── delay = min(delay * 2, opts.MaxBackoff)
│
└── listen(ctx, ps, handler):
    ├── for { ps.ReceiveMessage(ctx) }
    ├── on message → handler(ctx, msg.Channel, msg.Payload)
    └── on error → return error (triggers reconnect)
```

## Data Flow

```
[App code] → Publisher.Publish(ctx, channel, payload)
         → redis.Client.Publish(ctx, channel, payload)
         → Redis server
         → redis.PubSub.ReceiveMessage()
         → Subscriber dispatches → Handler(ctx, channel, payload)
         → [App code processes message]
```

Cross-process:

```
Process A (Go web app)          Redis Server          Process B (Go worker / PHP admin)
─────────────────────          ────────────           ──────────────────────────────────
Publisher.Publish("ch", msg) → PUBLISH ch msg       → Subscriber receives → Handler(ctx, "ch", msg)
```

## Configuration

No new environment variables. The `*redis.Client` is created externally via `cache.NewRedisClient()` which already reads `REDIS_*` env vars.

Subscriber behavior is configured via `SubscriberOptions` at construction time:

| Field | Type | Default | Description |
|---|---|---|---|
| `MinBackoff` | `time.Duration` | 500ms | Initial reconnect delay |
| `MaxBackoff` | `time.Duration` | 30s | Maximum reconnect delay (caps exponential growth) |
| `PingEvery` | `time.Duration` | 30s | Interval for health-check pings on the subscription connection |
| `Logger` | `*slog.Logger` | `slog.Default()` | Structured logger for reconnect/error events |

## Security Considerations

- **No auth changes**: Redis authentication is handled by the `*redis.Client` passed in (configured via `REDIS_PASSWORD` env var through `cache.NewRedisClient()`)
- **Channel names**: The package does not validate or restrict channel names. Apps should use well-defined prefixes (e.g., `myapp:events:`) to avoid collisions
- **Payload content**: Payloads are raw strings — no deserialization happens inside the package. Apps must validate/sanitize payloads in their handlers
- **No TLS in scope**: TLS for Redis connections is a client-level concern, not pubsub-level

## Trade-offs & Alternatives

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| **New `core/pubsub` package** | Clean separation, zero coupling to `core/events`, matches package-per-concern pattern | One more package to discover | ✅ Selected |
| Add to `core/events` | Single place for all event-like things | Leaks go-redis dep to all consumers, mixes in-process and network semantics | ❌ Breaks zero-dep contract |
| Add to `core/cache` | Redis-related code stays together | Cache and pubsub are different concerns; bloats cache package | ❌ Wrong abstraction |
| New `core/redis` package | Houses all Redis primitives | Overlaps with `core/cache` which already owns the client factory; forces moving code | ❌ Disrupts existing API |

## Explicit Limitations (document loudly)

1. **At-most-once delivery** — Messages published while a subscriber is disconnected are lost. Redis pub/sub has no persistence.
2. **No ordering guarantees** beyond what Redis provides (single-publisher order preserved per channel).
3. **No Redis Streams** — This is pub/sub only. Streams (durable, consumer groups) are a different feature for the future.
4. **No PSUBSCRIBE** in v1 — Pattern subscriptions can be added later without breaking changes.

## Next

Create tasks doc → `67-redis-pubsub-tasks.md`
