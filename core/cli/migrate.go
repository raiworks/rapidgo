package cli

import (
	"fmt"

	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/database/migrations"
	"github.com/RAiWorks/RGo/database/models"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		application := NewApp()
		db := container.MustMake[*gorm.DB](application.Container, "db")

		// Step 1: AutoMigrate all models
		if err := db.AutoMigrate(models.All()...); err != nil {
			return fmt.Errorf("auto-migrate failed: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "AutoMigrate complete.")

		// Step 2: Run pending file-based migrations
		migrator, err := migrations.NewMigrator(db)
		if err != nil {
			return err
		}
		n, err := migrator.Run()
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
