# 🧪 Test Plan: Views & Templates

> **Feature**: `17` — Views & Templates
> **Architecture**: [`17-views-templates-architecture.md`](17-views-templates-architecture.md)
> **Status**: ⬜ NOT RUN
> **Result**: ⬜ NOT RUN

---

## Test Files

- `core/router/view_test.go` — DefaultFuncMap and template rendering tests
- `http/controllers/controllers_test.go` — Updated Home controller test

---

## Test Cases

### TC-01: `TestDefaultFuncMap_ContainsRoute`
**What**: `DefaultFuncMap()` returns a FuncMap with a `route` key.
**How**: Call `DefaultFuncMap()`. Assert key `"route"` exists and is not nil.
**Pass**: FuncMap contains `route` function.

### TC-02: `TestDefaultFuncMap_RouteWorks`
**What**: The `route` function in the FuncMap resolves named routes.
**How**: Register a named route `"home"` → `"/"`. Get FuncMap. Call `route("home")`. Assert returns `"/"`.
**Pass**: Returns correct path.

### TC-03: `TestLoadTemplates_RenderHTML`
**What**: Templates loaded via `LoadTemplates` can be rendered with `c.HTML()`.
**How**: Create router. Set FuncMap. Load templates from a test fixtures dir. Register GET route that calls `c.HTML()`. Make request. Assert 200 with expected HTML content.
**Pass**: HTML body contains rendered template data.

### TC-04: `TestStatic_ServesFiles`
**What**: `Static` method serves files from a directory.
**How**: Create router. Create temp dir with a test file. Call `r.Static("/assets", tempDir)`. Make GET request to `/assets/test.txt`. Assert 200 with file content.
**Pass**: File content returned.

### TC-05: `TestStaticFile_ServesSingleFile`
**What**: `StaticFile` method serves a single file.
**How**: Create router. Create temp file. Call `r.StaticFile("/robots.txt", tempFile)`. Make GET request. Assert 200 with file content.
**Pass**: File content returned.

### TC-06: `TestHome_RendersHTML` (updated existing)
**What**: Home controller renders HTML template instead of JSON.
**How**: Create router with templates loaded. Register Home route. Make GET request. Assert 200 with Content-Type `text/html` and body contains "Welcome to RGo".
**Pass**: HTML response with expected content.

---

## Test Summary

| ID | Test Name | Type | File | Scope |
|---|---|---|---|---|
| TC-01 | `TestDefaultFuncMap_ContainsRoute` | Unit | `view_test.go` | FuncMap creation |
| TC-02 | `TestDefaultFuncMap_RouteWorks` | Unit | `view_test.go` | FuncMap route function |
| TC-03 | `TestLoadTemplates_RenderHTML` | Integration | `view_test.go` | Template rendering |
| TC-04 | `TestStatic_ServesFiles` | Integration | `view_test.go` | Static directory serving |
| TC-05 | `TestStaticFile_ServesSingleFile` | Integration | `view_test.go` | Static single file |
| TC-06 | `TestHome_RendersHTML` | Integration | `controllers_test.go` | Home controller HTML |

**Total new tests**: 5 (in `view_test.go`) + 1 updated (in `controllers_test.go`)
**Net new test count**: 5 new tests (1 is a replacement of existing)
**Expected total**: 172 + 5 = 177
