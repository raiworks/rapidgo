package main

import (
	"fmt"
	"log/slog"

	"github.com/RAiWorks/RGo/core/config"
	"github.com/RAiWorks/RGo/core/logger"
)

func main() {
	config.Load()
	logger.Setup()

	appName := config.Env("APP_NAME", "RGo")
	appPort := config.Env("APP_PORT", "8080")
	appEnv := config.AppEnv()

	fmt.Println("=================================")
	fmt.Printf("  %s Framework\n", appName)
	fmt.Println("  github.com/RAiWorks/RGo")
	fmt.Println("=================================")
	fmt.Printf("  Environment: %s\n", appEnv)
	fmt.Printf("  Port: %s\n", appPort)
	fmt.Printf("  Debug: %v\n", config.IsDebug())
	fmt.Println("=================================")

	slog.Info("server initialized",
		"app", appName,
		"port", appPort,
		"env", appEnv,
	)
}
