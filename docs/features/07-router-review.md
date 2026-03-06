# 🪞 Review: Router & Routing

> **Feature**: `07` — Router & Routing
> **Branch**: `feature/07-router`
> **Merged**: 2026-03-06
> **Duration**: 2026-03-06 → 2026-03-06

---

## Result

**Status**: ✅ Shipped

**Summary**: Implemented the `core/router` package wrapping Gin v1.12.0 with a framework-level API. Router struct with HTTP method helpers, route groups with prefix/middleware isolation, RESTful resource routes (7-route `Resource()` and 5-route `APIResource()`), thread-safe named route registry with URL generation, and `RouterProvider` for service container integration. 25/25 tests pass. 88 total tests across all packages — zero regressions.

---

## What Went Well ✅

- Architecture doc mapped cleanly to implementation — 4 source files matching the 4 documented components
- `gin.New()` instead of `gin.Default()` keeps full middleware control with the framework
- `setGinMode()` bridges `APP_ENV` → Gin modes seamlessly using existing config infrastructure
- `ResourceController` interface is straightforward — 7 `gin.HandlerFunc` methods, no unnecessary abstraction
- Named route registry with `sync.RWMutex` is thread-safe without overcomplication
- Provider integration worked first try — `RouterProvider` follows the same Register/Boot pattern as Config and Logger

## What Went Wrong ❌

- TC-14 initially failed: test expected `GET /users/create` to return 404 when `APIResource` is used, but Gin's parameterized route `GET /users/:id` matches `/users/create` with `id=create`. Test expectation was adjusted.
- Go version in `go.mod` upgraded from 1.21 to 1.25.0 — forced by Gin v1.12.0's minimum Go version requirement. Not a problem in practice but was unexpected.

## What Was Learned 📚

- Gin parameterized routes (`:id`) are greedy — they match any single path segment including literal words like "create". Tests must account for this behavior.
- When adding framework-level dependencies (Gin), `go mod tidy` can auto-upgrade the Go version. Always verify `go.mod` after adding major dependencies.
- Wrapping a third-party router requires careful API design — expose enough flexibility without leaking the underlying engine's internals.

## What To Do Differently Next Time 🔄

- Verify Gin routing semantics for parameterized routes *before* writing tests that assume 404 for masked paths
- Check Go version compatibility of new dependencies before adding them

## Metrics

| Metric | Value |
|---|---|
| Tasks planned | 22 |
| Tasks completed | 22 |
| Tests planned | 25 |
| Tests passed | 25 |
| Deviations from plan | 1 (TC-14 assertion adjusted) |
| Files created | 6 |
| Files modified | 5 |
| Commits on branch | 1 |
| Total project tests | 88 |

## Follow-ups

- [ ] Middleware pipeline (Feature #08) — `Use()` is ready, needs middleware registration/ordering system
- [ ] Route model binding — deferred to later feature per architecture doc
- [ ] Gin error middleware using `AppError` from Feature #04
- [ ] Content negotiation for error responses (JSON vs HTML)
