# 🏗️ Architecture: Redis DB Selection & Multi-Client Helper

> **Feature**: `66` — Redis DB Selection & Multi-Client Helper
> **Discussion**: [`66-redis-db-multi-client-discussion.md`](66-redis-db-multi-client-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-04-20

---

## Overview

Fix the `REDIS_DB` env var being silently ignored across `core/cache` and `core/session`, and export a reusable `NewRedisClient()` helper that reads all Redis env vars (host, port, password, DB, pool size, timeouts) with an optional DB override. This allows apps to register multiple named Redis clients without duplicating env parsing.

## File Structure

```
core/cache/
├── cache.go            # MODIFY — refactor newRedisClient() to call NewRedisClient(nil)
├── redis.go            # No changes
├── redis_helper.go     # NEW — exported NewRedisClient(dbOverride *int) + env helpers
├── redis_helper_test.go # NEW — tests for NewRedisClient, env parsing
├── cache_test.go       # No changes (existing tests still pass)
└── file.go             # No changes

core/session/
└── factory.go          # MODIFY — Redis branch uses cache.NewRedisClient() instead of inline

database/seeders/
└── (no changes)
```

## Component Design

### NewRedisClient (core/cache/redis_helper.go)

**Responsibility**: Build a `*redis.Client` from env vars with optional DB override.
**Package**: `core/cache`
**File**: `core/cache/redis_helper.go`

```
Exported API:
└── NewRedisClient(dbOverride *int) → *redis.Client
```

#### Environment Variables Read

| Env Var | Type | Default | Description |
|---|---|---|---|
| `REDIS_HOST` | string | `"localhost"` | Redis server host |
| `REDIS_PORT` | string | `"6379"` | Redis server port |
| `REDIS_PASSWORD` | string | `""` | Redis AUTH password |
| `REDIS_DB` | int | `0` | Logical database number (0-15) |
| `REDIS_POOL_SIZE` | int | `10` | Max connections in pool |
| `REDIS_DIAL_TIMEOUT` | duration | `5s` | Connection timeout |
| `REDIS_READ_TIMEOUT` | duration | `3s` | Read timeout |
| `REDIS_WRITE_TIMEOUT` | duration | `3s` | Write timeout |

#### Logic

```go
func NewRedisClient(dbOverride *int) *redis.Client {
    db := envInt("REDIS_DB", 0)
    if dbOverride != nil {
        db = *dbOverride
    }
    return redis.NewClient(&redis.Options{
        Addr:         envStr("REDIS_HOST", "localhost") + ":" + envStr("REDIS_PORT", "6379"),
        Password:     os.Getenv("REDIS_PASSWORD"),
        DB:           db,
        PoolSize:     envInt("REDIS_POOL_SIZE", 10),
        DialTimeout:  envDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
        ReadTimeout:  envDuration("REDIS_READ_TIMEOUT", 3*time.Second),
        WriteTimeout: envDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
    })
}
```

#### Internal Helpers (unexported)

```go
func envStr(key, fallback string) string
func envInt(key string, fallback int) int       // returns fallback if value is non-numeric or negative
func envDuration(key string, fallback time.Duration) time.Duration
```

**`envInt` behavior**: If `REDIS_DB` (or any int env var) is empty, non-numeric, or negative, the fallback is returned. This prevents negative DB indices which Redis does not support.

### Changes to cache.go

Replace `newRedisClient()` body to delegate to `NewRedisClient(nil)`:

```go
func newRedisClient() (*redis.Client, error) {
    return NewRedisClient(nil), nil
}
```

### Changes to session/factory.go

Replace the inline Redis client creation (lines 36-49) with:

```go
case "redis":
    client := cache.NewRedisClient(nil)
    prefix := os.Getenv("SESSION_PREFIX")
    if prefix == "" {
        prefix = "session:"
    }
    return NewRedisStore(client, prefix), nil
```

This adds an import of `github.com/raiworks/rapidgo/v2/core/cache` to `session/factory.go` and removes the `redis/go-redis/v9` direct import (since `cache.NewRedisClient` handles it).

## Data Flow

```
App startup → Provider calls cache.NewRedisClient(&db) → reads env → returns *redis.Client
                                                              ↓
                                                    Registered in container
                                                              ↓
                                              App resolves via container.MustMake
```

## Configuration

No new env vars required. `REDIS_DB` already exists in `.env.example` — it just works now.

New **optional** env vars (documented, have sensible defaults):

| Key | Type | Default | Description |
|---|---|---|---|
| `REDIS_POOL_SIZE` | int | `10` | Max connections in pool |
| `REDIS_DIAL_TIMEOUT` | duration string | `5s` | Connection dial timeout |
| `REDIS_READ_TIMEOUT` | duration string | `3s` | Per-command read timeout |
| `REDIS_WRITE_TIMEOUT` | duration string | `3s` | Per-command write timeout |

## Security Considerations

- `REDIS_PASSWORD` remains env-only, never logged or exposed.
- DB override is `*int` (not user input from HTTP), so no injection risk.
- Pool size capped at sensible defaults to prevent resource exhaustion.

## Trade-offs & Alternatives

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| Option 2a — helper in `core/cache` | No new packages, minimal diff, backward compatible | `cache` package gains Redis-specific helpers | ✅ Selected |
| Option 2b — new `core/redis` package | Cleaner separation, future home for pubsub | New package for a patch release is heavy | ❌ Deferred |

## Next

Create tasks doc → `66-redis-db-multi-client-tasks.md`
