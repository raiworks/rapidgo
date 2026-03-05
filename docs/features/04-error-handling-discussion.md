# 💬 Discussion: Error Handling

> **Feature**: `04` — Error Handling
> **Status**: 🟢 COMPLETE
> **Branch**: `feature/04-error-handling`
> **Depends On**: #01 (Project Setup ✅), #03 (Logging ✅)
> **Date Started**: 2026-03-05
> **Date Completed**: 2026-03-05

---

## Summary

Implement a centralized error handling package that defines structured error types, helper constructors, and consistent error formatting for the RGo framework. This is the **foundation layer** — pure Go error types with no HTTP framework dependency. The Gin middleware integration will be built in a later HTTP/Router feature.

---

## Functional Requirements

- As a **framework developer**, I want a structured `AppError` type so that errors carry HTTP status codes, user-facing messages, and internal error details in a single value
- As a **framework developer**, I want error constructor helpers (`NotFound`, `BadRequest`, `Internal`, etc.) so that creating typed errors is concise and consistent
- As a **framework developer**, I want `AppError` to implement the `error` interface so it works seamlessly with Go's standard error handling (`errors.Is`, `errors.As`, `%w` wrapping)
- As a **framework developer**, I want a debug-aware error formatting function so that development shows full details while production shows only safe messages
- As a **framework developer**, I want all errors to be loggable with structured context via `slog` so that error observability is built-in from the start

## Current State / Reference

### What Exists
- **Configuration** (#02): `config.IsDebug()`, `config.AppEnv()` — available for debug-mode error formatting
- **Logging** (#03): `slog` globally configured — available for structured error logging
- **Blueprint reference**: Shows a Gin error middleware pattern, but the core error types are HTTP-framework-independent
- **Framework reference doc** (`docs/framework/core/error-handling.md`): Defines `AppError` struct, constructors, middleware pattern, security rules

### What Works Well
- `config.IsDebug()` already returns `true`/`false` based on `APP_DEBUG` env var — perfect for toggling error detail visibility
- `slog` is already set up globally — errors can log structured context immediately

### What Needs Improvement
- No error types exist yet — all errors are raw `error` values with no structure
- No standard way to attach HTTP status codes to errors
- No safe/unsafe formatting distinction for dev vs prod

## Proposed Approach

Create a `core/errors` package with:

1. **`AppError` struct** — carries `Code` (HTTP status), `Message` (user-safe), `Err` (internal, wrapped)
2. **Constructor helpers** — `NotFound()`, `BadRequest()`, `Internal()`, `Unauthorized()`, `Forbidden()`, `Conflict()`, `Unprocessable()` — each returns `*AppError` with the correct status code
3. **`error` interface compliance** — `Error()` returns `Message`, `Unwrap()` returns `Err` for `errors.As`/`errors.Is` support
4. **`ErrorResponse` helper** — returns a map suitable for JSON responses, debug-aware (includes internal error details only when `APP_DEBUG=true`)
5. **Unit tests** — cover all constructors, unwrapping, and debug/production formatting

**NOT in scope for this feature** (deferred to HTTP/Router feature):
- Gin middleware (`ErrorHandler()`)
- Content negotiation (JSON vs HTML)
- Panic recovery middleware
- HTTP response writing

## Edge Cases & Risks

- [x] `AppError.Err` can be nil (e.g., `NotFound("user not found")` with no underlying error) — `Unwrap()` must handle nil gracefully
- [x] `ErrorResponse()` must NEVER expose internal error details when `APP_DEBUG=false` — security critical
- [x] Package naming conflict: Go stdlib has `errors` package — using `core/errors` is fine because it's imported by full path (`github.com/RAiWorks/RGo/core/errors`), but internal code must use the full import path to avoid shadowing
- [x] `AppError` must work with `errors.As()` — requires pointer receiver pattern

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Feature #01 — Project Setup | Feature | ✅ Done |
| Feature #03 — Logging | Feature | ✅ Done |
| Feature #02 — Configuration | Feature | ✅ Done (needed for `config.IsDebug()`) |
| `log/slog` | Stdlib | ✅ Available |
| `errors` | Stdlib | ✅ Available |

## Open Questions

_All resolved during discussion._

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-05 | Scope to core error types only, no Gin middleware | Gin is not yet a dependency; middleware comes with HTTP/Router feature |
| 2026-03-05 | Use `core/errors` package path | Matches framework convention (`core/config`, `core/logger`); full import path avoids stdlib shadowing |
| 2026-03-05 | Include `Unwrap()` for `errors.As`/`errors.Is` support | Idiomatic Go error wrapping — essential for downstream error inspection |
| 2026-03-05 | `ErrorResponse()` uses `config.IsDebug()` for detail toggle | Single source of truth for debug mode; already implemented in Feature #02 |
| 2026-03-05 | 7 constructor helpers covering common HTTP status codes | Covers 400, 401, 403, 404, 409, 422, 500 — the most common API error scenarios |

## Discussion Complete ✅

**Summary**: Feature #04 implements a `core/errors` package with structured `AppError` type, 7 HTTP-status-aware constructors, `error`/`Unwrap` interface compliance, and debug-aware response formatting. Scoped to pure Go types — no Gin dependency. Middleware integration deferred to HTTP/Router feature.
**Completed**: 2026-03-05
**Next**: Create architecture doc → `04-error-handling-architecture.md`
