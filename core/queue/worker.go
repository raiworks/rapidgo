package queue

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"
)

// WorkerConfig configures the worker pool.
type WorkerConfig struct {
	Queues       []string
	Concurrency  int
	PollInterval time.Duration
	MaxAttempts  uint
	RetryDelay   time.Duration
	Timeout      time.Duration
}

// Worker manages a pool of goroutines that process jobs.
type Worker struct {
	driver Driver
	config WorkerConfig
	log    *slog.Logger
}

// NewWorker creates a worker with the given driver and config.
func NewWorker(driver Driver, config WorkerConfig, log *slog.Logger) *Worker {
	if len(config.Queues) == 0 {
		config.Queues = []string{"default"}
	}
	if config.Concurrency < 1 {
		config.Concurrency = 1
	}
	if config.PollInterval == 0 {
		config.PollInterval = 3 * time.Second
	}
	if config.MaxAttempts == 0 {
		config.MaxAttempts = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 30 * time.Second
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}
	return &Worker{driver: driver, config: config, log: log}
}

// Run starts the worker pool and blocks until ctx is cancelled.
func (w *Worker) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	for i := 0; i < w.config.Concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			w.loop(ctx, workerID)
		}(i + 1)
	}

	wg.Wait()
	return nil
}

// loop is the main processing loop for a single worker goroutine.
func (w *Worker) loop(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		processed := false
		for _, q := range w.config.Queues {
			job, err := w.driver.Pop(ctx, q)
			if err != nil {
				w.log.Error("pop error", "worker", workerID, "queue", q, "error", err)
				break
			}
			if job == nil {
				continue
			}

			processed = true
			w.process(ctx, workerID, job)
		}

		if !processed {
			select {
			case <-ctx.Done():
				return
			case <-time.After(w.config.PollInterval):
			}
		}
	}
}

// process handles a single job with timeout, panic recovery, and retry logic.
func (w *Worker) process(ctx context.Context, workerID int, job *Job) {
	handler := ResolveHandler(job.Type)
	if handler == nil {
		w.log.Error("no handler for job type", "worker", workerID, "type", job.Type, "id", job.ID)
		_ = w.driver.Fail(ctx, job, fmt.Errorf("no handler registered for type %q", job.Type))
		return
	}

	maxAttempts := job.MaxAttempts
	if maxAttempts == 0 {
		maxAttempts = w.config.MaxAttempts
	}

	jobCtx, cancel := context.WithTimeout(ctx, w.config.Timeout)
	defer cancel()

	err := w.safeExecute(jobCtx, handler, job)

	if err == nil {
		_ = w.driver.Delete(ctx, job)
		return
	}

	w.log.Error("job failed", "worker", workerID, "id", job.ID, "type", job.Type,
		"attempt", job.Attempts+1, "max_attempts", maxAttempts, "error", err)

	if job.Attempts+1 < maxAttempts {
		_ = w.driver.Release(ctx, job, w.config.RetryDelay)
	} else {
		_ = w.driver.Fail(ctx, job, err)
	}
}

// safeExecute calls the handler with panic recovery.
func (w *Worker) safeExecute(ctx context.Context, handler HandlerFunc, job *Job) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			err = fmt.Errorf("panic: %v\n%s", r, stack)
		}
	}()
	return handler(ctx, job.Payload)
}
