---
title: "Build and Run"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Build and Run

## Abstract

This document covers the application entrypoint, HTTP server
configuration, graceful shutdown, and build commands.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Application Entrypoint](#2-application-entrypoint)
3. [Server Configuration](#3-server-configuration)
4. [Graceful Shutdown](#4-graceful-shutdown)
5. [Build Commands](#5-build-commands)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. Application Entrypoint

`cmd/main.go`:

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "yourframework/core/config"
    "yourframework/core/router"
    "yourframework/database"
)

func main() {
    config.Load()

    db, err := database.Connect()
    if err != nil {
        log.Fatal(err)
    }

    r := router.SetupRouter()

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

    srv := &http.Server{
        Addr:         ":" + port,
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        log.Printf("Server starting on :%s", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Forced shutdown: %v", err)
    }

    sqlDB, _ := db.DB()
    sqlDB.Close()

    log.Println("Server stopped")
}
```

## 3. Server Configuration

| Setting | Value | Purpose |
|---------|-------|---------|
| `Addr` | `":APP_PORT"` | Listen address from `.env` (default `8080`) |
| `ReadTimeout` | 15 s | Max time to read request (headers + body) |
| `WriteTimeout` | 15 s | Max time to write response |
| `IdleTimeout` | 60 s | Max time for keep-alive connections |

All timeout values **SHOULD** be tuned for production workloads.

## 4. Graceful Shutdown

The server handles graceful shutdown on `SIGINT` or `SIGTERM`:

1. Stop accepting new connections
2. Wait up to 30 seconds for active connections to finish
3. Close database connection pool
4. Exit cleanly

```text
Signal received (Ctrl+C or SIGTERM)
       │
       ▼
  srv.Shutdown(ctx)    ← 30-second deadline
       │
       ▼
  sqlDB.Close()
       │
       ▼
  "Server stopped"
```

## 5. Build Commands

### Development

```bash
go run cmd/main.go
```

### Production Binary

```bash
go build -o server ./cmd
./server
```

### Static Binary (for containers)

```bash
CGO_ENABLED=0 GOOS=linux go build -o server ./cmd
```

## 6. Security Considerations

- Server timeouts **MUST** be set to prevent slowloris attacks.
- The `APP_PORT` **SHOULD NOT** be exposed directly in production —
  use a reverse proxy (see [Caddy](caddy.md)).
- Database connections **MUST** be closed on shutdown to prevent
  connection leaks.

## 7. References

- [Application Lifecycle](../architecture/application-lifecycle.md)
- [Configuration](../core/configuration.md)
- [Docker](docker.md)
- [Caddy Integration](caddy.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
