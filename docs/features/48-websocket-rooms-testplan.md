# 🧪 Test Plan: WebSocket Rooms / Channels

> **Feature**: `48` — WebSocket Rooms / Channels
> **Status**: 🟡 IN PROGRESS
> **File**: `core/websocket/hub_test.go`
> **Date**: 2026-03-07

---

## Test Cases

| TC | Function | Description | Approach |
|----|----------|-------------|----------|
| T01 | `TestNewHub` | NewHub returns non-nil Hub | Call `NewHub()`, assert not nil |
| T02 | `TestHub_Handler_Connects` | Handler upgrades and calls onConnect | Create Hub, register Handler, dial WS, assert onConnect called |
| T03 | `TestHub_Handler_AssignsClientID` | Client receives a non-empty ID | In onConnect, assert `client.ID != ""` |
| T04 | `TestHub_Handler_RemovesOnDisconnect` | Client removed from Hub when connection closes | Dial, close client, assert Hub has 0 clients |
| T05 | `TestHub_Join` | Client joins a room | Join room, assert `Clients(room)` contains the client |
| T06 | `TestHub_Join_CreatesRoom` | Joining a room creates it | Join room, assert `Rooms()` contains the room name |
| T07 | `TestHub_Join_Idempotent` | Joining same room twice is a no-op | Join twice, assert `Clients(room)` has exactly 1 entry |
| T08 | `TestHub_Leave` | Client leaves a room | Join then leave, assert `Clients(room)` is empty |
| T09 | `TestHub_Leave_RemovesEmptyRoom` | Empty room removed after last client leaves | Join, leave, assert `Rooms()` does not contain the room |
| T10 | `TestHub_Broadcast` | Broadcast sends to all clients in room | 2 clients join room, broadcast, both receive message |
| T11 | `TestHub_BroadcastOthers` | BroadcastOthers skips sender | 2 clients in room, client A broadcasts, only client B receives |
| T12 | `TestHub_Broadcast_IgnoresOtherRooms` | Broadcast to room A doesn't reach room B | Client in room A, client in room B, broadcast to A, only A receives |
| T13 | `TestHub_Send` | Send delivers to specific client | 2 clients, send to client A by ID, only A receives |
| T14 | `TestHub_Send_NotFound` | Send returns error for unknown ID | Call `Send("nonexistent", ...)`, assert error returned |
| T15 | `TestHub_Remove` | Remove cleans client from all rooms | Client joins 2 rooms, Remove, assert both rooms empty |
| T16 | `TestHub_Clients_Snapshot` | Clients returns a copy, not the internal map | Get clients, modify returned slice, assert hub data unchanged |
| T17 | `TestHub_Rooms_Snapshot` | Rooms returns all active room names | Create 3 rooms, assert `Rooms()` returns all 3 |
| T18 | `TestHub_Concurrent` | Concurrent join/leave/broadcast is safe | Launch goroutines doing join/leave/broadcast, run with `-race`, no panics |

---

## Test Strategy

- Tests T02–T13, T15 use real HTTP test servers (`httptest.NewServer`) with WebSocket connections, same as existing `websocket_test.go`
- T01, T14, T16, T17 are unit tests that don't need a server
- T18 is a concurrency stress test using `sync.WaitGroup` and `-race` detector
- Reuse existing `newWSServer` pattern but with `hub.Handler` instead of `Upgrader`
- All tests are in `package websocket` (same package — access to unexported fields for assertions)

---

## Helper Functions

The existing test file provides:
- `newWSServer(handler Handler, opts *Options) *httptest.Server`
- `wsURL(s *httptest.Server) string`

A new helper will be added:
- `newHubServer(hub *Hub, onConnect func(*Client)) *httptest.Server` — creates a test server using `hub.Handler(onConnect)`

---

## Pass Criteria

- All 18 tests pass: `go test ./core/websocket/ -count=1 -v`
- Race detector clean: `go test ./core/websocket/ -count=1 -race`
- Full regression: `go test ./... -count=1`
- Binary builds: `go build -o bin/rapidgo.exe ./cmd`
