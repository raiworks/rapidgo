---
title: "Services Layer"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Services Layer

## Abstract

This document covers the services layer — business logic classes that
sit between controllers and models, with the controller → service →
model delegation pattern.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Purpose](#2-purpose)
3. [Service Location](#3-service-location)
4. [UserService Example](#4-userservice-example)
5. [Controller Integration](#5-controller-integration)
6. [Generating Services](#6-generating-services)
7. [Security Considerations](#7-security-considerations)
8. [References](#8-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Service** — A struct containing business logic that operates on
  models via GORM.

## 2. Purpose

Services separate **business logic** from **HTTP concerns**:

| Layer | Responsibility |
|-------|---------------|
| Controller | Parse request, call service, format response |
| Service | Business logic, validation, orchestration |
| Model | Data schema, relationships, hooks |

Benefits:
- Services are testable without HTTP context.
- Business rules live in one place, not scattered across controllers.
- Multiple controllers (web, API, CLI) can share the same service.

## 3. Service Location

Services live in `app/services/`.

## 4. UserService Example

```go
package services

import (
    "errors"

    "yourframework/database/models"

    "gorm.io/gorm"
)

type UserService struct {
    DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{DB: db}
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
    var user models.User
    if err := s.DB.First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *UserService) Create(name, email, password string) (*models.User, error) {
    existing := &models.User{}
    if err := s.DB.Where("email = ?", email).First(existing).Error; err == nil {
        return nil, errors.New("email already exists")
    }

    user := &models.User{
        Name:     name,
        Email:    email,
        Password: password, // hash before saving in real code
    }
    if err := s.DB.Create(user).Error; err != nil {
        return nil, err
    }
    return user, nil
}

func (s *UserService) Update(id uint, updates map[string]interface{}) (*models.User, error) {
    var user models.User
    if err := s.DB.First(&user, id).Error; err != nil {
        return nil, err
    }
    if err := s.DB.Model(&user).Updates(updates).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *UserService) Delete(id uint) error {
    return s.DB.Delete(&models.User{}, id).Error
}
```

## 5. Controller Integration

```go
func CreateUser(c *gin.Context) {
    var req requests.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        responses.Error(c, 422, err.Error())
        return
    }

    userSvc := services.NewUserService(db)
    user, err := userSvc.Create(req.Name, req.Email, req.Password)
    if err != nil {
        responses.Error(c, 400, err.Error())
        return
    }
    responses.Created(c, user)
}
```

## 6. Generating Services

Use the CLI scaffolding command:

```text
framework make:service OrderService
```

This generates `app/services/orderservice.go` from the service
template:

```go
package services

import "gorm.io/gorm"

type OrderService struct {
    DB *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
    return &OrderService{DB: db}
}

// Add service methods here
```

## 7. Security Considerations

- Services **SHOULD** validate business rules (e.g., unique email)
  before database operations.
- Passwords **MUST** be hashed in services or model hooks — never
  stored as plain text.
- Services **SHOULD** use parameterized queries (`Where("email = ?",
  email)`) to prevent SQL injection.

## 8. References

- [Controllers](../http/controllers.md)
- [Models](../data/models.md)
- [Transactions](../data/transactions.md)
- [Service Container](../core/service-container.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
