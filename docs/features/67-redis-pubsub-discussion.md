# 💬 Discussion: Redis Pub/Sub Package

> **Feature**: `67` — Redis Pub/Sub Package
> **Status**: 🟡 IN PROGRESS
> **Branch**: `feature/67-redis-pubsub`
> **Depends On**: #66 (NewRedisClient helper)
> **Date Started**: 2026-04-20
> **Date Completed**: —

---

## Summary

Add a new `core/pubsub` package that provides cross-process publish/subscribe messaging over Redis. The package exposes clean `Publisher` and `Subscriber` interfaces backed by Redis PUBLISH/SUBSCRIBE, with automatic reconnect and exponential backoff. This fills the gap between `core/events` (in-process only) and the need for multi-process communication in scaled deployments.

## Functional Requirements

- As a developer, I want to publish messages to a Redis channel so that other processes can receive them
- As a developer, I want to subscribe to Redis channels with automatic reconnect so that I don't have to implement retry/backoff logic myself
- As a developer, I want configurable backoff and ping intervals so that I can tune reconnect behavior for my environment
- As a developer, I want a clean separation between in-process events (`core/events`) and cross-process pub/sub (`core/pubsub`) so that I can choose the right tool for each use case

## Current State / Reference

### What Exists

- **`core/events`** — Pure in-process event dispatcher. `Dispatch()` (async goroutines) and `DispatchSync()` (sequential). Uses `map[string][]Handler` with `sync.RWMutex`. Zero external deps. Works well for single-process scenarios.
- **`core/queue`** — Job queue with a `RedisDriver` that accepts a pre-built `*redis.Client`. Uses Redis lists (LPUSH/RPOP), not pub/sub. Pattern: `NewRedisDriver(client *redis.Client)`.
- **`core/cache`** — `NewRedisClient(dbOverride *int)` (exported in v2.7.3) builds a `*redis.Client` from env vars. This is the standard Redis client factory.
- **go-redis v9.18.0** — Already in `go.mod`. Supports `Subscribe()`, `PSubscribe()`, and `*PubSub` with `ReceiveMessage()`.

### What Works Well

- `core/events` is clean, simple, and has zero dependencies — perfect for in-process use
- `core/queue.RedisDriver` proves the pattern of accepting an external `*redis.Client` works well
- `cache.NewRedisClient()` centralizes Redis connection config so pubsub doesn't need to re-invent env parsing

### What Needs Improvement

- No framework-level primitive for cross-process messaging
- Developers using raw `*redis.Client.Subscribe()` must re-implement reconnect, backoff, context cancellation, ping health checks, and structured logging every time
- No guidance on when to use `core/events` vs Redis pub/sub

## Proposed Approach

Create a new `core/pubsub` package with:

1. **Interfaces**: `Publisher` (Publish) and `Subscriber` (Subscribe with reconnect)
2. **Redis implementation**: `NewRedisPublisher(client)` and `NewRedisSubscriber(client, opts)`
3. **Reconnect logic**: Exponential backoff (configurable min/max), auto-reconnect on disconnect, ping interval for dead connection detection
4. **Handler type**: `func(ctx context.Context, channel, payload string)` — matches the ergonomics of `core/events.Handler` but adds context and channel info
5. **Accepts pre-built `*redis.Client`** — same pattern as `core/queue.RedisDriver`. The pubsub package does NOT create its own Redis clients; callers use `cache.NewRedisClient()` or bring their own.

### Why a new package (not in `core/events`)

- `core/events` is in-process with zero external deps — adding a Redis import leaks it to all consumers
- Pub/sub has fundamentally different semantics: network-bound, at-most-once, requires reconnect logic
- Matches existing package-per-concern pattern: `core/cache`, `core/session`, `core/queue`
- Keeps the "which tool do I use?" question answerable: events = in-process, pubsub = cross-process

## Edge Cases & Risks

- [ ] **At-most-once delivery**: Redis pub/sub has no persistence. If a subscriber is disconnected when a message is published, that message is lost. Must document loudly.
- [ ] **Reconnect during high-throughput**: Messages published during the reconnect window are silently dropped by Redis. Acceptable for the use cases this targets (cache invalidation, UI refresh signals) but must be documented.
- [ ] **Subscriber goroutine lifecycle**: `Subscribe()` blocks until context is cancelled. Callers must manage goroutine shutdown properly. Provide clear examples.
- [ ] **Multiple channels**: A single subscriber should be able to listen on multiple channels in one call (Redis supports this natively).
- [ ] **go-redis PubSub object cleanup**: Must `Close()` the underlying `*redis.PubSub` on context cancel and reconnect. Leaking subscriptions exhausts Redis client connections.
- [ ] **Ping interval vs Redis timeout**: If `PingEvery` exceeds the Redis server's `timeout` config, the server kills the connection before the client detects it. Document that `PingEvery` should be shorter than server timeout.

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Feature #66 — NewRedisClient helper | Feature | ✅ Done |
| go-redis/v9 v9.18.0 | Library | ✅ Already in go.mod |
| alicebob/miniredis/v2 | Test dep | ✅ Already in go.mod (supports PUBLISH/SUBSCRIBE) |

## Scope

### In Scope (from spec Change 3 + Change 4 + Change 5)

- `core/pubsub` package: `Publisher` interface, `Subscriber` interface, `Handler` type
- `NewRedisPublisher(client *redis.Client) Publisher`
- `NewRedisSubscriber(client *redis.Client, opts SubscriberOptions) Subscriber`
- Reconnect with exponential backoff (configurable min/max backoff, ping interval)
- Structured logging via `log/slog`
- Tests using miniredis
- New `docs/framework/infrastructure/pubsub.md` documentation
- Starter update: pubsub usage example in `redis_provider.go` (commented `redis.pubsub` singleton already present from #66)

### Out of Scope

- No Kafka/NATS/SNS adapters — Redis only
- No message ordering guarantees beyond what Redis pub/sub offers
- No persistence / durability (no Redis Streams)
- No bridge between `core/events` and `core/pubsub` — future opinion
- No pattern subscriptions (PSUBSCRIBE) in v1 — can be added later if needed

## Open Questions

- [x] Should this be `core/pubsub` or `core/redis`? → `core/pubsub` — cleaner separation, `core/cache` already owns the Redis client factory
- [x] Should we use `log/slog` or the framework's `core/logger`? → `log/slog` — keeps the package lightweight with zero internal deps, same as `core/events`
- [x] Version: v2.7.4 or v2.8.0? → **v2.8.0** — new package is a MINOR bump per semver

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-04-20 | Package location: `core/pubsub` | Keeps `core/cache` focused on caching; matches package-per-concern pattern |
| 2026-04-20 | Accept `*redis.Client` (not create one) | Same proven pattern as `core/queue.RedisDriver` |
| 2026-04-20 | Use `log/slog` for logging | Zero internal deps, lightweight, standard library |
| 2026-04-20 | Release as v2.8.0 | New public package = new feature = semver MINOR bump |
| 2026-04-20 | No PSUBSCRIBE in v1 | Keep scope tight; can add pattern subscriptions later |

## Discussion Complete ✅

**Summary**: New `core/pubsub` package providing Redis-backed cross-process pub/sub with `Publisher`/`Subscriber` interfaces, automatic reconnect with exponential backoff, and structured logging. Accepts pre-built `*redis.Client` from `cache.NewRedisClient()`. Ships as v2.8.0 with documentation and starter example.
**Completed**: 2026-04-20
**Next**: Create architecture doc → `67-redis-pubsub-architecture.md`
