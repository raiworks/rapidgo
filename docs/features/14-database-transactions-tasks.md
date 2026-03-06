# ✅ Tasks: Database Transactions

> **Feature**: `14` — Database Transactions
> **Architecture**: [`14-database-transactions-architecture.md`](14-database-transactions-architecture.md)
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/3 phases complete

---

## Phase A — Transaction Helper & Example

> Create `database/transaction.go` and `database/transaction_example.go`.

- [ ] **A.1** — Create `database/transaction.go`: `TxFunc` type, `WithTransaction()` function
- [ ] **A.2** — Create `database/transaction_example.go`: `TransferCredits()` function
- [ ] **A.3** — `go build ./database/...` clean
- [ ] **A.4** — `go vet ./database/...` clean
- [ ] 📍 **Checkpoint A** — Transaction helper and example compile

## Phase B — Testing

> Create tests for transaction system.

- [ ] **B.1** — Create `database/transaction_test.go` with test cases from testplan
- [ ] **B.2** — `go test ./database/... -v` — all pass
- [ ] **B.3** — `go test ./... -count=1` — full regression, 0 failures
- [ ] 📍 **Checkpoint B** — All new tests pass, no regressions

## Phase C — Changelog & Self-Review

- [ ] **C.1** — Update `14-database-transactions-changelog.md` with build log and deviations
- [ ] **C.2** — Cross-check: verify code matches architecture doc
- [ ] 📍 **Checkpoint C** — Changelog complete, architecture consistent
