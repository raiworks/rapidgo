---
title: "Service Container Diagram"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Service Container Diagram

## Abstract

This diagram illustrates the service container's internal structure,
the provider registration/boot lifecycle, and how services are
resolved at runtime.

## Container Architecture

```mermaid
graph TB
    subgraph "Service Container"
        Bindings["bindings<br/>map[string]Factory"]
        Instances["instances<br/>map[string]interface{}"]
        Mutex["sync.RWMutex<br/>(thread-safe)"]
    end

    subgraph "Registration Methods"
        Bind["Bind(name, factory)<br/>Transient — new instance each time"]
        Singleton["Singleton(name, factory)<br/>Shared — created once, reused"]
        Instance["Instance(name, obj)<br/>Pre-created — stored directly"]
    end

    subgraph "Resolution Methods"
        Make["Make(name)<br/>Returns interface{}"]
        MustMake["MustMake[T](name)<br/>Generic typed resolution"]
        Has["Has(name)<br/>Check if registered"]
    end

    Bind --> Bindings
    Singleton --> Bindings
    Singleton -.->|"first call stores"| Instances
    Instance --> Instances

    Make --> Instances
    Make --> Bindings
    MustMake --> Make
    Has --> Bindings
    Has --> Instances
```

## Provider Lifecycle

```mermaid
sequenceDiagram
    participant Main as main()
    participant App as App
    participant Container as Container
    participant DBP as DatabaseProvider
    participant SP as SessionProvider
    participant CP as CacheProvider
    participant MP as MailProvider
    participant EP as EventProvider

    Note over Main,EP: Phase 1: Registration

    Main->>App: Register(DatabaseProvider)
    App->>DBP: Register(container)
    DBP->>Container: Singleton("db", dbFactory)

    Main->>App: Register(SessionProvider)
    App->>SP: Register(container)
    SP->>Container: Singleton("session", sessionFactory)

    Main->>App: Register(CacheProvider)
    App->>CP: Register(container)
    CP->>Container: Singleton("cache", cacheFactory)

    Main->>App: Register(MailProvider)
    App->>MP: Register(container)
    MP->>Container: Singleton("mail", mailFactory)

    Main->>App: Register(EventProvider)
    App->>EP: Register(container)
    EP->>Container: Singleton("events", eventsFactory)

    Note over Main,EP: Phase 2: Boot

    Main->>App: Boot()
    App->>DBP: Boot(container)
    Note over DBP: Run auto-migrations<br/>if enabled
    App->>SP: Boot(container)
    App->>CP: Boot(container)
    App->>MP: Boot(container)
    App->>EP: Boot(container)
```

## Resolution Flow

```mermaid
flowchart TD
    A["Make('db')"] --> B{Instance exists?}
    B -->|Yes| C[Return cached instance]
    B -->|No| D{Binding exists?}
    D -->|No| E["panic: service not found"]
    D -->|Yes| F[Call factory function]
    F --> G{Is Singleton?}
    G -->|Yes| H[Store in instances map]
    H --> I[Return instance]
    G -->|No| I[Return instance]
```

## Binding Types

```mermaid
graph LR
    subgraph "Bind (Transient)"
        B1["Make('logger')"] --> B2["factory(c)"]
        B3["Make('logger')"] --> B4["factory(c)"]
        B2 --> B5["Instance A"]
        B4 --> B6["Instance B"]
        style B5 fill:#f9f,stroke:#333
        style B6 fill:#f9f,stroke:#333
    end

    subgraph "Singleton (Shared)"
        S1["Make('db')"] --> S2["factory(c)"]
        S2 --> S3["Instance X"]
        S4["Make('db')"] --> S3
        style S3 fill:#9ff,stroke:#333
    end

    subgraph "Instance (Pre-created)"
        I1["Instance('config', obj)"] --> I2["obj stored directly"]
        I3["Make('config')"] --> I2
        style I2 fill:#9f9,stroke:#333
    end
```

## Built-in Service Bindings

| Service Name | Type | Provider | Returns |
|-------------|------|----------|---------|
| `"db"` | Singleton | `DatabaseProvider` | `*gorm.DB` |
| `"session"` | Singleton | `SessionProvider` | `*session.Manager` |
| `"cache"` | Singleton | `CacheProvider` | `cache.Store` (Redis or Memory) |
| `"mail"` | Singleton | `MailProvider` | `*mail.Mailer` |
| `"events"` | Singleton | `EventProvider` | `*events.Dispatcher` |

## Custom Provider Example

```mermaid
sequenceDiagram
    participant User as Developer
    participant CLI as CLI (make:provider)
    participant File as providers/paymentprovider.go
    participant Main as main.go
    participant Container

    User->>CLI: framework make:provider PaymentProvider
    CLI->>File: Generate boilerplate

    User->>Main: application.Register(&PaymentProvider{})
    Main->>Container: Singleton("payment", stripeFactory)

    Note over Container: Available via Make("payment")<br/>or MustMake[*StripeGateway](c, "payment")
```

## References

- [Service Container](../../core/service-container.md)
- [Service Providers](../../core/service-providers.md)
- [Application Lifecycle](../application-lifecycle.md)
- [System Overview](system-overview.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
