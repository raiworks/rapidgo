# Feature #33 — Pagination: Discussion

## What problem does this solve?

Listing endpoints must return data in bounded pages rather than unbounded result sets. A generic pagination helper avoids duplicating offset/limit/count logic in every controller.

## Why now?

Database (#09), models (#11), and response helpers (#16, including `Paginated()`) are shipped. The missing piece is the GORM query-side helper that does the counting and slicing.

## What does the blueprint specify?

- `PaginateResult` struct with `Page`, `PerPage`, `Total`, `TotalPages`.
- `Paginate(db *gorm.DB, page, perPage int, dest interface{}) (*PaginateResult, error)` function.
- Clamps `page` (min 1) and `perPage` (1–100, default 15).
- Calls `db.Count(&total)`, then `db.Offset().Limit().Find(dest)`.
- Computes `TotalPages` with ceiling division.
- Lives in `app/helpers/` (same package as existing helpers).

## Design decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Package | `app/helpers` | Blueprint-specified; helpers package already exists |
| GORM dependency | Import `gorm.io/gorm` | Already in go.mod from #09 |
| Input clamping | page < 1 → 1, perPage < 1 or > 100 → 15 | Blueprint defaults; prevents abuse |
| Zero-total edge case | TotalPages = 0 when total = 0 | Avoids returning 1 page for empty sets |
| Return error | From `Find()` only | `Count()` error is visible through GORM's chain; Find error covers both |

## What is out of scope?

- Cursor-based / keyset pagination (different pattern, app-level).
- Response helper changes (`Paginated()` already ships in #16).
- `PaginationRequest` struct (validation — already covered by #23).
