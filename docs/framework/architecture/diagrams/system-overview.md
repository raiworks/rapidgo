---
title: "System Overview Diagram"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# System Overview Diagram

## Abstract

This diagram shows the high-level system architecture, illustrating
the major layers (HTTP, Application, Core, Data, Infrastructure) and
how they connect.

## Diagram

```mermaid
graph TB
    subgraph "Clients"
        Browser["🌐 Browser"]
        API_Client["📱 API Client"]
        WS_Client["⚡ WebSocket Client"]
    end

    subgraph "Entry Points"
        Caddy["Caddy<br/>(Optional Reverse Proxy)"]
        Static["Static Files<br/>/static, /uploads"]
    end

    subgraph "HTTP Layer"
        Gin["Gin Router"]
        MW_Auth["Auth Middleware"]
        MW_CSRF["CSRF Middleware"]
        MW_CORS["CORS Middleware"]
        MW_Rate["Rate Limiter"]
        MW_ReqID["Request ID"]
        MW_Session["Session Middleware"]
    end

    subgraph "Application Layer"
        Controllers["Controllers<br/>(HTTP handlers)"]
        Services["Services<br/>(Business logic)"]
        Helpers["Helpers<br/>(Utilities)"]
        Validation["Validation<br/>(Built-in + Struct)"]
        Responses["Response Helpers<br/>(JSON envelope)"]
        Views["Views<br/>(html/template)"]
    end

    subgraph "Core Framework"
        Container["Service Container<br/>(Bind, Singleton, Make)"]
        Providers["Service Providers<br/>(Register, Boot)"]
        Config["Configuration<br/>(.env, Env detection)"]
        Logger["Logger<br/>(slog JSON)"]
        Crypto["Crypto<br/>(AES, HMAC, Hash)"]
        Events["Event Dispatcher<br/>(Listen, Dispatch)"]
        I18n["i18n<br/>(Translations)"]
    end

    subgraph "Data Layer"
        GORM["GORM ORM"]
        SessionMgr["Session Manager<br/>(5 backends)"]
        CacheMgr["Cache Manager<br/>(Redis, Memory)"]
        StorageMgr["File Storage<br/>(Local, S3)"]
        Mailer["Mailer<br/>(SMTP)"]
        WS_Hub["WebSocket Hub"]
    end

    subgraph "Infrastructure"
        PG[("PostgreSQL")]
        MySQL[("MySQL")]
        SQLite[("SQLite")]
        Redis[("Redis")]
        S3[("S3 / Disk")]
        SMTP[("SMTP Server")]
    end

    Browser --> Caddy
    API_Client --> Caddy
    WS_Client --> Caddy
    Caddy --> Static
    Caddy --> Gin

    Gin --> MW_Auth
    Gin --> MW_CSRF
    Gin --> MW_CORS
    Gin --> MW_Rate
    Gin --> MW_ReqID
    Gin --> MW_Session

    MW_Auth --> Controllers
    MW_Session --> Controllers
    Controllers --> Services
    Controllers --> Validation
    Controllers --> Responses
    Controllers --> Views
    Services --> Helpers
    Services --> GORM
    Services --> CacheMgr
    Services --> Events
    Services --> Mailer
    Services --> StorageMgr

    Container --> Providers
    Config --> Container

    GORM --> PG
    GORM --> MySQL
    GORM --> SQLite
    SessionMgr --> PG
    SessionMgr --> Redis
    CacheMgr --> Redis
    StorageMgr --> S3
    Mailer --> SMTP
    WS_Hub --> WS_Client
```

## Layer Summary

| Layer | Responsibility |
|-------|---------------|
| **Clients** | Browsers, mobile apps, API consumers, WebSocket clients |
| **Entry Points** | Caddy reverse proxy (optional), static file serving |
| **HTTP Layer** | Gin router, middleware pipeline (auth, CSRF, CORS, rate limit, request ID, session) |
| **Application Layer** | Controllers, services, helpers, validation, response formatting, views |
| **Core Framework** | Service container, providers, configuration, logging, crypto, events, i18n |
| **Data Layer** | GORM ORM, session manager, cache manager, file storage, mailer, WebSocket hub |
| **Infrastructure** | PostgreSQL, MySQL, SQLite, Redis, S3/disk, SMTP server |

## References

- [Architecture Overview](../overview.md)
- [Request Lifecycle](request-lifecycle.md)
- [Service Container](service-container.md)
- [Data Flow](data-flow.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
