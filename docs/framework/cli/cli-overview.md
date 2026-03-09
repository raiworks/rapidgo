---
title: "CLI Overview"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# CLI Overview

## Abstract

This document describes the command-line interface built with Cobra,
listing every available command, the scaffolding system, database
management commands, and how to add custom commands.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Library](#2-library)
3. [Available Commands](#3-available-commands)
4. [Scaffolding Commands](#4-scaffolding-commands)
5. [Database Commands](#5-database-commands)
6. [Server Command](#6-server-command)
7. [Adding Custom Commands](#7-adding-custom-commands)
8. [Security Considerations](#8-security-considerations)
9. [References](#9-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Scaffolding** — Generating boilerplate code files from templates.
- **Command** — A CLI subcommand registered with Cobra.

## 2. Library

| Library | Import | Purpose |
|---------|--------|---------|
| Cobra | `github.com/spf13/cobra` | CLI command framework |
| Viper | `github.com/spf13/viper` | Configuration (pairs with Cobra) |

## 3. Available Commands

```text
framework new project          Create a new project
framework make:controller      Create a new controller
framework make:model           Create a new model
framework make:service         Create a new service
framework make:middleware      Create a new middleware
framework make:provider        Create a new service provider
framework make:request         Create a new request validation struct
framework migrate              Run database migrations
framework db:seed              Seed the database
framework serve                Start the HTTP server
```

## 4. Scaffolding Commands

Every `make:*` command generates a boilerplate `.go` file from a
template. See [Code Generation](code-generation.md) for full template
sources and the `generate()` implementation.

| Command | Output Path | Description |
|---------|------------|-------------|
| `make:controller PostController` | `http/controllers/postcontroller.go` | Controller with CRUD methods |
| `make:model User` | `database/models/user.go` | Model with `BaseModel` embed |
| `make:service OrderService` | `app/services/orderservice.go` | Service with DB injection |
| `make:provider PaymentProvider` | `app/providers/paymentprovider.go` | Provider with Register + Boot |
| `make:middleware RateLimit` | `http/middleware/ratelimit.go` | Middleware function |
| `make:request CreateOrderRequest` | `http/requests/createorderrequest.go` | Validation struct |

### Example

```bash
$ framework make:controller PostController
Controller created: http/controllers/postcontroller.go

$ framework make:service OrderService
Service created: app/services/orderservice.go
```

## 5. Database Commands

### `migrate`

Runs GORM `AutoMigrate` for all registered models:

```bash
$ framework migrate
```

See [Migrations](../data/migrations.md) for details on auto-migration
strategy and manual migration files.

### `db:seed`

Runs all seeder functions:

```bash
$ framework db:seed
```

See [Seeders](../data/seeders.md) for seeder implementation.

## 6. Server Command

### `serve`

Starts the HTTP server on the port from `.env` (`APP_PORT`, default
`8080`):

```bash
$ framework serve
Server starting on :8080
```

Equivalent to `go run cmd/main.go`.

## 7. Adding Custom Commands

Register new commands with Cobra in `cmd/`:

```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var customCmd = &cobra.Command{
    Use:   "greet [name]",
    Short: "Greet a user",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Hello, %s!\n", args[0])
    },
}

func init() {
    rootCmd.AddCommand(customCmd)
}
```

## 8. Security Considerations

- Generated files use templates with safe defaults — review
  generated code before committing.
- `db:seed` **SHOULD** only be used in development and staging
  environments.
- CLI commands that modify the database **SHOULD** require
  confirmation in production.

## 9. References

- [Cobra](https://github.com/spf13/cobra)
- [Code Generation](code-generation.md)
- [Migrations](../data/migrations.md)
- [Seeders](../data/seeders.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
