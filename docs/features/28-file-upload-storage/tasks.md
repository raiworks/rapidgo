# Feature #28 — File Upload & Storage: Tasks

## Prerequisites

- [x] Core infrastructure shipped (#02, #05)
- [x] `storage/uploads/` directory exists with `.gitkeep`

## Implementation tasks

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 1 | Create `Driver` interface + `NewDriver()` | `core/storage/storage.go` | ⬜ |
| 2 | Create `LocalDriver` with path traversal guard | `core/storage/local.go` | ⬜ |
| 3 | Write tests | `core/storage/storage_test.go` | ⬜ |
| 4 | Full regression + `go vet` | — | ⬜ |
| 5 | Commit, merge, review doc, roadmap update | — | ⬜ |

## Acceptance criteria

- `Driver` interface defines `Put`, `Get`, `Delete`, `URL`.
- `LocalDriver` implements all 4 methods.
- `NewDriver()` returns a `LocalDriver` when `STORAGE_DRIVER=local` or unset.
- `NewDriver()` returns error for unsupported driver.
- Path traversal attempts are rejected with an error.
- Files can be written, read back, and deleted.
- `URL()` returns `BaseURL + "/" + path`.
- All existing tests pass (regression).
