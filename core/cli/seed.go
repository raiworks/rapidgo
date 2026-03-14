package cli

import (
	"fmt"

	"github.com/raiworks/rapidgo/v2/core/container"
	"github.com/raiworks/rapidgo/v2/core/service"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var dbSeedCmd = &cobra.Command{
	Use:   "db:seed",
	Short: "Seed the database with records",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle --list flag (no DB connection needed)
		if list, _ := cmd.Flags().GetBool("list"); list {
			if seederListFn == nil {
				fmt.Fprintln(cmd.OutOrStdout(), "No seeder list registered. Call cli.SetSeederList() to enable.")
				return nil
			}
			names := seederListFn()
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No seeders registered.")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Available seeders:")
			for _, name := range names {
				fmt.Fprintln(cmd.OutOrStdout(), "  - "+name)
			}
			return nil
		}

		application := NewApp(service.ModeAll)
		db := container.MustMake[*gorm.DB](application.Container, "db")

		if seederFn == nil {
			return fmt.Errorf("no seeder registered — call cli.SetSeeder() in main.go")
		}

		name, _ := cmd.Flags().GetString("seeder")
		if err := seederFn(db, name); err != nil {
			return err
		}
		if name != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Seeder %s complete.\n", name)
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "Database seeding complete.")
		}
		return nil
	},
}

func init() {
	dbSeedCmd.Flags().String("seeder", "", "Run a specific seeder by name")
	dbSeedCmd.Flags().Bool("list", false, "List available seeders")
}
