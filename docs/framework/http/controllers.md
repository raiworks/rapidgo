---
title: "Controllers"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Controllers

## Abstract

Controllers handle incoming HTTP requests and return responses. This
document covers the MVC controller pattern, struct-based resource
controllers, the controller → service delegation pattern, and
accessing framework services from the container.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Basic Controllers](#2-basic-controllers)
3. [Struct-based Controllers](#3-struct-based-controllers)
4. [Controller → Service Pattern](#4-controller--service-pattern)
5. [Request Data Access](#5-request-data-access)
6. [Rendering Responses](#6-rendering-responses)
7. [Container Service Access](#7-container-service-access)
8. [Security Considerations](#8-security-considerations)
9. [References](#9-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Controller** — A function or method that receives a Gin context,
  processes the request, and returns a response.
- **Resource controller** — A struct implementing the
  `ResourceController` interface for full CRUD operations.

## 2. Basic Controllers

Controllers live in `http/controllers/`. The simplest form is a plain
function:

```go
package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
    c.HTML(http.StatusOK, "home.html", gin.H{
        "title": "Welcome",
    })
}

func Users(c *gin.Context) {
    users := []string{"Alice", "Bob"}

    c.JSON(http.StatusOK, gin.H{
        "data": users,
    })
}
```

Register in routes:

```go
r.GET("/", controllers.Home)
r.GET("/api/users", controllers.Users)
```

## 3. Struct-based Controllers

For resource routes, controllers **SHOULD** be structs implementing
the `ResourceController` interface:

```go
package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type PostController struct{}

func (ctrl *PostController) Index(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "PostController index"})
}

func (ctrl *PostController) Create(c *gin.Context) {
    // Render create form (SSR only)
    c.HTML(http.StatusOK, "posts/create.html", gin.H{})
}

func (ctrl *PostController) Store(c *gin.Context) {
    c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (ctrl *PostController) Show(c *gin.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, gin.H{"id": id})
}

func (ctrl *PostController) Edit(c *gin.Context) {
    id := c.Param("id")
    c.HTML(http.StatusOK, "posts/edit.html", gin.H{"id": id})
}

func (ctrl *PostController) Update(c *gin.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, gin.H{"id": id, "message": "updated"})
}

func (ctrl *PostController) Destroy(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
```

Register with resource routes:

```go
router.Resource(&r.RouterGroup, "/posts", &controllers.PostController{})
router.APIResource(api, "/posts", &controllers.PostController{})
```

## 4. Controller → Service Pattern

Controllers **SHOULD NOT** contain business logic. They delegate to
**services** for domain operations:

```go
func CreateUser(c *gin.Context) {
    var req requests.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        responses.Error(c, 422, err.Error())
        return
    }

    userSvc := services.NewUserService(db)
    user, err := userSvc.Create(req.Name, req.Email, req.Password)
    if err != nil {
        responses.Error(c, 400, err.Error())
        return
    }
    responses.Created(c, user)
}
```

### Responsibilities

| Layer | Responsibility |
|-------|---------------|
| **Controller** | Parse request, call service, format response |
| **Service** | Business logic, validation rules, orchestration |
| **Model** | Data schema, relationships, hooks |

### Pattern Summary

```text
Request → Controller → Service → Model → DB
                                       ↓
Response ← Controller ← Service ← Result
```

## 5. Request Data Access

Gin provides several methods for reading request data:

```go
// URL parameters
id := c.Param("id")          // /users/:id

// Query strings
page := c.Query("page")      // ?page=2
page := c.DefaultQuery("page", "1")

// Form data
name := c.PostForm("name")

// JSON body
var req CreateUserRequest
c.ShouldBindJSON(&req)

// File upload
file, header, err := c.Request.FormFile("file")
```

## 6. Rendering Responses

### HTML (SSR)

```go
func Show(c *gin.Context) {
    c.HTML(http.StatusOK, "users/show.html", gin.H{
        "title": "User Profile",
        "user":  user,
    })
}
```

### JSON (API)

Use the response helpers for consistent API envelopes:

```go
func Show(c *gin.Context) {
    responses.Success(c, user)   // 200 + success envelope
}

func Store(c *gin.Context) {
    responses.Created(c, user)   // 201 + success envelope
}

func Store(c *gin.Context) {
    responses.Error(c, 422, "validation failed") // error envelope
}
```

### Redirect

```go
c.Redirect(http.StatusFound, router.Route("users.show", fmt.Sprint(user.ID)))
```

## 7. Container Service Access

Resolve services from the container within controllers:

```go
func Dashboard(c *gin.Context) {
    // Via the global app container
    cache := container.MustMake[cache.Store](app.Container, "cache")
    mailer := container.MustMake[*mail.Mailer](app.Container, "mail")

    // Use resolved services
    val, _ := cache.Get("stats")
    // ...
}
```

For frequently used services, inject them into the controller struct
or resolve once at setup time rather than on every request.

## 8. Security Considerations

- Controllers **MUST NOT** trust user input. Always validate before
  processing.
- Passwords **MUST** be hashed before storage — never stored in
  plain text.
- JSON responses **MUST NOT** include sensitive fields. Use
  `json:"-"` tags on model fields like `Password`.
- Controllers behind auth routes **MUST** verify the user has
  permission for the requested action.

## 9. References

- [Routing](routing.md)
- [Views](views.md)
- [Responses](responses.md)
- [Requests & Validation](requests-validation.md)
- [Services Layer](../infrastructure/services-layer.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
