package storage

import (
	"fmt"
	"io"
	"os"
)

// Driver defines the interface for file storage backends.
type Driver interface {
	Put(path string, content io.Reader) (string, error)
	Get(path string) (io.ReadCloser, error)
	Delete(path string) error
	URL(path string) string
}

// NewDriver creates a storage driver based on the STORAGE_DRIVER env var.
// Supported values: "local" (default), "s3".
func NewDriver() (Driver, error) {
	driver := os.Getenv("STORAGE_DRIVER")
	if driver == "" {
		driver = "local"
	}

	switch driver {
	case "local":
		basePath := os.Getenv("STORAGE_LOCAL_PATH")
		if basePath == "" {
			basePath = "storage/uploads"
		}
		return &LocalDriver{
			BasePath: basePath,
			BaseURL:  "/uploads",
		}, nil
	case "s3":
		return NewS3Driver()
	default:
		return nil, fmt.Errorf("unsupported storage driver: %s", driver)
	}
}
