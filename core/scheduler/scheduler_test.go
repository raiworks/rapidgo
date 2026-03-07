package scheduler

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// newTestLogger creates a logger that writes to a buffer for assertion.
func newTestLogger() (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	h := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	return slog.New(h), &buf
}

// --- Construction ---

func TestNewScheduler(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)
	if s == nil {
		t.Fatal("New() returned nil")
	}
	if len(s.Tasks()) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(s.Tasks()))
	}
}

func TestNewSchedulerNilLogger(t *testing.T) {
	s := New(nil)
	if s == nil {
		t.Fatal("New(nil) returned nil")
	}
}

// --- Task Registration ---

func TestAddTask(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)

	err := s.Add("*/5 * * * *", "cleanup", func(ctx context.Context) error { return nil })
	if err != nil {
		t.Fatalf("Add() returned error: %v", err)
	}
	tasks := s.Tasks()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Name != "cleanup" {
		t.Errorf("expected name 'cleanup', got %q", tasks[0].Name)
	}
	if tasks[0].Schedule != "*/5 * * * *" {
		t.Errorf("expected schedule '*/5 * * * *', got %q", tasks[0].Schedule)
	}
}

func TestAddMultipleTasks(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)

	names := []string{"task-a", "task-b", "task-c"}
	for _, name := range names {
		if err := s.Add("@every 1m", name, func(ctx context.Context) error { return nil }); err != nil {
			t.Fatalf("Add(%s) error: %v", name, err)
		}
	}
	if len(s.Tasks()) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(s.Tasks()))
	}
	for i, task := range s.Tasks() {
		if task.Name != names[i] {
			t.Errorf("task %d: expected name %q, got %q", i, names[i], task.Name)
		}
	}
}

func TestAddInvalidCronExpression(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)

	err := s.Add("not-a-cron", "bad", func(ctx context.Context) error { return nil })
	if err == nil {
		t.Fatal("expected error for invalid cron expression, got nil")
	}
}

func TestAddEmptyName(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)

	err := s.Add("@every 1m", "", func(ctx context.Context) error { return nil })
	if err != nil {
		t.Fatalf("Add with empty name should succeed, got: %v", err)
	}
}

func TestAddDescriptors(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)

	for _, expr := range []string{"@every 5m", "@daily", "@hourly", "@weekly"} {
		if err := s.Add(expr, "desc-"+expr, func(ctx context.Context) error { return nil }); err != nil {
			t.Errorf("Add(%q) should succeed, got: %v", expr, err)
		}
	}
}

// --- Task Execution ---

func TestTaskExecutes(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)

	var counter atomic.Int32
	if err := s.Add("@every 1s", "counter", func(ctx context.Context) error {
		counter.Add(1)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return within timeout")
	}

	if counter.Load() < 1 {
		t.Errorf("expected counter >= 1, got %d", counter.Load())
	}
}

func TestMultipleTasksExecute(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)

	var counterA, counterB atomic.Int32
	s.Add("@every 1s", "task-a", func(ctx context.Context) error { counterA.Add(1); return nil })
	s.Add("@every 1s", "task-b", func(ctx context.Context) error { counterB.Add(1); return nil })

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run() error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() timeout")
	}

	if counterA.Load() < 1 {
		t.Errorf("counterA expected >= 1, got %d", counterA.Load())
	}
	if counterB.Load() < 1 {
		t.Errorf("counterB expected >= 1, got %d", counterB.Load())
	}
}

func TestTaskReceivesContext(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)

	var gotCtx atomic.Int32
	s.Add("@every 1s", "ctx-check", func(ctx context.Context) error {
		if ctx != nil {
			gotCtx.Store(1)
		}
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}

	if gotCtx.Load() != 1 {
		t.Error("task did not receive a non-nil context")
	}
}

// --- Error Handling ---

func TestTaskErrorDoesNotStopScheduler(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)

	var counter atomic.Int32
	s.Add("@every 1s", "error-task", func(ctx context.Context) error {
		counter.Add(1)
		return errors.New("task error")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}

	if counter.Load() < 2 {
		t.Errorf("expected error task to fire at least twice, got %d", counter.Load())
	}
}

func TestPanicRecovery(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)

	var safeCounter atomic.Int32

	s.Add("@every 1s", "panic-task", func(ctx context.Context) error {
		panic("test panic")
	})
	s.Add("@every 1s", "safe-task", func(ctx context.Context) error {
		safeCounter.Add(1)
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run() error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}

	if safeCounter.Load() < 1 {
		t.Error("safe task should have executed despite panic in other task")
	}
}

// --- Graceful Shutdown ---

func TestRunReturnsOnContextCancel(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)

	s.Add("@every 1h", "slow", func(ctx context.Context) error { return nil })

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	// Cancel immediately.
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run() did not return within 2s after cancel")
	}
}

func TestRunWaitsForRunningTask(t *testing.T) {
	log, _ := newTestLogger()
	s := New(log)

	var completed atomic.Int32
	s.Add("@every 1s", "long-task", func(ctx context.Context) error {
		time.Sleep(500 * time.Millisecond)
		completed.Store(1)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	// Wait for the task to start, then cancel.
	time.Sleep(1200 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for Run to return")
	}

	if completed.Load() != 1 {
		t.Error("expected running task to complete before shutdown")
	}
}

// --- Logging ---

func TestLogsTaskStarted(t *testing.T) {
	log, buf := newTestLogger()
	s := New(log)

	s.Add("@every 1s", "log-start", func(ctx context.Context) error { return nil })

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}

	if !strings.Contains(buf.String(), "task started") {
		t.Errorf("expected 'task started' in log, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "log-start") {
		t.Errorf("expected task name 'log-start' in log, got: %s", buf.String())
	}
}

func TestLogsTaskCompleted(t *testing.T) {
	log, buf := newTestLogger()
	s := New(log)

	s.Add("@every 1s", "log-done", func(ctx context.Context) error { return nil })

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}

	output := buf.String()
	if !strings.Contains(output, "task completed") {
		t.Errorf("expected 'task completed' in log, got: %s", output)
	}
	if !strings.Contains(output, "duration") {
		t.Errorf("expected 'duration' in log, got: %s", output)
	}
}

func TestLogsTaskFailed(t *testing.T) {
	log, buf := newTestLogger()
	s := New(log)

	s.Add("@every 1s", "log-fail", func(ctx context.Context) error {
		return errors.New("something went wrong")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}

	output := buf.String()
	if !strings.Contains(output, "task failed") {
		t.Errorf("expected 'task failed' in log, got: %s", output)
	}
	if !strings.Contains(output, "something went wrong") {
		t.Errorf("expected error message in log, got: %s", output)
	}
}

func TestLogsTaskPanicked(t *testing.T) {
	log, buf := newTestLogger()
	s := New(log)

	s.Add("@every 1s", "log-panic", func(ctx context.Context) error {
		panic("boom")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}

	output := buf.String()
	if !strings.Contains(output, "task panicked") {
		t.Errorf("expected 'task panicked' in log, got: %s", output)
	}
	if !strings.Contains(output, "boom") {
		t.Errorf("expected panic value 'boom' in log, got: %s", output)
	}
}
