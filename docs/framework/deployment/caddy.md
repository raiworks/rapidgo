---
title: "Caddy Integration"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Caddy Integration

## Abstract

This document covers two approaches to using Caddy with the
framework — as an embedded Go library or as an external reverse
proxy — including auto-HTTPS, static file serving, and Docker
configuration.

## Table of Contents

1. [Terminology](#1-terminology)
2. [When to Use Caddy](#2-when-to-use-caddy)
3. [Option A: Embedded Caddy](#3-option-a-embedded-caddy)
4. [Option B: External Reverse Proxy](#4-option-b-external-reverse-proxy)
5. [Docker Compose with Caddy](#5-docker-compose-with-caddy)
6. [Choosing an Approach](#6-choosing-an-approach)
7. [Security Considerations](#7-security-considerations)
8. [References](#8-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. When to Use Caddy

Caddy provides automatic HTTPS via Let's Encrypt. Use it when:

- You need TLS in production without manual certificate management
- You want static file serving offloaded from Go
- You need a reverse proxy with gzip/logging

Caddy is **OPTIONAL** — the framework runs standalone on any port.

## 3. Option A: Embedded Caddy

Use `caddyserver/caddy/v2` as a Go library. The Go app starts Caddy
programmatically.

### `StartWithCaddy`

```go
package server

import (
    "fmt"

    "github.com/caddyserver/caddy/v2"
    "github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func StartWithCaddy(appPort, domain string) error {
    caddyConfig := fmt.Sprintf(`
        %s {
            reverse_proxy localhost:%s
        }
    `, domain, appPort)

    adapter := caddyfile.Adapter{}
    cfgJSON, _, err := adapter.Adapt([]byte(caddyConfig), nil)
    if err != nil {
        return fmt.Errorf("caddy config error: %w", err)
    }

    return caddy.Load(cfgJSON, false)
}
```

### Usage in `main.go`

```go
func main() {
    config.Load()

    go func() {
        r := router.SetupRouter()
        r.Run(":8081") // Internal port
    }()

    if os.Getenv("CADDY_ENABLED") == "true" {
        domain := os.Getenv("CADDY_DOMAIN")
        if err := server.StartWithCaddy("8081", domain); err != nil {
            log.Fatal(err)
        }
        select {} // Block forever
    }
}
```

### `.env` Variables

```env
CADDY_ENABLED=false
CADDY_DOMAIN=localhost
```

## 4. Option B: External Reverse Proxy

Run Caddy separately — simpler deployment, no Go dependency on Caddy.

### `Caddyfile`

```caddyfile
# Production — automatic HTTPS via Let's Encrypt
example.com {
    reverse_proxy localhost:8080
    encode gzip

    # Serve static files directly (bypass Go app)
    handle_path /static/* {
        root * ./resources/static
        file_server
    }

    handle_path /uploads/* {
        root * ./storage/uploads
        file_server
    }

    log {
        output file /var/log/caddy/access.log
    }
}

# Development
# localhost:80 {
#     reverse_proxy localhost:8080
# }
```

### Benefits

- Static files served directly by Caddy (faster, no Go overhead)
- Automatic TLS certificate provisioning and renewal
- Gzip compression
- Access logging

## 5. Docker Compose with Caddy

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    env_file: .env
    depends_on:
      - db

  caddy:
    image: caddy:2-alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - ./resources/static:/srv/static
      - ./storage/uploads:/srv/uploads
      - caddy_data:/data
      - caddy_config:/config
    depends_on:
      - app

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASS}
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
  caddy_data:
  caddy_config:
```

## 6. Choosing an Approach

| Factor | Embedded | External |
|--------|----------|----------|
| Complexity | Higher (Go library) | Lower (separate process) |
| Deployment | Single binary | Two processes |
| Static files | Served by Go | Served by Caddy |
| TLS config | Programmatic | Caddyfile |
| Recommended for | Single-binary deploys | Standard deployments |

## 7. Security Considerations

- Caddy handles TLS automatically — **MUST NOT** store private keys
  manually.
- When using external Caddy, the Go app **SHOULD** only listen on
  `localhost` or an internal network.
- `CADDY_DOMAIN` **MUST** match your actual domain for TLS to work.

## 8. References

- [Caddy](https://caddyserver.com/)
- [Build and Run](build-and-run.md)
- [Docker](docker.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
