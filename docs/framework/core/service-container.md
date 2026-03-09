---
title: "Service Container"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Service Container

## Abstract

The service container is the foundation for dependency injection and
extensibility. It provides a register/resolve mechanism for managing
service instances, factories, and singletons — similar to Laravel's
IoC container or CI4's Services class.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Overview](#2-overview)
3. [Container Struct](#3-container-struct)
4. [Registration Methods](#4-registration-methods)
5. [Resolution Methods](#5-resolution-methods)
6. [Thread Safety](#6-thread-safety)
7. [Usage Examples](#7-usage-examples)
8. [Security Considerations](#8-security-considerations)
9. [References](#9-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Binding** — A factory function registered under a name that creates
  a service instance when resolved.
- **Singleton** — A binding that is resolved only once; subsequent
  resolutions return the cached instance.
- **Transient** — A binding that creates a new instance on every
  resolution.
- **Instance** — A pre-created object stored directly in the container.

## 2. Overview

The container provides three ways to register services and three ways
to resolve them:

| Registration | Behavior |
|-------------|----------|
| `Bind(name, factory)` | Transient — calls factory on every `Make()` |
| `Singleton(name, factory)` | Shared — calls factory once, caches result |
| `Instance(name, obj)` | Pre-created — stores object directly |

| Resolution | Return Type |
|-----------|------------|
| `Make(name)` | `interface{}` — caller must type-assert |
| `MustMake[T](c, name)` | `T` — generic typed resolution |
| `Has(name)` | `bool` — checks if service is registered |

## 3. Container Struct

```go
package container

import (
    "fmt"
    "sync"
)

type Factory func(c *Container) interface{}

type Container struct {
    mu        sync.RWMutex
    bindings  map[string]Factory
    instances map[string]interface{}
}

func New() *Container {
    return &Container{
        bindings:  make(map[string]Factory),
        instances: make(map[string]interface{}),
    }
}
```

The container maintains two internal maps:

- **`bindings`** — Maps service names to factory functions.
- **`instances`** — Maps service names to resolved singleton instances
  or pre-created objects.

## 4. Registration Methods

### Bind (Transient)

Registers a factory function that is called every time `Make()` is
invoked. Each call returns a new instance.

```go
func (c *Container) Bind(name string, factory Factory) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.bindings[name] = factory
}
```

**Use case:** Services that should not be shared (e.g., request-scoped
loggers, per-request validators).

### Singleton (Shared Instance)

Registers a factory that is called only once. The first `Make()` call
creates the instance; all subsequent calls return the same cached
instance.

```go
func (c *Container) Singleton(name string, factory Factory) {
    c.Bind(name, func(cont *Container) interface{} {
        cont.mu.RLock()
        if inst, ok := cont.instances[name]; ok {
            cont.mu.RUnlock()
            return inst
        }
        cont.mu.RUnlock()

        inst := factory(cont)
        cont.mu.Lock()
        cont.instances[name] = inst
        cont.mu.Unlock()
        return inst
    })
}
```

**Use case:** Most framework services — database connections, session
managers, cache stores, mailers, event dispatchers. These are expensive
to create and should be shared.

### Instance (Pre-created)

Registers an already-created object directly, bypassing the factory
pattern entirely.

```go
func (c *Container) Instance(name string, instance interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.instances[name] = instance
}
```

**Use case:** Configuration objects, test doubles, or services that
were created outside the container.

## 5. Resolution Methods

### Make

Resolves a service by name. Checks instances first, then bindings.
Panics if the service is not registered.

```go
func (c *Container) Make(name string) interface{} {
    c.mu.RLock()
    if inst, ok := c.instances[name]; ok {
        c.mu.RUnlock()
        return inst
    }
    factory, ok := c.bindings[name]
    c.mu.RUnlock()
    if !ok {
        panic(fmt.Sprintf("service not found: %s", name))
    }
    return factory(c)
}
```

**Resolution order:**

1. Check `instances` map → return if found.
2. Check `bindings` map → call factory if found.
3. Panic if neither exists.

### MustMake[T] (Generic Typed Resolution)

A generic helper that resolves and casts in one step. Panics if the
type assertion fails.

```go
func MustMake[T any](c *Container, name string) T {
    return c.Make(name).(T)
}
```

**Usage:**

```go
db := container.MustMake[*gorm.DB](c, "db")
mailer := container.MustMake[*mail.Mailer](c, "mail")
dispatcher := container.MustMake[*events.Dispatcher](c, "events")
```

### Has

Checks whether a service is registered (either as a binding or an
instance).

```go
func (c *Container) Has(name string) bool {
    c.mu.RLock()
    defer c.mu.RUnlock()
    _, hasBinding := c.bindings[name]
    _, hasInstance := c.instances[name]
    return hasBinding || hasInstance
}
```

## 6. Thread Safety

The container uses `sync.RWMutex` for concurrent access:

- **Read operations** (`Make`, `Has`) acquire a read lock — multiple
  goroutines can resolve services simultaneously.
- **Write operations** (`Bind`, `Singleton`, `Instance`) acquire an
  exclusive write lock.

This ensures the container is safe per Go's concurrent usage patterns.
All framework services are registered during startup (single-threaded)
and resolved during request handling (concurrent).

## 7. Usage Examples

### Registering services in a provider

```go
type DatabaseProvider struct{}

func (p *DatabaseProvider) Register(c *container.Container) {
    c.Singleton("db", func(c *container.Container) interface{} {
        db, err := database.Connect()
        if err != nil {
            panic("database connection failed: " + err.Error())
        }
        return db
    })
}
```

### Resolving in a controller

```go
func CreateUser(c *gin.Context) {
    db := container.MustMake[*gorm.DB](appContainer, "db")
    userSvc := services.NewUserService(db)
    // ...
}
```

### Swapping a service for testing

```go
func TestCreateUser(t *testing.T) {
    c := container.New()
    c.Instance("db", mockDB) // inject mock
    // ...
}
```

### Bind vs Singleton vs Instance

| Method | When Factory Runs | Use Case |
|--------|-------------------|----------|
| `Bind` | Every `Make()` call | Request-scoped, stateless |
| `Singleton` | First `Make()` only | Shared, expensive to create |
| `Instance` | Never (pre-created) | Config, test mocks |

## 8. Security Considerations

- The container stores service instances in memory. Sensitive services
  (e.g., database connections with credentials) are not exposed beyond
  the Go process boundary.
- Service factories **SHOULD NOT** log or expose credentials. Use
  environment variables accessed within the factory closure.
- The panic behavior on missing services is intentional — it fails
  fast during development. In production, all services **MUST** be
  registered before `Boot()` is called.

## 9. References

- [Service Providers](service-providers.md)
- [Application Lifecycle](../architecture/application-lifecycle.md)
- [Service Container Diagram](../architecture/diagrams/service-container.md)
- [Design Principles](../architecture/design-principles.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
