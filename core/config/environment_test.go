package config

import (
	"testing"
)

// T-022: APP_ENV=local → IsLocal() = true
func TestIsLocal_Local(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	if !IsLocal() {
		t.Error("IsLocal() = false, want true for APP_ENV=local")
	}
}

// T-023: APP_ENV=development → IsLocal() = true
func TestIsLocal_Development(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	if !IsLocal() {
		t.Error("IsLocal() = false, want true for APP_ENV=development")
	}
}

// T-024: APP_ENV=production → IsLocal() = false
func TestIsLocal_Production(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	if IsLocal() {
		t.Error("IsLocal() = true, want false for APP_ENV=production")
	}
}

// T-025: APP_ENV=testing → IsLocal() = false
func TestIsLocal_Testing(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	if IsLocal() {
		t.Error("IsLocal() = true, want false for APP_ENV=testing")
	}
}
