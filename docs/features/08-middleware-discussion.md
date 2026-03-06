# 💬 Discussion: Middleware Pipeline

> **Feature**: `08` — Middleware Pipeline
> **Status**: 🟢 COMPLETE
> **Branch**: `docs/08-middleware`
> **Depends On**: #07 (Router & Routing ✅)
> **Date Started**: 2026-03-06
> **Date Completed**: 2026-03-06

---

## Summary

Implement the middleware pipeline for the RGo framework. This feature provides a middleware registry with named aliases and groups, a set of built-in middleware (Recovery, RequestID, CORS, Error Handler), and integration with the router's `Use()` infrastructure. The registry enables middleware to be referenced by name and applied to route groups in bulk, following the patterns established in the blueprint.

---

## Functional Requirements

- As a **framework developer**, I want a middleware registry (`RegisterAlias`, `Resolve`) so that middleware can be referenced by string name throughout the application
- As a **framework developer**, I want middleware groups (`RegisterGroup`, `ResolveGroup`) so that bundles of middleware can be applied to route groups in a single call
- As a **framework developer**, I want a `Recovery` middleware so that panics are caught and converted to 500 responses without crashing the server
- As a **framework developer**, I want a `RequestID` middleware so that every request gets a unique identifier in the `X-Request-ID` header for tracing
- As a **framework developer**, I want a `CORS` middleware so that cross-origin requests are handled with configurable allowed origins, methods, and headers
- As a **framework developer**, I want an `ErrorHandler` middleware so that `AppError` values set during request handling are formatted into proper JSON responses
- As a **framework user**, I want to register custom middleware aliases at boot time so my application-specific middleware integrates with the framework registry
- As a **framework user**, I want to define custom middleware groups so I can apply different middleware stacks to "web" vs "api" route groups

## Current State / Reference

### What Exists
- **Router (#07 ✅)**: `Router.Use()` and `RouteGroup.Use()` accept `...gin.HandlerFunc` — middleware can already be applied to routes
- **Error Handling (#04 ✅)**: `AppError` with HTTP status codes, `ErrorResponse()` with debug-aware formatting
- **Config (#02 ✅)**: `config.AppEnv()`, `config.IsDebug()`, `config.IsTesting()`, `config.IsProduction()`
- **`core/middleware/`**: Empty directory with `.gitkeep`
- **Gin engine**: Created via `gin.New()` — no default middleware (no Recovery, no Logger)

### Blueprint Reference
The blueprint shows:
1. `AuthMiddleware()` — checks Authorization header, aborts with 401 if missing
2. Middleware registry — `RegisterAlias()`, `RegisterGroup()`, `Resolve()`, `ResolveGroup()` with global maps
3. Registration at boot — built-in aliases (`auth`, `csrf`, `cors`, `rate`, `requestid`) and groups (`web`, `api`)
4. Usage in routes — `middleware.ResolveGroup("web")...` spread into `r.Group()` calls
5. Error middleware — `r.Use(middleware.ErrorHandler())` referenced in config section
6. Conditional middleware — `if !config.IsTesting() { r.Use(middleware.RateLimitMiddleware()) }`

### What the Blueprint Does NOT Show
- How Recovery middleware replaces `gin.Recovery()` with framework-controlled behavior
- How RequestID middleware generates and propagates IDs
- Exact CORS configuration options or defaults
- How ErrorHandler middleware reads `AppError` from the Gin context
- How the middleware registry interacts with the provider lifecycle

## Proposed Approach

### Middleware Registry (`core/middleware/registry.go`)

Central registry with two maps: aliases (single middleware) and groups (middleware bundles).

```go
var (
    routeMiddleware   = map[string]gin.HandlerFunc{}
    middlewareGroups  = map[string][]gin.HandlerFunc{}
)
```

**Functions:**
- `RegisterAlias(name, handler)` — store a middleware under a string key
- `RegisterGroup(name, ...handlers)` — store a slice of middleware under a string key
- `Resolve(name) gin.HandlerFunc` — retrieve by alias (panics if not found)
- `ResolveGroup(name) []gin.HandlerFunc` — retrieve group (returns nil if not found)
- `ResetRegistry()` — clear both maps for test isolation

### Recovery Middleware (`core/middleware/recovery.go`)

Catches panics, logs the stack trace via `slog`, and returns a 500 JSON response. Replaces `gin.Recovery()` with framework-integrated behavior.

```go
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                slog.Error("panic recovered", "error", err, "path", c.Request.URL.Path)
                c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
            }
        }()
        c.Next()
    }
}
```

### RequestID Middleware (`core/middleware/request_id.go`)

Generates or propagates a UUID for each request. Checks for an existing `X-Request-ID` header first.

```go
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.GetHeader("X-Request-ID")
        if id == "" {
            id = generateID()  // UUID v4
        }
        c.Set("request_id", id)
        c.Header("X-Request-ID", id)
        c.Next()
    }
}
```

Uses `crypto/rand` for UUID generation — no external dependency needed.

### CORS Middleware (`core/middleware/cors.go`)

Configurable CORS with sensible defaults. Uses a `CORSConfig` struct:

```go
type CORSConfig struct {
    AllowOrigins []string
    AllowMethods []string
    AllowHeaders []string
    MaxAge       int  // seconds
}
```

- `CORS(config ...CORSConfig) gin.HandlerFunc` — creates CORS middleware with optional custom config
- Default: allow all origins (`*`), standard methods, common headers, 12-hour max age
- Handles preflight `OPTIONS` requests by aborting with 204 after setting headers

### ErrorHandler Middleware (`core/middleware/error_handler.go`)

Processes errors set on the Gin context. After `c.Next()`, checks for `AppError` in `c.Errors` and formats the response using `ErrorResponse()`.

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) == 0 {
            return
        }
        // Check for AppError and format response
    }
}
```

### MiddlewareProvider (`app/providers/middleware_provider.go`)

Registers built-in middleware aliases and default groups at boot time.

**Register**: no-op (no singleton needed in container)
**Boot**: Registers built-in aliases (`recovery`, `requestid`, `cors`, `error_handler`) and default groups

## Scoping Decisions

### IN Scope
| Item | Rationale |
|---|---|
| Middleware registry (aliases + groups) | Blueprint core — enables middleware-by-name |
| `RegisterAlias` / `Resolve` | Blueprint core |
| `RegisterGroup` / `ResolveGroup` | Blueprint core |
| Recovery middleware | Required — `gin.New()` has no recovery; server must not crash on panic |
| RequestID middleware | Blueprint lists as built-in |
| CORS middleware | Blueprint lists as built-in |
| ErrorHandler middleware | Blueprint references; uses AppError from #04 |
| MiddlewareProvider | Natural extension of provider pattern from #06 |
| `ResetRegistry()` for test cleanup | Test isolation |

### NOT in Scope
| Item | Rationale | Future Feature |
|---|---|---|
| Auth middleware (`AuthMiddleware`) | Depends on JWT/auth system (#20) | #20 |
| CSRF middleware | Depends on session system (#19) | #19 |
| Rate limiting middleware | Depends on cache/storage system | #28 |
| Session middleware | Separate feature | #19 |
| Gin Logger middleware replacement | Framework uses `slog` already (#03); request logging may come later | #29 |
| `make:middleware` CLI generator | Depends on CLI foundation (#10) | #10 |
| Route model binding middleware (`BindModel`) | Depends on GORM/Database (#09, #11) | #11 or #15 |

### Key Decision: Built-in Middleware Set (4 middleware)
We implement exactly **4 built-in middleware** for this feature:
1. **Recovery** — essential for server stability (replaces gin.Recovery)
2. **RequestID** — essential for tracing/debugging
3. **CORS** — common API requirement, blueprint-listed
4. **ErrorHandler** — connects AppError (#04) to HTTP responses

Auth, CSRF, Rate Limiting, and Session middleware are deferred to their respective features which provide the underlying functionality they depend on.

### Key Decision: UUID Generation with `crypto/rand`
Use `crypto/rand` for UUID v4 generation instead of adding a UUID library dependency. The implementation is ~10 lines and avoids a new external dependency for a simple operation.

### Key Decision: CORS Config Struct vs Environment Variables
CORS configuration uses a Go struct with defaults, not environment variables. This allows type-safe configuration and programmatic customization. Environment-based CORS config can be added later when the config system supports struct binding.

## Edge Cases & Risks

- [x] Recovery must not double-write response if handler already wrote before panicking — use `c.IsAborted()` check
- [x] RequestID should preserve caller-provided IDs (for distributed tracing chains)
- [x] CORS preflight must return 204 and abort (not call `c.Next()`)
- [x] ErrorHandler must run after all other handlers (`c.Next()` first, then check errors)
- [x] Middleware registry `Resolve()` panics on unknown name — same as blueprint's explicit behavior
- [x] `ResolveGroup()` returns nil for unknown names — graceful, allows `...` spread without panic
- [x] Registry maps are not thread-safe — registration happens at boot time (single goroutine), resolution happens during request handling (read-only). This is safe without mutex since maps are only written during boot before serving starts.
- [x] ErrorHandler should handle both `*AppError` and generic `error` types from Gin's error list
