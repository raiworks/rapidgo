# Feature #34 — Events / Hooks System: Architecture

## Component overview

```
core/events/events.go
    │
    ├── Handler type: func(payload interface{})
    │
    ├── Dispatcher struct
    │       sync.RWMutex + map[string][]Handler
    │
    ├── NewDispatcher() *Dispatcher
    │
    ├── Listen(event, handler)
    │       Appends handler to named event's listener list
    │
    ├── Dispatch(event, payload)
    │       Calls all handlers in goroutines (async)
    │
    ├── DispatchSync(event, payload)
    │       Calls all handlers sequentially (sync)
    │
    └── Has(event) bool
            Returns true if any listener is registered for the event
```

## New file

| File | Purpose |
|------|---------|
| `core/events/events.go` | `Handler`, `Dispatcher`, `NewDispatcher()`, `Listen()`, `Dispatch()`, `DispatchSync()`, `Has()` |

## Removed

| File | Reason |
|------|--------|
| `core/events/.gitkeep` | Replaced by real implementation |

## Types

```go
// Handler is a function that handles an event.
type Handler func(payload interface{})

// Dispatcher manages event listeners and dispatches events.
type Dispatcher struct {
    mu        sync.RWMutex
    listeners map[string][]Handler
}
```

## Functions

| Function | Signature | Behaviour |
|----------|-----------|-----------|
| `NewDispatcher()` | `func NewDispatcher() *Dispatcher` | Returns empty dispatcher |
| `Listen()` | `func (d *Dispatcher) Listen(event string, handler Handler)` | Appends handler under write lock |
| `Dispatch()` | `func (d *Dispatcher) Dispatch(event string, payload interface{})` | Fires each handler in a new goroutine |
| `DispatchSync()` | `func (d *Dispatcher) DispatchSync(event string, payload interface{})` | Fires each handler sequentially |
| `Has()` | `func (d *Dispatcher) Has(event string) bool` | Returns true if any listener exists for event |

## Usage

```go
d := events.NewDispatcher()

d.Listen("user.created", func(p interface{}) {
    fmt.Println("User created:", p)
})

d.Dispatch("user.created", user)      // async
d.DispatchSync("user.created", user)   // sync
d.Has("user.created")                  // true
```
