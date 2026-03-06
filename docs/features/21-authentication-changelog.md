# 📝 Changelog: Authentication

> **Feature**: `21` — Authentication
> **Branch**: `feature/21-authentication`
> **Started**: 2026-03-06
> **Completed**: 2026-03-06

---

## Log

### 2026-03-06

- **[Added]**: `github.com/golang-jwt/jwt/v5` v5.3.1 dependency
- **[Added]**: `core/auth/jwt.go` — `GenerateToken`, `ValidateToken`
- **[Added]**: `core/middleware/auth.go` — `AuthMiddleware()` with Bearer token validation
- **[Changed]**: `app/providers/middleware_provider.go` — registered `"auth"` alias
- **[Changed]**: `database/models/user.go` — added `BeforeCreate` GORM hook for password hashing
- **[Added]**: `core/auth/auth_test.go` — 9 tests (TC-01 to TC-09)
- **[Added]**: 5 tests in `core/middleware/middleware_test.go` (TC-10 to TC-14)
- **[Added]**: 2 tests in `database/models/models_test.go` (TC-15 to TC-16)

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| jwt/v5 version | "latest" | v5.3.1 | Resolved by `go get` |
| No deviations | — | Implementation matches architecture exactly | — |

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
| `"auth"` alias in MiddlewareProvider | Registered alongside existing aliases (recovery, requestid, cors, error_handler) | 2026-03-06 |
