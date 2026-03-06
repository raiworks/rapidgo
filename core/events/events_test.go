package events

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TC-01: Listen + DispatchSync invokes handler.
func TestDispatchSync_InvokesHandler(t *testing.T) {
	d := NewDispatcher()
	called := false
	d.Listen("test", func(p interface{}) { called = true })
	d.DispatchSync("test", nil)
	if !called {
		t.Fatal("handler not called")
	}
}

// TC-02: Multiple listeners all invoked.
func TestDispatchSync_MultipleListeners(t *testing.T) {
	d := NewDispatcher()
	var count int
	for i := 0; i < 3; i++ {
		d.Listen("test", func(p interface{}) { count++ })
	}
	d.DispatchSync("test", nil)
	if count != 3 {
		t.Fatalf("count = %d, want 3", count)
	}
}

// TC-03: Dispatch fires handlers asynchronously.
func TestDispatch_Async(t *testing.T) {
	d := NewDispatcher()
	var wg sync.WaitGroup
	wg.Add(1)
	d.Listen("async", func(p interface{}) {
		defer wg.Done()
	})
	d.Dispatch("async", nil)
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("async handler did not complete in time")
	}
}

// TC-04: Dispatching unknown event is a no-op.
func TestDispatch_UnknownEvent(t *testing.T) {
	d := NewDispatcher()
	// Should not panic.
	d.Dispatch("nope", nil)
	d.DispatchSync("nope", nil)
}

// TC-05: Has returns true for registered event.
func TestHas_True(t *testing.T) {
	d := NewDispatcher()
	d.Listen("x", func(p interface{}) {})
	if !d.Has("x") {
		t.Fatal("Has(x) = false, want true")
	}
}

// TC-06: Has returns false for unregistered event.
func TestHas_False(t *testing.T) {
	d := NewDispatcher()
	if d.Has("y") {
		t.Fatal("Has(y) = true, want false")
	}
}

// TC-07: Payload passed correctly.
func TestDispatchSync_Payload(t *testing.T) {
	d := NewDispatcher()
	var received interface{}
	d.Listen("data", func(p interface{}) { received = p })
	d.DispatchSync("data", "hello")
	if received != "hello" {
		t.Fatalf("payload = %v, want %q", received, "hello")
	}
}

// TC-08: Concurrent Listen + Dispatch is safe.
func TestConcurrentListenDispatch(t *testing.T) {
	d := NewDispatcher()
	var ops atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			d.Listen("concurrent", func(p interface{}) {
				ops.Add(1)
			})
		}()
		go func() {
			defer wg.Done()
			d.DispatchSync("concurrent", nil)
		}()
	}
	wg.Wait()
}
