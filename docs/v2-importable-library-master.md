# RapidGo v2 — Importable Library Split: Master Document

> **Project**: RapidGo Framework  
> **Version**: v2.0.0 (target)  
> **Base**: v1.0.0 (tagged 2026-03-07, commit `9c0a22a`)  
> **Author**: RAi Works  
> **Date**: 2026-03-08  
> **Status**: Approved — ready for implementation  

---

## 1. Objective

Transform RapidGo from a **monolithic starter** (clone-and-build-inside) into an **importable Go library** (`go get`) with a companion **starter template** (`RapidGo-starter`).

**v1.0.0** remains frozen on `main`. All v2 work happens on the `v2` branch, which becomes the default branch when complete.

---

## 2. Two-Repository Architecture

```
github.com/RAiWorks/RapidGo           ← Importable framework library (go get)
github.com/RAiWorks/RapidGo-starter   ← Scaffold project (clone → build inside)
```

### What Each Repo Contains

| Repo | Contains | Purpose |
|------|----------|---------|
| **RapidGo** (library) | `core/*`, `database/connection.go`, `database/transaction.go`, `database/migrations/migrator.go`, `database/models/base.go`, `testing/testutil/`, `go.mod`, `LICENSE`, `README.md` | Importable framework — `go get github.com/RAiWorks/RapidGo` |
| **RapidGo-starter** (template) | `cmd/`, `app/`, `routes/`, `http/`, `database/models/` (app models), `database/migrations/` (app migrations), `database/seeders/`, `resources/`, `storage/`, `plugins/`, `.env.example`, `Dockerfile`, `docker-compose.yml`, `Caddyfile`, `Makefile` | Clone or scaffold → customize → build |

---

## 3. Coupling Analysis — Complete (Verified Against Codebase)

### 3.1 All 7 True Coupling Points

The original split plan claimed 5 coupling points. A full codebase grep found 13 imports from `core/` to app-specific packages, but several are **not real coupling** because the imported engine stays in the library.

| # | File | Hard Import(s) | What It Does | Solution |
|---|------|---------------|--------------|----------|
| C1 | `core/cli/root.go` | `app/providers` | `NewApp()` hard-codes all 8 provider registrations | `SetBootstrap(fn)` callback |
| C2 | `core/cli/serve.go` | `routes` | `applyRoutesForMode()` calls `routes.RegisterWeb/API/WS()` | `SetRoutes(fn)` callback |
| C3 | `core/cli/work.go` | `app/jobs`, `app/providers` | Worker bootstrap + `jobs.RegisterJobs()` | `SetJobRegistrar(fn)` callback + bootstrap via `SetBootstrap` |
| C4 | `core/cli/schedule_run.go` | `app/providers`, `app/schedule` | Scheduler bootstrap + `schedule.RegisterSchedule()` | `SetScheduleRegistrar(fn)` callback + bootstrap via `SetBootstrap` |
| C5 | `core/cli/migrate.go` | `database/models` | `models.All()` for AutoMigrate | `SetModelRegistry(fn)` callback |
| C6 | `core/cli/seed.go` | `database/seeders` | `seeders.RunByName()` / `seeders.RunAll()` | `SetSeeder(fn)` callback |
| C7 | `core/audit/audit.go` | `database/models` | Uses `models.AuditLog` struct directly | Move `AuditLog` into `core/audit/` as framework-owned model |

### 3.1.1 Imports That Are NOT Coupling Points

These `core/` files import `database/migrations` or `database/seeders`, but the engine code stays in the library. These imports are **valid after the split** and need no hooks:

| File | Import | Why It's Fine |
|------|--------|---------------|
| `core/cli/migrate.go` | `database/migrations` | Imports `migrations.NewMigrator()` — the engine stays in the library |
| `core/cli/migrate_rollback.go` | `database/migrations` | Same — uses engine's `Rollback()` method |
| `core/cli/migrate_status.go` | `database/migrations` | Same — uses engine's `Status()` method |
| `core/cli/seed.go` | `database/seeders` | Imports `seeders.RunAll()` / `RunByName()` — the engine stays in the library |

App-specific migration files and seeder implementations move to the starter. They register themselves by calling the library's `migrations.Register()` / `seeders.Register()` via `init()` + blank import in the starter's `main.go`.

### 3.2 App-Side Coupling (Moves to Starter — No Refactoring Needed)

| File | Hard Import | Resolution |
|------|------------|------------|
| `app/providers/router_provider.go` | `routes` | Moves entirely to starter |
| `app/plugins.go` | `plugins/example` | Moves entirely to starter |
| `database/models/registry.go` | `database/models/*` (same package) | Moves entirely to starter |
| `database/models/user.go` | `app/helpers` (for `HashPassword`) | Both move to starter — no issue |
| `routes/web.go`, `routes/api.go` | `http/controllers` | Both move to starter — no issue |
| `app/services/user_service.go` | `database/models` | Both move to starter — no issue |

### 3.2.1 Test Files With App-Model Coupling (Must Refactor During Split)

These test files stay in the library (they test library code) but reference app-specific models (`User`, `Post`, `AuditLog`). They **will not compile** after app models are removed.

| Test File | Uses | Fix |
|-----------|------|-----|
| `database/models/models_test.go` | `User{}`, `Post{}` | Split: `BaseModel` tests stay (use test-only struct), `User`/`Post` tests move to starter |
| `database/models/scopes_test.go` | `User{}`, `Post{}` | Refactor: replace with test-only `testModel` struct that embeds `BaseModel` |
| `database/migrations/migrations_test.go` | `models.User`, `models.Post` | Refactor: use test-only model struct |
| `database/seeders/seeders_test.go` | `models.User` (in `setupTestDB()`) | Refactor: `setupTestDB()` uses test-only model; engine tests already use `mockSeeder` |
| `core/audit/audit_test.go` | `models.AuditLog` | Fixed by C7 (AuditLog moves to `core/audit/`) |

### 3.3 Packages Verified Clean (No Changes Needed)

All `core/` packages below have **zero** imports from `app/`, `routes/`, `http/`, `plugins/`, or `database/models/`:

| Package | Depends On |
|---------|-----------|
| `core/app` | `core/container` |
| `core/container` | (none) |
| `core/config` | `godotenv` |
| `core/logger` | `log/slog` |
| `core/errors` | (none) |
| `core/router` | `gin` |
| `core/middleware` | `gin`, `core/auth`, `core/session` |
| `core/auth` | `golang-jwt` |
| `core/crypto` | `golang.org/x/crypto` |
| `core/session` | `gorilla/sessions` |
| `core/cache` | `go-redis` |
| `core/mail` | `go-mail` |
| `core/events` | (none) |
| `core/i18n` | (none) |
| `core/health` | `gorm`, `core/router` |
| `core/server` | (none) |
| `core/websocket` | `coder/websocket`, `google/uuid` |
| `core/validation` | `gin` |
| `core/storage` | `aws-sdk-go-v2/s3` |
| `core/queue` | `gorm`, `go-redis` |
| `core/scheduler` | `robfig/cron` |
| `core/graphql` | `graphql-go` |
| `core/totp` | `pquerna/otp` |
| `core/plugin` | `core/container`, `cobra` |
| `core/service` | (none) |
| `core/metrics` | `prometheus/client_golang` |
| `core/oauth` | `golang.org/x/oauth2` |
| `database/connection.go` | `core/config`, GORM drivers |
| `database/transaction.go` | `gorm` |
| `database/resolver.go` | `gorm` |
| `database/migrations/migrator.go` | `gorm` (engine only — no app migrations) |
| `database/models/base.go` | `gorm` |
| `database/models/scopes.go` | `gorm` |

---

## 4. Callback Hook Design

### 4.1 All 6 Hooks Required

```go
// core/cli/hooks.go — Central hook registry

package cli

import (
    "github.com/RAiWorks/RapidGo/core/app"
    "github.com/RAiWorks/RapidGo/core/container"
    "github.com/RAiWorks/RapidGo/core/router"
    "github.com/RAiWorks/RapidGo/core/scheduler"
    "github.com/RAiWorks/RapidGo/core/service"
    "gorm.io/gorm"
)

// BootstrapFunc registers service providers on the application.
type BootstrapFunc func(a *app.App, mode service.Mode)

// RouteRegistrar registers routes on the router for a given mode.
type RouteRegistrar func(r *router.Router, c *container.Container, mode service.Mode)

// JobRegistrar registers application job handlers with the queue dispatcher.
type JobRegistrar func()

// ScheduleRegistrar registers scheduled tasks.
type ScheduleRegistrar func(s *scheduler.Scheduler, a *app.App)

// ModelRegistryFunc returns all model structs for AutoMigrate.
type ModelRegistryFunc func() []interface{}

// SeederFunc runs database seeders. If name is empty, runs all seeders.
type SeederFunc func(db *gorm.DB, name string) error

var (
    bootstrapFn       BootstrapFunc
    routeRegistrar    RouteRegistrar
    jobRegistrar      JobRegistrar
    scheduleRegistrar ScheduleRegistrar
    modelRegistryFn   ModelRegistryFunc
    seederFn          SeederFunc
)

func SetBootstrap(fn BootstrapFunc)             { bootstrapFn = fn }
func SetRoutes(fn RouteRegistrar)               { routeRegistrar = fn }
func SetJobRegistrar(fn JobRegistrar)           { jobRegistrar = fn }
func SetScheduleRegistrar(fn ScheduleRegistrar) { scheduleRegistrar = fn }
func SetModelRegistry(fn ModelRegistryFunc)     { modelRegistryFn = fn }
func SetSeeder(fn SeederFunc)                   { seederFn = fn }
```

### 4.1.1 Why NOT 8 Hooks

The original draft proposed 8 hooks. Two were removed after cross-checking the actual code:

| Removed Hook | Why |
|-------------|-----|
| `SetMigrationRegistrar` | Not needed. The migration engine (`migrator.go`) stays in the library. `migrate.go`, `migrate_rollback.go`, `migrate_status.go` import `database/migrations` for the engine — this is a valid library-internal import. App migration files self-register via `init()` + blank import in the starter's `main.go`. |
| `SetPluginRegistrar` | Not needed. Zero coupling from `core/` → `plugins/`. The plugin system (`core/plugin/`) has no imports from `app/` or `plugins/`. `app/plugins.go` imports `core/plugin` inward — no hook required. |

### 4.2 How Each Coupling Point Uses Its Hook

#### C1 — `root.go` → `SetBootstrap`

```go
// BEFORE (v1):
func NewApp(mode service.Mode) *app.App {
    application := app.New()
    application.Register(&providers.ConfigProvider{})  // hard import
    application.Register(&providers.LoggerProvider{})   // hard import
    // ... 6 more hard-coded providers
    application.Boot()
    return application
}

// AFTER (v2):
func NewApp(mode service.Mode) *app.App {
    application := app.New()
    if bootstrapFn != nil {
        bootstrapFn(application, mode)
    }
    application.Boot()
    return application
}
```

#### C2 — `serve.go` → `SetRoutes`

```go
// BEFORE (v1):
func applyRoutesForMode(r *router.Router, c *container.Container, m service.Mode) {
    if m.Has(service.ModeWeb) {
        // ... template/static setup ...
        routes.RegisterWeb(r)  // hard import
    }
    if m.Has(service.ModeAPI) {
        routes.RegisterAPI(r)  // hard import
    }
    if m.Has(service.ModeWS) {
        routes.RegisterWS(r)  // hard import
    }
    // health check stays
}

// AFTER (v2):
func applyRoutesForMode(r *router.Router, c *container.Container, m service.Mode) {
    if routeRegistrar != nil {
        routeRegistrar(r, c, m)
    }
    // Health check stays in framework — it's generic
    if c.Has("db") {
        health.Routes(r, func() *gorm.DB {
            return container.MustMake[*gorm.DB](c, "db")
        })
    }
}
```

#### C3 — `work.go` → `SetBootstrap` + `SetJobRegistrar`

```go
// BEFORE (v1):
application := app.New()
application.Register(&providers.ConfigProvider{})   // hard import
application.Register(&providers.LoggerProvider{})    // hard import
application.Register(&providers.DatabaseProvider{})  // hard import
application.Register(&providers.RedisProvider{})     // hard import
application.Register(&providers.QueueProvider{})     // hard import
application.Boot()
jobs.RegisterJobs()  // hard import

// AFTER (v2):
application := NewApp(service.ModeAll)  // uses SetBootstrap
if jobRegistrar != nil {
    jobRegistrar()
}
```

> **Note**: `work.go` currently has its own manual bootstrap (5 providers). After v2, it reuses `NewApp()` which calls `bootstrapFn`. The starter's bootstrap function handles mode-dependent provider selection.

#### C4 — `schedule_run.go` → `SetBootstrap` + `SetScheduleRegistrar`

```go
// BEFORE (v1):
application := app.New()
application.Register(&providers.ConfigProvider{})   // hard import
// ... 4 more hard-coded providers
application.Boot()
schedule.RegisterSchedule(s, application)  // hard import

// AFTER (v2):
application := NewApp(service.ModeAll)  // uses SetBootstrap
s := scheduler.New(slog.Default())
if scheduleRegistrar != nil {
    scheduleRegistrar(s, application)
}
```

#### C5 — `seed.go` → `SetSeeder`

```go
// BEFORE (v1):
seeders.RunByName(db, name)  // hard import
seeders.RunAll(db)           // hard import

// AFTER (v2):
if seederFn != nil {
    return seederFn(db, name)  // name="" means run all
}
```

#### C5 — `migrate.go` → `SetModelRegistry` (Only the `models.All()` import is coupling)

```go
// BEFORE (v1 — migrate.go):
db.AutoMigrate(models.All()...)           // hard import → COUPLING (needs hook)
migrator, _ := migrations.NewMigrator(db) // NOT coupling (engine stays in library)

// AFTER (v2 — migrate.go):
if modelRegistryFn != nil {
    db.AutoMigrate(modelRegistryFn()...)
}
migrator := migrations.NewMigrator(db)  // engine stays in library — no hook needed
migrator.Run()
```

> **Note**: `migrate_rollback.go` and `migrate_status.go` only import `database/migrations` for the engine. Since the engine stays in the library, these files need **zero changes**. App-specific migration files (the `init()` functions that call `migrations.Register()`) move to the starter and self-register via blank import in the starter's `main.go`.

#### C7 — `core/audit/audit.go` → Move `AuditLog` Into `core/audit/`

```go
// BEFORE (v1 — core/audit/audit.go):
import "github.com/RAiWorks/RapidGo/database/models"
record := models.AuditLog{...}           // hard import
var logs []models.AuditLog               // hard import

// AFTER (v2 — core/audit/audit.go):
// AuditLog struct defined HERE in core/audit/ package (framework-owned)
type AuditLog struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"index;not null;default:0" json:"user_id"`
    Action    string    `gorm:"size:50;not null;index" json:"action"`
    ModelType string    `gorm:"size:100;not null;index" json:"model_type"`
    ModelID   uint      `gorm:"not null;index" json:"model_id"`
    OldValues string    `gorm:"type:text" json:"old_values,omitempty"`
    NewValues string    `gorm:"type:text" json:"new_values,omitempty"`
    Metadata  string    `gorm:"type:text" json:"metadata,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}

record := AuditLog{...}    // no external import
var logs []AuditLog         // no external import
```

The starter's `database/models/registry.go` includes `&audit.AuditLog{}` if needed for AutoMigrate, or the framework's migrate command handles it automatically.

---

## 5. Starter `main.go` — The Single Wiring Point

After the split, the starter's entry point wires everything:

```go
package main

import (
    "github.com/RAiWorks/RapidGo/core/app"
    "github.com/RAiWorks/RapidGo/core/cli"
    "github.com/RAiWorks/RapidGo/core/container"
    "github.com/RAiWorks/RapidGo/core/router"
    "github.com/RAiWorks/RapidGo/core/scheduler"
    "github.com/RAiWorks/RapidGo/core/service"
    "github.com/RAiWorks/RapidGo/database/migrations"
    "gorm.io/gorm"

    "myapp/app/jobs"
    "myapp/app/providers"
    "myapp/app/schedule"
    "myapp/database/models"
    appmigrations "myapp/database/migrations"
    "myapp/database/seeders"
    "myapp/routes"
)

func main() {
    // 1. Bootstrap — register service providers
    cli.SetBootstrap(func(a *app.App, mode service.Mode) {
        a.Register(&providers.ConfigProvider{})
        a.Register(&providers.LoggerProvider{})
        a.Register(&providers.DatabaseProvider{})
        a.Register(&providers.RedisProvider{})
        a.Register(&providers.QueueProvider{})
        if mode.Has(service.ModeWeb) {
            a.Register(&providers.SessionProvider{})
        }
        a.Register(&providers.MiddlewareProvider{Mode: mode})
        a.Register(&providers.RouterProvider{Mode: mode})
    })

    // 2. Routes
    cli.SetRoutes(routes.Register)

    // 3. Background jobs
    cli.SetJobRegistrar(jobs.RegisterJobs)

    // 4. Scheduled tasks
    cli.SetScheduleRegistrar(schedule.RegisterSchedule)

    // 5. Database models (for AutoMigrate)
    cli.SetModelRegistry(models.All)

    // 6. Database seeders
    cli.SetSeeder(seeders.Run)

    // Blank-import app migrations so init() registers them
    _ = appmigrations.Migrations // or use blank import: _ "myapp/database/migrations"

    // Run CLI
    cli.Execute()
}
```

---

## 6. Branch Strategy

### 6.1 Branch Structure

```
main ─────────────────────────── v1.0.0 (frozen, monolithic, kept as-is)
  │
  └── v2 ────────────────────── default development branch for v2
       │
       ├── docs/v2-01-master-doc                (this document)
       │
       ├── feature/v2-01-hooks-foundation       Phase A: Add hooks.go + callback types
       ├── feature/v2-02-audit-decouple         Phase A: Move AuditLog → core/audit/
       ├── feature/v2-03-root-decouple          Phase B: root.go uses SetBootstrap
       ├── feature/v2-04-serve-decouple         Phase B: serve.go uses SetRoutes
       ├── feature/v2-05-worker-decouple        Phase B: work.go + schedule_run.go use hooks
       ├── feature/v2-06-migrate-decouple       Phase B: migrate/seed use hooks
       ├── feature/v2-07-remove-app-code        Phase C: Delete app/, routes/, etc.
       ├── feature/v2-08-starter-repo           Phase C: Create RapidGo-starter
       ├── feature/v2-09-rapidgo-new-cmd        Phase D: CLI scaffolder
       ├── feature/v2-10-library-readme         Phase D: Documentation
       │
       └── v2.0.0 tag when complete
```

### 6.2 Branch Rules

- Each feature branch creates from `v2`, merges back to `v2`
- Feature branches are deleted after merge (v2 follows clean branch strategy)
- `main` stays at v1.0.0 — no changes
- When v2 is complete, `v2` becomes the default branch on GitHub

---

## 7. Implementation Phases

### Phase A — Foundation (No Breaking Changes)

> These steps add new code alongside existing code. The monolith continues to work.

| Step | Branch | Work | Files Changed | Tests |
|------|--------|------|--------------|-------|
| A1 | `feature/v2-01-hooks-foundation` | Create `core/cli/hooks.go` with all 6 type definitions and `Set*()` functions. No existing code changes. | +1 new file | Unit test: each `Set*()` stores the function, default nil |
| A2 | `feature/v2-02-audit-decouple` | Move `AuditLog` struct from `database/models/audit_log.go` into `core/audit/model.go`. Update `core/audit/audit.go` to use local type. Keep `database/models/audit_log.go` as a type alias or re-export for backward compatibility. | `core/audit/audit.go`, `core/audit/audit_test.go`, +`core/audit/model.go`, `database/models/audit_log.go` | Existing audit tests pass |

### Phase B — Decouple (Break Hard Imports)

> These steps replace hard imports with callback calls. After each step, the app still works because `main.go` wires the callbacks.

| Step | Branch | Work | Files Changed | Tests |
|------|--------|------|--------------|-------|
| B1 | `feature/v2-03-root-decouple` | Replace `NewApp()` hard-coded providers with `bootstrapFn()`. Remove `app/providers` import. Update `cmd/main.go` to call `cli.SetBootstrap()`. | `core/cli/root.go`, `cmd/main.go` | Build + `go test ./...` pass |
| B2 | `feature/v2-04-serve-decouple` | Replace `routes.*` calls in `applyRoutesForMode()` with `routeRegistrar()`. Remove `routes` import. Update `cmd/main.go` to call `cli.SetRoutes()`. | `core/cli/serve.go`, `cmd/main.go` | Serve command starts correctly |
| B3 | `feature/v2-05-worker-decouple` | Replace manual bootstrap in `work.go` and `schedule_run.go` with `NewApp()` + hooks. Remove `app/jobs`, `app/providers`, `app/schedule` imports. Update `cmd/main.go`. | `core/cli/work.go`, `core/cli/schedule_run.go`, `cmd/main.go` | Build passes |
| B4 | `feature/v2-06-migrate-decouple` | Replace `models.All()` with `modelRegistryFn()` in `migrate.go`. Replace `seeders.*` calls with `seederFn()` in `seed.go`. Remove `database/models` import from `migrate.go`. Remove `database/seeders` import from `seed.go`. Note: `database/migrations` imports in `migrate*.go` stay — the engine is library code. Update `cmd/main.go`. Also refactor test files: `models_test.go`, `scopes_test.go`, `migrations_test.go`, `seeders_test.go` to use test-only model structs instead of app models. | `core/cli/migrate.go`, `core/cli/seed.go`, `cmd/main.go`, `database/models/models_test.go`, `database/models/scopes_test.go`, `database/migrations/migrations_test.go`, `database/seeders/seeders_test.go` | Migrate, seed commands work. All refactored test files pass. |

### Phase C — Split (Separate Repos)

> After Phase B, `core/` has zero imports from `app/`, `routes/`, `http/`, `plugins/`, `database/models/` (except `base.go`), or `database/seeders/`.

| Step | Branch | Work | Files Changed | Tests |
|------|--------|------|--------------|-------|
| C1 | `feature/v2-07-remove-app-code` | Delete from library: `app/`, `routes/`, `http/`, `plugins/`, `resources/`, `storage/`, `database/models/` (except `base.go`, `scopes.go`), `database/migrations/2026*` files (keep `migrator.go`), `database/seeders/` (except `seeder.go` engine), `.env`, `Dockerfile`, `docker-compose.yml`, `Caddyfile`, `Makefile`, `tests/`. Clean up `go.mod` (remove unused deps). Update `cmd/main.go` to minimal library CLI. | Many deletions | `go build ./...` and `go test ./...` pass on library alone |
| C2 | `feature/v2-08-starter-repo` | Create `RapidGo-starter` repo with all deleted code. New `go.mod` importing `github.com/RAiWorks/RapidGo`. New `cmd/main.go` with full wiring. `.env.example` (not `.env`). Verify builds and runs. | New repo | Starter builds and starts |

### Phase D — Polish

| Step | Branch | Work | Files Changed |
|------|--------|------|--------------|
| D1 | `feature/v2-09-rapidgo-new-cmd` | `rapidgo new myapp` command — downloads starter template, replaces module name, runs `go mod tidy` | +`core/cli/new.go` |
| D2 | `feature/v2-10-library-readme` | Library README with `go get` examples, package index, quick start. Starter README with clone + getting-started guide. | `README.md` (both repos) |

---

## 8. Verification Checkpoints

After each phase, these must pass:

| Checkpoint | Phase A | Phase B | Phase C | Phase D |
|-----------|:---:|:---:|:---:|:---:|
| `go build ./...` passes | Yes | Yes | Yes (library only) | Yes |
| `go test ./...` passes | Yes | Yes | Yes (library only) | Yes |
| `go vet ./...` passes | Yes | Yes | Yes | Yes |
| All 30 test packages pass | Yes | Yes | Library subset | Yes |
| `rapidgo serve` starts | Yes | Yes | N/A (cli binary removed) | Via starter |
| No `app/` imports in `core/` | No | Yes | Yes | Yes |
| No `routes/` imports in `core/` | No | Yes | Yes | Yes |
| No `database/models/` imports in `core/` (except base) | No | Yes | Yes | Yes |

> **Note**: `core/cli/seed.go` importing `database/seeders` is NOT a coupling point. The seeder engine (`RunAll()`, `RunByName()`, `Register()`) stays in the library. Only app-specific seeder implementations move to the starter.
| Starter builds and runs | N/A | N/A | Yes | Yes |

---

## 9. File Disposition Table

### Files That Stay in Library (RapidGo)

| File/Directory | Notes |
|---------------|-------|
| `core/` (all packages) | Framework internals — zero app-specific code after Phase B |
| `core/cli/hooks.go` | NEW — callback type definitions and setters |
| `core/audit/model.go` | NEW — `AuditLog` struct (moved from `database/models/`) |
| `database/connection.go` | Generic DB connection factory |
| `database/resolver.go` | Read/write split resolver |
| `database/transaction.go` | Transaction helpers |
| `database/transaction_example.go` | Example — could also move to starter |
| `database/migrations/migrator.go` | Migration engine (generic) |
| `database/models/base.go` | `BaseModel` with common fields |
| `database/models/scopes.go` | `WithTrashed()`, `OnlyTrashed()` — generic GORM scopes |
| `testing/testutil/` | Test utilities for user apps |
| `go.mod` | `module github.com/RAiWorks/RapidGo` |
| `go.sum` | Auto-generated by Go toolchain |
| `.gitignore` | Library-specific ignore rules |
| `LICENSE` | MIT |
| `README.md` | Library-focused (how to `go get`) |
| `database/database_test.go` | Generic DB connection tests |
| `database/resolver_test.go` | Generic resolver tests |
| `database/transaction_test.go` | Generic transaction tests |

### Files That Move to Starter (RapidGo-starter)

| File/Directory | Notes |
|---------------|-------|
| `cmd/main.go` | Rewritten with all `cli.Set*()` wiring |
| `app/providers/` | All 8 providers — app-specific |
| `app/helpers/` | App-specific utilities |
| `app/services/` | Business logic |
| `app/jobs/` | Background job handlers |
| `app/schedule/` | Scheduled task definitions |
| `app/plugins.go` | Plugin registration |
| `routes/web.go`, `api.go`, `ws.go` | App-specific route definitions |
| `http/controllers/` | App-specific request handlers |
| `http/requests/` | App-specific form requests |
| `http/responses/` | App-specific response formatting |
| `database/models/user.go` | App-specific model |
| `database/models/post.go` | App-specific model |
| `database/models/audit_log.go` | Re-export of `core/audit.AuditLog` (or removed) |
| `database/models/registry.go` | `All()` — app-specific model list |
| `database/migrations/2026*` | App-specific migration files |
| `database/migrations/migrations_test.go` | App-specific migration tests |
| `database/seeders/` | Seed data |
| `resources/views/` | HTML templates |
| `resources/lang/` | Translation files |
| `resources/static/` | CSS, JS, images |
| `storage/` | Runtime directories (uploads, cache, sessions, logs) |
| `plugins/example/` | Example plugin |
| `.env` → `.env.example` | Config template (not committed as `.env`) |
| `Dockerfile` | App-specific deployment |
| `docker-compose.yml` | App-specific deployment |
| `Caddyfile` | App-specific web server config |
| `Makefile` | App-specific build commands |
| `.gitignore` | Starter-specific ignore rules (separate from library) |
| `.dockerignore` | Docker build ignore rules |
| `app/providers/providers_test.go` | Provider tests |
| `app/helpers/helpers_test.go` | Helper function tests |
| `app/helpers/pagination_test.go` | Pagination tests |
| `app/services/user_service_test.go` | Service tests |
| `http/controllers/controllers_test.go` | Controller tests |
| `http/responses/response_test.go` | Response tests |
| `tests/integration/` | Integration test container |
| `tests/unit/` | Unit test container |

### Files That May Be Split or Duplicated

| File/Directory | Library | Starter | Notes |
|---------------|---------|---------|-------|
| `docs/` | Framework feature docs | Getting-started guide | Library gets `docs/framework/`, starter gets a focused README |
| `testing/` | `testutil/` stays | `tests/` integration tests | Starter's `tests/` has examples using library's `testutil` |
| `database/models/base.go` | Yes | Imported via framework | Starter models embed `models.BaseModel` from library |

### Library Test Files Requiring Refactoring

These test files stay in the library (they test library engines) but currently reference app-specific models (`User`, `Post`). They **must be refactored** during Phase B to use test-only model structs.

| Test File | Current App Deps | Required Fix |
|-----------|------------------|--------------|
| `database/models/models_test.go` | Uses `User{}`, `Post{}` directly | Split: `BaseModel` tests stay (use local `testModel` struct), `User`/`Post` tests move to starter |
| `database/models/scopes_test.go` | Uses `User{}`, `Post{}` in scope tests | Replace with `testModel struct { BaseModel; Name string }` |
| `database/migrations/migrations_test.go` | Imports `database/models` for `User`, `Post` | Replace with local test model struct |
| `database/seeders/seeders_test.go` | `setupTestDB()` uses `models.User` | Refactor `setupTestDB()` to use local test model; core engine tests already use `mockSeeder` (fine) |

---

## 10. Migration Registrar Design (Special Case)

The current migration system uses `init()` side effects — each migration file registers itself on package import:

```go
// database/migrations/20260307000001_create_jobs_tables.go (current)
func init() {
    Register(Migration{
        Version: "20260307000001_create_jobs_tables",
        Up:      func(db *gorm.DB) error { ... },
        Down:    func(db *gorm.DB) error { ... },
    })
}
```

After the split, migration files live in the starter. The `init()` approach still works because the starter's `go.mod` imports the library's `migrations` package:

```go
// Starter: database/migrations/20260307000001_create_jobs_tables.go
package migrations

import "github.com/RAiWorks/RapidGo/database/migrations"

func init() {
    migrations.Register(migrations.Migration{
        Version: "20260307000001_create_jobs_tables",
        Up:      func(db *gorm.DB) error { ... },
        Down:    func(db *gorm.DB) error { ... },
    })
}
```

**Decision**: The `init()` + blank import pattern is sufficient. No `SetMigrationRegistrar` hook is needed because:
1. The migration engine (`migrator.go` with `Register()`, `NewMigrator()`, `Run()`, `Rollback()`, `Status()`) stays in the library.
2. App migration files move to the starter and self-register by calling the library's `migrations.Register()` in their `init()` functions.
3. The starter's `main.go` does a blank import of its migrations package to trigger registration.

The same `init()` pattern already works for seeders via `seeders.Register()`.

---

## 11. The `database/models` Package After Split

### In the Library

```go
// database/models/base.go (unchanged)
package models

type BaseModel struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// database/models/scopes.go (unchanged)
func WithTrashed(db *gorm.DB) *gorm.DB { ... }
func OnlyTrashed(db *gorm.DB) *gorm.DB { ... }
```

### In the Starter

```go
// database/models/user.go
package models

import "github.com/RAiWorks/RapidGo/database/models"

type User struct {
    models.BaseModel  // embeds framework's BaseModel
    Name     string   `gorm:"size:255" json:"name"`
    Email    string   `gorm:"size:255;uniqueIndex" json:"email"`
    Password string   `gorm:"size:255" json:"-"`
    // ...
}

// database/models/registry.go
func All() []interface{} {
    return []interface{}{&User{}, &Post{}}
}
```

**Import conflict note**: The starter's `database/models` package has the same name as the library's. This is fine — Go handles it with import aliases:

```go
import (
    fwmodels "github.com/RAiWorks/RapidGo/database/models"
    "myapp/database/models"
)
```

---

## 12. Risk Register

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Package name collision (`database/models` in both repos) | Medium | Low | Import aliases work. Consider renaming library's to `database/base` if confusing. |
| `init()` migration registration breaks after split | Medium | Medium | Support both `init()` and explicit callback. Document clearly. |
| `AuditLog` moving to `core/audit/` breaks starter imports | High | Low | Starter updates import path. One-line change. |
| `work.go` / `schedule_run.go` have custom bootstrap (not using `NewApp`) | Known | Low | Refactor to use `NewApp()` in Phase B — `bootstrapFn` handles mode-specific logic. |
| Starter falls behind framework versions | Medium | Medium | CI in starter tests against `@latest` library. Semantic versioning. |
| `go.sum` churn when library adds dependencies | Low | Low | Normal Go module behavior. Unavoidable. |
| Some `core/cli/` tests depend on app-specific code | Medium | Medium | Audit during Phase B. Refactor tests to use stubs or mocks. |

---

## 13. Success Criteria

v2.0.0 is tagged when ALL of the following are true:

- [ ] `core/` has **zero** imports from `app/`, `routes/`, `http/`, `plugins/`
- [ ] `core/` has **zero** imports from `database/models/` except `base.go` and `scopes.go`
- [ ] `core/cli/hooks.go` defines all 6 hooks with `Set*()` functions
- [ ] Library `go build ./...` and `go test ./...` pass standalone (all test files refactored to use test-only model structs)
- [ ] Library has no `app/`, `routes/`, `http/`, `plugins/` directories
- [ ] Starter repo builds and runs with `go run cmd/main.go serve`
- [ ] Starter imports framework at `github.com/RAiWorks/RapidGo@v2.0.0`
- [ ] `rapidgo new myapp` scaffolds a working project
- [ ] Library README documents `go get` import path and package index
- [ ] Starter README documents getting started + every `Set*()` hook

---

## 14. Appendix: Current vs Target Import Graph

### v1.0.0 (Current) — Circular Dependencies

```
cmd/main.go → core/cli → app/providers → routes
                       → app/jobs
                       → app/schedule
                       → database/seeders
                       → database/models
                       → database/migrations
             core/audit → database/models
```

### v2.0.0 (Target) — Clean Dependency Flow

```
Library (github.com/RAiWorks/RapidGo):
  core/cli → core/*         (hooks.go provides callback types — no app imports)
  core/audit → core/audit   (AuditLog model is local)
  database/ → gorm          (connection, transaction, migrator — generic)

Starter (github.com/RAiWorks/RapidGo-starter):
  cmd/main.go → cli.Set*()  (wires app code into framework)
              → app/providers, routes, jobs, schedule, models, etc.
```

All coupling flows **one direction**: starter → library. Never library → starter.

---

> **Next Step**: Create `v2` branch from `main`, then start with `feature/v2-01-hooks-foundation`.
