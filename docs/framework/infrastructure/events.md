---
title: "Events"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Events

## Abstract

This document covers the lightweight publish-subscribe event system
— the dispatcher, registering listeners, synchronous and asynchronous
dispatch, and usage patterns.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Dispatcher](#2-dispatcher)
3. [Registering Listeners](#3-registering-listeners)
4. [Dispatching Events](#4-dispatching-events)
5. [Provider Registration](#5-provider-registration)
6. [Usage Patterns](#6-usage-patterns)
7. [Security Considerations](#7-security-considerations)
8. [References](#8-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Event** — A named occurrence in the application (e.g.,
  `user.created`).
- **Listener** — A function invoked when an event is dispatched.
- **Payload** — Data passed to listeners when an event fires.

## 2. Dispatcher

```go
package events

import "sync"

type Handler func(payload interface{})

type Dispatcher struct {
    mu        sync.RWMutex
    listeners map[string][]Handler
}

func NewDispatcher() *Dispatcher {
    return &Dispatcher{listeners: make(map[string][]Handler)}
}

func (d *Dispatcher) Listen(event string, handler Handler) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.listeners[event] = append(d.listeners[event], handler)
}

func (d *Dispatcher) Dispatch(event string, payload interface{}) {
    d.mu.RLock()
    defer d.mu.RUnlock()
    for _, h := range d.listeners[event] {
        go h(payload) // async by default
    }
}

func (d *Dispatcher) DispatchSync(event string, payload interface{}) {
    d.mu.RLock()
    defer d.mu.RUnlock()
    for _, h := range d.listeners[event] {
        h(payload)
    }
}
```

### Dispatch Modes

| Method | Execution | Use Case |
|--------|-----------|----------|
| `Dispatch` | Async (goroutines) | Non-blocking side effects (email, logging) |
| `DispatchSync` | Synchronous | When ordering or error handling matters |

## 3. Registering Listeners

Register listeners at boot time, typically in a provider:

```go
events.Listen("user.created", func(payload interface{}) {
    user := payload.(*models.User)
    mailer.Send(user.Email, "Welcome!", "<h1>Welcome!</h1>")
})

events.Listen("user.created", func(payload interface{}) {
    slog.Info("new user registered",
        "email", payload.(*models.User).Email)
})
```

Multiple listeners can be registered for the same event.

## 4. Dispatching Events

From services or controllers:

```go
dispatcher.Dispatch("user.created", &user)
```

## 5. Provider Registration

```go
type EventProvider struct{}

func (p *EventProvider) Register(c *container.Container) {
    c.Singleton("events", func(c *container.Container) interface{} {
        return events.NewDispatcher()
    })
}

func (p *EventProvider) Boot(c *container.Container) {}
```

## 6. Usage Patterns

### Common Events

| Event | Payload | Typical Listeners |
|-------|---------|-------------------|
| `user.created` | `*models.User` | Welcome email, audit log |
| `user.deleted` | `*models.User` | Cleanup, notification |
| `order.placed` | `*models.Order` | Confirmation email, inventory update |
| `login.failed` | `LoginAttempt` | Security logging, lockout check |

### Decoupling Components

Events let you add behavior without modifying existing code:

```go
// In UserService — just dispatch, doesn't know about email
func (s *UserService) Create(...) (*models.User, error) {
    // ... create user ...
    s.Events.Dispatch("user.created", user)
    return user, nil
}

// In a provider — registers side effects
events.Listen("user.created", sendWelcomeEmail)
events.Listen("user.created", logNewUser)
events.Listen("user.created", updateAnalytics)
```

## 7. Security Considerations

- Async listeners run in goroutines — ensure they handle panics
  gracefully (use `recover`).
- Event payloads **SHOULD NOT** include sensitive data that shouldn't
  be accessible to all listeners.
- Security-critical events (login failures, permission changes)
  **SHOULD** always be logged.

## 8. References

- [Service Providers](../core/service-providers.md)
- [Mail](mail.md)
- [Logging](../core/logging.md)
- [Data Flow Diagram](../architecture/diagrams/data-flow.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
