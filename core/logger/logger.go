package logger

import (
	"log/slog"
	"os"

	"github.com/RAiWorks/RGo/core/config"
)

// logFile holds the open log file handle (if LOG_OUTPUT=file).
var logFile *os.File

// Setup initializes the global slog logger based on config values.
// Reads LOG_LEVEL, LOG_FORMAT, LOG_OUTPUT from environment.
// Sets slog.SetDefault() so that slog.Info(), slog.Error() etc. work globally.
// Returns the configured logger instance.
func Setup() *slog.Logger {
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
	return logger
}

// Close closes the log file if one was opened.
// Should be called on application shutdown.
func Close() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}
