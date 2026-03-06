# Feature #31 — WebSocket Support: Architecture

## Component overview

```
core/websocket/websocket.go
    │
    ├── Upgrader(handler Handler, opts *Options) gin.HandlerFunc
    │       Upgrades HTTP → WebSocket, calls user-provided handler
    │
    ├── Handler type: func(*websocket.Conn, context.Context)
    │       User callback for handling the connection
    │
    ├── Options struct
    │       AcceptOptions wrapper (origin patterns, etc.)
    │
    └── Echo(conn, ctx)
            Built-in echo handler for testing
```

## New file

| File | Purpose |
|------|---------|
| `core/websocket/websocket.go` | `Upgrader()`, `Handler` type, `Options`, `Echo()` |

## Removed

| File | Reason |
|------|--------|
| `core/websocket/.gitkeep` | Replaced by real implementation |

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/coder/websocket` | WebSocket protocol implementation |

## Types

```go
// Handler is a callback for handling a WebSocket connection.
type Handler func(conn *websocket.Conn, ctx context.Context)

// Options configures the WebSocket upgrade.
type Options struct {
    OriginPatterns []string // Allowed origin patterns (default: allow all)
}
```

## Functions

| Function | Signature | Behaviour |
|----------|-----------|-----------|
| `Upgrader()` | `func Upgrader(handler Handler, opts *Options) gin.HandlerFunc` | Upgrades request to WebSocket, calls handler, closes on return |
| `Echo()` | `func Echo(conn *websocket.Conn, ctx context.Context)` | Reads messages and writes them back (echo loop) |

## Usage

```go
// In routes:
r.GET("/ws", ws.Upgrader(ws.Echo, nil))

// Custom handler:
r.GET("/ws/chat", ws.Upgrader(func(conn *websocket.Conn, ctx context.Context) {
    // custom logic
}, nil))
```
