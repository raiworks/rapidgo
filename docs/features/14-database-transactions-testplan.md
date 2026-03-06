# 🧪 Test Plan: Database Transactions

> **Feature**: `14` — Database Transactions
> **Architecture**: [`14-database-transactions-architecture.md`](14-database-transactions-architecture.md)
> **Status**: ⬜ NOT RUN
> **Result**: ⬜ NOT RUN

---

## Test File

`database/transaction_test.go`

All tests use SQLite `:memory:` via GORM — no external database required.

---

## Test Cases

### TC-01: `TestWithTransaction_Commit`
**What**: `WithTransaction` commits when callback returns `nil`.
**How**: Create a users table with a row. Call `WithTransaction` with a callback that updates the row. Verify the update persisted.
**Pass**: Row updated in database.

### TC-02: `TestWithTransaction_Rollback`
**What**: `WithTransaction` rolls back when callback returns an error.
**How**: Create a users table with a row. Call `WithTransaction` with a callback that updates the row then returns an error. Verify the row is unchanged.
**Pass**: Row unchanged, error returned to caller.

### TC-03: `TestWithTransaction_PanicRollback`
**What**: `WithTransaction` rolls back on panic inside callback.
**How**: Create a users table with a row. Call `WithTransaction` with a callback that updates the row then panics. Recover the panic. Verify the row is unchanged.
**Pass**: Row unchanged, panic recovered.

### TC-04: `TestWithTransaction_ErrorPropagation`
**What**: The error from the callback is returned to the caller.
**How**: Call `WithTransaction` with a callback that returns a specific sentinel error. Verify the returned error wraps/matches the sentinel.
**Pass**: Caller receives the callback's error.

### TC-05: `TestTransferCredits_Success`
**What**: `TransferCredits` atomically moves credits between users.
**How**: Create users table with `credits` column. Insert user A (credits=200) and user B (credits=50). Call `TransferCredits(db, A.ID, B.ID, 100)`. Verify A has 100, B has 150.
**Pass**: Credits transferred correctly.

### TC-06: `TestTransferCredits_SourceNotFound`
**What**: `TransferCredits` fails if source user doesn't exist.
**How**: Create users table. Insert only user B. Call `TransferCredits(db, 999, B.ID, 100)`. Verify error returned and B's credits unchanged.
**Pass**: Error returned, no credits moved.

### TC-07: `TestTransferCredits_DestNotFound`
**What**: `TransferCredits` fails if destination user doesn't exist.
**How**: Create users table with user A (credits=200). Call `TransferCredits(db, A.ID, 999, 100)`. Verify error returned and A's credits unchanged.
**Pass**: Error returned, no credits moved.

---

## Test Summary

| ID | Test Name | Type | Scope |
|---|---|---|---|
| TC-01 | `TestWithTransaction_Commit` | Unit | Commit path |
| TC-02 | `TestWithTransaction_Rollback` | Unit | Rollback path |
| TC-03 | `TestWithTransaction_PanicRollback` | Unit | Panic recovery |
| TC-04 | `TestWithTransaction_ErrorPropagation` | Unit | Error forwarding |
| TC-05 | `TestTransferCredits_Success` | Integration | Happy path |
| TC-06 | `TestTransferCredits_SourceNotFound` | Integration | Source missing |
| TC-07 | `TestTransferCredits_DestNotFound` | Integration | Dest missing |

**Total**: 7 test cases
**Expected new test count**: 148 + 7 = 155
