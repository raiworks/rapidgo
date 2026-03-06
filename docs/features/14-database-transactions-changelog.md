# 📝 Changelog: Database Transactions

> **Feature**: `14` — Database Transactions
> **Branch**: `feature/14-database-transactions`
> **Started**: —
> **Completed**: —

---

## Log

- **Phase A** — Created `database/transaction.go` (`TxFunc`, `WithTransaction`) and `database/transaction_example.go` (`TransferCredits`). Build + vet clean.
- **Phase B** — Created `database/transaction_test.go` with 7 tests (TC-01 through TC-07). All pass. Full regression: 155 tests, 0 failures.
- **Phase C** — Cross-check passed. One minor deviation noted.

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| `TransferCredits` existence check | `tx.First(&struct{ID uint}{}, id)` | `tx.Table("users").First(&struct{ID uint}{}, id)` | Explicit table targeting for clarity when not using a model |

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
| `testUser` local struct in tests | Tests need a users table with `credits` column, but can't import `models.User` (no `credits` field); local struct keeps package independent | 2026-03-06 |
