# ✅ Tasks: Redis Pub/Sub Package

> **Feature**: `67` — Redis Pub/Sub Package
> **Architecture**: [`67-redis-pubsub-architecture.md`](67-redis-pubsub-architecture.md)
> **Branch**: `feature/67-redis-pubsub`
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/16 tasks complete

---

## Pre-Flight

- [ ] Discussion doc is marked COMPLETE
- [ ] Architecture doc is FINALIZED
- [ ] Feature branch created from main
- [ ] Dependent features are merged to main (#66 ✅)

---

## Phase A — Interfaces & Publisher

> Define contracts and implement the simple publisher.

- [ ] **A.1** — Create `core/pubsub/pubsub.go`
  - [ ] `Handler` type: `func(ctx context.Context, channel string, payload string)`
  - [ ] `Publisher` interface with `Publish(ctx, channel, payload) error`
  - [ ] `Subscriber` interface with `Subscribe(ctx, channels, handler) error`
  - [ ] `SubscriberOptions` struct with `MinBackoff`, `MaxBackoff`, `PingEvery`, `Logger`
- [ ] **A.2** — Create `core/pubsub/publisher.go`
  - [ ] `redisPublisher` struct (unexported) with `client *redis.Client`
  - [ ] `NewRedisPublisher(client *redis.Client) Publisher` constructor
  - [ ] `Publish()` delegates to `client.Publish(ctx, channel, payload).Err()`
- [ ] 📍 **Checkpoint A** — `pubsub.go` and `publisher.go` compile, interfaces defined

---

## Phase B — Subscriber with Reconnect

> Core subscriber logic with exponential backoff and ping health checks.

- [ ] **B.1** — Create `core/pubsub/subscriber.go`
  - [ ] `redisSubscriber` struct (unexported) with `client` and `opts`
  - [ ] `NewRedisSubscriber(client *redis.Client, opts SubscriberOptions) Subscriber` constructor
  - [ ] `applyDefaults()` fills zero-value SubscriberOptions with defaults (500ms/30s/30s/slog.Default())
- [ ] **B.2** — Implement `Subscribe()` outer loop
  - [ ] Call `client.Subscribe(ctx, channels...)` to get `*redis.PubSub`
  - [ ] Start ping ticker (`opts.PingEvery`)
  - [ ] Call inner `listen()` loop
  - [ ] On error: close PubSub, log warning, sleep with backoff, double delay (capped at MaxBackoff)
  - [ ] On context cancel: close PubSub, return nil (clean shutdown)
- [ ] **B.3** — Implement `listen()` inner loop
  - [ ] `ReceiveMessage(ctx)` in a loop
  - [ ] On message: call `handler(ctx, msg.Channel, msg.Payload)`
  - [ ] On error: return error (triggers reconnect in outer loop)
- [ ] 📍 **Checkpoint B** — Subscriber compiles, `go vet ./core/pubsub/...` clean

---

## Phase C — Testing

> Hermetic tests using miniredis.

- [ ] **C.1** — Create `core/pubsub/subscriber_test.go`
  - [ ] TC-01: Publish → Subscribe → handler receives message
  - [ ] TC-02: Multiple channels — subscriber receives from all
  - [ ] TC-03: Context cancel — Subscribe returns nil (clean exit)
  - [ ] TC-04: Reconnect — kill miniredis, restart, publish post-reconnect → handler still receives
  - [ ] TC-05: Publisher returns error on failed publish
- [ ] 📍 **Checkpoint C** — `go test ./core/pubsub/... -v` all pass

---

## Phase D — Documentation

> Framework docs and starter update.

- [ ] **D.1** — Create `docs/framework/infrastructure/pubsub.md`
  - [ ] When to use `core/events` vs `core/pubsub`
  - [ ] API reference (Publisher, Subscriber, SubscriberOptions)
  - [ ] Usage examples (publish, subscribe, graceful shutdown)
  - [ ] Reconnect semantics and backoff tuning
  - [ ] Explicit limitations (at-most-once, no persistence)
- [ ] **D.2** — Update `docs/framework/infrastructure/caching.md`
  - [ ] Add cross-reference to pubsub.md in §7.1
- [ ] **D.3** — Update `rapidgo-starter/app/providers/redis_provider.go`
  - [ ] Uncomment `redis.pubsub` singleton or add pubsub usage comment in Boot()
- [ ] 📍 **Checkpoint D** — Docs complete, starter builds

---

## Phase E — Cleanup & Ship

> Final verification and release.

- [ ] **E.1** — Run full framework test suite: `go test ./... -count=1`
- [ ] **E.2** — Update changelog doc with final summary
- [ ] **E.3** — Update project roadmap with Feature #67
- [ ] 📍 **Checkpoint E** — Self-review all diffs, everything clean

---

## Ship 🚀

- [ ] All phases complete
- [ ] Final commit with descriptive message
- [ ] Push to feature branch
- [ ] Human approval received
- [ ] Merge to main
- [ ] Push main
- [ ] Tag v2.8.0, create GitHub release
- [ ] **Keep the feature branch** — do not delete
- [ ] Create review doc → `67-redis-pubsub-review.md`
