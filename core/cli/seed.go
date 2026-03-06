package cli

import (
	"fmt"

	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/database/seeders"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var dbSeedCmd = &cobra.Command{
	Use:   "db:seed",
	Short: "Seed the database with records",
	RunE: func(cmd *cobra.Command, args []string) error {
		application := NewApp()
		db := container.MustMake[*gorm.DB](application.Container, "db")

		name, _ := cmd.Flags().GetString("seeder")
		if name != "" {
			if err := seeders.RunByName(db, name); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Seeder %s complete.\n", name)
			return nil
		}

		if err := seeders.RunAll(db); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Database seeding complete.")
		return nil
	},
}

func init() {
	dbSeedCmd.Flags().String("seeder", "", "Run a specific seeder by name")
}
