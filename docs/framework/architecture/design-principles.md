---
title: "Design Principles"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Design Principles

## Abstract

This document explains the core design decisions and philosophy behind
the framework. Understanding these principles helps developers work
with the framework's conventions rather than against them, and guides
contributors when adding new features.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Laravel DX, Go Performance](#2-laravel-dx-go-performance)
3. [Driver/Interface Pattern](#3-driverinterface-pattern)
4. [Service Container & Providers](#4-service-container--providers)
5. [Convention Over Configuration](#5-convention-over-configuration)
6. [Zero-Dependency Built-ins](#6-zero-dependency-built-ins)
7. [Optional Features](#7-optional-features)
8. [Environment-Based Behavior](#8-environment-based-behavior)
9. [Separation of Concerns](#9-separation-of-concerns)
10. [Security Considerations](#10-security-considerations)
11. [References](#11-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. Laravel DX, Go Performance

The framework's mission is to combine **Laravel-style developer
experience with Go performance**. This means:

- **Familiar patterns** — MVC, service container, service providers,
  middleware groups, resource routes, named routes, and flash messages
  follow conventions established by Laravel and CodeIgniter.
- **Go idioms** — Under the hood, the framework uses Go interfaces,
  goroutines, channels, and the standard library where appropriate.
  It does not try to replicate PHP patterns verbatim.
- **Batteries included** — Common web application needs (auth, sessions,
  validation, caching, mail, file storage) are provided out of the box.

## 3. Driver/Interface Pattern

Every subsystem that supports multiple backends **MUST** define a
common Go interface. Implementations (drivers) are swappable via
configuration without changing application code.

This pattern is applied consistently across:

| Subsystem | Interface | Drivers |
|-----------|-----------|---------|
| Sessions | `session.Store` | Database, Redis, File, Memory, Cookie |
| Caching | `cache.Store` | Redis, Memory |
| File Storage | `storage.Driver` | Local filesystem, Amazon S3 |
| Database | GORM `Dialector` | PostgreSQL, MySQL, SQLite |

**Why this matters:**

- A developer can start with `memory` sessions in development and
  switch to `redis` in production by changing one `.env` variable.
- Custom drivers can be created by implementing the interface and
  registering them via a service provider.
- Testing becomes easier — use in-memory implementations for fast,
  isolated tests.

## 4. Service Container & Providers

The **service container** is the foundation for extensibility. It
provides dependency injection (DI) through a simple register/resolve
pattern — similar to Laravel's IoC container or CI4's Services class.

### Why a container?

- **Swappability** — Any registered service can be replaced by a custom
  implementation. Override the default mailer, cache, or session store
  without modifying framework code.
- **Testability** — Swap real services with mocks or stubs during testing.
- **Lifecycle management** — Singletons are created once and shared.
  Transient bindings create a new instance on each resolution.

### Provider lifecycle

Service providers **MUST** follow a two-phase lifecycle:

1. **Register** — Bind factories into the container. No other services
   **SHOULD** be resolved during this phase.
2. **Boot** — Run after all providers are registered. Services **MAY**
   resolve other services here (e.g., register event listeners that
   depend on the event dispatcher).

Every major framework component has its own provider:
`DatabaseProvider`, `SessionProvider`, `CacheProvider`, `MailProvider`,
`EventProvider`.

See: [Service Container](../core/service-container.md),
[Service Providers](../core/service-providers.md)

## 5. Convention Over Configuration

The framework relies on consistent naming and directory conventions so
developers spend less time configuring and more time building:

### Directory conventions

| Directory | Purpose |
|-----------|---------|
| `cmd/` | Application entrypoint (`main.go`) |
| `core/` | Framework internals — container, router, middleware, session, etc. |
| `app/` | User-managed code — providers, services, helpers |
| `http/` | Controllers, request structs, response helpers |
| `routes/` | Route definitions (`web.go`, `api.go`) |
| `database/` | Connection, migrations, seeders, models |
| `resources/` | Views, translations, static assets |
| `storage/` | Runtime data — uploads, cache, sessions, logs |
| `tests/` | Unit and integration tests |

### Naming conventions

- Controllers **SHOULD** be named `{Resource}Controller`
  (e.g., `UserController`, `PostController`).
- Services **SHOULD** be named `{Resource}Service`
  (e.g., `UserService`, `OrderService`).
- Models **SHOULD** match the database entity name in singular form
  (e.g., `User`, `Post`).
- Providers **SHOULD** be named `{Feature}Provider`
  (e.g., `PaymentProvider`, `NotificationProvider`).

### Code generation

The `make:*` CLI commands generate files in the correct directories
with the correct naming conventions, eliminating boilerplate.

See: [Project Structure](project-structure.md),
[Code Generation](../cli/code-generation.md)

## 6. Zero-Dependency Built-ins

The framework includes **built-in implementations** for common needs
that require zero external dependencies:

- **Validation engine** (`core/validation/`) — Required, MinLength,
  MaxLength, Email, URL, IP, Matches, In, Confirmed. No external
  library needed for common validation scenarios.
- **Crypto utilities** (`core/crypto/`) — RandomBytes, RandomHex,
  SHA256Hash, HMACSign/Verify, AES-256-GCM Encrypt/Decrypt. Built
  entirely on Go's `crypto/*` standard library packages.
- **Helpers** (`app/helpers/`) — Password hashing (bcrypt), slugify,
  truncate, string manipulation, time formatting, struct-to-map
  conversion.

For advanced scenarios, the framework also integrates with external
libraries (e.g., `go-playground/validator` for struct-based validation
via Gin's binding system).

## 7. Optional Features

Not every deployment needs every feature. The following components
**MAY** be enabled or disabled based on project requirements:

- **Docker** — Not required for running the framework. Small apps can
  run directly as a compiled binary.
- **Caddy** — Optional as either an embedded Go library or an external
  reverse proxy. Provides automatic HTTPS via Let's Encrypt.
- **Events / Hooks** — A lightweight pub-sub system. Can be skipped
  for simple applications.
- **Localization / i18n** — Translation support via JSON files. Only
  needed for multi-language applications.

Optional features **MUST NOT** cause errors or panics when disabled.

## 8. Environment-Based Behavior

The framework **MUST** support three environments controlled by the
`APP_ENV` variable:

| Environment | `APP_ENV` | Behavior |
|-------------|-----------|----------|
| Development | `development` | Verbose logging, detailed error messages, debug mode |
| Production | `production` | JSON logging, generic error messages, security hardened |
| Testing | `testing` | In-memory backends, middleware skipped where appropriate |

Environment detection helpers (`IsProduction()`, `IsDevelopment()`,
`IsTesting()`, `IsDebug()`) **MUST** be used to control:

- Error detail level (detailed in dev, generic in production)
- Logger output format (pretty-print in dev, JSON in production)
- Middleware activation (rate limiting may be skipped in testing)
- Session/cache driver defaults (memory in dev, Redis in production)

See: [Configuration](../core/configuration.md)

## 9. Separation of Concerns

Each layer has a clear, single responsibility:

- **Controllers** handle HTTP concerns only: parse request, call
  service, return response. They **MUST NOT** contain business logic.
- **Services** contain business logic and domain rules. They **MUST NOT**
  access HTTP request/response objects directly.
- **Models** define data schema and relationships. Complex queries
  **SHOULD** live in services, not models.
- **Helpers** are pure, stateless utility functions. They **MUST NOT**
  maintain state or depend on request context.
- **Middleware** handles cross-cutting concerns (auth, CSRF, logging).
  It **MUST NOT** contain business logic.

This separation ensures testability — services can be tested without
HTTP, models without services, helpers without any dependencies.

## 10. Security Considerations

Design principles directly impact security:

- The driver/interface pattern ensures session and cache backends can
  be hardened per environment without code changes.
- Environment-based behavior prevents accidental exposure of debug
  information in production.
- Zero-dependency crypto utilities reduce the attack surface from
  third-party dependencies.
- The service container enables swapping security components (auth,
  encryption) without framework modification.

## 11. References

- [Architecture Overview](overview.md)
- [Project Structure](project-structure.md)
- [Service Container](../core/service-container.md)
- [Configuration](../core/configuration.md)
- [Crypto Utilities](../security/crypto.md)
- [Requests & Validation](../http/requests-validation.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
