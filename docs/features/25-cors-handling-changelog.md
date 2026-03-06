# 📝 Changelog: CORS Handling

> **Feature**: `25` — CORS Handling
> **Branch**: `feature/25-cors-handling`
> **Started**: 2026-03-06
> **Completed**: 2026-03-06

---

## Log

- **BUILD**: All 9 CORS tests pass (3 existing + 6 new), full regression green, `go vet` clean
- **BUILD**: Added 6 tests (TC-26 to TC-31) to `core/middleware/middleware_test.go`
- **BUILD**: Enhanced `core/middleware/cors.go` — AllowCredentials, ExposeHeaders, env config, X-CSRF-Token
- **BUILD**: Created feature branch `feature/25-cors-handling`

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| None | — | — | Implementation matched architecture exactly |
