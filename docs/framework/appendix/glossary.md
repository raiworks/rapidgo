---
title: "Glossary"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Glossary

## Abstract

Definitions of all domain-specific terms used across the
architecture documentation.

## Terms

**AES-GCM** — Advanced Encryption Standard in Galois/Counter Mode.
Authenticated encryption algorithm used for data-at-rest encryption.

**Binding** — Registering a constructor or value in the service
container for later resolution.

**CORS** — Cross-Origin Resource Sharing. HTTP headers controlling
which origins can access API endpoints.

**CRUD** — Create, Read, Update, Delete. The four basic data
operations.

**CSRF** — Cross-Site Request Forgery. An attack where a malicious
site submits a request on behalf of an authenticated user.

**DI** — Dependency Injection. Providing dependencies to a component
rather than having it construct them.

**DSN** — Data Source Name. A connection string identifying a database.

**Factory** — A container binding that creates a new instance on every
resolution.

**Fallback Locale** — The default language used when a translation
key is missing in the requested locale.

**Flash Message** — Session data that persists for exactly one request,
used for success/error feedback after redirects.

**GORM** — The Go ORM library (`gorm.io/gorm`) used for database
interaction.

**Graceful Shutdown** — Stopping a server by finishing in-flight
requests before exiting.

**Handler** — A function that processes an HTTP request and writes a
response. In Gin: `gin.HandlerFunc`.

**Health Check** — An endpoint that reports application and dependency
status for monitoring/orchestration.

**HMAC** — Hash-based Message Authentication Code. A keyed hash for
verifying message integrity and authenticity.

**Hook** — A GORM callback triggered before or after a database
operation (e.g., `BeforeCreate`).

**i18n** — Internationalization. Supporting multiple languages in an
application.

**IoC** — Inversion of Control. A design principle where the
framework controls the flow, not the user code.

**JWT** — JSON Web Token. A compact, self-contained token for
stateless authentication.

**Liveness Probe** — A health check that verifies the process is
running.

**Locale** — A language/region identifier (e.g., `en`, `es`).

**Middleware** — A function that wraps an HTTP handler, executing
logic before and/or after the handler.

**Migration** — A database schema change, typically run via
`AutoMigrate` or manual SQL files.

**MVC** — Model-View-Controller. An architectural pattern separating
data, presentation, and input logic.

**ORM** — Object-Relational Mapping. A technique for querying
databases using Go structs instead of raw SQL.

**Provider** — See Service Provider.

**Readiness Probe** — A health check verifying the app can serve
requests (e.g., database is connected).

**Route Group** — A set of routes sharing a common prefix and
middleware stack.

**Seeder** — A function that populates the database with initial or
test data.

**Service Container** — A dependency injection container that stores
and resolves application services.

**Service Provider** — A struct implementing `Register()` and
`Boot()` that configures services in the container.

**Singleton** — A container binding that creates one shared instance,
reused on every resolution.

**SPA** — Single-Page Application. A frontend architecture where the
browser handles routing.

**SSR** — Server-Side Rendering. Generating HTML on the server and
sending it to the browser.

## References

- [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119) — Key words for
  use in RFCs

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
