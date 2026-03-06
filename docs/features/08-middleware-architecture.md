# 🏗️ Architecture: Middleware Pipeline

> **Feature**: `08` — Middleware Pipeline
> **Discussion**: [`08-middleware-discussion.md`](08-middleware-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

The Middleware Pipeline provides a registry for naming and grouping middleware, plus four built-in middleware functions: Recovery, RequestID, CORS, and ErrorHandler. All middleware are standard `gin.HandlerFunc` values, compatible with the Router's existing `Use()` and `Group()` methods. A `MiddlewareProvider` integrates the registry into the service provider lifecycle.

## File Structure

```
core/middleware/
├── registry.go         # RegisterAlias, RegisterGroup, Resolve, ResolveGroup
├── recovery.go         # Recovery middleware — catches panics
├── request_id.go       # RequestID middleware — adds X-Request-ID
├── cors.go             # CORS middleware — handles cross-origin requests
└── error_handler.go    # ErrorHandler middleware — formats AppError responses

core/middleware/
└── middleware_test.go   # Tests for all middleware + registry

app/providers/
└── middleware_provider.go  # MiddlewareProvider — registers aliases & groups
```

### Files Created (6)
| File | Package | Lines (est.) |
|---|---|---|
| `core/middleware/registry.go` | `middleware` | ~50 |
| `core/middleware/recovery.go` | `middleware` | ~25 |
| `core/middleware/request_id.go` | `middleware` | ~35 |
| `core/middleware/cors.go` | `middleware` | ~60 |
| `core/middleware/error_handler.go` | `middleware` | ~35 |
| `app/providers/middleware_provider.go` | `providers` | ~30 |

### Files Modified (0)
No existing files need modification. The Router's `Use()` method already supports `gin.HandlerFunc`, and `MiddlewareProvider` is a new provider added alongside existing ones.

---

## Component Design

### Middleware Registry (`core/middleware/registry.go`)

**Responsibility**: Store and retrieve middleware by name. Aliases map a string to a single `gin.HandlerFunc`; groups map a string to a slice of `gin.HandlerFunc`.
**Package**: `middleware`

```go
package middleware

import "github.com/gin-gonic/gin"

var (
	routeMiddleware  = map[string]gin.HandlerFunc{}
	middlewareGroups = map[string][]gin.HandlerFunc{}
)

// RegisterAlias registers a named middleware that can be referenced by string.
func RegisterAlias(name string, handler gin.HandlerFunc) {
	routeMiddleware[name] = handler
}

// RegisterGroup registers a named group of middleware.
func RegisterGroup(name string, handlers ...gin.HandlerFunc) {
	middlewareGroups[name] = handlers
}

// Resolve returns a middleware handler by alias name.
// Panics if the alias is not registered.
func Resolve(name string) gin.HandlerFunc {
	if h, ok := routeMiddleware[name]; ok {
		return h
	}
	panic("middleware not found: " + name)
}

// ResolveGroup returns all middleware in a named group.
// Returns nil if the group is not registered.
func ResolveGroup(name string) []gin.HandlerFunc {
	if g, ok := middlewareGroups[name]; ok {
		return g
	}
	return nil
}

// ResetRegistry clears all registered aliases and groups. For testing only.
func ResetRegistry() {
	routeMiddleware = map[string]gin.HandlerFunc{}
	middlewareGroups = map[string][]gin.HandlerFunc{}
}
```

**Design notes**:
- Maps are written only at boot time (single goroutine) and read during request handling (concurrent but read-only). No mutex needed.
- `Resolve` panics on missing alias — this is a programming error (misconfigured routes), not a runtime condition. Matches blueprint behavior.
- `ResolveGroup` returns nil for missing groups — allows safe `...` spread: `r.Group("/api", middleware.ResolveGroup("api")...)`
- `ResetRegistry` prevents test pollution across test functions.

### Recovery Middleware (`core/middleware/recovery.go`)

**Responsibility**: Catch panics during request handling, log the error, return a 500 JSON response.
**Package**: `middleware`

```go
package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery returns middleware that catches panics and returns a 500 error.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered",
					"error", err,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
				)

				if !c.Writer.Written() {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": "internal server error",
					})
				} else {
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}
```

**Design notes**:
- Uses `slog.Error` — integrates with framework's logging (#03)
- Checks `c.Writer.Written()` — if headers already sent, can't write JSON response, just abort
- Logs method + path for debugging context
- Replaces `gin.Recovery()` which writes to stdout with non-framework formatting

### RequestID Middleware (`core/middleware/request_id.go`)

**Responsibility**: Assign a unique identifier to every request. Propagate existing IDs for distributed tracing.
**Package**: `middleware`

```go
package middleware

import (
	"crypto/rand"
	"fmt"

	"github.com/gin-gonic/gin"
)

const requestIDHeader = "X-Request-ID"

// RequestID returns middleware that assigns a unique ID to each request.
// If the request already has an X-Request-ID header, it is preserved.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = generateUUID()
		}
		c.Set("request_id", id)
		c.Header(requestIDHeader, id)
		c.Next()
	}
}

// generateUUID produces a UUID v4 string using crypto/rand.
func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
```

**Design notes**:
- `crypto/rand` for cryptographically secure UUIDs — no external dependency
- Preserves incoming `X-Request-ID` for distributed tracing chains
- Sets both the Gin context (`c.Set`) and response header (`c.Header`)
- Header constant avoids typos across middleware

### CORS Middleware (`core/middleware/cors.go`)

**Responsibility**: Handle Cross-Origin Resource Sharing headers and preflight requests.
**Package**: `middleware`

```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds configuration for the CORS middleware.
type CORSConfig struct {
	AllowOrigins []string // Default: ["*"]
	AllowMethods []string // Default: ["GET","POST","PUT","DELETE","PATCH","OPTIONS"]
	AllowHeaders []string // Default: ["Origin","Content-Type","Accept","Authorization","X-Request-ID"]
	MaxAge       int      // Default: 43200 (12 hours), in seconds
}

// defaultCORSConfig returns the default CORS configuration.
func defaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		MaxAge:       43200,
	}
}

// CORS returns middleware that handles cross-origin requests.
// If no config is provided, sensible defaults are used.
func CORS(configs ...CORSConfig) gin.HandlerFunc {
	cfg := defaultCORSConfig()
	if len(configs) > 0 {
		cfg = configs[0]
	}

	origins := strings.Join(cfg.AllowOrigins, ", ")
	methods := strings.Join(cfg.AllowMethods, ", ")
	headers := strings.Join(cfg.AllowHeaders, ", ")
	maxAge := fmt.Sprintf("%d", cfg.MaxAge)

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origins)
		c.Header("Access-Control-Allow-Methods", methods)
		c.Header("Access-Control-Allow-Headers", headers)
		c.Header("Access-Control-Max-Age", maxAge)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
```

**Design notes**:
- Optional config parameter — zero-arg call uses safe defaults
- Pre-joins header strings outside the handler for performance (computed once, not per-request)
- Preflight `OPTIONS` returns 204 and aborts — no further handler execution
- `X-Request-ID` included in default allowed headers for RequestID interop
- Import `"fmt"` needed for `Sprintf` on `MaxAge`

### ErrorHandler Middleware (`core/middleware/error_handler.go`)

**Responsibility**: After all handlers run, check Gin's error list for `AppError` values and format the JSON response.
**Package**: `middleware`

```go
package middleware

import (
	"net/http"

	"github.com/RAiWorks/RGo/core/errors"
	"github.com/gin-gonic/gin"
)

// ErrorHandler returns middleware that processes errors added to the Gin context.
// If an error is an *errors.AppError, it uses the status code and ErrorResponse().
// Other errors are treated as 500 Internal Server Error.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		// Use the last error (most specific)
		lastErr := c.Errors.Last().Err

		if appErr, ok := lastErr.(*errors.AppError); ok {
			c.JSON(appErr.Code, appErr.ErrorResponse())
			return
		}

		// Generic error — wrap as 500
		wrapped := errors.Internal(lastErr)
		c.JSON(wrapped.Code, wrapped.ErrorResponse())
	}
}
```

**Design notes**:
- Calls `c.Next()` first — this ensures all downstream handlers run before error processing
- Uses `c.Errors.Last()` — in a chain of middleware, the most specific error is typically last
- Leverages `AppError.ErrorResponse()` — debug-aware formatting from Feature #04 (includes internal details when `APP_DEBUG=true`)
- Fallback to `errors.Internal(err)` for unexpected error types — always returns valid JSON

### MiddlewareProvider (`app/providers/middleware_provider.go`)

**Responsibility**: Register built-in middleware aliases and default groups at boot time.
**Package**: `providers`

```go
package providers

import (
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/core/middleware"
)

// MiddlewareProvider registers built-in middleware aliases and groups.
type MiddlewareProvider struct{}

// Register is a no-op — middleware has no singleton to register.
func (p *MiddlewareProvider) Register(c *container.Container) {}

// Boot registers built-in middleware aliases and default groups.
func (p *MiddlewareProvider) Boot(c *container.Container) {
	// Built-in aliases
	middleware.RegisterAlias("recovery", middleware.Recovery())
	middleware.RegisterAlias("requestid", middleware.RequestID())
	middleware.RegisterAlias("cors", middleware.CORS())
	middleware.RegisterAlias("error_handler", middleware.ErrorHandler())

	// Default groups
	middleware.RegisterGroup("global",
		middleware.Recovery(),
		middleware.RequestID(),
	)
}
```

**Design notes**:
- `Register()` is a no-op — middleware doesn't need a container singleton
- `Boot()` runs after all providers register, so router is available
- Default `"global"` group includes Recovery + RequestID — minimum viable middleware stack
- Auth, CSRF, Rate Limiting aliases are NOT registered here — deferred to their features
- Users can add custom aliases in their own providers or in route files

---

## Data Flow

### Request Lifecycle with Middleware

```
HTTP Request
  → Recovery middleware (defer/recover)
    → RequestID middleware (assign/propagate ID)
      → CORS middleware (set headers, handle preflight)
        → Route handler
      ← ErrorHandler middleware (process c.Errors)
    ← RequestID (response header already set)
  ← Recovery (catch panics from anywhere in chain)
HTTP Response
```

### Middleware Registration Flow

```
App.Boot()
  → MiddlewareProvider.Boot()
    → RegisterAlias("recovery", Recovery())
    → RegisterAlias("requestid", RequestID())
    → RegisterAlias("cors", CORS())
    → RegisterAlias("error_handler", ErrorHandler())
    → RegisterGroup("global", Recovery(), RequestID())
  → RouterProvider.Boot()
    → routes.RegisterWeb(r) — user applies middleware via Resolve/ResolveGroup
    → routes.RegisterAPI(r) — user applies middleware via Resolve/ResolveGroup
```

### Provider Order

The `MiddlewareProvider` should be registered **before** `RouterProvider` so that middleware aliases are available when routes are defined in `Boot()`:

```go
application.Register(&providers.ConfigProvider{})      // 1
application.Register(&providers.LoggerProvider{})       // 2
application.Register(&providers.MiddlewareProvider{})   // 3 — register aliases
application.Register(&providers.RouterProvider{})       // 4 — routes can use Resolve()
```

---

## Dependencies

### Internal
| Package | Used For |
|---|---|
| `core/errors` | `AppError` type in ErrorHandler |
| `core/config` | `IsDebug()` via `AppError.ErrorResponse()` |
| `log/slog` | Panic logging in Recovery |
| `crypto/rand` | UUID generation in RequestID |

### External
None — no new external dependencies. All middleware use only the standard library and `gin-gonic/gin` (already imported by #07).

---

## Contract & Interface

### Middleware Type
All middleware are `gin.HandlerFunc`:
```go
type HandlerFunc func(*gin.Context)
```

### Registry API
```go
func RegisterAlias(name string, handler gin.HandlerFunc)
func RegisterGroup(name string, handlers ...gin.HandlerFunc)
func Resolve(name string) gin.HandlerFunc       // panics on unknown
func ResolveGroup(name string) []gin.HandlerFunc // nil on unknown
func ResetRegistry()                            // test cleanup
```

### CORS Config
```go
type CORSConfig struct {
    AllowOrigins []string
    AllowMethods []string
    AllowHeaders []string
    MaxAge       int
}
```

### ErrorHandler Contract
Reads errors via `c.Errors`. Handlers signal errors using:
```go
c.Error(errors.NotFound("user not found"))
c.Abort()
```

---

## Testing Strategy

- **Unit tests per middleware**: Each middleware tested in isolation with `httptest.NewRecorder`
- **Registry tests**: Register/Resolve round-trip, panic on unknown, group operations
- **Integration test**: Multiple middleware in chain, verify execution order
- **Provider test**: Compile-time interface check + verify aliases registered after Boot

Tests use the same `httptest` + Gin engine pattern established in Feature #07's router tests.
