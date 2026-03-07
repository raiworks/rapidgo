# RapidGo Framework

A modern Go web framework inspired by Laravel and CodeIgniter — combining
**Laravel-style developer experience** with **Go performance**.

## Overview

RapidGo is an opinionated, full-stack Go web framework built on proven
libraries (Gin, GORM, Cobra) that provides everything needed for web
applications, REST APIs, and real-time WebSocket services.

### Key Features

- **MVC + Services + Helpers** architecture
- **Service Container & Providers** — Laravel-style IoC for extensibility
- **Multi-database support** — PostgreSQL, MySQL, SQLite via GORM
- **Session management** — DB, Redis, File, Memory, Cookie backends
- **Built-in validation** — zero-dependency validator + struct-based (go-playground)
- **Built-in crypto** — AES-256-GCM, HMAC-SHA256, secure random tokens
- **JWT & session-based authentication**
- **CLI scaffolding** — `make:controller`, `make:model`, `make:service`, etc.
- **Middleware registry** — aliases, groups, custom middleware
- **WebSocket support** — via `coder/websocket`
- **Caching** — Redis and in-memory backends
- **Mail** — SMTP via `go-mail`
- **File storage** — local disk and S3
- **Events / hooks** — pub-sub event dispatcher
- **i18n / localization** — JSON translation files
- **Caddy integration** — embedded or reverse proxy (optional)
- **Docker support** — multi-stage builds (optional)
- **Graceful shutdown**, health checks, CSRF, CORS, rate limiting

## Project Status

**Phase:** All 56 features shipped (Phases 1–6 complete)

The framework is feature-complete with full test coverage across 39+ packages. All features were built following the [Mastery](docs/mastery.md) development process.

## Documentation

All documentation lives in the [`docs/`](docs/) directory:

- **[Project Context](docs/project-context.md)** — What RapidGo is, technology stack, architecture
- **[Project Roadmap](docs/project-roadmap.md)** — Feature plan, dependencies, progress
- **[Mastery Process](docs/mastery.md)** — Development workflow and standards
- **[Framework Reference](docs/framework/README.md)** — Complete framework documentation

## Quick Start

```bash
# Setup
cp .env.example .env   # then edit .env with your credentials

# Build
make build

# Run
make run

# Test
make test
```

## Repository

- **GitHub:** [RAiWorks/RapidGo](https://github.com/RAiWorks/RapidGo)

## Tech Stack

| Component | Library |
|-----------|---------|
| HTTP Router | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) |
| CLI | [Cobra](https://github.com/spf13/cobra) |
| Config | [godotenv](https://github.com/joho/godotenv) |
| JWT | [golang-jwt](https://github.com/golang-jwt/jwt) |
| WebSocket | [coder/websocket](https://github.com/coder/websocket) |
| Redis | [go-redis](https://github.com/redis/go-redis) |
| Mail | [go-mail](https://github.com/wneessen/go-mail) |
| S3 | [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) |
| Scheduler | [robfig/cron](https://github.com/robfig/cron) |
| TOTP | [pquerna/otp](https://github.com/pquerna/otp) |
| GraphQL | [graphql-go](https://github.com/graphql-go/graphql) |
| OAuth2 | [x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) |
| Metrics | [prometheus](https://github.com/prometheus/client_golang) |

## License

MIT License. See [LICENSE](LICENSE) for details.

Copyright (c) 2026 RAi Works (https://rai.works)
