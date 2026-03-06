# ✅ Tasks: Middleware Pipeline

> **Feature**: `08` — Middleware Pipeline
> **Architecture**: [`08-middleware-architecture.md`](08-middleware-architecture.md)
> **Branch**: `feature/08-middleware`
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/21 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [ ] Feature branch created from latest `main`
- [x] Dependent features are merged to `main`
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — Setup

> Remove placeholder, verify project builds clean.

- [ ] **A.1** — Remove `core/middleware/.gitkeep` (will be replaced by real files)
- [ ] 📍 **Checkpoint A** — `go build ./...` succeeds

---

## Phase B — Middleware Registry

> Central registry for naming and grouping middleware.

- [ ] **B.1** — Create `core/middleware/registry.go` with package declaration and map variables
- [ ] **B.2** — Implement `RegisterAlias(name, handler)`
- [ ] **B.3** — Implement `RegisterGroup(name, ...handlers)`
- [ ] **B.4** — Implement `Resolve(name)` with panic on unknown
- [ ] **B.5** — Implement `ResolveGroup(name)` returning nil for unknown
- [ ] **B.6** — Implement `ResetRegistry()` for test cleanup
- [ ] 📍 **Checkpoint B** — Registry compiles, `go vet ./core/middleware/...` clean

---

## Phase C — Recovery Middleware

> Catch panics and return structured 500 responses.

- [ ] **C.1** — Create `core/middleware/recovery.go` with `Recovery()` function
- [ ] **C.2** — Implement defer/recover with `slog.Error` logging
- [ ] **C.3** — Handle already-written responses (check `c.Writer.Written()`)
- [ ] 📍 **Checkpoint C** — Recovery compiles, `go vet` clean

---

## Phase D — RequestID Middleware

> Assign unique IDs to every request for tracing.

- [ ] **D.1** — Create `core/middleware/request_id.go` with `RequestID()` function
- [ ] **D.2** — Implement `generateUUID()` using `crypto/rand`
- [ ] **D.3** — Preserve incoming `X-Request-ID` header if present
- [ ] 📍 **Checkpoint D** — RequestID compiles, `go vet` clean

---

## Phase E — CORS Middleware

> Handle cross-origin requests with configurable options.

- [ ] **E.1** — Create `core/middleware/cors.go` with `CORSConfig` struct and `defaultCORSConfig()`
- [ ] **E.2** — Implement `CORS(configs ...CORSConfig)` with optional custom config
- [ ] **E.3** — Handle preflight `OPTIONS` requests (abort with 204)
- [ ] 📍 **Checkpoint E** — CORS compiles, `go vet` clean

---

## Phase F — ErrorHandler Middleware

> Convert AppError values into JSON responses.

- [ ] **F.1** — Create `core/middleware/error_handler.go` with `ErrorHandler()` function
- [ ] **F.2** — Implement `c.Next()` + check `c.Errors` for `*AppError`
- [ ] **F.3** — Handle generic (non-AppError) errors as 500
- [ ] 📍 **Checkpoint F** — ErrorHandler compiles, imports `core/errors`, `go vet` clean

---

## Phase G — MiddlewareProvider

> Integrate middleware registration with the provider lifecycle.

- [ ] **G.1** — Create `app/providers/middleware_provider.go` with `MiddlewareProvider`
- [ ] **G.2** — Implement `Register()` as no-op
- [ ] **G.3** — Implement `Boot()` — register 4 built-in aliases + default group
- [ ] **G.4** — Update `cmd/main.go` — insert `MiddlewareProvider` before `RouterProvider` (provider #3)
- [ ] 📍 **Checkpoint G** — Provider compiles, `go vet` clean

---

## Phase H — Testing

> Comprehensive test suite for all middleware and registry functionality.

- [ ] **H.1** — Create `core/middleware/middleware_test.go` with all test cases
- [ ] **H.2** — Run `go test ./core/middleware/...` — all tests pass
- [ ] **H.3** — Run `go test ./...` — no regressions across all packages
- [ ] **H.4** — Run `go vet ./...` — no issues
- [ ] 📍 **Checkpoint H** — All tests pass, zero vet warnings

---

## Phase I — Documentation & Cleanup

> Changelog, self-review.

- [ ] **I.1** — Update changelog doc with implementation summary
- [ ] **I.2** — Self-review all diffs — code is clean, idiomatic Go
- [ ] 📍 **Checkpoint I** — Clean code, complete docs, ready to ship

---

## Ship 🚀

- [ ] All phases complete
- [ ] All checkpoints verified
- [ ] Final commit with descriptive message
- [ ] Merge to `main`
- [ ] Push `main`
- [ ] **Keep the feature branch** — do not delete
- [ ] Update project roadmap progress
- [ ] Create review doc → `08-middleware-review.md`
