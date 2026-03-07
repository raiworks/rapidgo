# 💬 Discussion: Soft Deletes

> **Feature**: `52` — Soft Deletes
> **Status**: � COMPLETE
> **Date**: 2026-03-07

---

## What Are We Building?

Framework-level soft delete support that marks records as deleted (via a `deleted_at` timestamp) instead of permanently removing them from the database. GORM queries automatically exclude soft-deleted records, while explicit `Unscoped()` calls retrieve or permanently remove them. The feature adds a `DeletedAt` field to `BaseModel`, query scope helpers, service-level restore/hard-delete methods, a migration for existing tables, and comprehensive tests.

---

## Why?

- **Data recovery**: accidental deletions can be reversed without backups or manual SQL
- **Audit trail**: knowing *when* a record was deleted is valuable for compliance and debugging
- **Referential integrity**: related records still have a valid FK target, preventing cascade issues when a parent is soft-deleted
- **Industry standard**: Laravel, Django, Rails, and GORM all provide soft delete out of the box — users expect it
- **Planned since #11**: the `BaseModel` in `database/models/base.go` explicitly excluded `DeletedAt` with a comment reserving it for Feature #52

---

## Prior Art

| System | Approach | Notes |
|---|---|---|
| GORM (Go) | `gorm.DeletedAt` field + `Unscoped()` | Built-in, automatic query filtering, index-friendly |
| Laravel (PHP) | `SoftDeletes` trait + `trashed()` / `withTrashed()` / `restore()` | Trait-based, adds `deleted_at` column |
| Django (Python) | `django-soft-delete` / custom managers | No built-in; community packages provide queryset filtering |
| Rails (Ruby) | `paranoia` gem / `discard` gem | `paranoia` overrides `delete`; `discard` uses explicit `discard`/`undiscard` |

---

## Constraints

1. **Leverage GORM's built-in `gorm.DeletedAt`** — no custom soft delete engine; GORM handles query filtering, timestamps, and `Unscoped()` natively
2. **Add to `BaseModel`** — all models get the `DeletedAt` field; soft deletes are framework-wide, not per-model opt-in (consistent with the convention that all models embed `BaseModel`)
3. **Breaking change acknowledged** — existing `Delete()` calls become soft deletes; `HardDelete()` method provides explicit permanent removal
4. **Migration required** — existing `users` and `posts` tables need an `ALTER TABLE ADD COLUMN deleted_at` migration
5. **Scope helpers** — provide `WithTrashed` and `OnlyTrashed` GORM scope functions for ergonomic querying
6. **Service layer** — `UserService` gets `Restore()` and `HardDelete()` methods alongside the existing `Delete()` which becomes a soft delete
7. **No cascading soft deletes** — soft-deleting a user does NOT auto-soft-delete their posts; each entity is deleted independently (matches GORM default)
8. **Index on `deleted_at`** — required for query performance since GORM filters `WHERE deleted_at IS NULL` on every query

---

## Decision Log

| # | Decision | Rationale |
|---|---|---|
| 1 | Use GORM's built-in `gorm.DeletedAt` | Battle-tested, zero custom code, automatic query filtering, works with all GORM features (Preload, joins, scopes) |
| 2 | Add `DeletedAt` to `BaseModel`, not individual models | Single source of truth; all models behave consistently; matches the existing pattern where all models embed `BaseModel` |
| 3 | Provide `WithTrashed` / `OnlyTrashed` scope helpers | Laravel-inspired naming; idiomatic GORM scopes (`func(db *gorm.DB) *gorm.DB`); avoids repeating `Unscoped()` + `Where()` everywhere |
| 4 | Explicit `HardDelete()` method on services | Clear intent — developer must opt in to permanent deletion; prevents accidental data loss |
| 5 | `Restore()` method on services | Makes recovery discoverable and tested; uses `Unscoped().Model().Update("deleted_at", nil)` pattern |
| 6 | No cascading soft deletes | Keeps behavior simple and predictable; cascading is a separate concern (can be added later via hooks if needed) |
| 7 | JSON tag `json:"deleted_at,omitempty"` on `DeletedAt` | Omitted from API responses when nil (normal records); visible when populated (trashed records queried explicitly) |
| 8 | Index on `deleted_at` column | Every GORM query adds `WHERE deleted_at IS NULL`; index prevents full table scans |

---

## Open Questions

_None — all resolved._

---

## Next

Architecture doc → `52-soft-deletes-architecture.md`
