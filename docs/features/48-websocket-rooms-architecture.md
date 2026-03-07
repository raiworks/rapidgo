# 🏗️ Architecture: WebSocket Rooms / Channels

> **Feature**: `48` — WebSocket Rooms / Channels
> **Status**: 🟡 IN PROGRESS
> **Package**: `core/websocket`
> **Date**: 2026-03-07

---

## Component Overview

```
core/websocket/
    │
    ├── websocket.go          (EXISTING — unchanged)
    │   ├── Handler type
    │   ├── Options struct
    │   ├── Upgrader()
    │   └── Echo()
    │
    └── hub.go                (NEW)
        ├── Client struct
        │   ├── ID       string
        │   ├── Conn     *websocket.Conn
        │   └── hub      *Hub
        │
        ├── Hub struct
        │   ├── clients  map[string]*Client
        │   ├── rooms    map[string]map[string]*Client
        │   └── mu       sync.RWMutex
        │
        ├── NewHub() *Hub
        ├── (*Hub).Handler(onConnect func(*Client)) gin.HandlerFunc
        ├── (*Hub).Join(client *Client, room string)
        ├── (*Hub).Leave(client *Client, room string)
        ├── (*Hub).Broadcast(room string, msgType websocket.MessageType, data []byte)
        ├── (*Hub).BroadcastOthers(room string, sender *Client, msgType websocket.MessageType, data []byte)
        ├── (*Hub).Send(clientID string, msgType websocket.MessageType, data []byte) error
        ├── (*Hub).Remove(client *Client)
        ├── (*Hub).Clients(room string) []*Client
        └── (*Hub).Rooms() []string
```

---

## Structs

### Client

```go
// Client represents a single WebSocket connection registered with a Hub.
type Client struct {
    ID   string
    Conn *websocket.Conn
    hub  *Hub
}
```

| Field | Type | Visibility | Description |
|-------|------|------------|-------------|
| `ID` | `string` | Exported | Unique identifier (UUID v4), generated on connect |
| `Conn` | `*websocket.Conn` | Exported | The underlying WebSocket connection |
| `hub` | `*Hub` | Unexported | Back-reference to the owning Hub |

### Hub

```go
// Hub manages WebSocket clients and rooms.
type Hub struct {
    clients map[string]*Client
    rooms   map[string]map[string]*Client
    mu      sync.RWMutex
}
```

| Field | Type | Visibility | Description |
|-------|------|------------|-------------|
| `clients` | `map[string]*Client` | Unexported | All connected clients, keyed by `Client.ID` |
| `rooms` | `map[string]map[string]*Client` | Unexported | Room name → set of clients (keyed by `Client.ID`) |
| `mu` | `sync.RWMutex` | Unexported | Protects `clients` and `rooms` |

---

## Functions

### NewHub

```go
func NewHub() *Hub
```

Creates an empty Hub with initialized maps. No goroutines spawned.

### (*Hub).Handler

```go
func (h *Hub) Handler(onConnect func(*Client)) gin.HandlerFunc
```

Returns a `gin.HandlerFunc` that:
1. Accepts the WebSocket upgrade (same as `Upgrader`, uses `InsecureSkipVerify: true` default)
2. Creates a `Client` with a generated UUID
3. Registers the client in `h.clients`
4. Calls `onConnect(client)` — the developer's callback for handling the connection
5. On return: calls `h.Remove(client)` to clean up, then closes the connection

The `onConnect` callback is where the developer reads messages, joins rooms, and implements business logic. It runs synchronously — one goroutine per connection (managed by Gin/HTTP server).

### (*Hub).Join

```go
func (h *Hub) Join(client *Client, room string)
```

Adds the client to the named room. Creates the room if it doesn't exist. No-op if the client is already in the room.

### (*Hub).Leave

```go
func (h *Hub) Leave(client *Client, room string)
```

Removes the client from the named room. Deletes the room map entry if the room becomes empty.

### (*Hub).Broadcast

```go
func (h *Hub) Broadcast(room string, msgType websocket.MessageType, data []byte)
```

Sends the message to **all** clients in the room. Skips clients whose `conn.Write` returns an error (connection broken). Uses `RLock` for concurrent read access.

### (*Hub).BroadcastOthers

```go
func (h *Hub) BroadcastOthers(room string, sender *Client, msgType websocket.MessageType, data []byte)
```

Same as `Broadcast` but skips the `sender` client. Useful for chat-style "echo to others" patterns.

### (*Hub).Send

```go
func (h *Hub) Send(clientID string, msgType websocket.MessageType, data []byte) error
```

Sends a message to a single client by ID. Returns an error if the client is not found.

### (*Hub).Remove

```go
func (h *Hub) Remove(client *Client)
```

Removes the client from all rooms and from `h.clients`. Called automatically by `Handler` on disconnect, but can also be called explicitly.

### (*Hub).Clients

```go
func (h *Hub) Clients(room string) []*Client
```

Returns a snapshot (slice copy) of all clients in the named room. Returns nil if the room doesn't exist.

### (*Hub).Rooms

```go
func (h *Hub) Rooms() []string
```

Returns a snapshot (slice) of all room names. Order is not guaranteed.

---

## Dependencies

| Package | Purpose | Status |
|---------|---------|--------|
| `github.com/coder/websocket` | WebSocket protocol (already in go.mod) | Existing |
| `github.com/gin-gonic/gin` | HTTP handler integration | Existing |
| `github.com/google/uuid` | Client ID generation | **NEW** — needs `go get` |

> **Note**: `crypto/rand` could generate UUIDs manually, but `google/uuid` is the standard library for UUID v4 in Go. It's already transitively available as a dependency of Gin. We can use it directly.

---

## UUID Dependency Check

```
google/uuid is already an indirect dependency via Gin → can be promoted to direct.
```

If `google/uuid` is not already in `go.mod`, it will be added automatically by `go get`. No license concern — BSD-3-Clause.

---

## Files Changed

| File | Action | Description |
|------|--------|-------------|
| `core/websocket/hub.go` | NEW | Hub, Client, all methods |
| `core/websocket/hub_test.go` | NEW | All test cases |

No existing files modified. `websocket.go` remains unchanged.

---

## Usage Example

```go
// In routes/ws.go:
import (
    ws "github.com/RAiWorks/RapidGo/core/websocket"
    "github.com/coder/websocket"
)

hub := ws.NewHub()

r.Get("/ws/chat", hub.Handler(func(client *ws.Client) {
    hub.Join(client, "general")
    defer hub.Leave(client, "general")

    ctx := client.Conn.CloseRead(context.Background())
    for {
        msgType, data, err := client.Conn.Read(ctx)
        if err != nil {
            return
        }
        hub.BroadcastOthers("general", client, msgType, data)
    }
}))
```

---

## Thread Safety Model

All Hub methods acquire the mutex appropriately:

| Method | Lock Type | Reason |
|--------|-----------|--------|
| `Handler` (register/remove) | `Lock` (write) | Modifies `clients` map |
| `Join` | `Lock` (write) | Modifies `rooms` map |
| `Leave` | `Lock` (write) | Modifies `rooms` map |
| `Broadcast` | `RLock` (read) | Reads `rooms` map, writes to individual conns |
| `BroadcastOthers` | `RLock` (read) | Same as Broadcast |
| `Send` | `RLock` (read) | Reads `clients` map, writes to one conn |
| `Remove` | `Lock` (write) | Modifies both maps |
| `Clients` | `RLock` (read) | Reads `rooms` map |
| `Rooms` | `RLock` (read) | Reads `rooms` map |
