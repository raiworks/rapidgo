# 🏗️ Architecture: Services Layer

> **Feature**: `18` — Services Layer
> **Discussion**: [`18-services-layer-discussion.md`](18-services-layer-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

Feature #18 adds the services layer pattern with a `UserService` example in `app/services/`. The service encapsulates CRUD business logic for users, taking a `*gorm.DB` dependency and operating on the `models.User` model.

## File Structure

```
app/services/
├── user_service.go      # UserService struct, constructor, CRUD methods
└── user_service_test.go  # Tests with SQLite in-memory
```

### Files Created (1)
| File | Package | Lines (est.) |
|---|---|---|
| `app/services/user_service.go` | `services` | ~55 |

### Files Modified (0)
No existing files need modification.

---

## Component Design

### UserService (`app/services/user_service.go`)

**Responsibility**: Business logic for user CRUD operations.
**Package**: `services`

```go
package services

import (
	"errors"

	"github.com/RAiWorks/RGo/database/models"

	"gorm.io/gorm"
)

// UserService handles business logic for user operations.
type UserService struct {
	DB *gorm.DB
}

// NewUserService creates a new UserService with the given database connection.
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// GetByID retrieves a user by their ID.
func (s *UserService) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a new user after checking for duplicate email.
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

// Update updates a user's fields by ID.
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

// Delete deletes a user by ID.
func (s *UserService) Delete(id uint) error {
	return s.DB.Delete(&models.User{}, id).Error
}
```

**Design notes**:
- Exact match to blueprint code (lines 2535–2595)
- Module path `github.com/RAiWorks/RGo/database/models` instead of blueprint's `yourframework/database/models`
- `Create` checks for duplicate email before insertion (business rule)
- `Password` stored as-is with comment — hashing belongs to Feature #22
- `Update` uses `map[string]interface{}` for flexible partial updates
- `Delete` uses GORM's `Delete` with model type and ID

---

## Data Flow

### Service usage pattern
```
Controller → services.NewUserService(db)
          → svc.Create(name, email, password)
          → GORM: check duplicate email → create user
          → return (*models.User, error)
```

### GetByID
```
svc.GetByID(42) → db.First(&user, 42) → *User or error
```

### Update
```
svc.Update(42, map{"name": "New"}) → db.First → db.Model.Updates → *User
```

### Delete
```
svc.Delete(42) → db.Delete(&User{}, 42) → error
```

---

## Dependencies

| Dependency | Type | Usage |
|---|---|---|
| `gorm.io/gorm` | external | Database operations |
| `database/models` | internal | `User` model |
| `errors` | stdlib | `errors.New` for business rule violations |

---

## Impact on Existing Code

| Component | Impact |
|---|---|
| `app/services/` | `.gitkeep` stays; `user_service.go` added |
| No other files modified | Services are standalone — no wiring needed |
