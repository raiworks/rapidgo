package models

import (
	"strings"

	"github.com/RAiWorks/RGo/app/helpers"
	"gorm.io/gorm"
)

// User represents an application user.
type User struct {
	BaseModel
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Role     string `gorm:"size:50;default:user" json:"role"`
	Active   bool   `gorm:"default:true" json:"active"`
	Posts    []Post `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

// BeforeCreate hashes the password before inserting into the database.
// Skips if the password is already bcrypt-hashed (starts with "$2a$" or "$2b$").
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" && !strings.HasPrefix(u.Password, "$2a$") && !strings.HasPrefix(u.Password, "$2b$") {
		hashed, err := helpers.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashed
	}
	return nil
}
