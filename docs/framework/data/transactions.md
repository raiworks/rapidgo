---
title: "Transactions"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Transactions

## Abstract

This document covers GORM database transactions — wrapping multi-step
operations in atomic units with automatic commit and rollback.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Transaction Pattern](#2-transaction-pattern)
3. [Example: Credit Transfer](#3-example-credit-transfer)
4. [Commit and Rollback](#4-commit-and-rollback)
5. [Usage in Services](#5-usage-in-services)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Transaction** — A sequence of database operations executed as a
  single atomic unit. Either all succeed (commit) or all are undone
  (rollback).

## 2. Transaction Pattern

GORM provides a `Transaction` method that handles commit/rollback
automatically:

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // All operations inside use `tx`, not `db`
    if err := tx.Create(&record1).Error; err != nil {
        return err // rollback
    }
    if err := tx.Create(&record2).Error; err != nil {
        return err // rollback
    }
    return nil // commit
})
```

- **Return `nil`** → transaction commits.
- **Return an error** → transaction rolls back automatically.

## 3. Example: Credit Transfer

A multi-step operation that **MUST** be atomic:

```go
func TransferCredits(db *gorm.DB, fromID, toID uint, amount int) error {
    return db.Transaction(func(tx *gorm.DB) error {
        var from, to models.User
        if err := tx.First(&from, fromID).Error; err != nil {
            return err
        }
        if err := tx.First(&to, toID).Error; err != nil {
            return err
        }
        if err := tx.Model(&from).Update("credits",
            gorm.Expr("credits - ?", amount)).Error; err != nil {
            return err
        }
        if err := tx.Model(&to).Update("credits",
            gorm.Expr("credits + ?", amount)).Error; err != nil {
            return err
        }
        return nil // commit
    })
}
```

> **Important:** Use `gorm.Expr()` for atomic SQL operations like
> `credits - ?` to prevent race conditions.

## 4. Commit and Rollback

| Scenario | Result |
|----------|--------|
| Callback returns `nil` | All changes committed |
| Callback returns error | All changes rolled back |
| Callback panics | GORM catches the panic and rolls back |

## 5. Usage in Services

Transactions **SHOULD** be used in the service layer for operations
that modify multiple records:

```go
func (s *OrderService) PlaceOrder(userID uint, items []OrderItem) error {
    return s.DB.Transaction(func(tx *gorm.DB) error {
        // Create order
        order := models.Order{UserID: userID}
        if err := tx.Create(&order).Error; err != nil {
            return err
        }
        // Create order items
        for _, item := range items {
            item.OrderID = order.ID
            if err := tx.Create(&item).Error; err != nil {
                return err
            }
        }
        // Deduct inventory
        for _, item := range items {
            if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
                Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
```

## 6. Security Considerations

- Use `gorm.Expr()` for arithmetic operations to prevent race
  conditions in concurrent environments.
- Transactions **SHOULD** be kept short to minimize lock contention.
- Always use the `tx` handle inside the callback, never the outer
  `db` — otherwise operations run outside the transaction.

## 7. References

- [Database](database.md)
- [Models](models.md)
- [Services Layer](../infrastructure/services-layer.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
