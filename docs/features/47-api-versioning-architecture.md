# 📐 Architecture: API Versioning

> **Feature**: `47` — API Versioning
> **Discussion**: [`47-api-versioning-discussion.md`](47-api-versioning-discussion.md)
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## Overview

Feature #47 adds two methods to the `Router` struct in `core/router/version.go`: `Version()` creates a route group prefixed with `/api/{version}`, and `DeprecatedVersion()` does the same but automatically injects middleware that sets `Sunset` and `X-API-Deprecated` response headers on every request. Both methods return the existing `*RouteGroup` type, so all existing route registration methods (Get, Post, APIResource, Group, Use, etc.) work seamlessly on versioned groups.

---

## File Structure

```
core/
  router/
    version.go       ← NEW — Version(), DeprecatedVersion(), deprecation middleware
    version_test.go  ← NEW — unit tests (10 tests)
```

No existing files are modified. No migrations. No model changes.

---

## Component Design

### 1. Version Method

**Package**: `core/router`
**File**: `version.go`

```go
// Version creates a route group prefixed with /api/{version}.
// The returned RouteGroup supports all existing route registration methods.
//
// Example:
//
//	v1 := r.Version("v1") // prefix: /api/v1
//	v1.Get("/users", listUsers)
//	v1.APIResource("/posts", &PostController{})
func (r *Router) Version(version string) *RouteGroup {
	return r.Group("/api/" + version)
}
```

| Parameter | Type | Purpose |
|---|---|---|
| `version` | `string` | Version identifier (e.g. "v1", "v2", "2024-01") |
| **Returns** | `*RouteGroup` | Route group at `/api/{version}` |

---

### 2. DeprecatedVersion Method

```go
// DeprecatedVersion creates a versioned route group (like Version) but injects
// middleware that adds deprecation headers to every response:
//   - Sunset: {sunsetDate} — RFC 8594 sunset date in HTTP-date format
//   - X-API-Deprecated: true — simple boolean signal for clients/monitoring
//
// Example:
//
//	v1 := r.DeprecatedVersion("v1", "Sat, 01 Jun 2026 00:00:00 GMT")
//	v1.Get("/users", listUsersV1)
func (r *Router) DeprecatedVersion(version, sunsetDate string) *RouteGroup {
	g := r.Version(version)
	g.Use(deprecationHeaders(sunsetDate))
	return g
}
```

| Parameter | Type | Purpose |
|---|---|---|
| `version` | `string` | Version identifier (e.g. "v1") |
| `sunsetDate` | `string` | Sunset date in HTTP-date format (RFC 7231 §7.1.1.1) |
| **Returns** | `*RouteGroup` | Route group at `/api/{version}` with deprecation middleware |

---

### 3. Deprecation Middleware

```go
// deprecationHeaders returns middleware that sets Sunset and X-API-Deprecated
// headers on every response.
func deprecationHeaders(sunsetDate string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Sunset", sunsetDate)
		c.Header("X-API-Deprecated", "true")
		c.Next()
	}
}
```

**Behavior**:
- Sets `Sunset` header with the provided date string (not validated — the developer is responsible for correct HTTP-date format)
- Sets `X-API-Deprecated: true` as a simple boolean marker
- Calls `c.Next()` to continue the middleware chain
- Unexported function — not part of the public API

---

## Public API Summary

| Function | Signature | Purpose |
|---|---|---|
| `Version` | `(r *Router) Version(version string) *RouteGroup` | Create a versioned API route group |
| `DeprecatedVersion` | `(r *Router) DeprecatedVersion(version, sunsetDate string) *RouteGroup` | Create a deprecated versioned API route group |

---

## Usage Example

```go
// routes/api.go
package routes

import (
	"github.com/RAiWorks/RapidGo/core/router"
	"github.com/RAiWorks/RapidGo/http/controllers"
)

func RegisterAPI(r *router.Router) {
	// Current version
	v2 := r.Version("v2")
	v2.APIResource("/posts", &controllers.PostController{})
	v2.Get("/users", controllers.ListUsersV2)

	// Deprecated version — sunset June 1, 2026
	v1 := r.DeprecatedVersion("v1", "Sat, 01 Jun 2026 00:00:00 GMT")
	v1.Get("/users", controllers.ListUsersV1)
}
```

Requests to `/api/v1/users` will include:
```
Sunset: Sat, 01 Jun 2026 00:00:00 GMT
X-API-Deprecated: true
```

Requests to `/api/v2/users` will have no deprecation headers.

---

## Dependencies

None. Uses only existing framework packages (`core/router`, `github.com/gin-gonic/gin`).

---

## Environment Variables

None.

---

## Next

Tasks → `47-api-versioning-tasks.md`
