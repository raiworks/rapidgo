package cache

import (
	"os"
	"testing"
	"time"
)

func TestNewRedisClient_DefaultDB(t *testing.T) {
	os.Unsetenv("REDIS_DB")
	client := NewRedisClient(nil)
	if client.Options().DB != 0 {
		t.Errorf("expected DB 0, got %d", client.Options().DB)
	}
}

func TestNewRedisClient_EnvDB(t *testing.T) {
	t.Setenv("REDIS_DB", "3")
	client := NewRedisClient(nil)
	if client.Options().DB != 3 {
		t.Errorf("expected DB 3, got %d", client.Options().DB)
	}
}

func TestNewRedisClient_DBOverride(t *testing.T) {
	t.Setenv("REDIS_DB", "3")
	db := 5
	client := NewRedisClient(&db)
	if client.Options().DB != 5 {
		t.Errorf("expected DB 5 (override), got %d", client.Options().DB)
	}
}

func TestNewRedisClient_InvalidDB(t *testing.T) {
	t.Setenv("REDIS_DB", "abc")
	client := NewRedisClient(nil)
	if client.Options().DB != 0 {
		t.Errorf("expected DB 0 (fallback), got %d", client.Options().DB)
	}
}

func TestNewRedisClient_NegativeDB(t *testing.T) {
	t.Setenv("REDIS_DB", "-1")
	client := NewRedisClient(nil)
	if client.Options().DB != 0 {
		t.Errorf("expected DB 0 (fallback for negative), got %d", client.Options().DB)
	}
}

func TestNewRedisClient_PoolSize(t *testing.T) {
	t.Setenv("REDIS_POOL_SIZE", "20")
	client := NewRedisClient(nil)
	if client.Options().PoolSize != 20 {
		t.Errorf("expected PoolSize 20, got %d", client.Options().PoolSize)
	}
}

func TestNewRedisClient_PoolSizeDefault(t *testing.T) {
	os.Unsetenv("REDIS_POOL_SIZE")
	client := NewRedisClient(nil)
	if client.Options().PoolSize != 10 {
		t.Errorf("expected default PoolSize 10, got %d", client.Options().PoolSize)
	}
}

func TestNewRedisClient_DialTimeout(t *testing.T) {
	t.Setenv("REDIS_DIAL_TIMEOUT", "10s")
	client := NewRedisClient(nil)
	if client.Options().DialTimeout != 10*time.Second {
		t.Errorf("expected DialTimeout 10s, got %v", client.Options().DialTimeout)
	}
}

func TestNewRedisClient_InvalidTimeout(t *testing.T) {
	t.Setenv("REDIS_DIAL_TIMEOUT", "invalid")
	client := NewRedisClient(nil)
	if client.Options().DialTimeout != 5*time.Second {
		t.Errorf("expected default DialTimeout 5s, got %v", client.Options().DialTimeout)
	}
}

func TestNewRedisClient_ReadWriteTimeout(t *testing.T) {
	t.Setenv("REDIS_READ_TIMEOUT", "7s")
	t.Setenv("REDIS_WRITE_TIMEOUT", "8s")
	client := NewRedisClient(nil)
	if client.Options().ReadTimeout != 7*time.Second {
		t.Errorf("expected ReadTimeout 7s, got %v", client.Options().ReadTimeout)
	}
	if client.Options().WriteTimeout != 8*time.Second {
		t.Errorf("expected WriteTimeout 8s, got %v", client.Options().WriteTimeout)
	}
}

func TestNewRedisClient_HostPort(t *testing.T) {
	t.Setenv("REDIS_HOST", "redis.example.com")
	t.Setenv("REDIS_PORT", "6380")
	client := NewRedisClient(nil)
	if client.Options().Addr != "redis.example.com:6380" {
		t.Errorf("expected addr redis.example.com:6380, got %s", client.Options().Addr)
	}
}

func TestNewRedisClient_DefaultHostPort(t *testing.T) {
	os.Unsetenv("REDIS_HOST")
	os.Unsetenv("REDIS_PORT")
	client := NewRedisClient(nil)
	if client.Options().Addr != "localhost:6379" {
		t.Errorf("expected addr localhost:6379, got %s", client.Options().Addr)
	}
}
