package cli

import (
	"fmt"

	"github.com/raiworks/rapidgo/v2/core/config"
	"github.com/raiworks/rapidgo/v2/core/container"
	"github.com/raiworks/rapidgo/v2/core/service"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var dbWipeCmd = &cobra.Command{
	Use:   "db:wipe",
	Short: "Truncate all database tables (preserves migration tracking)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.AppEnv() == "production" {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				return fmt.Errorf("refusing to wipe database in production without --force")
			}
		}

		application := NewApp(service.ModeAll)
		db := container.MustMake[*gorm.DB](application.Container, "db")

		tables, err := db.Migrator().GetTables()
		if err != nil {
			return fmt.Errorf("failed to get tables: %w", err)
		}
		for _, table := range tables {
			if table == "migrations" {
				continue // preserve migration tracking
			}
			if err := truncateTable(db, table); err != nil {
				return fmt.Errorf("failed to truncate table %s: %w", table, err)
			}
		}
		fmt.Fprintln(cmd.OutOrStdout(), "All tables truncated.")
		return nil
	},
}

func init() {
	dbWipeCmd.Flags().Bool("force", false, "Force the operation in production")
}

// truncateTable truncates a table using the appropriate SQL for the dialect.
// Table names come from GORM's Migrator().GetTables() — trusted DB metadata.
func truncateTable(db *gorm.DB, table string) error {
	switch db.Dialector.Name() {
	case "postgres":
		return db.Exec(fmt.Sprintf(`TRUNCATE TABLE %q CASCADE`, table)).Error
	case "mysql":
		return db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", table)).Error
	default: // sqlite
		return db.Exec(fmt.Sprintf(`DELETE FROM "%s"`, table)).Error
	}
}
