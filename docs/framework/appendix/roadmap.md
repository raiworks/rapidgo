---
title: "Roadmap"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Roadmap

## Abstract

This document lists planned advanced features, recommended libraries,
and the development timeline for the framework.

## Table of Contents

1. [Advanced Features (Planned)](#1-advanced-features-planned)
2. [Development Timeline](#2-development-timeline)
3. [Scope Estimate](#3-scope-estimate)
4. [References](#4-references)

## 1. Advanced Features (Planned)

| Feature | Recommended Library | Priority | Scope |
|---------|-------------------|----------|-------|
| Queue workers / background jobs | `github.com/hibiken/asynq` or `github.com/riverqueue/river` | High | Async task processing |
| Task scheduler / cron | `github.com/robfig/cron/v3` | High | Periodic job execution |
| Plugin / module system | Dynamic loading via service providers | Medium | Extensibility |
| GraphQL support | `github.com/99designs/gqlgen` | Medium | Alternative API layer |
| Admin panel / dashboard scaffolding | — | Medium | Auto-generated CRUD UI |
| API versioning | URL prefix or header-based | Medium | Multiple API versions |
| WebSocket rooms / channels | Broadcast groups | Medium | Real-time features |
| OAuth2 / social login | Google, GitHub, etc. | High | Third-party authentication |
| Two-factor authentication (TOTP) | — | High | Enhanced security |
| Audit logging | — | Medium | Who changed what |
| Soft deletes | GORM `DeletedAt` | Low | Recoverable deletions |
| Database read/write splitting | — | Low | Performance scaling |
| Prometheus metrics endpoint | — | Medium | Observability |

## 2. Development Timeline

| Phase | Scope | Estimate |
|-------|-------|----------|
| Core skeleton | Router, config, DB, middleware, CLI, logging, errors, service container | 2–3 weeks |
| MVC + Auth | Controllers, services, helpers, models, JWT, sessions, flash messages | 3–4 weeks |
| Web essentials | CSRF, validation, CORS, rate limiting, file upload, mail, middleware registry | 3–4 weeks |
| Caching + Events | Cache layer, events/hooks, i18n, pagination, env detection | 2–3 weeks |
| Caddy + Deploy | Caddy integration, Docker, health checks, testing, code generation | 2–3 weeks |
| Advanced | Queues, scheduler, plugins, GraphQL, OAuth2 | 6–12 months |
| Full ecosystem | Feature parity with Laravel | 1–2+ years |

## 3. Scope Estimate

A focused, well-designed Go framework can be effective at **20–35k
lines of code** for the core + auth + sessions + CLI + ORM + caching +
mail + file storage layers.

## 4. References

- [Architecture Overview](../architecture/overview.md)
- [Design Principles](../architecture/design-principles.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
