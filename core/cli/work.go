package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/RAiWorks/RapidGo/app/jobs"
	"github.com/RAiWorks/RapidGo/app/providers"
	"github.com/RAiWorks/RapidGo/core/app"
	"github.com/RAiWorks/RapidGo/core/config"
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/queue"
	"github.com/spf13/cobra"
)

var workQueues string
var workWorkers int
var workTimeout int

var workCmd = &cobra.Command{
	Use:   "work",
	Short: "Start the queue worker to process background jobs",
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

		// Register application job handlers.
		jobs.RegisterJobs()

		// Resolve dispatcher from container.
		dispatcher := container.MustMake[*queue.Dispatcher](application.Container, "queue")

		// Build worker config from flags and env vars.
		queues := strings.Split(workQueues, ",")
		if workQueues == "" {
			queues = []string{config.Env("QUEUE_DEFAULT", "default")}
		}

		concurrency := workWorkers
		if concurrency == 0 {
			concurrency = config.EnvInt("QUEUE_WORKERS", 1)
		}

		timeout := time.Duration(workTimeout) * time.Second
		if workTimeout == 0 {
			timeout = time.Duration(config.EnvInt("QUEUE_TIMEOUT", 60)) * time.Second
		}

		wkr := queue.NewWorker(dispatcher.Driver(), queue.WorkerConfig{
			Queues:       queues,
			Concurrency:  concurrency,
			PollInterval: time.Duration(config.EnvInt("QUEUE_POLL_INTERVAL", 3)) * time.Second,
			MaxAttempts:  uint(config.EnvInt("QUEUE_MAX_ATTEMPTS", 3)),
			RetryDelay:   time.Duration(config.EnvInt("QUEUE_RETRY_DELAY", 30)) * time.Second,
			Timeout:      timeout,
		}, slog.Default())

		driverName := config.Env("QUEUE_DRIVER", "database")

		fmt.Println("=================================")
		fmt.Println("  RapidGo Queue Worker")
		fmt.Println("=================================")
		fmt.Printf("  Driver:  %s\n", driverName)
		fmt.Printf("  Queues:  %s\n", strings.Join(queues, ", "))
		fmt.Printf("  Workers: %d\n", concurrency)
		fmt.Printf("  Timeout: %s\n", timeout)
		fmt.Println("=================================")

		slog.Info("queue worker starting",
			"driver", driverName,
			"queues", queues,
			"workers", concurrency,
			"timeout", timeout,
		)

		// Graceful shutdown on SIGINT/SIGTERM.
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		return wkr.Run(ctx)
	},
}

func init() {
	workCmd.Flags().StringVarP(&workQueues, "queues", "q", "", "comma-separated queue names (default: from QUEUE_DEFAULT)")
	workCmd.Flags().IntVarP(&workWorkers, "workers", "w", 0, "number of concurrent workers (default: 1)")
	workCmd.Flags().IntVar(&workTimeout, "timeout", 0, "max job processing time in seconds (default: from QUEUE_TIMEOUT)")
}
