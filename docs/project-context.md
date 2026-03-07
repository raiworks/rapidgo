# 🎯 Project Context — RapidGo Framework

> **Project**: RapidGo
> **Type**: Go Web Framework
> **Repository**: https://github.com/RAiWorks/RapidGo
> **Author**: RAiWorks
> **Created**: 2026-03-05

---

## What Is RapidGo?

RapidGo is a full-featured, opinionated Go web framework that combines **Laravel-style developer experience with Go performance**. It supports web applications, REST APIs, WebSockets, and CLI tools — all built on proven Go libraries and idiomatic patterns.

The framework follows the **MVC + Services + Helpers** architectural pattern, providing everything needed to build production-ready applications without cobbling together dozens of unrelated packages.

---

## Mission Statement

Provide a Go web framework that:

1. **Feels familiar** to developers coming from Laravel, CodeIgniter, or similar frameworks
2. **Stays idiomatic** in Go — no magic, no reflection abuse, clear error handling
3. **Bundles essentials** so developers don't assemble a framework from scratch for every project
4. **Performs at Go speed** without compromising developer experience
5. **Scales from prototypes to production** with the same codebase

---

## Target Capabilities

### Core Application

| Capability | Description |
|---|---|
| Web applications | MVC + Services + Helpers pattern |
| REST APIs | JSON responses with standardized envelope |
| WebSockets | Real-time communication via Coder WebSocket |
| CLI tools | Scaffolding, migrations, seeds via Cobra |

### Data & Storage

| Capability | Description |
|---|---|
| ORM | GORM with PostgreSQL, MySQL, SQLite support |
| Migrations | Schema management with up/down migrations |
| Seeders | Database seeding for dev/test environments |
| Transactions | GORM transaction patterns (auto, manual, nested) |
| Pagination | Configurable page-based pagination helper |
| File storage | Local filesystem and Amazon S3 via unified interface |
| Caching | Redis, in-memory, and file-based cache with TTL support |

### HTTP & Middleware

| Capability | Description |
|---|---|
| Routing | Routes, groups, resource routes, named routes, route model binding |
| Controllers | MVC controllers with ResourceController interface (7 CRUD actions) |
| Middleware | Pipeline with registry, aliases, and groups |
| Validation | Built-in engine + struct-based validation via `go-playground/validator` |
| Responses | API response helpers with success/error/paginated envelope |
| Views | Go `html/template` with layouts, partials, and template functions |
| Static files | CSS, JS, images served from `resources/static/` |
| WebSocket | Upgrade, hub, broadcast, and per-connection messaging |

### Security

| Capability | Description |
|---|---|
| Authentication | JWT (stateless) + session-based (stateful) |
| Sessions | 5 backends — database, Redis, file, memory, cookie |
| CSRF | Double-submit cookie pattern with per-request tokens |
| CORS | Per-origin, per-method, per-header configuration |
| Rate limiting | Token bucket with per-IP and per-route limits |
| Crypto | AES-256-GCM encryption, bcrypt hashing, HMAC, secure random tokens |
| Request ID | Unique identifier per request for tracing |

### Infrastructure

| Capability | Description |
|---|---|
| Service container | Dependency injection with singleton and transient bindings |
| Service providers | Register/Boot lifecycle, built-in + custom providers |
| Configuration | `.env` loading via godotenv |
| Environment detection | Development, production, testing modes |
| Logging | Structured JSON logging via `log/slog` |
| Error handling | Centralized error middleware with JSON/HTML responses |
| Mail | SMTP email sending via `go-mail` |
| Events | Publish-subscribe system with sync/async dispatch |
| i18n | JSON-based translation files with locale detection |
| Health checks | Liveness and readiness probe endpoints |
| Graceful shutdown | Signal handling with connection draining |

### Testing

| Capability | Description |
|---|---|
| Unit tests | Service, helper, and model tests using `testing` package |
| Integration tests | HTTP handler tests with test server and database |

> **Test location**: All tests are **co-located** next to the code they test (e.g. `core/auth/auth_test.go`). The top-level `tests/` directory is reserved for future end-to-end or cross-package integration tests; `testing/testutil/` holds shared test helpers.

### Deployment

| Capability | Description |
|---|---|
| Caddy | Optional — embedded library or external reverse proxy with auto-HTTPS |
| Docker | Multi-stage Dockerfile + docker-compose orchestration |
| Build & run | Single binary entrypoint with graceful shutdown |

### Advanced Features (Shipped)

| Feature | Status |
|---|---|
| Queue workers / background jobs | ✅ Shipped |
| Task scheduler / cron | ✅ Shipped |
| Plugin / module system | ✅ Shipped |
| GraphQL support | ✅ Shipped |
| Admin panel scaffolding | ✅ Shipped |
| API versioning | ✅ Shipped |
| WebSocket rooms / channels | ✅ Shipped |
| OAuth2 / social login | ✅ Shipped |
| Two-factor authentication (TOTP) | ✅ Shipped |
| Audit logging | ✅ Shipped |
| Soft deletes | ✅ Shipped |
| Database read/write splitting | ✅ Shipped |
| Prometheus metrics | ✅ Shipped |

---

## Technology Stack

| Component | Library | Import Path |
|---|---|---|
| Language | Go | 1.25+ |
| HTTP Router | Gin | `github.com/gin-gonic/gin` |
| ORM | GORM | `gorm.io/gorm` |
| CLI | Cobra | `github.com/spf13/cobra` |
| Configuration | godotenv | `github.com/joho/godotenv` |
| JWT | golang-jwt | `github.com/golang-jwt/jwt/v5` |
| WebSocket | coder/websocket | `github.com/coder/websocket` |
| Redis | go-redis | `github.com/redis/go-redis/v9` |
| CORS | gin-contrib/cors | `github.com/gin-contrib/cors` |
| Rate Limiting | ulule/limiter | `github.com/ulule/limiter/v3` |
| Email | go-mail | `github.com/wneessen/go-mail` |
| S3 Storage | aws-sdk-go-v2 | `github.com/aws/aws-sdk-go-v2` |
| Web Server | Caddy | `github.com/caddyserver/caddy/v2` |
| Logging | slog | `log/slog` (standard library) |
| Validation | validator | `github.com/go-playground/validator/v10` |
| Password Hashing | bcrypt | `golang.org/x/crypto/bcrypt` |
| Task Scheduler | robfig/cron | `github.com/robfig/cron/v3` |
| TOTP | pquerna/otp | `github.com/pquerna/otp` |
| GraphQL | graphql-go | `github.com/graphql-go/graphql` |
| OAuth2 | x/oauth2 | `golang.org/x/oauth2` |
| Metrics | prometheus | `github.com/prometheus/client_golang` |
| UUID | google/uuid | `github.com/google/uuid` |

---

## Architectural Pattern

```
┌──────────────────────────────────────────────────────┐
│                    HTTP Request                       │
├──────────────────────────────────────────────────────┤
│               Middleware Pipeline                     │
│   (auth, CSRF, CORS, rate-limit, request-id, session)│
├──────────────────────────────────────────────────────┤
│                   Router (Gin)                        │
│          (resource routes, named routes,              │
│           route model binding)                        │
├──────────────────────────────────────────────────────┤
│                  Controllers                          │
│          (HTTP concerns only — parse request,         │
│           call service, return response)              │
├──────────────────────────────────────────────────────┤
│                   Services                            │
│          (business logic, domain rules,               │
│           orchestration)                              │
├──────────────────────────────────────────────────────┤
│               Models (GORM)                           │
│          (data schema, relationships,                 │
│           hooks, database queries)                    │
├──────────────────────────────────────────────────────┤
│                  Database                             │
│          (PostgreSQL, MySQL, SQLite)                  │
└──────────────────────────────────────────────────────┘
```

**Layer rules**:

- **Controllers** handle HTTP concerns: parse request data, call services, return responses
- **Services** contain business logic — they **MUST NOT** access HTTP request/response objects
- **Models** define GORM data schemas, relationships, and hooks
- **Helpers** provide stateless utility functions used across all layers

---

## Project Structure

```
yourframework/
├── cmd/                        # Application entrypoint
│   └── main.go
├── core/                       # Framework internals
│   ├── app/                    # Application container & lifecycle
│   ├── audit/                  # Audit logging
│   ├── auth/                   # JWT authentication
│   ├── cache/                  # Cache manager
│   ├── cli/                    # CLI commands (Cobra)
│   ├── config/                 # Config loader (.env)
│   ├── container/              # Service container (DI)
│   ├── crypto/                 # Hashing, encryption, tokens
│   ├── errors/                 # Error types & middleware
│   ├── events/                 # Event dispatcher & listeners
│   ├── graphql/                # GraphQL handler & playground
│   ├── health/                 # Health check endpoint
│   ├── i18n/                   # Localization / translations
│   ├── logger/                 # Structured logging
│   ├── mail/                   # Email sender
│   ├── metrics/                # Prometheus metrics
│   ├── middleware/             # Middleware definitions & registry
│   ├── oauth/                  # OAuth2 / social login providers
│   ├── plugin/                 # Plugin / module system
│   ├── queue/                  # Queue workers / background jobs
│   ├── router/                 # Router engine, groups, resources
│   ├── scheduler/              # Task scheduler / cron
│   ├── server/                 # HTTP server & Caddy integration
│   ├── service/                # Service mode (multi-port serving)
│   ├── session/                # Session manager & store backends
│   ├── storage/                # File storage (local, S3)
│   ├── totp/                   # Two-factor authentication (TOTP)
│   ├── validation/             # Validation engine
│   └── websocket/              # WebSocket upgrader, hub & rooms
├── database/
│   ├── connection.go           # DB connection factory
│   ├── migrations/             # Migration files
│   ├── seeders/                # Seed data
│   ├── models/                 # GORM model structs
│   └── querybuilder/           # Query builder helpers
├── app/
│   ├── helpers/                # Utility functions
│   ├── jobs/                   # Queue job definitions
│   ├── plugins.go              # Plugin registration
│   ├── providers/              # Service providers
│   ├── schedule/               # Scheduled task definitions
│   └── services/               # Business logic
├── http/
│   ├── controllers/            # Request handlers
│   │   └── admin/              # Admin panel controllers (scaffolded)
│   ├── requests/               # Validation structs
│   └── responses/              # Response helpers
├── routes/
│   ├── web.go                  # Web (HTML) routes
│   ├── api.go                  # API routes
│   └── ws.go                   # WebSocket routes
├── plugins/
│   └── example/                # Example plugin demo
├── testing/
│   └── testutil/               # Test helper utilities
├── resources/
│   ├── views/                  # HTML templates
│   ├── lang/                   # Translation files
│   └── static/                 # CSS, JS, images
├── storage/
│   ├── uploads/                # User-uploaded files
│   ├── cache/                  # File-based cache
│   ├── sessions/               # File-based sessions
│   └── logs/                   # Log files
├── tests/                      # Reserved for cross-package / E2E tests
│   ├── unit/
│   └── integration/
├── docs/                       # This documentation
├── .env.example                # Environment config template (copy to .env)
├── go.mod
└── go.sum
```

---

## Scope Boundaries

### In Scope (Current Framework)

Everything listed in the **Target Capabilities** section above. This is the scope defined by the framework blueprint and documented in `docs/framework/`.

### Out of Scope (Not Now)

All planned features (Phases 1–6) have been shipped. No features are currently out of scope.

### Non-Goals

| Item | Reason |
|---|---|
| Frontend framework | RapidGo serves HTML templates; it does not include a JS framework |
| Microservices framework | RapidGo is a monolith-first framework; gRPC/service mesh is out of scope |
| ORM replacement | RapidGo uses GORM — building a custom ORM is not a goal |
| Package manager | Go modules (`go.mod`) handles dependency management |

---

## Design Principles

1. **Convention over configuration** — sensible defaults, override when needed
2. **Explicit over magic** — no hidden behavior, clear function signatures
3. **Composition over inheritance** — Go interfaces and embedding, not class hierarchies
4. **Fail fast, fail loud** — validate early, return errors immediately
5. **Single responsibility** — each package does one thing well
6. **Zero-config start** — `go run cmd/main.go` works out of the box with sensible defaults

---

## Reference Documents

| Document | Location |
|---|---|
| Blueprint (source of truth) | `reference/docs/go_web_framework_blueprint.md` |
| Framework reference docs | `docs/framework/README.md` |
| Development process | `docs/mastery.md` |
| Feature roadmap | `docs/project-roadmap.md` |

---

> *"Laravel-style developer experience. Go performance. No compromises."*
