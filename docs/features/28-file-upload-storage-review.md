# Feature #28 — File Upload & Storage: Review

## Summary

Implemented driver-based file storage in `core/storage/`. Created `Driver` interface with `Put`, `Get`, `Delete`, `URL` methods and `LocalDriver` implementation with path traversal protection. Factory function `NewDriver()` reads `STORAGE_DRIVER` env var.

## Files changed

| File | Change |
|------|--------|
| `core/storage/storage.go` | New — `Driver` interface + `NewDriver()` factory |
| `core/storage/local.go` | New — `LocalDriver` with `safePath()` guard |
| `core/storage/storage_test.go` | New — 12 tests (TC-01 to TC-12) |

## Test results

| TC | Description | Result |
|----|-------------|--------|
| TC-01 | Put writes file to disk | ✅ PASS |
| TC-02 | Put creates intermediate directories | ✅ PASS |
| TC-03 | Get returns file content | ✅ PASS |
| TC-04 | Get returns error for missing file | ✅ PASS |
| TC-05 | Delete removes file | ✅ PASS |
| TC-06 | Delete returns error for missing file | ✅ PASS |
| TC-07 | URL returns BaseURL + path | ✅ PASS |
| TC-08 | Path traversal in Put rejected | ✅ PASS |
| TC-09 | Path traversal in Get rejected | ✅ PASS |
| TC-10 | Path traversal in Delete rejected | ✅ PASS |
| TC-11 | NewDriver returns LocalDriver by default | ✅ PASS |
| TC-12 | NewDriver returns error for unknown driver | ✅ PASS |

## Regression

- All 24 packages pass.
- `go vet` clean.

## Deviation log

| # | Blueprint | Ours | Reason |
|---|-----------|------|--------|
| 1 | No path traversal guard | `safePath()` on all methods | Security: prevents directory escape |
| 2 | S3Driver included | Deferred | `aws-sdk-go-v2` is heavy; not needed for core |
| 3 | Upload controller shown | Out of scope | App-level code, not core framework |

## Commit

`0cf0182` — merged to `main`.
