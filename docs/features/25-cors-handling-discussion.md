# 💬 Discussion: CORS Handling

> **Feature**: `25` — CORS Handling
> **Status**: FINAL
> **Date**: 2026-03-06

---

## 1. What Are We Building?

Enhancing the existing CORS middleware (from #08) to match the full blueprint specification. The base `CORS()` function and `CORSConfig` struct already exist; this feature adds environment-based configuration, `AllowCredentials`, `ExposeHeaders`, and updates default headers to include `X-CSRF-Token`.

## 2. Why?

The #08 implementation provides basic CORS headers but is missing:
- **Environment configuration** — `CORS_ALLOWED_ORIGINS` env var for deployment flexibility
- **Credentials support** — `Access-Control-Allow-Credentials` header for cookie-based auth
- **Expose headers** — `Access-Control-Expose-Headers` for `Content-Length`, `X-Request-ID`
- **X-CSRF-Token** — now that CSRF (#24) is shipped, it must be in allowed headers

## 3. Scope

### In Scope

- Add `AllowCredentials bool` to `CORSConfig`
- Add `ExposeHeaders []string` to `CORSConfig`
- Update `defaultCORSConfig()` to read `CORS_ALLOWED_ORIGINS` env var
- Add `X-CSRF-Token` to default `AllowHeaders`
- Set `Access-Control-Allow-Credentials` header when true
- Set `Access-Control-Expose-Headers` header
- Default `ExposeHeaders`: `Content-Length`, `X-Request-ID`
- Default `AllowCredentials`: `true`

### Out of Scope

- Switching to `gin-contrib/cors` library (blueprint recommends it, but custom impl is already working and tested — no reason to add a dependency)
- Per-route CORS configuration
- Dynamic origin matching with wildcards/subdomains

## 4. Dependencies

- **#08 (Middleware)** — existing CORS implementation to enhance
- **#24 (CSRF)** — `X-CSRF-Token` header support

Both shipped.

## 5. Key Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Keep custom impl | Don't add `gin-contrib/cors` | Already working; fewer dependencies; full control |
| Env reading | `os.Getenv("CORS_ALLOWED_ORIGINS")` | Blueprint pattern; deployment flexibility |
| Default credentials | `true` | Blueprint specification |
| Default expose | `Content-Length`, `X-Request-ID` | Blueprint specification |
