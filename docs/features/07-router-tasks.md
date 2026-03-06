# ✅ Tasks: Router & Routing

> **Feature**: `07` — Router & Routing
> **Architecture**: [`07-router-architecture.md`](07-router-architecture.md)
> **Branch**: `feature/07-router`
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/22 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [ ] Feature branch created from latest `main`
- [x] Dependent features are merged to `main`
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — Dependencies & Setup

> Add Gin dependency, verify project compiles.

- [ ] **A.1** — Run `go get github.com/gin-gonic/gin` to add Gin dependency
- [ ] **A.2** — Verify `go mod tidy` completes without errors
- [ ] **A.3** — Remove `core/router/.gitkeep` (will be replaced by real files)
- [ ] 📍 **Checkpoint A** — `go build ./...` succeeds with Gin dependency

---

## Phase B — Router Core

> Router struct, HTTP method helpers, Gin mode integration.

- [ ] **B.1** — Create `core/router/router.go` with `Router` struct wrapping `*gin.Engine`
- [ ] **B.2** — Implement `New()` — creates Gin engine, sets mode from `APP_ENV`
- [ ] **B.3** — Implement `setGinMode()` — maps APP_ENV to gin.SetMode
- [ ] **B.4** — Implement HTTP method helpers: `Get`, `Post`, `Put`, `Delete`, `Patch`, `Options`
- [ ] **B.5** — Implement `Engine()`, `ServeHTTP()`, `Use()`, `Run()`
- [ ] **B.6** — Implement `Group()` — returns `*RouteGroup`
- [ ] 📍 **Checkpoint B** — Router compiles, `go vet ./core/router/...` clean

---

## Phase C — Route Groups

> RouteGroup struct with same method set as Router.

- [ ] **C.1** — Create `core/router/group.go` with `RouteGroup` struct wrapping `*gin.RouterGroup`
- [ ] **C.2** — Implement HTTP method helpers on `RouteGroup`: `Get`, `Post`, `Put`, `Delete`, `Patch`, `Options`
- [ ] **C.3** — Implement `Group()` on `RouteGroup` for nesting
- [ ] **C.4** — Implement `Use()` on `RouteGroup`
- [ ] 📍 **Checkpoint C** — Groups compile, nested groups work

---

## Phase D — Resource Routes

> ResourceController interface, Resource() and APIResource() on both Router and RouteGroup.

- [ ] **D.1** — Create `core/router/resource.go` with `ResourceController` interface (7 methods)
- [ ] **D.2** — Implement `Resource()` on `Router` — registers 7 RESTful routes
- [ ] **D.3** — Implement `APIResource()` on `Router` — registers 5 RESTful routes
- [ ] **D.4** — Implement `Resource()` on `RouteGroup` — registers 7 RESTful routes on group
- [ ] **D.5** — Implement `APIResource()` on `RouteGroup` — registers 5 RESTful routes on group
- [ ] 📍 **Checkpoint D** — Resource routes compile, correct paths registered

---

## Phase E — Named Routes

> Named route registry with URL generation.

- [ ] **E.1** — Create `core/router/named.go` with named route registry (sync.RWMutex)
- [ ] **E.2** — Implement `Name(name, pattern)`
- [ ] **E.3** — Implement `Route(name, params...) string` with parameter substitution
- [ ] **E.4** — Implement `ResetNamedRoutes()` for test cleanup
- [ ] 📍 **Checkpoint E** — Named routes compile, URL generation works

---

## Phase F — Provider & Route Files

> RouterProvider, updated routes/web.go and routes/api.go, updated main.go.

- [ ] **F.1** — Create `app/providers/router_provider.go` with `RouterProvider`
- [ ] **F.2** — Update `routes/web.go` — add `RegisterWeb(r *router.Router)`
- [ ] **F.3** — Update `routes/api.go` — add `RegisterAPI(r *router.Router)`
- [ ] **F.4** — Update `cmd/main.go` — register `RouterProvider`, start HTTP server
- [ ] 📍 **Checkpoint F** — `go run cmd/main.go` starts server, responds on port

---

## Phase G — Testing

> Comprehensive test suite for all router functionality.

- [ ] **G.1** — Create `core/router/router_test.go` with all test cases
- [ ] **G.2** — Run `go test ./core/router/...` — all tests pass
- [ ] **G.3** — Run `go test ./...` — no regressions across all packages
- [ ] **G.4** — Run `go vet ./...` — no issues
- [ ] 📍 **Checkpoint G** — All tests pass, zero vet warnings

---

## Phase H — Documentation & Cleanup

> Changelog, roadmap, self-review.

- [ ] **H.1** — Update changelog doc with implementation summary
- [ ] **H.2** — Self-review all diffs — code is clean, idiomatic Go
- [ ] 📍 **Checkpoint H** — Clean code, complete docs, ready to ship

---

## Ship 🚀

- [ ] All phases complete
- [ ] All checkpoints verified
- [ ] Final commit with descriptive message
- [ ] Merge to `main`
- [ ] Push `main`
- [ ] **Keep the feature branch** — do not delete
- [ ] Update project roadmap progress
- [ ] Create review doc → `07-router-review.md`
