package database

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// T01: NewResolver stores writer and reader correctly
func TestNewResolver(t *testing.T) {
	w := &gorm.DB{}
	r := &gorm.DB{}
	res := NewResolver(w, r)
	if res.writer != w {
		t.Fatal("writer not stored")
	}
	if res.reader != r {
		t.Fatal("reader not stored")
	}
}

// T02: Writer returns the writer *gorm.DB
func TestWriter(t *testing.T) {
	w := &gorm.DB{}
	r := &gorm.DB{}
	res := NewResolver(w, r)
	if res.Writer() != w {
		t.Fatal("Writer() did not return writer")
	}
}

// T03: Reader returns the reader *gorm.DB
func TestReader(t *testing.T) {
	w := &gorm.DB{}
	r := &gorm.DB{}
	res := NewResolver(w, r)
	if res.Reader() != r {
		t.Fatal("Reader() did not return reader")
	}
}

// T04: When same *gorm.DB passed for both, Reader returns same instance as Writer
func TestReaderFallback(t *testing.T) {
	w := &gorm.DB{}
	res := NewResolver(w, w)
	if res.Reader() != res.Writer() {
		t.Fatal("Reader() should be same instance as Writer() when both are the same")
	}
}

// T05: NewReadDBConfig reads DB_READ_* env vars
func TestNewReadDBConfig_Explicit(t *testing.T) {
	t.Setenv("DB_READ_DRIVER", "mysql")
	t.Setenv("DB_READ_HOST", "replica-host")
	t.Setenv("DB_READ_PORT", "3307")
	t.Setenv("DB_READ_NAME", "replica_db")
	t.Setenv("DB_READ_USER", "reader")
	t.Setenv("DB_READ_PASSWORD", "readpass")
	t.Setenv("DB_READ_SSL_MODE", "require")
	t.Setenv("DB_READ_MAX_OPEN_CONNS", "100")
	t.Setenv("DB_READ_MAX_IDLE_CONNS", "50")
	t.Setenv("DB_READ_CONN_MAX_LIFETIME", "15")
	t.Setenv("DB_READ_CONN_MAX_IDLE_TIME", "8")

	cfg := NewReadDBConfig()

	if cfg.Driver != "mysql" {
		t.Fatalf("expected Driver 'mysql', got '%s'", cfg.Driver)
	}
	if cfg.Host != "replica-host" {
		t.Fatalf("expected Host 'replica-host', got '%s'", cfg.Host)
	}
	if cfg.Port != "3307" {
		t.Fatalf("expected Port '3307', got '%s'", cfg.Port)
	}
	if cfg.Name != "replica_db" {
		t.Fatalf("expected Name 'replica_db', got '%s'", cfg.Name)
	}
	if cfg.User != "reader" {
		t.Fatalf("expected User 'reader', got '%s'", cfg.User)
	}
	if cfg.Password != "readpass" {
		t.Fatalf("expected Password 'readpass', got '%s'", cfg.Password)
	}
	if cfg.SSLMode != "require" {
		t.Fatalf("expected SSLMode 'require', got '%s'", cfg.SSLMode)
	}
	if cfg.MaxOpenConns != 100 {
		t.Fatalf("expected MaxOpenConns 100, got %d", cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns != 50 {
		t.Fatalf("expected MaxIdleConns 50, got %d", cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime != 15*time.Minute {
		t.Fatalf("expected ConnMaxLifetime 15m, got %v", cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime != 8*time.Minute {
		t.Fatalf("expected ConnMaxIdleTime 8m, got %v", cfg.ConnMaxIdleTime)
	}
}

// T06: Without DB_READ_*, NewReadDBConfig falls back to DB_* values
func TestNewReadDBConfig_Fallback(t *testing.T) {
	t.Setenv("DB_DRIVER", "postgres")
	t.Setenv("DB_HOST", "primary-host")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_NAME", "primary_db")
	t.Setenv("DB_USER", "primary_user")
	t.Setenv("DB_PASSWORD", "primary_pass")
	t.Setenv("DB_SSL_MODE", "verify-full")
	t.Setenv("DB_MAX_OPEN_CONNS", "30")
	t.Setenv("DB_MAX_IDLE_CONNS", "15")
	t.Setenv("DB_CONN_MAX_LIFETIME", "12")
	t.Setenv("DB_CONN_MAX_IDLE_TIME", "6")

	cfg := NewReadDBConfig()

	if cfg.Driver != "postgres" {
		t.Fatalf("expected Driver 'postgres', got '%s'", cfg.Driver)
	}
	if cfg.Host != "primary-host" {
		t.Fatalf("expected Host 'primary-host', got '%s'", cfg.Host)
	}
	if cfg.Port != "5433" {
		t.Fatalf("expected Port '5433', got '%s'", cfg.Port)
	}
	if cfg.Name != "primary_db" {
		t.Fatalf("expected Name 'primary_db', got '%s'", cfg.Name)
	}
	if cfg.User != "primary_user" {
		t.Fatalf("expected User 'primary_user', got '%s'", cfg.User)
	}
	if cfg.Password != "primary_pass" {
		t.Fatalf("expected Password 'primary_pass', got '%s'", cfg.Password)
	}
	if cfg.SSLMode != "verify-full" {
		t.Fatalf("expected SSLMode 'verify-full', got '%s'", cfg.SSLMode)
	}
	if cfg.MaxOpenConns != 30 {
		t.Fatalf("expected MaxOpenConns 30, got %d", cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns != 15 {
		t.Fatalf("expected MaxIdleConns 15, got %d", cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime != 12*time.Minute {
		t.Fatalf("expected ConnMaxLifetime 12m, got %v", cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime != 6*time.Minute {
		t.Fatalf("expected ConnMaxIdleTime 6m, got %v", cfg.ConnMaxIdleTime)
	}
}

// T07: Without any env vars, NewReadDBConfig returns same defaults as NewDBConfig
func TestNewReadDBConfig_Defaults(t *testing.T) {
	cfg := NewReadDBConfig()
	def := NewDBConfig()

	if cfg.Driver != def.Driver {
		t.Fatalf("Driver: got '%s', want '%s'", cfg.Driver, def.Driver)
	}
	if cfg.Host != def.Host {
		t.Fatalf("Host: got '%s', want '%s'", cfg.Host, def.Host)
	}
	if cfg.Port != def.Port {
		t.Fatalf("Port: got '%s', want '%s'", cfg.Port, def.Port)
	}
	if cfg.Name != def.Name {
		t.Fatalf("Name: got '%s', want '%s'", cfg.Name, def.Name)
	}
	if cfg.User != def.User {
		t.Fatalf("User: got '%s', want '%s'", cfg.User, def.User)
	}
	if cfg.Password != def.Password {
		t.Fatalf("Password: got '%s', want '%s'", cfg.Password, def.Password)
	}
	if cfg.SSLMode != def.SSLMode {
		t.Fatalf("SSLMode: got '%s', want '%s'", cfg.SSLMode, def.SSLMode)
	}
	if cfg.MaxOpenConns != def.MaxOpenConns {
		t.Fatalf("MaxOpenConns: got %d, want %d", cfg.MaxOpenConns, def.MaxOpenConns)
	}
	if cfg.MaxIdleConns != def.MaxIdleConns {
		t.Fatalf("MaxIdleConns: got %d, want %d", cfg.MaxIdleConns, def.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime != def.ConnMaxLifetime {
		t.Fatalf("ConnMaxLifetime: got %v, want %v", cfg.ConnMaxLifetime, def.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime != def.ConnMaxIdleTime {
		t.Fatalf("ConnMaxIdleTime: got %v, want %v", cfg.ConnMaxIdleTime, def.ConnMaxIdleTime)
	}
}

// T08: Two distinct SQLite connections via Resolver are independent and functional
func TestResolverWithSQLite(t *testing.T) {
	wDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open writer: %v", err)
	}
	rDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open reader: %v", err)
	}

	res := NewResolver(wDB, rDB)

	// Writer and reader should be distinct instances
	if res.Writer() == res.Reader() {
		t.Fatal("writer and reader should be distinct connections")
	}

	// Both should be functional — execute a simple query
	var result int
	if err := res.Writer().Raw("SELECT 1").Scan(&result).Error; err != nil {
		t.Fatalf("writer query failed: %v", err)
	}
	if result != 1 {
		t.Fatalf("writer: expected 1, got %d", result)
	}

	result = 0
	if err := res.Reader().Raw("SELECT 2").Scan(&result).Error; err != nil {
		t.Fatalf("reader query failed: %v", err)
	}
	if result != 2 {
		t.Fatalf("reader: expected 2, got %d", result)
	}
}

// T09: Same SQLite connection for both — queries work through both accessors
func TestResolverSameConnection(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	res := NewResolver(db, db)

	// Create a table through Writer
	if err := res.Writer().Exec("CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT)").Error; err != nil {
		t.Fatalf("create table failed: %v", err)
	}

	// Insert through Writer
	if err := res.Writer().Exec("INSERT INTO items (name) VALUES ('test')").Error; err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	// Read through Reader — should see the same data
	var name string
	if err := res.Reader().Raw("SELECT name FROM items WHERE id = 1").Scan(&name).Error; err != nil {
		t.Fatalf("reader query failed: %v", err)
	}
	if name != "test" {
		t.Fatalf("expected 'test', got '%s'", name)
	}
}

// T10: DB_READ_* pool settings are read independently from DB_*
func TestNewReadDBConfig_PoolSettings(t *testing.T) {
	// Set primary pool settings
	t.Setenv("DB_MAX_OPEN_CONNS", "25")
	t.Setenv("DB_MAX_IDLE_CONNS", "10")
	t.Setenv("DB_CONN_MAX_LIFETIME", "5")
	t.Setenv("DB_CONN_MAX_IDLE_TIME", "3")

	// Set different read pool settings
	t.Setenv("DB_READ_MAX_OPEN_CONNS", "50")
	t.Setenv("DB_READ_MAX_IDLE_CONNS", "20")
	t.Setenv("DB_READ_CONN_MAX_LIFETIME", "10")
	t.Setenv("DB_READ_CONN_MAX_IDLE_TIME", "7")

	cfg := NewReadDBConfig()

	if cfg.MaxOpenConns != 50 {
		t.Fatalf("expected MaxOpenConns 50, got %d", cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns != 20 {
		t.Fatalf("expected MaxIdleConns 20, got %d", cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime != 10*time.Minute {
		t.Fatalf("expected ConnMaxLifetime 10m, got %v", cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime != 7*time.Minute {
		t.Fatalf("expected ConnMaxIdleTime 7m, got %v", cfg.ConnMaxIdleTime)
	}
}
