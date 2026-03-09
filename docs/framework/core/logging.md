---
title: "Logging"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Logging

## Abstract

This document covers the framework's structured logging system built
on Go's standard `log/slog` package (Go 1.21+). It describes setup,
configuration, log levels, structured output, and environment-specific
behavior.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Overview](#2-overview)
3. [Logger Setup](#3-logger-setup)
4. [Log Levels](#4-log-levels)
5. [Structured Logging](#5-structured-logging)
6. [Environment-Specific Output](#6-environment-specific-output)
7. [Alternative Libraries](#7-alternative-libraries)
8. [Security Considerations](#8-security-considerations)
9. [References](#9-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

## 2. Overview

The framework uses `log/slog` from Go's standard library for structured
logging. This provides:

- **JSON output** — Machine-parseable logs for production.
- **Structured fields** — Key-value pairs instead of string formatting.
- **Log levels** — Debug, Info, Warn, Error.
- **Zero dependencies** — Part of the Go standard library.

## 3. Logger Setup

```go
package logger

import (
    "log/slog"
    "os"
)

func Setup() *slog.Logger {
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    })
    logger := slog.New(handler)
    slog.SetDefault(logger)
    return logger
}
```

The `Setup()` function:

1. Creates a `JSONHandler` that writes to stdout.
2. Sets the minimum log level to `Info`.
3. Sets this logger as the global default via `slog.SetDefault()`.
4. Returns the logger for dependency injection if needed.

## 4. Log Levels

| Level | Constant | When to Use |
|-------|----------|-------------|
| Debug | `slog.LevelDebug` | Detailed diagnostic information — disabled in production |
| Info | `slog.LevelInfo` | General operational messages — server started, request handled |
| Warn | `slog.LevelWarn` | Potential issues — deprecated feature used, slow query |
| Error | `slog.LevelError` | Failures — database error, external service failure |

The minimum level **SHOULD** be configured per environment:

- **Development:** `slog.LevelDebug` — see everything.
- **Production:** `slog.LevelInfo` — skip debug noise.
- **Testing:** `slog.LevelWarn` — only problems.

## 5. Structured Logging

Use key-value pairs with every log call:

```go
slog.Info("server started", "port", os.Getenv("APP_PORT"))
slog.Error("database error", "err", err)
slog.Info("user created", "user_id", user.ID, "email", user.Email)
slog.Warn("slow query", "duration_ms", elapsed, "query", queryName)
```

This produces JSON output:

```json
{"level":"INFO","msg":"server started","port":"8080"}
{"level":"ERROR","msg":"database error","err":"connection refused"}
{"level":"INFO","msg":"user created","user_id":42,"email":"user@example.com"}
```

### With request context

Include the request ID for correlation:

```go
slog.Error("request error",
    "path", c.Request.URL.Path,
    "method", c.Request.Method,
    "request_id", c.GetString("request_id"),
    "err", err,
)
```

## 6. Environment-Specific Output

### Development

- Output: Pretty-printed to **stdout**.
- Level: `Debug` and above.
- Format: Use `slog.NewTextHandler()` for human-readable output.

```go
if config.IsDevelopment() {
    handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    })
}
```

### Production

- Output: **JSON** to stdout (for log aggregation) or file.
- Level: `Info` and above.
- Format: `slog.NewJSONHandler()` — structured, parseable.

```go
if config.IsProduction() {
    handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    })
}
```

## 7. Alternative Libraries

While `log/slog` is the recommended default, the following alternatives
are well-suited for Go web applications:

| Library | Import Path | Notes |
|---------|-------------|-------|
| **zerolog** | `github.com/rs/zerolog` | Zero-allocation, extremely fast |
| **zap** | `go.uber.org/zap` | Structured, high-performance |

These **MAY** be used as drop-in replacements when performance
profiling indicates logging is a bottleneck.

## 8. Security Considerations

- **MUST NOT** log sensitive data: passwords, API keys, JWT tokens,
  credit card numbers, or personal identifiable information (PII).
- **MUST NOT** log full request bodies that may contain credentials.
- Log files in `storage/logs/` **SHOULD** have restricted filesystem
  permissions (0600 or 0640).
- In production, logs **SHOULD** be forwarded to a centralized
  logging system rather than stored on disk indefinitely.

## 9. References

- [Error Handling](error-handling.md)
- [Configuration](configuration.md)
- [Request ID](../security/request-id.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
