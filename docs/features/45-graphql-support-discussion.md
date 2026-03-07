# 💬 Discussion: GraphQL Support

> **Feature**: `45` — GraphQL Support
> **Status**: ✅ COMPLETE
> **Date**: 2026-03-07

---

## What Are We Building?

A `core/graphql` package that provides an HTTP handler for serving GraphQL queries over the Gin router. The package wraps the `graphql-go/graphql` execution engine, handling request parsing (POST JSON body and GET query params), query execution, and JSON response formatting per the GraphQL over HTTP specification. A built-in Playground handler serves an interactive GraphiQL IDE for development.

The framework provides the HTTP transport layer — schema definition and resolver logic remain application-level concerns using the `graphql-go/graphql` library directly.

---

## Why?

- **API flexibility**: GraphQL lets clients request exactly the data they need, reducing over-fetching and under-fetching compared to REST
- **Roadmap**: Feature #45 in the project roadmap, depends on #07 (Router) — shipped
- **Framework gap**: No GraphQL capability exists today — developers must integrate a library and write boilerplate HTTP handling themselves
- **Complement to REST**: RapidGo already supports REST via resource routes; GraphQL adds a second API paradigm without replacing REST

---

## Prior Art

| Framework | Approach |
|---|---|
| **Laravel (Lighthouse)** | Schema-first with SDL files; auto-generates resolvers from Eloquent models; custom directives |
| **Django (Graphene)** | Code-first; Python classes map to GraphQL types; integrates with Django ORM |
| **Rails (graphql-ruby)** | Code-first with class-based types; schema stitching; built-in Relay support |
| **Go (gqlgen)** | Schema-first with code generation; generates Go resolver stubs from SDL |
| **Go (graphql-go)** | Code-first; runtime schema definition in Go; no code generation required |

**Our approach**: Provide a thin HTTP handler wrapping `graphql-go/graphql`. No code generation, no SDL parsing, no ORM integration. The framework handles HTTP transport; the developer defines schemas and resolvers in application code using the graphql-go API directly. This keeps the core package small, testable, and free of magic.

---

## Constraints

1. **HTTP transport only** — the core package handles request/response, not schema design or resolver wiring
2. **Single query per request** — no batched query support (array of operations)
3. **No file uploads** — multipart GraphQL requests are out of scope
4. **No subscriptions** — WebSocket-based GraphQL subscriptions are separate from HTTP queries (Feature #48 covers WebSocket rooms)
5. **No ORM integration** — resolvers interact with the database through existing services/models; no automatic type generation from GORM models
6. **CDN-hosted Playground** — the GraphiQL IDE loads JavaScript from unpkg.com CDN; suitable for development only
7. **Dependencies**: `graphql-go/graphql` (MIT) — runtime GraphQL execution engine

---

## Decision Log

| # | Decision | Rationale |
|---|---|---|
| D1 | Use `graphql-go/graphql` over `gqlgen` | No code generation step required. Runtime schema definition is simpler for framework integration. Developers can use `gqlgen` if they prefer, but the core package doesn't depend on it. |
| D2 | HTTP handler only, not schema builder | Keeps the core package small and focused. Schema definition is highly application-specific — the framework shouldn't abstract it. |
| D3 | Support both POST (JSON) and GET | Per the GraphQL over HTTP spec. POST for mutations, GET for queries (cacheable). |
| D4 | Pass `gin.Context` through `context.Context` | Resolvers need access to request context (auth, headers, request ID). A `FromContext` helper makes this ergonomic. |
| D5 | GraphiQL from CDN, not embedded | Avoids embedding ~2MB of JavaScript in the binary. CDN-loaded is standard practice in Go GraphQL libraries. |
| D6 | No new service mode | GraphQL is an API concern — served under the existing `ModeAPI` flag, not a separate mode. Keeps the service mode system simple. |

---

## Open Questions

| # | Question | Answer |
|---|---|---|
| Q1 | Should we add a `ModeGraphQL` service mode? | ✅ No — GraphQL is an API-level concern. If the user wants GraphQL on a separate port, they can create a custom mode or mount it under the existing API routes. |
| Q2 | Should we auto-generate types from GORM models? | ✅ No — automatic type generation adds significant complexity and tight ORM coupling. Developers define GraphQL types explicitly using `graphql-go/graphql`. |
| Q3 | Should the Playground be disabled in production? | ✅ No — the Playground handler is just a function. The developer decides whether to register it in their routes. They can gate it behind `APP_ENV` checks in `routes/graphql.go` if desired. |
| Q4 | Should we support query complexity limits? | ✅ No — query complexity analysis is an advanced concern. Developers can use `graphql-go/graphql` middleware or custom validation in their resolvers. Rate limiting already exists via Feature #26. |

---

## Next

Architecture → `45-graphql-support-architecture.md`
