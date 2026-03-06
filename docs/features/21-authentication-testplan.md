# 🧪 Test Plan: Authentication

> **Feature**: `21` — Authentication
> **Total Test Cases**: 16

---

## JWT Package — `core/auth/auth_test.go`

| ID | Test | Input | Expected |
|-----|------|-------|----------|
| TC-01 | GenerateToken returns valid JWT | `userID=1`, `JWT_SECRET=test-secret` | Non-empty string, no error |
| TC-02 | ValidateToken parses valid token | Token from TC-01 | Claims contain `user_id=1`, `exp`, `iat` |
| TC-03 | ValidateToken rejects expired token | Token with `exp` in past | Error returned |
| TC-04 | ValidateToken rejects malformed token | `"not.a.jwt"` | Error returned |
| TC-05 | GenerateToken fails without secret | `JWT_SECRET=""` | Error: "JWT_SECRET is not set" |
| TC-06 | ValidateToken fails without secret | `JWT_SECRET=""` | Error: "JWT_SECRET is not set" |
| TC-07 | ValidateToken rejects wrong secret | Token signed with key A, validated with key B | Error returned |
| TC-08 | GenerateToken respects JWT_EXPIRY | `JWT_EXPIRY=7200` | Token `exp` claim ~7200s from now |
| TC-09 | ValidateToken user_id claim type | Token with `user_id=42` | `claims["user_id"]` is `float64(42)` (JSON number) |

## Auth Middleware — `core/middleware/middleware_test.go`

| ID | Test | Input | Expected |
|-----|------|-------|----------|
| TC-10 | Rejects missing auth header | No Authorization header | 401 + error message |
| TC-11 | Rejects invalid token | `Authorization: Bearer invalid-token` | 401 + error message |
| TC-12 | Sets user_id on valid token | Valid Bearer token | 200 + `user_id` in context |
| TC-13 | Rejects non-Bearer scheme | `Authorization: Basic abc123` | 401 + error message |
| TC-14 | Auth alias is resolvable | `middleware.Resolve("auth")` | Returns non-nil handler |

## User Model Hook — `database/models/models_test.go`

| ID | Test | Input | Expected |
|-----|------|-------|----------|
| TC-15 | BeforeCreate hashes plaintext | `Password: "mypassword"` | Password starts with `$2a$` |
| TC-16 | BeforeCreate skips hashed | `Password: "$2a$10$..."` | Password unchanged |

---

## Acceptance Criteria

1. All 16 tests pass
2. Full regression (`go test ./... -count=1`) — 0 failures
3. `go vet ./...` — clean
4. No new env vars required (JWT_SECRET, JWT_EXPIRY already exist)
