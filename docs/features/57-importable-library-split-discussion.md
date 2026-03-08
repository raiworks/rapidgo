# 💬 Discussion: Importable Library Split

> **Feature**: `57` — Importable Library Split
> **Status**: 🟢 COMPLETE
> **Branch**: `feature/57-importable-library-split`
> **Depends On**: All features #01–#56 (v1.0.0 complete)
> **Date Started**: 2026-03-07
> **Date Completed**: 2026-03-08

---

## Summary

Transform RapidGo from a monolithic starter project (clone-and-build-inside) into an **importable Go library** (`go get github.com/RAiWorks/RapidGo`) with a companion **starter template** (`RapidGo-starter`). This enables users to create new projects via `rapidgo new myapp` or `go get`, rather than cloning and modifying the framework repo directly.

---

## Functional Requirements

- As a developer, I want to `go get github.com/RAiWorks/RapidGo` so that I can import the framework as a dependency in my own Go module
- As a developer, I want to run `rapidgo new myapp` so that I can scaffold a new project from a clean starter template
- As a developer, I want the library to have zero imports from application-specific code (`app/`, `routes/`, `http/`, `plugins/`) so that it compiles standalone
- As a developer, I want backward compatibility during the transition so that the monolith keeps working throughout the refactor

## Current State / Reference

### What Exists

RapidGo v1.0.0 is a complete monolithic framework with 56 features shipped. The framework code (`core/`, `database/`) is interleaved with application code (`app/`, `routes/`, `http/`, `plugins/`, `resources/`, `storage/`). The `core/cli/` package hard-imports application packages in 7 places, making the framework impossible to `go get` independently.

### What Works Well

- The `database/migrations/migrator.go` engine uses a global registry with `Register()` + `init()` pattern — already library-friendly
- The `database/seeders/seeder.go` engine uses the same registry pattern — already library-friendly
- All `core/` packages except `core/cli/` and `core/audit/` have zero app-specific imports
- The plugin system (`core/plugin/`) has zero coupling to `plugins/` — purely inward imports

### What Needs Improvement

- `core/cli/root.go` hard-codes 8 provider registrations from `app/providers`
- `core/cli/serve.go` hard-imports `routes` package for `RegisterWeb/API/WS()`
- `core/cli/work.go` hard-imports `app/jobs` and `app/providers` with manual bootstrap
- `core/cli/schedule_run.go` hard-imports `app/providers` and `app/schedule` with manual bootstrap
- `core/cli/migrate.go` hard-imports `database/models` for `models.All()`
- `core/cli/seed.go` hard-imports `database/seeders` for `RunByName()`/`RunAll()`
- `core/audit/audit.go` hard-imports `database/models` for the `AuditLog` struct
- 4 test files reference app models (`User`, `Post`) that won't exist after the split

## Proposed Approach

**Callback hook architecture**: Introduce 6 `Set*()` functions in `core/cli/hooks.go` that accept function callbacks. The starter's `main.go` calls these setters before `cli.Execute()`, wiring application code into library commands without hard imports. The approach is:

1. **Phase A (Foundation)**: Create `hooks.go` with 6 callback types. Move `AuditLog` into `core/audit/`. No existing code changes.
2. **Phase B (Decouple)**: Replace hard imports in all 6 `core/cli/` files with hook callbacks. Wire hooks in `cmd/main.go`. Refactor 4 test files.
3. **Phase C (Split)**: Delete app code from library repo. Create `RapidGo-starter` repo with the removed code.
4. **Phase D (Polish)**: Add `rapidgo new` CLI command. Write READMEs for both repos. Tag v2.0.0.

Each phase is non-breaking until Phase C. The monolith keeps working throughout Phases A and B because `cmd/main.go` wires the same code via hooks.

## Edge Cases & Risks

- [x] **Circular import risk**: `database/models/audit_log.go` type alias imports `core/audit` — verify no cycle exists (audit → models → audit). Confirmed: no cycle because the alias replaces the struct, it doesn't import back.
- [x] **Test isolation**: 4 test files (`models_test.go`, `scopes_test.go`, `migrations_test.go`, `seeders_test.go`) use `User`/`Post` models that move to starter. Fix: replace with test-only model structs.
- [x] **Nil hook safety**: Commands must handle nil hooks gracefully (e.g., `if bootstrapFn != nil`). Without this, a bare library install would panic.
- [x] **`go mod tidy` after split**: Library may have unused dependencies after removing app code. Run `go mod tidy` to clean.
- [x] **Blank import pattern for migrations**: Starter `main.go` needs `_ "module/database/migrations"` to trigger `init()` registration. This is standard Go practice.

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Features #01–#56 (v1.0.0) | Feature | ✅ Done |
| v1.0.0 tag (commit `9c0a22a`) | Milestone | ✅ Tagged |
| GitHub repo `RapidGo-starter` | Infrastructure | 🔴 Needs creation (Phase C) |

## Open Questions

- [x] How many coupling points exist? → **7** (C1–C7, verified by codebase grep)
- [x] Should migration/seeder engines stay in library or move? → **Stay** (they are generic engines; app code registers via `init()`)
- [x] Is `SetMigrationRegistrar` needed? → **No** (engine stays, `Register()` + `init()` pattern suffices)
- [x] Is `SetPluginRegistrar` needed? → **No** (zero `core/ → plugins/` coupling exists)
- [x] How many hooks are needed? → **6** (SetBootstrap, SetRoutes, SetJobRegistrar, SetScheduleRegistrar, SetModelRegistry, SetSeeder)
- [x] What about `work.go`/`schedule_run.go` manual bootstrap? → Replace with `NewApp(service.ModeAll)` which uses `bootstrapFn`

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-07 | 6 hooks, not 8 | `SetMigrationRegistrar` and `SetPluginRegistrar` are unnecessary — engines stay in library, plugins have no outward coupling |
| 2026-03-07 | 7 coupling points, not 5 or 9 | Audit (C7) and seed (C6) are true coupling; migrate_rollback/status are engine imports (not coupling) |
| 2026-03-07 | AuditLog moves to `core/audit/model.go` | Breaks the only `core/ → database/models` dependency; type alias preserves backward compat |
| 2026-03-07 | worker/scheduler use `NewApp()` instead of manual bootstrap | Eliminates duplicate provider registration; `NewApp()` uses `bootstrapFn` after Phase B |
| 2026-03-08 | Two repos: `RapidGo` (library) + `RapidGo-starter` (template) | Clean separation; library is `go get`-able; starter is clone-and-customize |
| 2026-03-08 | Branch strategy: `v2` branch (not a new repo) | Keeps git history; `v2` becomes default branch when done |

## Discussion Complete ✅

**Summary**: RapidGo will be split into an importable library and a starter template using a 6-hook callback architecture across 4 phases (A–D, 10 steps), with zero breaking changes until Phase C.
**Completed**: 2026-03-08
**Next**: Create architecture doc → `57-importable-library-split-architecture.md`
