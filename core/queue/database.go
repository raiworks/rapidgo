package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// jobModel is the GORM model for the jobs table.
type jobModel struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement"`
	Queue       string     `gorm:"size:255;not null;index"`
	Type        string     `gorm:"size:255;not null"`
	Payload     string     `gorm:"type:text;not null"`
	Attempts    uint       `gorm:"not null;default:0"`
	MaxAttempts uint       `gorm:"not null;default:3"`
	AvailableAt time.Time  `gorm:"not null;index"`
	ReservedAt  *time.Time `gorm:"index"`
	CreatedAt   time.Time  `gorm:"not null"`
}

// failedJobModel is the GORM model for the failed_jobs table.
type failedJobModel struct {
	ID       uint64    `gorm:"primaryKey;autoIncrement"`
	Queue    string    `gorm:"size:255;not null"`
	Type     string    `gorm:"size:255;not null"`
	Payload  string    `gorm:"type:text;not null"`
	Error    string    `gorm:"type:text;not null"`
	FailedAt time.Time `gorm:"not null"`
}

// DatabaseDriver implements Driver using GORM.
type DatabaseDriver struct {
	db          *gorm.DB
	table       string
	failedTable string
}

// NewDatabaseDriver creates a database-backed queue driver.
func NewDatabaseDriver(db *gorm.DB, table, failedTable string) *DatabaseDriver {
	return &DatabaseDriver{
		db:          db,
		table:       table,
		failedTable: failedTable,
	}
}

func (d *DatabaseDriver) Push(_ context.Context, job *Job) error {
	m := jobModel{
		Queue:       job.Queue,
		Type:        job.Type,
		Payload:     string(job.Payload),
		Attempts:    job.Attempts,
		MaxAttempts: job.MaxAttempts,
		AvailableAt: job.AvailableAt,
		CreatedAt:   job.CreatedAt,
	}
	result := d.db.Table(d.table).Create(&m)
	if result.Error != nil {
		return fmt.Errorf("queue: database push failed: %w", result.Error)
	}
	job.ID = m.ID
	return nil
}

func (d *DatabaseDriver) Pop(_ context.Context, queue string) (*Job, error) {
	var m jobModel
	now := time.Now()

	err := d.db.Transaction(func(tx *gorm.DB) error {
		// Try FOR UPDATE SKIP LOCKED first (MySQL, Postgres).
		// Falls back to simple query for SQLite.
		result := tx.Table(d.table).
			Where("queue = ? AND available_at <= ? AND reserved_at IS NULL", queue, now).
			Order("id ASC").
			Limit(1).
			Set("gorm:query_option", "FOR UPDATE SKIP LOCKED").
			First(&m)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				return nil
			}
			// Retry without SKIP LOCKED (SQLite compatibility).
			result = tx.Table(d.table).
				Where("queue = ? AND available_at <= ? AND reserved_at IS NULL", queue, now).
				Order("id ASC").
				Limit(1).
				First(&m)
			if result.Error != nil {
				if result.Error == gorm.ErrRecordNotFound {
					return nil
				}
				return result.Error
			}
		}

		if m.ID == 0 {
			return nil
		}

		// Reserve the job.
		return tx.Table(d.table).Where("id = ?", m.ID).Update("reserved_at", now).Error
	})

	if err != nil {
		return nil, fmt.Errorf("queue: database pop failed: %w", err)
	}

	if m.ID == 0 {
		return nil, nil
	}

	reserved := now
	return &Job{
		ID:          m.ID,
		Queue:       m.Queue,
		Type:        m.Type,
		Payload:     json.RawMessage(m.Payload),
		Attempts:    m.Attempts,
		MaxAttempts: m.MaxAttempts,
		AvailableAt: m.AvailableAt,
		ReservedAt:  &reserved,
		CreatedAt:   m.CreatedAt,
	}, nil
}

func (d *DatabaseDriver) Delete(_ context.Context, job *Job) error {
	result := d.db.Table(d.table).Where("id = ?", job.ID).Delete(&jobModel{})
	if result.Error != nil {
		return fmt.Errorf("queue: database delete failed: %w", result.Error)
	}
	return nil
}

func (d *DatabaseDriver) Release(_ context.Context, job *Job, delay time.Duration) error {
	result := d.db.Table(d.table).Where("id = ?", job.ID).Updates(map[string]interface{}{
		"reserved_at":  nil,
		"available_at": time.Now().Add(delay),
		"attempts":     job.Attempts + 1,
	})
	if result.Error != nil {
		return fmt.Errorf("queue: database release failed: %w", result.Error)
	}
	return nil
}

func (d *DatabaseDriver) Fail(_ context.Context, job *Job, jobErr error) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		// Insert into failed_jobs.
		fm := failedJobModel{
			Queue:    job.Queue,
			Type:     job.Type,
			Payload:  string(job.Payload),
			Error:    jobErr.Error(),
			FailedAt: time.Now(),
		}
		if err := tx.Table(d.failedTable).Create(&fm).Error; err != nil {
			return fmt.Errorf("queue: database fail insert failed: %w", err)
		}

		// Delete from jobs.
		if err := tx.Table(d.table).Where("id = ?", job.ID).Delete(&jobModel{}).Error; err != nil {
			return fmt.Errorf("queue: database fail delete failed: %w", err)
		}

		return nil
	})
}

func (d *DatabaseDriver) Size(_ context.Context, queue string) (int64, error) {
	var count int64
	result := d.db.Table(d.table).
		Where("queue = ? AND reserved_at IS NULL", queue).
		Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("queue: database size failed: %w", result.Error)
	}
	return count, nil
}
