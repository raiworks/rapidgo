# ✅ Tasks: Importable Library Split

> **Feature**: `57` — Importable Library Split
> **Architecture**: [`57-importable-library-split-architecture.md`](57-importable-library-split-architecture.md)
> **Branch**: `v2` (integration) with sub-branches per step
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/40 tasks complete

---

## Pre-Flight Checklist

- [ ] Discussion doc is marked COMPLETE
- [ ] Architecture doc is FINALIZED
- [ ] `v2` branch created from latest `main` (v1.0.0)
- [ ] All 56 prior features merged to `main`
- [ ] Test plan doc created
- [ ] Changelog doc created (empty)

---

## Phase A — Foundation (No Breaking Changes)

> Create hook system and decouple AuditLog. Purely additive — no existing file logic changes.

### Step A1: Create `core/cli/hooks.go`
**Branch**: `feature/v2-01-hooks-foundation` (from `v2`)

- [ ] **A1.1** — Create `core/cli/hooks.go` with 6 callback types:
  - `BootstrapFunc`, `RouteRegistrar`, `JobRegistrar`, `ScheduleRegistrar`, `ModelRegistryFunc`, `SeederFunc`
  - 6 package-level vars (all defaulting to nil)
  - 6 `Set*()` setter functions
- [ ] **A1.2** — Create `core/cli/hooks_test.go`:
  - `TestHooksDefaultNil` — all 6 vars default to nil
  - `TestSetBootstrapStoresFunction` — setter stores and callback executes
  - `TestSetRoutesStoresFunction` — setter stores and callback executes
  - `TestSetJobRegistrarStoresFunction` — setter stores and callback executes
  - `TestSetModelRegistryStoresFunction` — setter stores and returns values
- [ ] **A1.3** — Run `go build ./...` — compiles with no errors
- [ ] **A1.4** — Run `go test ./core/cli/ -run TestHooks -v` + `go test ./core/cli/ -run TestSet -v` — pass
- [ ] **A1.5** — Run `go test ./... -count=1` — full suite passes (no regressions)
- [ ] **A1.6** — Commit: `feat(cli): add hook system for importable library split`
- [ ] **A1.7** — Push and merge into `v2`

### Step A2: Move AuditLog to `core/audit/`
**Branch**: `feature/v2-02-audit-decouple` (from `v2`, after A1 merged)

- [ ] **A2.1** — Create `core/audit/model.go` with `AuditLog` struct (copied from `database/models/audit_log.go`)
- [ ] **A2.2** — Modify `core/audit/audit.go`: remove `database/models` import, replace `models.AuditLog` → `AuditLog`
- [ ] **A2.3** — Modify `core/audit/audit_test.go`: remove `database/models` import, replace `models.AuditLog` → `AuditLog`
- [ ] **A2.4** — Modify `database/models/audit_log.go`: replace struct with `type AuditLog = audit.AuditLog` alias
- [ ] **A2.5** — Verify `grep "database/models" core/audit/audit.go` returns empty
- [ ] **A2.6** — Run `go test ./core/audit/ -v` — passes
- [ ] **A2.7** — Run `go test ./database/models/ -v` — passes (type alias works)
- [ ] **A2.8** — Run `go test ./... -count=1` — full suite passes
- [ ] **A2.9** — Commit: `refactor(audit): move AuditLog struct to core/audit/model.go`
- [ ] **A2.10** — Push and merge into `v2`

- [ ] 📍 **Checkpoint A** — `hooks.go` exists, AuditLog in `core/audit/`, `go build` + `go test` pass

---

## Phase B — Decouple (Break Hard Imports)

> Replace all app-specific imports in `core/` with hook callbacks. Wire hooks in `cmd/main.go`. Monolith still works.

### Step B1: Decouple `root.go`
**Branch**: `feature/v2-03-root-decouple` (from `v2`)

- [ ] **B1.1** — Modify `core/cli/root.go`: remove `app/providers` import, replace `NewApp()` body with `bootstrapFn(app, mode)` call
- [ ] **B1.2** — Modify `cmd/main.go`: add `cli.SetBootstrap()` with all 8 provider registrations
- [ ] **B1.3** — Verify `grep "app/providers" core/cli/root.go` returns empty
- [ ] **B1.4** — Run `go build ./...` + `go test ./... -count=1` — passes
- [ ] **B1.5** — Commit: `refactor(cli): decouple root.go from app/providers via SetBootstrap`
- [ ] **B1.6** — Push and merge into `v2`

### Step B2: Decouple `serve.go`
**Branch**: `feature/v2-04-serve-decouple` (from `v2`, after B1)

- [ ] **B2.1** — Modify `core/cli/serve.go`: remove `routes` import, replace `routes.Register*()` calls with `routeRegistrar(r, c, m)` callback
- [ ] **B2.2** — Modify `cmd/main.go`: add `cli.SetRoutes()` with route registration logic
- [ ] **B2.3** — Verify `grep "\".*routes\"" core/cli/serve.go` returns only `routeRegistrar` references
- [ ] **B2.4** — Run `go build ./...` + `go test ./... -count=1` — passes
- [ ] **B2.5** — Commit: `refactor(cli): decouple serve.go from routes via SetRoutes`
- [ ] **B2.6** — Push and merge into `v2`

### Step B3: Decouple `work.go` and `schedule_run.go`
**Branch**: `feature/v2-05-worker-decouple` (from `v2`, after B1)

- [ ] **B3.1** — Modify `core/cli/work.go`: remove `app/jobs` + `app/providers` imports, replace manual bootstrap with `NewApp(service.ModeAll)`, replace `jobs.RegisterJobs()` with `jobRegistrar()` callback
- [ ] **B3.2** — Modify `core/cli/schedule_run.go`: remove `app/providers` + `app/schedule` imports, replace manual bootstrap with `NewApp(service.ModeAll)`, replace `schedule.RegisterSchedule()` with `scheduleRegistrar()` callback
- [ ] **B3.3** — Modify `cmd/main.go`: add `cli.SetJobRegistrar(jobs.RegisterJobs)` + `cli.SetScheduleRegistrar(schedule.RegisterSchedule)`
- [ ] **B3.4** — Verify `grep "app/jobs\|app/providers\|app/schedule" core/cli/work.go core/cli/schedule_run.go` returns empty
- [ ] **B3.5** — Run `go build ./...` + `go test ./... -count=1` — passes
- [ ] **B3.6** — Commit: `refactor(cli): decouple work.go and schedule_run.go via hooks`
- [ ] **B3.7** — Push and merge into `v2`

### Step B4: Decouple `migrate.go` + `seed.go` + Test Refactoring
**Branch**: `feature/v2-06-migrate-decouple` (from `v2`, after B1)

- [ ] **B4.1** — Modify `core/cli/migrate.go`: remove `database/models` import, replace `models.All()` with `modelRegistryFn()` guarded by nil check
- [ ] **B4.2** — Modify `core/cli/seed.go`: remove `database/seeders` import, replace `seeders.RunByName()`/`RunAll()` with `seederFn(db, name)` call
- [ ] **B4.3** — Modify `cmd/main.go`: add `cli.SetModelRegistry(models.All)` + `cli.SetSeeder(...)` wrapper
- [ ] **B4.4** — Refactor `database/models/models_test.go`: replace `User{}`/`Post{}` with test-only `testModel` struct
- [ ] **B4.5** — Refactor `database/models/scopes_test.go`: replace `User{}`/`Post{}` with test-only `testScopesModel` struct
- [ ] **B4.6** — Refactor `database/migrations/migrations_test.go`: remove `database/models` import, use `testMigrationModel` struct
- [ ] **B4.7** — Refactor `database/seeders/seeders_test.go`: remove `database/models` import, use `testSeederModel` struct
- [ ] **B4.8** — Verify `grep "database/models" core/cli/migrate.go` + `grep "database/seeders" core/cli/seed.go` return empty
- [ ] **B4.9** — Run `go test ./database/models/ -v` + `go test ./database/migrations/ -v` + `go test ./database/seeders/ -v` — pass
- [ ] **B4.10** — Run `go test ./... -count=1` — full suite passes
- [ ] **B4.11** — Commit: `refactor(cli): decouple migrate.go and seed.go via hooks; refactor test models`
- [ ] **B4.12** — Push and merge into `v2`

- [ ] 📍 **Checkpoint B** — Zero `app/`/`routes/`/`http/`/`plugins/` imports in `core/`. Verify with:
  ```
  grep -rn "RAiWorks/RapidGo/app\|RAiWorks/RapidGo/routes\|RAiWorks/RapidGo/http\|RAiWorks/RapidGo/plugins" core/
  ```
  Must return **zero results**. `rapidgo serve` starts and serves routes.

---

## Phase C — Split (Separate Repos)

> Delete app code from library. Create starter repo. **This phase is breaking.**

### Step C1: Remove App Code from Library
**Branch**: `feature/v2-07-remove-app-code` (from `v2`)

- [ ] **C1.1** — Delete directories: `app/`, `routes/`, `http/`, `plugins/`, `resources/`, `storage/`, `tests/`, `reference/`
- [ ] **C1.2** — Delete app-specific database files: `database/models/user.go`, `post.go`, `audit_log.go`, `registry.go`; `database/transaction_example.go`
- [ ] **C1.3** — Delete app-specific migrations: all `20260*` files in `database/migrations/`
- [ ] **C1.4** — Delete app-specific seeders: `database/seeders/user_seeder.go`
- [ ] **C1.5** — Delete root files: `Dockerfile`, `docker-compose.yml`, `Caddyfile`, `Makefile`, `.dockerignore`, `.env.example`
- [ ] **C1.6** — Simplify `cmd/main.go` to minimal library entrypoint
- [ ] **C1.7** — Update `.gitignore` for library (remove app-specific ignores)
- [ ] **C1.8** — Run `go mod tidy` to remove unused dependencies
- [ ] **C1.9** — Run `go build ./...` — passes
- [ ] **C1.10** — Run `go test ./... -count=1` — passes
- [ ] **C1.11** — Run `go vet ./...` — passes
- [ ] **C1.12** — Verify no app directories remain: `app/`, `routes/`, `http/`, `plugins/`
- [ ] **C1.13** — Verify key files exist: `database/models/base.go`, `database/models/scopes.go`, `database/migrations/migrator.go`, `database/seeders/seeder.go`, `core/cli/hooks.go`, `core/audit/model.go`
- [ ] **C1.14** — Commit: `refactor(core): remove application code — library is standalone`
- [ ] **C1.15** — Push and merge into `v2`

### Step C2: Create RapidGo-starter Repository
**Branch**: `feature/v2-08-starter-repo` (from `v2`)

- [ ] **C2.1** — Create GitHub repo `RAiWorks/RapidGo-starter`
- [ ] **C2.2** — Initialize `go.mod` with `module github.com/RAiWorks/RapidGo-starter`
- [ ] **C2.3** — Copy all deleted app code into starter structure (see architecture doc for layout)
- [ ] **C2.4** — Update all import paths from `github.com/RAiWorks/RapidGo/app/...` → local paths
- [ ] **C2.5** — Update `database/models/*.go` to embed `fwmodels.BaseModel` from library
- [ ] **C2.6** — Update `database/migrations/*.go` to use `fwmigrations.Register()` with `init()`
- [ ] **C2.7** — Update `database/seeders/*.go` to use `fwseeders.Register()` with `init()`
- [ ] **C2.8** — Create `cmd/main.go` with full hook wiring (6 `cli.Set*()` calls)
- [ ] **C2.9** — Run `go mod tidy` in starter
- [ ] **C2.10** — Run `go build ./...` in starter — passes
- [ ] **C2.11** — Run `go test ./... -count=1` in starter — passes
- [ ] **C2.12** — Verify `go run cmd/main.go version` works
- [ ] **C2.13** — Commit and push starter repo
- [ ] **C2.14** — Verify library still builds: `go build ./...` + `go test ./... -count=1`

- [ ] 📍 **Checkpoint C** — Library and starter both build + test independently. Starter `serve`/`migrate`/`db:seed` commands work.

---

## Phase D — Polish (CLI + Documentation)

> Add scaffolding command and finalize documentation. Tag v2.0.0.

### Step D1: `rapidgo new` CLI Command
**Branch**: `feature/v2-09-rapidgo-new-cmd` (from `v2`)

- [ ] **D1.1** — Create `core/cli/new.go` with `newCmd`:
  - Download starter zip from GitHub
  - Extract with zip-slip protection
  - Replace module name in all `.go` files + `go.mod`
  - Run `go mod tidy`
- [ ] **D1.2** — Create `core/cli/new_test.go` — test validation, error cases
- [ ] **D1.3** — Modify `core/cli/root.go`: add `rootCmd.AddCommand(newCmd)` in `init()`
- [ ] **D1.4** — Run `go build ./...` + `go test ./core/cli/ -run TestNew -v` — passes
- [ ] **D1.5** — Integration test: `go run cmd/main.go new testproject` → `cd testproject && go build ./...`
- [ ] **D1.6** — Commit: `feat(cli): add 'rapidgo new' scaffolding command`
- [ ] **D1.7** — Push and merge into `v2`

### Step D2: Documentation (READMEs)
**Branch**: `feature/v2-10-library-readme` (from `v2`, independent of D1)

- [ ] **D2.1** — Rewrite library `README.md`: install instructions, package index, hook system docs
- [ ] **D2.2** — Create starter `README.md`: getting-started guide, project structure, hook wiring
- [ ] **D2.3** — Verify all links in both READMEs
- [ ] **D2.4** — Commit: `docs: write library and starter READMEs for v2`
- [ ] **D2.5** — Push and merge into `v2`

- [ ] 📍 **Checkpoint D** — `rapidgo new myapp` creates a working project. Both READMEs render correctly.

---

## Ship 🚀

- [ ] All phases complete (A, B, C, D)
- [ ] All checkpoints verified
- [ ] Tag `v2.0.0` on `v2` branch: `git tag -a v2.0.0 -m "v2.0.0 — Importable library split"`
- [ ] Push tag: `git push origin v2.0.0`
- [ ] Set `v2` as default branch on GitHub
- [ ] Tag `v1.0.0` on starter repo
- [ ] Verify `go get github.com/RAiWorks/RapidGo@v2.0.0` works
- [ ] Create GitHub releases on both repos
- [ ] **Keep all feature branches** — do not delete
- [ ] Create review doc → `57-importable-library-split-review.md`
