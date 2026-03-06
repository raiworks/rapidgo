# 💬 Discussion: Controllers

> **Feature**: `15` — Controllers
> **Status**: 🟢 COMPLETE
> **Date**: 2026-03-06

---

## What Are We Building?

MVC controllers in `http/controllers/` that handle HTTP requests and return responses. This includes a `Home` function controller for the root route and a `PostController` struct implementing the `ResourceController` interface for full CRUD. Route registration is wired in `routes/web.go` and `routes/api.go`.

## Blueprint References

The blueprint specifies:

1. **Directory**: `http/controllers/` (line 82)
2. **Function controller**: `Home(c *gin.Context)` — returns HTML (line 414–418)
3. **Function controller**: `Users(c *gin.Context)` — returns JSON (line 420–425)
4. **MVC Controller Example section** (lines 403–425)
5. **ResourceController interface**: already shipped in Feature #07 (`core/router/resource.go`)
6. **PostController struct**: example in framework doc implementing all 7 CRUD methods

## Scope for Feature #15

### In Scope
- `http/controllers/home_controller.go` — `Home()` function returning JSON welcome message
- `http/controllers/post_controller.go` — `PostController` struct implementing `ResourceController` (all 7 methods)
- `routes/web.go` — register `Home` on `GET /`
- `routes/api.go` — register `PostController` via `APIResource` on `/api/posts`
- Tests verifying controller responses and route registration

### Out of Scope (deferred)
- `make:controller` scaffolding — CLI code generation is a future feature
- Views / HTML templates — Feature #17
- Services layer delegation — Feature #18
- Request validation — future feature
- Response helpers (`responses.Success`, etc.) — Feature #16
- Database queries in controllers — controllers return static/mock data for now

## Key Design Decisions

### 1. JSON-Only Responses (No HTML Yet)
The blueprint shows `Home` returning `c.HTML()`, but Views (#17) aren't built yet. `Home` returns a JSON welcome message instead. `PostController` methods all return JSON. HTML rendering will be added when #17 ships.

### 2. PostController as the Example
The blueprint and framework doc both use `PostController` as the primary example. We implement all 7 `ResourceController` methods with placeholder/static responses — no database interaction yet.

### 3. Home as a Function, PostController as a Struct
Demonstrates both patterns from the blueprint: simple function handlers and struct-based resource controllers.

### 4. APIResource Registration (Not Full Resource)
`PostController` is registered via `APIResource` (5 routes, no Create/Edit form routes) since we're API-only for now. Full `Resource` registration with form routes will be relevant when Views (#17) ships.

## Dependencies

| Dependency | Status | Notes |
|---|---|---|
| Feature #07 — Router | ✅ Done | Provides Router, RouteGroup, ResourceController, APIResource |
| Feature #08 — Middleware | ✅ Done | Middleware stack available for route groups |

## Discussion Complete ✅
