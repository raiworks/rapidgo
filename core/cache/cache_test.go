package cache

import (
	"sync"
	"testing"
	"time"
)

// TC-01: Set and Get returns value.
func TestSetAndGet(t *testing.T) {
	c := NewMemoryCache()
	if err := c.Set("greeting", "hello", time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}
	val, err := c.Get("greeting")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "hello" {
		t.Fatalf("Get = %q, want %q", val, "hello")
	}
}

// TC-02: Get missing key returns empty string.
func TestGetMissing(t *testing.T) {
	c := NewMemoryCache()
	val, err := c.Get("nope")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "" {
		t.Fatalf("Get = %q, want empty", val)
	}
}

// TC-03: Get expired key returns empty string.
func TestGetExpired(t *testing.T) {
	c := NewMemoryCache()
	_ = c.Set("temp", "data", time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	val, err := c.Get("temp")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "" {
		t.Fatalf("Get = %q, want empty (expired)", val)
	}

	// Verify lazy delete removed the entry.
	c.mu.RLock()
	_, exists := c.items["temp"]
	c.mu.RUnlock()
	if exists {
		t.Fatal("expired entry not cleaned up")
	}
}

// TC-04: Delete removes key.
func TestDelete(t *testing.T) {
	c := NewMemoryCache()
	_ = c.Set("k", "v", time.Minute)
	if err := c.Delete("k"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	val, _ := c.Get("k")
	if val != "" {
		t.Fatalf("Get after Delete = %q, want empty", val)
	}
}

// TC-05: Flush clears all keys.
func TestFlush(t *testing.T) {
	c := NewMemoryCache()
	_ = c.Set("a", "1", time.Minute)
	_ = c.Set("b", "2", time.Minute)
	_ = c.Set("c", "3", time.Minute)
	if err := c.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}
	for _, key := range []string{"a", "b", "c"} {
		val, _ := c.Get(key)
		if val != "" {
			t.Fatalf("Get(%q) after Flush = %q, want empty", key, val)
		}
	}
}

// TC-06: Concurrent access is safe.
func TestConcurrentAccess(t *testing.T) {
	c := NewMemoryCache()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(3)
		key := "key"
		go func() {
			defer wg.Done()
			_ = c.Set(key, "v", time.Minute)
		}()
		go func() {
			defer wg.Done()
			_, _ = c.Get(key)
		}()
		go func() {
			defer wg.Done()
			_ = c.Delete(key)
		}()
	}
	wg.Wait()
}

// TC-07: NewStore "memory" returns a valid Store.
func TestNewStoreMemory(t *testing.T) {
	store, err := NewStore("memory", "")
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if store == nil {
		t.Fatal("NewStore returned nil")
	}
}

// TC-08: NewStore unknown driver returns error.
func TestNewStoreUnknown(t *testing.T) {
	store, err := NewStore("redis", "")
	if err == nil {
		t.Fatal("NewStore with unknown driver should return error")
	}
	if store != nil {
		t.Fatal("NewStore should return nil on error")
	}
}

// TC-09: Prefix applied to keys.
func TestPrefixStore(t *testing.T) {
	store, err := NewStore("memory", "app:")
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	_ = store.Set("k", "hello", time.Minute)
	val, _ := store.Get("k")
	if val != "hello" {
		t.Fatalf("Get = %q, want %q", val, "hello")
	}

	// Verify the underlying key is prefixed by accessing the inner store.
	ps := store.(*prefixStore)
	inner := ps.store.(*MemoryCache)
	inner.mu.RLock()
	_, plain := inner.items["k"]
	_, prefixed := inner.items["app:k"]
	inner.mu.RUnlock()
	if plain {
		t.Fatal("key stored without prefix")
	}
	if !prefixed {
		t.Fatal("key not stored with prefix")
	}
}

// TC-10: Set overwrites existing key.
func TestSetOverwrite(t *testing.T) {
	c := NewMemoryCache()
	_ = c.Set("k", "old", time.Minute)
	_ = c.Set("k", "new", time.Minute)
	val, _ := c.Get("k")
	if val != "new" {
		t.Fatalf("Get = %q, want %q", val, "new")
	}
}
