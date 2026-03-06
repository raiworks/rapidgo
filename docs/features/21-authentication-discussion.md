# 💬 Discussion: Authentication

> **Feature**: `21` — Authentication
> **Status**: COMPLETE
> **Date**: 2026-03-06

---

## 1. What Are We Building?

JWT and session-based authentication for the RGo framework. The blueprint defines two strategies:

- **JWT (stateless)** — for REST APIs. Uses `golang-jwt/jwt/v5` with HMAC-SHA256 signing.
- **Session-based** — for SSR web apps. Uses the existing session manager (#20) to store user state.

This feature adds:
1. A `core/auth/` package with `GenerateToken` and `ValidateToken`
2. An `AuthMiddleware` in `core/middleware/` that validates Bearer tokens
3. Registration of the `"auth"` alias in the middleware registry
4. A `User` model `BeforeCreate` hook for automatic password hashing (deferred gap from #11 audit)

## 2. Current State

### Already Built (dependencies satisfied)
- **#11 Models**: `User` model with `Password`, `Role`, `Active` fields
- **#19 Helpers**: `HashPassword()`, `CheckPassword()`, `RandomString()` — bcrypt + crypto/rand
- **#20 Session Management**: Full session system — `Manager`, `Start/Save/Destroy`, flash messages, middleware
- **#08 Middleware**: Registry with `RegisterAlias()`, `Resolve()`, `RegisterGroup()`, `ResolveGroup()`
- **`.env`**: `JWT_SECRET` and `JWT_EXPIRY` placeholders already exist

### Not Yet Built
- No `core/auth/` package
- No JWT dependency (`golang-jwt/jwt/v5`)
- No `AuthMiddleware` in middleware
- No `BeforeCreate` hook on User model

## 3. Dependency Analysis

The roadmap lists #22 (Crypto & Security Utilities) as a dependency. However, #21 does **not** functionally require #22's primitives:
- JWT crypto is handled by the `golang-jwt/jwt` library itself
- Password hashing uses existing `helpers.HashPassword` (bcrypt)
- Session crypto uses existing AES-256-GCM in cookie store

**Decision**: Proceed with #21 without #22. The crypto utilities package is standalone and will be built independently.

## 4. Blueprint Scope

From the blueprint, Feature #21 covers:

| Item | Blueprint Reference | Notes |
|------|---------------------|-------|
| `GenerateToken(userID uint)` | Lines 1199–1208 | HMAC-SHA256, `JWT_SECRET`, configurable expiry (default 1h) |
| `ValidateToken(tokenStr string)` | Lines 1210–1225 | Parse + validate claims |
| `AuthMiddleware()` | Lines 1075–1093 | Bearer token extraction + validation |
| `"auth"` alias registration | Lines 1149–1150 | `middleware.RegisterAlias("auth", ...)` |
| Session-based login pattern | Framework doc | Uses existing session manager |

### Not In Scope (deferred to future features)
- Refresh tokens — blueprint mentions but no code provided
- `"admin"` / `"verified"` aliases — user-defined, not framework-provided
- CSRF protection on session routes — Feature #24
- Login/register controllers — app-level, not framework-level
- OAuth / social login — not in blueprint

## 5. Approach

### JWT Package (`core/auth/`)
- `GenerateToken(userID uint) (string, error)` — reads `JWT_SECRET` and `JWT_EXPIRY` from env
- `ValidateToken(tokenStr string) (jwt.MapClaims, error)` — parses and validates token
- Use `jwt.SigningMethodHS256` — never allow `none` algorithm (security requirement from blueprint)

### Auth Middleware (`core/middleware/auth.go`)
- Extract `Authorization: Bearer <token>` header
- Call `auth.ValidateToken()` to validate
- Set `user_id` in Gin context via `c.Set("user_id", ...)`
- Abort with 401 on missing/invalid token

### User Model Enhancement (`database/models/user.go`)
- Add `BeforeCreate` GORM hook to auto-hash password via `helpers.HashPassword`
- This was identified as a gap in the #20 audit — actionable now since helpers exist

### Middleware Registration
- Register `"auth"` alias in `MiddlewareProvider` boot or in route setup

## 6. Edge Cases

1. **Empty JWT_SECRET**: Should fail loudly — don't silently use empty string for signing
2. **Expired tokens**: `jwt.Parse` handles expiry checking via `exp` claim
3. **Malformed Bearer header**: Must handle "Bearer" without token, no "Bearer " prefix
4. **Already-hashed password in BeforeCreate**: Must detect if password is already bcrypt-hashed to avoid double-hashing on updates

## 7. Dependencies

| Existing | New |
|----------|-----|
| `golang.org/x/crypto` (bcrypt) | `github.com/golang-jwt/jwt/v5` |
| `core/session/` (#20) | — |
| `app/helpers/` (#19) | — |
| `database/models/` (#11) | — |
| `core/middleware/registry.go` (#08) | — |

## 8. Open Questions — RESOLVED

| Question | Resolution |
|----------|-----------|
| Build #22 first? | No — #21 doesn't functionally need crypto utilities |
| Where to register auth alias? | In middleware provider or route setup — TBD in architecture |
| UserID type in claims? | `uint` — matches GORM `BaseModel.ID` type |
| Double-hash prevention? | Check if password starts with `$2a$` or `$2b$` bcrypt prefix |

---

**Summary**: Focused feature — JWT auth package, auth middleware, user model hook, alias registration. One new dependency (`golang-jwt/jwt/v5`). Builds on existing session, helpers, and middleware infrastructure.
