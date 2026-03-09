---
title: "WebSocket"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# WebSocket

## Abstract

This document covers WebSocket support using the `coder/websocket`
library — connection handling, the read/write loop, and Gin
integration.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Library Choice](#2-library-choice)
3. [WebSocket Handler](#3-websocket-handler)
4. [Gin Integration](#4-gin-integration)
5. [API Reference](#5-api-reference)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **WebSocket** — A full-duplex communication protocol over a single
  TCP connection (RFC 6455).
- **Upgrade** — The HTTP handshake that switches the protocol from
  HTTP to WebSocket.

## 2. Library Choice

The framework uses **coder/websocket**
(`github.com/coder/websocket`).

> **Note:** The original `gorilla/websocket` is archived and
> unmaintained. `coder/websocket` (formerly `nhooyr.io/websocket`) is
> the actively maintained successor.

Key benefits:
- Actively maintained
- Supports `context.Context` natively
- Minimal API surface
- Works with `net/http` handlers directly

## 3. WebSocket Handler

The core handler accepts WebSocket connections, reads messages, and
echoes them back:

```go
package websocket

import (
    "context"
    "log"
    "net/http"

    "github.com/coder/websocket"
)

func HandleWS(w http.ResponseWriter, r *http.Request) {
    conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
        // Configure allowed origins in production
    })
    if err != nil {
        log.Printf("websocket accept error: %v", err)
        return
    }
    defer conn.CloseNow()

    ctx := r.Context()

    for {
        msgType, msg, err := conn.Read(ctx)
        if err != nil {
            log.Printf("websocket read error: %v", err)
            return
        }
        if err := conn.Write(ctx, msgType, msg); err != nil {
            log.Printf("websocket write error: %v", err)
            return
        }
    }
}
```

### Connection Lifecycle

```text
Client                           Server
  |                                |
  |--- HTTP Upgrade Request ------>|
  |<-- 101 Switching Protocols ----|
  |                                |
  |--- Message (text/binary) ---->|
  |<-- Message (text/binary) -----|
  |                                |
  |--- Close Frame -------------->|
  |<-- Close Frame ---------------|
```

## 4. Gin Integration

Gin handlers receive `*gin.Context`, but `coder/websocket` accepts
standard `http.ResponseWriter` and `*http.Request`. Use the Gin
context to access both:

```go
r.GET("/ws", func(c *gin.Context) {
    websocket.HandleWS(c.Writer, c.Request)
})
```

### With Authentication

Apply middleware before the WebSocket route:

```go
ws := r.Group("/ws", middleware.Resolve("auth"))
{
    ws.GET("/chat", func(c *gin.Context) {
        websocket.HandleWS(c.Writer, c.Request)
    })
}
```

## 5. API Reference

### `websocket.Accept`

Upgrades the HTTP connection to a WebSocket connection.

```go
conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
    // InsecureSkipVerify: true, // only for development
})
```

### `conn.Read`

Reads the next message from the connection. Blocks until a message
is received or the context is cancelled.

```go
msgType, msg, err := conn.Read(ctx)
// msgType: websocket.MessageText or websocket.MessageBinary
```

### `conn.Write`

Writes a message to the connection.

```go
err := conn.Write(ctx, websocket.MessageText, []byte("hello"))
```

### `conn.Close`

Sends a close frame with a status code and reason.

```go
conn.Close(websocket.StatusNormalClosure, "goodbye")
```

### `conn.CloseNow`

Immediately closes the connection without a close handshake.
Typically used in `defer`.

```go
defer conn.CloseNow()
```

## 6. Security Considerations

- **Origin validation:** In production, configure
  `AcceptOptions.OriginPatterns` to restrict allowed origins and
  prevent Cross-Site WebSocket Hijacking.
- **Authentication:** WebSocket routes **SHOULD** be protected by
  auth middleware. Token validation **MUST** happen before the
  upgrade.
- **Rate limiting:** Consider limiting the number of concurrent
  WebSocket connections per client.
- **Message size:** Set maximum message size limits to prevent
  memory exhaustion attacks.
- **TLS:** WebSocket connections **SHOULD** use `wss://` (TLS) in
  production.

## 7. References

- [coder/websocket repository](https://github.com/coder/websocket)
- [RFC 6455 — The WebSocket Protocol](https://tools.ietf.org/html/rfc6455)
- [Routing](routing.md)
- [Middleware](middleware.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
