# RapidGo — Framework Comparison

How RapidGo compares to Go HTTP routers and application frameworks.

---

## Understanding the Categories

Go frameworks fall into two categories:

| Category | Purpose | Examples |
|---|---|---|
| **HTTP Router / Toolkit** | Handles HTTP requests and routing. You assemble everything else yourself. | Gin, Echo, Fiber, Chi, gorilla/mux |
| **Application Framework** | Full application stack: routing + ORM + auth + sessions + queues + CLI + mail + events + ... | **RapidGo**, Beego, Buffalo (archived) |

**RapidGo is built ON Gin** — the same relationship as Laravel on Symfony, NestJS on Express, or Rails on Rack. They are not competitors; they operate at different abstraction levels.

---

## Feature Matrix

| Feature | **RapidGo** | Gin | Echo | Fiber | Go Kit |
|---|---|---|---|---|---|
| **Core** | | | | | |
| HTTP Router | ✅ (via Gin) | ✅ | ✅ | ✅ | ✅ |
| Middleware Pipeline | ✅ | ✅ | ✅ | ✅ | ✅ |
| Dependency Injection Container | ✅ | ❌ | ❌ | ❌ | ❌ |
| Service Providers (Register/Boot) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Configuration (.env) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Structured Logging (slog) | ✅ | ❌ | ❌ | ❌ | ✅ |
| Plugin / Module System | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Data & Database** | | | | | |
| ORM (GORM) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Migrations (up/down) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Seeders | ✅ | ❌ | ❌ | ❌ | ❌ |
| Transactions (auto + nested) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Pagination | ✅ | ❌ | ❌ | ❌ | ❌ |
| Soft Deletes | ✅ | ❌ | ❌ | ❌ | ❌ |
| Read/Write Splitting | ✅ | ❌ | ❌ | ❌ | ❌ |
| **HTTP & Routing** | | | | | |
| Route Groups | ✅ | ✅ | ✅ | ✅ | ❌ |
| Resource Routes (7 CRUD) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Named Routes + URL Generation | ✅ | ❌ | ✅ | ❌ | ❌ |
| API Versioning | ✅ | ❌ | ❌ | ❌ | ❌ |
| MVC Controllers | ✅ | ❌ | ❌ | ❌ | ❌ |
| Input Validation | ✅ | Partial | ✅ | ❌ | ❌ |
| Response Helpers (envelope) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Views & Templates | ✅ | ✅ | ✅ | ✅ | ❌ |
| Static File Serving | ✅ | ✅ | ✅ | ✅ | ❌ |
| GraphQL | ✅ | ❌ | ❌ | ❌ | ❌ |
| WebSocket | ✅ | ❌ | ❌ | ✅ | ❌ |
| WebSocket Rooms/Channels | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Security & Auth** | | | | | |
| JWT Authentication | ✅ | ❌ | ❌ | ❌ | ❌ |
| Session-Based Auth (5 backends) | ✅ | ❌ | ❌ | ❌ | ❌ |
| OAuth2 / Social Login | ✅ | ❌ | ❌ | ❌ | ❌ |
| TOTP Two-Factor Auth | ✅ | ❌ | ❌ | ❌ | ❌ |
| CSRF Protection | ✅ | ❌ | ✅ | ✅ | ❌ |
| CORS Configuration | ✅ | ❌* | ✅ | ✅ | ❌ |
| Rate Limiting | ✅ | ❌ | ✅ | ✅ | ✅ |
| Crypto Utilities (AES, bcrypt) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Audit Logging | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Infrastructure** | | | | | |
| Queue Workers (4 drivers) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Task Scheduler (cron) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Event System (pub-sub) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Caching (3 backends) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Mail (SMTP) | ✅ | ❌ | ❌ | ❌ | ❌ |
| File Storage (local + S3) | ✅ | ❌ | ❌ | ❌ | ❌ |
| i18n / Localization | ✅ | ❌ | ❌ | ❌ | ❌ |
| Prometheus Metrics | ✅ | ❌ | ❌ | ✅ | ✅ |
| **CLI & DX** | | | | | |
| Code Generation (make:*) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Database CLI (migrate, seed) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Admin Panel Scaffolding | ✅ | ❌ | ❌ | ❌ | ❌ |
| Testing Utilities | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Deployment** | | | | | |
| Graceful Shutdown | ✅ | ✅ | ✅ | ✅ | ❌ |
| Health Checks | ✅ | ❌ | ❌ | ❌ | ❌ |
| Docker Support | ✅ | ❌ | ❌ | ❌ | ❌ |
| Caddy Integration | ✅ | ❌ | ❌ | ❌ | ❌ |
| Multi-Port Serving | ✅ | ❌ | ❌ | ❌ | ❌ |

\* Gin requires the separate `gin-contrib/cors` package.

---

## Feature Count Comparison

| Framework | Built-in Features | Category |
|---|---|---|
| **RapidGo** | **56** | Application Framework |
| Gin | ~8 | HTTP Router |
| Echo | ~12 | HTTP Router |
| Fiber | ~10 | HTTP Router |
| Go Kit | ~6 | Microservice Toolkit |
| Beego | ~20 | Application Framework |
| Buffalo | ~18 | Application Framework (archived) |

---

## Cross-Language Comparison

RapidGo's feature set is comparable to full application frameworks in other languages:

| Feature | RapidGo (Go) | Laravel (PHP) | NestJS (Node) | Django (Python) | Rails (Ruby) |
|---|---|---|---|---|---|
| DI Container | ✅ | ✅ | ✅ | ❌ | ❌ |
| ORM + Migrations | ✅ | ✅ | ✅ | ✅ | ✅ |
| Auth (JWT + Sessions) | ✅ | ✅ | ✅ | ✅ | ✅ |
| OAuth2 / Social | ✅ | ✅ | ✅ | ✅ | ✅ |
| TOTP 2FA | ✅ | ✅ | ❌ | ❌ | ❌ |
| Queue Workers | ✅ | ✅ | ✅ | ✅ | ✅ |
| Task Scheduler | ✅ | ✅ | ✅ | ❌* | ❌* |
| Event System | ✅ | ✅ | ✅ | ✅ | ❌ |
| Mail | ✅ | ✅ | ✅ | ✅ | ✅ |
| Cache (multi-backend) | ✅ | ✅ | ✅ | ✅ | ✅ |
| File Storage (S3) | ✅ | ✅ | ❌ | ❌ | ✅ |
| GraphQL | ✅ | ❌* | ✅ | ✅ | ❌* |
| WebSocket Rooms | ✅ | ✅ | ✅ | ❌* | ✅ |
| Plugin System | ✅ | ✅ | ✅ | ✅ | ✅ |
| CLI Scaffolding | ✅ | ✅ | ✅ | ✅ | ✅ |
| Admin Scaffolding | ✅ | ❌* | ❌ | ✅ | ❌* |
| Prometheus Metrics | ✅ | ❌ | ❌ | ❌ | ❌ |
| i18n | ✅ | ✅ | ✅ | ✅ | ✅ |
| Audit Logging | ✅ | ❌ | ❌ | ❌ | ❌ |

\* Available via community packages, not built-in.

---

## When to Use RapidGo

**Use RapidGo when you need:**
- A full web application in Go (not just an API)
- Authentication, sessions, queues, mail, events out of the box
- Convention-over-configuration (like Rails or Laravel)
- One `go get` instead of wiring 15+ packages
- CLI scaffolding for rapid development
- Go's performance with a productive developer experience

**Consider alternatives when you need:**
- A minimal HTTP API with 2–3 routes → use Gin or Echo directly
- A distributed microservice mesh → use Go Kit, Kratos, or Dapr
- Maximum control over every dependency → assemble packages manually
- A non-web application (pure CLI tool, data pipeline) → use Cobra directly

---

## How RapidGo Uses Gin

RapidGo does **not replace** Gin. It wraps and extends it:

```
RapidGo Application Framework
├── core/router     → wraps gin.Engine (adds groups, resources, named routes)
├── core/middleware  → wraps gin.HandlerFunc (adds registry, aliases)
├── core/server     → wraps http.Server (adds graceful shutdown, multi-port)
└── ... 50+ more packages that Gin doesn't provide
```

Every Gin feature remains accessible. RapidGo adds the application layer that Gin intentionally omits.
