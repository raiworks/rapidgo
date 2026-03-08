# RapidGo v2 — Phase B: Decouple

> **Phase**: B — Break Hard Imports  
> **Steps**: B1 (root), B2 (serve), B3 (worker/scheduler), B4 (migrate/seed + test refactoring)  
> **Branches**: `feature/v2-03-root-decouple` through `feature/v2-06-migrate-decouple`  
> **Pre-requisite**: Phase A complete (hooks.go exists, AuditLog moved)  
> **Post-condition**: Zero `app/`, `routes/`, `http/`, `plugins/` imports in `core/`. Monolith still works.  

---

## Step B1: Decouple `root.go`

### Branch

`feature/v2-03-root-decouple` (from `v2`)

### Objective

Replace the hard-coded provider registrations in `NewApp()` with the `bootstrapFn` callback. Remove `app/providers` import from `core/cli/root.go`. Wire `cli.SetBootstrap()` in `cmd/main.go`.

### Files Changed

| Action | File | What Changes |
|--------|------|-------------|
| MODIFY | `core/cli/root.go` | Remove `app/providers` import, use `bootstrapFn` in `NewApp()` |
| MODIFY | `cmd/main.go` | Add `cli.SetBootstrap()` call + all provider imports |

### Modification: `core/cli/root.go`

**Before:**
```go
import (
	"fmt"
	"os"

	"github.com/RAiWorks/RapidGo/app/providers"    // ← REMOVE
	"github.com/RAiWorks/RapidGo/core/app"
	"github.com/RAiWorks/RapidGo/core/service"
	"github.com/spf13/cobra"
)

func NewApp(mode service.Mode) *app.App {
	application := app.New()

	application.Register(&providers.ConfigProvider{})
	application.Register(&providers.LoggerProvider{})
	if mode.Has(service.ModeWeb) || mode.Has(service.ModeAPI) || mode.Has(service.ModeWS) {
		application.Register(&providers.DatabaseProvider{})
	}
	application.Register(&providers.RedisProvider{})
	application.Register(&providers.QueueProvider{})
	if mode.Has(service.ModeWeb) {
		application.Register(&providers.SessionProvider{})
	}
	application.Register(&providers.MiddlewareProvider{Mode: mode})
	application.Register(&providers.RouterProvider{Mode: mode})

	application.Boot()
	return application
}
```

**After:**
```go
import (
	"fmt"
	"os"

	"github.com/RAiWorks/RapidGo/core/app"
	"github.com/RAiWorks/RapidGo/core/service"
	"github.com/spf13/cobra"
)

func NewApp(mode service.Mode) *app.App {
	application := app.New()

	if bootstrapFn != nil {
		bootstrapFn(application, mode)
	}

	application.Boot()
	return application
}
```

### Modification: `cmd/main.go`

**Before:**
```go
package main

import "github.com/RAiWorks/RapidGo/core/cli"

func main() {
	cli.Execute()
}
```

**After:**
```go
package main

import (
	"github.com/RAiWorks/RapidGo/app/providers"
	"github.com/RAiWorks/RapidGo/core/app"
	"github.com/RAiWorks/RapidGo/core/cli"
	"github.com/RAiWorks/RapidGo/core/service"
)

func main() {
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

	cli.Execute()
}
```

### Verification

```bash
go build ./...
go test ./core/cli/ -v
grep "app/providers" core/cli/root.go     # should return nothing
go test ./... -count=1
```

---

## Step B2: Decouple `serve.go`

### Branch

`feature/v2-04-serve-decouple` (from `v2`, after B1 merges)

### Objective

Replace `routes.RegisterWeb/API/WS()` calls in `applyRoutesForMode()` with the `routeRegistrar` callback. Remove `routes` import from `core/cli/serve.go`. Wire `cli.SetRoutes()` in `cmd/main.go`.

### Files Changed

| Action | File | What Changes |
|--------|------|-------------|
| MODIFY | `core/cli/serve.go` | Remove `routes` import, use `routeRegistrar` in `applyRoutesForMode()` |
| MODIFY | `cmd/main.go` | Add `cli.SetRoutes()` call + route imports |

### Modification: `core/cli/serve.go`

**Imports — remove `routes`:**
```go
// BEFORE:
import (
	...
	"github.com/RAiWorks/RapidGo/routes"    // ← REMOVE
	...
)

// AFTER:
import (
	...
	// "routes" import removed
	...
)
```

**`applyRoutesForMode()` — replace route calls:**
```go
// BEFORE:
func applyRoutesForMode(r *router.Router, c *container.Container, m service.Mode) {
	if m.Has(service.ModeWeb) {
		r.SetFuncMap(router.DefaultFuncMap())
		viewsDir := filepath.Join("resources", "views")
		if info, err := os.Stat(viewsDir); err == nil && info.IsDir() {
			r.LoadTemplates(viewsDir)
		}
		if info, err := os.Stat("resources/static"); err == nil && info.IsDir() {
			r.Static("/static", "./resources/static")
		}
		if info, err := os.Stat("storage/uploads"); err == nil && info.IsDir() {
			r.Static("/uploads", "./storage/uploads")
		}
		routes.RegisterWeb(r)
	}
	if m.Has(service.ModeAPI) {
		routes.RegisterAPI(r)
	}
	if m.Has(service.ModeWS) {
		routes.RegisterWS(r)
	}
	// health check code stays...
}

// AFTER:
func applyRoutesForMode(r *router.Router, c *container.Container, m service.Mode) {
	if m.Has(service.ModeWeb) {
		r.SetFuncMap(router.DefaultFuncMap())
		viewsDir := filepath.Join("resources", "views")
		if info, err := os.Stat(viewsDir); err == nil && info.IsDir() {
			r.LoadTemplates(viewsDir)
		}
		if info, err := os.Stat("resources/static"); err == nil && info.IsDir() {
			r.Static("/static", "./resources/static")
		}
		if info, err := os.Stat("storage/uploads"); err == nil && info.IsDir() {
			r.Static("/uploads", "./storage/uploads")
		}
	}

	// Delegate route registration to the application callback
	if routeRegistrar != nil {
		routeRegistrar(r, c, m)
	}

	// Health check — each per-service router gets its own health endpoints
	if c.Has("db") {
		health.Routes(r, func() *gorm.DB {
			return container.MustMake[*gorm.DB](c, "db")
		})
	}
}
```

**Key design decision**: Static file/template setup stays in the library (it's framework configuration logic). Only the `routes.Register*()` calls move to the callback.

### Modification: `cmd/main.go`

Add to existing `main()` (after `SetBootstrap`):

```go
import (
	...
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/router"
	"github.com/RAiWorks/RapidGo/routes"
)

// In main():
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
```

### Verification

```bash
go build ./...
grep "routes" core/cli/serve.go   # only "routeRegistrar" — no "routes." package ref
go test ./... -count=1
# Manual: go run cmd/main.go serve  → verify routes respond
```

---

## Step B3: Decouple `work.go` and `schedule_run.go`

### Branch

`feature/v2-05-worker-decouple` (from `v2`, after B1 merges — independent of B2)

### Objective

Replace the manual provider bootstrap and hard imports in `work.go` and `schedule_run.go` with `NewApp()` (which now uses `bootstrapFn`) + the `jobRegistrar` and `scheduleRegistrar` callbacks.

### Files Changed

| Action | File | What Changes |
|--------|------|-------------|
| MODIFY | `core/cli/work.go` | Remove `app/jobs`, `app/providers` imports, use `NewApp()` + `jobRegistrar` |
| MODIFY | `core/cli/schedule_run.go` | Remove `app/providers`, `app/schedule` imports, use `NewApp()` + `scheduleRegistrar` |
| MODIFY | `cmd/main.go` | Add `cli.SetJobRegistrar()` + `cli.SetScheduleRegistrar()` |

### Modification: `core/cli/work.go`

**Before:**
```go
import (
	...
	"github.com/RAiWorks/RapidGo/app/jobs"        // ← REMOVE
	"github.com/RAiWorks/RapidGo/app/providers"    // ← REMOVE
	"github.com/RAiWorks/RapidGo/core/app"         // ← REMOVE (use NewApp instead)
	"github.com/RAiWorks/RapidGo/core/config"
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/queue"
	...
)

// Inside RunE:
config.Load()

// Minimal bootstrap — no HTTP providers needed.
application := app.New()
application.Register(&providers.ConfigProvider{})
application.Register(&providers.LoggerProvider{})
application.Register(&providers.DatabaseProvider{})
application.Register(&providers.RedisProvider{})
application.Register(&providers.QueueProvider{})
application.Boot()

// Register application job handlers.
jobs.RegisterJobs()
```

**After:**
```go
import (
	...
	"github.com/RAiWorks/RapidGo/core/config"
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/queue"
	"github.com/RAiWorks/RapidGo/core/service"
	...
)

// Inside RunE:
config.Load()

application := NewApp(service.ModeAll)

// Register application job handlers via callback.
if jobRegistrar != nil {
	jobRegistrar()
}
```

### Modification: `core/cli/schedule_run.go`

**Before:**
```go
import (
	...
	"github.com/RAiWorks/RapidGo/app/providers"    // ← REMOVE
	"github.com/RAiWorks/RapidGo/app/schedule"      // ← REMOVE
	"github.com/RAiWorks/RapidGo/core/app"           // ← REMOVE
	"github.com/RAiWorks/RapidGo/core/config"
	"github.com/RAiWorks/RapidGo/core/scheduler"
	...
)

// Inside RunE:
config.Load()

application := app.New()
application.Register(&providers.ConfigProvider{})
application.Register(&providers.LoggerProvider{})
application.Register(&providers.DatabaseProvider{})
application.Register(&providers.RedisProvider{})
application.Register(&providers.QueueProvider{})
application.Boot()

s := scheduler.New(slog.Default())
schedule.RegisterSchedule(s, application)
```

**After:**
```go
import (
	...
	"github.com/RAiWorks/RapidGo/core/config"
	"github.com/RAiWorks/RapidGo/core/scheduler"
	"github.com/RAiWorks/RapidGo/core/service"
	...
)

// Inside RunE:
config.Load()

application := NewApp(service.ModeAll)

s := scheduler.New(slog.Default())

// Register scheduled tasks via callback.
if scheduleRegistrar != nil {
	scheduleRegistrar(s, application)
}
```

### Modification: `cmd/main.go`

Add to existing `main()`:

```go
import (
	...
	"github.com/RAiWorks/RapidGo/app/jobs"
	"github.com/RAiWorks/RapidGo/app/schedule"
)

// In main():
cli.SetJobRegistrar(jobs.RegisterJobs)
cli.SetScheduleRegistrar(schedule.RegisterSchedule)
```

### Verification

```bash
go build ./...
grep "app/jobs\|app/providers\|app/schedule" core/cli/work.go core/cli/schedule_run.go  # empty
go test ./... -count=1
```

---

## Step B4: Decouple `migrate.go` + `seed.go` + Test Refactoring

### Branch

`feature/v2-06-migrate-decouple` (from `v2`, after B1 merges — independent of B2/B3)

### Objective

1. Replace `models.All()` in `migrate.go` with `modelRegistryFn()`. Remove `database/models` import.
2. Replace `seeders.RunByName/RunAll()` in `seed.go` with `seederFn()`. Remove `database/seeders` import.
3. Refactor 4 test files to use test-only model structs instead of app models (`User`, `Post`).

**Note**: `database/migrations` imports in `migrate.go`, `migrate_rollback.go`, `migrate_status.go` stay as-is — the migration engine is library code.

### Files Changed

| Action | File | What Changes |
|--------|------|-------------|
| MODIFY | `core/cli/migrate.go` | Remove `database/models` import, use `modelRegistryFn()` |
| MODIFY | `core/cli/seed.go` | Remove `database/seeders` import, use `seederFn()` |
| MODIFY | `cmd/main.go` | Add `cli.SetModelRegistry()` + `cli.SetSeeder()` |
| MODIFY | `database/models/models_test.go` | Replace `User`/`Post` with test-only model struct |
| MODIFY | `database/models/scopes_test.go` | Replace `User`/`Post` with test-only model struct |
| MODIFY | `database/migrations/migrations_test.go` | Replace `models.User` with test-only model struct |
| MODIFY | `database/seeders/seeders_test.go` | Replace `models.User` in `setupTestDB()` with test-only model |

### Modification: `core/cli/migrate.go`

**Before:**
```go
import (
	"fmt"

	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/service"
	"github.com/RAiWorks/RapidGo/database/migrations"
	"github.com/RAiWorks/RapidGo/database/models"       // ← REMOVE
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// Inside RunE:
if err := db.AutoMigrate(models.All()...); err != nil {
```

**After:**
```go
import (
	"fmt"

	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/service"
	"github.com/RAiWorks/RapidGo/database/migrations"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// Inside RunE:
if modelRegistryFn != nil {
	if err := db.AutoMigrate(modelRegistryFn()...); err != nil {
		return fmt.Errorf("auto-migrate failed: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "AutoMigrate complete.")
}
```

### Modification: `core/cli/seed.go`

**Before:**
```go
import (
	"fmt"

	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/service"
	"github.com/RAiWorks/RapidGo/database/seeders"    // ← REMOVE
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// Inside RunE:
name, _ := cmd.Flags().GetString("seeder")
if name != "" {
	if err := seeders.RunByName(db, name); err != nil {
		return err
	}
	...
}
if err := seeders.RunAll(db); err != nil {
```

**After:**
```go
import (
	"fmt"

	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/service"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// Inside RunE:
if seederFn == nil {
	return fmt.Errorf("no seeder registered — call cli.SetSeeder() in main.go")
}

name, _ := cmd.Flags().GetString("seeder")
if err := seederFn(db, name); err != nil {
	return err
}
if name != "" {
	fmt.Fprintf(cmd.OutOrStdout(), "Seeder %s complete.\n", name)
} else {
	fmt.Fprintln(cmd.OutOrStdout(), "Database seeding complete.")
}
```

### Modification: `cmd/main.go`

Add to existing `main()`:

```go
import (
	...
	"github.com/RAiWorks/RapidGo/database/models"
	"github.com/RAiWorks/RapidGo/database/seeders"
	"gorm.io/gorm"
)

// In main():
cli.SetModelRegistry(models.All)

cli.SetSeeder(func(db *gorm.DB, name string) error {
	if name != "" {
		return seeders.RunByName(db, name)
	}
	return seeders.RunAll(db)
})
```

### Test File Refactoring

#### `database/models/models_test.go`

Replace references to `User{}` and `Post{}` with a test-only model:

```go
// Add at top of test file:
type testModel struct {
	BaseModel
	Name  string `gorm:"size:255"`
	Email string `gorm:"size:255;uniqueIndex"`
}

// Replace all User{} → testModel{}
// Replace all Post{} → testModel{} (or create testPost if needed)
// Remove any references to User/Post specific fields not in testModel
```

#### `database/models/scopes_test.go`

Replace `User{}` / `Post{}` with test-only struct:

```go
type testScopesModel struct {
	BaseModel
	Name string `gorm:"size:255"`
}

// In setupTestDB:
// BEFORE: db.AutoMigrate(&User{})
// AFTER:  db.AutoMigrate(&testScopesModel{})

// In tests:
// BEFORE: db.Create(&User{...})
// AFTER:  db.Create(&testScopesModel{Name: "test"})
```

#### `database/migrations/migrations_test.go`

Replace `models.User` / `models.Post` with local test model:

```go
// Remove: "github.com/RAiWorks/RapidGo/database/models"

// Add:
type testMigrationModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}

// Replace: db.AutoMigrate(&models.User{})
// With:    db.AutoMigrate(&testMigrationModel{})
```

#### `database/seeders/seeders_test.go`

Replace `models.User` in `setupTestDB()`:

```go
// Remove: "github.com/RAiWorks/RapidGo/database/models"

// Add:
type testSeederModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}

// In setupTestDB():
// BEFORE: db.AutoMigrate(&models.User{})
// AFTER:  db.AutoMigrate(&testSeederModel{})
```

### Verification

```bash
go build ./...
grep "database/models" core/cli/migrate.go     # should NOT appear
grep "database/seeders" core/cli/seed.go        # should NOT appear
grep "database/models" core/cli/seed.go         # should NOT appear
go test ./database/models/ -v                   # refactored tests pass
go test ./database/migrations/ -v               # refactored tests pass
go test ./database/seeders/ -v                  # refactored tests pass
go test ./... -count=1                          # full suite passes
```

---

## Phase B Summary

### Import Coupling — Before vs After Phase B

| File | Before (v1) | After (Phase B) |
|------|-------------|----------------|
| `core/cli/root.go` | `app/providers` | **Clean** — uses `bootstrapFn` |
| `core/cli/serve.go` | `routes` | **Clean** — uses `routeRegistrar` |
| `core/cli/work.go` | `app/jobs`, `app/providers` | **Clean** — uses `NewApp()` + `jobRegistrar` |
| `core/cli/schedule_run.go` | `app/providers`, `app/schedule` | **Clean** — uses `NewApp()` + `scheduleRegistrar` |
| `core/cli/migrate.go` | `database/models`, `database/migrations` | **Clean** — `models` removed, `migrations` stays (engine) |
| `core/cli/seed.go` | `database/seeders` | **Clean** — uses `seederFn` |
| `core/cli/migrate_rollback.go` | `database/migrations` | **Unchanged** — engine stays in library |
| `core/cli/migrate_status.go` | `database/migrations` | **Unchanged** — engine stays in library |
| `core/audit/audit.go` | `database/models` | **Clean** (resolved in Phase A2) |

### Test Files — Refactored

| File | Before | After |
|------|--------|-------|
| `database/models/models_test.go` | Uses `User`, `Post` | Uses `testModel` |
| `database/models/scopes_test.go` | Uses `User`, `Post` | Uses `testScopesModel` |
| `database/migrations/migrations_test.go` | Imports `database/models` | Uses `testMigrationModel` |
| `database/seeders/seeders_test.go` | Imports `database/models` | Uses `testSeederModel` |

### `cmd/main.go` — Final State After Phase B

```go
package main

import (
	"github.com/RAiWorks/RapidGo/app/jobs"
	"github.com/RAiWorks/RapidGo/app/providers"
	"github.com/RAiWorks/RapidGo/app/schedule"
	"github.com/RAiWorks/RapidGo/core/app"
	"github.com/RAiWorks/RapidGo/core/cli"
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/router"
	"github.com/RAiWorks/RapidGo/core/service"
	"github.com/RAiWorks/RapidGo/database/models"
	"github.com/RAiWorks/RapidGo/database/seeders"
	"github.com/RAiWorks/RapidGo/routes"
	"gorm.io/gorm"

	_ "github.com/RAiWorks/RapidGo/database/migrations" // blank import for init() registration
)

func main() {
	// 1. Bootstrap
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

### Zero-Coupling Verification

After Phase B, run this one command to verify:

```bash
grep -rn "RAiWorks/RapidGo/app\|RAiWorks/RapidGo/routes\|RAiWorks/RapidGo/http\|RAiWorks/RapidGo/plugins" core/
```

This must return **zero results**.
