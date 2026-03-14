package config

import (
	"os"
	"testing"
	"time"
)

// TC-01: Load() with .env file present (handled by integration test — go run)

// TC-02: Load() without .env file — should not panic
func TestLoad_NoEnvFile(t *testing.T) {
	// Change to a temp dir where no .env exists
	original, _ := os.Getwd()
	tmp := t.TempDir()
	os.Chdir(tmp)
	defer os.Chdir(original)

	// Should not panic — just logs a message
	Load()
}

// TC-03: Env() with key present
func TestEnv_KeyPresent(t *testing.T) {
	t.Setenv("TEST_KEY", "hello")
	got := Env("TEST_KEY", "default")
	if got != "hello" {
		t.Errorf("Env() = %q, want %q", got, "hello")
	}
}

// TC-04: Env() with key absent (fallback)
func TestEnv_KeyAbsent(t *testing.T) {
	got := Env("TEST_MISSING_KEY_XYZ", "fallback_value")
	if got != "fallback_value" {
		t.Errorf("Env() = %q, want %q", got, "fallback_value")
	}
}

// TC-13: Env() with empty string value
func TestEnv_EmptyValue(t *testing.T) {
	t.Setenv("TEST_EMPTY", "")
	got := Env("TEST_EMPTY", "fallback")
	if got != "fallback" {
		t.Errorf("Env() = %q, want %q", got, "fallback")
	}
}

// TC-05: EnvInt() with valid integer
func TestEnvInt_ValidInt(t *testing.T) {
	t.Setenv("TEST_INT", "42")
	got := EnvInt("TEST_INT", 0)
	if got != 42 {
		t.Errorf("EnvInt() = %d, want %d", got, 42)
	}
}

// TC-06: EnvInt() with invalid string
func TestEnvInt_InvalidString(t *testing.T) {
	t.Setenv("TEST_INT_BAD", "not_a_number")
	got := EnvInt("TEST_INT_BAD", 99)
	if got != 99 {
		t.Errorf("EnvInt() = %d, want %d", got, 99)
	}
}

// EnvInt() with empty value (fallback)
func TestEnvInt_Empty(t *testing.T) {
	got := EnvInt("TEST_INT_MISSING_XYZ", 77)
	if got != 77 {
		t.Errorf("EnvInt() = %d, want %d", got, 77)
	}
}

// TC-07: EnvBool() truthy values
func TestEnvBool_Truthy(t *testing.T) {
	t.Setenv("TEST_BOOL_T", "true")
	if !EnvBool("TEST_BOOL_T", false) {
		t.Error("EnvBool(\"true\") should return true")
	}

	t.Setenv("TEST_BOOL_1", "1")
	if !EnvBool("TEST_BOOL_1", false) {
		t.Error("EnvBool(\"1\") should return true")
	}
}

// TC-08: EnvBool() falsy values
func TestEnvBool_Falsy(t *testing.T) {
	t.Setenv("TEST_BOOL_F", "false")
	if EnvBool("TEST_BOOL_F", true) {
		t.Error("EnvBool(\"false\") should return false")
	}

	t.Setenv("TEST_BOOL_0", "0")
	if EnvBool("TEST_BOOL_0", true) {
		t.Error("EnvBool(\"0\") should return false")
	}
}

// TC-14: EnvBool() with empty (fallback)
func TestEnvBool_Empty(t *testing.T) {
	got := EnvBool("TEST_BOOL_MISSING_XYZ", true)
	if !got {
		t.Error("EnvBool() with missing key should return fallback (true)")
	}
}

// TC-09: Environment detection functions
func TestEnvironmentDetection_Production(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	if !IsProduction() {
		t.Error("IsProduction() should return true")
	}
	if IsDevelopment() {
		t.Error("IsDevelopment() should return false")
	}
	if IsTesting() {
		t.Error("IsTesting() should return false")
	}
}

func TestEnvironmentDetection_Development(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	if !IsDevelopment() {
		t.Error("IsDevelopment() should return true")
	}
	if IsProduction() {
		t.Error("IsProduction() should return false")
	}
}

func TestEnvironmentDetection_Testing(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	if !IsTesting() {
		t.Error("IsTesting() should return true")
	}
	if IsProduction() {
		t.Error("IsProduction() should return false")
	}
}

func TestEnvironmentDetection_Default(t *testing.T) {
	t.Setenv("APP_ENV", "")
	if !IsDevelopment() {
		t.Error("AppEnv() should default to development")
	}
}

func TestIsDebug(t *testing.T) {
	t.Setenv("APP_DEBUG", "true")
	if !IsDebug() {
		t.Error("IsDebug() should return true when APP_DEBUG=true")
	}

	t.Setenv("APP_DEBUG", "false")
	if IsDebug() {
		t.Error("IsDebug() should return false when APP_DEBUG=false")
	}
}

// --- LoadConfig[T]() tests ---

type testConfig struct {
	Name    string        `env:"TC_NAME" default:"myapp"`
	Port    int           `env:"TC_PORT" default:"8080"`
	Debug   bool          `env:"TC_DEBUG" default:"false"`
	Rate    float64       `env:"TC_RATE" default:"1.5"`
	Timeout time.Duration `env:"TC_TIMEOUT" default:"30s"`
	Ignored string        // no env tag — should be skipped
}

// TC-C01: All basic types parsed from env
func TestLoadConfig_AllTypes(t *testing.T) {
	t.Setenv("TC_NAME", "testapp")
	t.Setenv("TC_PORT", "9090")
	t.Setenv("TC_DEBUG", "true")
	t.Setenv("TC_RATE", "2.5")
	t.Setenv("TC_TIMEOUT", "10s")

	cfg, err := LoadConfig[testConfig]()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Name != "testapp" {
		t.Errorf("Name = %q, want %q", cfg.Name, "testapp")
	}
	if cfg.Port != 9090 {
		t.Errorf("Port = %d, want %d", cfg.Port, 9090)
	}
	if !cfg.Debug {
		t.Error("Debug = false, want true")
	}
	if cfg.Rate != 2.5 {
		t.Errorf("Rate = %f, want %f", cfg.Rate, 2.5)
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 10*time.Second)
	}
}

// TC-C02: Default values used when env empty
func TestLoadConfig_Defaults(t *testing.T) {
	// Clear all env vars used by testConfig
	t.Setenv("TC_NAME", "")
	t.Setenv("TC_PORT", "")
	t.Setenv("TC_DEBUG", "")
	t.Setenv("TC_RATE", "")
	t.Setenv("TC_TIMEOUT", "")

	cfg, err := LoadConfig[testConfig]()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Name != "myapp" {
		t.Errorf("Name = %q, want default %q", cfg.Name, "myapp")
	}
	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want default %d", cfg.Port, 8080)
	}
	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want default %v", cfg.Timeout, 30*time.Second)
	}
}

// TC-C03: Validation failure for required field
func TestLoadConfig_ValidationFail(t *testing.T) {
	type requiredConfig struct {
		Name string `env:"TC_REQ_NAME" validate:"required"`
	}
	t.Setenv("TC_REQ_NAME", "")

	_, err := LoadConfig[requiredConfig]()
	if err == nil {
		t.Fatal("expected validation error for empty required field")
	}
}

// TC-C04: Required field present
func TestLoadConfig_RequiredPresent(t *testing.T) {
	type requiredConfig struct {
		Name string `env:"TC_REQ_NAME2" validate:"required"`
	}
	t.Setenv("TC_REQ_NAME2", "hello")

	cfg, err := LoadConfig[requiredConfig]()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Name != "hello" {
		t.Errorf("Name = %q, want %q", cfg.Name, "hello")
	}
}

// TC-C05: Duration parsing
func TestLoadConfig_Duration(t *testing.T) {
	type durConfig struct {
		Timeout time.Duration `env:"TC_DUR" default:"5m"`
	}
	t.Setenv("TC_DUR", "")

	cfg, err := LoadConfig[durConfig]()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Timeout != 5*time.Minute {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 5*time.Minute)
	}
}

// TC-C06: Slice string comma split
func TestLoadConfig_StringSlice(t *testing.T) {
	type sliceConfig struct {
		Hosts []string `env:"TC_HOSTS"`
	}
	t.Setenv("TC_HOSTS", "a, b, c")

	cfg, err := LoadConfig[sliceConfig]()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if len(cfg.Hosts) != 3 || cfg.Hosts[0] != "a" || cfg.Hosts[1] != "b" || cfg.Hosts[2] != "c" {
		t.Errorf("Hosts = %v, want [a b c]", cfg.Hosts)
	}
}

// TC-C07: Unsupported field type returns error
func TestLoadConfig_UnsupportedType(t *testing.T) {
	type badConfig struct {
		Val complex64 `env:"TC_COMPLEX" default:"1"`
	}
	t.Setenv("TC_COMPLEX", "")

	_, err := LoadConfig[badConfig]()
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

// TC-C08: Fields without env tag are skipped
func TestLoadConfig_NoEnvTag(t *testing.T) {
	type partialConfig struct {
		Tagged   string `env:"TC_TAGGED" default:"val"`
		Untagged string
	}
	t.Setenv("TC_TAGGED", "")

	cfg, err := LoadConfig[partialConfig]()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Tagged != "val" {
		t.Errorf("Tagged = %q, want %q", cfg.Tagged, "val")
	}
	if cfg.Untagged != "" {
		t.Errorf("Untagged = %q, want empty (skipped)", cfg.Untagged)
	}
}
