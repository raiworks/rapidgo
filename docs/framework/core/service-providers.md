---
title: "Service Providers"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Service Providers

## Abstract

Service providers are the central place to configure and register
services into the container. Every major framework component — database,
sessions, cache, mail, events — has its own provider. This document
covers the provider interface, lifecycle, built-in providers, the App
bootstrap struct, and how to create custom providers.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Provider Interface](#2-provider-interface)
3. [Provider Lifecycle](#3-provider-lifecycle)
4. [Built-in Providers](#4-built-in-providers)
5. [App Struct (Bootstrap)](#5-app-struct-bootstrap)
6. [Registration Order in main.go](#6-registration-order-in-maingo)
7. [Creating a Custom Provider](#7-creating-a-custom-provider)
8. [Register vs Boot](#8-register-vs-boot)
9. [Security Considerations](#9-security-considerations)
10. [References](#10-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. Provider Interface

```go
package container

// Provider defines the lifecycle hooks for service registration.
type Provider interface {
    // Register binds services into the container.
    // Called before the application boots. No other services
    // should be used here — only register bindings.
    Register(c *Container)

    // Boot runs after ALL providers have been registered.
    // You can use other services here (resolve from container).
    Boot(c *Container)
}
```

Every provider **MUST** implement both `Register()` and `Boot()`.

## 3. Provider Lifecycle

The lifecycle is split into two distinct phases:

### Phase 1: Register

- Called immediately when `app.Register(provider)` is invoked.
- **MUST** only bind factories into the container (`Bind`, `Singleton`,
  `Instance`).
- **SHOULD NOT** resolve other services — they may not be registered
  yet.

### Phase 2: Boot

- Called after all providers have been registered, when `app.Boot()`
  runs.
- Providers **MAY** resolve other services from the container.
- Used for: running auto-migrations, registering event listeners,
  setting up cross-service dependencies.

```text
┌─────────────────────────────────────────┐
│           Registration Phase             │
│                                          │
│  Provider A → Register(container)        │
│  Provider B → Register(container)        │
│  Provider C → Register(container)        │
│                                          │
│  (Only bind factories — no resolution)   │
├─────────────────────────────────────────┤
│           Boot Phase                     │
│                                          │
│  Provider A → Boot(container)            │
│  Provider B → Boot(container)            │
│  Provider C → Boot(container)            │
│                                          │
│  (May resolve services, run setup)       │
└─────────────────────────────────────────┘
```

## 4. Built-in Providers

### DatabaseProvider

Registers the database connection as a singleton.

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

func (p *DatabaseProvider) Boot(c *container.Container) {
    // Run auto-migrations if enabled
}
```

- **Service name:** `"db"`
- **Type:** `*gorm.DB`
- **Resolution:** `container.MustMake[*gorm.DB](c, "db")`

### SessionProvider

Registers the session manager, which depends on the database connection.

```go
type SessionProvider struct{}

func (p *SessionProvider) Register(c *container.Container) {
    c.Singleton("session", func(c *container.Container) interface{} {
        db := container.MustMake[*gorm.DB](c, "db")
        store, _ := session.NewStore(db)
        return session.NewManager(store)
    })
}

func (p *SessionProvider) Boot(c *container.Container) {}
```

- **Service name:** `"session"`
- **Type:** `*session.Manager`
- **Depends on:** `"db"` (for database session store)

### CacheProvider

Registers the cache store.

```go
type CacheProvider struct{}

func (p *CacheProvider) Register(c *container.Container) {
    c.Singleton("cache", func(c *container.Container) interface{} {
        return cache.NewMemoryCache() // swap for Redis in production
    })
}

func (p *CacheProvider) Boot(c *container.Container) {}
```

- **Service name:** `"cache"`
- **Type:** `cache.Store`
- **Backends:** `MemoryCache` (default), `RedisCache` (production)

### MailProvider

Registers the SMTP mailer.

```go
type MailProvider struct{}

func (p *MailProvider) Register(c *container.Container) {
    c.Singleton("mail", func(c *container.Container) interface{} {
        return mail.NewMailer()
    })
}

func (p *MailProvider) Boot(c *container.Container) {}
```

- **Service name:** `"mail"`
- **Type:** `*mail.Mailer`

### EventProvider

Registers the event dispatcher.

```go
type EventProvider struct{}

func (p *EventProvider) Register(c *container.Container) {
    c.Singleton("events", func(c *container.Container) interface{} {
        return events.NewDispatcher()
    })
}

func (p *EventProvider) Boot(c *container.Container) {}
```

- **Service name:** `"events"`
- **Type:** `*events.Dispatcher`

## 5. App Struct (Bootstrap)

The `App` struct orchestrates the provider lifecycle:

```go
package app

import "yourframework/core/container"

type App struct {
    Container *container.Container
    providers []container.Provider
}

func New() *App {
    return &App{
        Container: container.New(),
    }
}

// Register adds a provider to the application.
func (a *App) Register(provider container.Provider) {
    a.providers = append(a.providers, provider)
    provider.Register(a.Container)
}

// Boot calls Boot on all registered providers.
func (a *App) Boot() {
    for _, p := range a.providers {
        p.Boot(a.Container)
    }
}

// Make resolves a service.
func (a *App) Make(name string) interface{} {
    return a.Container.Make(name)
}
```

Key behaviors:

- `Register()` immediately calls `provider.Register()` — the binding
  is available right away.
- `Boot()` iterates all providers in **registration order** — earlier
  providers boot first.
- `Make()` is a convenience proxy to `a.Container.Make()`.

## 6. Registration Order in main.go

```go
func main() {
    config.Load()

    application := app.New()

    // Built-in providers (order matters for dependencies)
    application.Register(&providers.DatabaseProvider{})    // 1. DB first
    application.Register(&providers.SessionProvider{})     // 2. Sessions need DB
    application.Register(&providers.CacheProvider{})       // 3. Cache
    application.Register(&providers.MailProvider{})         // 4. Mail
    application.Register(&providers.EventProvider{})       // 5. Events

    // User custom providers
    application.Register(&providers.PaymentProvider{})
    application.Register(&providers.NotificationProvider{})

    // Boot all providers
    application.Boot()

    // Use services
    db := container.MustMake[*gorm.DB](application.Container, "db")
    _ = db

    // Start server...
}
```

**Order matters:** `SessionProvider` depends on `"db"` — it **MUST**
be registered after `DatabaseProvider`. Since `Singleton` factories are
lazy (resolved on first `Make()`), the actual creation order is
determined by resolution, but registration order affects `Boot()`.

## 7. Creating a Custom Provider

### Step 1: Generate with CLI

```bash
framework make:provider PaymentProvider
```

This creates `app/providers/paymentprovider.go`.

### Step 2: Implement the provider

```go
package providers

import "yourframework/core/container"

type PaymentProvider struct{}

func (p *PaymentProvider) Register(c *container.Container) {
    c.Singleton("payment", func(c *container.Container) interface{} {
        // Return your payment gateway service
        return NewStripeGateway(os.Getenv("STRIPE_KEY"))
    })
}

func (p *PaymentProvider) Boot(c *container.Container) {
    // Register event listeners, run setup, etc.
}
```

### Step 3: Register in main.go

```go
application.Register(&providers.PaymentProvider{})
```

### Step 4: Use in controllers/services

```go
gateway := container.MustMake[*StripeGateway](app.Container, "payment")
gateway.Charge(amount)
```

## 8. Register vs Boot

| Aspect | Register | Boot |
|--------|----------|------|
| **When called** | During `app.Register()` | During `app.Boot()` (after all Register) |
| **Purpose** | Bind factories | Initialize, resolve dependencies |
| **May resolve services?** | **NO** | Yes |
| **May register bindings?** | Yes | **SHOULD NOT** (too late for others) |
| **Called how many times?** | Once per provider | Once per provider |

**Rule of thumb:** If it's a `Bind`, `Singleton`, or `Instance` call,
it goes in `Register()`. Everything else goes in `Boot()`.

## 9. Security Considerations

- Provider factories that access credentials (database passwords, API
  keys) **SHOULD** read them from environment variables inside the
  factory closure, not store them as struct fields.
- Custom providers **MUST NOT** log sensitive configuration values
  during registration or boot.
- The panic behavior in `DatabaseProvider` on connection failure is
  intentional — the application cannot function without a database.

## 10. References

- [Service Container](service-container.md)
- [Application Lifecycle](../architecture/application-lifecycle.md)
- [Service Container Diagram](../architecture/diagrams/service-container.md)
- [Configuration](configuration.md)
- [Code Generation](../cli/code-generation.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
