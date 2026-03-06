# 🗺️ Project Roadmap — RGo Framework

> **Project**: RGo
> **Context**: [`project-context.md`](project-context.md)
> **Process**: [`mastery.md`](mastery.md)
> **Last Updated**: 2026-03-06

---

## How to Read This Document

Each feature is numbered in **dependency order** — lower numbers must be built first. Features are grouped into **phases** that align with the blueprint's development timeline.

Every feature follows the [Mastery lifecycle](mastery.md): Discuss → Design → Plan → Build → Ship → Reflect.

### Status Legend

| Symbol | Meaning |
|---|---|
| ⬜ | Not started |
| 🟡 | In progress |
| ✅ | Complete |
| 🔮 | Future (not in current scope) |

---

## Phase 1 — Core Skeleton

> Foundation layer. Everything else depends on this.

| # | Feature | Status | Depends On | Blueprint Sections |
|---|---|---|---|---|
| 01 | Project Setup & Structure | ✅ | — | Suggested Project Structure |
| 02 | Configuration System | ✅ | #01 | Configuration System |
| 03 | Logging | ✅ | #01, #02 | Logging |
| 04 | Error Handling | ✅ | #01, #03 | Error Handling |
| 05 | Service Container | ✅ | #01 | Service Container & Providers |
| 06 | Service Providers | ✅ | #05 | Service Container & Providers |
| 07 | Router & Routing | ✅ | #01, #05 | Router Layer (routes, groups, resource routes, named routes, route model binding) |
| 08 | Middleware Pipeline | ✅ | #07 | Middleware (Custom Registration) |
| 09 | Database Connection | ✅ | #02, #05 | Database Layer |
| 10 | CLI Foundation | ✅ | #01, #02 | CLI Tools (cobra setup, base commands) |

### Phase 1 — Blueprint Coverage

| Blueprint Section | Feature(s) |
|---|---|
| Suggested Project Structure | #01 |
| Configuration System | #02 |
| Logging | #03 |
| Error Handling | #04 |
| Service Container & Providers | #05, #06 |
| Router Layer | #07 |
| Middleware (Custom Registration) | #08 |
| Database Layer | #09 |
| CLI Tools (foundation) | #10 |

---

## Phase 2 — MVC + Auth

> Application layer. Controllers, services, models, authentication, sessions.

| # | Feature | Status | Depends On | Blueprint Sections |
|---|---|---|---|---|
| 11 | Models (GORM) | ✅ | #09 | Models (GORM) |
| 12 | Database Migrations | ✅ | #09, #10 | CLI Tools (migrate commands) |
| 13 | Database Seeding | ✅ | #09, #10, #11 | Database Seeding |
| 14 | Database Transactions | ✅ | #09, #11 | Database Transactions |
| 15 | Controllers | ✅ | #07, #08 | MVC Controller Example |
| 16 | Response Helpers | ✅ | #07 | Response Helpers |
| 17 | Views & Templates | ✅ | #07, #15 | View / Template Engine |
| 18 | Services Layer | ✅ | #05, #11 | Services Layer |
| 19 | Helpers | ✅ | #01 | Helpers, Built-in String & Data Helpers |
| 20 | Session Management | ✅ | #02, #08, #09 | Session Management (DB, Redis, file, memory, cookie — all 5 backends) |
| 21 | Authentication | ✅ | #20, #11, #19, #22 | Authentication (JWT + session-based) |
| 22 | Crypto & Security Utilities | ✅ | #01 | Built-in Crypto & Security Utilities |

### Phase 2 — Blueprint Coverage

| Blueprint Section | Feature(s) |
|---|---|
| Models (GORM) | #11 |
| CLI Tools (migrate, seed commands) | #12, #13 |
| Database Seeding | #13 |
| Database Transactions | #14 |
| MVC Controller Example | #15 |
| Response Helpers | #16 |
| View / Template Engine | #17 |
| Services Layer | #18 |
| Helpers | #19 |
| Built-in String & Data Helpers | #19 |
| Session Management | #20 |
| Authentication | #21 |
| Built-in Crypto & Security Utilities | #22 |

---

## Phase 3 — Web Essentials

> Security middleware, validation, file upload, mail.

| # | Feature | Status | Depends On | Blueprint Sections |
|---|---|---|---|---|
| 23 | Input Validation | ✅ | #07, #15 | Input Validation (Built-in) |
| 24 | CSRF Protection | ✅ | #08, #20 | CSRF Protection |
| 25 | CORS Handling | ✅ | #08 | CORS |
| 26 | Rate Limiting | ✅ | #08 | Rate Limiting |
| 27 | Request ID / Tracing | ✅ | #08 | Request ID / Tracing |
| 28 | File Upload & Storage | ✅ | #02, #05 | File Upload & Storage (local + S3) |
| 29 | Mail / Email | ✅ | #02, #05 | Mail / Email |
| 30 | Static File Serving | ✅ | #07 | Static File Serving |
| 31 | WebSocket Support | ✅ | #07, #08 | WebSocket Support |

### Phase 3 — Blueprint Coverage

| Blueprint Section | Feature(s) |
|---|---|
| Input Validation (Built-in) | #23 |
| CSRF Protection | #24 |
| CORS | #25 |
| Rate Limiting | #26 |
| Request ID / Tracing | #27 |
| File Upload & Storage | #28 |
| Mail / Email | #29 |
| Static File Serving | #30 |
| WebSocket Support | #31 |

---

## Phase 4 — Caching + Events

> Infrastructure services, pagination, i18n.

| # | Feature | Status | Depends On | Blueprint Sections |
|---|---|---|---|---|
| 32 | Caching | ⬜ | #02, #05 | Caching (Redis + memory + file) |
| 33 | Pagination | ⬜ | #09, #11 | Pagination |
| 34 | Events / Hooks System | ⬜ | #05 | Events / Hooks System |
| 35 | Localization / i18n | ⬜ | #02 | Localization / i18n |

### Phase 4 — Blueprint Coverage

| Blueprint Section | Feature(s) |
|---|---|
| Caching | #32 |
| Pagination | #33 |
| Events / Hooks System | #34 |
| Localization / i18n | #35 |

---

## Phase 5 — Deployment + Testing + DX

> Caddy integration, Docker, health checks, testing infrastructure, code generation.

| # | Feature | Status | Depends On | Blueprint Sections |
|---|---|---|---|---|
| 36 | Health Checks | ⬜ | #07, #09 | Health Check |
| 37 | Graceful Shutdown | ⬜ | #01 | Build and Run (with Graceful Shutdown) |
| 38 | Caddy Integration | ⬜ | #02 | Caddy Web Server (Optional) |
| 39 | Docker Deployment | ⬜ | #37 | Docker (Optional) |
| 40 | Testing Infrastructure | ⬜ | #01, #05 | Testing |
| 41 | Code Generation (CLI) | ⬜ | #10 | CLI Tools (make:controller, make:model, etc.) |

### Phase 5 — Blueprint Coverage

| Blueprint Section | Feature(s) |
|---|---|
| Health Check | #36 |
| Build and Run (with Graceful Shutdown) | #37 |
| Caddy Web Server (Optional) | #38 |
| Docker (Optional) | #39 |
| Testing | #40 |
| CLI Tools (code generation) | #41 |

---

## Phase 6 — Advanced (Future)

> Not in current scope. Tracked here for planning visibility.

| # | Feature | Status | Depends On | Blueprint Section |
|---|---|---|---|---|
| 42 | Queue Workers / Background Jobs | 🔮 | #05, #09 | Advanced Features (Planned) |
| 43 | Task Scheduler / Cron | 🔮 | #05 | Advanced Features (Planned) |
| 44 | Plugin / Module System | 🔮 | #05, #06 | Advanced Features (Planned) |
| 45 | GraphQL Support | 🔮 | #07 | Advanced Features (Planned) |
| 46 | Admin Panel Scaffolding | 🔮 | #15, #17 | Advanced Features (Planned) |
| 47 | API Versioning | 🔮 | #07 | Advanced Features (Planned) |
| 48 | WebSocket Rooms / Channels | 🔮 | #31 | Advanced Features (Planned) |
| 49 | OAuth2 / Social Login | 🔮 | #21 | Advanced Features (Planned) |
| 50 | Two-Factor Authentication (TOTP) | 🔮 | #21 | Advanced Features (Planned) |
| 51 | Audit Logging | 🔮 | #11, #03 | Advanced Features (Planned) |
| 52 | Soft Deletes | 🔮 | #11 | Advanced Features (Planned) |
| 53 | Database Read/Write Splitting | 🔮 | #09 | Advanced Features (Planned) |
| 54 | Prometheus Metrics | 🔮 | #07 | Advanced Features (Planned) |

---

## Full Blueprint Traceability Matrix

Every blueprint section mapped to its feature number. No gaps, no extras.

| Blueprint Section | Feature # | Phase |
|---|---|---|
| Suggested Project Structure | 01 | 1 |
| Configuration System | 02 | 1 |
| Logging | 03 | 1 |
| Error Handling | 04 | 1 |
| Service Container & Providers | 05, 06 | 1 |
| Router Layer | 07 | 1 |
| Middleware (Custom Registration) | 08 | 1 |
| Database Layer | 09 | 1 |
| CLI Tools (foundation) | 10 | 1 |
| Models (GORM) | 11 | 2 |
| CLI Tools (migrate commands) | 12 | 2 |
| Database Seeding | 13 | 2 |
| Database Transactions | 14 | 2 |
| MVC Controller Example | 15 | 2 |
| Response Helpers | 16 | 2 |
| View / Template Engine | 17 | 2 |
| Services Layer | 18 | 2 |
| Helpers + String & Data Helpers | 19 | 2 |
| Session Management | 20 | 2 |
| Authentication | 21 | 2 |
| Built-in Crypto & Security Utilities | 22 | 2 |
| Input Validation (Built-in) | 23 | 3 |
| CSRF Protection | 24 | 3 |
| CORS | 25 | 3 |
| Rate Limiting | 26 | 3 |
| Request ID / Tracing | 27 | 3 |
| File Upload & Storage | 28 | 3 |
| Mail / Email | 29 | 3 |
| Static File Serving | 30 | 3 |
| WebSocket Support | 31 | 3 |
| Caching | 32 | 4 |
| Pagination | 33 | 4 |
| Events / Hooks System | 34 | 4 |
| Localization / i18n | 35 | 4 |
| Health Check | 36 | 5 |
| Build and Run (Graceful Shutdown) | 37 | 5 |
| Caddy Web Server (Optional) | 38 | 5 |
| Docker (Optional) | 39 | 5 |
| Testing | 40 | 5 |
| CLI Tools (code generation) | 41 | 5 |
| Advanced Features (Planned) | 42–54 | 6 (future) |

---

## Progress Tracker

| Phase | Features | Complete | Remaining | Status |
|---|---|---|---|---|
| Phase 1 — Core Skeleton | 01–10 | 6/10 | 4 | 🟡 In progress |
| Phase 2 — MVC + Auth | 11–22 | 0/12 | 12 | ⬜ Not started |
| Phase 3 — Web Essentials | 23–31 | 0/9 | 9 | ⬜ Not started |
| Phase 4 — Caching + Events | 32–35 | 0/4 | 4 | ⬜ Not started |
| Phase 5 — Deploy + Testing + DX | 36–41 | 0/6 | 6 | ⬜ Not started |
| Phase 6 — Advanced (Future) | 42–54 | 0/13 | 13 | 🔮 Future |
| **Total (Current Scope)** | **01–41** | **4/41** | **37** | **🟡 In progress** |

---

## Post-Framework Showcase Projects

> These are applications built WITH the completed framework — not part of the framework itself.
> Do not start until Phases 1–5 are complete.

| Project | Description | Depends On |
|---|---|---|
| **RGo Docs App** | Documentation platform that parses and renders Markdown files, with admin panel for editing. Serves as the official RGo documentation site and proves the framework works end-to-end. | All Phases 1–5 |

### RGo Docs App — Planned Capabilities

- Markdown parsing and HTML rendering (via `goldmark`)
- YAML frontmatter extraction for metadata
- Auto-generated navigation from folder structure
- Full-text search across all docs
- Syntax-highlighted code blocks
- Version switching for docs across releases
- Admin panel with auth-protected CRUD editor for `.md` files
- SSR rendered via RGo's own Views & Templates system

**Reference inspiration**: VitePress (DX), Docusaurus (features), Stripe Docs (API docs UX).

---

## Notes

- **Feature numbering is final** — do not renumber. If a feature is added later, append it with the next available number.
- **Dependencies are strict** — do not start a feature until its dependencies are complete.
- **Phase 6 is out of scope** — tracked for visibility only. Do not create feature docs for Phase 6 until the current scope (Phases 1–5) is complete.
- **Each feature gets its own Mastery document set** — discussion, architecture, tasks, testplan, changelog, review (and api spec if applicable).
- **Update this document** after every feature merge — mark status, update counts.

---

> *"Ship what matters. Track what's next. Ignore what's not yet."*
