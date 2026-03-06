# 💬 Discussion: Router & Routing

> **Feature**: `07` — Router & Routing
> **Status**: 🟢 COMPLETE
> **Branch**: `docs/07-router`
> **Depends On**: #01 (Project Setup ✅), #05 (Service Container ✅)
> **Date Started**: 2026-03-06
> **Date Completed**: 2026-03-06

---

## Summary

Implement the Router layer using Gin as the underlying HTTP engine. This feature establishes route registration (GET, POST, PUT, DELETE, PATCH, OPTIONS), route groups with shared prefixes, resource routes (RESTful CRUD), named routes with URL generation, and integration with the service container via a `RouterProvider`. The router is the foundation for all HTTP handling in the framework.

---

## Functional Requirements

- As a **framework developer**, I want a `Router` struct that wraps `*gin.Engine` so that route registration is centralized and framework-controlled
- As a **framework developer**, I want HTTP method helpers (`Get`, `Post`, `Put`, `Delete`, `Patch`, `Options`) so that routes are registered with clean syntax
- As a **framework developer**, I want `Group(prefix, ...handlers)` so that routes can be organized under shared prefixes with optional middleware
- As a **framework developer**, I want `Resource(path, controller)` and `APIResource(path, controller)` so that RESTful CRUD routes are registered in a single call
- As a **framework developer**, I want `ResourceController` interface so that controllers can implement standard CRUD methods
- As a **framework developer**, I want `Name(name, pattern)` and `Route(name, params...)` so that URLs can be generated from named routes
- As a **framework developer**, I want a `RouterProvider` so that the Gin engine is bootstrapped through the provider lifecycle and registered in the container as `"router"`
- As a **framework developer**, I want `routes/web.go` and `routes/api.go` to have `Register(r)` functions so that route definitions are organized by purpose

## Current State / Reference

### What Exists
- **Service Container** (#05 ✅): `Container`, `Provider` interface, `App` struct — fully operational
- **Service Providers** (#06 ✅): `ConfigProvider`, `LoggerProvider` — App bootstrap pattern established
- **`core/router/`**: Empty directory with `.gitkeep`
- **`routes/web.go`**: Placeholder package declaration
- **`routes/api.go`**: Placeholder package declaration
- **`http/controllers/`**: Empty directory with `.gitkeep`
- **`core/middleware/`**: Empty directory with `.gitkeep`

### Blueprint Reference
The blueprint uses Gin (`github.com/gin-gonic/gin`) and shows:
1. `SetupRouter() *gin.Engine` — creates engine, registers routes directly
2. `ResourceController` interface — 7 CRUD methods (Index, Create, Store, Show, Edit, Update, Destroy)
3. `Resource()` and `APIResource()` — register all RESTful routes for a controller
4. Named routes — global `namedRoutes` map with `Name()` and `Route()` functions
5. Route model binding — middleware that loads GORM models by `:id`

### What the Blueprint Does NOT Show
- How the router integrates with the service container/provider lifecycle
- How Gin's mode is set based on `APP_ENV`
- How routes are organized in separate files (web vs API)

## Proposed Approach

### Router Wrapper (`core/router/router.go`)

Create a `Router` struct that wraps `*gin.Engine` and provides framework-level route registration methods. This gives us a clean API while keeping Gin's performance underneath.

```go
type Router struct {
    engine *gin.Engine
}
```

**Methods:**
- `New() *Router` — creates Gin engine, sets Gin mode based on `APP_ENV`
- `Get(path, ...handlers)`, `Post(...)`, `Put(...)`, `Delete(...)`, `Patch(...)`, `Options(...)` — HTTP method registration
- `Group(prefix, ...handlers) *RouteGroup` — creates a route group
- `Resource(path, controller)` — registers all 7 RESTful routes
- `APIResource(path, controller)` — registers 5 RESTful routes (no Create/Edit form routes)
- `Engine() *gin.Engine` — exposes underlying Gin engine (for server startup)
- `ServeHTTP(w, r)` — implements `http.Handler` for flexibility

### Route Groups (`core/router/group.go`)

Wrap `*gin.RouterGroup` with the same method set:

```go
type RouteGroup struct {
    group *gin.RouterGroup
}
```

- Same HTTP method helpers as Router
- `Group(prefix, ...handlers) *RouteGroup` — nested groups
- `Resource()` / `APIResource()` — resource routes on groups

### Resource Controller Interface (`core/router/resource.go`)

```go
type ResourceController interface {
    Index(c *gin.Context)
    Create(c *gin.Context)
    Store(c *gin.Context)
    Show(c *gin.Context)
    Edit(c *gin.Context)
    Update(c *gin.Context)
    Destroy(c *gin.Context)
}
```

### Named Routes (`core/router/named.go`)

Thread-safe named route registry with URL generation:
- `Name(name, pattern)` — register a named route
- `Route(name, params...) string` — generate URL from name with parameter substitution

### Router Provider (`app/providers/router_provider.go`)

```go
type RouterProvider struct{}
```

- `Register()` — creates `router.New()`, registers as `"router"` singleton in container
- `Boot()` — calls route registration functions (`routes.RegisterWeb`, `routes.RegisterAPI`)

### Route Files (`routes/web.go`, `routes/api.go`)

Updated from placeholders to actual route registration:
- `routes.RegisterWeb(r *router.Router)` — web (HTML) routes
- `routes.RegisterAPI(r *router.Router)` — API routes under `/api` group

### Updated main.go

Add `RouterProvider` to the bootstrap chain after Config and Logger.

## Scoping Decisions

### IN Scope
| Item | Rationale |
|---|---|
| Router struct wrapping Gin | Blueprint core |
| HTTP method helpers (6 methods) | Blueprint core |
| Route groups with prefix | Blueprint core |
| Resource routes (Resource + APIResource) | Blueprint core |
| ResourceController interface | Blueprint core |
| Named routes + URL generation | Blueprint core |
| RouterProvider for container integration | Natural extension of #05/#06 |
| Route file organization (web.go, api.go) | Blueprint structure |
| Gin mode from APP_ENV | Framework integration |

### NOT in Scope
| Item | Rationale | Future Feature |
|---|---|---|
| Route model binding (BindModel) | Depends on GORM/Database (#09, #11) | #15 or #11 |
| Middleware definitions (Auth, CSRF, etc.) | Feature #08 | #08 |
| Middleware registry (aliases, groups) | Feature #08 | #08 |
| Error handling middleware | Feature #08 | #08 |
| Controllers | Feature #15 | #15 |
| Response helpers | Feature #16 | #16 |
| Static file serving | Feature #30 | #30 |
| HTTP server start/listen | Needs server package | #07 (minimal — just `Run()`) |

### Key Decision: Gin as the Router Engine
The blueprint recommends Gin or Chi. We choose **Gin** because:
1. Blueprint examples all use Gin
2. Most widely adopted Go web framework
3. Fastest HTTP router available
4. Rich middleware ecosystem
5. Battle-tested in production

## Edge Cases & Risks

- [x] Gin mode must match APP_ENV — set `gin.SetMode()` before creating engine
- [x] Named routes are global — need thread-safe map (sync.RWMutex)
- [x] Resource routes need consistent path patterns (`/resource`, `/resource/:id`, etc.)
- [x] Route groups must support nesting (group within group)
- [x] `gin.HandlerFunc` is the handler type — framework handlers use Gin's context directly for now
- [x] Empty route files should not panic — RegisterWeb/RegisterAPI are no-ops until real routes exist

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Feature #01 — Project Setup | Feature | ✅ Done |
| Feature #05 — Service Container | Feature | ✅ Done |
| Feature #06 — Service Providers | Feature | ✅ Done |
| `github.com/gin-gonic/gin` | External | New dependency |

## Open Questions

_All resolved during discussion._

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-06 | Use Gin as router engine | Blueprint uses Gin in all examples; fastest, most adopted Go router |
| 2026-03-06 | Wrap Gin in Router struct | Framework-controlled API surface; swap engine later if needed |
| 2026-03-06 | RouterProvider registers router as `"router"` singleton | Follows container/provider pattern from #05/#06 |
| 2026-03-06 | Defer route model binding | Depends on GORM (#09/#11) — not available yet |
| 2026-03-06 | Defer middleware registry | Feature #08 — keep #07 focused on routing only |
| 2026-03-06 | Handler type is `gin.HandlerFunc` | Matches blueprint; abstract later if needed |
