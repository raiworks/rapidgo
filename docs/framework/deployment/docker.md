---
title: "Docker"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Docker

## Abstract

This document covers Docker containerization — multi-stage
Dockerfile, docker-compose setup with PostgreSQL and Redis, and
deployment commands. Docker is **OPTIONAL**.

## Table of Contents

1. [Terminology](#1-terminology)
2. [When to Use Docker](#2-when-to-use-docker)
3. [Dockerfile (Multi-Stage)](#3-dockerfile-multi-stage)
4. [Docker Compose](#4-docker-compose)
5. [Deployment Commands](#5-deployment-commands)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. When to Use Docker

Docker is **not required** for running the framework. Small apps can
run directly as a compiled binary. Use Docker when you need:

- Containerized deployment (cloud, CI/CD)
- Multi-service orchestration (app + database + Redis)
- Reproducible builds across environments

## 3. Dockerfile (Multi-Stage)

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/resources ./resources
COPY --from=builder /app/.env .env

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s \
    CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./server"]
```

### Build Stages

| Stage | Base Image | Purpose |
|-------|-----------|---------|
| Builder | `golang:1.22-alpine` | Compile Go binary |
| Runtime | `alpine:3.19` | Minimal production image |

### Key Decisions

- `CGO_ENABLED=0` — static binary, no C dependencies
- `ca-certificates` — required for HTTPS outbound calls
- `tzdata` — timezone support
- `HEALTHCHECK` — Docker-native liveness probe via `/health`

## 4. Docker Compose

Standalone setup (without Caddy):

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    env_file: .env
    volumes:
      - ./storage:/app/storage
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASS}
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    restart: unless-stopped

volumes:
  pgdata:
```

### Service Dependencies

The `app` service depends on `db` with `condition: service_healthy`,
ensuring the database is ready before the app starts.

## 5. Deployment Commands

```bash
# Start all services
docker compose up -d

# Rebuild and start
docker compose up -d --build

# View logs
docker compose logs -f app

# Stop all services
docker compose down

# Stop and remove volumes
docker compose down -v
```

## 6. Security Considerations

- `.env` files **MUST NOT** be baked into public images — use
  `env_file` in compose or runtime secrets.
- Database ports **SHOULD NOT** be exposed in production compose
  files — remove the `ports` mapping for `db` and `redis`.
- The runtime image **SHOULD** use a non-root user.
- Use `restart: unless-stopped` for production resilience.

## 7. References

- [Build and Run](build-and-run.md)
- [Caddy Integration](caddy.md)
- [Health Checks](health-checks.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
