---
title: "Routing"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Routing

## Abstract

This document covers the framework's routing system built on Gin —
basic route definition, route groups, resource routes (RESTful
convention), named routes with URL generation, and route model binding.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Router Setup](#2-router-setup)
3. [Basic Routes](#3-basic-routes)
4. [Route Groups](#4-route-groups)
5. [Resource Routes](#5-resource-routes)
6. [Named Routes & URL Generation](#6-named-routes--url-generation)
7. [Route Model Binding](#7-route-model-binding)
8. [Route Files](#8-route-files)
9. [Security Considerations](#9-security-considerations)
10. [References](#10-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Resource route** — A single function call that registers all 7
  RESTful CRUD routes for a controller.
- **Named route** — A route with a string identifier used for URL
  generation in templates and redirects.
- **Route model binding** — Middleware that automatically loads a
  GORM model from the `:id` URL parameter.

## 2. Router Setup

The framework uses **Gin** (`github.com/gin-gonic/gin`) as the HTTP
router.

```go
package router

import (
    "yourframework/http/controllers"
    "yourframework/core/middleware"

    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // Web routes
    r.GET("/", controllers.Home)

    // API routes
    api := r.Group("/api", middleware.AuthMiddleware())
    {
        api.GET("/users", controllers.Users)
    }

    return r
}
```

## 3. Basic Routes

Register routes using HTTP method functions:

```go
r.GET("/path", handler)
r.POST("/path", handler)
r.PUT("/path", handler)
r.PATCH("/path", handler)
r.DELETE("/path", handler)
r.OPTIONS("/path", handler)
```

With URL parameters:

```go
r.GET("/users/:id", controllers.ShowUser)      // required param
r.GET("/files/*filepath", controllers.ServeFile) // catch-all param
```

## 4. Route Groups

Group routes to share prefixes and middleware:

```go
// API group — CORS + rate limiting + request ID
api := r.Group("/api", middleware.ResolveGroup("api")...)
{
    api.GET("/users", controllers.ListUsers)
    api.POST("/users", controllers.CreateUser)
}

// Admin group — auth + admin middleware
admin := r.Group("/admin", middleware.Resolve("auth"), middleware.Resolve("admin"))
{
    admin.GET("/dashboard", controllers.AdminDashboard)
}
```

## 5. Resource Routes

### ResourceController Interface

Define all 7 CRUD actions in a single interface:

```go
// ResourceController defines the interface for RESTful controllers.
type ResourceController interface {
    Index(c *gin.Context)   // GET    /resource
    Create(c *gin.Context)  // GET    /resource/create  (SSR form)
    Store(c *gin.Context)   // POST   /resource
    Show(c *gin.Context)    // GET    /resource/:id
    Edit(c *gin.Context)    // GET    /resource/:id/edit (SSR form)
    Update(c *gin.Context)  // PUT    /resource/:id
    Destroy(c *gin.Context) // DELETE /resource/:id
}
```

### Resource (Full CRUD)

Registers all 7 routes — including `Create` and `Edit` form routes
for SSR applications:

```go
func Resource(group *gin.RouterGroup, path string, ctrl ResourceController) {
    group.GET(path, ctrl.Index)
    group.GET(path+"/create", ctrl.Create)
    group.POST(path, ctrl.Store)
    group.GET(path+"/:id", ctrl.Show)
    group.GET(path+"/:id/edit", ctrl.Edit)
    group.PUT(path+"/:id", ctrl.Update)
    group.DELETE(path+"/:id", ctrl.Destroy)
}
```

### APIResource (API only)

Registers 5 routes — without `Create` and `Edit` since APIs don't
serve HTML forms:

```go
func APIResource(group *gin.RouterGroup, path string, ctrl ResourceController) {
    group.GET(path, ctrl.Index)
    group.POST(path, ctrl.Store)
    group.GET(path+"/:id", ctrl.Show)
    group.PUT(path+"/:id", ctrl.Update)
    group.DELETE(path+"/:id", ctrl.Destroy)
}
```

### Usage

```go
// Web — full CRUD with create/edit forms
router.Resource(&r.RouterGroup, "/posts", &controllers.PostController{})

// API — no form routes
router.APIResource(api, "/users", &controllers.UserController{})
```

### Generated Routes

For `Resource(group, "/posts", ctrl)`:

| Method | URI | Action | Handler |
|--------|-----|--------|---------|
| GET | `/posts` | Index | `ctrl.Index` |
| GET | `/posts/create` | Create form | `ctrl.Create` |
| POST | `/posts` | Store | `ctrl.Store` |
| GET | `/posts/:id` | Show | `ctrl.Show` |
| GET | `/posts/:id/edit` | Edit form | `ctrl.Edit` |
| PUT | `/posts/:id` | Update | `ctrl.Update` |
| DELETE | `/posts/:id` | Destroy | `ctrl.Destroy` |

## 6. Named Routes & URL Generation

### Registering Named Routes

```go
var (
    namedRoutes = make(map[string]string)
    mu          sync.RWMutex
)

func Name(name, pattern string) {
    mu.Lock()
    defer mu.Unlock()
    namedRoutes[name] = pattern
}
```

Register names alongside route definitions:

```go
router.Name("home", "/")
router.Name("users.index", "/users")
router.Name("users.show", "/users/:id")
router.Name("posts.edit", "/posts/:id/edit")
```

### Generating URLs

The `Route()` function generates URLs from named routes with
parameter substitution:

```go
func Route(name string, params ...string) string {
    mu.RLock()
    pattern, ok := namedRoutes[name]
    mu.RUnlock()
    if !ok {
        return "/"
    }

    i := 0
    result := pattern
    for i < len(params) {
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
        i++
    }
    return result
}
```

**Examples:**

```go
router.Route("home")                    // "/"
router.Route("users.show", "42")        // "/users/42"
router.Route("posts.edit", "7")         // "/posts/7/edit"
```

### Template Usage

Pass `Route` as a template function:

```html
<a href="{{route `users.show` .User.ID}}">View Profile</a>
<a href="{{route `posts.edit` .Post.ID}}">Edit Post</a>
```

### Redirect Usage

```go
c.Redirect(http.StatusFound, router.Route("users.show", fmt.Sprint(user.ID)))
```

## 7. Route Model Binding

Automatically resolve a GORM model from the `:id` route parameter,
similar to Laravel's implicit model binding:

```go
func BindModel(db *gorm.DB, key string, model interface{}) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        result := db.First(model, id)
        if result.Error != nil {
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
                "error": "resource not found",
            })
            return
        }
        c.Set(key, model)
        c.Next()
    }
}
```

### Usage

Register as middleware on a route:

```go
r.GET("/users/:id",
    middleware.BindModel(db, "user", &models.User{}),
    controllers.ShowUser,
)
```

Retrieve the model in the controller:

```go
func ShowUser(c *gin.Context) {
    user, _ := c.Get("user")
    responses.Success(c, user)
}
```

The model is already loaded — no database query needed in the
controller.

## 8. Route Files

Routes **SHOULD** be organized into separate files:

### `routes/web.go`

Web routes use the `web` middleware group (session, CSRF, request ID):

```go
web := r.Group("/", middleware.ResolveGroup("web")...)
{
    web.GET("/", controllers.Home)
    router.Resource(web, "/posts", &controllers.PostController{})
    router.Resource(web, "/users", &controllers.UserController{})
}
```

### `routes/api.go`

API routes use the `api` middleware group (CORS, rate limit, request ID):

```go
api := r.Group("/api", middleware.ResolveGroup("api")...)
{
    api.POST("/login", controllers.Login)
    router.APIResource(api, "/users", &controllers.UserController{})
    router.APIResource(api, "/posts", &controllers.PostController{})
}
```

## 9. Security Considerations

- Route model binding uses `db.First(model, id)` which parameterizes
  the query — safe from SQL injection.
- Named routes use parameter substitution on pre-registered patterns,
  **not** user input as route patterns.
- API routes **SHOULD** always include rate limiting middleware.
- Admin routes **MUST** be protected by authentication and
  authorization middleware.

## 10. References

- [Controllers](controllers.md)
- [Middleware](middleware.md)
- [Views](views.md)
- [Request Lifecycle Diagram](../architecture/diagrams/request-lifecycle.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
