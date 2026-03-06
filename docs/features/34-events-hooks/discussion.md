# Feature #34 — Events / Hooks System: Discussion

## What problem does this solve?

Components often need to react to things that happen elsewhere (user created, order placed, password changed) without direct coupling. A publish-subscribe dispatcher lets any part of the app fire named events and any number of listeners respond independently.

## Why now?

Service container (#05) and providers (#06) are shipped. The event dispatcher can be registered as a singleton and listeners wired during provider boot.

## What does the blueprint specify?

- `Handler` type: `func(payload interface{})`.
- `Dispatcher` struct with `sync.RWMutex` + `map[string][]Handler`.
- `NewDispatcher()` constructor.
- `Listen(event, handler)` — registers a handler for a named event.
- `Dispatch(event, payload)` — calls all handlers asynchronously (`go h(payload)`).
- `DispatchSync(event, payload)` — calls all handlers synchronously.
- `EventProvider` registers a `"events"` singleton in the container.

## Design decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Package | `core/events/` | Blueprint-specified; existing stub directory |
| Handler type | `func(payload interface{})` | Blueprint-specified; simple and flexible |
| Async dispatch | `go h(payload)` in `Dispatch()` | Blueprint-specified; non-blocking |
| Sync dispatch | Sequential in `DispatchSync()` | Blueprint-specified; useful for ordering guarantees |
| Thread safety | `sync.RWMutex` on listeners map | Blueprint-specified; safe for concurrent Listen/Dispatch |
| Has method | Add `Has(event) bool` | Small addition — useful for checking if any listener registered |

## What is out of scope?

- Event priority / ordering (app-level concern).
- Event objects / typed events (can be added later over `interface{}`).
- Wildcard listeners (e.g., `"user.*"`).
- EventProvider registration (can be added when needed).
