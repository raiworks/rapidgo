---
title: "WebSocket Rooms & Channels"
version: "0.1.0"
status: "Final"
date: "2026-03-11"
last_updated: "2026-03-11"
authors:
  - "RAiWorks"
supersedes: ""
---

# WebSocket Rooms & Channels

## Abstract

A guide to building real-time features using the RapidGo Hub — room
management, broadcasting, heartbeat, and lifecycle callbacks.

## Table of Contents

1. [Overview](#1-overview)
2. [Creating a Hub](#2-creating-a-hub)
3. [Handling Connections](#3-handling-connections)
4. [Rooms (Join / Leave)](#4-rooms-join--leave)
5. [Broadcasting](#5-broadcasting)
6. [Heartbeat & Dead Connection Detection](#6-heartbeat--dead-connection-detection)
7. [Lifecycle Callbacks](#7-lifecycle-callbacks)
8. [Direct Messaging](#8-direct-messaging)
9. [Full Example: Chat Room](#9-full-example-chat-room)
10. [References](#10-references)

---

## 1. Overview

The `core/websocket` package provides a **Hub** that manages WebSocket
clients and rooms. Features:

- Automatic client tracking (connect/disconnect)
- Room-based grouping (join/leave)
- Broadcast to room or exclude sender
- Heartbeat ping/pong for dead connection detection
- OnJoin/OnLeave lifecycle callbacks

## 2. Creating a Hub

```go
import "github.com/RAiWorks/RapidGo/v2/core/websocket"

// Default: 30s ping, 10s pong timeout
hub := websocket.NewHub()

// Custom config
hub := websocket.NewHubWithConfig(websocket.HubConfig{
    PingInterval: 15 * time.Second,
    PongTimeout:  5 * time.Second,
})
```

Set `PingInterval: 0` to disable heartbeat entirely.

## 3. Handling Connections

Register the Hub handler on a Gin route:

```go
r.GET("/ws", hub.Handler(func(client *websocket.Client) {
    // This runs for the lifetime of the connection.
    // When this function returns, the client is automatically removed.
    for {
        _, msg, err := client.Conn.Read(ctx)
        if err != nil {
            break // disconnect
        }
        handleMessage(client, msg)
    }
}))
```

Each client gets a unique `client.ID` (UUID).

## 4. Rooms (Join / Leave)

```go
hub.Join(client, "chat:general")   // add client to room
hub.Leave(client, "chat:general")  // remove from room
```

- Rooms are created automatically on first Join.
- Rooms are deleted automatically when the last member leaves.
- Joining the same room twice is a no-op.
- When a client disconnects, `Remove()` is called automatically, which
  removes the client from **all** rooms.

### Inspect Rooms

```go
rooms := hub.Rooms()           // []string of active room names
clients := hub.Clients("chat") // []*Client in the room (nil if room doesn't exist)
```

## 5. Broadcasting

```go
// Send to everyone in the room
hub.Broadcast("chat:general", websocket.MessageText, []byte(`{"msg":"hello"}`))

// Send to everyone except the sender
hub.BroadcastOthers("chat:general", sender, websocket.MessageText, data)
```

## 6. Heartbeat & Dead Connection Detection

By default, the Hub sends a WebSocket Ping frame every 30 seconds. If
the client does not respond with a Pong within 10 seconds, the
connection is closed and the client is removed from all rooms.

This catches:
- Browser tab closed without clean close frame
- Network drops
- Mobile app backgrounded

Configure timing via `HubConfig`:

```go
hub := websocket.NewHubWithConfig(websocket.HubConfig{
    PingInterval: 10 * time.Second, // ping every 10s
    PongTimeout:  3 * time.Second,  // 3s to respond
})
```

## 7. Lifecycle Callbacks

React to room membership changes:

```go
hub.OnJoin(func(client *websocket.Client, room string) {
    log.Printf("client %s joined %s", client.ID, room)
    // e.g., broadcast "user joined" to the room
})

hub.OnLeave(func(client *websocket.Client, room string) {
    log.Printf("client %s left %s", client.ID, room)
    // e.g., broadcast "user left", update presence
})
```

- **OnJoin** fires once per room (not on duplicate joins).
- **OnLeave** fires on explicit `Leave()` and on `Remove()` (disconnect)
  for each room the client was in.

## 8. Direct Messaging

Send to a specific client by ID:

```go
err := hub.Send(clientID, websocket.MessageText, []byte(`{"type":"dm"}`))
if err != nil {
    // client not found
}
```

## 9. Full Example: Chat Room

```go
func setupChat(r *gin.Engine) {
    hub := websocket.NewHub()

    hub.OnJoin(func(c *websocket.Client, room string) {
        hub.Broadcast(room, websocket.MessageText,
            []byte(fmt.Sprintf(`{"event":"joined","client":"%s"}`, c.ID)))
    })

    hub.OnLeave(func(c *websocket.Client, room string) {
        hub.Broadcast(room, websocket.MessageText,
            []byte(fmt.Sprintf(`{"event":"left","client":"%s"}`, c.ID)))
    })

    r.GET("/ws/chat", hub.Handler(func(client *websocket.Client) {
        hub.Join(client, "chat:lobby")

        ctx := context.Background()
        for {
            _, msg, err := client.Conn.Read(ctx)
            if err != nil {
                break
            }
            hub.BroadcastOthers("chat:lobby", client, websocket.MessageText, msg)
        }
    }))
}
```

## 10. References

- [WebSocket Reference](../http/websocket.md) — low-level connection handling
- [Hub API](../../core/websocket/hub.go) — source code
