# Feature #32 — Caching: Review

## Summary

Implemented a driver-based caching subsystem with a `Store` interface, three backends (memory, file, Redis), key-prefix wrapper, and factory function.

## Files changed

| File | Action | Purpose |
|------|--------|---------|
| `core/cache/cache.go` | Updated | `Store` interface, `MemoryCache`, `prefixStore`, `NewStore()` with file + redis support |
| `core/cache/file.go` | Created | `FileCache` — disk-based cache with JSON files and TTL |
| `core/cache/redis.go` | Created | `RedisCache` — Redis-backed cache using go-redis/v9 |
| `core/cache/cache_test.go` | Updated | 25 test cases (was 10) |

## Blueprint compliance

| Blueprint item | Status | Notes |
|----------------|--------|-------|
| `Store` interface (Get/Set/Delete/Flush) | ✅ | Exact match |
| `MemoryCache` with sync.RWMutex | ✅ | Improved: lazy delete on Get() |
| `RedisCache` with go-redis/v9 | ✅ | Full implementation with TTL |
| `FileCache` with disk storage | ✅ | JSON files with TTL, lazy expiry on Get() |
| `CACHE_DRIVER` / `CACHE_PREFIX` env vars | ✅ | Factory reads driver + prefix |
| `NewMemoryCache()` constructor | ✅ | Exact match |
| `NewFileCache(path)` constructor | ✅ | Reads `CACHE_FILE_PATH` env (default: "storage/cache") |
| `NewRedisCache(client)` constructor | ✅ | Reads `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD` env |

## FileCache Details

- Each key stored as `{key}.cache` JSON file in configured directory
- JSON structure: `{"value": "...", "expires_at": "..."}` with Go time format
- Lazy expiry: expired files deleted on next `Get()` call
- `Flush()` removes all `.cache` files from directory
- `Delete()` is no-op for missing keys (idempotent)

## RedisCache Details

- Uses `github.com/redis/go-redis/v9`
- TTL handled natively by Redis SET with KEEPTTL
- `Flush()` calls `FLUSHDB` (clears current database)
- Missing keys return empty string (no error)

## Deviations

| # | Blueprint | Ours | Reason |
|---|-----------|------|--------|
| 1 | Prefix baked into each driver | `prefixStore` wrapper in factory | Cleaner separation; prefix is a cross-cutting concern |
| 2 | Get() holds RLock, doesn't clean up | MemoryCache Get() upgrades to write lock for lazy delete | Prevents unbounded memory growth from expired entries |
| 3 | Redis/File deferred | Now fully implemented | Dependencies added: go-redis/v9, miniredis/v2 (test) |

## Test results

| TC | Description | Result |
|----|-------------|--------|
| TC-01 | Set and Get returns value | ✅ PASS |
| TC-02 | Get missing key returns empty string | ✅ PASS |
| TC-03 | Get expired key returns empty string + lazy delete | ✅ PASS |
| TC-04 | Delete removes key | ✅ PASS |
| TC-05 | Flush clears all keys | ✅ PASS |
| TC-06 | Concurrent access is safe | ✅ PASS |
| TC-07 | NewStore "memory" returns valid Store | ✅ PASS |
| TC-08 | NewStore unknown driver returns error | ✅ PASS |
| TC-09 | Prefix applied to keys | ✅ PASS |
| TC-10 | Set overwrites existing key | ✅ PASS |
| TC-11 | FileCache set and get returns value | ✅ PASS |
| TC-12 | FileCache get missing key returns empty string | ✅ PASS |
| TC-13 | FileCache expired key returns empty string | ✅ PASS |
| TC-14 | FileCache delete removes key | ✅ PASS |
| TC-15 | FileCache delete missing key is no-op | ✅ PASS |
| TC-16 | FileCache flush clears all keys | ✅ PASS |
| TC-17 | FileCache set overwrites existing key | ✅ PASS |
| TC-18 | NewStore "file" returns valid Store | ✅ PASS |
| TC-19 | RedisCache set and get returns value | ✅ PASS |
| TC-20 | RedisCache get missing key returns empty string | ✅ PASS |
| TC-21 | RedisCache delete removes key | ✅ PASS |
| TC-22 | RedisCache flush clears all keys | ✅ PASS |
| TC-23 | RedisCache TTL expiry works | ✅ PASS |
| TC-24 | NewStore "redis" returns valid Store | ✅ PASS |
| TC-25 | NewStore "redis" with prefix wraps correctly | ✅ PASS |

## Regression

- All 30 packages pass (`go test ./...`)
- `go vet ./...` clean
- Dependencies: go-redis/v9, aws-sdk-go-v2, miniredis/v2 (test-only)
