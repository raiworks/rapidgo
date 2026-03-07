package queue

import (
	"context"
	"sync"
	"time"
)

// MemoryDriver implements Driver using in-process slices.
// No persistence — useful for testing and development.
type MemoryDriver struct {
	mu     sync.Mutex
	queues map[string][]*Job
	failed []*Job
	nextID uint64
}

// NewMemoryDriver creates an in-memory queue driver.
func NewMemoryDriver() *MemoryDriver {
	return &MemoryDriver{
		queues: make(map[string][]*Job),
	}
}

func (d *MemoryDriver) Push(_ context.Context, job *Job) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.nextID++
	job.ID = d.nextID
	d.queues[job.Queue] = append(d.queues[job.Queue], job)
	return nil
}

func (d *MemoryDriver) Pop(_ context.Context, queue string) (*Job, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	jobs := d.queues[queue]
	for i, job := range jobs {
		if job.ReservedAt == nil && !job.AvailableAt.After(now) {
			reserved := now
			job.ReservedAt = &reserved
			// Remove from slice (order preserved).
			d.queues[queue] = append(jobs[:i], jobs[i+1:]...)
			return job, nil
		}
	}
	return nil, nil
}

func (d *MemoryDriver) Delete(_ context.Context, job *Job) error {
	// Job was already removed from the queue slice in Pop.
	// Nothing more to do.
	return nil
}

func (d *MemoryDriver) Release(_ context.Context, job *Job, delay time.Duration) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	job.Attempts++
	job.ReservedAt = nil
	job.AvailableAt = time.Now().Add(delay)
	d.queues[job.Queue] = append(d.queues[job.Queue], job)
	return nil
}

func (d *MemoryDriver) Fail(_ context.Context, job *Job, _ error) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.failed = append(d.failed, job)
	return nil
}

func (d *MemoryDriver) Size(_ context.Context, queue string) (int64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var count int64
	for _, job := range d.queues[queue] {
		if job.ReservedAt == nil {
			count++
		}
	}
	return count, nil
}
