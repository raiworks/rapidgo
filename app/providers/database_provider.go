package providers

import (
	"github.com/RAiWorks/RapidGo/core/config"
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/database"
	"gorm.io/gorm"
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

	c.Singleton("db.resolver", func(c *container.Container) interface{} {
		writer := c.Make("db").(*gorm.DB)
		if config.Env("DB_READ_HOST", "") == "" {
			return database.NewResolver(writer, writer)
		}
		reader, err := database.ConnectWithConfig(database.NewReadDBConfig())
		if err != nil {
			panic("read database connection failed: " + err.Error())
		}
		return database.NewResolver(writer, reader)
	})
}

// Boot is a no-op. Migrations and seeding are future features.
func (p *DatabaseProvider) Boot(c *container.Container) {}
