package pubsub

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newTestClient(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return client, mr
}

// TC-01: Publish and Receive
func TestPublishAndReceive(t *testing.T) {
	client, _ := newTestClient(t)

	pub := NewRedisPublisher(client)
	sub := NewRedisSubscriber(client, SubscriberOptions{})

	var received struct {
		mu      sync.Mutex
		channel string
		payload string
	}
	done := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go sub.Subscribe(ctx, []string{"test"}, func(_ context.Context, ch, payload string) {
		received.mu.Lock()
		received.channel = ch
		received.payload = payload
		received.mu.Unlock()
		close(done)
	})

	// Give subscriber time to connect.
	time.Sleep(100 * time.Millisecond)

	if err := pub.Publish(ctx, "test", "hello"); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for message")
	}

	received.mu.Lock()
	defer received.mu.Unlock()
	if received.channel != "test" {
		t.Errorf("channel = %q, want %q", received.channel, "test")
	}
	if received.payload != "hello" {
		t.Errorf("payload = %q, want %q", received.payload, "hello")
	}
}

// TC-02: Multiple Channels
func TestMultipleChannels(t *testing.T) {
	client, _ := newTestClient(t)

	pub := NewRedisPublisher(client)
	sub := NewRedisSubscriber(client, SubscriberOptions{})

	var count atomic.Int32
	done := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go sub.Subscribe(ctx, []string{"ch1", "ch2"}, func(_ context.Context, ch, payload string) {
		if count.Add(1) == 2 {
			close(done)
		}
	})

	time.Sleep(100 * time.Millisecond)

	pub.Publish(ctx, "ch1", "msg1")
	pub.Publish(ctx, "ch2", "msg2")

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for messages on both channels")
	}

	if got := count.Load(); got != 2 {
		t.Errorf("handler called %d times, want 2", got)
	}
}

// TC-03: Context Cancellation (Clean Shutdown)
func TestContextCancel(t *testing.T) {
	client, _ := newTestClient(t)

	sub := NewRedisSubscriber(client, SubscriberOptions{})

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)

	go func() {
		errCh <- sub.Subscribe(ctx, []string{"test"}, func(_ context.Context, _, _ string) {})
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Subscribe returned error %v, want nil", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Subscribe did not return after context cancel")
	}
}

// TC-04: Reconnect After Disconnect
func TestReconnect(t *testing.T) {
	mr := miniredis.RunT(t)

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	sub := NewRedisSubscriber(client, SubscriberOptions{
		MinBackoff: 50 * time.Millisecond,
		MaxBackoff: 200 * time.Millisecond,
	})

	var count atomic.Int32
	done := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go sub.Subscribe(ctx, []string{"test"}, func(_ context.Context, _, payload string) {
		if payload == "after-reconnect" {
			count.Add(1)
			close(done)
		}
	})

	time.Sleep(100 * time.Millisecond)

	// Disconnect and restart on same address.
	mr.Close()
	time.Sleep(100 * time.Millisecond)
	mr.Restart()

	// Wait for reconnect.
	time.Sleep(500 * time.Millisecond)

	pub := NewRedisPublisher(redis.NewClient(&redis.Options{Addr: mr.Addr()}))
	if err := pub.Publish(ctx, "test", "after-reconnect"); err != nil {
		t.Fatalf("Publish after reconnect: %v", err)
	}

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for message after reconnect")
	}

	if got := count.Load(); got != 1 {
		t.Errorf("handler called %d times, want 1", got)
	}
}

// TC-05: Publisher Error on Failed Publish
func TestPublishError(t *testing.T) {
	client, mr := newTestClient(t)
	pub := NewRedisPublisher(client)
	mr.Close()

	err := pub.Publish(context.Background(), "test", "msg")
	if err == nil {
		t.Error("expected error from Publish after server close, got nil")
	}
}

// TC-06: SubscriberOptions Defaults
func TestSubscriberDefaults(t *testing.T) {
	opts := applyDefaults(SubscriberOptions{})

	if opts.MinBackoff != 500*time.Millisecond {
		t.Errorf("MinBackoff = %v, want 500ms", opts.MinBackoff)
	}
	if opts.MaxBackoff != 30*time.Second {
		t.Errorf("MaxBackoff = %v, want 30s", opts.MaxBackoff)
	}
	if opts.PingEvery != 30*time.Second {
		t.Errorf("PingEvery = %v, want 30s", opts.PingEvery)
	}
	if opts.Logger == nil {
		t.Error("Logger should not be nil after applyDefaults")
	}
}

// Edge case: Subscribe with empty channels
func TestSubscribeEmptyChannels(t *testing.T) {
	client, _ := newTestClient(t)
	sub := NewRedisSubscriber(client, SubscriberOptions{})

	err := sub.Subscribe(context.Background(), []string{}, func(_ context.Context, _, _ string) {})
	if err == nil {
		t.Error("expected error for empty channels, got nil")
	}
}

// Edge case: Handler panic does not crash subscriber
func TestHandlerPanic(t *testing.T) {
	client, _ := newTestClient(t)

	pub := NewRedisPublisher(client)
	sub := NewRedisSubscriber(client, SubscriberOptions{})

	var count atomic.Int32
	done := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go sub.Subscribe(ctx, []string{"test"}, func(_ context.Context, _, payload string) {
		if payload == "panic" {
			panic("test panic")
		}
		count.Add(1)
		close(done)
	})

	time.Sleep(100 * time.Millisecond)

	// First message triggers panic — subscriber should survive.
	pub.Publish(ctx, "test", "panic")
	time.Sleep(100 * time.Millisecond)

	// Second message should still be received.
	pub.Publish(ctx, "test", "ok")

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("subscriber did not survive handler panic")
	}

	if got := count.Load(); got != 1 {
		t.Errorf("handler called %d times after panic recovery, want 1", got)
	}
}
