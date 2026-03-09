---
title: "Extending the Framework"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Extending the Framework

## Abstract

This guide covers how to extend and customize the framework — creating
reusable packages, swapping built-in services, adding custom middleware
and CLI commands, and overriding defaults via providers.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Creating Reusable Packages](#2-creating-reusable-packages)
3. [Swapping Built-in Services](#3-swapping-built-in-services)
4. [Adding Custom Middleware](#4-adding-custom-middleware)
5. [Adding Custom CLI Commands](#5-adding-custom-cli-commands)
6. [Overriding Defaults via Providers](#6-overriding-defaults-via-providers)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. Creating Reusable Packages

Framework extensions follow Go's package conventions. A reusable
package **SHOULD**:

- Accept interfaces, not concrete types
- Provide a service provider for container registration
- Include its own configuration keys in `.env`

```text
mypackage/
├── mypackage.go      # Core logic
├── provider.go       # Service provider
└── config.go         # Configuration helpers
```

```go
// provider.go
package mypackage

import "yourframework/core/container"

type MyPackageProvider struct{}

func (p *MyPackageProvider) Register(c *container.Container) {
    c.Singleton(func() *MyService {
        return NewMyService()
    })
}

func (p *MyPackageProvider) Boot(c *container.Container) {}
```

Register in the app:

```go
app.Register(&mypackage.MyPackageProvider{})
```

## 3. Swapping Built-in Services

The framework uses interfaces for pluggable components. To swap a
service, register a new implementation for the same interface.

### Example: Custom Cache Driver

Implement the `cache.Store` interface:

```go
type MemcachedCache struct {
    client *memcache.Client
}

func (m *MemcachedCache) Get(key string) (string, error) { /* ... */ }
func (m *MemcachedCache) Set(key string, value string, ttl time.Duration) error { /* ... */ }
func (m *MemcachedCache) Delete(key string) error { /* ... */ }
func (m *MemcachedCache) Flush() error { /* ... */ }
```

Register it in a provider:

```go
func (p *CacheProvider) Register(c *container.Container) {
    c.Singleton(func() cache.Store {
        return NewMemcachedCache(os.Getenv("MEMCACHED_HOST"))
    })
}
```

Any code requesting `cache.Store` from the container will receive
the Memcached implementation instead of the default.

## 4. Adding Custom Middleware

Create middleware following the `gin.HandlerFunc` signature:

```go
package middleware

import "github.com/gin-gonic/gin"

func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start)
        slog.Info("request",
            "method", c.Request.Method,
            "path", c.Request.URL.Path,
            "status", c.Writer.Status(),
            "duration", duration,
        )
    }
}
```

Register as an alias in the middleware registry:

```go
registry.RegisterAlias("logger", middleware.RequestLogger())
```

Use on routes:

```go
r.Use(middleware.Resolve("logger"))
```

See [Middleware](../http/middleware.md).

## 5. Adding Custom CLI Commands

Register Cobra commands in `cmd/`:

```go
var importCmd = &cobra.Command{
    Use:   "import [file]",
    Short: "Import data from a CSV file",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        // Import logic
    },
}

func init() {
    rootCmd.AddCommand(importCmd)
}
```

See [CLI Overview](../cli/cli-overview.md).

## 6. Overriding Defaults via Providers

Providers registered later override earlier bindings. To override a
framework default, register your provider after the built-in ones:

```go
// Built-in providers (registered first)
app.Register(&providers.DatabaseProvider{})
app.Register(&providers.CacheProvider{})

// Your override (registered last — wins)
app.Register(&providers.CustomCacheProvider{})
```

The container resolves to the last binding for a given type.

## 7. References

- [Service Container](../core/service-container.md)
- [Service Providers](../core/service-providers.md)
- [Middleware](../http/middleware.md)
- [CLI Overview](../cli/cli-overview.md)
- [Caching](../infrastructure/caching.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
