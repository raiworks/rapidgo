# 📝 Changelog: Database Migrations

> **Feature**: `12` — Database Migrations
> **Branch**: `feature/12-database-migrations`
> **Started**: 2026-03-06
> **Completed**: 2026-03-06

---

## Log

- **Phase A** — Created `database/models/registry.go`: `All()` returns `[]interface{}{&User{}, &Post{}}`. Build clean.
- **Phase B** — Created `database/migrations/migrator.go`: `SchemaMigration`, `Migration`, `MigrationStatus`, `Register()`, `ResetRegistry()`, `NewMigrator()`, `Run()`, `Rollback()`, `Status()`. Build + vet clean.
- **Phase C** — Created 4 CLI commands (`migrate.go`, `migrate_rollback.go`, `migrate_status.go`, `make_migration.go`) + registered in `root.go`. Build + vet clean.
- **Phase D** — Created 10 tests (9 in `migrations_test.go`, 1 in `cli_test.go`). TC-08 failed initially (version sort order: `iota` < `theta`). Fixed by using prefixed versions (`a_create_theta`, `b_create_iota`). All 10 pass. Full regression: 141 tests, 0 failures.
- **Phase E** — Cross-check passed: all files/signatures match architecture. One minor addition documented below.

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| Added `ResetRegistry()` | Not in architecture | Exported function to clear global registry | Tests need isolated registry state between test cases |
| TC-08 version names | `create_theta` / `create_iota` | `a_create_theta` / `b_create_iota` | Original names sorted `iota` before `theta` alphabetically, breaking index-based assertions |

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
| Expose `ResetRegistry()` in migrator.go | Tests mutate the global registry; each test needs clean state | 2026-03-06 |
