# 🧪 Test Plan: Audit Logging

> **Feature**: `51` — Audit Logging
> **Tasks**: [`51-audit-logging-tasks.md`](51-audit-logging-tasks.md)
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## Test File

- `core/audit/audit_test.go`

---

## Unit Tests

### 1. Logger Construction

| # | Test | Expectation |
|---|---|---|
| T01 | `TestNewLogger_ReturnsLogger` | Returns non-nil `*Logger` |

### 2. Logging Entries

| # | Test | Expectation |
|---|---|---|
| T02 | `TestLog_CreateAction` | Persists entry with action "create"; record found in DB with correct UserID, Action, ModelType, ModelID |
| T03 | `TestLog_UpdateWithOldNewValues` | Entry with OldValues and NewValues stored as valid JSON strings in DB |
| T04 | `TestLog_DeleteAction` | Entry with action "delete" and OldValues persisted correctly |
| T05 | `TestLog_WithMetadata` | Metadata map serialized to JSON and stored in Metadata column |
| T06 | `TestLog_NilMapsStoreEmpty` | OldValues, NewValues, Metadata as nil → stored as empty strings (not "null") |
| T07 | `TestLog_ZeroUserID` | UserID 0 (system action) persisted without error |
| T08 | `TestLog_CustomAction` | Action "login" (non-standard verb) persisted correctly |

### 3. Querying Entries

| # | Test | Expectation |
|---|---|---|
| T09 | `TestFind_ByUserID` | Returns only entries matching the given user_id |
| T10 | `TestFind_OrderedNewestFirst` | Results are ordered by `created_at DESC` |
| T11 | `TestFind_NoResults` | Returns empty slice and nil error when no entries match |
| T12 | `TestForModel_ReturnsMatchingEntries` | Returns only entries for the given model_type and model_id |
| T13 | `TestForModel_DifferentModelTypes` | Two entries for different model types; ForModel returns only the correct one |

### 4. AuditLog Model

| # | Test | Expectation |
|---|---|---|
| T14 | `TestAuditLog_NoDeletedAt` | AuditLog struct has no `DeletedAt` field (audit logs are immutable) |

---

## Acceptance Criteria

1. All 14 tests pass
2. All existing tests across all packages still pass (`go test ./...`)
3. `core/audit` package provides `NewLogger`, `Log`, `Find`, `ForModel`
4. `AuditLog` model registered in `models.All()`
5. Migration file compiles and follows existing patterns
6. `go build` succeeds with no errors

---

## Next

Changelog → `51-audit-logging-changelog.md`
