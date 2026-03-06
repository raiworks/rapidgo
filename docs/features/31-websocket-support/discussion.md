# Feature #31 — WebSocket Support: Discussion

## What problem does this solve?

Real-time features (chat, notifications, live updates, dashboards) require persistent bidirectional connections. WebSockets provide full-duplex communication over a single TCP connection, replacing the need for polling.

## Why now?

HTTP infrastructure (#07), middleware (#08), and router (#07) are shipped. WebSocket support is the next real-time capability, completing the communication layer.

## What does the blueprint specify?

- Library: `github.com/coder/websocket` (successor to archived `gorilla/websocket`).
- `HandleWS(w, r)` function that accepts a WebSocket connection, reads messages in a loop, and echoes them back.
- Uses `websocket.Accept()` for upgrading, `conn.Read()` / `conn.Write()` for I/O.

## Design decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Library | `github.com/coder/websocket` | Blueprint-specified; actively maintained, stdlib-compatible API |
| Package location | `core/websocket/` | Existing stub directory |
| API design | `Upgrader` helper + `Handler` type | More reusable than a single echo function; app code can provide custom message handlers |
| Handler type | `func(conn, context)` callback | Lets app define behaviour; framework handles upgrade |
| Echo handler | Include as example/default | Matches blueprint; useful for testing |

## What is out of scope?

- Room/channel management (app-level concern).
- Broadcasting to multiple connections (app-level concern).
- Authentication within WebSocket (can use existing auth middleware on the HTTP route).
- Message serialization (JSON, protobuf) — app-level.
