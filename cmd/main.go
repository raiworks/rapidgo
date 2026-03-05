package main

import (
	"fmt"
	"log/slog"

	"github.com/RAiWorks/RGo/app/providers"
	"github.com/RAiWorks/RGo/core/app"
	"github.com/RAiWorks/RGo/core/config"
)

func main() {
	application := app.New()

	// Register providers (order matters)
	application.Register(&providers.ConfigProvider{})  // 1. Config first — loads .env
	application.Register(&providers.LoggerProvider{})  // 2. Logger — uses config in Boot

	// Boot all providers
	application.Boot()

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
