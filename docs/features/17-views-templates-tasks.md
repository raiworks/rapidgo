# ✅ Tasks: Views & Templates

> **Feature**: `17` — Views & Templates
> **Architecture**: [`17-views-templates-architecture.md`](17-views-templates-architecture.md)
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/3 phases complete

---

## Phase A — View Engine & Static Serving

> Add template and static methods to Router, create DefaultFuncMap, wire in provider.

- [ ] **A.1** — Create `core/router/view.go`: `DefaultFuncMap()` returning `template.FuncMap{"route": Route}`
- [ ] **A.2** — Add `SetFuncMap`, `LoadTemplates`, `Static`, `StaticFile` methods to `core/router/router.go`
- [ ] **A.3** — Update `app/providers/router_provider.go`: add `SetFuncMap`, `LoadTemplates`, `Static` calls in Boot
- [ ] **A.4** — Create `resources/views/home.html` sample template
- [ ] **A.5** — Update `http/controllers/home_controller.go`: `c.HTML()` instead of `c.JSON()`
- [ ] **A.6** — `go build ./...` clean
- [ ] **A.7** — `go vet ./...` clean
- [ ] 📍 **Checkpoint A** — Framework compiles with view engine support

## Phase B — Testing

> Update existing tests and add new tests for view engine features.

- [ ] **B.1** — Update `http/controllers/controllers_test.go`: Home test asserts HTML response
- [ ] **B.2** — Create `core/router/view_test.go`: test DefaultFuncMap + integration
- [ ] **B.3** — `go test ./core/router/... -v` — all pass
- [ ] **B.4** — `go test ./http/controllers/... -v` — all pass
- [ ] **B.5** — `go test ./... -count=1` — full regression, 0 failures
- [ ] 📍 **Checkpoint B** — All tests pass, no regressions

## Phase C — Changelog & Self-Review

- [ ] **C.1** — Update `17-views-templates-changelog.md` with build log and deviations
- [ ] **C.2** — Cross-check: verify code matches architecture doc
- [ ] 📍 **Checkpoint C** — Changelog complete, architecture consistent
