package cache

import (
	"fmt"
	"sync"
	"time"
)

// Store defines the cache contract.
type Store interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
	Delete(key string) error
	Flush() error
}

type memCacheEntry struct {
	Value     string
	ExpiresAt time.Time
}

// MemoryCache is an in-process cache backed by a map with TTL-based expiry.
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]memCacheEntry
}

// NewMemoryCache returns an empty in-memory store.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{items: make(map[string]memCacheEntry)}
}

func (c *MemoryCache) Get(key string) (string, error) {
	c.mu.RLock()
	entry, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return "", nil
	}
	if time.Now().After(entry.ExpiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return "", nil
	}
	return entry.Value, nil
}

func (c *MemoryCache) Set(key, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = memCacheEntry{Value: value, ExpiresAt: time.Now().Add(ttl)}
	return nil
}

func (c *MemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *MemoryCache) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]memCacheEntry)
	return nil
}

// prefixStore wraps a Store and prepends a prefix to all keys.
type prefixStore struct {
	store  Store
	prefix string
}

func (p *prefixStore) Get(key string) (string, error) {
	return p.store.Get(p.prefix + key)
}

func (p *prefixStore) Set(key, value string, ttl time.Duration) error {
	return p.store.Set(p.prefix+key, value, ttl)
}

func (p *prefixStore) Delete(key string) error {
	return p.store.Delete(p.prefix + key)
}

func (p *prefixStore) Flush() error {
	return p.store.Flush()
}

// NewStore creates a cache Store for the given driver.
// If prefix is non-empty, all keys are automatically prefixed.
func NewStore(driver, prefix string) (Store, error) {
	var store Store
	switch driver {
	case "memory":
		store = NewMemoryCache()
	default:
		return nil, fmt.Errorf("cache: unsupported driver %q", driver)
	}
	if prefix != "" {
		store = &prefixStore{store: store, prefix: prefix}
	}
	return store, nil
}
