# ✅ Tasks: Database Connection

> **Feature**: `09` — Database Connection
> **Architecture**: [`09-database-architecture.md`](09-database-architecture.md)
> **Branch**: `feature/09-database`
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/14 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [ ] Feature branch created from latest `main`
- [x] Dependent features are merged to `main`
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — Dependencies

> Add GORM and driver packages to the project.

- [ ] **A.1** — Run `go get gorm.io/gorm gorm.io/driver/postgres gorm.io/driver/mysql github.com/glebarez/sqlite`
- [ ] 📍 **Checkpoint A** — `go build ./...` succeeds, `go.mod` lists new dependencies

---

## Phase B — Connection Module

> Implement database config, DSN builder, and connection factory.

- [ ] **B.1** — Expand `database/connection.go`: add `DBConfig` struct + `NewDBConfig()` function
- [ ] **B.2** — Implement `DSN()` method on `DBConfig` (postgres, mysql, sqlite formats)
- [ ] **B.3** — Implement `newDialector()` internal function (switch on driver, return `gorm.Dialector`)
- [ ] **B.4** — Implement `ConnectWithConfig(cfg DBConfig) (*gorm.DB, error)`
- [ ] **B.5** — Implement `Connect() (*gorm.DB, error)` wrapper calling `NewDBConfig()` + `ConnectWithConfig()`
- [ ] 📍 **Checkpoint B** — `go build ./database/...` succeeds, `go vet ./database/...` clean

---

## Phase C — DatabaseProvider & main.go

> Integrate database connection with the provider lifecycle.

- [ ] **C.1** — Create `app/providers/database_provider.go` with `Register()` (Singleton) and `Boot()` (no-op)
- [ ] **C.2** — Update `cmd/main.go` — insert `DatabaseProvider` as provider #3 (Middleware → #4, Router → #5)
- [ ] 📍 **Checkpoint C** — `go build ./...` succeeds, `go vet ./...` clean

---

## Phase D — Testing

> Comprehensive test suite for config, DSN, connection, and provider.

- [ ] **D.1** — Create `database/database_test.go` with config, DSN, and connection tests
- [ ] **D.2** — Add provider tests to `app/providers/providers_test.go` (compile-time check + binding test)
- [ ] **D.3** — Run `go test ./database/...` — all tests pass
- [ ] **D.4** — Run `go test ./...` + `go vet ./...` — full regression, no failures
- [ ] 📍 **Checkpoint D** — All tests pass, zero vet warnings

---

## Phase E — Documentation & Cleanup

> Changelog, self-review.

- [ ] **E.1** — Update changelog doc with implementation summary
- [ ] **E.2** — Self-review all diffs — code is clean, idiomatic Go
- [ ] 📍 **Checkpoint E** — Clean code, complete docs, ready to ship

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
