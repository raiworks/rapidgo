# 🧪 Test Plan: Router & Routing

> **Feature**: `07` — Router & Routing
> **Tasks**: [`07-router-tasks.md`](07-router-tasks.md)
> **Date**: 2026-03-06

---

## Acceptance Criteria

- [ ] `Router` struct wraps `*gin.Engine` and provides HTTP method helpers
- [ ] `RouteGroup` supports sub-groups with shared prefix and middleware
- [ ] `ResourceController` interface has 7 CRUD methods
- [ ] `Resource()` registers 7 routes, `APIResource()` registers 5 routes
- [ ] Named routes support URL generation with parameter substitution
- [ ] `RouterProvider` registers router as `"router"` in the container
- [ ] `routes/web.go` and `routes/api.go` have `RegisterWeb`/`RegisterAPI` functions
- [ ] `cmd/main.go` starts HTTP server via the router
- [ ] Gin mode set correctly based on `APP_ENV`
- [ ] All tests pass with `go test ./core/router/...`
- [ ] `go vet ./...` reports no issues

---

## Test Cases

### TC-01: Router.New creates a valid router

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Set APP_ENV to "testing" → 2. Call `router.New()` → 3. Check router is not nil → 4. Check `Engine()` is not nil |
| **Expected Result** | Router and engine are both non-nil |
| **Status** | ⬜ Not Run |

### TC-02: Gin mode set to ReleaseMode in production

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | APP_ENV=production |
| **Steps** | 1. Set APP_ENV="production" → 2. Call `router.New()` → 3. Check `gin.Mode()` |
| **Expected Result** | `gin.Mode()` returns `"release"` |
| **Status** | ⬜ Not Run |

### TC-03: Gin mode set to TestMode in testing

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | APP_ENV=testing |
| **Steps** | 1. Set APP_ENV="testing" → 2. Call `router.New()` → 3. Check `gin.Mode()` |
| **Expected Result** | `gin.Mode()` returns `"test"` |
| **Status** | ⬜ Not Run |

### TC-04: Gin mode defaults to DebugMode for development

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | APP_ENV=development |
| **Steps** | 1. Set APP_ENV="development" → 2. Call `router.New()` → 3. Check `gin.Mode()` |
| **Expected Result** | `gin.Mode()` returns `"debug"` |
| **Status** | ⬜ Not Run |

### TC-05: GET route responds with 200

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Router created |
| **Steps** | 1. Create router → 2. Register `r.Get("/ping", handler)` → 3. Send GET `/ping` via httptest → 4. Check status |
| **Expected Result** | Status 200, response body matches |
| **Status** | ⬜ Not Run |

### TC-06: POST route responds with 200

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Router created |
| **Steps** | 1. Create router → 2. Register `r.Post("/data", handler)` → 3. Send POST `/data` via httptest |
| **Expected Result** | Status 200 |
| **Status** | ⬜ Not Run |

### TC-07: PUT route responds with 200

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Router created |
| **Steps** | 1. Register `r.Put("/items/:id", handler)` → 2. Send PUT `/items/1` via httptest |
| **Expected Result** | Status 200 |
| **Status** | ⬜ Not Run |

### TC-08: DELETE route responds with 200

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Router created |
| **Steps** | 1. Register `r.Delete("/items/:id", handler)` → 2. Send DELETE `/items/1` via httptest |
| **Expected Result** | Status 200 |
| **Status** | ⬜ Not Run |

### TC-09: PATCH route responds with 200

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Router created |
| **Steps** | 1. Register `r.Patch("/items/:id", handler)` → 2. Send PATCH `/items/1` via httptest |
| **Expected Result** | Status 200 |
| **Status** | ⬜ Not Run |

### TC-10: OPTIONS route responds with 200

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Router created |
| **Steps** | 1. Register `r.Options("/cors", handler)` → 2. Send OPTIONS `/cors` via httptest |
| **Expected Result** | Status 200 |
| **Status** | ⬜ Not Run |

### TC-11: Route group adds prefix to paths

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Router created |
| **Steps** | 1. Create group `/api` → 2. Register `group.Get("/users", handler)` → 3. Send GET `/api/users` |
| **Expected Result** | Status 200 |
| **Status** | ⬜ Not Run |

### TC-12: Nested route group combines prefixes

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Router created |
| **Steps** | 1. Create group `/api` → 2. Create sub-group `/v1` → 3. Register `sub.Get("/users", handler)` → 4. Send GET `/api/v1/users` |
| **Expected Result** | Status 200 |
| **Status** | ⬜ Not Run |

### TC-13: Resource registers all 7 RESTful routes

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Mock ResourceController, Router created |
| **Steps** | 1. Call `r.Resource("/posts", mockCtrl)` → 2. Send requests to all 7 routes → 3. Verify each returns 200 |
| **Expected Result** | GET /posts (Index), GET /posts/create (Create), POST /posts (Store), GET /posts/1 (Show), GET /posts/1/edit (Edit), PUT /posts/1 (Update), DELETE /posts/1 (Destroy) — all return 200 |
| **Status** | ⬜ Not Run |

### TC-14: APIResource registers 5 routes (no Create/Edit)

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Mock ResourceController, Router created |
| **Steps** | 1. Call `r.APIResource("/users", mockCtrl)` → 2. Send requests to 5 routes → 3. Verify each returns 200 → 4. Verify GET /users/create returns 404 → 5. Verify GET /users/1/edit returns 404 |
| **Expected Result** | 5 routes return 200, create/edit form routes return 404 |
| **Status** | ⬜ Not Run |

### TC-15: Resource routes on RouteGroup combine prefix

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Mock ResourceController, Group created |
| **Steps** | 1. Create group `/api` → 2. Call `group.Resource("/posts", mockCtrl)` → 3. Send GET `/api/posts` |
| **Expected Result** | Status 200 |
| **Status** | ⬜ Not Run |

### TC-16: APIResource on RouteGroup combines prefix

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Mock ResourceController, Group created |
| **Steps** | 1. Create group `/api` → 2. Call `group.APIResource("/users", mockCtrl)` → 3. Send GET `/api/users` |
| **Expected Result** | Status 200 |
| **Status** | ⬜ Not Run |

### TC-17: Named route — Name and Route generate correct URL

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. `Name("users.show", "/users/:id")` → 2. `Route("users.show", "42")` |
| **Expected Result** | Returns `"/users/42"` |
| **Status** | ⬜ Not Run |

### TC-18: Named route — multiple parameters

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. `Name("posts.comment", "/posts/:postId/comments/:commentId")` → 2. `Route("posts.comment", "5", "99")` |
| **Expected Result** | Returns `"/posts/5/comments/99"` |
| **Status** | ⬜ Not Run |

### TC-19: Named route — unknown name returns "/"

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | No routes registered |
| **Steps** | 1. `Route("nonexistent")` |
| **Expected Result** | Returns `"/"` |
| **Status** | ⬜ Not Run |

### TC-20: Named route — no parameters returns pattern as-is

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | None |
| **Steps** | 1. `Name("home", "/")` → 2. `Route("home")` |
| **Expected Result** | Returns `"/"` |
| **Status** | ⬜ Not Run |

### TC-21: Router implements http.Handler (ServeHTTP)

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Router created with route |
| **Steps** | 1. Create router with GET /test → 2. Use `httptest.NewRecorder` + `http.NewRequest` → 3. Call `r.ServeHTTP(w, req)` |
| **Expected Result** | Response recorded correctly (status 200) |
| **Status** | ⬜ Not Run |

### TC-22: Group middleware executes for group routes only

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Router with group and middleware |
| **Steps** | 1. Register root route `/public` → 2. Create group `/api` with header-setting middleware → 3. Register `/api/data` → 4. Send GET `/api/data` — check header present → 5. Send GET `/public` — check header absent |
| **Expected Result** | Middleware header only present on group routes |
| **Status** | ⬜ Not Run |

### TC-23: RouterProvider implements Provider interface

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | None |
| **Steps** | 1. Compile-time check: `var _ container.Provider = (*RouterProvider)(nil)` |
| **Expected Result** | Compiles without error |
| **Status** | ⬜ Not Run |

### TC-24: RouterProvider registers router as "router" in container

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Container created |
| **Steps** | 1. Create container → 2. Call `RouterProvider.Register(c)` → 3. `container.MustMake[*router.Router](c, "router")` |
| **Expected Result** | Returns non-nil `*Router`, no panic |
| **Status** | ⬜ Not Run |

### TC-25: Route parameter extraction via Gin context

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Router with parameterized route |
| **Steps** | 1. Register `r.Get("/users/:id", handler)` where handler reads `c.Param("id")` → 2. Send GET `/users/42` → 3. Check response body contains "42" |
| **Expected Result** | Parameter extracted correctly |
| **Status** | ⬜ Not Run |

---

## Edge Cases

| # | Scenario | Expected Behavior |
|---|---|---|
| 1 | Unknown named route | `Route()` returns `"/"` |
| 2 | Named route with no params | Returns pattern unchanged (static path) |
| 3 | Unregistered path request | Gin returns 404 automatically |
| 4 | Empty RegisterWeb/RegisterAPI | No-op, no panic |

## Test Summary

| Category | Total | Pass | Fail | Skip |
|---|---|---|---|---|
| Happy Path | 22 | — | — | — |
| Edge Cases | 3 | — | — | — |
| **Total** | **25** | — | — | — |

**Result**: ⬜ NOT RUN
