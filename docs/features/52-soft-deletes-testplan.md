# 🧪 Test Plan: Soft Deletes

> **Feature**: `52` — Soft Deletes
> **Tasks**: [`52-soft-deletes-tasks.md`](52-soft-deletes-tasks.md)
> **Status**: � COMPLETE
> **Date**: 2026-03-07

---

## Test Files

- `database/models/scopes_test.go` — scope helper tests
- `app/services/user_service_test.go` — updated + new service tests

---

## Test Helpers

Both test files use the existing `setupTestDB` pattern:

```go
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	// ...
	db.AutoMigrate(&models.User{}, &models.Post{})
	return db
}
```

---

## Unit Tests

### 1. Scope Helpers (`database/models/scopes_test.go`)

| # | Test | Expectation |
|---|---|---|
| T01 | `TestWithTrashed_IncludesDeletedRecords` | `db.Scopes(WithTrashed).Find()` returns both active and soft-deleted records |
| T02 | `TestWithTrashed_IncludesActiveRecords` | `db.Scopes(WithTrashed).Find()` still includes non-deleted records |
| T03 | `TestOnlyTrashed_ReturnsOnlyDeletedRecords` | `db.Scopes(OnlyTrashed).Find()` returns only records with `deleted_at IS NOT NULL` |
| T04 | `TestOnlyTrashed_ExcludesActiveRecords` | `db.Scopes(OnlyTrashed).Find()` does not return active records |
| T05 | `TestDefaultQuery_ExcludesDeletedRecords` | `db.Find()` without scopes excludes soft-deleted records |

### 2. Soft Delete Behavior (`app/services/user_service_test.go`)

| # | Test | Expectation |
|---|---|---|
| T06 | `TestDelete_SoftDeletesUser` | `Delete()` sets `deleted_at`; `GetByID()` returns error; `Unscoped().First()` still finds the record |
| T07 | `TestDelete_SetsDeletedAtTimestamp` | After `Delete()`, the record queried via `Unscoped()` has a non-nil `DeletedAt` |

### 3. Hard Delete (`app/services/user_service_test.go`)

| # | Test | Expectation |
|---|---|---|
| T08 | `TestHardDelete_PermanentlyRemovesUser` | `HardDelete()` removes the record; `Unscoped().First()` also returns error |
| T09 | `TestHardDelete_NonExistentID_NoError` | `HardDelete()` on non-existent ID returns no error (GORM no-op) |

### 4. Restore (`app/services/user_service_test.go`)

| # | Test | Expectation |
|---|---|---|
| T10 | `TestRestore_RecoversSoftDeletedUser` | After `Delete()` + `Restore()`, `GetByID()` returns the user with nil `DeletedAt` |
| T11 | `TestRestore_NonDeletedUser_NoError` | `Restore()` on an active user is a no-op (no error) |

### 5. Existing Test Updates (`app/services/user_service_test.go`)

| # | Test | Expectation |
|---|---|---|
| T12 | `TestDelete_RemovesUser` (TC-08 updated) | Existing test updated to verify soft delete: user hidden from normal queries but still exists via `Unscoped()` |

---

## Acceptance Criteria

1. All 12 new/updated tests pass
2. All existing tests in `app/services/`, `database/models/`, and `database/` still pass
3. `go test ./...` — all packages green, no regressions
4. `BaseModel` includes `DeletedAt` field with correct GORM + JSON tags
5. `WithTrashed` and `OnlyTrashed` scopes produce correct query behavior
6. `UserService.Delete()` soft-deletes (record recoverable)
7. `UserService.HardDelete()` permanently removes record
8. `UserService.Restore()` clears `deleted_at` and makes record queryable again

---

## Next

Changelog → `52-soft-deletes-changelog.md`
