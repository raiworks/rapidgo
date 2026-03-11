---
title: "Multi-Port Serving (Service Mode)"
version: "0.1.0"
status: "Final"
date: "2026-03-11"
last_updated: "2026-03-11"
authors:
  - "RAiWorks"
supersedes: ""
---

# Multi-Port Serving (Service Mode)

## Abstract

How to run Web, API, and WebSocket services on separate ports using the
`--mode` flag — suitable for microservice-style deployments while keeping
a single codebase.

## Table of Contents

1. [Overview](#1-overview)
2. [Modes](#2-modes)
3. [Single-Mode Usage](#3-single-mode-usage)
4. [Multi-Port Deployment](#4-multi-port-deployment)
5. [Port Configuration](#5-port-configuration)
6. [Docker / Process Manager](#6-docker--process-manager)
7. [References](#7-references)

---

## 1. Overview

RapidGo supports running all services in a single process (monolith) or
splitting them across separate processes on different ports. This is
controlled by the `--mode` flag on the `serve` command.

Modes are bitmask-based, so you can combine them:

| Mode | Routes Loaded |
|------|---------------|
| `all` | Web + API + WebSocket (default) |
| `web` | SSR templates, static files |
| `api` | JSON API endpoints |
| `ws` | WebSocket handlers |
| `api,ws` | API + WebSocket (no web) |

## 2. Modes

### Monolith (default)

```bash
go run cmd/main.go serve
# or
go run cmd/main.go serve --mode all
```

All routes registered on a single port (default `APP_PORT=8080`).

### Single Service

```bash
go run cmd/main.go serve --mode api
go run cmd/main.go serve --mode web
go run cmd/main.go serve --mode ws
```

Only the routes for that service are registered.

### Combined

```bash
go run cmd/main.go serve --mode api,ws
```

API and WebSocket routes on one process; Web on another.

## 3. Single-Mode Usage

When running a single mode, set the mode via flag or environment:

```bash
# Via CLI flag
go run cmd/main.go serve --mode api --port 8081

# Via environment variable
RAPIDGO_MODE=api APP_PORT=8081 go run cmd/main.go serve
```

Priority: CLI flag > `RAPIDGO_MODE` env var > default `all`.

## 4. Multi-Port Deployment

Run each service as a separate process:

**Terminal 1 — Web (SSR)**:
```bash
WEB_PORT=8080 go run cmd/main.go serve --mode web
```

**Terminal 2 — API**:
```bash
API_PORT=8081 go run cmd/main.go serve --mode api
```

**Terminal 3 — WebSocket**:
```bash
WS_PORT=8082 go run cmd/main.go serve --mode ws
```

Each process only loads its own routes and middleware.

## 5. Port Configuration

| Environment Variable | Used When | Default |
|---|---|---|
| `APP_PORT` | `--mode all` or fallback | `8080` |
| `WEB_PORT` | `--mode web` | `APP_PORT` |
| `API_PORT` | `--mode api` | `APP_PORT` |
| `WS_PORT` | `--mode ws` | `APP_PORT` |

The `--port` flag overrides everything for that process.

## 6. Docker / Process Manager

### Docker Compose

```yaml
services:
  web:
    build: .
    command: ["./app", "serve", "--mode", "web"]
    ports: ["8080:8080"]
    environment:
      WEB_PORT: "8080"

  api:
    build: .
    command: ["./app", "serve", "--mode", "api"]
    ports: ["8081:8081"]
    environment:
      API_PORT: "8081"

  ws:
    build: .
    command: ["./app", "serve", "--mode", "ws"]
    ports: ["8082:8082"]
    environment:
      WS_PORT: "8082"

  worker:
    build: .
    command: ["./app", "work"]
```

### Nginx Reverse Proxy

```nginx
upstream web  { server 127.0.0.1:8080; }
upstream api  { server 127.0.0.1:8081; }
upstream ws   { server 127.0.0.1:8082; }

server {
    listen 443 ssl;
    server_name yourapp.com;

    location / {
        proxy_pass http://web;
    }

    location /api/ {
        proxy_pass http://api;
    }

    location /ws {
        proxy_pass http://ws;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

## 7. References

- [Service Mode source](../../core/service/mode.go) — Mode type and parsing
- [Serve command](../../core/cli/serve.go) — CLI implementation
- [Service Mode Architecture](../service-mode-architecture.md) — design document
