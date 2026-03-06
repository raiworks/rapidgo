# Feature #31 — WebSocket Support: Test Plan

## Test cases

| TC | Description | Method | Expected |
|----|-------------|--------|----------|
| TC-01 | Upgrader upgrades HTTP to WebSocket | Dial /ws | Successful connection |
| TC-02 | Echo handler echoes text message | Write "hello" | Read back "hello" |
| TC-03 | Echo handler echoes binary message | Write binary | Read back same bytes |
| TC-04 | Connection closes cleanly on handler return | Close client | No error/panic |
| TC-05 | Upgrader with nil options uses defaults | Dial with nil opts | Successful connection |
| TC-06 | Custom handler receives connection | Dial /ws | Handler callback invoked |

## Notes

- Tests use `httptest.NewServer` to create a real HTTP server for WebSocket dial.
- `github.com/coder/websocket` provides `websocket.Dial()` for test clients.
- Echo tests send a message and verify the response matches.
