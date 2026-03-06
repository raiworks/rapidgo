# Feature #38 — Caddy Integration: Design

## Approach

External Caddyfile as a reverse proxy in front of the Go app.

## Caddyfile Template

```
# Production — automatic HTTPS via Let's Encrypt
{$CADDY_DOMAIN:localhost} {
    reverse_proxy localhost:{$APP_PORT:8080}
    encode gzip

    handle_path /static/* {
        root * ./resources/static
        file_server
    }

    handle_path /uploads/* {
        root * ./storage/uploads
        file_server
    }

    log {
        output stdout
    }
}
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CADDY_DOMAIN` | `localhost` | Domain for Caddy to serve (e.g., `example.com`) |
| `APP_PORT` | `8080` | Port the Go app listens on |

## File Layout

```
Caddyfile          — Template Caddyfile (project root)
```

## No Go Code Changes

This feature is purely configuration — no new Go packages or code modifications.
