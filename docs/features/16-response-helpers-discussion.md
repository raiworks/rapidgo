# 💬 Discussion: Response Helpers

> **Feature**: `16` — Response Helpers
> **Status**: 🟢 COMPLETE
> **Date**: 2026-03-06

---

## What Are We Building?

Standardized API response helpers in `http/responses/` that wrap Gin's `c.JSON()` with a consistent JSON envelope. All API responses share a common structure (`success`, `data`, `error`, `meta`) so clients get predictable, uniform JSON.

## Blueprint References

The blueprint specifies (lines 2406–2460):

1. **Directory**: `http/responses/` (line 84)
2. **Types**: `APIResponse` struct with `Success`, `Data`, `Error`, `Meta` fields
3. **Type**: `Meta` struct with `Page`, `PerPage`, `Total`, `TotalPages`
4. **Helpers**: `Success(c, data)`, `Created(c, data)`, `Error(c, status, message)`, `Paginated(c, data, page, perPage, total)`
5. **Pagination logic**: `totalPages` calculated with ceiling division

## Scope for Feature #16

### In Scope
- `APIResponse` struct — the standard envelope
- `Meta` struct — pagination metadata
- `Success(c, data)` — 200 with success envelope
- `Created(c, data)` — 201 with success envelope
- `Error(c, status, message)` — error envelope with given status
- `Paginated(c, data, page, perPage, total)` — 200 with data + meta
- Tests for all helpers

### Out of Scope (deferred)
- `Paginate` database helper (`helpers.Paginate`) — Feature #19 (Helpers)
- HTML rendering helpers — Feature #17 (Views)
- Validation error formatting — future feature
- NoContent (204) helper — not in blueprint, can be added later

## Key Design Decisions

### 1. Single File in `http/responses/`
All types and helpers fit in one file (`response.go`). The blueprint shows them together, and there's no need for separate files at this scope.

### 2. `interface{}` for Data
The blueprint uses `interface{}` for the `Data` field, allowing any JSON-serializable value. We follow this exactly — no generics needed.

### 3. Ceiling Division for TotalPages
`Paginated` computes `totalPages` with `int(total)/perPage` plus a remainder check, matching the blueprint exactly. This avoids importing `math.Ceil`.

### 4. `Error` Takes Status Code
Unlike `Success`/`Created` which have fixed status codes, `Error` accepts a status parameter so callers can use 400, 404, 422, 500, etc.

## Dependencies

| Dependency | Status | Notes |
|---|---|---|
| Feature #07 — Router | ✅ Done | Provides Gin context used by helpers |

## Discussion Complete ✅
