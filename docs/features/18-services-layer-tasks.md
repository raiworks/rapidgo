# ✅ Tasks: Services Layer

> **Feature**: `18` — Services Layer
> **Architecture**: [`18-services-layer-architecture.md`](18-services-layer-architecture.md)
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/3 phases complete

---

## Phase A — UserService

> Create the UserService with CRUD methods.

- [ ] **A.1** — Create `app/services/user_service.go`: `UserService`, `NewUserService`, `GetByID`, `Create`, `Update`, `Delete`
- [ ] **A.2** — `go build ./app/services/...` clean
- [ ] **A.3** — `go vet ./app/services/...` clean
- [ ] 📍 **Checkpoint A** — UserService compiles

## Phase B — Testing

> Create tests for all UserService methods using SQLite in-memory.

- [ ] **B.1** — Create `app/services/user_service_test.go` with test cases from testplan
- [ ] **B.2** — `go test ./app/services/... -v` — all pass
- [ ] **B.3** — `go test ./... -count=1` — full regression, 0 failures
- [ ] 📍 **Checkpoint B** — All new tests pass, no regressions

## Phase C — Changelog & Self-Review

- [ ] **C.1** — Update `18-services-layer-changelog.md` with build log and deviations
- [ ] **C.2** — Cross-check: verify code matches architecture doc
- [ ] 📍 **Checkpoint C** — Changelog complete, architecture consistent
