# 🏗️ Architecture: Response Helpers

> **Feature**: `16` — Response Helpers
> **Discussion**: [`16-response-helpers-discussion.md`](16-response-helpers-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

Feature #16 adds standardized API response helpers to `http/responses/`. The `APIResponse` envelope provides a consistent JSON structure with `success`, `data`, `error`, and `meta` fields. Four helper functions cover the common response patterns: success, created, error, and paginated.

## File Structure

```
http/responses/
├── response.go          # APIResponse, Meta, Success, Created, Error, Paginated
└── response_test.go     # Tests for all helpers
```

### Files Created (1)
| File | Package | Lines (est.) |
|---|---|---|
| `http/responses/response.go` | `responses` | ~55 |

### Files Modified (0)
No existing files need modification.

---

## Component Design

### Response Types & Helpers (`http/responses/response.go`)

**Responsibility**: Standard API response envelope and helper functions.
**Package**: `responses`

```go
package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is the standard JSON envelope for all API responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta holds pagination metadata.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Success responds with 200 and the data wrapped in a success envelope.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{Success: true, Data: data})
}

// Created responds with 201 and the data wrapped in a success envelope.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{Success: true, Data: data})
}

// Error responds with the given status code and an error message.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, APIResponse{Success: false, Error: message})
}

// Paginated responds with 200, data, and pagination metadata.
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

**Design notes**:
- Exact match to blueprint code (lines 2408–2460)
- `omitempty` on `Data`, `Error`, `Meta` — only present when set
- `Meta` is a pointer — `nil` when not paginated, omitted from JSON
- `totalPages` ceiling division without `math.Ceil`
- No abort — callers handle flow control after calling these helpers

---

## Data Flow

### Success Response
```
Controller → responses.Success(c, user)
          → c.JSON(200, APIResponse{Success: true, Data: user})
          → {"success": true, "data": {...}}
```

### Error Response
```
Controller → responses.Error(c, 404, "not found")
          → c.JSON(404, APIResponse{Success: false, Error: "not found"})
          → {"success": false, "error": "not found"}
```

### Paginated Response
```
Controller → responses.Paginated(c, users, 1, 15, 42)
          → totalPages = ceil(42/15) = 3
          → c.JSON(200, APIResponse{..., Meta: &Meta{1, 15, 42, 3}})
          → {"success": true, "data": [...], "meta": {"page": 1, ...}}
```

---

## Constraints & Invariants

1. `Success` always returns **200**, `Created` always returns **201**
2. `Error` returns the **caller-specified** status code
3. `Paginated` always returns **200** with a non-nil `Meta`
4. `Data` is `interface{}` — any JSON-serializable value
5. JSON tags use `omitempty` — absent fields are not rendered
