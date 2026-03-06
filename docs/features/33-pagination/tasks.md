# Feature #33 — Pagination: Tasks

## Prerequisites

- [x] Database infrastructure shipped (#09)
- [x] Models shipped (#11)
- [x] Response helpers with `Paginated()` shipped (#16)
- [x] `app/helpers/` package exists

## Implementation tasks

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 1 | Create `PaginateResult` struct | `app/helpers/pagination.go` | ⬜ |
| 2 | Implement `Paginate()` with clamping, count, offset/limit | `app/helpers/pagination.go` | ⬜ |
| 3 | Write tests | `app/helpers/helpers_test.go` (or new file) | ⬜ |
| 4 | Full regression + `go vet` | — | ⬜ |
| 5 | Commit, merge, review doc, roadmap update | — | ⬜ |

## Acceptance criteria

- `PaginateResult` has Page, PerPage, Total, TotalPages fields.
- `page < 1` clamped to 1.
- `perPage < 1` or `perPage > 100` clamped to 15.
- `TotalPages` uses ceiling division.
- `TotalPages = 0` when `total = 0`.
- Offset computed as `(page - 1) * perPage`.
- Error from GORM `Find()` is returned.
- All existing tests pass (regression).
