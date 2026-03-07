# 📝 Changelog: Audit Logging

> **Feature**: `51` — Audit Logging
> **Status**: ✅ SHIPPED
> **Date**: 2026-03-07
> **Commit**: `1f56346`

---

## Added

- `core/audit/audit.go` — Audit package: `NewLogger()`, `Log()`, `Find()`, `ForModel()`
- `core/audit/audit_test.go` — 14 unit tests for audit logging operations
- `database/models/audit_log.go` — AuditLog model (ID, UserID, Action, ModelType, ModelID, OldValues, NewValues, Metadata, CreatedAt)
- `database/migrations/20260308000003_create_audit_logs_table.go` — migration creating `audit_logs` table with indexes

## Changed

- `database/models/registry.go` — add `&AuditLog{}` to `All()` slice

## Files

| File | Action |
|---|---|
| `core/audit/audit.go` | NEW |
| `core/audit/audit_test.go` | NEW |
| `database/models/audit_log.go` | NEW |
| `database/models/registry.go` | MODIFIED |
| `database/migrations/20260308000003_create_audit_logs_table.go` | NEW |

## Migration Guide

- Run `migrate` to create the `audit_logs` table
- No new environment variables required
- No breaking changes — existing models and tables are unaffected
- Import `core/audit` and call `audit.NewLogger(db)` to start logging
