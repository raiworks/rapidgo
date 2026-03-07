package scheduler

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/robfig/cron/v3"
)

// TaskFunc is a function executed on schedule.
type TaskFunc func(ctx context.Context) error

// Task is a named scheduled task with a cron expression.
type Task struct {
	Name     string
	Schedule string
	Run      TaskFunc
}

// Scheduler wraps robfig/cron and manages named tasks.
type Scheduler struct {
	cron  *cron.Cron
	tasks []Task
	log   *slog.Logger
}

// New creates a scheduler with the standard 5-field cron parser plus descriptors (@every, @daily, etc.).
func New(log *slog.Logger) *Scheduler {
	if log == nil {
		log = slog.Default()
	}
	return &Scheduler{
		cron: cron.New(cron.WithParser(cron.NewParser(
			cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		))),
		log: log,
	}
}

// Add registers a named task with a cron expression.
// Returns an error if the cron expression is invalid.
func (s *Scheduler) Add(schedule, name string, fn TaskFunc) error {
	_, err := s.cron.AddFunc(schedule, s.wrap(name, fn))
	if err != nil {
		return err
	}
	s.tasks = append(s.tasks, Task{
		Name:     name,
		Schedule: schedule,
		Run:      fn,
	})
	return nil
}

// Tasks returns all registered tasks.
func (s *Scheduler) Tasks() []Task {
	return s.tasks
}

// Run starts the cron engine and blocks until ctx is cancelled.
// On shutdown it waits for any running tasks to complete before returning.
func (s *Scheduler) Run(ctx context.Context) error {
	s.cron.Start()
	<-ctx.Done()
	stopCtx := s.cron.Stop()
	<-stopCtx.Done()
	return nil
}

// wrap returns a cron-compatible func() that executes the task with
// structured logging, duration measurement, and panic recovery.
func (s *Scheduler) wrap(name string, fn TaskFunc) func() {
	return func() {
		start := time.Now()
		s.log.Info("task started", "task", name)

		defer func() {
			if r := recover(); r != nil {
				s.log.Error("task panicked", "task", name, "panic", r,
					"stack", string(debug.Stack()))
			}
		}()

		if err := fn(context.Background()); err != nil {
			s.log.Error("task failed", "task", name, "error", err,
				"duration", time.Since(start))
			return
		}

		s.log.Info("task completed", "task", name,
			"duration", time.Since(start))
	}
}
