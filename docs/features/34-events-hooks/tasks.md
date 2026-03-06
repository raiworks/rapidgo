# Feature #34 — Events / Hooks System: Tasks

## Prerequisites

- [x] Service container shipped (#05)
- [x] `core/events/` directory exists (stub with `.gitkeep`)

## Implementation tasks

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 1 | Define `Handler` type and `Dispatcher` struct | `core/events/events.go` | ⬜ |
| 2 | Implement `NewDispatcher()`, `Listen()`, `Has()` | `core/events/events.go` | ⬜ |
| 3 | Implement `Dispatch()` (async) and `DispatchSync()` | `core/events/events.go` | ⬜ |
| 4 | Remove `.gitkeep` | `core/events/.gitkeep` | ⬜ |
| 5 | Write tests | `core/events/events_test.go` | ⬜ |
| 6 | Full regression + `go vet` | — | ⬜ |
| 7 | Commit, merge, review doc, roadmap update | — | ⬜ |

## Acceptance criteria

- `Handler` type is `func(payload interface{})`.
- `Dispatcher` is thread-safe (`sync.RWMutex`).
- `Listen()` registers handler for named event.
- `Dispatch()` fires handlers asynchronously (`go h(payload)`).
- `DispatchSync()` fires handlers sequentially.
- `Has()` returns whether any listener is registered.
- Multiple listeners per event are supported.
- All existing tests pass (regression).
