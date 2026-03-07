package cli

import (
	"fmt"
	"os"

	"github.com/RAiWorks/RapidGo/app/providers"
	"github.com/RAiWorks/RapidGo/core/app"
	"github.com/RAiWorks/RapidGo/core/service"
	"github.com/spf13/cobra"
)

// Version is the current framework version.
const Version = "0.2.0"

var rootCmd = &cobra.Command{
	Use:   "rapidgo",
	Short: "RapidGo — A Go web framework with Laravel-style developer experience",
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
	rootCmd.AddCommand(workCmd)
	rootCmd.AddCommand(scheduleRunCmd)
}

// Execute runs the root command. Called from main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// RootCmd returns the root Cobra command, allowing plugins to add subcommands.
func RootCmd() *cobra.Command {
	return rootCmd
}

// NewApp creates and boots a RapidGo application configured for the given mode.
// Used by commands that need the application lifecycle (serve, migrate, etc.).
func NewApp(mode service.Mode) *app.App {
	application := app.New()

	// Always required
	application.Register(&providers.ConfigProvider{})
	application.Register(&providers.LoggerProvider{})

	// DB required for HTTP modes that may access data
	if mode.Has(service.ModeWeb) || mode.Has(service.ModeAPI) || mode.Has(service.ModeWS) {
		application.Register(&providers.DatabaseProvider{})
	}

	// Queue and Redis — lazy singletons, safe to register always
	application.Register(&providers.RedisProvider{})
	application.Register(&providers.QueueProvider{})

	// Session only needed for web mode (cookie-based auth)
	if mode.Has(service.ModeWeb) {
		application.Register(&providers.SessionProvider{})
	}

	// Middleware and Router for any HTTP mode
	application.Register(&providers.MiddlewareProvider{Mode: mode})
	application.Register(&providers.RouterProvider{Mode: mode})

	application.Boot()
	return application
}
