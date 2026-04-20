# ✅ Tasks: Redis DB Selection & Multi-Client Helper

> **Feature**: `66` — Redis DB Selection & Multi-Client Helper
> **Architecture**: [`66-redis-db-multi-client-architecture.md`](66-redis-db-multi-client-architecture.md)
> **Branch**: `feature/66-redis-db-multi-client`
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/14 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [x] Feature branch created from latest `main`
- [x] Dependent features are merged to `main`
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — Core Logic (NewRedisClient helper)

> New exported helper + env parsing helpers.

- [ ] **A.1** — Create `core/cache/redis_helper.go` with `envStr()`, `envInt()`, `envDuration()` unexported helpers
- [ ] **A.2** — Implement `NewRedisClient(dbOverride *int) *redis.Client` in `redis_helper.go`
- [ ] 📍 **Checkpoint A** — `NewRedisClient(nil)` builds a client with default env values, `NewRedisClient(&db)` overrides DB

---

## Phase B — Integration (cache + session)

> Wire existing code to use the new helper.

- [ ] **B.1** — Refactor `newRedisClient()` in `core/cache/cache.go` to delegate to `NewRedisClient(nil)`
- [ ] **B.2** — Refactor Redis branch in `core/session/factory.go` to use `cache.NewRedisClient(nil)` instead of inline client creation
- [ ] 📍 **Checkpoint B** — Existing cache and session Redis paths use the shared helper

---

## Phase C — Testing

> Tests for the new helper + verify existing tests still pass.

- [ ] **C.1** — Create `core/cache/redis_helper_test.go` with tests:
  - `TestNewRedisClient_DefaultDB` — no REDIS_DB set → DB 0
  - `TestNewRedisClient_EnvDB` — REDIS_DB=3 → DB 3
  - `TestNewRedisClient_DBOverride` — REDIS_DB=3 + override 5 → DB 5
  - `TestNewRedisClient_InvalidDB` — REDIS_DB=abc → DB 0
  - `TestNewRedisClient_PoolAndTimeouts` — REDIS_POOL_SIZE, timeout env vars honored
- [ ] **C.2** — Run `go test ./core/cache/...` — all pass (existing + new)
- [ ] **C.3** — Run `go test ./core/session/...` — all pass (session factory still works)
- [ ] **C.4** — Run `go test ./...` — full regression green
- [ ] 📍 **Checkpoint C** — All tests pass, no regressions

---

## Phase D — Starter Update

> Update `rapidgo-starter` to show the named-client pattern.

- [ ] **D.1** — Update `rapidgo-starter/app/providers/redis_provider.go` to use `cache.NewRedisClient(nil)` and show commented named-client examples
- [ ] **D.2** — Update `rapidgo-starter/.env.example` to document new optional env vars (`REDIS_POOL_SIZE`, `REDIS_DIAL_TIMEOUT`, `REDIS_READ_TIMEOUT`, `REDIS_WRITE_TIMEOUT`)
- [ ] 📍 **Checkpoint D** — Starter builds and reflects new pattern

---

## Phase E — Documentation & Cleanup

> Changelog, docs, self-review.

- [ ] **E.1** — Update `docs/framework/infrastructure/caching.md` with REDIS_DB and multi-client pattern
- [ ] **E.2** — Add v2.7.3 entry to `docs/CHANGELOG.md`
- [ ] **E.3** — Update changelog doc with final summary
- [ ] **E.4** — Self-review all diffs
- [ ] 📍 **Checkpoint E** — Clean code, complete docs, ready to ship

---

## Ship 🚀

- [ ] All phases complete
- [ ] All checkpoints verified
- [ ] Final commit with descriptive message
- [ ] Push to feature branch
- [ ] Merge to `main`
- [ ] Push `main`
- [ ] Tag `v2.7.3`
- [ ] Push tag
- [ ] Create GitHub release (marked as Latest)
- [ ] **Keep the feature branch** — do not delete
- [ ] Create review doc → `66-redis-db-multi-client-review.md`
