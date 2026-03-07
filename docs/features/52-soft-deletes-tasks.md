# 📋 Tasks: Soft Deletes

> **Feature**: `52` — Soft Deletes
> **Architecture**: [`52-soft-deletes-architecture.md`](52-soft-deletes-architecture.md)
> **Status**: � COMPLETE
> **Date**: 2026-03-07

---

## Phase A — BaseModel & Scopes

| # | Task | Detail |
|---|---|---|
| A1 | Add `DeletedAt` to `BaseModel` | Add `gorm.DeletedAt` field with `gorm:"index"` and `json:"deleted_at,omitempty"` tags to `database/models/base.go` |
| A2 | Update `BaseModel` comment | Reflect that soft deletes are now included |
| A3 | Create `scopes.go` | Add `WithTrashed` and `OnlyTrashed` scope functions in `database/models/scopes.go` |
| A4 | Create `scopes_test.go` | Test both scopes verify correct query behavior |

**Exit**: `BaseModel` has `DeletedAt`, scopes compile, scope tests pass.

---

## Phase B — Service Layer

| # | Task | Detail |
|---|---|---|
| B1 | Add `HardDelete` method | `UserService.HardDelete(id)` using `db.Unscoped().Delete()` in `app/services/user_service.go` |
| B2 | Add `Restore` method | `UserService.Restore(id)` using `db.Unscoped().Model().Update("deleted_at", nil)` |
| B3 | Update existing service tests | `TestDelete_RemovesUser` now verifies soft delete behavior (record hidden but exists via Unscoped) |
| B4 | Add new service tests | Tests for `HardDelete`, `Restore`, and `Delete` as soft delete |

**Exit**: `UserService` has `Delete` (soft), `HardDelete`, and `Restore`; all service tests pass.

---

## Phase C — Migration

| # | Task | Detail |
|---|---|---|
| C1 | Create migration file | `database/migrations/20260308000001_add_soft_deletes.go` with Up (add column + index) and Down (drop) |

**Exit**: Migration file compiles, follows existing migration patterns.

---

## Phase D — Verification

| # | Task | Detail |
|---|---|---|
| D1 | Run all tests | `go test ./...` — all packages pass |
| D2 | Verify no regressions | Existing model/service/database tests still pass with soft delete behavior |

**Exit**: All packages green, no regressions.

---

## Next

Test plan → `52-soft-deletes-testplan.md`
