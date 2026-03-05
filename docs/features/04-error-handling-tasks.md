# ✅ Tasks: Error Handling

> **Feature**: `04` — Error Handling
> **Architecture**: [`04-error-handling-architecture.md`](04-error-handling-architecture.md)
> **Branch**: `feature/04-error-handling`
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

## Phase A — Core Error Types

> AppError struct, error/Unwrap interface compliance.

- [ ] **A.1** — Create `core/errors/errors.go` with package declaration and imports
- [ ] **A.2** — Define `AppError` struct with `Code int`, `Message string`, `Err error` fields
- [ ] **A.3** — Implement `Error() string` method returning `Message`
- [ ] **A.4** — Implement `Unwrap() error` method returning `Err` (handles nil)
- [ ] 📍 **Checkpoint A** — `AppError` compiles, implements `error` interface, `go vet` clean

---

## Phase B — Constructor Helpers

> Helper functions for common HTTP error scenarios.

- [ ] **B.1** — Implement `NotFound(message string) *AppError` — 404
- [ ] **B.2** — Implement `BadRequest(message string) *AppError` — 400
- [ ] **B.3** — Implement `Internal(err error) *AppError` — 500 with wrapped error
- [ ] **B.4** — Implement `Unauthorized(message string) *AppError` — 401
- [ ] **B.5** — Implement `Forbidden(message string) *AppError` — 403
- [ ] **B.6** — Implement `Conflict(message string) *AppError` — 409
- [ ] **B.7** — Implement `Unprocessable(message string) *AppError` — 422
- [ ] 📍 **Checkpoint B** — All 7 constructors compile, return correct codes, `go vet` clean

---

## Phase C — Error Response Helper

> Debug-aware response formatting.

- [ ] **C.1** — Implement `ErrorResponse() map[string]any` on `*AppError` — returns `{"success": false, "error": message}`, adds `"internal"` key only when `config.IsDebug() && e.Err != nil`
- [ ] 📍 **Checkpoint C** — `ErrorResponse()` compiles, `go vet` clean

---

## Phase D — Testing

> Execute the test plan, verify all acceptance criteria.

- [ ] **D.1** — Create `core/errors/errors_test.go` with all test cases from test plan
- [ ] **D.2** — Run `go test ./core/errors/...` — all tests pass
- [ ] **D.3** — Run `go vet ./...` — no issues
- [ ] 📍 **Checkpoint D** — All test cases pass, zero vet warnings

---

## Phase E — Documentation & Cleanup

> Changelog, roadmap, self-review.

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
- [ ] Create review doc → `04-error-handling-review.md`
