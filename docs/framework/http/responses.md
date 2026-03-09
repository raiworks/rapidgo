---
title: "Responses"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Responses

## Abstract

This document covers the standardized API response envelope, response
helper functions, pagination metadata, and HTML template rendering
patterns.

## Table of Contents

1. [Terminology](#1-terminology)
2. [API Response Envelope](#2-api-response-envelope)
3. [Response Helpers](#3-response-helpers)
4. [Paginated Responses](#4-paginated-responses)
5. [HTML Rendering](#5-html-rendering)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. API Response Envelope

All API responses use a consistent JSON envelope defined in
`http/responses/`:

```go
package responses

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
    Page       int   `json:"page"`
    PerPage    int   `json:"per_page"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
}
```

### Success Response

```json
{
    "success": true,
    "data": {
        "id": 1,
        "name": "Alice",
        "email": "alice@example.com"
    }
}
```

### Error Response

```json
{
    "success": false,
    "error": "validation failed"
}
```

### Paginated Response

```json
{
    "success": true,
    "data": [...],
    "meta": {
        "page": 1,
        "per_page": 15,
        "total": 42,
        "total_pages": 3
    }
}
```

## 3. Response Helpers

### `Success`

Returns 200 with the data wrapped in a success envelope:

```go
func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, APIResponse{Success: true, Data: data})
}
```

Usage:

```go
responses.Success(c, user)
responses.Success(c, gin.H{"message": "done"})
```

### `Created`

Returns 201 for newly created resources:

```go
func Created(c *gin.Context, data interface{}) {
    c.JSON(http.StatusCreated, APIResponse{Success: true, Data: data})
}
```

Usage:

```go
responses.Created(c, newUser)
```

### `Error`

Returns the specified HTTP status with an error message:

```go
func Error(c *gin.Context, status int, message string) {
    c.JSON(status, APIResponse{Success: false, Error: message})
}
```

Usage:

```go
responses.Error(c, 404, "user not found")
responses.Error(c, 422, "validation failed")
responses.Error(c, 500, "internal server error")
```

### `Paginated`

Returns 200 with data and pagination metadata:

```go
func Paginated(c *gin.Context, data interface{}, page, perPage int, total int64) {
    totalPages := int(total) / perPage
    if int(total)%perPage != 0 {
        totalPages++
    }
    c.JSON(http.StatusOK, APIResponse{
        Success: true,
        Data:    data,
        Meta: &Meta{
            Page:       page,
            PerPage:    perPage,
            Total:      total,
            TotalPages: totalPages,
        },
    })
}
```

Usage:

```go
responses.Paginated(c, users, result.Page, result.PerPage, result.Total)
```

## 4. Paginated Responses

The full pagination flow combines the `Paginate` helper with the
`Paginated` response helper:

```go
func ListUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))

    var users []models.User
    result, err := helpers.Paginate(db.Model(&models.User{}), page, perPage, &users)
    if err != nil {
        responses.Error(c, 500, "failed to fetch users")
        return
    }
    responses.Paginated(c, users, result.Page, result.PerPage, result.Total)
}
```

See [Pagination](../data/pagination.md) for the `Paginate` helper
implementation.

## 5. HTML Rendering

For SSR responses, use Gin's `c.HTML()`:

```go
c.HTML(http.StatusOK, "users/show.html", gin.H{
    "title": "User Profile",
    "user":  user,
})
```

On validation failure:

```go
c.HTML(http.StatusUnprocessableEntity, "users/create.html", gin.H{
    "errors": v.Errors(),
    "old":    gin.H{"name": name, "email": email},
})
```

## 6. Security Considerations

- API error responses **MUST NOT** include stack traces or internal
  details in production.
- Model fields containing sensitive data (passwords, tokens)
  **MUST** use `json:"-"` to exclude them from responses.
- The `Error` helper always returns the exact message given — ensure
  messages are user-safe and do not leak implementation details.

## 7. References

- [Controllers](controllers.md)
- [Pagination](../data/pagination.md)
- [Error Handling](../core/error-handling.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
