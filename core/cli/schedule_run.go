package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/raiworks/rapidgo/v2/core/scheduler"
	"github.com/raiworks/rapidgo/v2/core/service"
	"github.com/spf13/cobra"
)

var scheduleRunCmd = &cobra.Command{
	Use:   "schedule:run",
	Short: "Start the task scheduler to run registered tasks on their cron schedules",
	RunE: func(cmd *cobra.Command, args []string) error {
		application := NewApp(service.ModeAll)

		// Create scheduler.
		s := scheduler.New(slog.Default())

		// Register application-defined tasks via callback.
		if scheduleRegistrar != nil {
			scheduleRegistrar(s, application)
		}

		// Print banner with registered tasks.
		fmt.Println("=================================")
		fmt.Println("  RapidGo Task Scheduler")
		fmt.Println("=================================")
		for _, t := range s.Tasks() {
			fmt.Printf("  [%s] %s\n", t.Schedule, t.Name)
		}
		fmt.Println("=================================")

		slog.Info("task scheduler starting", "tasks", len(s.Tasks()))

		// Graceful shutdown on SIGINT/SIGTERM.
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		return s.Run(ctx)
	},
}
