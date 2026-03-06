# 📋 Tasks: Authentication

> **Feature**: `21` — Authentication
> **Branch**: `feature/21-authentication`
> **Status**: NOT STARTED

---

## Phase 1 — JWT Package

- [ ] Add `github.com/golang-jwt/jwt/v5` dependency via `go get`
- [ ] Create `core/auth/jwt.go` with `GenerateToken` function
- [ ] Create `core/auth/jwt.go` with `ValidateToken` function (signing method check)
- [ ] Verify: `JWT_SECRET` and `JWT_EXPIRY` env vars already exist in `.env`

**Checkpoint**: `GenerateToken` and `ValidateToken` compile and work with env vars.

## Phase 2 — Auth Middleware

- [ ] Create `core/middleware/auth.go` with `AuthMiddleware()` function
- [ ] Handles missing/malformed Authorization header (401)
- [ ] Extracts Bearer token and calls `auth.ValidateToken()`
- [ ] Sets `user_id` in Gin context on success

- [ ] Register `"auth"` alias: `middleware.RegisterAlias("auth", AuthMiddleware())`

**Checkpoint**: AuthMiddleware compiles with correct imports. Alias registered.

## Phase 3 — User Model Hook

- [ ] Add `BeforeCreate` GORM hook to `database/models/user.go`
- [ ] Hook auto-hashes password via `helpers.HashPassword`
- [ ] Hook skips already-hashed passwords (`$2a$`/`$2b$` prefix check)

**Checkpoint**: User model compiles with new import.

## Phase 4 — Tests

- [ ] Create `core/auth/auth_test.go`
  - [ ] TC-01: GenerateToken returns valid JWT string
  - [ ] TC-02: ValidateToken parses valid token and returns claims
  - [ ] TC-03: ValidateToken rejects expired token
  - [ ] TC-04: ValidateToken rejects malformed token
  - [ ] TC-05: GenerateToken fails when JWT_SECRET is empty
  - [ ] TC-06: ValidateToken fails when JWT_SECRET is empty
  - [ ] TC-07: ValidateToken rejects token signed with wrong secret
  - [ ] TC-08: GenerateToken respects JWT_EXPIRY env var
  - [ ] TC-09: ValidateToken returns user_id claim as expected type
- [ ] Add auth middleware tests to `core/middleware/middleware_test.go`
  - [ ] TC-10: AuthMiddleware rejects request without Authorization header
  - [ ] TC-11: AuthMiddleware rejects request with invalid token
  - [ ] TC-12: AuthMiddleware sets user_id on valid token
  - [ ] TC-13: AuthMiddleware rejects non-Bearer auth scheme
  - [ ] TC-14: Auth alias is resolvable
- [ ] Add BeforeCreate test to `database/models/models_test.go`
  - [ ] TC-14: Auth alias is resolvable
- [ ] Add BeforeCreate test to `database/models/models_test.go`
  - [ ] TC-15: BeforeCreate hashes plaintext password
  - [ ] TC-16: BeforeCreate skips already-hashed password
- [ ] Run full `go test ./... -count=1` — all pass

**Checkpoint**: All 15 new tests pass. Full regression passes.

## Phase 5 — Finalize

- [ ] Update changelog with all changes
- [ ] Run `go vet ./...` — clean
- [ ] Commit and push to feature branch
