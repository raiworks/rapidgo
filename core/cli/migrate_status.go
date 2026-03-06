package cli

import (
	"fmt"

	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/database/migrations"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var migrateStatusCmd = &cobra.Command{
	Use:   "migrate:status",
	Short: "Show the status of all migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		application := NewApp()
		db := container.MustMake[*gorm.DB](application.Container, "db")

		migrator, err := migrations.NewMigrator(db)
		if err != nil {
			return err
		}
		statuses, err := migrator.Status()
		if err != nil {
			return err
		}
		if len(statuses) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No migrations registered.")
			return nil
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Migration                | Status  | Batch")
		fmt.Fprintln(cmd.OutOrStdout(), "-------------------------+---------+------")
		for _, s := range statuses {
			status := "Pending"
			batch := ""
			if s.Applied {
				status = "Applied"
				batch = fmt.Sprintf("%d", s.Batch)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%-25s| %-8s| %s\n", s.Version, status, batch)
		}
		return nil
	},
}
