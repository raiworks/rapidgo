package cli

import (
	"fmt"

	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/database/migrations"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var migrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last batch of migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		application := NewApp()
		db := container.MustMake[*gorm.DB](application.Container, "db")

		migrator, err := migrations.NewMigrator(db)
		if err != nil {
			return err
		}
		n, err := migrator.Rollback()
		if err != nil {
			return fmt.Errorf("rollback failed: %w", err)
		}
		if n == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Nothing to rollback.")
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Rolled back %d migration(s).\n", n)
		}
		return nil
	},
}
