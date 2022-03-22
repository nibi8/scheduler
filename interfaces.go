package scheduler

import (
	"context"

	"github.com/nibi8/scheduler/models"
)

// Runs distributed jobs.
// Job execution sets a lock for a total duration period and requires completion before the lock ends.
type Scheduler interface {
	// Runs job within new goroutine
	RunJob(ctx context.Context, job models.Job)
}

// Sets locks for jobs.
// You need implement this in your persistent storage.
type StorageProvider interface {
	// Returns LockRecord or error ErrNotFound if not found or unexpected error.
	GetLockRecord(ctx context.Context, jobName string) (lr models.LockRecord, err error)

	// Creates LockRecord or returns error ErrDuplicate if already exists or unexpected error.
	// This happens once.
	CreateLockRecord(ctx context.Context, lr models.LockRecord) (err error)

	// Updates LockRecord with new values or return error ErrNotFound if not found or unexpected error.
	UpdateLockRecord(ctx context.Context, jobName string, version string, patch models.LockRecordPatch) (err error)
}
