# Feature #28 — File Upload & Storage: Changelog

## [Unreleased]

### Added
- `core/storage/storage.go` — `Driver` interface (`Put`, `Get`, `Delete`, `URL`) and `NewDriver()` factory.
- `core/storage/local.go` — `LocalDriver` with path traversal protection.
- 12 test cases covering CRUD operations, path traversal security, and factory function.

### Deviation log
| # | Blueprint | Ours | Reason |
|---|-----------|------|--------|
| 1 | No path traversal guard | Added guard to all LocalDriver methods | Security: prevents `../` directory escape |
| 2 | S3Driver included | Deferred | `aws-sdk-go-v2` is a massive dependency; not needed for core framework |
| 3 | Upload controller shown | Out of scope | App-level code, not core framework |
