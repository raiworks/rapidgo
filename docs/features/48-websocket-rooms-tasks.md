# ✅ Tasks: WebSocket Rooms / Channels

> **Feature**: `48` — WebSocket Rooms / Channels
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## Implementation Tasks

### T1: Create `core/websocket/hub.go`

- [ ] Define `Client` struct (`ID string`, `Conn *websocket.Conn`, `hub *Hub`)
- [ ] Define `Hub` struct (`clients map[string]*Client`, `rooms map[string]map[string]*Client`, `mu sync.RWMutex`)
- [ ] Implement `NewHub() *Hub`
- [ ] Implement `(*Hub).Handler(onConnect func(*Client)) gin.HandlerFunc`
- [ ] Implement `(*Hub).Join(client *Client, room string)`
- [ ] Implement `(*Hub).Leave(client *Client, room string)`
- [ ] Implement `(*Hub).Broadcast(room string, msgType websocket.MessageType, data []byte)`
- [ ] Implement `(*Hub).BroadcastOthers(room string, sender *Client, msgType websocket.MessageType, data []byte)`
- [ ] Implement `(*Hub).Send(clientID string, msgType websocket.MessageType, data []byte) error`
- [ ] Implement `(*Hub).Remove(client *Client)`
- [ ] Implement `(*Hub).Clients(room string) []*Client`
- [ ] Implement `(*Hub).Rooms() []string`

### T2: Add UUID dependency

- [ ] Run `go get github.com/google/uuid` (or verify it's already available)
- [ ] Use `uuid.NewString()` in `Handler` to generate client IDs

### T3: Create `core/websocket/hub_test.go`

- [ ] Write all test cases from the test plan (T01–T18)
- [ ] Reuse existing `newWSServer` and `wsURL` test helpers where applicable
- [ ] All tests pass with `go test ./core/websocket/ -count=1`

### T4: Regression

- [ ] `go test ./... -count=1` — all packages pass
- [ ] `go build -o bin/rapidgo.exe ./cmd` — binary builds clean

---

## Acceptance Criteria

1. `NewHub()` creates a usable Hub with no panics
2. `Hub.Handler` upgrades connections and calls `onConnect`
3. Clients can join and leave rooms
4. `Broadcast` delivers messages to all room members
5. `BroadcastOthers` skips the sender
6. `Send` delivers to a specific client by ID
7. `Remove` cleans up the client from all rooms
8. `Clients` returns a snapshot of room members
9. `Rooms` returns all active room names
10. Empty rooms are automatically cleaned up after last client leaves
11. All operations are thread-safe (no data races under `go test -race`)
12. No existing tests broken
