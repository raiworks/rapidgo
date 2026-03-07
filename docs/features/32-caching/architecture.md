# Feature #32 — Caching: Architecture

## Component overview

```
core/cache/cache.go
    │
    ├── Store interface
    │       Get(key) (string, error)
    │       Set(key, value, ttl) error
    │       Delete(key) error
    │       Flush() error
    │
    ├── MemoryCache struct
    │       sync.RWMutex + map[string]memCacheEntry
    │       Lazy expiry on Get()
    │
    └── NewStore(driver, prefix string) (Store, error)
            Factory: reads driver name, returns prefixed store
```

## New file

| File | Purpose |
|------|---------|
| `core/cache/cache.go` | `Store` interface, `MemoryCache`, `NewStore()` factory |

## Removed

| File | Reason |
|------|--------|
| `core/cache/.gitkeep` | Replaced by real implementation |

## Types

```go
// Store defines the cache contract.
type Store interface {
    Get(key string) (string, error)
    Set(key string, value string, ttl time.Duration) error
    Delete(key string) error
    Flush() error
}

// MemoryCache is an in-process cache with TTL-based expiry.
type MemoryCache struct {
    mu    sync.RWMutex
    items map[string]memCacheEntry
}

type memCacheEntry struct {
    Value     string
    ExpiresAt time.Time
}
```

## Functions

| Function | Signature | Behaviour |
|----------|-----------|-----------|
| `NewMemoryCache()` | `func NewMemoryCache() *MemoryCache` | Returns empty in-memory store |
| `Get()` | `func (c *MemoryCache) Get(key string) (string, error)` | Returns value if exists and not expired; lazy-deletes expired entries |
| `Set()` | `func (c *MemoryCache) Set(key, value string, ttl time.Duration) error` | Stores value with absolute expiry |
| `Delete()` | `func (c *MemoryCache) Delete(key string) error` | Removes key from map |
| `Flush()` | `func (c *MemoryCache) Flush() error` | Replaces map with empty map |
| `NewStore()` | `func NewStore(driver, prefix string) (Store, error)` | Factory: `"memory"` → `MemoryCache` (with prefix wrapper), unknown → error |

## Environment variables

| Var | Default | Purpose |
|-----|---------|---------|
| `CACHE_DRIVER` | `memory` | Backend driver name |
| `CACHE_PREFIX` | `rapidgo_` | Key prefix for namespacing |
| `CACHE_TTL` | `3600` | Default TTL in seconds (informational; callers pass explicit TTL) |

## Usage

```go
store, _ := cache.NewStore("memory", "app:")
store.Set("users:count", "42", 5*time.Minute)
val, _ := store.Get("users:count") // "42"
store.Delete("users:count")
store.Flush()
```
