# 📝 Changelog: Router & Routing

> **Feature**: `07` — Router & Routing
> **Branch**: `feature/07-router`
> **Started**: 2026-03-06
> **Completed**: 2026-03-06

---

## Log

- **2026-03-06** — Added `github.com/gin-gonic/gin v1.12.0` dependency (requires Go 1.25.0, upgraded from Go 1.21)
- **2026-03-06** — Created `core/router/router.go`: `Router` struct wrapping `*gin.Engine`, `New()` with env-based Gin mode, HTTP method helpers, `Group()`, `Use()`, `Run()`, `ServeHTTP()`
- **2026-03-06** — Created `core/router/group.go`: `RouteGroup` struct wrapping `*gin.RouterGroup`, same HTTP helpers + nesting + middleware
- **2026-03-06** — Created `core/router/resource.go`: `ResourceController` interface (7 methods), `Resource()` (7 routes) and `APIResource()` (5 routes) on Router and RouteGroup
- **2026-03-06** — Created `core/router/named.go`: Thread-safe named route registry with `sync.RWMutex`, `Name()`, `Route()`, `ResetNamedRoutes()`
- **2026-03-06** — Created `app/providers/router_provider.go`: `RouterProvider` — `Register()` creates router, `Boot()` calls route registration
- **2026-03-06** — Updated `routes/web.go` and `routes/api.go`: Added `RegisterWeb()` / `RegisterAPI()` functions
- **2026-03-06** — Updated `cmd/main.go`: Added `RouterProvider`, router resolution from container, HTTP server startup
- **2026-03-06** — Created `core/router/router_test.go`: 23 test functions covering all router functionality
- **2026-03-06** — Updated `app/providers/providers_test.go`: 2 new tests (RouterProvider registration + full bootstrap with router)
- **2026-03-06** — All 88 tests pass across entire project, `go vet` clean

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| Go version upgraded to 1.25.0 | `go.mod` had `go 1.21` | Gin v1.12.0 requires Go 1.25.0, `go mod tidy` auto-upgraded | Gin's minimum Go version requirement |
| TC-14 adjusted for Gin routing behavior | Test expected `GET /users/create` → 404 when APIResource registered | Removed that assertion; Gin's `/users/:id` matches `/users/create` with `id=create` | Gin parameterized routes are greedy — this is expected framework behavior, not a bug |

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
| Use `gin.New()` not `gin.Default()` | Avoids Gin's default Logger/Recovery middleware — framework controls its own middleware stack | 2026-03-06 |
| Map `APP_ENV` to Gin modes in `setGinMode()` | production→release, testing→test, development/default→debug — consistent with framework's existing env pattern | 2026-03-06 |
| Global named route registry with `sync.RWMutex` | Matches blueprint specification; thread-safe for concurrent access | 2026-03-06 |
| `ResourceController` uses `gin.HandlerFunc` | Direct Gin handler type — no unnecessary abstraction layer before middleware feature | 2026-03-06 |
