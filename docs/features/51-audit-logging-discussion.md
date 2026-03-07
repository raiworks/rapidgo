# 💬 Discussion: Audit Logging

> **Feature**: `51` — Audit Logging
> **Status**: ✅ COMPLETE
> **Date**: 2026-03-07

---

## What Are We Building?

A `core/audit` package that records **who did what, to which record, and when** into an `audit_logs` database table. The package provides a lightweight `Logger` struct that application code calls explicitly — no magic GORM hooks or automatic interception. This keeps audit logging opt-in, predictable, and easy to test.

Each audit entry captures: the action (create, update, delete, or a custom verb), the model type and ID, the actor (user ID), an optional diff of old → new values (JSON), and optional metadata. The framework provides the building blocks — the application developer decides which operations to audit.

---

## Why?

- **Compliance**: Regulatory requirements (GDPR, SOC 2, HIPAA) often mandate change tracking for sensitive data
- **Debugging**: "Who changed this field and when?" is a common production question
- **Roadmap**: Feature #51 in the project roadmap, depends on #11 (Models) and #03 (Logging) — both shipped
- **Framework gap**: No existing audit trail capability — developers must build from scratch today

---

## Prior Art

| Framework | Approach |
|---|---|
| **Laravel (Spatie Activity Log)** | Explicit `activity()->log()` API; stores in `activity_log` table with `causer`, `subject`, `properties` JSON columns |
| **Django (django-auditlog)** | Model-level decorator/mixin; automatic field-level diff via `LogEntry` model |
| **Rails (PaperTrail)** | ActiveRecord callbacks; stores `item_type`, `item_id`, `event`, `object` (YAML/JSON serialized previous state) |

**Our approach**: Closest to Laravel's Spatie — explicit API calls, JSON properties column, no automatic hooks. This avoids the complexity and performance overhead of intercepting every GORM operation.

---

## Constraints

1. **Explicit only** — no automatic GORM callback registration; developers call `audit.Log()` where needed
2. **Same database** — audit logs stored in the application database (same connection pool)
3. **JSON storage** — old/new values stored as JSON text columns for flexibility across all DB drivers
4. **No query builder** — simple `Find` and `ForModel` query helpers, not a full query DSL
5. **No retention/cleanup** — framework provides the table; cleanup policies are application-level concerns
6. **Dependency**: `database/models` (BaseModel), `core/config` (env vars) — both shipped

---

## Decision Log

| # | Decision | Rationale |
|---|---|---|
| D1 | Explicit API, not GORM hooks | Hooks fire on every operation globally — too noisy, hard to control, impacts performance. Explicit logging lets developers audit only what matters. |
| D2 | JSON text columns for old/new values | Works across PostgreSQL, MySQL, SQLite without driver-specific JSON types. Simple `encoding/json` marshal/unmarshal. |
| D3 | Single `audit_logs` table | Simpler than per-model audit tables. Model type + ID columns allow filtering. Indexes on `model_type` and `model_id` for query performance. |
| D4 | `*gorm.DB` injected, not global | Logger receives DB via constructor — testable with SQLite in-memory, no global state. |
| D5 | No events integration | Keeping audit logging self-contained. Developers can dispatch events from their own code after calling audit.Log() if needed. |

---

## Open Questions

| # | Question | Answer |
|---|---|---|
| Q1 | Should we auto-diff old/new values? | ✅ No — caller provides old/new maps explicitly. Auto-diff requires reflection and model introspection which adds complexity. |
| Q2 | Should audit logs be soft-deletable? | ✅ No — audit logs should be immutable. Use only `ID` and `CreatedAt`; no `UpdatedAt` or `DeletedAt`. |
| Q3 | IP address / user agent tracking? | ✅ No — that's request-level context, not model-level. The `Metadata` JSON field lets apps store this if needed. |

---

## Next

Architecture → `51-audit-logging-architecture.md`
