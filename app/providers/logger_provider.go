package providers

import (
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/core/logger"
)

// LoggerProvider sets up structured logging via slog.
// Must be registered after ConfigProvider — logger reads
// LOG_LEVEL, LOG_FORMAT, LOG_OUTPUT from environment.
type LoggerProvider struct{}

// Register is a no-op. Logger setup requires config values.
func (p *LoggerProvider) Register(c *container.Container) {}

// Boot initializes the logger. Runs after all providers have
// registered, so config values are guaranteed available.
func (p *LoggerProvider) Boot(c *container.Container) {
	logger.Setup()
}
