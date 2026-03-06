# Feature #27 — Request ID / Tracing: Discussion

## What problem does this solve?

Every request needs a unique identifier for log correlation, debugging, and distributed tracing. When a request fails, the ID links log entries across middleware, handlers, and services.

## Why now?

This feature was already fully implemented in Feature #08 (Middleware Infrastructure). The roadmap listed it as a separate item because the blueprint has a dedicated "Request ID / Tracing" section, but our #08 implementation already covers 100% of the blueprint specification.

## What does the blueprint specify?

- `RequestIDMiddleware()` that reads `X-Request-ID` from the incoming request header.
- If missing, generates a random 16-byte hex string.
- Stores the ID in Gin context as `"request_id"`.
- Sets `X-Request-ID` response header.
- Calls `c.Next()`.

## What do we already have?

- `RequestID()` in `core/middleware/request_id.go` (shipped in #08).
- Reads `X-Request-ID` from incoming header — preserves if present.
- If missing, generates a **UUID v4** (superior to blueprint's raw hex — proper format, version/variant bits).
- Stores as `"request_id"` in Gin context.
- Sets `X-Request-ID` response header.
- Alias `"requestid"` registered in `middleware_provider.go`.
- Part of the `"global"` middleware group.
- 2 tests: TC-08 (generates UUID v4), TC-09 (preserves existing).
- `generateUUID()` helper tested separately (TC-14).

## Design decision

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Status | Already shipped (#08) | Full blueprint coverage, no gaps |
| ID format | UUID v4 | Exceeds blueprint (raw hex) — standard format, version/variant bits |
| Code changes | None needed | Implementation is complete |

## What is out of scope?

- Distributed tracing (OpenTelemetry, Jaeger) — not in blueprint.
- Logger integration for automatic request_id injection — not in blueprint.
