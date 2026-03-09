---
title: "Building a REST API"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Building a REST API

## Abstract

A guide to building a JSON REST API with JWT authentication, request
validation, pagination, CORS, and rate limiting.

## Table of Contents

1. [Overview](#1-overview)
2. [Create Model and Service](#2-create-model-and-service)
3. [Create API Controller](#3-create-api-controller)
4. [Register API Routes](#4-register-api-routes)
5. [Add JWT Authentication](#5-add-jwt-authentication)
6. [Add Request Validation](#6-add-request-validation)
7. [Add Pagination](#7-add-pagination)
8. [Add CORS and Rate Limiting](#8-add-cors-and-rate-limiting)
9. [References](#9-references)

## 1. Overview

API routes use the `api` middleware group (rate limiting + request ID)
and return JSON via the `APIResponse` helpers.

## 2. Create Model and Service

```bash
framework make:model Article
framework make:service ArticleService
```

Model:

```go
type Article struct {
    BaseModel
    Title  string `gorm:"size:200;not null" json:"title"`
    Body   string `gorm:"type:text;not null" json:"body"`
    UserID uint   `gorm:"not null" json:"user_id"`
}
```

Service with full CRUD — see [Services Layer](../infrastructure/services-layer.md).

## 3. Create API Controller

```go
package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "yourframework/app/services"
    "yourframework/core/responses"
    "yourframework/database/models"
)

type ArticleAPIController struct {
    Service *services.ArticleService
}

func (ctrl *ArticleAPIController) Index(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))
    articles, total, _ := ctrl.Service.GetAll(page, perPage)

    responses.Paginated(c, articles, int(total), page, perPage)
}

func (ctrl *ArticleAPIController) Show(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    article, err := ctrl.Service.GetByID(uint(id))
    if err != nil {
        responses.Error(c, http.StatusNotFound, "Article not found")
        return
    }
    responses.Success(c, article)
}

func (ctrl *ArticleAPIController) Store(c *gin.Context) {
    var input models.Article
    if err := c.ShouldBindJSON(&input); err != nil {
        responses.Error(c, http.StatusUnprocessableEntity, err.Error())
        return
    }
    if err := ctrl.Service.Create(&input); err != nil {
        responses.Error(c, http.StatusInternalServerError, "Failed to create")
        return
    }
    responses.Created(c, input)
}

func (ctrl *ArticleAPIController) Update(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    article, err := ctrl.Service.GetByID(uint(id))
    if err != nil {
        responses.Error(c, http.StatusNotFound, "Article not found")
        return
    }
    if err := c.ShouldBindJSON(article); err != nil {
        responses.Error(c, http.StatusUnprocessableEntity, err.Error())
        return
    }
    ctrl.Service.Update(article)
    responses.Success(c, article)
}

func (ctrl *ArticleAPIController) Destroy(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    ctrl.Service.Delete(uint(id))
    responses.Success(c, nil)
}
```

## 4. Register API Routes

In `routes/api.go`:

```go
api := r.Group("/api/v1")
api.Use(middleware.Resolve("rate"), middleware.Resolve("requestid"))
{
    // Public
    api.POST("/login", authCtrl.Login)

    // Protected
    auth := api.Group("/")
    auth.Use(middleware.Resolve("auth"))
    {
        router.APIResource(auth, "/articles", articleCtrl)
    }
}
```

`APIResource` registers all routes except the HTML form routes
(`create`, `edit`).

## 5. Add JWT Authentication

Login endpoint returns a token:

```go
func (ctrl *AuthController) Login(c *gin.Context) {
    var input struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        responses.Error(c, http.StatusUnprocessableEntity, err.Error())
        return
    }

    user, err := ctrl.UserService.FindByEmail(input.Email)
    if err != nil || !helpers.CheckPassword(user.Password, input.Password) {
        responses.Error(c, http.StatusUnauthorized, "Invalid credentials")
        return
    }

    token, _ := auth.GenerateToken(user.ID, user.Email)
    responses.Success(c, gin.H{"token": token})
}
```

Clients send the token:

```yaml
Authorization: Bearer <token>
```

See [Authentication](../security/authentication.md).

## 6. Add Request Validation

Use struct-based validation with binding tags:

```go
type CreateArticleRequest struct {
    Title string `json:"title" binding:"required,min=3,max=200"`
    Body  string `json:"body" binding:"required"`
}

func (ctrl *ArticleAPIController) Store(c *gin.Context) {
    var input CreateArticleRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        responses.Error(c, http.StatusUnprocessableEntity, err.Error())
        return
    }
    // ... create article from input
}
```

See [Requests & Validation](../http/requests-validation.md).

## 7. Add Pagination

Use the `Paginate` helper and `responses.Paginated`:

```go
func (ctrl *ArticleAPIController) Index(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))

    result := database.Paginate(db, page, perPage, &[]models.Article{})
    responses.Paginated(c, result.Data, int(result.Total), page, perPage)
}
```

Response:

```json
{
    "success": true,
    "data": [...],
    "meta": {
        "page": 1,
        "per_page": 15,
        "total": 42,
        "last_page": 3
    }
}
```

See [Pagination](../data/pagination.md) and
[Responses](../http/responses.md).

## 8. Add CORS and Rate Limiting

Both are applied via middleware groups or aliases:

```go
api := r.Group("/api/v1")
api.Use(
    middleware.Resolve("cors"),
    middleware.Resolve("rate"),
    middleware.Resolve("requestid"),
)
```

See [CORS](../security/cors.md) and
[Rate Limiting](../security/rate-limiting.md).

## 9. References

- [Routing — API Resources](../http/routing.md)
- [Controllers](../http/controllers.md)
- [Responses](../http/responses.md)
- [Authentication](../security/authentication.md)
- [Requests & Validation](../http/requests-validation.md)
- [Pagination](../data/pagination.md)
- [CORS](../security/cors.md)
- [Rate Limiting](../security/rate-limiting.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
