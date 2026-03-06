# Feature #27 — Request ID / Tracing: Changelog

## [Unreleased]

### Status
Feature already fully implemented in Feature #08 (Middleware Infrastructure).

### Existing (no changes)
- `core/middleware/request_id.go` — `RequestID()` with UUID v4 generation.
- `"requestid"` alias registered in `middleware_provider.go`.
- Included in `"global"` middleware group.
- 3 test cases (TC-08, TC-09, TC-14) covering all blueprint requirements.

### Added
- Feature #27 documentation confirming blueprint coverage.

### Deviation log
| # | Blueprint | Ours | Reason |
|---|-----------|------|--------|
| 1 | Raw hex ID (32 chars) | UUID v4 (36 chars with dashes) | Standard format, version/variant bits, better for human readability and tooling |
