# 🧪 Test Plan: Database Seeding

> **Feature**: `13` — Database Seeding
> **Architecture**: [`13-database-seeding-architecture.md`](13-database-seeding-architecture.md)
> **Status**: ⬜ NOT RUN
> **Result**: ⬜ NOT RUN

---

## Test File

`database/seeders/seeders_test.go`

All tests use SQLite `:memory:` via GORM — no external database required.

---

## Test Cases

### TC-01: `TestRegister_AddsSeeder`
**What**: `Register()` adds a seeder to the registry.
**How**: Reset registry, register a mock seeder, call `Names()`, assert length is 1 and name matches.
**Pass**: Registry contains the mock seeder.

### TC-02: `TestRunAll_ExecutesSeeders`
**What**: `RunAll()` calls `Seed()` on all registered seeders.
**How**: Register 2 mock seeders that track calls. Call `RunAll(db)`. Assert both were called.
**Pass**: Both seeders executed.

### TC-03: `TestRunAll_StopsOnError`
**What**: `RunAll()` stops on first error.
**How**: Register seeder A (succeeds) and seeder B (returns error). Call `RunAll(db)`. Assert error returned, seeder A ran, seeder B ran (it's the one that fails).
**Pass**: Error returned from failing seeder.

### TC-04: `TestRunByName_FindsSeeder`
**What**: `RunByName()` executes the named seeder.
**How**: Register 2 mock seeders. Call `RunByName(db, "second")`. Assert only the second was called.
**Pass**: Only named seeder executed.

### TC-05: `TestRunByName_NotFound`
**What**: `RunByName()` returns error for unknown name.
**How**: Reset registry, call `RunByName(db, "nonexistent")`.
**Pass**: Error contains "not found".

### TC-06: `TestUserSeeder_CreatesUsers`
**What**: `UserSeeder` creates admin and regular user.
**How**: AutoMigrate User model, run `UserSeeder.Seed(db)`. Query users. Assert 2 users exist with correct emails and roles.
**Pass**: Admin (admin@example.com, role=admin) and User (user@example.com, role=user) exist.

### TC-07: `TestUserSeeder_Idempotent`
**What**: Running `UserSeeder` twice doesn't create duplicates.
**How**: AutoMigrate, run `UserSeeder.Seed(db)` twice. Count users.
**Pass**: Still exactly 2 users.

---

## Test Summary

| ID | Test Name | Type | Scope |
|---|---|---|---|
| TC-01 | `TestRegister_AddsSeeder` | Unit | Registry |
| TC-02 | `TestRunAll_ExecutesSeeders` | Unit | RunAll orchestration |
| TC-03 | `TestRunAll_StopsOnError` | Unit | Error handling |
| TC-04 | `TestRunByName_FindsSeeder` | Unit | Name-based lookup |
| TC-05 | `TestRunByName_NotFound` | Unit | Error case |
| TC-06 | `TestUserSeeder_CreatesUsers` | Integration | User data creation |
| TC-07 | `TestUserSeeder_Idempotent` | Integration | Idempotency |

**Total**: 7 test cases
**Expected new test count**: 141 + 7 = 148
