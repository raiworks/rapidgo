---
title: "Data Flow Diagram"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Data Flow Diagram

## Abstract

This diagram traces data flow from the controller through the service
layer to the database and back, showing how data is transformed at
each stage for both API and SSR responses.

## Primary Data Flow

```mermaid
flowchart LR
    subgraph "HTTP Input"
        JSON["JSON Body<br/>c.ShouldBindJSON(&req)"]
        Form["Form Data<br/>c.PostForm('field')"]
        Query["Query Params<br/>c.DefaultQuery('page', '1')"]
        Param["URL Params<br/>c.Param('id')"]
    end

    subgraph "Controller"
        Parse["Parse & Validate"]
        Delegate["Call Service"]
        Format["Format Response"]
    end

    subgraph "Service"
        Logic["Business Logic"]
        Orchestrate["Orchestrate Models"]
    end

    subgraph "Model (GORM)"
        CRUD["Create / Read /<br/>Update / Delete"]
        Hooks["Hooks<br/>(BeforeCreate, etc.)"]
        Relations["Preload<br/>Relations"]
    end

    subgraph "Database"
        SQL["SQL Query"]
        Data[("Rows")]
    end

    subgraph "HTTP Output"
        JSONResp["JSON Response<br/>{success, data, meta}"]
        HTMLResp["HTML Response<br/>Rendered template"]
        Redirect["Redirect<br/>with flash message"]
    end

    JSON --> Parse
    Form --> Parse
    Query --> Parse
    Param --> Parse
    Parse --> Delegate
    Delegate --> Logic
    Logic --> Orchestrate
    Orchestrate --> CRUD
    CRUD --> Hooks
    CRUD --> Relations
    CRUD --> SQL
    SQL --> Data
    Data --> SQL
    SQL --> CRUD
    CRUD --> Orchestrate
    Orchestrate --> Logic
    Logic --> Delegate
    Delegate --> Format
    Format --> JSONResp
    Format --> HTMLResp
    Format --> Redirect
```

## API Create Flow (POST /api/users)

```mermaid
sequenceDiagram
    participant Client
    participant Controller as UserController
    participant Validator as Validator
    participant Service as UserService
    participant Model as User Model
    participant DB as Database

    Client->>Controller: POST /api/users<br/>{name, email, password}

    Controller->>Validator: ShouldBindJSON(&CreateUserRequest)
    alt Invalid
        Validator-->>Controller: binding error
        Controller-->>Client: 422 {error: "validation details"}
    end

    Controller->>Service: Create(name, email, password)

    Service->>Model: db.Where(email).First()
    Model->>DB: SELECT ... WHERE email = ?
    DB-->>Model: result
    alt Email exists
        Model-->>Service: found
        Service-->>Controller: error "email already exists"
        Controller-->>Client: 400 {error: "email already exists"}
    end

    Note over Model: BeforeCreate hook:<br/>HashPassword(password)

    Service->>Model: db.Create(&user)
    Model->>DB: INSERT INTO users ...
    DB-->>Model: created user with ID
    Model-->>Service: &user

    Service-->>Controller: &user, nil

    Controller-->>Client: 201 {success: true, data: user}
```

## SSR Form Flow (POST /users)

```mermaid
sequenceDiagram
    participant Browser
    participant Session as Session MW
    participant Controller as UserController
    participant Validator as Built-in Validator
    participant Service as UserService
    participant Flash as Flash Messages

    Browser->>Session: POST /users<br/>(with _csrf_token)
    Session->>Controller: Session data loaded

    Controller->>Validator: Required, MinLength, Email, Confirmed

    alt Invalid
        Controller->>Flash: FlashErrors(errors)
        Controller->>Flash: FlashOldInput({name, email})
        Controller-->>Browser: 302 Redirect → /users/create

        Note over Browser: Next GET /users/create
        Browser->>Session: GET /users/create
        Session->>Controller: Load session
        Controller->>Flash: GetFlash("_errors")
        Controller->>Flash: GetFlash("_old_input")
        Controller-->>Browser: HTML with errors + old values
    end

    Controller->>Service: Create(name, email, password)
    Service-->>Controller: &user, nil

    Controller->>Flash: Flash("success", "User created!")
    Controller-->>Browser: 302 Redirect → /users

    Note over Browser: Next GET /users
    Browser->>Session: GET /users
    Session->>Controller: Load session
    Controller->>Flash: GetFlash("success")
    Controller-->>Browser: HTML with success alert
```

## Pagination Data Flow

```mermaid
sequenceDiagram
    participant Client
    participant Controller
    participant Paginator as helpers.Paginate
    participant GORM
    participant DB

    Client->>Controller: GET /api/users?page=2&per_page=15

    Controller->>Paginator: Paginate(db.Model(&User{}), 2, 15, &users)

    Paginator->>GORM: db.Count(&total)
    GORM->>DB: SELECT COUNT(*) FROM users
    DB-->>GORM: total = 45

    Paginator->>GORM: db.Offset(15).Limit(15).Find(&users)
    GORM->>DB: SELECT * FROM users LIMIT 15 OFFSET 15
    DB-->>GORM: 15 rows

    Paginator-->>Controller: PaginateResult{Page:2, PerPage:15, Total:45, TotalPages:3}

    Controller-->>Client: 200 {success, data: [...], meta: {page:2, per_page:15, total:45, total_pages:3}}
```

## Caching Data Flow

```mermaid
flowchart TD
    A[Service needs data] --> B{Cache hit?}
    B -->|Yes| C[Return cached value]
    B -->|No| D[Query database via GORM]
    D --> E[Store result in cache]
    E --> F[Return data]

    subgraph "Cache Store"
        G["Redis<br/>(production)"]
        H["Memory<br/>(development)"]
    end

    B -.-> G
    B -.-> H
    E -.-> G
    E -.-> H
```

## Event Dispatch Flow

```mermaid
sequenceDiagram
    participant Service as UserService
    participant Dispatcher as EventDispatcher
    participant Mailer as Mail Listener
    participant Logger as Log Listener

    Service->>Dispatcher: Dispatch("user.created", &user)

    par Async (goroutines)
        Dispatcher->>Mailer: handler(user)
        Note over Mailer: Send welcome email

        Dispatcher->>Logger: handler(user)
        Note over Logger: slog.Info("new user", ...)
    end
```

## References

- [System Overview](system-overview.md)
- [Request Lifecycle](request-lifecycle.md)
- [Services Layer](../../infrastructure/services-layer.md)
- [Database](../../data/database.md)
- [Pagination](../../data/pagination.md)
- [Caching](../../infrastructure/caching.md)
- [Events](../../infrastructure/events.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
