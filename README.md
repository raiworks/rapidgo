# RapidGo

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Features](https://img.shields.io/badge/Features-56-success)](FEATURES.md)

**The Laravel of Go.** A batteries-included web framework with 56 features, built on [Gin](https://github.com/gin-gonic/gin) + [GORM](https://gorm.io) + [Cobra](https://cobra.dev).

---

## Why RapidGo?

Go has amazing HTTP routers (Gin, Echo, Fiber) — but **no full application framework**. Every Go project starts by assembling 15–20 packages for auth, ORM, sessions, queues, mail, validation, and more. RapidGo ships all of that out of the box.

**RapidGo is NOT another HTTP router.** It's built *on top of* Gin — the same way Laravel is built on Symfony, NestJS on Express, and Rails on Rack.

| What RapidGo Is | What RapidGo Is NOT |
|---|---|
| Full application framework (like Laravel/NestJS/Django) | An HTTP router (that's Gin's job — RapidGo uses Gin internally) |
| Batteries-included for real apps | A microservice mesh (use Go Kit, Kratos, or Dapr for that) |
| Convention-over-configuration | A minimal toolkit that requires assembly |
| 56 shipped features, production-ready | An experiment or proof-of-concept |

---

## Feature Highlights

### Core Application
- ✅ **Dependency injection container** — singleton and transient bindings with type-safe resolution
- ✅ **Service providers** — Register/Boot lifecycle pattern for modular bootstrapping
- ✅ **Configuration** — `.env` loading via godotenv with environment detection
- ✅ **Structured logging** — JSON logging via Go's `log/slog`
- ✅ **Error handling** — centralized middleware with JSON/HTML responses
- ✅ **Plugin / module system** — extensible plugin architecture with route, command, and event hooks

### HTTP & Routing
- ✅ **Gin-based router** — route groups, resource routes, named routes, route model binding
- ✅ **MVC controllers** — `ResourceController` interface with 7 CRUD actions
- ✅ **Middleware pipeline** — registry with aliases and groups (auth, CORS, CSRF, rate limit, sessions, request ID)
- ✅ **Input validation** — built-in engine + struct-based via `go-playground/validator`
- ✅ **API response helpers** — standardized success/error/paginated envelope
- ✅ **Views & templates** — `html/template` with layouts, partials, and template functions
- ✅ **WebSocket** — upgrade, hub, broadcast, rooms/channels via `coder/websocket`
- ✅ **API versioning** — version-prefixed route groups
- ✅ **GraphQL** — handler with GraphiQL playground via `graphql-go`
- ✅ **Static file serving** — CSS, JS, images

### Data & Database
- ✅ **GORM ORM** — PostgreSQL, MySQL, SQLite with models, relationships, hooks
- ✅ **Migrations** — schema management with up/down via CLI
- ✅ **Seeders** — interface-based registry with `RunAll()` and `RunByName()`
- ✅ **Transactions** — auto, manual, and nested transaction patterns
- ✅ **Pagination** — configurable page-based helper
- ✅ **Soft deletes** — `DeletedAt` field with `WithTrashed` and `OnlyTrashed` scopes
- ✅ **Read/write splitting** — separate read replica connections

### Security & Authentication
- ✅ **JWT authentication** — stateless token auth via `golang-jwt`
- ✅ **Session-based auth** — stateful with 5 backends (database, Redis, file, memory, cookie)
- ✅ **OAuth2 / social login** — Google, GitHub, Facebook, and custom providers via `x/oauth2`
- ✅ **TOTP two-factor auth** — with bcrypt-hashed backup codes via `pquerna/otp`
- ✅ **CSRF protection** — double-submit cookie pattern with per-request tokens
- ✅ **CORS** — per-origin, per-method, per-header configuration
- ✅ **Rate limiting** — token bucket with per-IP and per-route limits via `ulule/limiter`
- ✅ **Crypto utilities** — AES-256-GCM encryption, bcrypt hashing, HMAC-SHA256, secure random tokens
- ✅ **Audit logging** — who did what, when, with structured audit trail

### Infrastructure
- ✅ **Queue workers** — background jobs with 4 drivers (database, Redis, memory, sync)
- ✅ **Task scheduler** — cron-based scheduling via `robfig/cron`
- ✅ **Event system** — publish-subscribe with sync/async dispatch
- ✅ **Caching** — Redis, in-memory, and file-based with TTL support
- ✅ **Mail** — SMTP email via `go-mail`
- ✅ **File storage** — local filesystem and Amazon S3 via unified interface
- ✅ **i18n** — JSON-based translation files with locale detection
- ✅ **Prometheus metrics** — request duration, status codes, custom counters

### CLI & Developer Experience
- ✅ **Code generation** — `make:controller`, `make:model`, `make:service`, `make:provider`, `make:migration`
- ✅ **Database CLI** — `migrate`, `migrate:rollback`, `migrate:status`, `db:seed`
- ✅ **Server CLI** — `serve` with `--mode` flag (web, api, ws, all)
- ✅ **Worker CLI** — `work` with `--queues`, `--workers`, `--timeout`
- ✅ **Scheduler CLI** — `schedule:run`
- ✅ **Admin panel scaffolding** — generates admin CRUD controllers and views

### Deployment
- ✅ **Graceful shutdown** — signal handling with connection draining
- ✅ **Health checks** — liveness and readiness probe endpoints
- ✅ **Docker** — multi-stage Dockerfile + docker-compose
- ✅ **Caddy integration** — optional auto-HTTPS reverse proxy
- ✅ **Multi-port serving** — service mode for Web, API, WebSocket on separate ports

> **[See all 56 features with package paths →](FEATURES.md)**

---

## Quick Comparison

RapidGo vs HTTP routers — different categories entirely:

| Feature | Gin | Echo | Fiber | **RapidGo** |
|---|---|---|---|---|
| HTTP Router | ✅ | ✅ | ✅ | ✅ (via Gin) |
| DI Container | ❌ | ❌ | ❌ | ✅ |
| ORM + Migrations | ❌ | ❌ | ❌ | ✅ |
| Auth (JWT + Sessions) | ❌ | ❌ | ❌ | ✅ |
| OAuth2 + TOTP 2FA | ❌ | ❌ | ❌ | ✅ |
| Queue Workers | ❌ | ❌ | ❌ | ✅ |
| Task Scheduler | ❌ | ❌ | ❌ | ✅ |
| Event System | ❌ | ❌ | ❌ | ✅ |
| Plugin System | ❌ | ❌ | ❌ | ✅ |
| GraphQL | ❌ | ❌ | ❌ | ✅ |
| WebSocket Rooms | ❌ | ❌ | ❌ | ✅ |
| Cache (3 backends) | ❌ | ❌ | ❌ | ✅ |
| Mail | ❌ | ❌ | ❌ | ✅ |
| File Storage (S3) | ❌ | ❌ | ❌ | ✅ |
| Prometheus Metrics | ❌ | ❌ | ❌ | ✅ |
| CLI Scaffolding | ❌ | ❌ | ❌ | ✅ |
| Audit Logging | ❌ | ❌ | ❌ | ✅ |
| i18n | ❌ | ❌ | ❌ | ✅ |

> **[See the full comparison →](COMPARISON.md)**

---

## Architecture

```
Request → Middleware Pipeline → Router (Gin) → Controller → Service → Model → Database
               ↑                                                         ↓
          (auth, CSRF, CORS,                                      (PostgreSQL,
           rate-limit, session,                                    MySQL, SQLite)
           request-id, metrics)
```

**Pattern**: MVC + Services + Helpers

- **Controllers** — HTTP concerns only (parse request, call service, return response)
- **Services** — business logic and domain rules (no HTTP objects)
- **Models** — GORM data schemas, relationships, hooks
- **Helpers** — stateless utility functions across all layers

**Principles**: convention over configuration · explicit over magic · composition over inheritance · fail fast · single responsibility

---

## Quick Start

```bash
# Install the CLI
go install github.com/RAiWorks/RapidGo/v2/cmd/rapidgo@latest

# Create a new project
rapidgo new myapp
cd myapp
cp .env.example .env

# Start the server
go run cmd/main.go serve

# Generate code
rapidgo make:controller UserController
rapidgo make:model User
rapidgo make:service UserService
rapidgo make:migration create_posts_table

# Database operations
rapidgo migrate
rapidgo db:seed

# Background workers
rapidgo work --queues=default,emails --workers=4
rapidgo schedule:run
```

Or clone the starter: [RapidGo-starter](https://github.com/RAiWorks/RapidGo-starter)

---

## Technology Stack

| Component | Library | Version |
|---|---|---|
| Language | Go | 1.25+ |
| HTTP Router | [Gin](https://github.com/gin-gonic/gin) | v1.12.0 |
| ORM | [GORM](https://gorm.io) | v1.31.1 |
| CLI | [Cobra](https://github.com/spf13/cobra) | v1.10.2 |
| JWT | [golang-jwt](https://github.com/golang-jwt/jwt) | v5.3.1 |
| WebSocket | [coder/websocket](https://github.com/coder/websocket) | v1.8.14 |
| Redis | [go-redis](https://github.com/redis/go-redis) | v9.18.0 |
| Rate Limiting | [ulule/limiter](https://github.com/ulule/limiter) | v3.11.2 |
| Email | [go-mail](https://github.com/wneessen/go-mail) | v0.7.2 |
| S3 Storage | [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) | v1.41.3 |
| Validation | [validator](https://github.com/go-playground/validator) | v10.30.1 |
| Scheduler | [robfig/cron](https://github.com/robfig/cron) | v3.0.1 |
| TOTP | [pquerna/otp](https://github.com/pquerna/otp) | v1.5.0 |
| GraphQL | [graphql-go](https://github.com/graphql-go/graphql) | v0.8.1 |
| OAuth2 | [x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) | v0.35.0 |
| Metrics | [prometheus](https://github.com/prometheus/client_golang) | v1.23.2 |
| Config | [godotenv](https://github.com/joho/godotenv) | v1.5.1 |
| Logging | [slog](https://pkg.go.dev/log/slog) | stdlib |
| Crypto | [x/crypto](https://pkg.go.dev/golang.org/x/crypto) | v0.48.0 |

---

## Package Index

| Package | Import Path | Purpose |
|---------|-------------|---------|
| app | `core/app` | Application lifecycle and bootstrapping |
| audit | `core/audit` | Audit logging with AuditLog model |
| auth | `core/auth` | JWT authentication |
| cache | `core/cache` | File, Redis, and in-memory caching |
| cli | `core/cli` | Cobra CLI with scaffold commands |
| config | `core/config` | Configuration and environment loading |
| container | `core/container` | IoC service container and providers |
| crypto | `core/crypto` | AES-256-GCM, HMAC-SHA256, bcrypt, secure tokens |
| errors | `core/errors` | Error handling utilities |
| events | `core/events` | Pub-sub event dispatcher |
| graphql | `core/graphql` | GraphQL server with GraphiQL playground |
| health | `core/health` | Health check endpoints (liveness + readiness) |
| i18n | `core/i18n` | JSON-based localization |
| logger | `core/logger` | Structured logging via slog |
| mail | `core/mail` | SMTP email via go-mail |
| metrics | `core/metrics` | Prometheus metrics collection |
| middleware | `core/middleware` | Middleware registry (CORS, CSRF, rate limit, auth, session, request ID, recovery) |
| oauth | `core/oauth` | OAuth2 / social login providers |
| plugin | `core/plugin` | Plugin / module system |
| queue | `core/queue` | Background job queue (4 drivers: database, Redis, memory, sync) |
| router | `core/router` | Gin-based HTTP router with groups, resources, named routes |
| scheduler | `core/scheduler` | Cron-based task scheduling |
| server | `core/server` | HTTP server with graceful shutdown and multi-port serving |
| service | `core/service` | Service mode flags (Web, API, WS, All) |
| session | `core/session` | Session management (5 backends: DB, Redis, file, memory, cookie) |
| storage | `core/storage` | File storage (local disk + Amazon S3) |
| totp | `core/totp` | TOTP two-factor authentication with backup codes |
| validation | `core/validation` | Request validation engine |
| websocket | `core/websocket` | WebSocket support with rooms and channels |
| database | `database/` | Connection, transactions, resolver, read/write splitting |
| migrations | `database/migrations` | Migration engine with up/down |
| models | `database/models` | BaseModel with soft deletes and query scopes |
| seeders | `database/seeders` | Seeder engine with interface-based registry |

## Hook System

The `core/cli` package provides 6 hooks for wiring application code to the framework:

```go
cli.SetBootstrap(func(a *app.App, mode service.Mode) {
    a.Register(&providers.ConfigProvider{})
    a.Register(&providers.DatabaseProvider{})
    // ...
})

cli.SetRoutes(func(r *router.Router, c *container.Container, mode service.Mode) {
    routes.RegisterWeb(r)
    routes.RegisterAPI(r)
})

cli.SetJobRegistrar(jobs.RegisterJobs)
cli.SetScheduleRegistrar(schedule.RegisterSchedule)
cli.SetModelRegistry(models.All)
cli.SetSeeder(func(db *gorm.DB, name string) error {
    if name != "" {
        return seeders.RunByName(db, name)
    }
    return seeders.RunAll(db)
})
```

---

## Documentation

| Resource | Description |
|---|---|
| **[Complete Feature List](FEATURES.md)** | All 56 features with package paths |
| **[Framework Comparison](COMPARISON.md)** | RapidGo vs Gin, Echo, Fiber, Go Kit |
| **[Framework Reference](docs/framework/README.md)** | 59 RFC-style reference documents |
| **[Architecture Overview](docs/framework/architecture/overview.md)** | System design and patterns |
| **[Getting Started Guide](docs/framework/guides/getting-started.md)** | First project setup |
| **[Project Context](docs/project-context.md)** | Architecture and tech stack |

---

## License

MIT — Copyright (c) 2026 [RAi Works](https://rai.works)
