package providers

import (
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/database"
)

// DatabaseProvider registers the database connection in the service container.
type DatabaseProvider struct{}

// Register binds a *gorm.DB singleton. The connection is established lazily
// on first container.Make("db") call, not at registration time.
func (p *DatabaseProvider) Register(c *container.Container) {
	c.Singleton("db", func(c *container.Container) interface{} {
		db, err := database.Connect()
		if err != nil {
			panic("database connection failed: " + err.Error())
		}
		return db
	})
}

// Boot is a no-op. Migrations and seeding are future features.
func (p *DatabaseProvider) Boot(c *container.Container) {}
