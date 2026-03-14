package cli

import (
	"fmt"

	"github.com/raiworks/rapidgo/v2/core/config"
	"github.com/raiworks/rapidgo/v2/core/container"
	"github.com/raiworks/rapidgo/v2/core/service"
	"github.com/raiworks/rapidgo/v2/database/migrations"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var migrateFreshCmd = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Drop all tables and re-run all migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.AppEnv() == "production" {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				return fmt.Errorf("refusing to run migrate:fresh in production without --force")
			}
		}

		application := NewApp(service.ModeAll)
		db := container.MustMake[*gorm.DB](application.Container, "db")

		// Drop all tables
		migrator := db.Migrator()
		tables, err := migrator.GetTables()
		if err != nil {
			return fmt.Errorf("failed to get tables: %w", err)
		}
		for _, table := range tables {
			if err := migrator.DropTable(table); err != nil {
				return fmt.Errorf("failed to drop table %s: %w", table, err)
			}
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Dropped all tables.")

		// Re-run migrations (same as migrate command)
		if modelRegistryFn != nil {
			if err := db.AutoMigrate(modelRegistryFn()...); err != nil {
				return fmt.Errorf("auto-migrate failed: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "AutoMigrate complete.")
		}

		m, err := migrations.NewMigrator(db)
		if err != nil {
			return err
		}
		n, err := m.Run()
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
		if n == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Nothing to migrate.")
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Applied %d migration(s).\n", n)
		}
		return nil
	},
}

func init() {
	migrateFreshCmd.Flags().Bool("force", false, "Force the operation in production")
}
