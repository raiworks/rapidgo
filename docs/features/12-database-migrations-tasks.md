# ✅ Tasks: Database Migrations

> **Feature**: `12` — Database Migrations
> **Architecture**: [`12-database-migrations-architecture.md`](12-database-migrations-architecture.md)
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/5 phases complete

---

## Phase A — Model Registry

> Create `database/models/registry.go` with `All()` function.

- [ ] **A.1** — Create `database/models/registry.go` with `All()` returning `[]interface{}{&User{}, &Post{}}`
- [ ] **A.2** — `go build ./database/models/...` clean
- [ ] 📍 **Checkpoint A** — `All()` compiles, returns 2 models

## Phase B — Migration Engine

> Create `database/migrations/migrator.go` with Migrator struct and methods.

- [ ] **B.1** — Create `database/migrations/migrator.go`:
  - `SchemaMigration` struct (tracking table)
  - `Migration` struct (Version, Up, Down)
  - `MigrationFunc` type
  - `MigrationStatus` struct
  - Package-level `registry` + `Register()` function
  - `NewMigrator()` — creates Migrator, auto-creates schema_migrations table
  - `Run()` — applies pending migrations in version order
  - `Rollback()` — undoes last batch in reverse order
  - `Status()` — returns sorted list of migration statuses
- [ ] **B.2** — `go build ./database/migrations/...` clean
- [ ] 📍 **Checkpoint B** — Migrator builds, all exported types/functions present

## Phase C — CLI Commands

> Create 4 CLI commands and register in root.

- [ ] **C.1** — Create `core/cli/migrate.go`: `migrateCmd` — AutoMigrate + pending migrations
- [ ] **C.2** — Create `core/cli/migrate_rollback.go`: `migrateRollbackCmd` — rollback last batch
- [ ] **C.3** — Create `core/cli/migrate_status.go`: `migrateStatusCmd` — show status table
- [ ] **C.4** — Create `core/cli/make_migration.go`: `makeMigrationCmd` — generate migration file, include `toSnakeCase` helper
- [ ] **C.5** — Update `core/cli/root.go` — add 4 commands to `init()`
- [ ] **C.6** — `go build ./...` clean
- [ ] 📍 **Checkpoint C** — All 4 commands registered, build clean

## Phase D — Testing

> Create tests for migration engine and CLI commands.

- [ ] **D.1** — Create `database/migrations/migrations_test.go` with test cases from testplan
- [ ] **D.2** — `go test ./database/migrations/... -v` — all pass
- [ ] **D.3** — `go test ./... -count=1` — full regression, 0 failures
- [ ] 📍 **Checkpoint D** — All new tests pass, no regressions

## Phase E — Changelog & Self-Review

- [ ] **E.1** — Update `12-database-migrations-changelog.md` with build log and deviations
- [ ] **E.2** — Cross-check: verify code matches architecture doc
- [ ] 📍 **Checkpoint E** — Changelog complete, architecture consistent
