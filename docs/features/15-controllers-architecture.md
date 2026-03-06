# 🏗️ Architecture: Controllers

> **Feature**: `15` — Controllers
> **Discussion**: [`15-controllers-discussion.md`](15-controllers-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

Feature #15 adds MVC controllers to `http/controllers/`. A `Home` function controller handles the root route, and a `PostController` struct implements the `ResourceController` interface for RESTful CRUD. Routes are registered in `routes/web.go` and `routes/api.go`.

## File Structure

```
http/controllers/
├── home_controller.go     # Home() function controller
└── post_controller.go     # PostController struct (ResourceController)

routes/
├── web.go                 # (modified) Register Home route
└── api.go                 # (modified) Register PostController routes
```

### Files Created (2)
| File | Package | Lines (est.) |
|---|---|---|
| `http/controllers/home_controller.go` | `controllers` | ~15 |
| `http/controllers/post_controller.go` | `controllers` | ~50 |

### Files Modified (2)
| File | Change |
|---|---|
| `routes/web.go` | Register `GET /` → `controllers.Home` |
| `routes/api.go` | Register `APIResource("/posts", &PostController{})` |

---

## Component Design

### Home Controller (`http/controllers/home_controller.go`)

**Responsibility**: Handle the root route with a welcome response.
**Package**: `controllers`

```go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Home handles the root route and returns a welcome message.
func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to RGo",
	})
}
```

**Design notes**:
- Simple function handler — demonstrates the basic controller pattern
- Returns JSON since Views (#17) aren't built yet
- Blueprint shows `c.HTML()` — will be updated when templates ship

### Post Controller (`http/controllers/post_controller.go`)

**Responsibility**: RESTful resource controller for posts.
**Package**: `controllers`

```go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PostController handles CRUD operations for posts.
type PostController struct{}

func (ctrl *PostController) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PostController index"})
}

func (ctrl *PostController) Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PostController create form"})
}

func (ctrl *PostController) Store(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "PostController store"})
}

func (ctrl *PostController) Show(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "PostController show", "id": id})
}

func (ctrl *PostController) Edit(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "PostController edit form", "id": id})
}

func (ctrl *PostController) Update(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "PostController update", "id": id})
}

func (ctrl *PostController) Destroy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PostController destroy"})
}
```

**Design notes**:
- Implements all 7 `ResourceController` interface methods
- Placeholder responses — no database interaction yet (services layer is #18)
- `c.Param("id")` for URL parameters — standard Gin pattern
- `Store` returns `http.StatusCreated` (201), all others return 200
- Matches blueprint template output exactly

### Route Registration (`routes/web.go`)

```go
package routes

import (
	"github.com/RAiWorks/RGo/core/router"
	"github.com/RAiWorks/RGo/http/controllers"
)

func RegisterWeb(r *router.Router) {
	r.Get("/", controllers.Home)
}
```

### Route Registration (`routes/api.go`)

```go
package routes

import (
	"github.com/RAiWorks/RGo/core/router"
	"github.com/RAiWorks/RGo/http/controllers"
)

func RegisterAPI(r *router.Router) {
	api := r.Group("/api")
	api.APIResource("/posts", &controllers.PostController{})
}
```

---

## Data Flow

### `GET /`
```
Request → Router → controllers.Home(c)
       → c.JSON(200, {"message": "Welcome to RGo"})
       → Response
```

### `GET /api/posts`
```
Request → Router → api group → PostController.Index(c)
       → c.JSON(200, {"message": "PostController index"})
       → Response
```

### `GET /api/posts/:id`
```
Request → Router → api group → PostController.Show(c)
       → c.Param("id") → c.JSON(200, {"message": "...", "id": id})
       → Response
```

---

## Constraints & Invariants

1. Controllers **return JSON only** — no HTML until Views (#17) ships
2. Controllers **contain no business logic** — placeholder responses for now
3. `PostController` implements **all 7** `ResourceController` methods (interface compliance)
4. Routes use **existing Router API** — `r.Get()`, `api.APIResource()` — no new router code
5. Controller package has **no database imports** — pure HTTP layer
