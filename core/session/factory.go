package session

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewStore resolves the correct session backend from SESSION_DRIVER.
func NewStore(db *gorm.DB) (Store, error) {
	driver := os.Getenv("SESSION_DRIVER")

	switch driver {
	case "db":
		db.AutoMigrate(&SessionRecord{})
		return &DBStore{DB: db}, nil
	case "file":
		path := os.Getenv("SESSION_FILE_PATH")
		if path == "" {
			path = "storage/sessions"
		}
		return &FileStore{Path: path}, nil
	case "memory", "":
		return NewMemoryStore(), nil
	case "cookie":
		key := []byte(os.Getenv("APP_KEY"))
		store, err := NewCookieStore(key)
		if err != nil {
			return nil, fmt.Errorf("cookie session store: %w", err)
		}
		return store, nil
	case "redis":
		host := os.Getenv("REDIS_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("REDIS_PORT")
		if port == "" {
			port = "6379"
		}
		password := os.Getenv("REDIS_PASSWORD")
		client := redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: password,
		})
		prefix := os.Getenv("SESSION_PREFIX")
		if prefix == "" {
			prefix = "session:"
		}
		return NewRedisStore(client, prefix), nil
	default:
		return nil, fmt.Errorf("unsupported SESSION_DRIVER: %s", driver)
	}
}
