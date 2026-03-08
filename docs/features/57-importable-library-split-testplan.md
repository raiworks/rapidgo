# 🧪 Test Plan: Importable Library Split

> **Feature**: `57` — Importable Library Split
> **Tasks**: [`57-importable-library-split-tasks.md`](57-importable-library-split-tasks.md)
> **Date**: 2026-03-08

---

## Acceptance Criteria

- [ ] `core/cli/hooks.go` exists with 6 type definitions and 6 `Set*()` functions
- [ ] `core/audit/model.go` contains the `AuditLog` struct (no `database/models` import in `core/audit/`)
- [ ] Zero `app/`/`routes/`/`http/`/`plugins/` imports in any `core/` file
- [ ] `cmd/main.go` wires all 6 hooks before `cli.Execute()`
- [ ] All test files in `database/` use test-only model structs (no `User`/`Post` references)
- [ ] Library repo builds and tests standalone (no `app/`, `routes/`, `http/`, `plugins/` directories)
- [ ] `RapidGo-starter` repo builds and runs (`serve`, `migrate`, `db:seed`)
- [ ] `rapidgo new myapp` scaffolds a working project
- [ ] `go get github.com/RAiWorks/RapidGo@v2.0.0` succeeds
- [ ] v2.0.0 tag exists on library repo

---

## Test Cases

### Phase A — Foundation

#### TC-A01: Hook Defaults

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `hooks.go` created |
| **Steps** | 1. Reset all hook vars to nil → 2. Check each var is nil |
| **Expected Result** | All 6 hook variables default to nil |
| **Status** | ⬜ Not Run |
| **Notes** | `TestHooksDefaultNil` |

#### TC-A02: SetBootstrap Stores Callback

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `hooks.go` created |
| **Steps** | 1. Call `SetBootstrap(fn)` → 2. Verify `bootstrapFn` is not nil → 3. Call `bootstrapFn(nil, ModeAll)` → 4. Verify callback was invoked |
| **Expected Result** | Function is stored and callable |
| **Status** | ⬜ Not Run |
| **Notes** | `TestSetBootstrapStoresFunction` |

#### TC-A03: SetRoutes Stores Callback

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `hooks.go` created |
| **Steps** | 1. Call `SetRoutes(fn)` → 2. Verify `routeRegistrar` is not nil → 3. Invoke it → 4. Verify callback executed |
| **Expected Result** | Function is stored and callable |
| **Status** | ⬜ Not Run |
| **Notes** | `TestSetRoutesStoresFunction` |

#### TC-A04: SetJobRegistrar Stores Callback

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `hooks.go` created |
| **Steps** | 1. Call `SetJobRegistrar(fn)` → 2. Verify not nil → 3. Invoke → 4. Verify called |
| **Expected Result** | Function is stored and callable |
| **Status** | ⬜ Not Run |
| **Notes** | `TestSetJobRegistrarStoresFunction` |

#### TC-A05: SetModelRegistry Stores and Returns

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `hooks.go` created |
| **Steps** | 1. Call `SetModelRegistry(fn)` where fn returns `[]interface{}{"test"}` → 2. Invoke `modelRegistryFn()` → 3. Check length is 1 |
| **Expected Result** | Stored function returns expected values |
| **Status** | ⬜ Not Run |
| **Notes** | `TestSetModelRegistryStoresFunction` |

#### TC-A06: AuditLog in core/audit/

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `core/audit/model.go` created, `audit.go` + `audit_test.go` modified |
| **Steps** | 1. `go test ./core/audit/ -v` → 2. `grep "database/models" core/audit/audit.go` |
| **Expected Result** | Tests pass; grep returns empty (no `database/models` import) |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-A07: Type Alias Backward Compatibility

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `database/models/audit_log.go` updated to type alias |
| **Steps** | 1. `go test ./database/models/ -v` → 2. Code using `models.AuditLog` compiles |
| **Expected Result** | Type alias works transparently |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-A08: Full Suite Post-Phase-A

| Property | Value |
|---|---|
| **Category** | Regression |
| **Precondition** | Phase A complete |
| **Steps** | 1. `go build ./...` → 2. `go test ./... -count=1` |
| **Expected Result** | Zero errors, all tests pass |
| **Status** | ⬜ Not Run |
| **Notes** | Gate A verification |

### Phase B — Decouple

#### TC-B01: root.go Decoupled

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | B1 complete |
| **Steps** | 1. `grep "app/providers" core/cli/root.go` → 2. `go build ./...` → 3. `go test ./... -count=1` |
| **Expected Result** | No `app/providers` import; passes build + tests |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-B02: serve.go Decoupled

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | B2 complete |
| **Steps** | 1. `grep '".*routes"' core/cli/serve.go` → 2. `go build ./...` → 3. `go test ./... -count=1` |
| **Expected Result** | No `routes` package import; passes build + tests |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-B03: work.go Decoupled

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | B3 complete |
| **Steps** | 1. `grep "app/jobs\|app/providers" core/cli/work.go` → 2. Build + test |
| **Expected Result** | No app imports; passes |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-B04: schedule_run.go Decoupled

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | B3 complete |
| **Steps** | 1. `grep "app/providers\|app/schedule" core/cli/schedule_run.go` → 2. Build + test |
| **Expected Result** | No app imports; passes |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-B05: migrate.go Decoupled

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | B4 complete |
| **Steps** | 1. `grep "database/models" core/cli/migrate.go` → 2. Build + test |
| **Expected Result** | No `database/models` import; passes |
| **Status** | ⬜ Not Run |
| **Notes** | `database/migrations` import stays (engine) |

#### TC-B06: seed.go Decoupled

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | B4 complete |
| **Steps** | 1. `grep "database/seeders" core/cli/seed.go` → 2. Build + test |
| **Expected Result** | No `database/seeders` import; passes |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-B07: Test Files Refactored

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | B4 complete |
| **Steps** | 1. `go test ./database/models/ -v` → 2. `go test ./database/migrations/ -v` → 3. `go test ./database/seeders/ -v` → 4. Verify no `User`/`Post` references in test files (only test-only structs) |
| **Expected Result** | All tests pass; no app model references |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-B08: Zero Coupling Verification

| Property | Value |
|---|---|
| **Category** | Regression |
| **Precondition** | All Phase B complete |
| **Steps** | 1. `grep -rn "RAiWorks/RapidGo/app\|RAiWorks/RapidGo/routes\|RAiWorks/RapidGo/http\|RAiWorks/RapidGo/plugins" core/` |
| **Expected Result** | Zero results |
| **Status** | ⬜ Not Run |
| **Notes** | Gate B — the critical milestone |

#### TC-B09: Monolith Still Works

| Property | Value |
|---|---|
| **Category** | Regression |
| **Precondition** | All Phase B complete |
| **Steps** | 1. `go build ./...` → 2. `go test ./... -count=1` → 3. `go run cmd/main.go version` |
| **Expected Result** | Build passes, tests pass, version prints |
| **Status** | ⬜ Not Run |
| **Notes** | Monolith continuity check |

### Phase C — Split

#### TC-C01: Library Standalone Build

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | C1 complete (app code removed) |
| **Steps** | 1. `go build ./...` → 2. `go test ./... -count=1` → 3. `go vet ./...` |
| **Expected Result** | All pass with no app code present |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-C02: No App Directories in Library

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | C1 complete |
| **Steps** | 1. Check `app/` does not exist → 2. Check `routes/` does not exist → 3. Check `http/` does not exist → 4. Check `plugins/` does not exist |
| **Expected Result** | None exist |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-C03: Key Library Files Preserved

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | C1 complete |
| **Steps** | Check existence of: `database/models/base.go`, `database/models/scopes.go`, `database/migrations/migrator.go`, `database/seeders/seeder.go`, `core/cli/hooks.go`, `core/audit/model.go` |
| **Expected Result** | All exist |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-C04: Starter Build

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | C2 complete |
| **Steps** | 1. (in starter dir) `go build ./...` → 2. `go test ./... -count=1` |
| **Expected Result** | Passes |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-C05: Starter Commands

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | C2 complete with database configured |
| **Steps** | 1. `go run cmd/main.go version` → 2. `go run cmd/main.go serve` → 3. `go run cmd/main.go migrate` → 4. `go run cmd/main.go db:seed` |
| **Expected Result** | Version prints, server starts, migrations run, seeding works |
| **Status** | ⬜ Not Run |
| **Notes** | Requires database setup for 2–4 |

#### TC-C06: Library go.mod Clean

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | C1 complete |
| **Steps** | 1. Run `go mod tidy` → 2. Check `git diff go.mod go.sum` is empty (no further changes) |
| **Expected Result** | No unused dependencies remain |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-C07: Starter Module Resolution

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | C2 complete |
| **Steps** | 1. (in starter dir) `go mod tidy` → 2. Verify `go.mod` imports `github.com/RAiWorks/RapidGo` → 3. `go build ./...` |
| **Expected Result** | Module resolves correctly, builds |
| **Status** | ⬜ Not Run |
| **Notes** | — |
| **Notes** | — |

#### TC-C05: Starter Commands

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | C2 complete with database configured |
| **Steps** | 1. `go run cmd/main.go version` → 2. `go run cmd/main.go serve` → 3. `go run cmd/main.go migrate` → 4. `go run cmd/main.go db:seed` |
| **Expected Result** | Version prints, server starts, migrations run, seeding works |
| **Status** | ⬜ Not Run |
| **Notes** | Requires database setup for 2–4 |

### Phase D — Polish

#### TC-D01: `rapidgo new` — No Args

| Property | Value |
|---|---|
| **Category** | Error |
| **Precondition** | D1 complete |
| **Steps** | 1. Run `rapidgo new` (no args) |
| **Expected Result** | Error + usage displayed |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-D02: `rapidgo new myapp` — Success

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | D1 complete, internet access |
| **Steps** | 1. `rapidgo new testapp` → 2. `cd testapp && go build ./...` → 3. `grep -r "RapidGo-starter" .` |
| **Expected Result** | Project created, builds, no starter module references remain |
| **Status** | ⬜ Not Run |
| **Notes** | Integration test — network required |

#### TC-D03: `rapidgo new existing-dir` — Error

| Property | Value |
|---|---|
| **Category** | Error |
| **Precondition** | D1 complete, target directory exists |
| **Steps** | 1. `mkdir exists` → 2. `rapidgo new exists` |
| **Expected Result** | Error: "directory already exists" |
| **Status** | ⬜ Not Run |
| **Notes** | — |

#### TC-D04: `rapidgo new "bad/name"` — Validation

| Property | Value |
|---|---|
| **Category** | Error |
| **Precondition** | D1 complete |
| **Steps** | 1. `rapidgo new "bad/name"` |
| **Expected Result** | Error: "invalid project name" |
| **Status** | ⬜ Not Run |
| **Notes** | — |

---

## Edge Cases

| # | Scenario | Expected Behavior |
|---|---|---|
| 1 | `SetBootstrap` not called, then `rapidgo serve` | `bootstrapFn` is nil → no providers registered → app starts but no routes/db |
| 2 | `SetModelRegistry` not called, then `rapidgo migrate` | AutoMigrate step skipped (nil check), file migrations still run |
| 3 | `SetSeeder` not called, then `rapidgo db:seed` | Returns error: "no seeder registered" |
| 4 | Zip download fails in `rapidgo new` | Returns error with context: "download failed: ..." |
| 5 | Malicious zip with path traversal | Zip slip protection rejects illegal paths |

## Security Tests

| # | Test | Expected |
|---|---|---|
| 1 | Zip slip attack in `extractZip()` | Paths escaping target directory are rejected with error |
| 2 | Project name with path traversal chars | Rejected: "invalid project name" |
| 3 | No `.env` or secrets in library repo after split | Only starter has `.env.example`; library has none |

## Performance Considerations

| Metric | Target | Actual |
|---|---|---|
| `go build ./...` (library alone) | < 30s | — |
| `go test ./... -count=1` (library alone) | < 60s | — |
| `rapidgo new myapp` total time | < 30s | — |

---

## Test Summary

| Category | Total | Pass | Fail | Skip |
|---|---|---|---|---|
| Phase A | 8 | — | — | — |
| Phase B | 9 | — | — | — |
| Phase C | 7 | — | — | — |
| Phase D | 4 | — | — | — |
| Edge Cases | 5 | — | — | — |
| Security | 3 | — | — | — |
| **Total** | **36** | — | — | — |

**Result**: ⬜ NOT RUN
