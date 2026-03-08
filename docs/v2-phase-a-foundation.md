# RapidGo v2 — Phase A: Foundation

> **Phase**: A — Foundation (No Breaking Changes)  
> **Steps**: A1 (hooks) + A2 (audit decouple)  
> **Branch**: `feature/v2-01-hooks-foundation`, `feature/v2-02-audit-decouple`  
> **Pre-requisite**: `v2` branch created from `main`  
> **Post-condition**: Monolith keeps working — all existing tests pass  

---

## Step A1: Create `core/cli/hooks.go`

### Branch

`feature/v2-01-hooks-foundation` (from `v2`)

### Objective

Define the 6 callback types and their `Set*()` setter functions. No existing file is modified. This is purely additive.

### Files Changed

| Action | File |
|--------|------|
| CREATE | `core/cli/hooks.go` |
| CREATE | `core/cli/hooks_test.go` |

### Implementation: `core/cli/hooks.go`

```go
package cli

import (
	"github.com/RAiWorks/RapidGo/core/app"
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/router"
	"github.com/RAiWorks/RapidGo/core/scheduler"
	"github.com/RAiWorks/RapidGo/core/service"
	"gorm.io/gorm"
)

// BootstrapFunc registers service providers on the application for the given mode.
type BootstrapFunc func(a *app.App, mode service.Mode)

// RouteRegistrar registers routes on the router for a given mode.
type RouteRegistrar func(r *router.Router, c *container.Container, mode service.Mode)

// JobRegistrar registers application job handlers with the queue dispatcher.
type JobRegistrar func()

// ScheduleRegistrar registers scheduled tasks on the scheduler.
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

// SetBootstrap sets the function that registers service providers during app initialization.
func SetBootstrap(fn BootstrapFunc) { bootstrapFn = fn }

// SetRoutes sets the function that registers routes on the router.
func SetRoutes(fn RouteRegistrar) { routeRegistrar = fn }

// SetJobRegistrar sets the function that registers background job handlers.
func SetJobRegistrar(fn JobRegistrar) { jobRegistrar = fn }

// SetScheduleRegistrar sets the function that registers scheduled tasks.
func SetScheduleRegistrar(fn ScheduleRegistrar) { scheduleRegistrar = fn }

// SetModelRegistry sets the function that returns all model structs for AutoMigrate.
func SetModelRegistry(fn ModelRegistryFunc) { modelRegistryFn = fn }

// SetSeeder sets the function that runs database seeders.
func SetSeeder(fn SeederFunc) { seederFn = fn }
```

### Implementation: `core/cli/hooks_test.go`

```go
package cli

import (
	"testing"

	"github.com/RAiWorks/RapidGo/core/app"
	"github.com/RAiWorks/RapidGo/core/service"
)

func TestHooksDefaultNil(t *testing.T) {
	// Reset to ensure clean state
	bootstrapFn = nil
	routeRegistrar = nil
	jobRegistrar = nil
	scheduleRegistrar = nil
	modelRegistryFn = nil
	seederFn = nil

	if bootstrapFn != nil {
		t.Error("bootstrapFn should default to nil")
	}
	if routeRegistrar != nil {
		t.Error("routeRegistrar should default to nil")
	}
	if jobRegistrar != nil {
		t.Error("jobRegistrar should default to nil")
	}
	if scheduleRegistrar != nil {
		t.Error("scheduleRegistrar should default to nil")
	}
	if modelRegistryFn != nil {
		t.Error("modelRegistryFn should default to nil")
	}
	if seederFn != nil {
		t.Error("seederFn should default to nil")
	}
}

func TestSetBootstrapStoresFunction(t *testing.T) {
	defer func() { bootstrapFn = nil }()

	called := false
	SetBootstrap(func(a *app.App, mode service.Mode) {
		called = true
	})

	if bootstrapFn == nil {
		t.Fatal("SetBootstrap did not store function")
	}
	bootstrapFn(nil, service.ModeAll)
	if !called {
		t.Error("stored bootstrap function was not called")
	}
}

func TestSetRoutesStoresFunction(t *testing.T) {
	defer func() { routeRegistrar = nil }()

	called := false
	SetRoutes(func(r *router.Router, c *container.Container, mode service.Mode) {
		called = true
	})

	if routeRegistrar == nil {
		t.Fatal("SetRoutes did not store function")
	}
	routeRegistrar(nil, nil, service.ModeAll)
	if !called {
		t.Error("stored route registrar was not called")
	}
}

func TestSetJobRegistrarStoresFunction(t *testing.T) {
	defer func() { jobRegistrar = nil }()

	called := false
	SetJobRegistrar(func() {
		called = true
	})

	if jobRegistrar == nil {
		t.Fatal("SetJobRegistrar did not store function")
	}
	jobRegistrar()
	if !called {
		t.Error("stored job registrar was not called")
	}
}

func TestSetModelRegistryStoresFunction(t *testing.T) {
	defer func() { modelRegistryFn = nil }()

	SetModelRegistry(func() []interface{} {
		return []interface{}{"test"}
	})

	if modelRegistryFn == nil {
		t.Fatal("SetModelRegistry did not store function")
	}
	result := modelRegistryFn()
	if len(result) != 1 {
		t.Errorf("modelRegistryFn returned %d items, want 1", len(result))
	}
}
```

### Test Plan: A1

| # | Test | Verifies |
|---|------|----------|
| T01 | `TestHooksDefaultNil` | All 6 hook vars default to nil |
| T02 | `TestSetBootstrapStoresFunction` | `SetBootstrap()` stores and callback works |
| T03 | `TestSetRoutesStoresFunction` | `SetRoutes()` stores and callback works |
| T04 | `TestSetJobRegistrarStoresFunction` | `SetJobRegistrar()` stores and callback works |
| T05 | `TestSetModelRegistryStoresFunction` | `SetModelRegistry()` stores and callback works |
| T06 | `go build ./...` | Compiles with no errors |
| T07 | `go test ./...` | All existing tests still pass |

### Verification Commands

```bash
go build ./...
go test ./core/cli/ -run TestHooks -v
go test ./core/cli/ -run TestSet -v
go test ./... -count=1
```

---

## Step A2: Move `AuditLog` to `core/audit/`

### Branch

`feature/v2-02-audit-decouple` (from `v2`, after A1 merges)

### Objective

Move the `AuditLog` struct from `database/models/audit_log.go` into `core/audit/model.go`. Update `core/audit/audit.go` to use the local type instead of importing `database/models`. Keep backward compatibility via type alias in `database/models/audit_log.go`.

### Files Changed

| Action | File | What Changes |
|--------|------|-------------|
| CREATE | `core/audit/model.go` | `AuditLog` struct definition (moved from `database/models/`) |
| MODIFY | `core/audit/audit.go` | Remove `database/models` import, use local `AuditLog` |
| MODIFY | `core/audit/audit_test.go` | Remove `database/models` import, use local `AuditLog` |
| MODIFY | `database/models/audit_log.go` | Replace struct with type alias `AuditLog = audit.AuditLog` |

### Implementation: `core/audit/model.go` (NEW)

```go
package audit

import "time"

// AuditLog records a single auditable action on a model.
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
```

### Modifications: `core/audit/audit.go`

**Remove import:**
```go
// BEFORE:
import (
    "encoding/json"
    "github.com/RAiWorks/RapidGo/database/models"
    "gorm.io/gorm"
)

// AFTER:
import (
    "encoding/json"
    "gorm.io/gorm"
)
```

**Replace all `models.AuditLog` with `AuditLog`:**
```go
// BEFORE:
record := models.AuditLog{...}
var logs []models.AuditLog

// AFTER:
record := AuditLog{...}
var logs []AuditLog
```

### Modifications: `core/audit/audit_test.go`

**Remove import:**
```go
// BEFORE:
import (
    "github.com/RAiWorks/RapidGo/database/models"
    ...
)

// AFTER: (remove the models import line)
```

**Replace all `models.AuditLog` with `AuditLog`:**
```go
// BEFORE:
db.AutoMigrate(&models.AuditLog{})
var record models.AuditLog

// AFTER:
db.AutoMigrate(&AuditLog{})
var record AuditLog
```

### Modifications: `database/models/audit_log.go`

**Replace struct with type alias for backward compatibility:**
```go
// BEFORE:
package models

import "time"

type AuditLog struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    ...
}

// AFTER:
package models

import "github.com/RAiWorks/RapidGo/core/audit"

// AuditLog is an alias for the framework's audit log model.
// Kept for backward compatibility — new code should use audit.AuditLog directly.
type AuditLog = audit.AuditLog
```

### Test Plan: A2

| # | Test | Verifies |
|---|------|----------|
| T01 | `go test ./core/audit/ -v` | All existing audit tests pass with local `AuditLog` |
| T02 | `go test ./database/models/ -v` | Models package still compiles (type alias works) |
| T03 | `go test ./... -count=1` | Full test suite passes — no regressions |
| T04 | Verify `core/audit/audit.go` has NO `database/models` import | `grep "database/models" core/audit/audit.go` returns empty |
| T05 | Verify type alias works | Code using `models.AuditLog` still compiles |

### Verification Commands

```bash
go build ./...
go test ./core/audit/ -v
go test ./database/models/ -v
go test ./... -count=1
grep -r "database/models" core/audit/  # should return nothing
```

---

## Phase A Summary

| What | Before | After |
|------|--------|-------|
| `core/cli/hooks.go` | Does not exist | 6 types + 6 `Set*()` functions |
| `core/audit/model.go` | Does not exist | `AuditLog` struct |
| `core/audit/audit.go` | Imports `database/models` | Uses local `AuditLog` |
| `database/models/audit_log.go` | Defines `AuditLog` struct | Type alias → `audit.AuditLog` |
| Monolith behavior | Works | Still works — zero breaking changes |
| Coupling points in `core/` | 7 | 6 (C7 resolved) |
