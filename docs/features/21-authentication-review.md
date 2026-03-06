# 📋 Review: Authentication

> **Feature**: `21` — Authentication
> **Branch**: `feature/21-authentication`
> **Merged**: 2026-03-06
> **Commit**: `ca4f7a6` (impl)

---

## Summary

Feature #21 adds JWT authentication infrastructure to the framework: a `core/auth/` package with token generation/validation, an `AuthMiddleware` for route protection, and a `BeforeCreate` GORM hook on the User model for automatic password hashing.

## Files Changed

| File | Type | Description |
|---|---|---|
| `core/auth/jwt.go` | Created | `GenerateToken`, `ValidateToken` — HMAC-SHA256 JWT with configurable expiry |
| `core/middleware/auth.go` | Created | `AuthMiddleware` — Bearer token extraction, validation, sets `user_id` |
| `core/auth/auth_test.go` | Created | 9 tests for JWT functions |
| `database/models/user.go` | Modified | Added `BeforeCreate` hook — auto-hashes password via bcrypt |
| `app/providers/middleware_provider.go` | Modified | Registered `"auth"` alias |
| `core/middleware/middleware_test.go` | Modified | +5 auth middleware tests |
| `database/models/models_test.go` | Modified | +2 BeforeCreate hook tests |
| `go.mod` / `go.sum` | Modified | Added `golang-jwt/jwt/v5` v5.3.1 |

## Dependencies Added

| Package | Version | Purpose |
|---|---|---|
| `github.com/golang-jwt/jwt/v5` | v5.3.1 | JWT token creation and parsing |

## Test Results

- **New tests**: 16 (9 auth + 5 middleware + 2 model)
- **Total tests**: 253 — all pass
- **`go vet`**: clean

## Architecture Compliance

Implementation matches architecture document exactly. No deviations.

## Security Hardening

1. Signing method validation prevents `none` algorithm attack
2. Empty `JWT_SECRET` returns explicit error
3. Bcrypt double-hash prevention via `$2a$`/`$2b$` prefix check

## Status: ✅ SHIPPED
