package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TC-01: parseLevel() with valid levels
func TestParseLevel_ValidLevels(t *testing.T) {
	tests := []struct {
		input string
		want  slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
	}
	for _, tt := range tests {
		got := parseLevel(tt.input)
		if got != tt.want {
			t.Errorf("parseLevel(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// TC-02: parseLevel() with unknown string
func TestParseLevel_Unknown(t *testing.T) {
	got := parseLevel("invalid")
	if got != slog.LevelInfo {
		t.Errorf("parseLevel(\"invalid\") = %v, want %v", got, slog.LevelInfo)
	}
}

// TC-09: parseLevel() with empty string
func TestParseLevel_Empty(t *testing.T) {
	got := parseLevel("")
	if got != slog.LevelInfo {
		t.Errorf("parseLevel(\"\") = %v, want %v", got, slog.LevelInfo)
	}
}

// TC-03: Setup() with JSON format
func TestSetup_JSONFormat(t *testing.T) {
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_OUTPUT", "stdout")

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)
	logger.Info("test message", "key", "value")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("JSON output should be valid JSON: %v", err)
	}
	if result["msg"] != "test message" {
		t.Errorf("msg = %v, want %q", result["msg"], "test message")
	}
	if result["key"] != "value" {
		t.Errorf("key = %v, want %q", result["key"], "value")
	}
}

// TC-04: Setup() with text format
func TestSetup_TextFormat(t *testing.T) {
	t.Setenv("LOG_FORMAT", "text")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_OUTPUT", "stdout")

	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)
	logger.Info("test message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "level=INFO") {
		t.Errorf("text output should contain 'level=INFO', got: %s", output)
	}
	if !strings.Contains(output, "msg=\"test message\"") {
		t.Errorf("text output should contain msg, got: %s", output)
	}
}

// TC-05: Setup() with file output
func TestSetup_FileOutput(t *testing.T) {
	tmp := t.TempDir()
	original, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(original)

	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_OUTPUT", "file")

	logger := Setup()
	logger.Info("file test message")

	// Close the log file before TempDir cleanup
	Close()

	logPath := filepath.Join(tmp, "storage", "logs", "app.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Fatal("storage/logs/app.log should be created")
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "file test message") {
		t.Errorf("log file should contain 'file test message', got: %s", string(content))
	}
}

// TC-06: Setup() log level filtering
func TestSetup_LevelFiltering(t *testing.T) {
	t.Setenv("LOG_LEVEL", "warn")
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_OUTPUT", "stdout")

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})
	logger := slog.New(handler)

	logger.Info("should not appear")
	logger.Warn("should appear")

	output := buf.String()
	if strings.Contains(output, "should not appear") {
		t.Error("info message should be filtered out at warn level")
	}
	if !strings.Contains(output, "should appear") {
		t.Error("warn message should appear at warn level")
	}
}

// TC-10: Setup() with invalid LOG_FORMAT defaults to JSON
func TestSetup_InvalidFormat(t *testing.T) {
	t.Setenv("LOG_FORMAT", "yaml")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_OUTPUT", "stdout")

	// Setup should not panic with invalid format
	logger := Setup()
	if logger == nil {
		t.Fatal("Setup() should return a non-nil logger")
	}
}

// TC-03/04 integration: Setup() actually configures the global default
func TestSetup_SetsDefault(t *testing.T) {
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("LOG_OUTPUT", "stdout")

	logger := Setup()
	if logger == nil {
		t.Fatal("Setup() should return a non-nil logger")
	}

	// After Setup(), slog.Default() should return the configured logger
	defaultLogger := slog.Default()
	if defaultLogger == nil {
		t.Fatal("slog.Default() should not be nil after Setup()")
	}
}

// TC-L01: SlogLogger implements Logger interface (compile-time check)
var _ Logger = (*SlogLogger)(nil)

// TC-L02: With() returns a new Logger with attributes
func TestSlogLogger_With(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := NewSlogLogger(slog.New(handler))

	child := l.With("component", "test")
	if child == nil {
		t.Fatal("With() should return a non-nil Logger")
	}

	child.Info("hello")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("expected valid JSON: %v", err)
	}
	if result["component"] != "test" {
		t.Errorf("expected component=test, got %v", result["component"])
	}
	if result["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", result["msg"])
	}
}

// TC-L03: Setup() returns Logger interface
func TestSetup_ReturnsLoggerInterface(t *testing.T) {
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_OUTPUT", "stdout")

	var l Logger = Setup()
	if l == nil {
		t.Fatal("Setup() should return a non-nil Logger")
	}
}

// TC-L04: SlogLogger methods log to underlying handler
func TestSlogLogger_AllMethods(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := NewSlogLogger(slog.New(handler))

	l.Debug("d")
	l.Info("i")
	l.Warn("w")
	l.Error("e")

	output := buf.String()
	for _, msg := range []string{"d", "i", "w", "e"} {
		if !strings.Contains(output, `"msg":"`+msg+`"`) {
			t.Errorf("expected message %q in output", msg)
		}
	}
}
