# ✅ Tasks: Database Seeding

> **Feature**: `13` — Database Seeding
> **Architecture**: [`13-database-seeding-architecture.md`](13-database-seeding-architecture.md)
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/4 phases complete

---

## Phase A — Seeder Interface & Registry

> Create `database/seeders/seeder.go` with interface and registry functions.

- [ ] **A.1** — Create `database/seeders/seeder.go`: `Seeder` interface, `Register()`, `ResetRegistry()`, `RunAll()`, `RunByName()`, `Names()`
- [ ] **A.2** — `go build ./database/seeders/...` clean
- [ ] 📍 **Checkpoint A** — Seeder interface and registry compile

## Phase B — UserSeeder + CLI Command

> Create the sample seeder and CLI command.

- [ ] **B.1** — Create `database/seeders/user_seeder.go`: `UserSeeder` with `init()` registration
- [ ] **B.2** — Create `core/cli/seed.go`: `dbSeedCmd` with `--seeder` flag
- [ ] **B.3** — Update `core/cli/root.go`: add `dbSeedCmd` to `init()`
- [ ] **B.4** — `go build ./...` clean
- [ ] 📍 **Checkpoint B** — `db:seed` command registered, build clean

## Phase C — Testing

> Create tests for seeder system.

- [ ] **C.1** — Create `database/seeders/seeders_test.go` with test cases from testplan
- [ ] **C.2** — `go test ./database/seeders/... -v` — all pass
- [ ] **C.3** — `go test ./... -count=1` — full regression, 0 failures
- [ ] 📍 **Checkpoint C** — All new tests pass, no regressions

## Phase D — Changelog & Self-Review

- [ ] **D.1** — Update `13-database-seeding-changelog.md` with build log and deviations
- [ ] **D.2** — Cross-check: verify code matches architecture doc
- [ ] 📍 **Checkpoint D** — Changelog complete, architecture consistent
