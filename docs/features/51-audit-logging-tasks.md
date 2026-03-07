# 📋 Tasks: Audit Logging

> **Feature**: `51` — Audit Logging
> **Architecture**: [`51-audit-logging-architecture.md`](51-audit-logging-architecture.md)
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## Phase A — Data Layer

| # | Task | Detail |
|---|---|---|
| A1 | Create `database/models/audit_log.go` | `AuditLog` struct with ID, UserID, Action, ModelType, ModelID, OldValues, NewValues, Metadata, CreatedAt — NO BaseModel embed |
| A2 | Register `AuditLog` in `database/models/registry.go` | Add `&AuditLog{}` to `All()` slice |
| A3 | Create migration `database/migrations/20260308000003_create_audit_logs_table.go` | Up: `db.AutoMigrate(&AuditLog{})` with all fields and indexes. Down: `db.Migrator().DropTable("audit_logs")` |

**Exit**: `go build` succeeds; migration file compiles; AuditLog model in registry

---

## Phase B — Core Audit Package

| # | Task | Detail |
|---|---|---|
| B1 | Create `core/audit/audit.go` | `Logger` struct, `NewLogger(db)`, `Entry` struct, `Log(e)`, `Find(query, args...)`, `ForModel(modelType, modelID)` |

**Exit**: `go build` succeeds; package compiles with no errors

---

## Phase C — Tests & Verification

| # | Task | Detail |
|---|---|---|
| C1 | Create `core/audit/audit_test.go` | All tests from testplan (T01–T14) using SQLite in-memory |
| C2 | Run `go test ./core/audit/... -v` | All 14 tests pass |
| C3 | Run `go test ./... -count=1` | All packages pass (no regressions) |
| C4 | Run `go build -o bin/rapidgo.exe ./cmd` | Binary builds successfully |

**Exit**: All tests pass, binary builds, no regressions

---

## Next

Test plan → `51-audit-logging-testplan.md`
