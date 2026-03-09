---
title: "Code Generation"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Code Generation

## Abstract

This document details the `make:*` scaffolding commands — the shared
`generate()` function, each command definition, and the full template
source for every generated artifact.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Generator Function](#2-generator-function)
3. [Command Definitions](#3-command-definitions)
4. [Templates](#4-templates)
5. [Output Examples](#5-output-examples)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Template** — A Go `text/template` string used to produce
  generated source files.

## 2. Generator Function

All `make:*` commands delegate to a shared `generate()` function:

```go
package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "text/template"
)

func generate(kind, name, tpl, dir string) {
    filename := strings.ToLower(name) + ".go"
    path := filepath.Join(dir, filename)

    if err := os.MkdirAll(dir, 0755); err != nil {
        fmt.Printf("Error creating directory: %v\n", err)
        return
    }

    f, err := os.Create(path)
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        return
    }
    defer f.Close()

    t := template.Must(template.New(kind).Parse(tpl))
    t.Execute(f, map[string]string{"Name": name})

    fmt.Printf("%s created: %s\n", strings.Title(kind), path)
}
```

### Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `kind` | `string` | Artifact type (used as template name and label) |
| `name` | `string` | User-supplied name, passed as `{{.Name}}` |
| `tpl` | `string` | Go template string |
| `dir` | `string` | Target directory (created if absent) |

### Naming Convention

The output filename is always `strings.ToLower(name) + ".go"`.

| Input | Output Path |
|-------|-------------|
| `PostController` | `http/controllers/postcontroller.go` |
| `User` | `database/models/user.go` |
| `OrderService` | `app/services/orderservice.go` |
| `PaymentProvider` | `app/providers/paymentprovider.go` |

## 3. Command Definitions

Each command is a Cobra command that calls `generate()` with its kind,
template, and output directory.

### `make:controller`

```go
var makeControllerCmd = &cobra.Command{
    Use:   "make:controller [name]",
    Short: "Create a new controller",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        generate("controller", name, controllerTpl, "http/controllers")
    },
}
```

### `make:model`

```go
var makeModelCmd = &cobra.Command{
    Use:   "make:model [name]",
    Short: "Create a new model",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        generate("model", name, modelTpl, "database/models")
    },
}
```

### `make:service`

```go
var makeServiceCmd = &cobra.Command{
    Use:   "make:service [name]",
    Short: "Create a new service",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        generate("service", name, serviceTpl, "app/services")
    },
}
```

### `make:provider`

```go
var makeProviderCmd = &cobra.Command{
    Use:   "make:provider [name]",
    Short: "Create a new service provider",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        generate("provider", name, providerTpl, "app/providers")
    },
}
```

## 4. Templates

### Controller Template

```go
var controllerTpl = `package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type {{.Name}} struct{}

func (ctrl *{{.Name}}) Index(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "{{.Name}} index"})
}

func (ctrl *{{.Name}}) Show(c *gin.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, gin.H{"id": id})
}

func (ctrl *{{.Name}}) Store(c *gin.Context) {
    c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (ctrl *{{.Name}}) Update(c *gin.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, gin.H{"id": id, "message": "updated"})
}

func (ctrl *{{.Name}}) Destroy(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
`
```

Generated controllers implement all five CRUD methods matching the
`ResourceController` interface. Developers **SHOULD** inject a
service and replace the stub responses.

### Model Template

```go
var modelTpl = `package models

type {{.Name}} struct {
    BaseModel
    // Add fields here
}
`
```

Models embed `BaseModel` which provides `ID`, `CreatedAt`, and
`UpdatedAt` fields.

### Service Template

```go
var serviceTpl = `package services

import "gorm.io/gorm"

type {{.Name}} struct {
    DB *gorm.DB
}

func New{{.Name}}(db *gorm.DB) *{{.Name}} {
    return &{{.Name}}{DB: db}
}

// Add service methods here
`
```

Services receive `*gorm.DB` via constructor injection.

### Provider Template

```go
var providerTpl = `package providers

import "yourframework/core/container"

type {{.Name}} struct{}

func (p *{{.Name}}) Register(c *container.Container) {
    // Bind services into the container
}

func (p *{{.Name}}) Boot(c *container.Container) {
    // Run after all providers are registered
}
`
```

## 5. Output Examples

```bash
$ framework make:controller PostController
Controller created: http/controllers/postcontroller.go

$ framework make:model Post
Model created: database/models/post.go

$ framework make:service OrderService
Service created: app/services/orderservice.go

$ framework make:provider PaymentProvider
Provider created: app/providers/paymentprovider.go
```

## 6. Security Considerations

- Generated files are written using `os.Create` which truncates
  existing files — the command **SHOULD** check for existing files
  and warn before overwriting.
- Template content is trusted (compiled into the binary) — no
  injection risk.
- Developers **MUST** review generated code and add proper
  authorization before deploying.

## 7. References

- [CLI Overview](cli-overview.md)
- [Controllers](../http/controllers.md)
- [Models](../data/models.md)
- [Services Layer](../infrastructure/services-layer.md)
- [Service Providers](../core/service-providers.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
