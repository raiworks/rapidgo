package cli

import (
	"testing"
	"time"
)

func TestParseDuration_ValidString(t *testing.T) {
	d := parseDuration("30s", 15*time.Second)
	if d != 30*time.Second {
		t.Errorf("parseDuration(\"30s\") = %v, want 30s", d)
	}
}

func TestParseDuration_InvalidFallsBack(t *testing.T) {
	d := parseDuration("invalid", 15*time.Second)
	if d != 15*time.Second {
		t.Errorf("parseDuration(\"invalid\") = %v, want 15s", d)
	}
}

func TestParseDuration_EmptyFallsBack(t *testing.T) {
	d := parseDuration("", 15*time.Second)
	if d != 15*time.Second {
		t.Errorf("parseDuration(\"\") = %v, want 15s", d)
	}
}

func TestResolveServerTimeouts_Defaults(t *testing.T) {
	// Clear any env vars that might interfere
	t.Setenv("SERVER_READ_TIMEOUT", "")
	t.Setenv("SERVER_WRITE_TIMEOUT", "")
	t.Setenv("SERVER_IDLE_TIMEOUT", "")
	t.Setenv("SERVER_SHUTDOWN_TIMEOUT", "")

	read, write, idle, shutdown := resolveServerTimeouts()
	if read != 15*time.Second {
		t.Errorf("read = %v, want 15s", read)
	}
	if write != 15*time.Second {
		t.Errorf("write = %v, want 15s", write)
	}
	if idle != 60*time.Second {
		t.Errorf("idle = %v, want 60s", idle)
	}
	if shutdown != 30*time.Second {
		t.Errorf("shutdown = %v, want 30s", shutdown)
	}
}

func TestResolveServerTimeouts_CustomValues(t *testing.T) {
	t.Setenv("SERVER_READ_TIMEOUT", "30s")
	t.Setenv("SERVER_WRITE_TIMEOUT", "")
	t.Setenv("SERVER_IDLE_TIMEOUT", "2m")
	t.Setenv("SERVER_SHUTDOWN_TIMEOUT", "")

	read, write, idle, shutdown := resolveServerTimeouts()
	if read != 30*time.Second {
		t.Errorf("read = %v, want 30s", read)
	}
	if write != 15*time.Second {
		t.Errorf("write = %v, want 15s (default)", write)
	}
	if idle != 2*time.Minute {
		t.Errorf("idle = %v, want 2m", idle)
	}
	if shutdown != 30*time.Second {
		t.Errorf("shutdown = %v, want 30s (default)", shutdown)
	}
}
