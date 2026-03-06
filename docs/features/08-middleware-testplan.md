# 🧪 Test Plan: Middleware Pipeline

> **Feature**: `08` — Middleware Pipeline
> **Tasks**: [`08-middleware-tasks.md`](08-middleware-tasks.md)
> **Date**: 2026-03-06

---

## Acceptance Criteria

- [ ] Middleware registry supports alias registration and resolution
- [ ] Middleware registry supports group registration and resolution
- [ ] `Resolve()` panics for unknown aliases
- [ ] `ResolveGroup()` returns nil for unknown groups
- [ ] `ResetRegistry()` clears all registered aliases and groups
- [ ] Recovery middleware catches panics and returns 500 JSON
- [ ] Recovery middleware handles already-written responses
- [ ] RequestID middleware assigns UUID to requests
- [ ] RequestID middleware preserves incoming X-Request-ID
- [ ] CORS middleware sets correct headers with defaults
- [ ] CORS middleware accepts custom configuration
- [ ] CORS middleware handles preflight OPTIONS with 204
- [ ] ErrorHandler middleware formats AppError as JSON
- [ ] ErrorHandler middleware wraps generic errors as 500
- [ ] ErrorHandler middleware is a no-op when no errors exist
- [ ] MiddlewareProvider implements Provider interface
- [ ] MiddlewareProvider registers built-in aliases on Boot
- [ ] All tests pass with `go test ./core/middleware/...`
- [ ] `go vet ./...` reports no issues

---

## Test Cases

### TC-01: RegisterAlias and Resolve round-trip

**File**: `core/middleware/middleware_test.go`
**Function**: `TestRegisterAlias_AndResolve`

| Step | Action | Expected |
|---|---|---|
| 1 | `ResetRegistry()` | Clean state |
| 2 | `RegisterAlias("test", handler)` | No error |
| 3 | `Resolve("test")` | Returns the registered handler (not nil) |

---

### TC-02: Resolve panics on unknown alias

**File**: `core/middleware/middleware_test.go`
**Function**: `TestResolve_PanicsOnUnknown`

| Step | Action | Expected |
|---|---|---|
| 1 | `ResetRegistry()` | Clean state |
| 2 | `Resolve("nonexistent")` inside `defer recover` | Panics with "middleware not found: nonexistent" |

---

### TC-03: RegisterGroup and ResolveGroup round-trip

**File**: `core/middleware/middleware_test.go`
**Function**: `TestRegisterGroup_AndResolveGroup`

| Step | Action | Expected |
|---|---|---|
| 1 | `ResetRegistry()` | Clean state |
| 2 | `RegisterGroup("web", h1, h2)` | No error |
| 3 | `ResolveGroup("web")` | Returns slice of length 2 |

---

### TC-04: ResolveGroup returns nil for unknown group

**File**: `core/middleware/middleware_test.go`
**Function**: `TestResolveGroup_ReturnsNilOnUnknown`

| Step | Action | Expected |
|---|---|---|
| 1 | `ResetRegistry()` | Clean state |
| 2 | `ResolveGroup("nonexistent")` | Returns nil |

---

### TC-05: ResetRegistry clears all entries

**File**: `core/middleware/middleware_test.go`
**Function**: `TestResetRegistry_ClearsAll`

| Step | Action | Expected |
|---|---|---|
| 1 | `RegisterAlias("a", handler)` | Registered |
| 2 | `RegisterGroup("g", handler)` | Registered |
| 3 | `ResetRegistry()` | Both cleared |
| 4 | `ResolveGroup("g")` | Returns nil |
| 5 | `Resolve("a")` inside `defer recover` | Panics |

---

### TC-06: Recovery catches panic and returns 500

**File**: `core/middleware/middleware_test.go`
**Function**: `TestRecovery_CatchesPanic`

| Step | Action | Expected |
|---|---|---|
| 1 | Create Gin engine with `Recovery()` middleware | Engine created |
| 2 | Register route that calls `panic("test panic")` | Route registered |
| 3 | Send `GET` request | Response status 500 |
| 4 | Check response body | Contains `"error": "internal server error"` |

---

### TC-07: Recovery passes through normal requests

**File**: `core/middleware/middleware_test.go`
**Function**: `TestRecovery_PassesNormalRequest`

| Step | Action | Expected |
|---|---|---|
| 1 | Create Gin engine with `Recovery()` middleware | Engine created |
| 2 | Register route that returns 200 with `"ok"` | Route registered |
| 3 | Send `GET` request | Response status 200 |

---

### TC-08: RequestID generates UUID when no header present

**File**: `core/middleware/middleware_test.go`
**Function**: `TestRequestID_GeneratesUUID`

| Step | Action | Expected |
|---|---|---|
| 1 | Create Gin engine with `RequestID()` middleware | Engine created |
| 2 | Register route that reads `c.GetString("request_id")` | Route registered |
| 3 | Send `GET` request without `X-Request-ID` header | Response status 200 |
| 4 | Check `X-Request-ID` response header | Not empty, matches UUID format |
| 5 | Check response body | Contains the same ID |

---

### TC-09: RequestID preserves incoming X-Request-ID

**File**: `core/middleware/middleware_test.go`
**Function**: `TestRequestID_PreservesExisting`

| Step | Action | Expected |
|---|---|---|
| 1 | Create Gin engine with `RequestID()` middleware | Engine created |
| 2 | Register route that reads `c.GetString("request_id")` | Route registered |
| 3 | Send `GET` request with `X-Request-ID: my-trace-id-123` | Response status 200 |
| 4 | Check `X-Request-ID` response header | `"my-trace-id-123"` |
| 5 | Check response body | Contains `"my-trace-id-123"` |

---

### TC-10: CORS sets default headers

**File**: `core/middleware/middleware_test.go`
**Function**: `TestCORS_DefaultHeaders`

| Step | Action | Expected |
|---|---|---|
| 1 | Create Gin engine with `CORS()` middleware (no config) | Engine created |
| 2 | Register route returning 200 | Route registered |
| 3 | Send `GET` request | Response status 200 |
| 4 | Check `Access-Control-Allow-Origin` header | `"*"` |
| 5 | Check `Access-Control-Allow-Methods` header | Contains `"GET"`, `"POST"`, etc. |
| 6 | Check `Access-Control-Allow-Headers` header | Contains `"Authorization"`, `"Content-Type"` |

---

### TC-11: CORS handles preflight OPTIONS with 204

**File**: `core/middleware/middleware_test.go`
**Function**: `TestCORS_PreflightOptions`

| Step | Action | Expected |
|---|---|---|
| 1 | Create Gin engine with `CORS()` middleware | Engine created |
| 2 | Register `GET` route | Route registered |
| 3 | Send `OPTIONS` request to same path | Response status 204 |
| 4 | Check CORS headers are set | All present |

---

### TC-12: CORS accepts custom configuration

**File**: `core/middleware/middleware_test.go`
**Function**: `TestCORS_CustomConfig`

| Step | Action | Expected |
|---|---|---|
| 1 | Create `CORSConfig{AllowOrigins: ["https://example.com"]}` | Config created |
| 2 | Create Gin engine with `CORS(config)` middleware | Engine created |
| 3 | Register route returning 200 | Route registered |
| 4 | Send `GET` request | Response status 200 |
| 5 | Check `Access-Control-Allow-Origin` header | `"https://example.com"` |

---

### TC-13: ErrorHandler formats AppError as JSON

**File**: `core/middleware/middleware_test.go`
**Function**: `TestErrorHandler_FormatsAppError`

| Step | Action | Expected |
|---|---|---|
| 1 | Create Gin engine with `ErrorHandler()` middleware | Engine created |
| 2 | Register route that calls `c.Error(errors.NotFound("user not found"))` and `c.Abort()` | Route registered |
| 3 | Send `GET` request | Response status 404 |
| 4 | Check response body | Contains `"error": "user not found"` |

---

### TC-14: ErrorHandler wraps generic error as 500

**File**: `core/middleware/middleware_test.go`
**Function**: `TestErrorHandler_WrapsGenericError`

| Step | Action | Expected |
|---|---|---|
| 1 | Create Gin engine with `ErrorHandler()` middleware | Engine created |
| 2 | Register route that calls `c.Error(fmt.Errorf("something broke"))` and `c.Abort()` | Route registered |
| 3 | Send `GET` request | Response status 500 |
| 4 | Check response body | Contains `"error": "internal server error"` |

---

### TC-15: ErrorHandler is no-op when no errors

**File**: `core/middleware/middleware_test.go`
**Function**: `TestErrorHandler_NoOpWhenNoErrors`

| Step | Action | Expected |
|---|---|---|
| 1 | Create Gin engine with `ErrorHandler()` middleware | Engine created |
| 2 | Register route that returns 200 normally | Route registered |
| 3 | Send `GET` request | Response status 200 |
| 4 | Check response body | Contains handler's normal response |

---

### TC-16: MiddlewareProvider implements Provider interface

**File**: `app/providers/providers_test.go`
**Function**: Compile-time check: `var _ container.Provider = (*MiddlewareProvider)(nil)`

| Step | Action | Expected |
|---|---|---|
| 1 | Compile-time `var _ container.Provider = (*MiddlewareProvider)(nil)` | Compiles without error |

---

### TC-17: MiddlewareProvider Boot registers built-in aliases

**File**: `app/providers/providers_test.go`
**Function**: `TestMiddlewareProvider_RegistersAliases`

| Step | Action | Expected |
|---|---|---|
| 1 | `middleware.ResetRegistry()` | Clean state |
| 2 | Create container, call `MiddlewareProvider.Boot(c)` | No error |
| 3 | `middleware.Resolve("recovery")` | Not nil |
| 4 | `middleware.Resolve("requestid")` | Not nil |
| 5 | `middleware.Resolve("cors")` | Not nil |
| 6 | `middleware.Resolve("error_handler")` | Not nil |
| 7 | `middleware.ResolveGroup("global")` | Length 2 |

---

### TC-18: Middleware chain executes in correct order

**File**: `core/middleware/middleware_test.go`
**Function**: `TestMiddlewareChain_ExecutionOrder`

| Step | Action | Expected |
|---|---|---|
| 1 | Create Gin engine with Recovery + RequestID + test handler | Engine created |
| 2 | Send `GET` request | Response status 200 |
| 3 | Check `X-Request-ID` header is present | UUID set |
| 4 | Verify no panic recovery triggered (normal flow) | Clean 200 |

---

### TC-19: UUID format validation

**File**: `core/middleware/middleware_test.go`
**Function**: `TestGenerateUUID_Format`

| Step | Action | Expected |
|---|---|---|
| 1 | Call `generateUUID()` multiple times | Returns strings |
| 2 | Check each is 36 chars with hyphens at positions 8, 13, 18, 23 | Valid UUID format |
| 3 | Check version nibble (char at position 14) | `'4'` |
| 4 | Check two UUIDs are different | Unique |

---

## Test File Summary

| File | Test Functions | Test Cases |
|---|---|---|
| `core/middleware/middleware_test.go` | 17 | TC-01 through TC-15, TC-18, TC-19 |
| `app/providers/providers_test.go` | 2 (added) | TC-16, TC-17 |

**Total**: 19 test functions covering 19 test cases
