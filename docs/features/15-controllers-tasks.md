# ✅ Tasks: Controllers

> **Feature**: `15` — Controllers
> **Architecture**: [`15-controllers-architecture.md`](15-controllers-architecture.md)
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/3 phases complete

---

## Phase A — Controllers + Route Registration

> Create controller files and wire up routes.

- [ ] **A.1** — Create `http/controllers/home_controller.go`: `Home()` function
- [ ] **A.2** — Create `http/controllers/post_controller.go`: `PostController` struct with 7 methods
- [ ] **A.3** — Update `routes/web.go`: register `GET /` → `controllers.Home`
- [ ] **A.4** — Update `routes/api.go`: register `APIResource("/posts", &PostController{})`
- [ ] **A.5** — `go build ./...` clean
- [ ] **A.6** — `go vet ./...` clean
- [ ] 📍 **Checkpoint A** — Controllers and routes compile

## Phase B — Testing

> Create tests for controllers.

- [ ] **B.1** — Create `http/controllers/controllers_test.go` with test cases from testplan
- [ ] **B.2** — `go test ./http/controllers/... -v` — all pass
- [ ] **B.3** — `go test ./... -count=1` — full regression, 0 failures
- [ ] 📍 **Checkpoint B** — All new tests pass, no regressions

## Phase C — Changelog & Self-Review

- [ ] **C.1** — Update `15-controllers-changelog.md` with build log and deviations
- [ ] **C.2** — Cross-check: verify code matches architecture doc
- [ ] 📍 **Checkpoint C** — Changelog complete, architecture consistent
