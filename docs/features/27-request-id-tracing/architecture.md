# Feature #27 — Request ID / Tracing: Architecture

## Status: Already implemented in Feature #08

No new files or modifications required.

## Existing implementation

| File | Purpose | Shipped in |
|------|---------|------------|
| `core/middleware/request_id.go` | `RequestID()` middleware + `generateUUID()` helper | #08 |
| `app/providers/middleware_provider.go` | `"requestid"` alias + `"global"` group membership | #08 |

## Request flow

```
Incoming request
    │
    ├── Has X-Request-ID header? → Use it
    └── No header? → generateUUID() → UUID v4
    │
    ├── c.Set("request_id", id)     ← available to handlers
    ├── c.Header("X-Request-ID", id) ← echoed in response
    └── c.Next()
```

## Blueprint comparison

| Blueprint spec | Our implementation | Match |
|---|---|---|
| `RequestIDMiddleware()` | `RequestID()` | ✅ (name follows our convention) |
| `rand.Read(b)` → hex string | `generateUUID()` → UUID v4 | ✅ (exceeds spec) |
| `c.GetHeader("X-Request-ID")` | `c.GetHeader(requestIDHeader)` | ✅ |
| `c.Set("request_id", id)` | `c.Set("request_id", id)` | ✅ |
| `c.Header("X-Request-ID", id)` | `c.Header(requestIDHeader, id)` | ✅ |
| `c.Next()` | `c.Next()` | ✅ |
