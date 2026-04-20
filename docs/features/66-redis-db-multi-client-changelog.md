# 📝 Changelog: Redis DB Selection & Multi-Client Helper

> **Feature**: `66` — Redis DB Selection & Multi-Client Helper
> **Branch**: `feature/66-redis-db-multi-client`
> **Started**: 2026-04-20
> **Completed**: —

---

## Log

### 2026-04-20

- **Added**: `core/cache/redis_helper.go` — exported `NewRedisClient(dbOverride *int)` with env helpers (`envStr`, `envInt`, `envDuration`)
  - Reads `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`, `REDIS_POOL_SIZE`, `REDIS_DIAL_TIMEOUT`, `REDIS_READ_TIMEOUT`, `REDIS_WRITE_TIMEOUT`
  - `envInt` returns fallback for empty, non-numeric, or negative values
- **Changed**: `core/cache/cache.go` — `newRedisClient()` now delegates to `NewRedisClient(nil)`
- **Changed**: `core/session/factory.go` — Redis branch uses `cache.NewRedisClient(nil)` instead of inline creation. Import changed from `redis/go-redis/v9` to `core/cache`.
- **Added**: `core/cache/redis_helper_test.go` — 12 test functions covering DB selection, override, invalid input, pool/timeout, host/port
- **Fixed**: `REDIS_DB` env var now honored (was silently ignored since initial implementation)

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
