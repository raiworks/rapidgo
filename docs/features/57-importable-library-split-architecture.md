# 🏗️ Architecture: Importable Library Split

> **Feature**: `57` — Importable Library Split
> **Discussion**: [`57-importable-library-split-discussion.md`](57-importable-library-split-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-08

---

## Overview

Introduce a 6-hook callback system in `core/cli/hooks.go` that decouples the library from application code. Move `AuditLog` into `core/audit/`. Replace all hard imports in `core/cli/` commands with hook callbacks. Wire hooks in `cmd/main.go`. Delete application code from the library repo and create `RapidGo-starter` as a separate repository. Add `rapidgo new` scaffolding command.

**Reference**: Full diagrams, dependency graphs, and data flow in [`v2-architecture.md`](../v2-architecture.md). Detailed per-step code in [`v2-phase-a-foundation.md`](../v2-phase-a-foundation.md), [`v2-phase-b-decouple.md`](../v2-phase-b-decouple.md), [`v2-phase-c-split.md`](../v2-phase-c-split.md), [`v2-phase-d-polish.md`](../v2-phase-d-polish.md).

## File Structure

### Phase A — New Files

```
core/cli/
├── hooks.go             # CREATE — 6 callback types + Set*() functions
└── hooks_test.go        # CREATE — Tests for all hooks

core/audit/
└── model.go             # CREATE — AuditLog struct (moved from database/models/)
```

### Phase A — Modified Files

```
core/audit/
└── audit.go             # MODIFY — Remove database/models import, use local AuditLog
└── audit_test.go        # MODIFY — Remove database/models import, use local AuditLog

database/models/
└── audit_log.go         # MODIFY — Replace struct with type alias: AuditLog = audit.AuditLog
```

### Phase B — Modified Files

```
core/cli/
├── root.go              # MODIFY — Remove app/providers import, use bootstrapFn
├── serve.go             # MODIFY — Remove routes import, use routeRegistrar
├── work.go              # MODIFY — Remove app/jobs + app/providers, use NewApp() + jobRegistrar
├── schedule_run.go      # MODIFY — Remove app/providers + app/schedule, use NewApp() + scheduleRegistrar
├── migrate.go           # MODIFY — Remove database/models import, use modelRegistryFn
└── seed.go              # MODIFY — Remove database/seeders import, use seederFn

cmd/
└── main.go              # MODIFY — Wire all 6 cli.Set*() hooks

database/models/
├── models_test.go       # MODIFY — Replace User/Post with test-only model struct
└── scopes_test.go       # MODIFY — Replace User/Post with test-only model struct

database/migrations/
└── migrations_test.go   # MODIFY — Replace models.User with test-only model struct

database/seeders/
└── seeders_test.go      # MODIFY — Replace models.User with test-only model struct
```

### Phase C — Deletions from Library

```
DELETE (full directories):
  app/                   # 27 files — providers, helpers, services, jobs, schedule
  routes/                # 3 files — web.go, api.go, ws.go
  http/                  # 8 files — controllers, requests, responses
  plugins/               # 1 file — example plugin
  resources/             # 4 files — views, lang, static
  storage/               # 4 dirs — cache, logs, sessions, uploads
  tests/                 # 2 dirs — integration, unit
  reference/             # 3 files — reference docs

DELETE (individual files from database/):
  database/models/user.go
  database/models/post.go
  database/models/audit_log.go         # type alias no longer needed
  database/models/registry.go          # All() moves to starter
  database/transaction_example.go      # example code → starter

DELETE (app-specific migrations):
  database/migrations/20260307000001_create_jobs_tables.go
  database/migrations/20260308000001_add_soft_deletes.go
  database/migrations/20260308000002_add_totp_fields.go
  database/migrations/20260308000003_create_audit_logs_table.go

DELETE (app-specific seeders):
  database/seeders/user_seeder.go

DELETE (root files):
  Dockerfile, docker-compose.yml, Caddyfile, Makefile
  .dockerignore, .env.example
```

### Phase C — Starter Repo (new `RapidGo-starter` repository)

```
RapidGo-starter/
├── cmd/main.go                        # Full hook wiring
├── app/
│   ├── helpers/                       # All helper files
│   ├── jobs/example_job.go
│   ├── providers/                     # All 8 providers
│   ├── schedule/schedule.go
│   ├── services/
│   └── plugins.go
├── routes/web.go, api.go, ws.go
├── http/controllers/, requests/, responses/
├── database/
│   ├── models/user.go, post.go, registry.go  # Embed fwmodels.BaseModel
│   ├── migrations/                     # init() + fwmigrations.Register()
│   └── seeders/                        # init() + fwseeders.Register()
├── resources/, storage/, plugins/, tests/
├── .env.example, Dockerfile, Caddyfile, Makefile
├── go.mod                              # module github.com/RAiWorks/RapidGo-starter
└── README.md
```

### Phase D — New/Modified Files

```
core/cli/
├── new.go               # CREATE — rapidgo new command
├── new_test.go          # CREATE — Tests for new command
└── root.go              # MODIFY — Add newCmd to init()

README.md                # REWRITE — Library README with package index
```

## Data Model

No new database tables. The `AuditLog` struct moves from `database/models/` to `core/audit/model.go` (identical schema, same GORM tags). A type alias in `database/models/audit_log.go` preserves backward compatibility during Phase A–B.

## Component Design

### Hook System (`core/cli/hooks.go`)

**Responsibility**: Define 6 callback types that decouple CLI commands from application code.
**Package**: `core/cli`
**File**: `core/cli/hooks.go`

```
Exported API:
├── SetBootstrap(fn BootstrapFunc)                  # Register providers
├── SetRoutes(fn RouteRegistrar)                   # Register routes
├── SetJobRegistrar(fn JobRegistrar)               # Register jobs
├── SetScheduleRegistrar(fn ScheduleRegistrar)     # Register schedule
├── SetModelRegistry(fn ModelRegistryFunc)          # Provide models
└── SetSeeder(fn SeederFunc)                        # Run seeders
```

**Type signatures**:

| Type | Signature |
|------|-----------|
| `BootstrapFunc` | `func(a *app.App, mode service.Mode)` |
| `RouteRegistrar` | `func(r *router.Router, c *container.Container, mode service.Mode)` |
| `JobRegistrar` | `func()` |
| `ScheduleRegistrar` | `func(s *scheduler.Scheduler, a *app.App)` |
| `ModelRegistryFunc` | `func() []interface{}` |
| `SeederFunc` | `func(db *gorm.DB, name string) error` |

### `rapidgo new` Command (`core/cli/new.go`)

**Responsibility**: Scaffold new projects from `RapidGo-starter` template.
**Package**: `core/cli`
**File**: `core/cli/new.go`

```
Exported API:
└── newCmd (*cobra.Command)        # "rapidgo new [project-name]"

Internal:
├── downloadStarter() (string, error)              # Download zip from GitHub
├── extractZip(zipPath, targetDir string) error    # Extract with zip-slip protection
└── replaceModuleName(projectDir string) error     # Replace module paths in .go + go.mod
```

## Data Flow

### Serve Command (after decoupling)

```
cmd/main.go
  │ cli.SetBootstrap(...) → stores bootstrapFn
  │ cli.SetRoutes(...)    → stores routeRegistrar
  │ cli.Execute()
  │
  ▼
core/cli/root.go → NewApp(mode)
  │ bootstrapFn(app, mode)     ← calls starter's bootstrap function
  │ app.Boot()
  │
  ▼
core/cli/serve.go → applyRoutesForMode(r, c, mode)
  │ [template/static setup stays in library]
  │ routeRegistrar(r, c, mode) ← calls starter's route registrar
  │ [health check stays in library]
  │
  ▼
core/server → ListenAndServe()
```

### Migrate Command (after decoupling)

```
cmd/main.go
  │ cli.SetModelRegistry(models.All)
  │ cli.Execute()
  │
  ▼
core/cli/migrate.go
  │ NewApp(ModeAll) → bootstrapFn(app, ModeAll) → Boot()
  │ db := container.MustMake[*gorm.DB]("db")
  │
  │ modelRegistryFn() → returns []interface{}{&User{}, &Post{}, ...}
  │ db.AutoMigrate(models...)
  │
  │ migrations.NewMigrator(db) → runs registered file migrations
  ▼
```

## Configuration

No new environment variables.

## Security Considerations

- **`rapidgo new` zip extraction**: Zip slip protection in `extractZip()` — all paths validated against target directory using `filepath.Clean()` + prefix check
- **Project name validation**: `new` command rejects names with path traversal characters (`/\:*?"<>|`)
- **No secrets in library**: `.env` and `.env.example` move to starter repo

## Trade-offs & Alternatives

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| **Callback hooks** (`Set*()` functions) | Simple, explicit, zero reflection, easy to understand | Global mutable state (package-level vars) | ✅ Selected |
| Interface-based DI | More testable, no global state | Over-engineered for 6 hooks, requires interface definitions in library | ❌ Too complex |
| Config-based plugin loading (reflection) | Most flexible, runtime discovery | Magic, hard to debug, non-idiomatic Go | ❌ Anti-pattern in Go |
| Separate binary per command | No coupling at all | Bad DX, multiple binaries to install | ❌ Poor UX |

## Next

Create tasks doc → `57-importable-library-split-tasks.md`
