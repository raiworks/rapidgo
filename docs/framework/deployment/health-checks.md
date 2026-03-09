---
title: "Health Checks"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Health Checks

## Abstract

This document describes the health check endpoints for liveness and
readiness probes, usable with Docker and Kubernetes.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Endpoints](#2-endpoints)
3. [Implementation](#3-implementation)
4. [Docker Integration](#4-docker-integration)
5. [Kubernetes Integration](#5-kubernetes-integration)
6. [References](#6-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Liveness probe** — Checks if the process is alive.
- **Readiness probe** — Checks if the app can serve requests
  (dependencies healthy).

## 2. Endpoints

| Endpoint | Type | Checks | Success | Failure |
|----------|------|--------|---------|---------|
| `GET /health` | Liveness | Process alive | 200 | — |
| `GET /health/ready` | Readiness | DB connection | 200 | 503 |

## 3. Implementation

```go
func HealthRoutes(r *gin.Engine, db *gorm.DB) {
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    r.GET("/health/ready", func(c *gin.Context) {
        sqlDB, err := db.DB()
        if err != nil {
            c.JSON(503, gin.H{"status": "error", "db": err.Error()})
            return
        }
        if err := sqlDB.Ping(); err != nil {
            c.JSON(503, gin.H{"status": "error", "db": err.Error()})
            return
        }
        c.JSON(200, gin.H{"status": "ready", "db": "connected"})
    })
}
```

### Response Examples

**Healthy:**

```json
{"status": "ok"}
```

**Ready:**

```json
{"status": "ready", "db": "connected"}
```

**Not Ready:**

```json
{"status": "error", "db": "connection refused"}
```

## 4. Docker Integration

In the Dockerfile:

```dockerfile
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s \
    CMD wget -qO- http://localhost:8080/health || exit 1
```

| Parameter | Value | Purpose |
|-----------|-------|---------|
| `interval` | 30 s | Time between checks |
| `timeout` | 5 s | Max wait for response |
| `start-period` | 5 s | Grace period on startup |

## 5. Kubernetes Integration

```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
    - name: app
      livenessProbe:
        httpGet:
          path: /health
          port: 8080
        initialDelaySeconds: 5
        periodSeconds: 30
      readinessProbe:
        httpGet:
          path: /health/ready
          port: 8080
        initialDelaySeconds: 5
        periodSeconds: 10
```

## 6. References

- [Build and Run](build-and-run.md)
- [Docker](docker.md)
- [Database](../data/database.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
