# RapidGo v2 — Phase C: Split

> **Phase**: C — Separate Library and Starter Repos  
> **Steps**: C1 (remove app code from library) + C2 (create starter repo)  
> **Branches**: `feature/v2-07-remove-app-code`, `feature/v2-08-starter-repo`  
> **Pre-requisite**: Phase B complete (zero `app/` imports in `core/`)  
> **Post-condition**: Two separate repos, both build and test independently  

---

## Step C1: Remove App Code from Library

### Branch

`feature/v2-07-remove-app-code` (from `v2`)

### Objective

Delete all application-specific files from the library repo. After this step, `go build ./...` and `go test ./...` pass on the library alone with no references to `app/`, `routes/`, `http/`, or `plugins/`.

### Files to DELETE from Library

**Full directories (recursive delete):**

| Directory | Contents | Files |
|-----------|----------|:-----:|
| `app/` | Providers, helpers, services, jobs, schedule, plugins.go | 27 |
| `routes/` | web.go, api.go, ws.go | 3 |
| `http/` | Controllers, requests, responses | 8 |
| `plugins/` | example/ | 1 |
| `resources/` | views/, lang/, static/ | 4 |
| `storage/` | cache/, logs/, sessions/, uploads/ | 4 |
| `tests/` | integration/, unit/ | 2 |
| `reference/` | Reference documentation | 3 |

**Individual files from `database/`:**

| File | Why Remove |
|------|-----------|
| `database/models/user.go` | App-specific model |
| `database/models/post.go` | App-specific model |
| `database/models/audit_log.go` | Type alias (original struct now in `core/audit/model.go`) |
| `database/models/registry.go` | App-specific `All()` function |
| `database/models/models_test.go` | Tests for app models (after B4, tests BaseModel only — keep if generic, delete if still app-specific) |
| `database/models/scopes_test.go` | Stays if refactored to use test-only model in B4. Verify before deleting. |
| `database/models/.gitkeep` | Unnecessary (base.go/scopes.go exist) |
| `database/migrations/20260307000001_create_jobs_tables.go` | App-specific migration |
| `database/migrations/20260308000001_add_soft_deletes.go` | App-specific migration |
| `database/migrations/20260308000002_add_totp_fields.go` | App-specific migration |
| `database/migrations/20260308000003_create_audit_logs_table.go` | App-specific migration |
| `database/migrations/migrations_test.go` | Stays if refactored in B4. Verify. |
| `database/migrations/.gitkeep` | Unnecessary (migrator.go exists) |
| `database/seeders/user_seeder.go` | App-specific seeder |
| `database/seeders/seeders_test.go` | Stays if refactored in B4. Verify. |
| `database/seeders/.gitkeep` | Unnecessary (seeder.go exists) |
| `database/transaction_example.go` | Example code — move to starter |

**Root files:**

| File | Why Remove |
|------|-----------|
| `Dockerfile` | App-specific deployment |
| `docker-compose.yml` | App-specific deployment |
| `Caddyfile` | App-specific web server |
| `Makefile` | App-specific build commands |
| `.dockerignore` | App-specific Docker config |
| `.env.example` | App-specific env template |

### Files that STAY in Library

| File | Why Stay |
|------|----------|
| `core/` (all packages) | Framework internals |
| `core/cli/hooks.go` | Callback types (added in Phase A) |
| `core/audit/model.go` | AuditLog struct (moved in Phase A) |
| `database/connection.go` | Generic DB factory |
| `database/resolver.go` | Read/write splitting |
| `database/resolver_test.go` | Generic resolver tests |
| `database/transaction.go` | Transaction helpers |
| `database/transaction_test.go` | Generic transaction tests |
| `database/database_test.go` | Generic DB tests |
| `database/models/base.go` | `BaseModel` struct |
| `database/models/scopes.go` | `WithTrashed()`, `OnlyTrashed()` |
| `database/migrations/migrator.go` | Migration engine |
| `database/seeders/seeder.go` | Seeder engine (registry + run) |
| `testing/testutil/` | Test utilities |
| `go.mod` | Module definition |
| `go.sum` | Dependency checksums |
| `LICENSE` | MIT license |
| `README.md` | Will be rewritten in Phase D |
| `.gitignore` | Stays (update contents for library) |
| `docs/` | Framework documentation (decision: keep v2 docs, archive v1 feature docs) |

### Modifications After Deletion

#### `cmd/main.go` — Minimal Library CLI

After deleting all app code, `cmd/main.go` becomes a minimal entry point that serves as documentation for how to use the library:

```go
package main

import (
	"fmt"
	"os"

	"github.com/RAiWorks/RapidGo/core/cli"
)

func main() {
	// This is the library's built-in CLI.
	// Application projects should use cli.Set*() to wire their code.
	// See: https://github.com/RAiWorks/RapidGo-starter
	fmt.Fprintln(os.Stderr, "RapidGo is a library. Create a project with: rapidgo new myapp")
	fmt.Fprintln(os.Stderr, "Or see: https://github.com/RAiWorks/RapidGo-starter")
	cli.Execute()
}
```

#### `go.mod` — Clean Up Unused Dependencies

After deletions, run:
```bash
go mod tidy
```

This will remove dependencies only used by deleted code (e.g., if any were exclusive to app code). Most dependencies will stay because `core/` packages use them.

#### `.gitignore` — Update for Library

Remove lines for app-specific artifacts:
```
# Remove these lines:
storage/logs/*
storage/cache/*
storage/sessions/*
storage/uploads/*
.env
bin/

# Keep these lines:
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
vendor/
```

### Verification

```bash
# 1. Build check
go build ./...

# 2. Test check
go test ./... -count=1

# 3. Vet check
go vet ./...

# 4. Zero app imports in core/
grep -rn "RAiWorks/RapidGo/app\|RAiWorks/RapidGo/routes\|RAiWorks/RapidGo/http\|RAiWorks/RapidGo/plugins" core/
# Expected: no output

# 5. No app directories
test -d app && echo "FAIL: app/ exists" || echo "OK: app/ removed"
test -d routes && echo "FAIL: routes/ exists" || echo "OK: routes/ removed"
test -d http && echo "FAIL: http/ exists" || echo "OK: http/ removed"
test -d plugins && echo "FAIL: plugins/ exists" || echo "OK: plugins/ removed"

# 6. Key files still exist
test -f database/models/base.go && echo "OK: base.go exists"
test -f database/models/scopes.go && echo "OK: scopes.go exists"
test -f database/migrations/migrator.go && echo "OK: migrator.go exists"
test -f database/seeders/seeder.go && echo "OK: seeder.go exists"
test -f core/cli/hooks.go && echo "OK: hooks.go exists"
test -f core/audit/model.go && echo "OK: audit model.go exists"
```

---

## Step C2: Create RapidGo-starter Repository

### Branch

`feature/v2-08-starter-repo` (from `v2`, after C1 merges)

### Objective

Create the `RapidGo-starter` repository with all the application code removed from the library. Set up its own `go.mod` importing the library. Verify it builds and runs independently.

### Repository Setup

```bash
# On GitHub: create repo RAiWorks/RapidGo-starter
# Locally:
mkdir RapidGo-starter
cd RapidGo-starter
git init
go mod init github.com/RAiWorks/RapidGo-starter
```

### Starter Directory Structure

```
RapidGo-starter/
├── cmd/
│   └── main.go                    ← Full wiring with all cli.Set*() hooks
├── app/
│   ├── helpers/                   ← All helper files from library
│   ├── jobs/
│   │   └── example_job.go
│   ├── providers/                 ← All 8 providers
│   ├── schedule/
│   │   └── schedule.go
│   ├── services/
│   │   ├── user_service.go
│   │   └── user_service_test.go
│   └── plugins.go
├── routes/
│   ├── web.go
│   ├── api.go
│   └── ws.go
├── http/
│   ├── controllers/
│   │   ├── home_controller.go
│   │   ├── post_controller.go
│   │   └── controllers_test.go
│   ├── requests/
│   └── responses/
│       ├── response.go
│       └── response_test.go
├── database/
│   ├── models/
│   │   ├── user.go                ← Embeds fwmodels.BaseModel
│   │   ├── post.go                ← Embeds fwmodels.BaseModel
│   │   └── registry.go           ← All() returns []*User, *Post
│   ├── migrations/
│   │   ├── 20260307000001_create_jobs_tables.go
│   │   ├── 20260308000001_add_soft_deletes.go
│   │   ├── 20260308000002_add_totp_fields.go
│   │   └── 20260308000003_create_audit_logs_table.go
│   └── seeders/
│       └── user_seeder.go
├── resources/
│   ├── views/
│   │   └── home.html
│   ├── lang/
│   └── static/
├── storage/
│   ├── cache/
│   ├── logs/
│   ├── sessions/
│   └── uploads/
├── plugins/
│   └── example/
│       └── example.go
├── tests/
│   ├── integration/
│   └── unit/
├── .env.example
├── .gitignore
├── .dockerignore
├── Dockerfile
├── docker-compose.yml
├── Caddyfile
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

### Key File Modifications in Starter

#### `go.mod`

```
module github.com/RAiWorks/RapidGo-starter

go 1.25

require (
    github.com/RAiWorks/RapidGo v2.0.0
)
```

#### `cmd/main.go`

See the complete starter `main.go` in `v2-architecture.md` section 6.

#### `database/models/user.go` — Updated Import

```go
package models

import (
	fwmodels "github.com/RAiWorks/RapidGo/database/models"
	"github.com/RAiWorks/RapidGo-starter/app/helpers"
)

type User struct {
	fwmodels.BaseModel
	Name     string `gorm:"size:255" json:"name"`
	Email    string `gorm:"size:255;uniqueIndex" json:"email"`
	Password string `gorm:"size:255" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashed, err := helpers.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashed
	}
	return nil
}
```

#### `database/migrations/*.go` — Updated Imports

Each migration file needs to import the library's migrations package:

```go
package migrations

import (
	fwmigrations "github.com/RAiWorks/RapidGo/database/migrations"
	"gorm.io/gorm"
)

func init() {
	fwmigrations.Register(fwmigrations.Migration{
		Version: "20260307000001_create_jobs_tables",
		Up: func(db *gorm.DB) error {
			// ... same migration logic
		},
		Down: func(db *gorm.DB) error {
			// ... same migration logic
		},
	})
}
```

#### `database/seeders/user_seeder.go` — Updated Imports

```go
package seeders

import (
	fwseeders "github.com/RAiWorks/RapidGo/database/seeders"
	"github.com/RAiWorks/RapidGo-starter/database/models"
	"gorm.io/gorm"
)

type UserSeeder struct{}

func (s *UserSeeder) Name() string { return "UserSeeder" }
func (s *UserSeeder) Seed(db *gorm.DB) error {
	// ... same seeder logic
}

func init() {
	fwseeders.Register(&UserSeeder{})
}
```

#### `app/providers/*.go` — Updated Imports

All providers change their module path from `github.com/RAiWorks/RapidGo/app/providers` to match the starter's module. The `core/` imports remain the same since they point to the library:

```go
package providers

import (
	"github.com/RAiWorks/RapidGo/core/config"     // ← library import (unchanged)
	"github.com/RAiWorks/RapidGo/core/container"   // ← library import (unchanged)
)

type ConfigProvider struct{}
// ... same implementation
```

### Verification

```bash
# In starter directory:
go mod tidy
go build ./...
go test ./... -count=1
go vet ./...

# Functional tests:
go run cmd/main.go version
go run cmd/main.go serve     # verify routes respond
go run cmd/main.go migrate   # verify migrations run
go run cmd/main.go db:seed   # verify seeding works

# Back in library directory:
go build ./...
go test ./... -count=1
```

---

## Phase C Checklist

| # | Check | Command |
|---|-------|---------|
| 1 | Library has no `app/` directory | `test ! -d app` |
| 2 | Library has no `routes/` directory | `test ! -d routes` |
| 3 | Library has no `http/` directory | `test ! -d http` |
| 4 | Library has no `plugins/` directory | `test ! -d plugins` |
| 5 | Library `go build ./...` passes | `go build ./...` |
| 6 | Library `go test ./...` passes | `go test ./... -count=1` |
| 7 | Library `go vet ./...` passes | `go vet ./...` |
| 8 | Library `go.mod` has no unused deps | `go mod tidy` returns no changes |
| 9 | Starter `go build ./...` passes | (in starter dir) `go build ./...` |
| 10 | Starter `go test ./...` passes | (in starter dir) `go test ./... -count=1` |
| 11 | Starter `serve` works | `go run cmd/main.go serve` |
| 12 | Starter `migrate` works | `go run cmd/main.go migrate` |
| 13 | Starter `db:seed` works | `go run cmd/main.go db:seed` |
| 14 | Starter `work` starts | `go run cmd/main.go work` |
| 15 | Starter `schedule:run` starts | `go run cmd/main.go schedule:run` |
