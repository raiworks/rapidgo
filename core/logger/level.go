package logger

import "log/slog"

// parseLevel converts a string level name to a slog.Level.
// Returns slog.LevelInfo if the string is unrecognized.
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
