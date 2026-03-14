package logger

import (
	"log/slog"
	"os"

	"github.com/raiworks/rapidgo/v2/core/config"
)

// Logger defines the logging contract for RapidGo.
// The default implementation wraps slog. Users can implement this
// interface to plug in Zap, Zerolog, or a test spy.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
}

// SlogLogger wraps *slog.Logger to implement the Logger interface.
type SlogLogger struct {
	log *slog.Logger
}

// NewSlogLogger creates a Logger wrapping the given *slog.Logger.
func NewSlogLogger(l *slog.Logger) *SlogLogger {
	return &SlogLogger{log: l}
}

func (s *SlogLogger) Debug(msg string, args ...any) { s.log.Debug(msg, args...) }
func (s *SlogLogger) Info(msg string, args ...any)  { s.log.Info(msg, args...) }
func (s *SlogLogger) Warn(msg string, args ...any)  { s.log.Warn(msg, args...) }
func (s *SlogLogger) Error(msg string, args ...any) { s.log.Error(msg, args...) }

// With returns a new Logger with the given attributes attached to every log message.
func (s *SlogLogger) With(args ...any) Logger {
	return &SlogLogger{log: s.log.With(args...)}
}

// logFile holds the open log file handle (if LOG_OUTPUT=file).
var logFile *os.File

// Setup initializes the global slog logger based on config values.
// Reads LOG_LEVEL, LOG_FORMAT, LOG_OUTPUT from environment.
// Sets slog.SetDefault() so that slog.Info(), slog.Error() etc. work globally.
// Returns a Logger interface wrapping the configured slog instance.
func Setup() Logger {
	level := parseLevel(config.Env("LOG_LEVEL", "info"))
	format := config.Env("LOG_FORMAT", "json")
	output := config.Env("LOG_OUTPUT", "stdout")

	var writer *os.File
	if output == "file" {
		if err := os.MkdirAll("storage/logs", 0755); err != nil {
			slog.Warn("failed to create log directory, falling back to stdout", "err", err)
			writer = os.Stdout
		} else {
			f, err := os.OpenFile("storage/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				slog.Warn("failed to open log file, falling back to stdout", "err", err)
				writer = os.Stdout
			} else {
				logFile = f
				writer = f
			}
		}
	} else {
		writer = os.Stdout
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if format == "text" {
		handler = slog.NewTextHandler(writer, opts)
	} else {
		handler = slog.NewJSONHandler(writer, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return NewSlogLogger(logger)
}

// Close closes the log file if one was opened.
// Should be called on application shutdown.
func Close() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}
