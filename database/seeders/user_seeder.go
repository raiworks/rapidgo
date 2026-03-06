package seeders

import (
	"github.com/RAiWorks/RGo/database/models"
	"gorm.io/gorm"
)

func init() {
	Register(&UserSeeder{})
}

// UserSeeder creates default user accounts.
type UserSeeder struct{}

func (s *UserSeeder) Name() string { return "UserSeeder" }

func (s *UserSeeder) Seed(db *gorm.DB) error {
	users := []models.User{
		{Name: "Admin", Email: "admin@example.com", Password: "password123", Role: "admin"},
		{Name: "User", Email: "user@example.com", Password: "password123", Role: "user"},
	}
	for _, u := range users {
		// TODO: Hash password when Feature #19/#22 ships
		if err := db.FirstOrCreate(&u, models.User{Email: u.Email}).Error; err != nil {
			return err
		}
	}
	return nil
}
