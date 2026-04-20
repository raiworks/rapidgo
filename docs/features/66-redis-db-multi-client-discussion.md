# 💬 Discussion: Redis DB Selection & Multi-Client Helper

> **Feature**: `66` — Redis DB Selection & Multi-Client Helper
> **Status**: 🟢 COMPLETE
> **Branch**: `feature/66-redis-db-multi-client`
> **Depends On**: #32 (Caching), #20 (Session Management)
> **Date Started**: 2026-04-20
> **Date Completed**: 2026-04-20

---

## Summary

Two bugs/gaps in RapidGo's Redis integration:

1. **`REDIS_DB` env var is declared in `.env.example` but never consumed** — all Redis clients silently default to DB 0 regardless of config. This is a bugfix.
2. **No way to register multiple named Redis clients** — apps needing separate logical DBs (e.g., cache on DB 2, sessions on DB 0, pub/sub on DB 5) must hand-roll env parsing or use the broken `SELECT` anti-pattern. A new exported `NewRedisClient` helper with DB override solves this.

---

## Functional Requirements

- As a developer, I want `REDIS_DB` to be respected so that I can isolate cache/session/queue on different logical DBs.
- As a developer, I want a reusable `NewRedisClient` helper so that I can register multiple named Redis clients without duplicating env parsing.

## Current State / Reference

### What Exists
- `core/cache/cache.go` — private `newRedisClient()` creates Redis client from `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD` env vars. **Does NOT read `REDIS_DB`.**
- `core/session/factory.go` — inline Redis client creation. **Does NOT read `REDIS_DB`.**
- `rapidgo-starter/app/providers/redis_provider.go` — registers single `"redis"` singleton. **Does NOT read `REDIS_DB`.**

### What Works Well
- The `Store` interface pattern in `core/cache` is clean.
- Miniredis is already in `go.mod` for hermetic tests.

### What Needs Improvement
- `REDIS_DB` must be honored (silent footgun today).
- Need an exported helper for multi-client registration.

## Proposed Approach

1. **Change 1 (bugfix)**: Add `REDIS_DB` support to `newRedisClient()` in `core/cache/cache.go` and inline client creation in `core/session/factory.go`. Default remains `0`, so existing apps see zero behavior change.

2. **Change 2 (helper)**: Export `NewRedisClient(dbOverride *int)` in `core/cache` that builds a `*redis.Client` from env with optional DB override. Internal `newRedisClient()` refactored to call this. Apps can use it to register named clients without duplicating env parsing.

## Edge Cases & Risks

- [x] `REDIS_DB` out of range (< 0 or > 15) — clamp/ignore, use default 0
- [x] `REDIS_DB` non-integer — ignore, use default 0
- [x] Backwards compatibility — default is still DB 0, zero breakage

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Feature #32 — Caching | Feature | ✅ Done |
| Feature #20 — Session Management | Feature | ✅ Done |
| `redis/go-redis/v9` | External | ✅ Available |
| `alicebob/miniredis/v2` | External | ✅ Available |

## Open Questions

- [x] Option 2a (helper in `core/cache`) vs 2b (new `core/redis` pkg)? → **2a** — simpler for a patch release, no new packages.
- [x] Should `NewRedisClient` expose pool/timeout config? → **Yes** — `REDIS_POOL_SIZE`, `REDIS_DIAL_TIMEOUT`, `REDIS_READ_TIMEOUT`, `REDIS_WRITE_TIMEOUT`.

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-04-20 | Option 2a — helper in `core/cache` | Simpler for patch release, avoids new package |
| 2026-04-20 | Scope: Changes 1-2 only, defer pubsub to future minor | pubsub is a new feature, not appropriate for v2.7.3 |
| 2026-04-20 | Expose pool/timeout env vars | Production-grade defaults without requiring code changes |

## Discussion Complete ✅

**Summary**: Fix `REDIS_DB` silent ignore bug and export `NewRedisClient()` helper with DB override + pool/timeout config from env. Ship as v2.7.3 patch.
**Completed**: 2026-04-20
**Next**: Create architecture doc → `66-redis-db-multi-client-architecture.md`
