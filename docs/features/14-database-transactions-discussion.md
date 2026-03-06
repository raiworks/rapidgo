# 💬 Discussion: Database Transactions

> **Feature**: `14` — Database Transactions
> **Status**: 🟢 COMPLETE
> **Date**: 2026-03-06

---

## What Are We Building?

A thin transaction helper package in `database/` that wraps GORM's `db.Transaction()` with a framework-consistent API. The helper provides a `WithTransaction` function that accepts a callback, commits on `nil` return, and rolls back on error or panic. A `TransferCredits` example demonstrates the pattern.

## Blueprint References

The blueprint specifies (lines 3004–3026):

1. **Pattern**: `db.Transaction(func(tx *gorm.DB) error { ... })` — GORM's callback-based API
2. **Example**: `TransferCredits(db, fromID, toID, amount)` — atomic credit transfer
3. **Commit/Rollback**: Return `nil` → commit, return error → rollback
4. **Atomic SQL**: Uses `gorm.Expr("credits - ?", amount)` for race-safe arithmetic

## Scope for Feature #14

### In Scope
- `WithTransaction(db, fn)` — framework wrapper over GORM's `db.Transaction()`
- `TransferCredits` — example function demonstrating the pattern
- Tests proving commit, rollback, and error propagation

### Out of Scope (deferred)
- Nested transactions / savepoints — GORM supports these natively; no wrapper needed now
- Manual `Begin/Commit/Rollback` — the callback pattern is safer and simpler
- Transaction middleware — will be relevant when services layer (#29) ships

## Key Design Decisions

### 1. `WithTransaction` Wrapper vs Direct GORM
The blueprint shows raw `db.Transaction()`. We add a `WithTransaction` helper in the `database` package for:
- Framework-consistent API (users import `database.WithTransaction`)
- Single place to add logging/metrics later without changing call sites
- Matches the framework's pattern of thin wrappers over GORM

The wrapper delegates directly to `db.Transaction()` — zero overhead, purely organizational.

### 2. TransferCredits as Example
The blueprint's `TransferCredits` function serves as both the primary example and a testable function. It lives in the `database` package alongside `WithTransaction`. Uses `gorm.Expr()` for atomic SQL per the blueprint's security guidance.

### 3. No Separate Sub-Package
Transactions are a core database concern. `WithTransaction` lives in the existing `database` package (`database/transaction.go`) rather than a sub-package — it's tightly coupled to `*gorm.DB` and there's no registry/lifecycle to justify a separate package.

## Dependencies

| Dependency | Status | Notes |
|---|---|---|
| Feature #09 — Database Connection | ✅ Done | Provides `*gorm.DB` |
| Feature #11 — Models (GORM) | ✅ Done | Provides User model (for TransferCredits example) |

## Discussion Complete ✅
