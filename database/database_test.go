package database

import (
	"strings"
	"testing"
	"time"
)

// TC-01: NewDBConfig reads all env vars
func TestNewDBConfig_ReadsEnvVars(t *testing.T) {
	t.Setenv("DB_DRIVER", "postgres")
	t.Setenv("DB_HOST", "testhost")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")
	t.Setenv("DB_SSL_MODE", "require")
	t.Setenv("DB_MAX_OPEN_CONNS", "50")
	t.Setenv("DB_MAX_IDLE_CONNS", "20")
	t.Setenv("DB_CONN_MAX_LIFETIME", "10")
	t.Setenv("DB_CONN_MAX_IDLE_TIME", "7")

	cfg := NewDBConfig()

	if cfg.Driver != "postgres" {
		t.Fatalf("expected Driver 'postgres', got '%s'", cfg.Driver)
	}
	if cfg.Host != "testhost" {
		t.Fatalf("expected Host 'testhost', got '%s'", cfg.Host)
	}
	if cfg.Port != "3306" {
		t.Fatalf("expected Port '3306', got '%s'", cfg.Port)
	}
	if cfg.Name != "testdb" {
		t.Fatalf("expected Name 'testdb', got '%s'", cfg.Name)
	}
	if cfg.User != "testuser" {
		t.Fatalf("expected User 'testuser', got '%s'", cfg.User)
	}
	if cfg.Password != "testpass" {
		t.Fatalf("expected Password 'testpass', got '%s'", cfg.Password)
	}
	if cfg.SSLMode != "require" {
		t.Fatalf("expected SSLMode 'require', got '%s'", cfg.SSLMode)
	}
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

// TC-02: NewDBConfig uses defaults
func TestNewDBConfig_Defaults(t *testing.T) {
	// Clear all DB_ vars to ensure defaults
	t.Setenv("DB_DRIVER", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_SSL_MODE", "")
	t.Setenv("DB_MAX_OPEN_CONNS", "")
	t.Setenv("DB_MAX_IDLE_CONNS", "")
	t.Setenv("DB_CONN_MAX_LIFETIME", "")
	t.Setenv("DB_CONN_MAX_IDLE_TIME", "")

	cfg := NewDBConfig()

	if cfg.Driver != "" {
		t.Fatalf("expected empty Driver, got '%s'", cfg.Driver)
	}
	if cfg.Host != "localhost" {
		t.Fatalf("expected Host 'localhost', got '%s'", cfg.Host)
	}
	if cfg.Port != "5432" {
		t.Fatalf("expected Port '5432', got '%s'", cfg.Port)
	}
	if cfg.Name != "rgo_dev" {
		t.Fatalf("expected Name 'rgo_dev', got '%s'", cfg.Name)
	}
	if cfg.MaxOpenConns != 25 {
		t.Fatalf("expected MaxOpenConns 25, got %d", cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns != 10 {
		t.Fatalf("expected MaxIdleConns 10, got %d", cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime != 5*time.Minute {
		t.Fatalf("expected ConnMaxLifetime 5m, got %v", cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime != 3*time.Minute {
		t.Fatalf("expected ConnMaxIdleTime 3m, got %v", cfg.ConnMaxIdleTime)
	}
}

// TC-03: DSN postgres format
func TestDSN_Postgres(t *testing.T) {
	cfg := DBConfig{
		Driver:   "postgres",
		Host:     "db.example.com",
		User:     "admin",
		Password: "s3cret",
		Name:     "myapp",
		Port:     "5432",
		SSLMode:  "require",
	}

	expected := "host=db.example.com user=admin password=s3cret dbname=myapp port=5432 sslmode=require"
	if got := cfg.DSN(); got != expected {
		t.Fatalf("expected DSN:\n%s\ngot:\n%s", expected, got)
	}
}

// TC-04: DSN mysql format
func TestDSN_MySQL(t *testing.T) {
	cfg := DBConfig{
		Driver:   "mysql",
		Host:     "db.example.com",
		User:     "root",
		Password: "pass",
		Name:     "myapp",
		Port:     "3306",
	}

	expected := "root:pass@tcp(db.example.com:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local"
	if got := cfg.DSN(); got != expected {
		t.Fatalf("expected DSN:\n%s\ngot:\n%s", expected, got)
	}
}

// TC-05: DSN sqlite format
func TestDSN_SQLite(t *testing.T) {
	cfg := DBConfig{
		Driver: "sqlite",
		Name:   ":memory:",
	}

	if got := cfg.DSN(); got != ":memory:" {
		t.Fatalf("expected ':memory:', got '%s'", got)
	}
}

// TC-06: DSN unsupported driver
func TestDSN_UnsupportedDriver(t *testing.T) {
	cfg := DBConfig{Driver: "oracle"}

	if got := cfg.DSN(); got != "" {
		t.Fatalf("expected empty string for unsupported driver, got '%s'", got)
	}
}

// TC-07: ConnectWithConfig succeeds with SQLite in-memory
func TestConnectWithConfig_SQLiteMemory(t *testing.T) {
	cfg := DBConfig{
		Driver:          "sqlite",
		Name:            ":memory:",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: time.Minute,
		ConnMaxIdleTime: time.Minute,
	}

	db, err := ConnectWithConfig(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if db == nil {
		t.Fatal("expected non-nil *gorm.DB")
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("expected no error getting sql.DB, got: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("expected successful ping, got: %v", err)
	}
}

// TC-08: ConnectWithConfig fails with unsupported driver
func TestConnectWithConfig_UnsupportedDriver(t *testing.T) {
	cfg := DBConfig{Driver: "oracle"}

	db, err := ConnectWithConfig(cfg)
	if err == nil {
		t.Fatal("expected error for unsupported driver")
	}
	if db != nil {
		t.Fatal("expected nil *gorm.DB for unsupported driver")
	}
	if !strings.Contains(err.Error(), "unsupported DB_DRIVER: oracle") {
		t.Fatalf("expected error to contain 'unsupported DB_DRIVER: oracle', got: %v", err)
	}
}

// TC-09: Pool settings are applied
func TestConnectWithConfig_PoolSettings(t *testing.T) {
	cfg := DBConfig{
		Driver:          "sqlite",
		Name:            ":memory:",
		MaxOpenConns:    42,
		MaxIdleConns:    7,
		ConnMaxLifetime: 10 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}

	db, err := ConnectWithConfig(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("expected no error getting sql.DB, got: %v", err)
	}

	stats := sqlDB.Stats()
	if stats.MaxOpenConnections != 42 {
		t.Fatalf("expected MaxOpenConnections 42, got %d", stats.MaxOpenConnections)
	}
}
