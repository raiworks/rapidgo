# RapidGo

A batteries-included Go web framework with Laravel-style developer experience.

Built on [Gin](https://github.com/gin-gonic/gin), [GORM](https://gorm.io), and [Cobra](https://cobra.dev).

## Install

```bash
go get github.com/RAiWorks/RapidGo
```

## Quick Start

```bash
# Install the CLI
go install github.com/RAiWorks/RapidGo/cmd/rapidgo@latest

# Scaffold a new project
rapidgo new myapp
cd myapp
cp .env.example .env
go run cmd/main.go serve
```

Or clone the starter directly: [RapidGo-starter](https://github.com/RAiWorks/RapidGo-starter)

## Package Index

| Package | Import Path | Purpose |
|---------|-------------|---------|
| app | `core/app` | Application lifecycle and bootstrapping |
| audit | `core/audit` | Audit logging with AuditLog model |
| auth | `core/auth` | JWT authentication |
| cache | `core/cache` | File + Redis caching |
| cli | `core/cli` | Cobra CLI with scaffold commands |
| config | `core/config` | Configuration and environment loading |
| container | `core/container` | IoC service container and providers |
| crypto | `core/crypto` | AES-256-GCM, HMAC-SHA256, secure tokens |
| errors | `core/errors` | Error handling utilities |
| events | `core/events` | Pub-sub event dispatcher |
| graphql | `core/graphql` | GraphQL server integration |
| health | `core/health` | Health check endpoints |
| i18n | `core/i18n` | JSON-based localization |
| logger | `core/logger` | Structured logging |
| mail | `core/mail` | SMTP email via go-mail |
| metrics | `core/metrics` | Prometheus metrics |
| middleware | `core/middleware` | Middleware registry (CORS, CSRF, rate limit, etc.) |
| oauth | `core/oauth` | OAuth2 provider integration |
| plugin | `core/plugin` | Plugin system |
| queue | `core/queue` | Background job queue (Redis-backed) |
| router | `core/router` | Gin-based HTTP router |
| scheduler | `core/scheduler` | Cron-based task scheduling |
| server | `core/server` | HTTP server with graceful shutdown |
| service | `core/service` | Service mode flags (Web, API, WS, Worker) |
| session | `core/session` | Session management (DB, Redis, File, Memory, Cookie) |
| storage | `core/storage` | File storage (local disk + S3) |
| totp | `core/totp` | TOTP two-factor authentication |
| validation | `core/validation` | Request validation |
| websocket | `core/websocket` | WebSocket support via coder/websocket |
| database | `database/` | Connection, transactions, resolver |
| migrations | `database/migrations` | Migration engine |
| models | `database/models` | BaseModel and query scopes |
| seeders | `database/seeders` | Seeder engine |

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
    return seeders.RunAll(db)
})
```

## Documentation

- **[Mastery Process](docs/mastery.md)** — Development workflow and standards
- **[Framework Reference](docs/framework/README.md)** — Complete framework documentation
- **[Project Context](docs/project-context.md)** — Architecture and tech stack

## License

MIT

Copyright (c) 2026 RAi Works (https://rai.works)
