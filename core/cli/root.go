package cli

import (
	"fmt"
	"os"

	"github.com/raiworks/rapidgo/v2/core/app"
	"github.com/raiworks/rapidgo/v2/core/config"
	"github.com/raiworks/rapidgo/v2/core/service"
	"github.com/spf13/cobra"
)

// Version is the current framework version.
const Version = "2.4.0"

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
	rootCmd.AddCommand(makeAdminCmd)
	rootCmd.AddCommand(newCmd)
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
	config.Load()

	application := app.New()

	if bootstrapFn != nil {
		bootstrapFn(application, mode)
	}

	application.Boot()
	return application
}
