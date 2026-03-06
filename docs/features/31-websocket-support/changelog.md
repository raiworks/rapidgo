# Feature #31 — WebSocket Support: Changelog

## [Unreleased]

### Added
- `core/websocket/websocket.go` — `Handler` type, `Options` struct, `Upgrader()`, `Echo()`.
- `github.com/coder/websocket` dependency.
- 6 test cases (TC-01 to TC-06) for upgrade, echo, and connection handling.

### Removed
- `core/websocket/.gitkeep` — replaced by real implementation.

### Deviation log
| # | Blueprint | Ours | Reason |
|---|-----------|------|--------|
| 1 | Single `HandleWS` function | `Upgrader()` + `Handler` type + `Echo()` | More reusable — framework handles upgrade, app provides handler callback |
