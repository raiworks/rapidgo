---
title: "Models"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Models

## Abstract

This document covers GORM model definitions, the `BaseModel` struct,
field tags, relationships, eager loading, and model hooks.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Model Location](#2-model-location)
3. [BaseModel](#3-basemodel)
4. [Example Models](#4-example-models)
5. [Field Tags](#5-field-tags)
6. [Relationships](#6-relationships)
7. [Eager Loading](#7-eager-loading)
8. [Model Hooks](#8-model-hooks)
9. [Security Considerations](#9-security-considerations)
10. [References](#10-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Model** — A Go struct that maps to a database table via GORM.
- **Hook** — A method on a model that GORM calls automatically at
  specific lifecycle points (before create, after update, etc.).

## 2. Model Location

Models live in `database/models/`.

## 3. BaseModel

All models **SHOULD** embed `BaseModel` for consistent primary key
and timestamp fields:

```go
package models

import "time"

type BaseModel struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

GORM automatically manages `CreatedAt` and `UpdatedAt`.

## 4. Example Models

### User

```go
type User struct {
    BaseModel
    Name     string `gorm:"size:100;not null" json:"name"`
    Email    string `gorm:"size:255;uniqueIndex;not null" json:"email"`
    Password string `gorm:"size:255;not null" json:"-"`
    Role     string `gorm:"size:50;default:user" json:"role"`
    Active   bool   `gorm:"default:true" json:"active"`
    Posts    []Post  `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}
```

### Post

```go
type Post struct {
    BaseModel
    Title   string `gorm:"size:255;not null" json:"title"`
    Slug    string `gorm:"size:255;uniqueIndex" json:"slug"`
    Body    string `gorm:"type:text" json:"body"`
    UserID  uint   `gorm:"index;not null" json:"user_id"`
    User    User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

## 5. Field Tags

### GORM Tags

| Tag | Description |
|-----|-------------|
| `gorm:"primaryKey"` | Primary key |
| `gorm:"size:255"` | Column size |
| `gorm:"not null"` | NOT NULL constraint |
| `gorm:"uniqueIndex"` | Unique index |
| `gorm:"index"` | Standard index |
| `gorm:"default:value"` | Default value |
| `gorm:"type:text"` | Column type override |
| `gorm:"foreignKey:Field"` | Foreign key relationship |

### JSON Tags

| Tag | Description |
|-----|-------------|
| `json:"field_name"` | JSON key name |
| `json:"-"` | Exclude from JSON output |
| `json:"field,omitempty"` | Omit if empty/zero |

> **Important:** Use `json:"-"` on sensitive fields like `Password`
> to prevent accidental exposure in API responses.

## 6. Relationships

GORM supports standard relationship types:

| Type | Example |
|------|---------|
| `HasOne` | User has one Profile |
| `HasMany` | User has many Posts |
| `BelongsTo` | Post belongs to User |
| `Many2Many` | Post has many Tags, Tag has many Posts |

Define with struct fields and foreign key tags:

```go
// User HasMany Posts
type User struct {
    BaseModel
    Posts []Post `gorm:"foreignKey:UserID"`
}

// Post BelongsTo User
type Post struct {
    BaseModel
    UserID uint `gorm:"index;not null"`
    User   User `gorm:"foreignKey:UserID"`
}
```

## 7. Eager Loading

Use `Preload` to avoid N+1 queries:

```go
// Load user with their posts
var user User
db.Preload("Posts").First(&user, id)

// Nested preload
db.Preload("Posts.Comments").Find(&users)

// Conditional preload
db.Preload("Posts", "published = ?", true).Find(&users)
```

## 8. Model Hooks

Hooks execute automatically during GORM operations:

```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    hashed, err := helpers.HashPassword(u.Password)
    if err != nil {
        return err
    }
    u.Password = hashed
    return nil
}
```

Available hooks:

| Hook | When |
|------|------|
| `BeforeCreate` | Before INSERT |
| `AfterCreate` | After INSERT |
| `BeforeUpdate` | Before UPDATE |
| `AfterUpdate` | After UPDATE |
| `BeforeDelete` | Before DELETE |
| `AfterDelete` | After DELETE |
| `BeforeSave` | Before INSERT or UPDATE |
| `AfterSave` | After INSERT or UPDATE |

## 9. Security Considerations

- Passwords **MUST** be hashed before storage (use `BeforeCreate`
  hook).
- Sensitive fields **MUST** use `json:"-"` to prevent API exposure.
- User-supplied values used in `Where` clauses **MUST** use
  parameterized queries (`Where("email = ?", email)`), not string
  concatenation.

## 10. References

- [Database](database.md)
- [Transactions](transactions.md)
- [Seeders](seeders.md)
- [Services Layer](../infrastructure/services-layer.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
