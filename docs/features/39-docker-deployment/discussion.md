# Feature #39 — Docker Deployment: Discussion

## Overview

Provide production-ready Docker configuration for containerized deployment of the RGo app.

## Blueprint Reference

- Multi-stage Dockerfile: build with `golang:alpine`, run with `alpine`.
- `docker-compose.yml` with app, postgres, redis services.
- HEALTHCHECK using `/health` endpoint from Feature #36.

## Deliverables

1. `Dockerfile` — multi-stage build.
2. `docker-compose.yml` — app + postgres + redis.
3. `.dockerignore` — exclude unnecessary files from build context.
