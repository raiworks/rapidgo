# 📝 Changelog: Middleware Pipeline

> **Feature**: `08` — Middleware Pipeline
> **Branch**: `feature/08-middleware`
> **Started**: 2026-03-06
> **Completed**: 2026-03-06

---

## Log

- **2026-03-06** — Removed `core/middleware/.gitkeep`
- **2026-03-06** — Created `core/middleware/registry.go`: `RegisterAlias`, `RegisterGroup`, `Resolve`, `ResolveGroup`, `ResetRegistry`
- **2026-03-06** — Created `core/middleware/recovery.go`: `Recovery()` — catches panics, logs via `slog.Error`, returns 500 JSON
- **2026-03-06** — Created `core/middleware/request_id.go`: `RequestID()` — UUID v4 via `crypto/rand`, preserves incoming header
- **2026-03-06** — Created `core/middleware/cors.go`: `CORS()` with `CORSConfig` struct, preflight OPTIONS → 204
- **2026-03-06** — Created `core/middleware/error_handler.go`: `ErrorHandler()` — formats `AppError` as JSON, wraps generic errors as 500
- **2026-03-06** — Created `app/providers/middleware_provider.go`: registers 4 aliases + `global` group
- **2026-03-06** — Updated `cmd/main.go`: inserted `MiddlewareProvider` as provider #3 before `RouterProvider` (#4)
- **2026-03-06** — Created `core/middleware/middleware_test.go`: 17 test functions
- **2026-03-06** — Updated `app/providers/providers_test.go`: 2 new tests (compile-time check + aliases verification)
- **2026-03-06** — All 107 tests pass across entire project, `go vet` clean

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| None | — | — | Implementation matched architecture doc exactly |

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
| No mutex on registry maps | Maps written at boot (single goroutine), read-only during request handling | 2026-03-06 |
| `c.Writer.Written()` guard in Recovery | Prevents double-write if handler partially wrote before panicking | 2026-03-06 |
