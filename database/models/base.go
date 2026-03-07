package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel provides common fields for all models.
// Embed this in your model structs to get ID, CreatedAt,
// UpdatedAt, and soft delete support via DeletedAt.
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
