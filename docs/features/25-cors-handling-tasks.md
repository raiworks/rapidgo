# 📋 Tasks: CORS Handling

> **Feature**: `25` — CORS Handling
> **Branch**: `feature/25-cors-handling`
> **Status**: NOT STARTED

---

## Phase 1 — Implementation

- [ ] Add `AllowCredentials` and `ExposeHeaders` fields to `CORSConfig`
- [ ] Update `defaultCORSConfig()` — read `CORS_ALLOWED_ORIGINS` env, add `X-CSRF-Token`, expose headers, credentials
- [ ] Update `CORS()` handler — emit `Access-Control-Allow-Credentials` and `Access-Control-Expose-Headers` headers
- [ ] Add `os` import
- [ ] Verify: `go build ./core/middleware/...` compiles

**Checkpoint**: Enhanced CORS compiles, existing tests still pass.

## Phase 2 — Tests

- [ ] Add new CORS tests to `core/middleware/middleware_test.go`
  - [ ] TC-01: Default includes X-CSRF-Token in allowed headers
  - [ ] TC-02: AllowCredentials header set to "true" by default
  - [ ] TC-03: ExposeHeaders set by default
  - [ ] TC-04: Env-based CORS_ALLOWED_ORIGINS overrides default
  - [ ] TC-05: Custom config can disable credentials
  - [ ] TC-06: Preflight still returns 204 with new headers
- [ ] Existing CORS tests still pass
- [ ] Run full `go test ./... -count=1` — all pass

**Checkpoint**: All new + existing tests pass. Full regression green.

## Phase 3 — Finalize

- [ ] Update changelog
- [ ] Run `go vet ./...` — clean
- [ ] Commit and push
