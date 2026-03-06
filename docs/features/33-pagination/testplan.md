# Feature #33 — Pagination: Test Plan

## Test cases

| TC | Description | Method | Expected |
|----|-------------|--------|----------|
| TC-01 | Page < 1 clamped to 1 | Paginate(db, 0, 10, &dest) | result.Page == 1 |
| TC-02 | PerPage < 1 clamped to 15 | Paginate(db, 1, 0, &dest) | result.PerPage == 15 |
| TC-03 | PerPage > 100 clamped to 15 | Paginate(db, 1, 200, &dest) | result.PerPage == 15 |
| TC-04 | Valid inputs pass through | Paginate(db, 2, 10, &dest) | result.Page == 2, result.PerPage == 10 |
| TC-05 | TotalPages ceiling division | 25 total, perPage 10 | TotalPages == 3 |
| TC-06 | TotalPages exact division | 30 total, perPage 10 | TotalPages == 3 |
| TC-07 | TotalPages zero when empty | 0 total | TotalPages == 0 |
| TC-08 | Offset calculation | page=2, perPage=10 | Offset == 10 (second page of data) |
| TC-09 | Returns GORM error | Invalid DB operation | error != nil |

## Notes

- Tests use SQLite in-memory database (`:memory:`) with GORM.
- A simple test model (e.g., `Item{ID, Name}`) is used for seeding.
- TC-08 verifies by seeding 25 items, requesting page 2 with perPage 10, and checking dest has 10 items.
- TC-09 can use a closed/invalid DB connection or a model that triggers a GORM error.
