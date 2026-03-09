# RapidGo — Complete Feature List

56 features shipped across 6 phases. Built on Go 1.25+, Gin, GORM, Cobra.

---

## Feature Count by Category

| Category | Count | Phase(s) |
|---|---|---|
| Core Application | 10 | Phase 1 |
| MVC + Authentication | 12 | Phase 2 |
| Web Essentials | 9 | Phase 3 |
| Caching + Events | 4 | Phase 4 |
| Deployment + Testing + DX | 6 | Phase 5 |
| Advanced Features | 13 | Phase 6 |
| Infrastructure | 2 | Additional |
| **Total** | **56** | **All Complete** |

---

## Phase 1 — Core Skeleton (10 features)

### #01 Project Setup & Structure
- **Package**: `cmd/`, root structure
- **Capabilities**: Convention-based directory layout, entrypoint, `.env.example`

### #02 Configuration System
- **Package**: `core/config`
- **Capabilities**: `.env` loading via godotenv, environment detection (development/staging/production), type-safe getters (`GetString`, `GetInt`, `GetBool`)

### #03 Logging
- **Package**: `core/logger`
- **Capabilities**: Structured JSON logging via Go's `log/slog`, configurable log levels (debug/info/warn/error), field-based context

### #04 Error Handling
- **Package**: `core/errors`
- **Capabilities**: Centralized error types, middleware-based error recovery, JSON and HTML error responses

### #05 Service Container (Dependency Injection)
- **Package**: `core/container`
- **Capabilities**: IoC container with singleton and transient bindings, type-safe resolution via generics, `Bind()`, `Singleton()`, `Resolve()`

### #06 Service Providers
- **Package**: `core/container` (provider.go)
- **Capabilities**: `ServiceProvider` interface with `Register()` and `Boot()` lifecycle, modular bootstrapping

### #07 Router & Routing
- **Package**: `core/router`
- **Capabilities**: Gin-based router, route groups, resource routes (7 CRUD actions), named routes with URL generation, route model binding

### #08 Middleware Pipeline
- **Package**: `core/middleware`
- **Capabilities**: Registry with aliases and groups, `Use()` / `Group()` application, request lifecycle hooks

### #09 Database Connection
- **Package**: `database/`
- **Capabilities**: GORM-based connection to PostgreSQL, MySQL, SQLite, `DBConfig` struct, `Connect()`, `ConnectWithConfig()`

### #10 CLI Foundation
- **Package**: `core/cli`
- **Capabilities**: Cobra-based CLI, `serve`, `version`, extensible command registration via hooks

---

## Phase 2 — MVC + Authentication (12 features)

### #11 Models (GORM)
- **Package**: `database/models`
- **Capabilities**: `BaseModel` with `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, GORM relationships (HasOne, HasMany, BelongsTo, Many2Many)

### #12 Database Migrations
- **Package**: `database/migrations`, `core/cli`
- **Capabilities**: Schema management with `Up()` / `Down()`, CLI commands `migrate`, `migrate:rollback`, `migrate:status`, `make:migration`

### #13 Database Seeding
- **Package**: `database/seeders`, `core/cli`
- **Capabilities**: `Seeder` interface with `Name()` and `Seed(db)`, registry with `Register()`, `RunAll()`, `RunByName()`, CLI `db:seed --seeder=name`

### #14 Database Transactions
- **Package**: `database/` (transaction.go)
- **Capabilities**: Auto transactions (`WithTransaction`), manual `Begin`/`Commit`/`Rollback`, nested transactions via savepoints

### #15 Controllers
- **Package**: Application-level (`http/controllers`)
- **Capabilities**: `ResourceController` interface with `Index`, `Show`, `Create`, `Store`, `Edit`, `Update`, `Destroy`

### #16 Response Helpers
- **Package**: Application-level (`http/responses`)
- **Capabilities**: Standardized JSON envelope (`Success`, `Error`, `Paginated`), HTTP status code helpers

### #17 Views & Templates
- **Package**: Application-level (`resources/views`)
- **Capabilities**: Go `html/template` with layouts, partials, template functions, nested templates

### #18 Services Layer
- **Package**: Application-level (`app/services`)
- **Capabilities**: Business logic separation from controllers, testable service pattern

### #19 Helpers
- **Package**: Application-level (`app/helpers`)
- **Capabilities**: Stateless utility functions across all layers

### #20 Session Management
- **Package**: `core/session`
- **Capabilities**: 5 storage backends (database, Redis, file, memory, cookie), flash messages, session middleware

### #21 Authentication (JWT)
- **Package**: `core/auth`
- **Capabilities**: Stateless JWT auth via `golang-jwt/jwt/v5`, token generation, validation, claims extraction, auth middleware

### #22 Crypto & Security Utilities
- **Package**: `core/crypto`
- **Capabilities**: AES-256-GCM encryption/decryption, bcrypt password hashing, HMAC-SHA256 signing, secure random token generation

---

## Phase 3 — Web Essentials (9 features)

### #23 Input Validation
- **Package**: `core/validation`
- **Capabilities**: Built-in validation engine, struct-based validation via `go-playground/validator/v10`, custom rules, error message formatting

### #24 CSRF Protection
- **Package**: `core/middleware` (csrf.go)
- **Capabilities**: Double-submit cookie pattern, per-request token generation, automatic form injection, configurable exempt routes

### #25 CORS Handling
- **Package**: `core/middleware` (cors.go)
- **Capabilities**: Per-origin, per-method, per-header configuration, preflight request handling

### #26 Rate Limiting
- **Package**: `core/middleware` (ratelimit.go)
- **Capabilities**: Token bucket algorithm via `ulule/limiter/v3`, per-IP and per-route limits, configurable window and rate

### #27 Request ID / Tracing
- **Package**: `core/middleware` (request_id.go)
- **Capabilities**: UUID generation per request, header propagation (`X-Request-ID`), correlation for logging

### #28 File Upload & Storage
- **Package**: `core/storage`
- **Capabilities**: Unified storage interface, local filesystem driver, Amazon S3 driver via `aws-sdk-go-v2`, `Put()`, `Get()`, `Delete()`, `Exists()`

### #29 Mail / Email
- **Package**: `core/mail`
- **Capabilities**: SMTP email via `go-mail/v2`, configurable sender, HTML and plain text bodies, attachments

### #30 Static File Serving
- **Package**: `core/router`
- **Capabilities**: CSS, JS, image serving, configurable static directory, cache headers

### #31 WebSocket Support
- **Package**: `core/websocket`
- **Capabilities**: HTTP upgrade via `coder/websocket`, connection hub, broadcast, message handling

---

## Phase 4 — Caching + Events (4 features)

### #32 Caching
- **Package**: `core/cache`
- **Capabilities**: 3 backends (Redis, in-memory, file-based), TTL support, `Get()`, `Set()`, `Delete()`, `Has()`, `Flush()`

### #33 Pagination
- **Package**: `database/`
- **Capabilities**: Page-based pagination helper, configurable page size, metadata (total, per_page, current_page, last_page)

### #34 Events / Hooks System
- **Package**: `core/events`
- **Capabilities**: Publish-subscribe dispatcher, `Listen()`, `Dispatch()`, sync and async event handling, wildcard listeners

### #35 Localization / i18n
- **Package**: `core/i18n`
- **Capabilities**: JSON-based translation files, locale detection, `T()` translation function, parameter interpolation

---

## Phase 5 — Deployment + Testing + DX (6 features)

### #36 Health Checks
- **Package**: `core/health`
- **Capabilities**: Liveness (`/healthz`) and readiness (`/readyz`) probe endpoints, custom check registration, aggregated status

### #37 Graceful Shutdown
- **Package**: `core/server`
- **Capabilities**: OS signal handling (SIGINT, SIGTERM), connection draining, configurable shutdown timeout

### #38 Caddy Integration
- **Package**: `core/server`
- **Capabilities**: Optional Caddy reverse proxy configuration, automatic HTTPS via Caddy

### #39 Docker Deployment
- **Package**: `Dockerfile`, `docker-compose.yml`
- **Capabilities**: Multi-stage Docker build, docker-compose for app + database + Redis, production-ready container

### #40 Testing Infrastructure
- **Package**: `testing/testutil`
- **Capabilities**: Test utilities and helpers, isolated test setup, test database support

### #41 Code Generation (CLI Scaffolding)
- **Package**: `core/cli`
- **Capabilities**: `make:controller`, `make:model`, `make:service`, `make:provider`, `make:migration`, `make:admin` — generates Go source files from templates

---

## Phase 6 — Advanced (13 features)

### #42 Queue Workers / Background Jobs
- **Package**: `core/queue`
- **Capabilities**: 4 drivers (database, Redis, memory, sync), `Enqueue()`, `Dequeue()`, retry with backoff, `work` CLI command with `--queues`, `--workers`, `--timeout`

### #43 Task Scheduler / Cron
- **Package**: `core/scheduler`
- **Capabilities**: Cron-based scheduling via `robfig/cron/v3`, `EveryMinute()`, `Hourly()`, `Daily()`, `Weekly()`, custom cron expressions, `schedule:run` CLI

### #44 Plugin / Module System
- **Package**: `core/plugin`
- **Capabilities**: `Plugin` interface with `Name()`, `Boot()`, route hooks, command hooks, event hooks, plugin isolation

### #45 GraphQL Support
- **Package**: `core/graphql`
- **Capabilities**: GraphQL handler via `graphql-go/graphql`, GraphiQL playground, custom schema and resolver registration

### #46 Admin Panel Scaffolding
- **Package**: `core/cli`, application-level controllers
- **Capabilities**: `make:admin` CLI command, generates admin CRUD controllers and views, admin middleware

### #47 API Versioning
- **Package**: `core/router`
- **Capabilities**: Version-prefixed route groups (`/api/v1/`, `/api/v2/`), version negotiation

### #48 WebSocket Rooms / Channels
- **Package**: `core/websocket`
- **Capabilities**: Room-based message routing, join/leave channels, room broadcast, per-room connection tracking

### #49 OAuth2 / Social Login
- **Package**: `core/oauth`
- **Capabilities**: Google, GitHub, Facebook providers via `x/oauth2`, custom provider support, token exchange, user profile retrieval

### #50 Two-Factor Authentication (TOTP)
- **Package**: `core/totp`
- **Capabilities**: TOTP via `pquerna/otp`, QR code generation, bcrypt-hashed backup codes, verify and validate

### #51 Audit Logging
- **Package**: `core/audit`
- **Capabilities**: `AuditLog` model with actor, action, entity, old/new values, IP address, structured audit trail

### #52 Soft Deletes
- **Package**: `database/models`
- **Capabilities**: GORM `DeletedAt` field, `WithTrashed` scope (include deleted), `OnlyTrashed` scope (only deleted), automatic query filtering

### #53 Database Read/Write Splitting
- **Package**: `database/` (resolver.go)
- **Capabilities**: Separate read replica connection, GORM `DBResolver` integration, automatic read/write routing

### #54 Prometheus Metrics
- **Package**: `core/metrics`
- **Capabilities**: Request duration histogram, status code counters, custom metric registration, `/metrics` endpoint via `prometheus/client_golang`

---

## Additional Features

### #55 Framework Rename (RGo → RapidGo)
- **Scope**: All packages
- **Capabilities**: Module path migration, import path updates, CLI binary rename

### #56 Service Mode (Multi-Port Serving)
- **Package**: `core/service`
- **Capabilities**: `Mode` type (Web, API, WS, All), `--mode` flag on `serve` command, separate port binding per service type

---

## Phase Summary

| Phase | Name | Features | Status |
|---|---|---|---|
| 1 | Core Skeleton | #01 – #10 (10) | ✅ Complete |
| 2 | MVC + Auth | #11 – #22 (12) | ✅ Complete |
| 3 | Web Essentials | #23 – #31 (9) | ✅ Complete |
| 4 | Caching + Events | #32 – #35 (4) | ✅ Complete |
| 5 | Deploy + Testing + DX | #36 – #41 (6) | ✅ Complete |
| 6 | Advanced | #42 – #54 (13) | ✅ Complete |
| — | Additional | #55 – #56 (2) | ✅ Complete |
| | **Total** | **56** | **All Shipped** |
