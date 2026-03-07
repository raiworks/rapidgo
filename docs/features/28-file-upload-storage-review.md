# Feature #28 — File Upload & Storage: Review

## Summary

Implemented driver-based file storage in `core/storage/`. Created `Driver` interface with `Put`, `Get`, `Delete`, `URL` methods. Two backends: `LocalDriver` (filesystem with path traversal protection) and `S3Driver` (AWS S3 / S3-compatible services). Factory function `NewDriver()` reads `STORAGE_DRIVER` env var.

## Files changed

| File | Change |
|------|--------|
| `core/storage/storage.go` | `Driver` interface + `NewDriver()` factory (supports "local" and "s3") |
| `core/storage/local.go` | `LocalDriver` with `safePath()` guard |
| `core/storage/s3.go` | New — `S3Driver` with AWS SDK v2, path validation, custom endpoint support |
| `core/storage/storage_test.go` | 20 tests (TC-01 to TC-20) |

## S3 Driver Details

- Uses `aws-sdk-go-v2/service/s3` (official AWS SDK v2)
- Reads env vars: `S3_BUCKET`, `S3_REGION`, `S3_KEY`, `S3_SECRET` (all required)
- Optional: `S3_ENDPOINT` for S3-compatible services (MinIO, DigitalOcean Spaces, etc.)
- Path traversal protection via `safePath()` (rejects `../`, empty paths)
- `URL()` returns standard AWS URL or custom endpoint URL
- `UsePathStyle = true` when custom endpoint is set (required for MinIO)

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
| TC-13 | S3 URL returns standard AWS URL | ✅ PASS |
| TC-14 | S3 URL with custom endpoint | ✅ PASS |
| TC-15 | S3 safePath rejects path traversal | ✅ PASS |
| TC-16 | S3 safePath rejects empty path | ✅ PASS |
| TC-17 | S3 safePath cleans valid paths | ✅ PASS |
| TC-18 | NewS3Driver fails without required env vars | ✅ PASS |
| TC-19 | NewS3Driver succeeds with env vars set | ✅ PASS |
| TC-20 | NewDriver returns S3Driver when STORAGE_DRIVER=s3 | ✅ PASS |

## Regression

- All 30 packages pass.
- `go vet` clean.

## Deviation log

| # | Blueprint | Ours | Reason |
|---|-----------|------|--------|
| 1 | No path traversal guard | `safePath()` on all methods | Security: prevents directory escape |
| 2 | S3Driver deferred | Now implemented | `aws-sdk-go-v2` added as dependency |
| 3 | Upload controller shown | Out of scope | App-level code, not core framework |
