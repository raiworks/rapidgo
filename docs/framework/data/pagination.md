---
title: "Pagination"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Pagination

## Abstract

This document covers the generic pagination helper for GORM queries
and its integration with the API response envelope.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Paginate Helper](#2-paginate-helper)
3. [Usage in Controllers](#3-usage-in-controllers)
4. [Response Format](#4-response-format)
5. [Query Parameters](#5-query-parameters)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. Paginate Helper

Located in `app/helpers/`. Works with any GORM query:

```go
package helpers

import "gorm.io/gorm"

type PaginateResult struct {
    Page       int   `json:"page"`
    PerPage    int   `json:"per_page"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
}

func Paginate(db *gorm.DB, page, perPage int, dest interface{}) (*PaginateResult, error) {
    if page < 1 {
        page = 1
    }
    if perPage < 1 || perPage > 100 {
        perPage = 15
    }

    var total int64
    db.Count(&total)

    offset := (page - 1) * perPage
    err := db.Offset(offset).Limit(perPage).Find(dest).Error

    totalPages := int(total) / perPage
    if int(total)%perPage != 0 {
        totalPages++
    }

    return &PaginateResult{
        Page:       page,
        PerPage:    perPage,
        Total:      total,
        TotalPages: totalPages,
    }, err
}
```

### Defaults and Bounds

| Parameter | Default | Bounds |
|-----------|---------|--------|
| `page` | 1 | minimum 1 |
| `perPage` | 15 | 1–100 |

Values outside bounds are clamped to defaults.

## 3. Usage in Controllers

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

### With Filters

Chain GORM scopes before paginating:

```go
query := db.Model(&models.Post{}).Where("published = ?", true).Order("created_at DESC")
result, err := helpers.Paginate(query, page, perPage, &posts)
```

## 4. Response Format

The `Paginated` response helper produces:

```json
{
    "success": true,
    "data": [
        {"id": 1, "name": "Alice"},
        {"id": 2, "name": "Bob"}
    ],
    "meta": {
        "page": 1,
        "per_page": 15,
        "total": 42,
        "total_pages": 3
    }
}
```

## 5. Query Parameters

Clients request pages via query parameters:

```text
GET /api/users?page=2&per_page=20
```

For struct-based validation of pagination parameters:

```go
type PaginationRequest struct {
    Page    int `form:"page" binding:"omitempty,min=1"`
    PerPage int `form:"per_page" binding:"omitempty,min=1,max=100"`
}
```

## 6. Security Considerations

- `perPage` **MUST** have an upper bound (default: 100) to prevent
  denial-of-service via excessively large page sizes.
- Page and perPage values **MUST** be validated as positive integers.

## 7. References

- [Responses — Paginated](../http/responses.md#4-paginated-responses)
- [Database](database.md)
- [Controllers](../http/controllers.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
