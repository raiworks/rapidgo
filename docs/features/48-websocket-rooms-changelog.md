# 📝 Changelog: WebSocket Rooms / Channels

> **Feature**: `48` — WebSocket Rooms / Channels
> **Status**: ✅ SHIPPED
> **Date**: 2026-03-07

---

## Added

- `core/websocket/hub.go` — `Hub`, `Client`, room management, broadcast, send
- `core/websocket/hub_test.go` — 18 unit tests for hub, rooms, broadcast, concurrency

## Dependencies

| Package | Version | License | Purpose |
|---------|---------|---------|---------|
| `github.com/google/uuid` | (already indirect via Gin) | BSD-3-Clause | Client ID generation |

## Files

| File | Action |
|------|--------|
| `core/websocket/hub.go` | NEW |
| `core/websocket/hub_test.go` | NEW |

## Migration Guide

- No migrations required
- No new environment variables
- No breaking changes — existing `Upgrader`, `Handler`, `Echo` remain unchanged
- Use `NewHub()` + `hub.Handler(onConnect)` for room-based WebSocket routing
- Use `hub.Join/Leave/Broadcast/Send` inside the `onConnect` callback
