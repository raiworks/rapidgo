package queue

import (
	"context"
	"fmt"
	"time"
)

// SyncDriver executes jobs immediately when dispatched.
// No worker needed. Useful for local development and testing.
type SyncDriver struct{}

// NewSyncDriver creates a synchronous queue driver.
func NewSyncDriver() *SyncDriver {
	return &SyncDriver{}
}

func (d *SyncDriver) Push(ctx context.Context, job *Job) error {
	handler := ResolveHandler(job.Type)
	if handler == nil {
		return fmt.Errorf("queue: no handler registered for type %q", job.Type)
	}
	return handler(ctx, job.Payload)
}

func (d *SyncDriver) Pop(_ context.Context, _ string) (*Job, error) {
	return nil, nil
}

func (d *SyncDriver) Delete(_ context.Context, _ *Job) error {
	return nil
}

func (d *SyncDriver) Release(_ context.Context, _ *Job, _ time.Duration) error {
	return nil
}

func (d *SyncDriver) Fail(_ context.Context, _ *Job, _ error) error {
	return nil
}

func (d *SyncDriver) Size(_ context.Context, _ string) (int64, error) {
	return 0, nil
}
