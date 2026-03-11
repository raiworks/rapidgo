package queue

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ─── Handler Registry ─────────────────────────────────────────────────────────

func TestRegisterAndResolveHandler(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	called := false
	RegisterHandler("email", func(_ context.Context, _ json.RawMessage) error {
		called = true
		return nil
	})

	h := ResolveHandler("email")
	if h == nil {
		t.Fatal("expected handler, got nil")
	}
	_ = h(context.Background(), nil)
	if !called {
		t.Error("handler was not called")
	}
}

func TestResolveUnknownHandler(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	if h := ResolveHandler("nonexistent"); h != nil {
		t.Error("expected nil for unknown handler")
	}
}

func TestResetHandlers(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	RegisterHandler("test", func(_ context.Context, _ json.RawMessage) error { return nil })
	ResetHandlers()

	if h := ResolveHandler("test"); h != nil {
		t.Error("expected nil after reset")
	}
}

// ─── Dispatcher ───────────────────────────────────────────────────────────────

func TestDispatcherDispatch(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	RegisterHandler("test_job", func(_ context.Context, _ json.RawMessage) error { return nil })

	mem := NewMemoryDriver()
	d := NewDispatcher(mem)

	err := d.Dispatch(context.Background(), "default", "test_job", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}

	job, _ := mem.Pop(context.Background(), "default")
	if job == nil {
		t.Fatal("expected job, got nil")
	}
	if job.Queue != "default" {
		t.Errorf("queue = %q, want %q", job.Queue, "default")
	}
	if job.Type != "test_job" {
		t.Errorf("type = %q, want %q", job.Type, "test_job")
	}
	if job.Attempts != 0 {
		t.Errorf("attempts = %d, want 0", job.Attempts)
	}

	var payload map[string]string
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if payload["key"] != "value" {
		t.Errorf("payload[key] = %q, want %q", payload["key"], "value")
	}
}

func TestDispatcherDispatchUnknownType(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	d := NewDispatcher(NewMemoryDriver())
	err := d.Dispatch(context.Background(), "default", "unknown", nil)
	if err == nil {
		t.Error("expected error for unknown handler")
	}
}

func TestDispatcherDispatchDelayed(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	RegisterHandler("delayed", func(_ context.Context, _ json.RawMessage) error { return nil })

	mem := NewMemoryDriver()
	d := NewDispatcher(mem)

	err := d.DispatchDelayed(context.Background(), "default", "delayed", "payload", 5*time.Minute)
	if err != nil {
		t.Fatalf("dispatch delayed failed: %v", err)
	}

	// Should not be poppable yet (available in future).
	job, _ := mem.Pop(context.Background(), "default")
	if job != nil {
		t.Error("delayed job should not be available yet")
	}
}

// ─── MemoryDriver ─────────────────────────────────────────────────────────────

func TestMemoryDriverLifecycle(t *testing.T) {
	d := NewMemoryDriver()
	ctx := context.Background()

	job := &Job{Queue: "q", Type: "t", Payload: json.RawMessage(`{}`), AvailableAt: time.Now(), CreatedAt: time.Now()}
	if err := d.Push(ctx, job); err != nil {
		t.Fatalf("push: %v", err)
	}

	size, _ := d.Size(ctx, "q")
	if size != 1 {
		t.Errorf("size = %d, want 1", size)
	}

	popped, err := d.Pop(ctx, "q")
	if err != nil {
		t.Fatalf("pop: %v", err)
	}
	if popped == nil {
		t.Fatal("expected job, got nil")
	}
	if popped.ReservedAt == nil {
		t.Error("expected ReservedAt to be set")
	}

	if err := d.Delete(ctx, popped); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestMemoryDriverRelease(t *testing.T) {
	d := NewMemoryDriver()
	ctx := context.Background()

	job := &Job{Queue: "q", Type: "t", Payload: json.RawMessage(`{}`), AvailableAt: time.Now(), CreatedAt: time.Now()}
	_ = d.Push(ctx, job)
	popped, _ := d.Pop(ctx, "q")

	if err := d.Release(ctx, popped, 0); err != nil {
		t.Fatalf("release: %v", err)
	}

	again, _ := d.Pop(ctx, "q")
	if again == nil {
		t.Fatal("released job should be poppable again")
	}
	if again.Attempts != 1 {
		t.Errorf("attempts = %d, want 1", again.Attempts)
	}
}

func TestMemoryDriverFail(t *testing.T) {
	d := NewMemoryDriver()
	ctx := context.Background()

	job := &Job{Queue: "q", Type: "t", Payload: json.RawMessage(`{}`), AvailableAt: time.Now(), CreatedAt: time.Now()}
	_ = d.Push(ctx, job)
	popped, _ := d.Pop(ctx, "q")

	if err := d.Fail(ctx, popped, errors.New("boom")); err != nil {
		t.Fatalf("fail: %v", err)
	}

	size, _ := d.Size(ctx, "q")
	if size != 0 {
		t.Errorf("size after fail = %d, want 0", size)
	}
}

func TestMemoryDriverDelayedJob(t *testing.T) {
	d := NewMemoryDriver()
	ctx := context.Background()

	job := &Job{Queue: "q", Type: "t", Payload: json.RawMessage(`{}`), AvailableAt: time.Now().Add(time.Hour), CreatedAt: time.Now()}
	_ = d.Push(ctx, job)

	popped, _ := d.Pop(ctx, "q")
	if popped != nil {
		t.Error("delayed job should not be poppable yet")
	}

	size, _ := d.Size(ctx, "q")
	if size != 1 {
		t.Errorf("size = %d, want 1 (still pending)", size)
	}
}

// ─── SyncDriver ───────────────────────────────────────────────────────────────

func TestSyncDriverImmediateExecution(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	called := false
	RegisterHandler("sync_test", func(_ context.Context, _ json.RawMessage) error {
		called = true
		return nil
	})

	d := NewSyncDriver()
	job := &Job{Type: "sync_test", Payload: json.RawMessage(`{}`)}
	if err := d.Push(context.Background(), job); err != nil {
		t.Fatalf("push: %v", err)
	}
	if !called {
		t.Error("handler should have been called immediately")
	}
}

func TestSyncDriverReturnsHandlerError(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	RegisterHandler("sync_err", func(_ context.Context, _ json.RawMessage) error {
		return errors.New("handler error")
	})

	d := NewSyncDriver()
	job := &Job{Type: "sync_err", Payload: json.RawMessage(`{}`)}
	err := d.Push(context.Background(), job)
	if err == nil || err.Error() != "handler error" {
		t.Errorf("expected 'handler error', got %v", err)
	}
}

func TestSyncDriverUnknownHandler(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	d := NewSyncDriver()
	job := &Job{Type: "unknown_sync", Payload: json.RawMessage(`{}`)}
	err := d.Push(context.Background(), job)
	if err == nil {
		t.Error("expected error for unknown handler")
	}
}

// ─── DatabaseDriver ───────────────────────────────────────────────────────────

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	// Create tables for testing.
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		queue TEXT NOT NULL,
		type TEXT NOT NULL,
		payload TEXT NOT NULL,
		attempts INTEGER NOT NULL DEFAULT 0,
		max_attempts INTEGER NOT NULL DEFAULT 3,
		available_at DATETIME NOT NULL,
		reserved_at DATETIME,
		created_at DATETIME NOT NULL
	)`).Error; err != nil {
		t.Fatalf("create jobs table: %v", err)
	}

	if err := db.Exec(`CREATE TABLE IF NOT EXISTS failed_jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		queue TEXT NOT NULL,
		type TEXT NOT NULL,
		payload TEXT NOT NULL,
		error TEXT NOT NULL,
		failed_at DATETIME NOT NULL
	)`).Error; err != nil {
		t.Fatalf("create failed_jobs table: %v", err)
	}

	return db
}

func TestDatabaseDriverLifecycle(t *testing.T) {
	db := setupTestDB(t)
	d := NewDatabaseDriver(db, "jobs", "failed_jobs")
	ctx := context.Background()

	job := &Job{Queue: "default", Type: "test", Payload: json.RawMessage(`{"a":1}`), AvailableAt: time.Now(), CreatedAt: time.Now(), MaxAttempts: 3}
	if err := d.Push(ctx, job); err != nil {
		t.Fatalf("push: %v", err)
	}
	if job.ID == 0 {
		t.Error("expected assigned ID")
	}

	popped, err := d.Pop(ctx, "default")
	if err != nil {
		t.Fatalf("pop: %v", err)
	}
	if popped == nil {
		t.Fatal("expected job, got nil")
	}
	if popped.Type != "test" {
		t.Errorf("type = %q, want %q", popped.Type, "test")
	}

	if err := d.Delete(ctx, popped); err != nil {
		t.Fatalf("delete: %v", err)
	}

	// Verify deleted.
	size, _ := d.Size(ctx, "default")
	if size != 0 {
		t.Errorf("size after delete = %d, want 0", size)
	}
}

func TestDatabaseDriverRelease(t *testing.T) {
	db := setupTestDB(t)
	d := NewDatabaseDriver(db, "jobs", "failed_jobs")
	ctx := context.Background()

	job := &Job{Queue: "default", Type: "test", Payload: json.RawMessage(`{}`), AvailableAt: time.Now(), CreatedAt: time.Now(), MaxAttempts: 3}
	_ = d.Push(ctx, job)
	popped, _ := d.Pop(ctx, "default")

	if err := d.Release(ctx, popped, 0); err != nil {
		t.Fatalf("release: %v", err)
	}

	// Should be available again.
	again, _ := d.Pop(ctx, "default")
	if again == nil {
		t.Fatal("released job should be poppable")
	}
	if again.Attempts != 1 {
		t.Errorf("attempts = %d, want 1", again.Attempts)
	}
}

func TestDatabaseDriverFail(t *testing.T) {
	db := setupTestDB(t)
	d := NewDatabaseDriver(db, "jobs", "failed_jobs")
	ctx := context.Background()

	job := &Job{Queue: "default", Type: "test", Payload: json.RawMessage(`{}`), AvailableAt: time.Now(), CreatedAt: time.Now(), MaxAttempts: 3}
	_ = d.Push(ctx, job)
	popped, _ := d.Pop(ctx, "default")

	if err := d.Fail(ctx, popped, errors.New("job failed")); err != nil {
		t.Fatalf("fail: %v", err)
	}

	// Verify removed from jobs.
	size, _ := d.Size(ctx, "default")
	if size != 0 {
		t.Errorf("jobs size = %d, want 0", size)
	}

	// Verify in failed_jobs.
	var count int64
	db.Table("failed_jobs").Count(&count)
	if count != 1 {
		t.Errorf("failed_jobs count = %d, want 1", count)
	}
}

// ─── RedisDriver ──────────────────────────────────────────────────────────────

func setupMiniRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return client, mr
}

func TestRedisDriverLifecycle(t *testing.T) {
	client, _ := setupMiniRedis(t)
	d := NewRedisDriver(client)
	ctx := context.Background()

	job := &Job{Queue: "default", Type: "test", Payload: json.RawMessage(`{"x":1}`), AvailableAt: time.Now(), CreatedAt: time.Now()}
	if err := d.Push(ctx, job); err != nil {
		t.Fatalf("push: %v", err)
	}

	size, _ := d.Size(ctx, "default")
	if size != 1 {
		t.Errorf("size = %d, want 1", size)
	}

	popped, err := d.Pop(ctx, "default")
	if err != nil {
		t.Fatalf("pop: %v", err)
	}
	if popped == nil {
		t.Fatal("expected job, got nil")
	}
	if popped.Type != "test" {
		t.Errorf("type = %q, want %q", popped.Type, "test")
	}

	size, _ = d.Size(ctx, "default")
	if size != 0 {
		t.Errorf("size after pop = %d, want 0", size)
	}
}

func TestRedisDriverRelease(t *testing.T) {
	client, _ := setupMiniRedis(t)
	d := NewRedisDriver(client)
	ctx := context.Background()

	job := &Job{Queue: "default", Type: "test", Payload: json.RawMessage(`{}`), AvailableAt: time.Now(), CreatedAt: time.Now()}
	_ = d.Push(ctx, job)
	popped, _ := d.Pop(ctx, "default")

	if err := d.Release(ctx, popped, 0); err != nil {
		t.Fatalf("release: %v", err)
	}

	again, _ := d.Pop(ctx, "default")
	if again == nil {
		t.Fatal("released job should be poppable")
	}
	if again.Attempts != 1 {
		t.Errorf("attempts = %d, want 1", again.Attempts)
	}
}

func TestRedisDriverFail(t *testing.T) {
	client, _ := setupMiniRedis(t)
	d := NewRedisDriver(client)
	ctx := context.Background()

	job := &Job{Queue: "default", Type: "test", Payload: json.RawMessage(`{}`), AvailableAt: time.Now(), CreatedAt: time.Now()}
	_ = d.Push(ctx, job)
	popped, _ := d.Pop(ctx, "default")

	if err := d.Fail(ctx, popped, errors.New("redis fail")); err != nil {
		t.Fatalf("fail: %v", err)
	}

	// Verify in failed list.
	failedLen := client.LLen(ctx, d.failedKey("default")).Val()
	if failedLen != 1 {
		t.Errorf("failed list len = %d, want 1", failedLen)
	}
}

func TestRedisDriverSize(t *testing.T) {
	client, _ := setupMiniRedis(t)
	d := NewRedisDriver(client)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_ = d.Push(ctx, &Job{Queue: "q", Type: "t", Payload: json.RawMessage(`{}`), AvailableAt: time.Now(), CreatedAt: time.Now()})
	}

	size, _ := d.Size(ctx, "q")
	if size != 3 {
		t.Errorf("size = %d, want 3", size)
	}

	_, _ = d.Pop(ctx, "q")
	size, _ = d.Size(ctx, "q")
	if size != 2 {
		t.Errorf("size after pop = %d, want 2", size)
	}
}

// ─── Worker ───────────────────────────────────────────────────────────────────

func testLogger() *slog.Logger {
	return slog.Default()
}

func TestWorkerProcessesJobs(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	var processed atomic.Int32
	RegisterHandler("work_test", func(_ context.Context, _ json.RawMessage) error {
		processed.Add(1)
		return nil
	})

	mem := NewMemoryDriver()
	disp := NewDispatcher(mem)
	for i := 0; i < 3; i++ {
		_ = disp.Dispatch(context.Background(), "default", "work_test", "hello")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	w := NewWorker(mem, WorkerConfig{
		Queues:       []string{"default"},
		Concurrency:  1,
		PollInterval: 50 * time.Millisecond,
		MaxAttempts:  3,
		RetryDelay:   0,
		Timeout:      5 * time.Second,
	}, testLogger())

	_ = w.Run(ctx)

	if processed.Load() != 3 {
		t.Errorf("processed = %d, want 3", processed.Load())
	}
}

func TestWorkerRetryOnFailure(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	var calls atomic.Int32
	RegisterHandler("retry_test", func(_ context.Context, _ json.RawMessage) error {
		n := calls.Add(1)
		if n < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	mem := NewMemoryDriver()
	disp := NewDispatcher(mem)
	_ = disp.Dispatch(context.Background(), "default", "retry_test", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w := NewWorker(mem, WorkerConfig{
		Queues:       []string{"default"},
		Concurrency:  1,
		PollInterval: 50 * time.Millisecond,
		MaxAttempts:  5,
		RetryDelay:   time.Millisecond,
		Timeout:      5 * time.Second,
	}, testLogger())

	_ = w.Run(ctx)

	if calls.Load() != 3 {
		t.Errorf("calls = %d, want 3", calls.Load())
	}
}

func TestWorkerFailAfterMaxAttempts(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	RegisterHandler("max_fail", func(_ context.Context, _ json.RawMessage) error {
		return errors.New("always fail")
	})

	mem := NewMemoryDriver()
	// Push directly with MaxAttempts=0 so worker config (2) is used.
	_ = mem.Push(context.Background(), &Job{
		Queue: "default", Type: "max_fail", Payload: json.RawMessage(`{}`),
		MaxAttempts: 0, AvailableAt: time.Now(), CreatedAt: time.Now(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w := NewWorker(mem, WorkerConfig{
		Queues:       []string{"default"},
		Concurrency:  1,
		PollInterval: 50 * time.Millisecond,
		MaxAttempts:  2,
		RetryDelay:   time.Millisecond,
		Timeout:      5 * time.Second,
	}, testLogger())

	_ = w.Run(ctx)

	// Should be in failed storage.
	if len(mem.failed) != 1 {
		t.Errorf("failed count = %d, want 1", len(mem.failed))
	}
}

func TestWorkerPanicRecovery(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	RegisterHandler("panic_job", func(_ context.Context, _ json.RawMessage) error {
		panic("job exploded")
	})

	mem := NewMemoryDriver()
	// Push directly with MaxAttempts=0 so worker config (1) is used.
	_ = mem.Push(context.Background(), &Job{
		Queue: "default", Type: "panic_job", Payload: json.RawMessage(`{}`),
		MaxAttempts: 0, AvailableAt: time.Now(), CreatedAt: time.Now(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	w := NewWorker(mem, WorkerConfig{
		Queues:       []string{"default"},
		Concurrency:  1,
		PollInterval: 50 * time.Millisecond,
		MaxAttempts:  1,
		RetryDelay:   time.Millisecond,
		Timeout:      5 * time.Second,
	}, testLogger())

	// Should not panic.
	_ = w.Run(ctx)

	if len(mem.failed) != 1 {
		t.Errorf("failed count = %d, want 1 (panic should be caught)", len(mem.failed))
	}
}

func TestWorkerGracefulShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	w := NewWorker(NewMemoryDriver(), WorkerConfig{
		Queues:       []string{"default"},
		Concurrency:  2,
		PollInterval: 50 * time.Millisecond,
		Timeout:      5 * time.Second,
	}, testLogger())

	done := make(chan error, 1)
	go func() {
		done <- w.Run(ctx)
	}()

	// Cancel immediately.
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("worker did not shut down in time")
	}
}

func TestWorkerJobTimeout(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	RegisterHandler("slow_job", func(ctx context.Context, _ json.RawMessage) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Second):
			return nil
		}
	})

	mem := NewMemoryDriver()
	// Push directly with MaxAttempts=0 so worker config (1) is used.
	_ = mem.Push(context.Background(), &Job{
		Queue: "default", Type: "slow_job", Payload: json.RawMessage(`{}`),
		MaxAttempts: 0, AvailableAt: time.Now(), CreatedAt: time.Now(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w := NewWorker(mem, WorkerConfig{
		Queues:       []string{"default"},
		Concurrency:  1,
		PollInterval: 50 * time.Millisecond,
		MaxAttempts:  1,
		RetryDelay:   time.Millisecond,
		Timeout:      100 * time.Millisecond,
	}, testLogger())

	_ = w.Run(ctx)

	if len(mem.failed) != 1 {
		t.Errorf("failed count = %d, want 1 (timeout should cause failure)", len(mem.failed))
	}
}

// ─── DispatchWithBackoff ──────────────────────────────────────────────────────

func TestDispatchWithBackoff(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	RegisterHandler("backoff_test", func(_ context.Context, _ json.RawMessage) error {
		return nil
	})

	mem := NewMemoryDriver()
	disp := NewDispatcher(mem)
	backoff := []uint{5, 30, 120}
	err := disp.DispatchWithBackoff(context.Background(), "default", "backoff_test", "data", backoff)
	if err != nil {
		t.Fatalf("DispatchWithBackoff error: %v", err)
	}

	job, _ := mem.Pop(context.Background(), "default")
	if job == nil {
		t.Fatal("expected job, got nil")
	}
	if job.MaxAttempts != 4 {
		t.Errorf("MaxAttempts = %d, want 4 (len(backoff)+1)", job.MaxAttempts)
	}
	if len(job.BackoffSeconds) != 3 {
		t.Errorf("BackoffSeconds len = %d, want 3", len(job.BackoffSeconds))
	}
}

func TestDispatchWithBackoff_UnknownHandler(t *testing.T) {
	ResetHandlers()
	defer ResetHandlers()

	mem := NewMemoryDriver()
	disp := NewDispatcher(mem)
	err := disp.DispatchWithBackoff(context.Background(), "default", "nope", "data", []uint{5})
	if err == nil {
		t.Fatal("expected error for unknown handler")
	}
}

// ─── RetryDelay (backoff) ─────────────────────────────────────────────────────

func TestRetryDelay_UsesBackoffSlice(t *testing.T) {
	w := NewWorker(NewMemoryDriver(), WorkerConfig{
		RetryDelay: 30 * time.Second,
	}, testLogger())

	job := &Job{BackoffSeconds: []uint{5, 30, 120}}

	// attempt 0 → backoff[0] = 5s
	job.Attempts = 0
	if d := w.retryDelay(job); d != 5*time.Second {
		t.Errorf("attempt 0: got %v, want 5s", d)
	}

	// attempt 1 → backoff[1] = 30s
	job.Attempts = 1
	if d := w.retryDelay(job); d != 30*time.Second {
		t.Errorf("attempt 1: got %v, want 30s", d)
	}

	// attempt 2 → backoff[2] = 120s
	job.Attempts = 2
	if d := w.retryDelay(job); d != 120*time.Second {
		t.Errorf("attempt 2: got %v, want 120s", d)
	}

	// attempt 3 → beyond slice → use last = 120s
	job.Attempts = 3
	if d := w.retryDelay(job); d != 120*time.Second {
		t.Errorf("attempt 3: got %v, want 120s (last)", d)
	}
}

func TestRetryDelay_FallsBackToConfig(t *testing.T) {
	w := NewWorker(NewMemoryDriver(), WorkerConfig{
		RetryDelay: 45 * time.Second,
	}, testLogger())

	job := &Job{} // no BackoffSeconds
	job.Attempts = 0
	if d := w.retryDelay(job); d != 45*time.Second {
		t.Errorf("got %v, want 45s (config fallback)", d)
	}
}
