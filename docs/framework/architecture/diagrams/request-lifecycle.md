---
title: "Request Lifecycle Diagram"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Request Lifecycle Diagram

## Abstract

This diagram traces the complete journey of an HTTP request through the
framework — from initial receipt to final response — showing every
middleware, service, and component involved.

## Request Flow Diagram

```mermaid
sequenceDiagram
    participant Client
    participant Gin as Gin Router
    participant ReqID as RequestID MW
    participant CORS as CORS MW
    participant Rate as RateLimit MW
    participant Session as Session MW
    participant CSRF as CSRF MW
    participant Auth as Auth MW
    participant Controller
    participant Validation
    participant Service
    participant Model as Model (GORM)
    participant DB as Database
    participant Response as Response Helper

    Client->>Gin: HTTP Request
    Note over Gin: Route matching

    alt API Route (/api/*)
        Gin->>CORS: Apply CORS headers
        CORS->>Rate: Check rate limit
        Rate->>ReqID: Assign X-Request-ID
        ReqID->>Auth: Validate JWT token
        Auth->>Controller: Authorized request
    else Web Route (/)
        Gin->>Session: Load session from store
        Session->>CSRF: Validate CSRF token
        CSRF->>ReqID: Assign X-Request-ID
        ReqID->>Controller: Request with session
    end

    Controller->>Validation: Validate input
    alt Validation fails
        Validation-->>Controller: Errors
        Controller-->>Client: 422 (JSON errors or redirect with flash)
    else Validation passes
        Controller->>Service: Business logic
        Service->>Model: Query / mutate
        Model->>DB: SQL
        DB-->>Model: Result
        Model-->>Service: Data
        Service-->>Controller: Result
        Controller->>Response: Format response
        alt JSON API
            Response-->>Client: JSON { success, data, meta }
        else SSR Web
            Response-->>Client: HTML (rendered template)
        end
    end

    Note over Session: Session MW saves data<br/>after handler returns
```

## Web Route Flow (Detailed)

For a typical web form submission (`POST /users`):

```mermaid
flowchart TD
    A[POST /users] --> B[Session Middleware]
    B --> C[Load session from store]
    C --> D[CSRF Middleware]
    D --> E{CSRF token valid?}
    E -->|No| F[403 CSRF token mismatch]
    E -->|Yes| G[Request ID Middleware]
    G --> H[Controller: StoreUser]
    H --> I[Built-in Validator]
    I --> J{Valid?}
    J -->|No| K[Flash errors + old input]
    K --> L[Redirect to /users/create]
    J -->|Yes| M[UserService.Create]
    M --> N[GORM: db.Create]
    N --> O{DB error?}
    O -->|Yes| P[Flash error message]
    P --> L
    O -->|No| Q[Flash success message]
    Q --> R[Redirect to /users]
    R --> S[Session Middleware saves data]
```

## API Route Flow (Detailed)

For a typical API request (`GET /api/users`):

```mermaid
flowchart TD
    A[GET /api/users] --> B[CORS Middleware]
    B --> C[Rate Limit Middleware]
    C --> D{Under limit?}
    D -->|No| E[429 Too Many Requests]
    D -->|Yes| F[Request ID Middleware]
    F --> G[Auth Middleware]
    G --> H{Valid JWT?}
    H -->|No| I[401 Unauthorized]
    H -->|Yes| J[Controller: ListUsers]
    J --> K[Parse query params: page, per_page]
    K --> L[UserService.List]
    L --> M[helpers.Paginate]
    M --> N[GORM: Offset + Limit + Count]
    N --> O[responses.Paginated]
    O --> P["200 {success, data, meta}"]
```

## Middleware Execution Order

Middleware executes in registration order. The framework defines two
default groups:

### Web Group

```text
1. SessionMiddleware  → Load/save session per request
2. CSRFMiddleware     → Generate/validate CSRF tokens
3. RequestIDMiddleware → Assign X-Request-ID header
```

### API Group

```text
1. CORSMiddleware      → Set CORS headers
2. RateLimitMiddleware  → Enforce rate limits
3. RequestIDMiddleware  → Assign X-Request-ID header
```

Additional middleware (e.g., `auth`, `admin`) is applied per route
or per route group.

## References

- [System Overview](system-overview.md)
- [Routing](../../http/routing.md)
- [Middleware](../../http/middleware.md)
- [Controllers](../../http/controllers.md)
- [Responses](../../http/responses.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
