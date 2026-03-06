# 🏗️ Architecture: Router & Routing

> **Feature**: `07` — Router & Routing
> **Discussion**: [`07-router-discussion.md`](07-router-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

The Router layer wraps Gin (`github.com/gin-gonic/gin`) in a framework-controlled API. It provides HTTP method helpers, route groups, RESTful resource routes, named routes with URL generation, and integrates with the service container via `RouterProvider`. Route definitions live in `routes/web.go` and `routes/api.go`.

## File Structure

```
core/router/
├── router.go           # Router struct, New(), HTTP method helpers, ServeHTTP
├── group.go            # RouteGroup struct, HTTP method helpers, nested groups
├── resource.go         # ResourceController interface, Resource(), APIResource()
└── named.go            # Named route registry, Name(), Route()

core/router/
└── router_test.go      # Tests for router, groups, resources, named routes

app/providers/
└── router_provider.go  # RouterProvider — creates router, registers in container

routes/
├── web.go              # MODIFIED — RegisterWeb(r *router.Router)
└── api.go              # MODIFIED — RegisterAPI(r *router.Router)

cmd/
└── main.go             # MODIFIED — adds RouterProvider, starts HTTP server
```

### Files Created (5)
| File | Package | Lines (est.) |
|---|---|---|
| `core/router/router.go` | `router` | ~80 |
| `core/router/group.go` | `router` | ~60 |
| `core/router/resource.go` | `router` | ~50 |
| `core/router/named.go` | `router` | ~55 |
| `app/providers/router_provider.go` | `providers` | ~25 |

### Files Modified (3)
| File | Change |
|---|---|
| `routes/web.go` | Add `RegisterWeb(r)` function |
| `routes/api.go` | Add `RegisterAPI(r)` function |
| `cmd/main.go` | Add RouterProvider, start server |

---

## Component Design

### Router (`core/router/router.go`)

**Responsibility**: Wrap `*gin.Engine`, provide HTTP method registration, expose engine for server startup
**Package**: `router`

```go
package router

import (
	"net/http"

	"github.com/RAiWorks/RGo/core/config"
	"github.com/gin-gonic/gin"
)

// Router wraps the Gin engine and provides framework-level route registration.
type Router struct {
	engine *gin.Engine
}

// New creates a new Router with Gin mode set based on APP_ENV.
func New() *Router {
	setGinMode()
	engine := gin.New()
	return &Router{engine: engine}
}

// setGinMode configures Gin's mode based on the APP_ENV environment variable.
func setGinMode() {
	switch config.AppEnv() {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "testing":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

// Engine returns the underlying Gin engine.
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// ServeHTTP implements http.Handler, delegating to the Gin engine.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}

// --- HTTP Method Helpers ---

// Get registers a GET route.
func (r *Router) Get(path string, handlers ...gin.HandlerFunc) {
	r.engine.GET(path, handlers...)
}

// Post registers a POST route.
func (r *Router) Post(path string, handlers ...gin.HandlerFunc) {
	r.engine.POST(path, handlers...)
}

// Put registers a PUT route.
func (r *Router) Put(path string, handlers ...gin.HandlerFunc) {
	r.engine.PUT(path, handlers...)
}

// Delete registers a DELETE route.
func (r *Router) Delete(path string, handlers ...gin.HandlerFunc) {
	r.engine.DELETE(path, handlers...)
}

// Patch registers a PATCH route.
func (r *Router) Patch(path string, handlers ...gin.HandlerFunc) {
	r.engine.PATCH(path, handlers...)
}

// Options registers an OPTIONS route.
func (r *Router) Options(path string, handlers ...gin.HandlerFunc) {
	r.engine.OPTIONS(path, handlers...)
}

// Group creates a new route group with a shared prefix and optional middleware.
func (r *Router) Group(prefix string, handlers ...gin.HandlerFunc) *RouteGroup {
	return &RouteGroup{group: r.engine.Group(prefix, handlers...)}
}

// Use adds global middleware to the router.
func (r *Router) Use(middleware ...gin.HandlerFunc) {
	r.engine.Use(middleware...)
}

// Run starts the HTTP server on the given address.
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
```

**Design notes**:
- `gin.New()` instead of `gin.Default()` — no default middleware attached. Middleware is Feature #08's responsibility.
- `setGinMode()` reads `config.AppEnv()` — requires config to be loaded first (enforced by provider order)
- `ServeHTTP` implements `http.Handler` for testing with `httptest.NewRecorder()`
- `Run()` wraps `gin.Engine.Run()` for convenience; the server package (#38) may replace this later

### RouteGroup (`core/router/group.go`)

**Responsibility**: Wrap `*gin.RouterGroup`, provide identical HTTP method helpers for grouped routes
**Package**: `router`

```go
package router

import "github.com/gin-gonic/gin"

// RouteGroup wraps a Gin router group for sub-route registration.
type RouteGroup struct {
	group *gin.RouterGroup
}

// Get registers a GET route in the group.
func (g *RouteGroup) Get(path string, handlers ...gin.HandlerFunc) {
	g.group.GET(path, handlers...)
}

// Post registers a POST route in the group.
func (g *RouteGroup) Post(path string, handlers ...gin.HandlerFunc) {
	g.group.POST(path, handlers...)
}

// Put registers a PUT route in the group.
func (g *RouteGroup) Put(path string, handlers ...gin.HandlerFunc) {
	g.group.PUT(path, handlers...)
}

// Delete registers a DELETE route in the group.
func (g *RouteGroup) Delete(path string, handlers ...gin.HandlerFunc) {
	g.group.DELETE(path, handlers...)
}

// Patch registers a PATCH route in the group.
func (g *RouteGroup) Patch(path string, handlers ...gin.HandlerFunc) {
	g.group.PATCH(path, handlers...)
}

// Options registers an OPTIONS route in the group.
func (g *RouteGroup) Options(path string, handlers ...gin.HandlerFunc) {
	g.group.OPTIONS(path, handlers...)
}

// Group creates a nested sub-group with a shared prefix and optional middleware.
func (g *RouteGroup) Group(prefix string, handlers ...gin.HandlerFunc) *RouteGroup {
	return &RouteGroup{group: g.group.Group(prefix, handlers...)}
}

// Use adds middleware to the group.
func (g *RouteGroup) Use(middleware ...gin.HandlerFunc) {
	g.group.Use(middleware...)
}
```

**Design notes**:
- Mirrors Router's method set so groups and root have identical APIs
- Supports nesting: `g.Group("/sub")` returns another `*RouteGroup`
- `Use()` adds middleware only to this group and its children

### ResourceController Interface & Resource Routes (`core/router/resource.go`)

**Responsibility**: Define the RESTful controller contract and register all CRUD routes in one call
**Package**: `router`

```go
package router

import "github.com/gin-gonic/gin"

// ResourceController defines the interface for RESTful controllers.
// Implement all 7 methods for full CRUD with form routes,
// or use APIResource to skip Create/Edit form routes.
type ResourceController interface {
	Index(c *gin.Context)   // GET    /resource
	Create(c *gin.Context)  // GET    /resource/create  (SSR form)
	Store(c *gin.Context)   // POST   /resource
	Show(c *gin.Context)    // GET    /resource/:id
	Edit(c *gin.Context)    // GET    /resource/:id/edit (SSR form)
	Update(c *gin.Context)  // PUT    /resource/:id
	Destroy(c *gin.Context) // DELETE /resource/:id
}

// Resource registers all 7 RESTful routes for a controller on the router.
func (r *Router) Resource(path string, ctrl ResourceController) {
	r.engine.GET(path, ctrl.Index)
	r.engine.GET(path+"/create", ctrl.Create)
	r.engine.POST(path, ctrl.Store)
	r.engine.GET(path+"/:id", ctrl.Show)
	r.engine.GET(path+"/:id/edit", ctrl.Edit)
	r.engine.PUT(path+"/:id", ctrl.Update)
	r.engine.DELETE(path+"/:id", ctrl.Destroy)
}

// APIResource registers 5 RESTful routes (no Create/Edit form routes) on the router.
func (r *Router) APIResource(path string, ctrl ResourceController) {
	r.engine.GET(path, ctrl.Index)
	r.engine.POST(path, ctrl.Store)
	r.engine.GET(path+"/:id", ctrl.Show)
	r.engine.PUT(path+"/:id", ctrl.Update)
	r.engine.DELETE(path+"/:id", ctrl.Destroy)
}

// Resource registers all 7 RESTful routes for a controller on the group.
func (g *RouteGroup) Resource(path string, ctrl ResourceController) {
	g.group.GET(path, ctrl.Index)
	g.group.GET(path+"/create", ctrl.Create)
	g.group.POST(path, ctrl.Store)
	g.group.GET(path+"/:id", ctrl.Show)
	g.group.GET(path+"/:id/edit", ctrl.Edit)
	g.group.PUT(path+"/:id", ctrl.Update)
	g.group.DELETE(path+"/:id", ctrl.Destroy)
}

// APIResource registers 5 RESTful routes (no Create/Edit form routes) on the group.
func (g *RouteGroup) APIResource(path string, ctrl ResourceController) {
	g.group.GET(path, ctrl.Index)
	g.group.POST(path, ctrl.Store)
	g.group.GET(path+"/:id", ctrl.Show)
	g.group.PUT(path+"/:id", ctrl.Update)
	g.group.DELETE(path+"/:id", ctrl.Destroy)
}
```

**Design notes**:
- Blueprint defines `Resource()` as a standalone function taking `*gin.RouterGroup` — we make it a method on both `Router` and `RouteGroup` for API consistency
- `ResourceController` interface matched exactly from blueprint
- Form routes (Create, Edit) only in `Resource()`, not `APIResource()` — per blueprint

### Named Routes (`core/router/named.go`)

**Responsibility**: Thread-safe named route registry with URL generation via parameter substitution
**Package**: `router`

```go
package router

import (
	"strings"
	"sync"
)

var (
	namedRoutes = make(map[string]string)
	namedMu     sync.RWMutex
)

// Name registers a route name mapped to a path pattern.
func Name(name, pattern string) {
	namedMu.Lock()
	defer namedMu.Unlock()
	namedRoutes[name] = pattern
}

// Route generates a URL from a named route with parameter substitution.
// Parameters replace :param placeholders in registration order.
// Returns "/" if the route name is not found.
//
// Example: Route("users.show", "42") → "/users/42"
func Route(name string, params ...string) string {
	namedMu.RLock()
	pattern, ok := namedRoutes[name]
	namedMu.RUnlock()
	if !ok {
		return "/"
	}

	result := pattern
	for i := 0; i < len(params); i++ {
		idx := strings.Index(result, ":")
		if idx == -1 {
			break
		}
		end := strings.IndexAny(result[idx:], "/")
		if end == -1 {
			result = result[:idx] + params[i]
		} else {
			result = result[:idx] + params[i] + result[idx+end:]
		}
	}
	return result
}

// ResetNamedRoutes clears all named routes. Used in tests only.
func ResetNamedRoutes() {
	namedMu.Lock()
	defer namedMu.Unlock()
	namedRoutes = make(map[string]string)
}
```

**Design notes**:
- Global registry with `sync.RWMutex` — matched from blueprint
- `Route()` parameter substitution logic matches blueprint exactly
- Added `ResetNamedRoutes()` for test isolation — named routes are package-level globals
- No validation on duplicate names — last-write-wins (consistent with container pattern)

### RouterProvider (`app/providers/router_provider.go`)

**Responsibility**: Create the router, register it in the container, and trigger route registration
**Package**: `providers`

```go
package providers

import (
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/core/router"
	"github.com/RAiWorks/RGo/routes"
)

// RouterProvider creates the router and registers route definitions.
type RouterProvider struct{}

// Register creates a new Router and registers it as "router" in the container.
func (p *RouterProvider) Register(c *container.Container) {
	c.Instance("router", router.New())
}

// Boot loads route definitions from the routes package.
func (p *RouterProvider) Boot(c *container.Container) {
	r := container.MustMake[*router.Router](c, "router")
	routes.RegisterWeb(r)
	routes.RegisterAPI(r)
}
```

**Design notes**:
- `Instance()` not `Singleton()` — router is created once eagerly, not lazily
- `Register()` creates the router so other providers can access it during their `Boot()` phase
- `Boot()` calls route registration — routes may reference services that are only available after all providers register
- Route definitions in separate functions for clean organization

### Updated Routes Files

**`routes/web.go`**:
```go
package routes

import "github.com/RAiWorks/RGo/core/router"

// RegisterWeb defines web (HTML) routes.
func RegisterWeb(r *router.Router) {
	// Web routes will be added here.
	// Example:
	// r.Get("/", controllers.Home)
}
```

**`routes/api.go`**:
```go
package routes

import "github.com/RAiWorks/RGo/core/router"

// RegisterAPI defines API routes under the /api prefix.
func RegisterAPI(r *router.Router) {
	// api := r.Group("/api")
	// API routes will be added here.
	// Example:
	// api.Get("/users", controllers.ListUsers)
}
```

### Updated main.go (`cmd/main.go`)

**Changes**:
- Add `RouterProvider` registration after LoggerProvider
- Add HTTP server startup using `router.Run()`

```go
package main

import (
	"fmt"
	"log/slog"

	"github.com/RAiWorks/RGo/app/providers"
	"github.com/RAiWorks/RGo/core/app"
	"github.com/RAiWorks/RGo/core/config"
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/core/router"
)

func main() {
	application := app.New()

	// Register providers (order matters)
	application.Register(&providers.ConfigProvider{})  // 1. Config first — loads .env
	application.Register(&providers.LoggerProvider{})   // 2. Logger — uses config in Boot
	application.Register(&providers.RouterProvider{})   // 3. Router — creates Gin engine

	// Boot all providers
	application.Boot()

	appName := config.Env("APP_NAME", "RGo")
	appPort := config.Env("APP_PORT", "8080")
	appEnv := config.AppEnv()

	fmt.Println("=================================")
	fmt.Printf("  %s Framework\n", appName)
	fmt.Println("  github.com/RAiWorks/RGo")
	fmt.Println("=================================")
	fmt.Printf("  Environment: %s\n", appEnv)
	fmt.Printf("  Port: %s\n", appPort)
	fmt.Printf("  Debug: %v\n", config.IsDebug())
	fmt.Println("=================================")

	slog.Info("server starting",
		"app", appName,
		"port", appPort,
		"env", appEnv,
	)

	r := container.MustMake[*router.Router](application.Container, "router")
	if err := r.Run(":" + appPort); err != nil {
		slog.Error("server failed to start", "err", err)
	}
}
```

---

## Data Flow

```
main.go
  │
  ├── app.New()                          # Create app + container
  ├── Register(ConfigProvider)           # Loads .env
  ├── Register(LoggerProvider)           # (no-op in Register)
  ├── Register(RouterProvider)           # Creates Router, stores as "router"
  │     └── router.New()
  │           └── setGinMode()           # Reads APP_ENV → sets gin.ReleaseMode/etc
  │           └── gin.New()              # Creates bare Gin engine
  │
  ├── Boot()
  │     ├── ConfigProvider.Boot()        # no-op
  │     ├── LoggerProvider.Boot()        # logger.Setup()
  │     └── RouterProvider.Boot()        # Calls RegisterWeb + RegisterAPI
  │           ├── routes.RegisterWeb(r)  # Registers web routes
  │           └── routes.RegisterAPI(r)  # Registers API routes
  │
  └── r.Run(":8080")                     # Start HTTP server
```

---

## Dependency Graph

```
core/router/router.go    → config (AppEnv), gin
core/router/group.go     → gin
core/router/resource.go  → gin
core/router/named.go     → (no deps — stdlib only)
app/providers/router_provider.go → core/container, core/router, routes
routes/web.go            → core/router
routes/api.go            → core/router
cmd/main.go              → app/providers, core/app, core/config, core/container, core/router
```

No circular dependencies. `core/router` depends on `core/config` (to set Gin mode) and `gin`. Route files depend on `core/router`. Provider bridges container and routes.

---

## External Dependencies

| Package | Version | Purpose |
|---|---|---|
| `github.com/gin-gonic/gin` | latest | HTTP router engine |

This is the first major external dependency beyond godotenv. Gin will pull in its own dependencies (e.g., `go-playground/validator`, `ugorji/go/codec`).
