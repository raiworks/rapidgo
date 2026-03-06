# Feature #31 — WebSocket Support: Tasks

## Prerequisites

- [x] Router infrastructure shipped (#07)
- [x] `core/websocket/` directory exists (stub with `.gitkeep`)

## Implementation tasks

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 1 | `go get github.com/coder/websocket` | `go.mod`, `go.sum` | ⬜ |
| 2 | Create `Handler`, `Options`, `Upgrader()`, `Echo()` | `core/websocket/websocket.go` | ⬜ |
| 3 | Remove `.gitkeep` | `core/websocket/.gitkeep` | ⬜ |
| 4 | Write tests | `core/websocket/websocket_test.go` | ⬜ |
| 5 | Full regression + `go vet` | — | ⬜ |
| 6 | Commit, merge, review doc, roadmap update | — | ⬜ |

## Acceptance criteria

- `Handler` type defined as `func(*websocket.Conn, context.Context)`.
- `Upgrader()` returns a `gin.HandlerFunc` that upgrades and calls the handler.
- `Echo()` reads and writes messages in a loop.
- `Options` allows configuring origin patterns.
- Tests verify upgrade, echo, and connection close.
- All existing tests pass (regression).
