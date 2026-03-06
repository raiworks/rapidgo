package cli

import (
	"fmt"
	"os"

	"github.com/RAiWorks/RGo/app/providers"
	"github.com/RAiWorks/RGo/core/app"
	"github.com/spf13/cobra"
)

// Version is the current framework version.
const Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "rgo",
	Short: "RGo — A Go web framework with Laravel-style developer experience",
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(migrateRollbackCmd)
	rootCmd.AddCommand(migrateStatusCmd)
	rootCmd.AddCommand(makeMigrationCmd)
	rootCmd.AddCommand(dbSeedCmd)
	rootCmd.AddCommand(makeControllerCmd)
	rootCmd.AddCommand(makeModelCmd)
	rootCmd.AddCommand(makeServiceCmd)
	rootCmd.AddCommand(makeProviderCmd)
}

// Execute runs the root command. Called from main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// NewApp creates and boots a fully configured RGo application.
// Used by commands that need the application lifecycle (serve, migrate, etc.).
func NewApp() *app.App {
	application := app.New()
	application.Register(&providers.ConfigProvider{})
	application.Register(&providers.LoggerProvider{})
	application.Register(&providers.DatabaseProvider{})
	application.Register(&providers.SessionProvider{})
	application.Register(&providers.MiddlewareProvider{})
	application.Register(&providers.RouterProvider{})
	application.Boot()
	return application
}
