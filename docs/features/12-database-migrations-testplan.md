# 🧪 Test Plan: Database Migrations

> **Feature**: `12` — Database Migrations
> **Architecture**: [`12-database-migrations-architecture.md`](12-database-migrations-architecture.md)
> **Status**: ⬜ NOT RUN
> **Result**: ⬜ NOT RUN

---

## Test File

`database/migrations/migrations_test.go`

All tests use SQLite `:memory:` via GORM — no external database required.

---

## Test Cases

### TC-01: `TestModelsAll`
**What**: `models.All()` returns all registered models.
**How**: Call `models.All()`, assert length is 2, assert types are `*User` and `*Post`.
**Pass**: Returns 2 models of correct types.

### TC-02: `TestNewMigrator_CreatesTable`
**What**: `NewMigrator` auto-creates the `schema_migrations` table.
**How**: Open SQLite `:memory:`, call `NewMigrator(db)`, verify table exists via `db.Migrator().HasTable(&SchemaMigration{})`.
**Pass**: No error, table exists.

### TC-03: `TestMigrator_RunPending`
**What**: `Run()` applies pending migrations.
**How**: Register 2 test migrations (create tables via raw SQL). Call `Run()`. Assert returns 2. Verify tables exist. Verify 2 rows in `schema_migrations` with batch 1.
**Pass**: 2 applied, tables exist, tracking rows present.

### TC-04: `TestMigrator_RunIdempotent`
**What**: `Run()` skips already-applied migrations.
**How**: Register 2 migrations. Call `Run()` twice. Second call should return 0.
**Pass**: First call returns 2, second returns 0.

### TC-05: `TestMigrator_RunOrder`
**What**: Migrations run in version-sorted order.
**How**: Register migrations with versions "20260306_b" and "20260306_a" (out of order). Track execution order via a slice. Call `Run()`. Assert "a" ran before "b".
**Pass**: Migrations execute in sorted version order.

### TC-06: `TestMigrator_Rollback`
**What**: `Rollback()` undoes the last batch.
**How**: Register 2 migrations. `Run()`. `Rollback()`. Assert returns 2. Verify tables dropped. Verify `schema_migrations` is empty.
**Pass**: 2 rolled back, tables gone, tracking empty.

### TC-07: `TestMigrator_RollbackBatches`
**What**: Rollback only undoes the last batch, not all.
**How**: Register migration A. `Run()` (batch 1). Register migration B. `Run()` (batch 2). `Rollback()`. Assert returns 1 (only B). Assert A's table still exists.
**Pass**: Only batch 2 rolled back.

### TC-08: `TestMigrator_Status`
**What**: `Status()` returns correct applied/pending status.
**How**: Register 2 migrations. `Run()` only the first (by calling Run with 1 registered, then register 2nd). Call `Status()`. Assert first is Applied (batch 1), second is Pending.
**Pass**: Correct statuses for both.

### TC-09: `TestMigrator_RollbackEmpty`
**What**: `Rollback()` with no applied migrations returns 0.
**How**: Create migrator, call `Rollback()` with no migrations applied.
**Pass**: Returns 0, no error.

### TC-10: `TestToSnakeCase`
**File**: `core/cli/cli_test.go` (same package as `make_migration.go`)
**What**: `toSnakeCase` converts PascalCase to snake_case.
**How**: Test cases: `"CreateUsersTable"` → `"create_users_table"`, `"addEmailIndex"` → `"add_email_index"`, `"simple"` → `"simple"`, `"ABCTest"` → `"a_b_c_test"`.
**Pass**: All conversions correct.

---

## Test Summary

| ID | Test Name | File | Type | Scope |
|---|---|---|---|---|
| TC-01 | `TestModelsAll` | `database/migrations/migrations_test.go` | Unit | `models.All()` |
| TC-02 | `TestNewMigrator_CreatesTable` | `database/migrations/migrations_test.go` | Integration | Schema tracking |
| TC-03 | `TestMigrator_RunPending` | `database/migrations/migrations_test.go` | Integration | Migration execution |
| TC-04 | `TestMigrator_RunIdempotent` | `database/migrations/migrations_test.go` | Integration | Idempotency |
| TC-05 | `TestMigrator_RunOrder` | `database/migrations/migrations_test.go` | Integration | Version ordering |
| TC-06 | `TestMigrator_Rollback` | `database/migrations/migrations_test.go` | Integration | Full rollback |
| TC-07 | `TestMigrator_RollbackBatches` | `database/migrations/migrations_test.go` | Integration | Batch isolation |
| TC-08 | `TestMigrator_Status` | `database/migrations/migrations_test.go` | Integration | Status reporting |
| TC-09 | `TestMigrator_RollbackEmpty` | `database/migrations/migrations_test.go` | Integration | Edge case |
| TC-10 | `TestToSnakeCase` | `core/cli/cli_test.go` | Unit | String conversion |

**Total**: 10 test cases
**Expected new test count**: 131 + 10 = 141
