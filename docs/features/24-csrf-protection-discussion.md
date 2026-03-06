# 💬 Discussion: CSRF Protection

> **Feature**: `24` — CSRF Protection
> **Status**: FINAL
> **Date**: 2026-03-06

---

## 1. What Are We Building?

A CSRF (Cross-Site Request Forgery) protection middleware that generates per-session tokens and validates them on state-changing HTTP requests (POST, PUT, PATCH, DELETE). Essential for any SSR application rendering HTML forms.

## 2. Why?

Without CSRF protection, a malicious site can trick an authenticated user's browser into making unwanted requests. The framework needs built-in CSRF middleware so SSR applications are secure by default.

## 3. Scope

### In Scope

- `CSRFMiddleware()` function in `core/middleware/csrf.go`
- Token generation via `crypto/rand` (32 random bytes → 64-char hex)
- Token stored in session data under `_csrf_token` key
- Token exposed to templates via `c.Set("csrf_token", token)`
- Safe methods skip validation: GET, HEAD, OPTIONS
- State-changing methods validate: POST, PUT, PATCH, DELETE
- Token accepted from form field `_csrf_token` or header `X-CSRF-Token`
- 403 Forbidden with JSON error on mismatch
- Register `"csrf"` alias in middleware provider

### Out of Scope

- Double-submit cookie pattern (session-based is sufficient)
- Per-request token rotation (per-session is the blueprint approach)
- SameSite cookie configuration (separate concern)

## 4. Dependencies

- **#08 (Middleware)** — middleware registration system
- **#20 (Session Management)** — session data stored in `c.Get("session")` as `map[string]interface{}`

Both are shipped.

## 5. Key Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Token storage | Session (`_csrf_token` key) | Blueprint pattern; ties to existing session infrastructure |
| Token size | 32 bytes (64 hex chars) | Sufficient entropy for CSRF tokens |
| Token source | `crypto/rand` | Cryptographically secure |
| Validation sources | Form field + header | Supports both SSR forms and AJAX/API calls |
| Safe methods | GET, HEAD, OPTIONS | RFC 7231 safe methods |
| Error response | 403 JSON `{"error": "CSRF token mismatch"}` | Blueprint specification |
| Alias | `"csrf"` | Consistent with existing aliases |
