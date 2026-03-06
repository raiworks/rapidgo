# Feature #34 — Events / Hooks System: Changelog

## [Unreleased]

### Added
- `core/events/events.go` — `Handler` type, `Dispatcher` struct, `NewDispatcher()`, `Listen()`, `Dispatch()`, `DispatchSync()`, `Has()`.
- 8 test cases (TC-01 to TC-08) for listener invocation, async/sync dispatch, unknown event, Has, payload, concurrency.

### Removed
- `core/events/.gitkeep` — replaced by real implementation.

### Deviation log
| # | Blueprint | Ours | Reason |
|---|-----------|------|--------|
| 1 | No `Has()` method | Added `Has(event) bool` | Useful utility for checking if listeners exist before dispatching |
