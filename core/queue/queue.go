package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Job represents a unit of work to be processed by a worker.
type Job struct {
	ID          uint64
	Queue       string
	Type        string
	Payload     json.RawMessage
	Attempts    uint
	MaxAttempts uint
	AvailableAt time.Time
	ReservedAt  *time.Time
	CreatedAt   time.Time
}

// HandlerFunc processes a job. Receives the raw JSON payload.
type HandlerFunc func(ctx context.Context, payload json.RawMessage) error

// handler registry
var (
	handlersMu sync.RWMutex
	handlers   = make(map[string]HandlerFunc)
)

// RegisterHandler maps a type name to a handler function.
func RegisterHandler(typeName string, handler HandlerFunc) {
	handlersMu.Lock()
	defer handlersMu.Unlock()
	handlers[typeName] = handler
}

// ResolveHandler returns the handler for a type name, or nil.
func ResolveHandler(typeName string) HandlerFunc {
	handlersMu.RLock()
	defer handlersMu.RUnlock()
	return handlers[typeName]
}

// ResetHandlers clears the registry. For testing only.
func ResetHandlers() {
	handlersMu.Lock()
	defer handlersMu.Unlock()
	handlers = make(map[string]HandlerFunc)
}

// Driver is the storage backend for the queue system.
type Driver interface {
	// Push adds a job to the queue.
	Push(ctx context.Context, job *Job) error

	// Pop retrieves and reserves the next available job from the given queue.
	// Returns nil, nil if no job is available.
	Pop(ctx context.Context, queue string) (*Job, error)

	// Delete removes a completed job.
	Delete(ctx context.Context, job *Job) error

	// Release puts a reserved job back into the queue for retry.
	Release(ctx context.Context, job *Job, delay time.Duration) error

	// Fail moves a job to the failed jobs storage.
	Fail(ctx context.Context, job *Job, jobErr error) error

	// Size returns the number of pending jobs in a queue.
	Size(ctx context.Context, queue string) (int64, error)
}

// Dispatcher is the public API for dispatching jobs.
type Dispatcher struct {
	driver Driver
}

// NewDispatcher creates a dispatcher with the given driver.
func NewDispatcher(driver Driver) *Dispatcher {
	return &Dispatcher{driver: driver}
}

// Dispatch pushes a job onto the named queue.
func (d *Dispatcher) Dispatch(ctx context.Context, queue, typeName string, payload interface{}) error {
	if handler := ResolveHandler(typeName); handler == nil {
		return fmt.Errorf("queue: no handler registered for type %q", typeName)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("queue: failed to marshal payload: %w", err)
	}

	job := &Job{
		Queue:       queue,
		Type:        typeName,
		Payload:     data,
		Attempts:    0,
		MaxAttempts: 3,
		AvailableAt: time.Now(),
		CreatedAt:   time.Now(),
	}

	return d.driver.Push(ctx, job)
}

// DispatchDelayed pushes a job onto the queue with a delay.
func (d *Dispatcher) DispatchDelayed(ctx context.Context, queue, typeName string, payload interface{}, delay time.Duration) error {
	if handler := ResolveHandler(typeName); handler == nil {
		return fmt.Errorf("queue: no handler registered for type %q", typeName)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("queue: failed to marshal payload: %w", err)
	}

	job := &Job{
		Queue:       queue,
		Type:        typeName,
		Payload:     data,
		Attempts:    0,
		MaxAttempts: 3,
		AvailableAt: time.Now().Add(delay),
		CreatedAt:   time.Now(),
	}

	return d.driver.Push(ctx, job)
}

// Driver returns the underlying driver.
func (d *Dispatcher) Driver() Driver {
	return d.driver
}
