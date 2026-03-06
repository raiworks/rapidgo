# 📋 Review: Services Layer

> **Feature**: `18` — Services Layer
> **Branch**: `feature/18-services-layer`
> **Merged**: 2026-03-06 → `main` (`1a7f436`)
> **Status**: ✅ SHIPPED

---

## Summary

Added the services layer pattern with `UserService` in `app/services/`. The service encapsulates CRUD business logic for users, taking a `*gorm.DB` dependency and operating on `models.User`. This establishes the Controller → Service → Model delegation pattern for the framework.

## Files Changed

| File | Action | Purpose |
|---|---|---|
| `app/services/user_service.go` | Created | UserService with GetByID, Create, Update, Delete |
| `app/services/user_service_test.go` | Created | 8 tests with SQLite in-memory |
| `docs/features/18-services-layer-changelog.md` | Updated | Build log |
| `docs/project-roadmap.md` | Updated | #18 → ✅ |

## Test Results

- **New tests**: 8 (TC-01 through TC-08)
- **Total tests**: 185
- **Failures**: 0

## Deviations from Architecture

None. Implementation matches architecture doc exactly.

## Key Design Points

- `UserService.Create` checks duplicate email before insertion
- `Password` stored as-is with comment — hashing belongs to Feature #22
- `Update` uses `map[string]interface{}` for flexible partial updates
- Services are standalone — no container registration, instantiated directly
