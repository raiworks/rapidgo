---
title: "Project Structure"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Project Structure

## Abstract

This document describes the framework's directory layout, file
organization conventions, and the purpose of each top-level directory.
It serves as a map for developers navigating the codebase and
understanding where to place new code.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Full Directory Tree](#2-full-directory-tree)
3. [Directory Reference](#3-directory-reference)
4. [Framework vs User Code](#4-framework-vs-user-code)
5. [Root Files](#5-root-files)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. Full Directory Tree

```text
yourframework/
│
├── cmd/                        # Application entrypoint
│   └── main.go                 # package main, func main()
│
├── core/
│   ├── app/                    # Application container & lifecycle
│   ├── container/              # Service container (DI)
│   ├── router/                 # Router engine setup
│   ├── middleware/             # Middleware definitions & registry
│   ├── config/                 # Config loader (reads .env / YAML)
│   ├── logger/                 # Structured logging
│   ├── errors/                 # Error types & error middleware
│   ├── session/                # Session manager & store backends
│   ├── validation/             # Built-in form & data validation
│   ├── crypto/                 # Hashing, encryption, tokens
│   ├── cache/                  # Cache manager (Redis, memory)
│   ├── mail/                   # Email sender
│   ├── events/                 # Event dispatcher & listeners
│   ├── i18n/                   # Localization / translations
│   ├── server/                 # HTTP server & Caddy integration
│   └── websocket/              # WebSocket upgrader & hub
│
├── database/
│   ├── connection.go           # DB connection factory
│   ├── migrations/             # Migration files
│   ├── seeders/                # Database seed data
│   ├── models/                 # GORM model structs
│   └── querybuilder/           # Query builder helpers
│
├── app/
│   ├── providers/              # Service providers (register & boot)
│   ├── services/               # Business logic (custom services)
│   └── helpers/                # Utility / helper functions
│
├── http/
│   ├── controllers/            # Request handlers
│   ├── requests/               # Request validation structs
│   └── responses/              # Response helpers & types
│
├── routes/
│   ├── web.go                  # Web (HTML) routes
│   └── api.go                  # API routes
│
├── resources/
│   ├── views/                  # HTML templates
│   ├── lang/                   # Translation files (en.json, etc.)
│   └── static/                 # CSS, JS, images
│
├── storage/
│   ├── uploads/                # User-uploaded files
│   ├── cache/                  # File-based cache
│   ├── sessions/               # File-based sessions
│   └── logs/                   # Log files
│
├── tests/                      # Test files
│   ├── unit/
│   └── integration/
│
├── .env                        # Environment variables
├── Caddyfile                   # (optional) Caddy reverse proxy config
├── Dockerfile                  # (optional) Container build
├── docker-compose.yml          # (optional) Multi-service orchestration
├── go.mod
└── go.sum
```

## 3. Directory Reference

### `cmd/` — Application Entrypoint

Contains `main.go` — the single entry point for the application.
This file **MUST**:

1. Load configuration (`config.Load()`)
2. Create the application container (`app.New()`)
3. Register all service providers
4. Boot all providers
5. Set up the router with middleware
6. Start the HTTP server with graceful shutdown

Developers **SHOULD NOT** place business logic in `cmd/`. It is purely
an orchestration layer.

### `core/` — Framework Internals

The heart of the framework. Contains 16 subdirectories, each with a
focused responsibility:

| Package | Purpose |
|---------|---------|
| `core/app/` | `App` struct — orchestrates provider lifecycle (Register → Boot) |
| `core/container/` | Service container — `Bind()`, `Singleton()`, `Instance()`, `Make()`, `MustMake[T]()`, `Has()` |
| `core/router/` | Gin router setup, resource routes, named routes, route model binding |
| `core/middleware/` | Middleware registry (aliases, groups), built-in middleware (auth, CSRF, CORS, rate limit, request ID, session) |
| `core/config/` | `.env` loader (godotenv), `Env()` helper, environment detection (`IsProduction()`, `IsDevelopment()`, `IsTesting()`, `IsDebug()`) |
| `core/logger/` | `log/slog` setup with JSON handler |
| `core/errors/` | Error middleware — catches panics, returns JSON/HTML errors |
| `core/session/` | Session manager, store interface, 5 backends (DB, Redis, File, Memory, Cookie), flash messages |
| `core/validation/` | Built-in validator — Required, MinLength, MaxLength, Email, URL, IP, Matches, In, Confirmed |
| `core/crypto/` | RandomBytes, RandomHex, SHA256Hash, HMACSign/Verify, AES-256-GCM Encrypt/Decrypt |
| `core/cache/` | Cache store interface, Redis and memory backends |
| `core/mail/` | SMTP mailer using `go-mail` |
| `core/events/` | Event dispatcher — `Listen()`, `Dispatch()`, `DispatchSync()` |
| `core/i18n/` | Translator — loads JSON translation files, parameter substitution, fallback locale |
| `core/server/` | HTTP server configuration, embedded Caddy integration |
| `core/websocket/` | WebSocket handler using `coder/websocket` |

Application code **MUST NOT** modify `core/` packages directly. To
customize behavior, use service providers to swap implementations.

### `database/` — Data Layer

| File/Directory | Purpose |
|----------------|---------|
| `connection.go` | Multi-driver database connection factory (PostgreSQL, MySQL, SQLite) with connection pooling |
| `migrations/` | Migration files — GORM `AutoMigrate` based |
| `seeders/` | Database seed data — `RunAll()` orchestrates all seeders |
| `models/` | GORM model structs — `BaseModel`, `User`, `Post`, etc. |
| `querybuilder/` | Query builder helpers |

### `app/` — User-Managed Code

This is where developers place their application-specific code. The
framework generates scaffolding into these directories:

| Directory | Purpose |
|-----------|---------|
| `app/providers/` | Custom service providers (e.g., `PaymentProvider`) |
| `app/services/` | Business logic services (e.g., `UserService`, `OrderService`) |
| `app/helpers/` | Utility functions (password hashing, slugify, string helpers, data helpers) |

The `app/` directory represents **user code** while `core/` represents
**framework code**. This separation ensures framework upgrades do not
conflict with user modifications.

### `http/` — HTTP Layer

| Directory | Purpose |
|-----------|---------|
| `http/controllers/` | Request handlers — both function-based and struct-based (for `ResourceController`) |
| `http/requests/` | Request validation structs with `binding` tags (for struct-based validation) |
| `http/responses/` | `APIResponse` struct, `Success()`, `Created()`, `Error()`, `Paginated()` helpers |

### `routes/` — Route Definitions

| File | Purpose |
|------|---------|
| `web.go` | Web (HTML/SSR) routes — uses `web` middleware group (session, CSRF, request ID) |
| `api.go` | API (JSON) routes — uses `api` middleware group (CORS, rate limit, request ID) |

### `resources/` — Frontend Assets

| Directory | Purpose |
|-----------|---------|
| `resources/views/` | HTML templates rendered by `html/template` or Templ |
| `resources/lang/` | Translation JSON files (`en.json`, `es.json`, etc.) |
| `resources/static/` | CSS, JavaScript, images — served via `r.Static("/static", ...)` |

### `storage/` — Runtime Data

| Directory | Purpose |
|-----------|---------|
| `storage/uploads/` | User-uploaded files (when using local storage driver) |
| `storage/cache/` | File-based cache data |
| `storage/sessions/` | File-based session data (when `SESSION_DRIVER=file`) |
| `storage/logs/` | Application log files |

The `storage/` directory **MUST** be writable by the application
process. It **SHOULD** be excluded from version control (except the
directory structure itself).

### `tests/` — Test Files

| Directory | Purpose |
|-----------|---------|
| `tests/unit/` | Unit tests for services, helpers, and isolated functions |
| `tests/integration/` | HTTP handler integration tests using `httptest` |

Run tests with:

```bash
go test ./tests/... -v
go test ./... -cover
```

## 4. Framework vs User Code

The separation between `core/` and `app/` is a fundamental design
decision:

```text
┌─────────────────────────────────────┐
│           core/ (Framework)          │
│  Container, Router, Middleware,      │
│  Session, Validation, Crypto,        │
│  Cache, Mail, Events, i18n,         │
│  Logger, Errors, Server, WebSocket  │
├─────────────────────────────────────┤
│           app/ (User Code)           │
│  Providers, Services, Helpers        │
├─────────────────────────────────────┤
│           http/ (HTTP Bridge)        │
│  Controllers, Requests, Responses    │
└─────────────────────────────────────┘
```

- **`core/`** — Framework code. Developers extend it through interfaces
  and service providers, never by modifying it directly.
- **`app/`** — User code. All custom providers, services, and helpers
  live here. This is where the `make:*` CLI commands generate files.
- **`http/`** — The bridge between HTTP and business logic. Controllers
  parse requests, delegate to services, and format responses.

## 5. Root Files

| File | Required | Purpose |
|------|----------|---------|
| `.env` | **REQUIRED** | Environment variables (DB credentials, app config, secrets) |
| `go.mod` | **REQUIRED** | Go module definition |
| `go.sum` | **REQUIRED** | Go dependency checksums |
| `Caddyfile` | OPTIONAL | Caddy reverse proxy configuration |
| `Dockerfile` | OPTIONAL | Multi-stage container build |
| `docker-compose.yml` | OPTIONAL | Multi-service orchestration |

## 6. Security Considerations

- The `.env` file **MUST NOT** be committed to version control. It
  contains secrets (database passwords, JWT secret, APP_KEY).
- The `storage/` directory **SHOULD** have restricted filesystem
  permissions (0700 or 0750).
- User uploads in `storage/uploads/` **MUST** be validated before
  serving. Do not serve executable files.
- The `core/crypto/` package `APP_KEY` **MUST** be exactly 32 bytes
  for AES-256-GCM encryption.

## 7. References

- [Architecture Overview](overview.md)
- [Design Principles](design-principles.md)
- [Application Lifecycle](application-lifecycle.md)
- [CLI Overview](../cli/cli-overview.md)
- [Code Generation](../cli/code-generation.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
