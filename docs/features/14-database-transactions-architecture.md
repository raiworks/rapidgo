# 🏗️ Architecture: Database Transactions

> **Feature**: `14` — Database Transactions
> **Discussion**: [`14-database-transactions-discussion.md`](14-database-transactions-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

Feature #14 adds a `WithTransaction` helper function to the `database` package and a `TransferCredits` example function. The helper wraps GORM's `db.Transaction()` for framework consistency while the example demonstrates atomic multi-step operations with `gorm.Expr()`.

## File Structure

```
database/
├── connection.go        # (existing) DB connection management
├── transaction.go       # WithTransaction helper
├── transaction_example.go  # TransferCredits example
└── database_test.go     # (existing) connection tests

database/
└── transaction_test.go  # Tests for transactions
```

### Files Created (2)
| File | Package | Lines (est.) |
|---|---|---|
| `database/transaction.go` | `database` | ~15 |
| `database/transaction_example.go` | `database` | ~30 |

### Files Modified (0)
No existing files need modification.

---

## Component Design

### Transaction Helper (`database/transaction.go`)

**Responsibility**: Framework-consistent wrapper over GORM's transaction API.
**Package**: `database`

```go
package database

import "gorm.io/gorm"

// TxFunc is the callback signature for transactional operations.
// Return nil to commit, return an error to rollback.
type TxFunc func(tx *gorm.DB) error

// WithTransaction wraps fn in a database transaction.
// If fn returns nil the transaction commits.
// If fn returns an error or panics the transaction rolls back.
func WithTransaction(db *gorm.DB, fn TxFunc) error {
	return db.Transaction(fn)
}
```

**Design notes**:
- `TxFunc` type alias makes signatures cleaner in calling code
- Direct delegation to `db.Transaction()` — zero overhead
- GORM handles panic recovery and rollback internally
- Single choke-point for future enhancements (logging, metrics, tracing)

### TransferCredits Example (`database/transaction_example.go`)

**Responsibility**: Demonstrate the transaction pattern with an atomic credit transfer.
**Package**: `database`

```go
package database

import (
	"gorm.io/gorm"
)

// TransferCredits atomically moves credits from one user to another.
// Both users must exist and the operation uses gorm.Expr for race-safe SQL.
func TransferCredits(db *gorm.DB, fromID, toID uint, amount int) error {
	return WithTransaction(db, func(tx *gorm.DB) error {
		// Verify source user exists
		if err := tx.First(&struct{ ID uint }{}, fromID).Error; err != nil {
			return err
		}
		// Verify destination user exists
		if err := tx.First(&struct{ ID uint }{}, toID).Error; err != nil {
			return err
		}
		// Deduct from source
		if err := tx.Table("users").Where("id = ?", fromID).
			Update("credits", gorm.Expr("credits - ?", amount)).Error; err != nil {
			return err
		}
		// Add to destination
		if err := tx.Table("users").Where("id = ?", toID).
			Update("credits", gorm.Expr("credits + ?", amount)).Error; err != nil {
			return err
		}
		return nil
	})
}
```

**Design notes**:
- Uses `WithTransaction` (not raw `db.Transaction`) — dogfoods the helper
- `gorm.Expr()` for atomic `credits - ?` / `credits + ?` — prevents race conditions per blueprint
- Uses `tx.Table("users")` and `tx.First` with anonymous struct for existence check — avoids hard dependency on `models.User` (the `database` package should not import `database/models` to avoid circular concerns)
- The blueprint shows `tx.First(&from, fromID)` with `models.User` — we adapt slightly to avoid the import cycle while preserving the same SQL semantics

---

## Data Flow

### `WithTransaction` — Success Path
```
caller → WithTransaction(db, fn)
       → db.Transaction(fn)
           → BEGIN
           → fn(tx) returns nil
           → COMMIT
       → returns nil
```

### `WithTransaction` — Error Path
```
caller → WithTransaction(db, fn)
       → db.Transaction(fn)
           → BEGIN
           → fn(tx) returns error
           → ROLLBACK
       → returns error
```

### `TransferCredits` — Happy Path
```
TransferCredits(db, 1, 2, 100)
  → WithTransaction(db, fn)
      → BEGIN
      → SELECT * FROM users WHERE id = 1  (verify exists)
      → SELECT * FROM users WHERE id = 2  (verify exists)
      → UPDATE users SET credits = credits - 100 WHERE id = 1
      → UPDATE users SET credits = credits + 100 WHERE id = 2
      → COMMIT
```

---

## Constraints & Invariants

1. `WithTransaction` delegates to `db.Transaction()` — **no custom BEGIN/COMMIT logic**
2. GORM handles panic recovery inside `db.Transaction()` automatically
3. `TransferCredits` uses `gorm.Expr()` — **never string interpolation for SQL arithmetic**
4. The `database` package **does not import `database/models`** — avoids import direction issues
5. Transaction callback **must use `tx`**, never the outer `db`
