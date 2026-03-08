# RapidGo v2 — Architecture Document

> **Project**: RapidGo Framework  
> **Target**: v2.0.0  
> **Base**: v1.0.0  
> **Date**: 2026-03-08  

---

## 1. High-Level Architecture

### v1.0.0 — Monolithic Starter

```
┌─────────────────────────────────────────────────────────┐
│                   RapidGo (single repo)                 │
│                                                         │
│  cmd/main.go ──→ core/cli.Execute()                     │
│                    │                                    │
│  ┌─────────────────┼──────── HARD IMPORTS ────────────┐ │
│  │                 ▼                                  │ │
│  │  core/cli/root.go ──→ app/providers (8 providers)  │ │
│  │  core/cli/serve.go ──→ routes (RegisterWeb/API/WS) │ │
│  │  core/cli/work.go ──→ app/jobs + app/providers     │ │
│  │  core/cli/schedule_run.go ──→ app/schedule          │ │
│  │  core/cli/migrate.go ──→ database/models            │ │
│  │  core/cli/seed.go ──→ database/seeders              │ │
│  │  core/audit/audit.go ──→ database/models/AuditLog   │ │
│  └────────────────────────────────────────────────────┘ │
│                                                         │
│  app/        ← User application code                    │
│  routes/     ← Route definitions                        │
│  http/       ← Controllers, requests, responses         │
│  database/   ← Models, migrations, seeders              │
│  resources/  ← Views, lang, static                      │
│  plugins/    ← Plugin implementations                   │
│  storage/    ← Runtime dirs (logs, cache, sessions)     │
│                                                         │
│  Problem: core/ can't compile without app/ and routes/  │
└─────────────────────────────────────────────────────────┘
```

### v2.0.0 — Two-Repository Architecture

```
┌──────────────────────────────────────────────────────────┐
│              RapidGo Library (importable)                 │
│              go get github.com/RAiWorks/RapidGo          │
│                                                          │
│  core/                                                   │
│  ├── app/          Application lifecycle                 │
│  ├── auth/         JWT authentication                    │
│  ├── audit/        Audit logging + AuditLog model        │
│  ├── cache/        File + Redis caching                  │
│  ├── cli/          Cobra commands + hooks.go             │
│  ├── config/       Configuration + environment           │
│  ├── container/    Service container (IoC)               │
│  ├── crypto/       Encryption utilities                  │
│  ├── errors/       Error handling                        │
│  ├── events/       Event dispatcher                      │
│  ├── graphql/      GraphQL support                       │
│  ├── health/       Health checks                         │
│  ├── i18n/         Internationalization                  │
│  ├── logger/       Structured logging                    │
│  ├── mail/         Email sending                         │
│  ├── metrics/      Prometheus metrics                    │
│  ├── middleware/    HTTP middleware                       │
│  ├── oauth/        OAuth2 social login                   │
│  ├── plugin/       Plugin system                         │
│  ├── queue/        Job queue + workers                   │
│  ├── router/       Gin-based router                      │
│  ├── scheduler/    Cron task scheduler                   │
│  ├── server/       HTTP server                           │
│  ├── service/      Service mode (Web/API/WS)             │
│  ├── session/      Session management                    │
│  ├── storage/      File storage (local + S3)             │
│  ├── totp/         TOTP 2FA                              │
│  ├── validation/   Input validation                      │
│  └── websocket/    WebSocket rooms                       │
│                                                          │
│  database/                                               │
│  ├── connection.go     DB connection factory             │
│  ├── resolver.go       Read/write splitting              │
│  ├── transaction.go    Transaction helpers               │
│  ├── models/                                             │
│  │   ├── base.go       BaseModel (ID, timestamps, soft) │
│  │   └── scopes.go     WithTrashed, OnlyTrashed          │
│  └── migrations/                                         │
│      └── migrator.go   Migration engine (generic)        │
│                                                          │
│  testing/testutil/     Test utilities                    │
│                                                          │
│  ZERO imports from app/, routes/, http/, plugins/        │
└──────────────────────────────────────────────────────────┘
        ▲
        │  go get
        │
┌──────────────────────────────────────────────────────────┐
│            RapidGo-Starter (template project)            │
│            github.com/RAiWorks/RapidGo-starter           │
│                                                          │
│  cmd/main.go ───→ cli.Set*() hooks wire app → framework │
│                                                          │
│  app/                                                    │
│  ├── providers/    8 service providers                   │
│  ├── helpers/      App-specific utilities                │
│  ├── services/     Business logic                        │
│  ├── jobs/         Queue job handlers                    │
│  ├── schedule/     Scheduled tasks                       │
│  └── plugins.go    Plugin registration                   │
│                                                          │
│  routes/           Web, API, WS route definitions        │
│  http/             Controllers, requests, responses      │
│                                                          │
│  database/                                               │
│  ├── models/       App models (User, Post, registry)     │
│  ├── migrations/   App migration files (init() + Register) │
│  └── seeders/      App seeder implementations            │
│                                                          │
│  resources/        Views, lang, static assets            │
│  storage/          Runtime dirs (logs, cache, uploads)    │
│  plugins/          Plugin implementations                │
│                                                          │
│  .env.example, Dockerfile, docker-compose.yml            │
│  Caddyfile, Makefile                                     │
└──────────────────────────────────────────────────────────┘
```

---

## 2. Hook Architecture

The 6 callback hooks in `core/cli/hooks.go` are the bridge between library and application. The library defines the interfaces; the starter wires in the concrete implementations.

```
┌─────────────────────────────────────────────────┐
│                 Starter main.go                  │
│                                                  │
│  cli.SetBootstrap(...)       ─┐                  │
│  cli.SetRoutes(...)           │ 6 hooks          │
│  cli.SetJobRegistrar(...)     │ wired at startup  │
│  cli.SetScheduleRegistrar(...)│                  │
│  cli.SetModelRegistry(...)    │                  │
│  cli.SetSeeder(...)          ─┘                  │
│                                                  │
│  cli.Execute()                                   │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│               Library core/cli/                  │
│                                                  │
│  hooks.go:                                       │
│    var bootstrapFn       BootstrapFunc           │
│    var routeRegistrar    RouteRegistrar           │
│    var jobRegistrar      JobRegistrar             │
│    var scheduleRegistrar ScheduleRegistrar        │
│    var modelRegistryFn   ModelRegistryFunc        │
│    var seederFn          SeederFunc               │
│                                                  │
│  Commands call hooks:                            │
│    root.go:         bootstrapFn(app, mode)       │
│    serve.go:        routeRegistrar(r, c, mode)   │
│    work.go:         jobRegistrar()               │
│    schedule_run.go: scheduleRegistrar(s, app)    │
│    migrate.go:      modelRegistryFn()            │
│    seed.go:         seederFn(db, name)           │
└─────────────────────────────────────────────────┘
```

### Hook Type Definitions

| Hook | Signature | Called By | Purpose |
|------|-----------|-----------|---------|
| `SetBootstrap` | `func(a *app.App, mode service.Mode)` | `root.go` → `NewApp()` | Register service providers |
| `SetRoutes` | `func(r *router.Router, c *container.Container, mode service.Mode)` | `serve.go` → `applyRoutesForMode()` | Register web/API/WS routes |
| `SetJobRegistrar` | `func()` | `work.go` | Register queue job handlers |
| `SetScheduleRegistrar` | `func(s *scheduler.Scheduler, a *app.App)` | `schedule_run.go` | Register scheduled tasks |
| `SetModelRegistry` | `func() []interface{}` | `migrate.go` | Return model list for AutoMigrate |
| `SetSeeder` | `func(db *gorm.DB, name string) error` | `seed.go` | Run database seeders |

---

## 3. Package Dependency Graph

### Before (v1 — circular dependencies)

```
cmd/main.go
  └→ core/cli
       ├→ core/app
       ├→ core/config
       ├→ core/container
       ├→ core/router
       ├→ core/server
       ├→ core/service
       ├→ core/health
       │
       ├→ app/providers  ←── COUPLING (8 providers)
       │    ├→ core/config
       │    ├→ core/logger
       │    ├→ core/session
       │    ├→ core/router
       │    └→ routes/     ←── nested coupling
       │
       ├→ app/jobs       ←── COUPLING (RegisterJobs)
       ├→ app/schedule   ←── COUPLING (RegisterSchedule)
       ├→ routes/        ←── COUPLING (RegisterWeb/API/WS)
       │
       ├→ database/models     ←── COUPLING (models.All())
       ├→ database/migrations ←── OK (engine stays)
       └→ database/seeders    ←── OK (engine stays)

  core/audit
       └→ database/models  ←── COUPLING (AuditLog struct)
```

### After (v2 — clean dependency flow)

```
Starter: cmd/main.go
  ├→ core/cli           (Set*() hooks)
  ├→ app/providers      (passed to SetBootstrap)
  ├→ routes/            (passed to SetRoutes)
  ├→ app/jobs           (passed to SetJobRegistrar)
  ├→ app/schedule       (passed to SetScheduleRegistrar)
  ├→ database/models    (passed to SetModelRegistry)
  └→ database/seeders   (passed to SetSeeder)

Library: core/cli
  ├→ core/app
  ├→ core/config
  ├→ core/container
  ├→ core/router
  ├→ core/server
  ├→ core/service
  ├→ core/health
  ├→ database/migrations  (engine — NewMigrator, Run, Rollback)
  └→ [hooks]              (callbacks — NO hard imports to app/)

Library: core/audit
  └→ core/audit/model.go  (AuditLog is local — no database/models import)

Library: database/
  ├→ connection.go   → gorm
  ├→ resolver.go     → gorm
  ├→ transaction.go  → gorm
  ├→ models/base.go  → gorm  (BaseModel)
  ├→ models/scopes.go → gorm (WithTrashed, OnlyTrashed)
  └→ migrations/migrator.go → gorm  (migration engine)
```

---

## 4. Public API Surface

These are the packages and types that starter apps will import:

### Core Packages (imported as `github.com/RAiWorks/RapidGo/core/...`)

| Package | Key Exports | Used By Starter |
|---------|------------|:--------------:|
| `app` | `App`, `New()` | Yes (via SetBootstrap) |
| `auth` | JWT functions | Yes |
| `audit` | `Logger`, `NewLogger()`, `Entry`, `AuditLog` | Yes |
| `cache` | `Cache`, `FileStore`, `RedisStore` | Yes |
| `cli` | `Execute()`, `RootCmd()`, `Set*()` hooks | Yes (main.go) |
| `config` | `Load()`, `Env()`, `Get()` | Yes |
| `container` | `Container`, `MustMake[]()`, `Singleton()` | Yes |
| `crypto` | `Encrypt()`, `Decrypt()`, `Hash()` | Yes |
| `errors` | `AppError`, `Errors` | Yes |
| `events` | `Dispatcher`, `Listen()`, `Dispatch()` | Yes |
| `graphql` | `Handler()`, `NewSchema()` | Optional |
| `health` | `Routes()` | Auto-registered |
| `i18n` | `T()`, `Translator` | Yes |
| `logger` | `New()`, `Logger` | Yes |
| `mail` | `Mailer`, `Send()` | Yes |
| `metrics` | `Collector`, `Routes()` | Optional |
| `middleware` | All middleware functions | Yes |
| `oauth` | `Provider`, `Redirect()`, `Callback()` | Optional |
| `plugin` | `Plugin`, `Manager` | Optional |
| `queue` | `Dispatcher`, `Job`, `Worker`, drivers | Yes |
| `router` | `Router`, `Route`, `Group()` | Yes (via SetRoutes) |
| `scheduler` | `Scheduler`, `New()`, `Add()` | Yes (via SetScheduleRegistrar) |
| `server` | `ListenAndServe()`, `ListenAndServeMulti()` | Internal to cli |
| `service` | `Mode`, `ModeWeb`, `ModeAPI`, `ModeWS` | Yes (via SetBootstrap) |
| `session` | `Store`, sessions | Yes |
| `storage` | `Disk`, `Local`, `S3` | Yes |
| `totp` | `GenerateSecret()`, `Validate()` | Optional |
| `validation` | `Validate()`, rules | Yes |
| `websocket` | `Hub`, `Room`, `Client` | Optional |

### Database Packages (imported as `github.com/RAiWorks/RapidGo/database/...`)

| Package | Key Exports | Used By Starter |
|---------|------------|:--------------:|
| `database` | `Connect()`, `NewResolver()`, `Transaction()` | Via providers |
| `models` | `BaseModel` | Yes (embedded in app models) |
| `models` | `WithTrashed()`, `OnlyTrashed()` | Yes (GORM scopes) |
| `migrations` | `Register()`, `Migration`, `Migrator` | Yes (migration files) |
| `seeders` | `Register()`, `Seeder`, `RunAll()`, `RunByName()` | Yes (seeder files) |

### Testing Package

| Package | Key Exports | Used By Starter |
|---------|------------|:--------------:|
| `testing/testutil` | Test helpers | Yes (test files) |

---

## 5. Data Flow Diagrams

### Serve Command Flow (v2)

```
main.go                    Library (core/cli)              App Code
────────                   ──────────────────              ────────
                          
cli.SetBootstrap(fn)  ───→ stores bootstrapFn
cli.SetRoutes(fn)     ───→ stores routeRegistrar
cli.Execute()         ───→ rootCmd.Execute()
                           │
                           ▼
                      serveCmd.RunE()
                           │
                      NewApp(mode)
                           │
                      bootstrapFn(app, mode)  ──────→  registers 8 providers
                           │                            (Config, Logger, DB,
                      app.Boot()                         Redis, Queue, Session,
                           │                             Middleware, Router)
                      applyRoutesForMode(r,c,mode)
                           │
                      routeRegistrar(r,c,mode) ────→  registers web/api/ws routes
                           │
                      server.ListenAndServe()
```

### Migrate Command Flow (v2)

```
main.go                    Library (core/cli)              App Code
────────                   ──────────────────              ────────

cli.SetBootstrap(fn) ────→ stores bootstrapFn
cli.SetModelRegistry(fn)─→ stores modelRegistryFn
cli.Execute()        ────→ rootCmd.Execute()
                           │
                           ▼
                      migrateCmd.RunE()
                           │
                      NewApp(ModeAll)
                           │
                      bootstrapFn(app, ModeAll) ──→  registers providers
                           │
                      app.Boot()
                           │
                      db := resolve "db"
                           │
                      modelRegistryFn() ──────────→  returns []interface{}{&User{}, &Post{}}
                           │
                      db.AutoMigrate(models...)
                           │
                      migrations.NewMigrator(db)
                           │
                      migrator.Run()   ←── reads global registry
                                           (populated by init() in starter's
                                            migration files via blank import)
```

### Worker Command Flow (v2)

```
main.go                    Library (core/cli)              App Code
────────                   ──────────────────              ────────

cli.SetBootstrap(fn) ────→ stores bootstrapFn
cli.SetJobRegistrar(fn)──→ stores jobRegistrar
cli.Execute()        ────→ rootCmd.Execute()
                           │
                           ▼
                      workCmd.RunE()
                           │
                      NewApp(ModeAll)
                           │
                      bootstrapFn(app, ModeAll) ──→  registers providers
                           │                         (Config, Logger, DB,
                      app.Boot()                      Redis, Queue)
                           │
                      jobRegistrar() ─────────────→  registers job handlers
                           │                          with queue dispatcher
                      resolve dispatcher
                           │
                      worker.Run(ctx)
```

---

## 6. Starter `main.go` — Target State

```go
package main

import (
    "github.com/RAiWorks/RapidGo/core/app"
    "github.com/RAiWorks/RapidGo/core/cli"
    "github.com/RAiWorks/RapidGo/core/container"
    "github.com/RAiWorks/RapidGo/core/router"
    "github.com/RAiWorks/RapidGo/core/scheduler"
    "github.com/RAiWorks/RapidGo/core/service"
    "gorm.io/gorm"

    "myapp/app/jobs"
    "myapp/app/providers"
    "myapp/app/schedule"
    "myapp/database/models"
    "myapp/database/seeders"
    "myapp/routes"

    // Blank import to trigger init() migration registration
    _ "myapp/database/migrations"
)

func main() {
    // 1. Bootstrap — register service providers
    cli.SetBootstrap(func(a *app.App, mode service.Mode) {
        a.Register(&providers.ConfigProvider{})
        a.Register(&providers.LoggerProvider{})
        if mode.Has(service.ModeWeb) || mode.Has(service.ModeAPI) || mode.Has(service.ModeWS) {
            a.Register(&providers.DatabaseProvider{})
        }
        a.Register(&providers.RedisProvider{})
        a.Register(&providers.QueueProvider{})
        if mode.Has(service.ModeWeb) {
            a.Register(&providers.SessionProvider{})
        }
        a.Register(&providers.MiddlewareProvider{Mode: mode})
        a.Register(&providers.RouterProvider{Mode: mode})
    })

    // 2. Routes
    cli.SetRoutes(func(r *router.Router, c *container.Container, mode service.Mode) {
        if mode.Has(service.ModeWeb) {
            routes.RegisterWeb(r)
        }
        if mode.Has(service.ModeAPI) {
            routes.RegisterAPI(r)
        }
        if mode.Has(service.ModeWS) {
            routes.RegisterWS(r)
        }
    })

    // 3. Background jobs
    cli.SetJobRegistrar(jobs.RegisterJobs)

    // 4. Scheduled tasks
    cli.SetScheduleRegistrar(schedule.RegisterSchedule)

    // 5. Database models
    cli.SetModelRegistry(models.All)

    // 6. Database seeders
    cli.SetSeeder(func(db *gorm.DB, name string) error {
        if name != "" {
            return seeders.RunByName(db, name)
        }
        return seeders.RunAll(db)
    })

    cli.Execute()
}
```

---

## 7. Import Alias Convention

The starter's `database/models` package shares a name with the library's. Go handles this cleanly with import aliases:

```go
import (
    fwmodels "github.com/RAiWorks/RapidGo/database/models"  // library's BaseModel
    "myapp/database/models"                                   // app's User, Post
)

type User struct {
    fwmodels.BaseModel
    Name  string `gorm:"size:255" json:"name"`
    Email string `gorm:"size:255;uniqueIndex" json:"email"`
}
```

If this becomes confusing, the library's package could be renamed to `database/base` in a future minor version. For v2.0.0, import aliases are sufficient.
