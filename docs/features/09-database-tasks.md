# ✅ Tasks: Database Connection

> **Feature**: `09` — Database Connection
> **Architecture**: [`09-database-architecture.md`](09-database-architecture.md)
> **Branch**: `feature/09-database`
> **Status**: � COMPLETE
> **Progress**: 15/15 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [x] Feature branch created from latest `main`
- [x] Dependent features are merged to `main`
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — Dependencies

> Add GORM and driver packages to the project.

- [x] **A.1** — Run `go get gorm.io/gorm gorm.io/driver/postgres gorm.io/driver/mysql github.com/glebarez/sqlite`
- [x] 📍 **Checkpoint A** — `go build ./...` succeeds, `go.mod` lists new dependencies

---

## Phase B — Connection Module

> Implement database config, DSN builder, and connection factory.

- [x] **B.1** — Expand `database/connection.go`: add `DBConfig` struct + `NewDBConfig()` function
- [x] **B.2** — Implement `DSN()` method on `DBConfig` (postgres, mysql, sqlite formats)
- [x] **B.3** — Implement `newDialector()` internal function (switch on driver, return `gorm.Dialector`)
- [x] **B.4** — Implement `ConnectWithConfig(cfg DBConfig) (*gorm.DB, error)`
- [x] **B.5** — Implement `Connect() (*gorm.DB, error)` wrapper calling `NewDBConfig()` + `ConnectWithConfig()`
- [x] **B.6** — Update `.env` — add commented-out pool tuning variables (`DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME`, `DB_CONN_MAX_IDLE_TIME`)
- [x] 📍 **Checkpoint B** — `go build ./database/...` succeeds, `go vet ./database/...` clean

---

## Phase C — DatabaseProvider & main.go

> Integrate database connection with the provider lifecycle.

- [x] **C.1** — Create `app/providers/database_provider.go` with `Register()` (Singleton) and `Boot()` (no-op)
- [x] **C.2** — Update `cmd/main.go` — insert `DatabaseProvider` as provider #3 (Middleware → #4, Router → #5)
- [x] 📍 **Checkpoint C** — `go build ./...` succeeds, `go vet ./...` clean

---

## Phase D — Testing

> Comprehensive test suite for config, DSN, connection, and provider.

- [x] **D.1** — Create `database/database_test.go` with config, DSN, and connection tests
- [x] **D.2** — Add provider tests to `app/providers/providers_test.go` (compile-time check + binding test)
- [x] **D.3** — Run `go test ./database/...` — all tests pass
- [x] **D.4** — Run `go test ./...` + `go vet ./...` — full regression, no failures
- [x] 📍 **Checkpoint D** — All tests pass, zero vet warnings

---

## Phase E — Documentation & Cleanup

> Changelog, self-review.

- [x] **E.1** — Update changelog doc with implementation summary
- [x] **E.2** — Self-review all diffs — code is clean, idiomatic Go
- [x] 📍 **Checkpoint E** — Clean code, complete docs, ready to ship

---

## Ship 🚀

- [ ] All phases complete
- [ ] All checkpoints verified
- [ ] Final commit with descriptive message
- [ ] Merge to `main`
- [ ] Push `main`
- [ ] **Keep the feature branch** — do not delete
- [ ] Update project roadmap progress
- [ ] Create review doc → `09-database-review.md`
