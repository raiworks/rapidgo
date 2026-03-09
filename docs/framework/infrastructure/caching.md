---
title: "Caching"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Caching

## Abstract

This document covers the driver-based cache system — the cache
interface, Redis and memory backends, configuration, and usage
patterns.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Configuration](#2-configuration)
3. [Cache Interface](#3-cache-interface)
4. [Redis Cache](#4-redis-cache)
5. [Memory Cache](#5-memory-cache)
6. [Usage Patterns](#6-usage-patterns)
7. [Provider Registration](#7-provider-registration)
8. [Security Considerations](#8-security-considerations)
9. [References](#9-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Cache store** — A backend implementing the `Store` interface for
  key-value storage with TTL.

## 2. Configuration

`.env`:

```env
CACHE_DRIVER=redis
CACHE_PREFIX=app:
CACHE_TTL=3600
```

## 3. Cache Interface

All backends implement this contract in `core/cache/`:

```go
package cache

import "time"

type Store interface {
    Get(key string) (string, error)
    Set(key string, value string, ttl time.Duration) error
    Delete(key string) error
    Flush() error
}
```

## 4. Redis Cache

```go
type RedisCache struct {
    Client *redis.Client
    Prefix string
}

func (c *RedisCache) Get(key string) (string, error) {
    val, err := c.Client.Get(context.Background(), c.Prefix+key).Result()
    if err == redis.Nil {
        return "", nil
    }
    return val, err
}

func (c *RedisCache) Set(key, value string, ttl time.Duration) error {
    return c.Client.Set(context.Background(), c.Prefix+key, value, ttl).Err()
}

func (c *RedisCache) Delete(key string) error {
    return c.Client.Del(context.Background(), c.Prefix+key).Err()
}

func (c *RedisCache) Flush() error {
    return c.Client.FlushDB(context.Background()).Err()
}
```

## 5. Memory Cache

For development and testing — data is lost on restart:

```go
type MemoryCache struct {
    mu    sync.RWMutex
    items map[string]memCacheEntry
}

func NewMemoryCache() *MemoryCache {
    return &MemoryCache{items: make(map[string]memCacheEntry)}
}
```

The memory cache automatically handles TTL by checking expiration
on every `Get`.

## 6. Usage Patterns

### Cache a Database Query

```go
val, _ := cacheStore.Get("users:count")
if val == "" {
    var count int64
    db.Model(&models.User{}).Count(&count)
    cacheStore.Set("users:count", strconv.FormatInt(count, 10), 5*time.Minute)
    val = strconv.FormatInt(count, 10)
}
```

### Invalidate on Write

```go
func (s *UserService) Create(...) (*models.User, error) {
    // ... create user ...
    cacheStore.Delete("users:count")
    return user, nil
}
```

### Cache Key Patterns

| Pattern | Example |
|---------|---------|
| Entity count | `users:count` |
| Single entity | `users:42` |
| Query result | `posts:published:page:1` |
| Config value | `settings:site_name` |

## 7. Provider Registration

```go
type CacheProvider struct{}

func (p *CacheProvider) Register(c *container.Container) {
    c.Singleton("cache", func(c *container.Container) interface{} {
        return cache.NewMemoryCache() // swap for Redis in production
    })
}

func (p *CacheProvider) Boot(c *container.Container) {}
```

## 8. Security Considerations

- Cache keys **MUST NOT** include unvalidated user input to prevent
  cache poisoning.
- Sensitive data in cache **SHOULD** be encrypted or have short TTLs.
- `Flush()` clears the entire cache and **SHOULD** be restricted to
  admin operations.
- Redis connections **MUST** use authentication in production.

## 9. References

- [Configuration](../core/configuration.md)
- [Service Providers](../core/service-providers.md)
- [Data Flow Diagram](../architecture/diagrams/data-flow.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
