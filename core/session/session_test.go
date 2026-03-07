package session

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// --- MemoryStore ---

func TestMemoryStore_WriteRead(t *testing.T) {
	s := NewMemoryStore()
	data := map[string]interface{}{"user": "alice"}
	if err := s.Write("s1", data, 10*time.Minute); err != nil {
		t.Fatal(err)
	}
	got, err := s.Read("s1")
	if err != nil {
		t.Fatal(err)
	}
	if got["user"] != "alice" {
		t.Fatalf("expected 'alice', got %v", got["user"])
	}
}

func TestMemoryStore_ReadMissing(t *testing.T) {
	s := NewMemoryStore()
	got, err := s.Read("nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestMemoryStore_Destroy(t *testing.T) {
	s := NewMemoryStore()
	s.Write("s1", map[string]interface{}{"k": "v"}, 10*time.Minute)
	s.Destroy("s1")
	got, _ := s.Read("s1")
	if len(got) != 0 {
		t.Fatal("expected empty map after destroy")
	}
}

func TestMemoryStore_GC(t *testing.T) {
	s := NewMemoryStore()
	// Write with very short lifetime
	s.Write("expired", map[string]interface{}{"k": "v"}, 1*time.Millisecond)
	s.Write("valid", map[string]interface{}{"k": "v"}, 10*time.Minute)
	time.Sleep(5 * time.Millisecond)
	s.GC(1 * time.Millisecond)

	got, _ := s.Read("expired")
	if len(got) != 0 {
		t.Fatal("expected expired session to be removed")
	}
	got, _ = s.Read("valid")
	if len(got) == 0 {
		t.Fatal("expected valid session to remain")
	}
}

// --- FileStore ---

func TestFileStore_WriteRead(t *testing.T) {
	dir := t.TempDir()
	s := &FileStore{Path: dir}
	data := map[string]interface{}{"user": "bob"}
	if err := s.Write("s1", data, 10*time.Minute); err != nil {
		t.Fatal(err)
	}
	got, err := s.Read("s1")
	if err != nil {
		t.Fatal(err)
	}
	if got["user"] != "bob" {
		t.Fatalf("expected 'bob', got %v", got["user"])
	}
}

func TestFileStore_ReadMissing(t *testing.T) {
	dir := t.TempDir()
	s := &FileStore{Path: dir}
	got, err := s.Read("nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestFileStore_Destroy(t *testing.T) {
	dir := t.TempDir()
	s := &FileStore{Path: dir}
	s.Write("s1", map[string]interface{}{"k": "v"}, 10*time.Minute)
	s.Destroy("s1")
	got, _ := s.Read("s1")
	if len(got) != 0 {
		t.Fatal("expected empty map after destroy")
	}
}

// --- CookieStore ---

func TestCookieStore_WriteRead(t *testing.T) {
	key := []byte("01234567890123456789012345678901") // 32 bytes
	s, err := NewCookieStore(key)
	if err != nil {
		t.Fatal(err)
	}
	data := map[string]interface{}{"user": "carol"}
	if err := s.Write("s1", data, 10*time.Minute); err != nil {
		t.Fatal(err)
	}
	got, err := s.Read("s1")
	if err != nil {
		t.Fatal(err)
	}
	if got["user"] != "carol" {
		t.Fatalf("expected 'carol', got %v", got["user"])
	}
}

func TestCookieStore_ReadMissing(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	s, _ := NewCookieStore(key)
	got, err := s.Read("nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestCookieStore_InvalidKey(t *testing.T) {
	_, err := NewCookieStore([]byte("short"))
	if err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestCookieStore_Destroy(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	s, _ := NewCookieStore(key)
	s.Write("s1", map[string]interface{}{"k": "v"}, 10*time.Minute)
	s.Destroy("s1")
	got, _ := s.Read("s1")
	if len(got) != 0 {
		t.Fatal("expected empty map after destroy")
	}
}

// --- Manager ---

func newTestManager() *Manager {
	store := NewMemoryStore()
	return &Manager{
		Store:      store,
		CookieName: "test_session",
		Lifetime:   120 * time.Minute,
		Path:       "/",
		HTTPOnly:   true,
		SameSite:   http.SameSiteLaxMode,
	}
}

func TestManager_StartNewSession(t *testing.T) {
	mgr := newTestManager()
	req := httptest.NewRequest("GET", "/", nil)
	id, data, err := mgr.Start(req)
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("expected non-empty session ID")
	}
	if len(data) != 0 {
		t.Fatal("expected empty data for new session")
	}
}

func TestManager_StartExistingSession(t *testing.T) {
	mgr := newTestManager()
	// Pre-populate data
	mgr.Store.Write("existing-id", map[string]interface{}{"user": "dave"}, 120*time.Minute)
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "test_session", Value: "existing-id"})

	id, data, err := mgr.Start(req)
	if err != nil {
		t.Fatal(err)
	}
	if id != "existing-id" {
		t.Fatalf("expected 'existing-id', got %q", id)
	}
	if data["user"] != "dave" {
		t.Fatalf("expected 'dave', got %v", data["user"])
	}
}

func TestManager_Save(t *testing.T) {
	mgr := newTestManager()
	w := httptest.NewRecorder()
	data := map[string]interface{}{"user": "eve"}
	if err := mgr.Save(w, "save-id", data); err != nil {
		t.Fatal(err)
	}
	// Check cookie was set
	cookies := w.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "test_session" && c.Value == "save-id" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected session cookie to be set")
	}
	// Check data was persisted
	got, _ := mgr.Store.Read("save-id")
	if got["user"] != "eve" {
		t.Fatalf("expected 'eve', got %v", got["user"])
	}
}

func TestManager_Destroy(t *testing.T) {
	mgr := newTestManager()
	mgr.Store.Write("del-id", map[string]interface{}{"k": "v"}, 120*time.Minute)
	w := httptest.NewRecorder()
	if err := mgr.Destroy(w, "del-id"); err != nil {
		t.Fatal(err)
	}
	// Session data should be gone
	got, _ := mgr.Store.Read("del-id")
	if len(got) != 0 {
		t.Fatal("expected session data to be removed")
	}
	// Cookie should be cleared
	cookies := w.Result().Cookies()
	for _, c := range cookies {
		if c.Name == "test_session" && c.MaxAge != -1 {
			t.Fatal("expected cookie MaxAge=-1")
		}
	}
}

// --- Flash Messages ---

func TestManager_Flash(t *testing.T) {
	mgr := newTestManager()
	data := make(map[string]interface{})
	mgr.Flash(data, "success", "Item created!")
	flashes, _ := data["_flashes"].(map[string]interface{})
	if flashes["success"] != "Item created!" {
		t.Fatalf("expected flash 'Item created!', got %v", flashes["success"])
	}
}

func TestManager_GetFlash(t *testing.T) {
	mgr := newTestManager()
	data := make(map[string]interface{})
	mgr.Flash(data, "success", "Done!")
	val, ok := mgr.GetFlash(data, "success")
	if !ok || val != "Done!" {
		t.Fatalf("expected 'Done!', got %v", val)
	}
	// Should be removed after reading
	_, ok = mgr.GetFlash(data, "success")
	if ok {
		t.Fatal("expected flash to be removed after reading")
	}
}

func TestManager_GetFlash_Missing(t *testing.T) {
	mgr := newTestManager()
	data := make(map[string]interface{})
	val, ok := mgr.GetFlash(data, "nonexistent")
	if ok || val != nil {
		t.Fatal("expected nil, false for missing flash")
	}
}

func TestManager_FlashErrors(t *testing.T) {
	mgr := newTestManager()
	data := make(map[string]interface{})
	errs := map[string][]string{"email": {"required"}}
	mgr.FlashErrors(data, errs)
	flashes, _ := data["_flashes"].(map[string]interface{})
	if flashes["_errors"] == nil {
		t.Fatal("expected _errors flash to be set")
	}
}

func TestManager_FlashOldInput(t *testing.T) {
	mgr := newTestManager()
	data := make(map[string]interface{})
	input := map[string]string{"name": "Alice"}
	mgr.FlashOldInput(data, input)
	flashes, _ := data["_flashes"].(map[string]interface{})
	if flashes["_old_input"] == nil {
		t.Fatal("expected _old_input flash to be set")
	}
}

// --- Factory ---

func TestNewStore_Memory(t *testing.T) {
	t.Setenv("SESSION_DRIVER", "memory")
	store, err := NewStore(nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := store.(*MemoryStore); !ok {
		t.Fatal("expected *MemoryStore")
	}
}

func TestNewStore_File(t *testing.T) {
	t.Setenv("SESSION_DRIVER", "file")
	t.Setenv("SESSION_FILE_PATH", t.TempDir())
	store, err := NewStore(nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := store.(*FileStore); !ok {
		t.Fatal("expected *FileStore")
	}
}

func TestNewStore_Unsupported(t *testing.T) {
	t.Setenv("SESSION_DRIVER", "invalid")
	_, err := NewStore(nil)
	if err == nil {
		t.Fatal("expected error for unsupported driver")
	}
}

func TestNewStore_Redis(t *testing.T) {
	s := miniredis.RunT(t)
	t.Setenv("SESSION_DRIVER", "redis")
	t.Setenv("REDIS_HOST", s.Host())
	t.Setenv("REDIS_PORT", s.Port())
	store, err := NewStore(nil)
	if err != nil {
		t.Fatalf("expected redis store, got error: %v", err)
	}
	if _, ok := store.(*RedisStore); !ok {
		t.Fatalf("expected *RedisStore, got %T", store)
	}
}

func TestNewStore_Default(t *testing.T) {
	os.Unsetenv("SESSION_DRIVER")
	store, err := NewStore(nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := store.(*MemoryStore); !ok {
		t.Fatal("expected *MemoryStore as default")
	}
}

// --- RedisStore ---

func newTestRedisStore(t *testing.T) (*RedisStore, *miniredis.Miniredis) {
	t.Helper()
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	return NewRedisStore(client, "test:"), s
}

func TestRedisStore_WriteRead(t *testing.T) {
	store, _ := newTestRedisStore(t)
	data := map[string]interface{}{"user": "redis-alice"}
	if err := store.Write("s1", data, 10*time.Minute); err != nil {
		t.Fatal(err)
	}
	got, err := store.Read("s1")
	if err != nil {
		t.Fatal(err)
	}
	if got["user"] != "redis-alice" {
		t.Fatalf("expected 'redis-alice', got %v", got["user"])
	}
}

func TestRedisStore_ReadMissing(t *testing.T) {
	store, _ := newTestRedisStore(t)
	got, err := store.Read("nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestRedisStore_Destroy(t *testing.T) {
	store, _ := newTestRedisStore(t)
	store.Write("s1", map[string]interface{}{"k": "v"}, 10*time.Minute)
	if err := store.Destroy("s1"); err != nil {
		t.Fatal(err)
	}
	got, _ := store.Read("s1")
	if len(got) != 0 {
		t.Fatal("expected empty map after destroy")
	}
}

func TestRedisStore_GC(t *testing.T) {
	store, _ := newTestRedisStore(t)
	// GC is a no-op for Redis (TTL handles expiry) — just ensure it doesn't error.
	if err := store.GC(10 * time.Minute); err != nil {
		t.Fatalf("GC should not error: %v", err)
	}
}

func TestRedisStore_TTLExpiry(t *testing.T) {
	store, mr := newTestRedisStore(t)
	store.Write("exp", map[string]interface{}{"k": "v"}, 2*time.Second)
	// Fast-forward miniredis time
	mr.FastForward(3 * time.Second)
	got, err := store.Read("exp")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatal("expected empty map after TTL expiry")
	}
}

func TestRedisStore_Prefix(t *testing.T) {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	store := NewRedisStore(client, "myapp:")
	store.Write("abc", map[string]interface{}{"k": "v"}, time.Minute)
	// Key in Redis should be prefixed
	if !s.Exists("myapp:abc") {
		t.Fatal("expected key 'myapp:abc' to exist in Redis")
	}
}
