# 💬 Discussion: API Versioning

> **Feature**: `47` — API Versioning
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## What Are We Building?

A `core/router` extension that adds a `Version` method to the `Router` struct, returning a `*RouteGroup` scoped to a URL path prefix like `/api/v1`. This provides a clean, one-line API for grouping versioned routes — a thin convenience on top of the existing `Group()` method. A companion `DeprecatedVersion` method does the same but injects a middleware that sets a `Sunset` response header and a custom `X-API-Deprecated` header, signaling to clients that the version is deprecated per RFC 8594.

The framework provides the versioning primitives — the developer decides which versions to create, what routes belong to each, and when to deprecate.

---

## Why?

- **API evolution**: As APIs evolve, breaking changes need to coexist with older versions during migration periods
- **Roadmap**: Feature #47 in the project roadmap, depends on #07 (Router) — shipped
- **Framework gap**: Developers currently use `r.Group("/api/v1")` manually — no deprecation signaling, no standardized pattern
- **Convention**: A `Version` method establishes a clear project convention, making versioned APIs self-documenting

---

## Prior Art

| Framework | Approach |
|---|---|
| **Laravel** | URL prefix versioning via route groups; no built-in deprecation headers |
| **Django REST** | URL or header-based versioning via `DEFAULT_VERSIONING_CLASS`; negotiated per-request |
| **Rails** | Namespace-based routing with `scope` and `constraints`; URL prefix is the convention |
| **ASP.NET** | NuGet package with URL, query-string, header, and media-type strategies; `Sunset` header support |
| **Express (Node)** | Manual prefix groups or `express-api-versioning` middleware |

**Our approach**: URL prefix only (`/api/v1`, `/api/v2`). This is the most common, most visible, and easiest to reason about. Header-based or query-string versioning adds complexity with minimal benefit for most APIs. Deprecation uses standard `Sunset` and custom `X-API-Deprecated` headers.

---

## Constraints

1. **URL prefix only** — no header-based, query-string, or content-type versioning
2. **Thin wrapper** — `Version()` returns a `*RouteGroup`; all existing `RouteGroup` methods (Get, Post, APIResource, etc.) work on versioned groups
3. **No version negotiation** — the client chooses the version by URL; the server doesn't negotiate
4. **No automatic route copying** — each version has its own route registrations; the framework doesn't clone routes between versions
5. **Deprecation is opt-in** — developers create deprecated versions with `DeprecatedVersion()` which adds Sunset/deprecation headers via middleware
6. **No existing files modified** — new methods added to `core/router/version.go` only (keeps router.go clean)

---

## Decision Log

| # | Decision | Rationale |
|---|---|---|
| D1 | URL prefix versioning only | Most common convention. Clear, cacheable, visible in logs/docs. Header-based adds complexity without proportional benefit. |
| D2 | Return `*RouteGroup` from `Version()` | Developers get full access to all existing route registration methods (Get, Post, Group, APIResource, Use, etc.) — zero new API to learn. |
| D3 | Separate `DeprecatedVersion()` | Explicit is better than implicit. A developer must actively mark a version as deprecated — no configuration flags or date-based logic needed. |
| D4 | `Sunset` header per RFC 8594 | Industry standard for communicating API deprecation. Value is the sunset date in HTTP-date format. |
| D5 | `X-API-Deprecated` header | Simple boolean signal (`true`) that clients/monitoring can detect without parsing dates. |
| D6 | New file `version.go` | Keeps the versioning methods separate from the core router.go to maintain clean file organization. |

---

## Open Questions

| # | Question | Answer |
|---|---|---|
| Q1 | Should `Version()` enforce a naming format (e.g., must start with "v")? | ✅ No — the developer passes the full version string. `Version("v1")` produces `/api/v1`. No enforcement keeps it flexible (e.g., `Version("2024-01")` for date-based versioning). |
| Q2 | Should deprecated versions log a warning on each request? | ✅ No — header-based signaling is sufficient. Logging on every request would be too noisy in production. Developers can add custom middleware if they want logging. |
| Q3 | Should we store version metadata (sunset date, etc.) for discovery endpoints? | ✅ No — out of scope. A version listing endpoint is an application-level concern. The framework provides the routing primitives. |

---

## Next

Architecture → `47-api-versioning-architecture.md`
