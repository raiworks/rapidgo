package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type fileCacheEntry struct {
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}

// FileCache is a cache backend that stores each key as a JSON file on disk.
type FileCache struct {
	Path string
}

// NewFileCache returns a file-based cache store.
func NewFileCache(path string) *FileCache {
	return &FileCache{Path: path}
}

func (c *FileCache) keyPath(key string) string {
	return filepath.Join(c.Path, key+".cache")
}

func (c *FileCache) Get(key string) (string, error) {
	raw, err := os.ReadFile(c.keyPath(key))
	if err != nil {
		return "", nil // missing file = cache miss
	}
	var entry fileCacheEntry
	if err := json.Unmarshal(raw, &entry); err != nil {
		return "", nil
	}
	if time.Now().After(entry.ExpiresAt) {
		os.Remove(c.keyPath(key))
		return "", nil
	}
	return entry.Value, nil
}

func (c *FileCache) Set(key, value string, ttl time.Duration) error {
	if err := os.MkdirAll(c.Path, 0755); err != nil {
		return err
	}
	entry := fileCacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
	raw, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(c.keyPath(key), raw, 0600)
}

func (c *FileCache) Delete(key string) error {
	err := os.Remove(c.keyPath(key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (c *FileCache) Flush() error {
	entries, err := os.ReadDir(c.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, e := range entries {
		if !e.IsDir() {
			os.Remove(filepath.Join(c.Path, e.Name()))
		}
	}
	return nil
}
