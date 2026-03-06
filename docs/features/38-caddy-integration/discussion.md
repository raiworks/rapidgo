# Feature #38 — Caddy Integration: Discussion

## Overview

Provide ready-to-use Caddy web server configuration for reverse-proxying to the RGo app, including automatic HTTPS, gzip, and static file serving.

## Blueprint Reference

Two options:
- **Option A (Embedded Caddy)**: Import `caddyserver/caddy/v2` as a Go library. Heavy dependency, complex.
- **Option B (External Caddyfile)**: Run Caddy as a separate process/container. Simple config files.

## Decision

**Option B** — external Caddyfile. Rationale:
- No new Go dependency (Caddy v2 pulls ~100+ transitive deps).
- Standard deployment pattern for production Go apps.
- Users can swap Caddy for nginx/traefik without touching Go code.
- Blueprint provides both; Option B is the practical choice.

## Deliverables

1. `Caddyfile` template in project root.
2. Documentation on Caddy setup (development + production).
3. `.env` additions: `CADDY_DOMAIN`.
