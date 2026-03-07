package providers

import (
	"github.com/RAiWorks/RapidGo/core/config"
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/queue"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// QueueProvider registers the queue Dispatcher in the container.
type QueueProvider struct{}

// Register binds a *queue.Dispatcher singleton. Driver is selected via QUEUE_DRIVER env var.
func (p *QueueProvider) Register(c *container.Container) {
	c.Singleton("queue", func(c *container.Container) interface{} {
		driver := config.Env("QUEUE_DRIVER", "database")
		switch driver {
		case "database":
			db := container.MustMake[*gorm.DB](c, "db")
			table := config.Env("QUEUE_TABLE", "jobs")
			failedTable := config.Env("QUEUE_FAILED_TABLE", "failed_jobs")
			return queue.NewDispatcher(queue.NewDatabaseDriver(db, table, failedTable))
		case "redis":
			client := container.MustMake[*redis.Client](c, "redis")
			return queue.NewDispatcher(queue.NewRedisDriver(client))
		case "memory":
			return queue.NewDispatcher(queue.NewMemoryDriver())
		case "sync":
			return queue.NewDispatcher(queue.NewSyncDriver())
		default:
			panic("queue: unsupported driver: " + driver)
		}
	})
}

// Boot is a no-op.
func (p *QueueProvider) Boot(c *container.Container) {}
