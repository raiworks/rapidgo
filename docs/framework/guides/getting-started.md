---
title: "Getting Started"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Getting Started

## Abstract

A step-by-step guide to creating your first project, configuring it,
running migrations, and serving your first page.

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [Create a New Project](#2-create-a-new-project)
3. [Configure Environment](#3-configure-environment)
4. [Run Migrations](#4-run-migrations)
5. [Start the Server](#5-start-the-server)
6. [Visit Your First Page](#6-visit-your-first-page)
7. [Next Steps](#7-next-steps)
8. [References](#8-references)

## 1. Prerequisites

| Requirement | Version | Notes |
|-------------|---------|-------|
| Go | 1.21+ | [golang.org/dl](https://golang.org/dl/) |
| Database | — | PostgreSQL, MySQL, or SQLite |
| Git | any | For version control |

Optional:

- **Redis** — for session/cache drivers
- **Docker** — for containerized deployment

## 2. Create a New Project

```bash
framework new myapp
cd myapp
```

Or manually:

```bash
mkdir myapp && cd myapp
go mod init myapp
```

Install dependencies:

```bash
go mod tidy
```

## 3. Configure Environment

Copy the example `.env` file and edit it:

```bash
cp .env.example .env
```

Minimum configuration:

```env
APP_NAME=MyApp
APP_PORT=8080
APP_ENV=development
APP_DEBUG=true
APP_KEY=base64:your-32-byte-random-key-here

DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=secret
DB_NAME=myapp

JWT_SECRET=change-me-in-production
```

See [Environment Variables Reference](../reference/env-reference.md)
for all available variables.

## 4. Run Migrations

Create your database tables:

```bash
framework migrate
```

This runs GORM `AutoMigrate` for all registered models. See
[Migrations](../data/migrations.md) for details.

## 5. Start the Server

```bash
framework serve
```

Or:

```bash
go run cmd/main.go
```

Output:

```text
Server starting on :8080
```

## 6. Visit Your First Page

Open your browser to:

```text
http://localhost:8080/health
```

You should see:

```json
{"status": "ok"}
```

## 7. Next Steps

| Goal | Guide |
|------|-------|
| Build a CRUD app | [Creating a CRUD Application](creating-crud.md) |
| Build a REST API | [Building a REST API](building-api.md) |
| Add SSR forms | [SSR Forms with Validation](ssr-forms.md) |
| Create custom services | [Custom Service & Provider](custom-service.md) |

## 8. References

- [Project Structure](../architecture/project-structure.md)
- [Configuration](../core/configuration.md)
- [Application Lifecycle](../architecture/application-lifecycle.md)
- [CLI Overview](../cli/cli-overview.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
