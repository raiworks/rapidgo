# 📋 Tasks: CSRF Protection

> **Feature**: `24` — CSRF Protection
> **Branch**: `feature/24-csrf-protection`
> **Status**: NOT STARTED

---

## Phase 1 — Implementation

- [ ] Create `core/middleware/csrf.go` with `CSRFMiddleware()` function
- [ ] Add `"csrf"` alias in `app/providers/middleware_provider.go`
- [ ] Verify: `go build ./core/middleware/...` compiles

**Checkpoint**: CSRF middleware compiles, alias registered.

## Phase 2 — Tests

- [ ] Add CSRF tests to `core/middleware/middleware_test.go`
  - [ ] TC-01: GET request — token generated and set in context
  - [ ] TC-02: POST with valid form token — passes
  - [ ] TC-03: POST with valid header token — passes
  - [ ] TC-04: POST with missing token — 403
  - [ ] TC-05: POST with wrong token — 403
  - [ ] TC-06: HEAD request skips validation
  - [ ] TC-07: OPTIONS request skips validation
  - [ ] TC-08: PUT with valid token — passes
  - [ ] TC-09: DELETE with missing token — 403
  - [ ] TC-10: Token persists across requests (same session)
  - [ ] TC-11: "csrf" alias is resolvable
- [ ] Run full `go test ./... -count=1` — all pass

**Checkpoint**: All 11 tests pass. Full regression green.

## Phase 3 — Finalize

- [ ] Update changelog
- [ ] Run `go vet ./...` — clean
- [ ] Commit and push
