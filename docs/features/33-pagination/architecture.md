# Feature #33 — Pagination: Architecture

## Component overview

```
app/helpers/pagination.go
    │
    ├── PaginateResult struct
    │       Page, PerPage, Total, TotalPages
    │
    └── Paginate(db, page, perPage, dest) (*PaginateResult, error)
            Counts total, applies offset/limit, returns metadata
```

## New file

| File | Purpose |
|------|---------|
| `app/helpers/pagination.go` | `PaginateResult` struct and `Paginate()` function |

## Dependencies

| Package | Purpose | Status |
|---------|---------|--------|
| `gorm.io/gorm` | Query builder | Already in go.mod |

## Types

```go
// PaginateResult holds pagination metadata returned by Paginate().
type PaginateResult struct {
    Page       int   `json:"page"`
    PerPage    int   `json:"per_page"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
}
```

## Functions

| Function | Signature | Behaviour |
|----------|-----------|-----------|
| `Paginate()` | `func Paginate(db *gorm.DB, page, perPage int, dest interface{}) (*PaginateResult, error)` | Clamps inputs, counts total, applies Offset+Limit, computes TotalPages |

## Input clamping

| Input | Rule |
|-------|------|
| `page < 1` | Set to 1 |
| `perPage < 1` or `perPage > 100` | Set to 15 |

## Integration with existing code

```go
// Controller:
var users []models.User
result, err := helpers.Paginate(db.Model(&models.User{}), page, perPage, &users)
responses.Paginated(c, users, result.Page, result.PerPage, result.Total)
```

The `responses.Paginated()` helper was shipped in Feature #16.
