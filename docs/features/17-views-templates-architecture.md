# 🏗️ Architecture: Views & Templates

> **Feature**: `17` — Views & Templates
> **Discussion**: [`17-views-templates-discussion.md`](17-views-templates-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

Feature #17 adds server-side rendering support to the framework by wiring Go's `html/template` engine into the Router, registering template functions (including the `route` helper), configuring static file serving, providing a sample home template, and updating the Home controller to render HTML.

## File Structure

```
core/router/
├── router.go           # + LoadTemplates, SetFuncMap, Static, StaticFile methods
├── view.go             # DefaultFuncMap() with route helper
http/controllers/
├── home_controller.go  # Updated: c.HTML() instead of c.JSON()
app/providers/
├── router_provider.go  # Updated: template + static setup in Boot
resources/views/
├── home.html           # Sample home template
```

### Files Created (2)
| File | Package | Lines (est.) |
|---|---|---|
| `core/router/view.go` | `router` | ~20 |
| `resources/views/home.html` | — | ~15 |

### Files Modified (3)
| File | Change |
|---|---|
| `core/router/router.go` | Add `LoadTemplates`, `SetFuncMap`, `Static`, `StaticFile` methods |
| `http/controllers/home_controller.go` | Use `c.HTML()` with `home.html` template |
| `app/providers/router_provider.go` | Add template + static setup in `Boot` |

---

## Component Design

### Template Function Map (`core/router/view.go`)

**Responsibility**: Provide the default template function map with the `route` helper.
**Package**: `router`

```go
package router

import "html/template"

// DefaultFuncMap returns the template function map with framework helpers.
func DefaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"route": Route,
	}
}
```

**Design notes**:
- `Route` function already exists in `core/router/named.go`
- Separate file keeps view-related code organized
- Returns `template.FuncMap` (standard Go type) — no custom types needed

### Router Methods (`core/router/router.go` additions)

**Responsibility**: Thin wrappers exposing Gin's template and static methods.

```go
// SetFuncMap sets the template function map on the Gin engine.
// Must be called before LoadTemplates.
func (r *Router) SetFuncMap(funcMap template.FuncMap) {
	r.engine.SetFuncMap(funcMap)
}

// LoadTemplates loads HTML templates matching the glob pattern.
func (r *Router) LoadTemplates(pattern string) {
	r.engine.LoadHTMLGlob(pattern)
}

// Static serves files from a local directory under the given URL path.
func (r *Router) Static(urlPath, dirPath string) {
	r.engine.Static(urlPath, dirPath)
}

// StaticFile serves a single file at the given URL path.
func (r *Router) StaticFile(urlPath, filePath string) {
	r.engine.StaticFile(urlPath, filePath)
}
```

**Design notes**:
- `SetFuncMap` must be called **before** `LoadTemplates` (Gin requirement)
- Consistent with existing Router wrapper pattern (Get, Post, etc.)
- Import `html/template` added to router.go

### Router Provider Update (`app/providers/router_provider.go`)

**Responsibility**: Wire template engine and static serving during boot.

```go
func (p *RouterProvider) Boot(c *container.Container) {
	r := container.MustMake[*router.Router](c, "router")

	// Template engine setup
	r.SetFuncMap(router.DefaultFuncMap())
	r.LoadTemplates("resources/views/**/*")

	// Static file serving
	r.Static("/static", "./resources/static")
	r.Static("/uploads", "./storage/uploads")

	// Route definitions
	routes.RegisterWeb(r)
	routes.RegisterAPI(r)
}
```

**Design notes**:
- `SetFuncMap` → `LoadTemplates` order is correct (Gin requires FuncMap first)
- Static paths match blueprint and framework doc
- Removed `StaticFile("/favicon.ico", ...)` — no favicon.ico exists yet; can be added when assets are created

### Home Controller Update (`http/controllers/home_controller.go`)

```go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Home renders the home page using the home.html template.
func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "Welcome to RGo",
	})
}
```

**Design notes**:
- Changes from `c.JSON()` to `c.HTML()` with template name and data
- Passes `title` as template data
- Template name `home.html` matches file in `resources/views/`

### Sample Template (`resources/views/home.html`)

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .title }}</title>
</head>
<body>
    <h1>{{ .title }}</h1>
    <p>RGo Framework is running.</p>
</body>
</html>
```

**Design notes**:
- Minimal HTML5 template proving the pipeline works
- Uses `{{ .title }}` data binding
- No external CSS/JS dependencies — keeps it self-contained

---

## Data Flow

### Template Rendering
```
Request GET / → Router → Home controller
  → c.HTML(200, "home.html", gin.H{"title": "Welcome to RGo"})
  → Gin renders html/template with data
  → HTML response to client
```

### Static File Serving
```
Request GET /static/css/app.css → Gin Static handler
  → Serves ./resources/static/css/app.css
```

---

## Dependencies

| Dependency | Type | Usage |
|---|---|---|
| `html/template` | stdlib | FuncMap type in router.go |
| `core/router/named.go` | internal | `Route` function exposed as template func |
| Gin `engine.SetFuncMap` | framework | Template function registration |
| Gin `engine.LoadHTMLGlob` | framework | Template loading |
| Gin `engine.Static` | framework | Static file serving |

---

## Impact on Existing Code

| Component | Impact |
|---|---|
| `core/router/router.go` | 4 new methods added (non-breaking) |
| `http/controllers/home_controller.go` | `c.JSON()` → `c.HTML()` (changes response type) |
| `app/providers/router_provider.go` | Template + static setup added to Boot (non-breaking) |
| `http/controllers/controllers_test.go` | Home test needs update (JSON → HTML assertion) |
| `resources/views/` | `.gitkeep` stays; `home.html` added |
