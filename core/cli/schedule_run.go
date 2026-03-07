package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/RAiWorks/RapidGo/app/providers"
	"github.com/RAiWorks/RapidGo/app/schedule"
	"github.com/RAiWorks/RapidGo/core/app"
	"github.com/RAiWorks/RapidGo/core/config"
	"github.com/RAiWorks/RapidGo/core/scheduler"
	"github.com/spf13/cobra"
)

var scheduleRunCmd = &cobra.Command{
	Use:   "schedule:run",
	Short: "Start the task scheduler to run registered tasks on their cron schedules",
	RunE: func(cmd *cobra.Command, args []string) error {
		config.Load()

		// Minimal bootstrap — no HTTP providers needed.
		application := app.New()
		application.Register(&providers.ConfigProvider{})
		application.Register(&providers.LoggerProvider{})
		application.Register(&providers.DatabaseProvider{})
		application.Register(&providers.RedisProvider{})
		application.Register(&providers.QueueProvider{})
		application.Boot()

		// Create scheduler.
		s := scheduler.New(slog.Default())

		// Register application-defined tasks.
		schedule.RegisterSchedule(s, application)

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
