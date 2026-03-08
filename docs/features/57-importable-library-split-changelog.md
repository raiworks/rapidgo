# 📝 Changelog: Importable Library Split

> **Feature**: `57` — Importable Library Split
> **Branch**: `v2` (integration) with sub-branches per step
> **Started**: —
> **Completed**: —

---

## Log

<!-- Add entries as you work. Most recent first. -->

### Step D2 — `feature/v2-10-library-readme` → `v2`
- Rewrote `README.md` for importable library: package index (33 packages), hook system docs, quick start

### Step D1 — `feature/v2-09-rapidgo-new-cmd` → `v2` (commit `e0835e7`)
- Created `core/cli/new.go`: `rapidgo new` command — downloads starter zip, extracts, replaces module name, runs `go mod tidy`
- Created `core/cli/new_test.go`: 6 tests (extract, zip slip, module rename, invalid name, existing dir)
- Added `newCmd` to `root.go` init block

### Step C1 — `feature/v2-07-remove-app-code` → `v2` (commit `53304dc`)
- Deleted 8 app directories: `app/`, `routes/`, `http/`, `plugins/`, `resources/`, `storage/`, `tests/`, `reference/`
- Deleted `database/querybuilder/`, app migration files, `user.go`, `post.go`, `audit_log.go`, `registry.go`, `user_seeder.go`, `transaction_example.go`
- Deleted root files: `Dockerfile`, `docker-compose.yml`, `Caddyfile`, `Makefile`, `.dockerignore`, `.env.example`
- Rewrote `models_test.go` with generic `testItem` (no User/Post refs)
- Rewrote `scopes_test.go` with `testScopesItem` (no User refs)
- Removed TC-06/TC-07 from `seeders_test.go` (UserSeeder)
- Removed TC-05/TC-06/TC-07 from `transaction_test.go` (TransferCredits)
- Updated `cmd/main.go` to minimal library CLI
- 76 files changed, 3058 deletions. All 34 packages build and pass tests.

### Gate B — PASSED
- Zero `app/`, `routes/`, `http/`, `plugins/`, `database/models`, `database/seeders` imports in `core/`
- All 40 packages pass, all 7 coupling points resolved

### Step B4 — `feature/v2-06-migrate-decouple` → `v2` (commit `500d9f5`)
- `core/cli/migrate.go`: removed `database/models`, uses `modelRegistryFn()`
- `core/cli/seed.go`: removed `database/seeders`, uses `seederFn()`
- `cmd/main.go`: wired `SetModelRegistry(models.All)` + `SetSeeder()`
- `database/migrations/migrations_test.go`: replaced `models.User/Post/AuditLog` with `testMigrationModel`
- `database/seeders/seeders_test.go`: simplified `setupTestDB`, added `setupUserSeederDB`

### Step B3 — `feature/v2-05-worker-decouple` → `v2` (commit `034d062`)
- `core/cli/work.go`: removed `app/jobs`, `app/providers`, uses `NewApp()` + `jobRegistrar`
- `core/cli/schedule_run.go`: removed `app/providers`, `app/schedule`, uses `NewApp()` + `scheduleRegistrar`
- `cmd/main.go`: wired `SetJobRegistrar(jobs.RegisterJobs)` + `SetScheduleRegistrar(schedule.RegisterSchedule)`

### Step B2 — `feature/v2-04-serve-decouple` → `v2` (commit `d58019c`)
- `core/cli/serve.go`: removed `routes` import, uses `routeRegistrar` callback
- Static file/template setup stays in library
- `cmd/main.go`: wired `SetRoutes()` with mode-conditional Register calls

### Step B1 — `feature/v2-03-root-decouple` → `v2` (commit `a62995c`)
- `core/cli/root.go`: removed `app/providers`, `NewApp()` uses `bootstrapFn`
- `core/cli/cli_test.go`: updated `TestNewApp_ReturnsBootedApp` to use `SetBootstrap` with test bindings
- `cmd/main.go`: wired `SetBootstrap()` with all provider registrations

### Step A2 — `feature/v2-02-audit-decouple` → `v2` (commit `932fdc0`)
- Created `core/audit/model.go` with canonical `AuditLog` struct
- `core/audit/audit.go` + `audit_test.go`: removed `database/models` import
- `database/models/audit_log.go`: replaced struct with `type AuditLog = audit.AuditLog`

### Step A1 — `feature/v2-01-hooks-foundation` → `v2` (commit `2e93063`)
- Created `core/cli/hooks.go`: 6 callback types, 6 package-level vars, 6 `Set*()` functions
- Created `core/cli/hooks_test.go`: 7 tests (nil defaults + setter storage)

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
