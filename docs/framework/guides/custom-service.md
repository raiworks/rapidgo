---
title: "Creating a Custom Service & Provider"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Creating a Custom Service & Provider

## Abstract

A guide to creating a custom service, registering it through a
service provider, and resolving it in controllers.

## Table of Contents

1. [Overview](#1-overview)
2. [Generate Provider and Service](#2-generate-provider-and-service)
3. [Implement Service Logic](#3-implement-service-logic)
4. [Register in Provider](#4-register-in-provider)
5. [Register Provider in App](#5-register-provider-in-app)
6. [Resolve and Use in Controllers](#6-resolve-and-use-in-controllers)
7. [References](#7-references)

## 1. Overview

The framework uses a service container with providers for dependency
management:

1. **Generate** scaffolding with CLI
2. **Implement** service logic
3. **Register** service in a provider (as Singleton or Bind)
4. **Boot** the provider (optional setup, event listeners)
5. **Resolve** the service in controllers via the container

## 2. Generate Provider and Service

```bash
framework make:provider PaymentProvider
framework make:service PaymentGateway
```

This creates:
- `app/providers/paymentprovider.go`
- `app/services/paymentgateway.go`

## 3. Implement Service Logic

Edit `app/services/paymentgateway.go`:

```go
package services

import "gorm.io/gorm"

type PaymentGateway struct {
    DB     *gorm.DB
    APIKey string
}

func NewPaymentGateway(db *gorm.DB, apiKey string) *PaymentGateway {
    return &PaymentGateway{DB: db, APIKey: apiKey}
}

func (pg *PaymentGateway) Charge(userID uint, amount float64) error {
    // Payment processing logic
    return nil
}

func (pg *PaymentGateway) Refund(transactionID string) error {
    // Refund logic
    return nil
}
```

## 4. Register in Provider

Edit `app/providers/paymentprovider.go`:

```go
package providers

import (
    "os"

    "yourframework/app/services"
    "yourframework/core/container"
)

type PaymentProvider struct{}

func (p *PaymentProvider) Register(c *container.Container) {
    c.Singleton(func() *services.PaymentGateway {
        db := container.MustMake[*gorm.DB](c)
        apiKey := os.Getenv("PAYMENT_API_KEY")
        return services.NewPaymentGateway(db, apiKey)
    })
}

func (p *PaymentProvider) Boot(c *container.Container) {
    // Optional: listen for events, run setup
    dispatcher := container.MustMake[*events.Dispatcher](c)
    dispatcher.Listen("order.completed", func(payload interface{}) {
        // Post-payment processing
    })
}
```

Using `Singleton` ensures one shared instance across the application.

## 5. Register Provider in App

In `cmd/main.go` (or your app bootstrap):

```go
app.Register(&providers.PaymentProvider{})
```

The container calls `Register()` on all providers first, then
`Boot()` on all providers.

## 6. Resolve and Use in Controllers

```go
func (ctrl *OrderController) Checkout(c *gin.Context) {
    gateway := container.MustMake[*services.PaymentGateway](app.Container)

    err := gateway.Charge(userID, totalAmount)
    if err != nil {
        response.Error(c, http.StatusPaymentRequired, "Payment failed")
        return
    }
    response.Success(c, gin.H{"status": "paid"})
}
```

## 7. References

- [Service Container](../core/service-container.md)
- [Service Providers](../core/service-providers.md)
- [Services Layer](../infrastructure/services-layer.md)
- [Events](../infrastructure/events.md)
- [Code Generation](../cli/code-generation.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
