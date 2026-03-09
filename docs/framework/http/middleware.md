---
title: "Middleware"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Middleware

## Abstract

This document covers the middleware system — the central registry for
aliases and groups, built-in middleware (auth, CSRF, CORS, rate
limiting, request ID, sessions), and creating custom middleware.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Middleware Registry](#2-middleware-registry)
3. [Built-in Aliases](#3-built-in-aliases)
4. [Built-in Groups](#4-built-in-groups)
5. [Registering at Boot Time](#5-registering-at-boot-time)
6. [Using Middleware in Routes](#6-using-middleware-in-routes)
7. [Auth Middleware](#7-auth-middleware)
8. [Session Middleware](#8-session-middleware)
9. [Request ID Middleware](#9-request-id-middleware)
10. [Custom Middleware](#10-custom-middleware)
11. [Execution Order](#11-execution-order)
12. [Security Considerations](#12-security-considerations)
13. [References](#13-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Middleware** — A function that intercepts requests before and/or
  after the handler runs.
- **Alias** — A string name mapped to a single middleware handler.
- **Group** — A string name mapped to an ordered list of middleware.

## 2. Middleware Registry

The registry provides a central place to define middleware aliases and
groups, similar to Laravel's `$routeMiddleware` and
`$middlewareGroups`:

```go
package middleware

import "github.com/gin-gonic/gin"

var (
    // Route-level middleware (use by alias)
    routeMiddleware = map[string]gin.HandlerFunc{}

    // Middleware groups
    middlewareGroups = map[string][]gin.HandlerFunc{}
)

// RegisterAlias registers a named middleware that can be referenced
// by string.
func RegisterAlias(name string, handler gin.HandlerFunc) {
    routeMiddleware[name] = handler
}

// RegisterGroup registers a named group of middleware.
func RegisterGroup(name string, handlers ...gin.HandlerFunc) {
    middlewareGroups[name] = handlers
}

// Resolve returns a middleware handler by alias name.
func Resolve(name string) gin.HandlerFunc {
    if h, ok := routeMiddleware[name]; ok {
        return h
    }
    panic("middleware not found: " + name)
}

// ResolveGroup returns all middleware in a group.
func ResolveGroup(name string) []gin.HandlerFunc {
    if g, ok := middlewareGroups[name]; ok {
        return g
    }
    return nil
}
```

## 3. Built-in Aliases

| Alias | Middleware | Purpose |
|-------|-----------|---------|
| `auth` | `AuthMiddleware()` | JWT/session authentication |
| `csrf` | `CSRFMiddleware()` | CSRF token validation |
| `cors` | `CORSMiddleware()` | Cross-Origin Resource Sharing |
| `rate` | `RateLimitMiddleware()` | Request rate limiting |
| `requestid` | `RequestIDMiddleware()` | Unique request ID tracing |

## 4. Built-in Groups

### `web` Group

For server-rendered HTML routes:

| Order | Middleware | Purpose |
|-------|-----------|---------|
| 1 | `SessionMiddleware` | Load/save session per request |
| 2 | `CSRFMiddleware` | Token generation and validation |
| 3 | `RequestIDMiddleware` | Attach request ID |

### `api` Group

For JSON API routes:

| Order | Middleware | Purpose |
|-------|-----------|---------|
| 1 | `CORSMiddleware` | Handle preflight and CORS headers |
| 2 | `RateLimitMiddleware` | Throttle requests |
| 3 | `RequestIDMiddleware` | Attach request ID |

## 5. Registering at Boot Time

Register all aliases and groups during application bootstrap:

```go
// Built-in aliases
middleware.RegisterAlias("auth", AuthMiddleware())
middleware.RegisterAlias("csrf", CSRFMiddleware())
middleware.RegisterAlias("cors", CORSMiddleware())
middleware.RegisterAlias("rate", RateLimitMiddleware())
middleware.RegisterAlias("requestid", RequestIDMiddleware())

// User-defined aliases
middleware.RegisterAlias("admin", myapp.AdminOnlyMiddleware())
middleware.RegisterAlias("verified", myapp.EmailVerifiedMiddleware())

// Groups
middleware.RegisterGroup("web",
    SessionMiddleware(sessionMgr),
    CSRFMiddleware(),
    RequestIDMiddleware(),
)
middleware.RegisterGroup("api",
    CORSMiddleware(),
    RateLimitMiddleware(),
    RequestIDMiddleware(),
)
```

## 6. Using Middleware in Routes

### Apply a Group

```go
web := r.Group("/", middleware.ResolveGroup("web")...)
api := r.Group("/api", middleware.ResolveGroup("api")...)
```

### Apply by Alias

```go
admin := web.Group("/admin",
    middleware.Resolve("auth"),
    middleware.Resolve("admin"),
)
```

### Inline

```go
r.GET("/users/:id",
    middleware.BindModel(db, "user", &models.User{}),
    controllers.ShowUser,
)
```

## 7. Auth Middleware

Validates JWT tokens from the `Authorization` header:

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")

        if header == "" || !strings.HasPrefix(header, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "missing or invalid Authorization header",
            })
            return
        }

        // Validate JWT token and set user in context
        // token := strings.TrimPrefix(header, "Bearer ")
        // claims, err := auth.ValidateToken(token)

        c.Next()
    }
}
```

## 8. Session Middleware

Automatically loads and saves session data per request:

```go
func SessionMiddleware(mgr *session.Manager) gin.HandlerFunc {
    return func(c *gin.Context) {
        id, data, err := mgr.Start(c.Request)
        if err != nil {
            c.AbortWithStatus(500)
            return
        }

        c.Set("session_id", id)
        c.Set("session", data)

        c.Next()

        // Persist session after the handler runs
        updated, _ := c.Get("session")
        mgr.Save(c.Writer, id, updated.(map[string]interface{}))
    }
}
```

## 9. Request ID Middleware

Attaches a unique request ID for tracing across logs:

```go
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.GetHeader("X-Request-ID")
        if id == "" {
            b := make([]byte, 16)
            rand.Read(b)
            id = hex.EncodeToString(b)
        }
        c.Set("request_id", id)
        c.Header("X-Request-ID", id)
        c.Next()
    }
}
```

If the client sends an `X-Request-ID` header, the framework reuses
it. Otherwise, it generates a cryptographically random 32-character
hex string.

## 10. Custom Middleware

Create custom middleware following the Gin handler pattern:

```go
func AdminOnlyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check user role from context (set by auth middleware)
        role, exists := c.Get("user_role")
        if !exists || role != "admin" {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": "admin access required",
            })
            return
        }
        c.Next()
    }
}
```

Register as an alias:

```go
middleware.RegisterAlias("admin", AdminOnlyMiddleware())
```

## 11. Execution Order

Middleware executes in the order it is registered. For nested groups,
outer middleware runs first:

```text
Request
  → Group middleware (web/api)
    → Route-level middleware (auth, admin)
      → Controller handler
    ← Route-level middleware (after c.Next())
  ← Group middleware (after c.Next())
Response
```

## 12. Security Considerations

- Auth middleware **MUST** run before any handler that requires
  authentication.
- CSRF middleware **MUST** be applied to all web routes that handle
  state-changing requests (POST, PUT, DELETE).
- Rate limiting **SHOULD** be applied to all API routes and
  authentication endpoints.
- Middleware **MUST NOT** leak internal error details in production
  responses.

## 13. References

- [CSRF Protection](../security/csrf.md)
- [CORS](../security/cors.md)
- [Rate Limiting](../security/rate-limiting.md)
- [Authentication](../security/authentication.md)
- [Sessions](../security/sessions.md)
- [Routing](routing.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
