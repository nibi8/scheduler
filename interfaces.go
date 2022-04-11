package scheduler

import (
	"context"

	"github.com/nibi8/scheduler/models"
)

// Runs distributed jobs.
// Job execution sets a lock for a total duration period and requires completion before the lock ends.
type Scheduler interface {
	// Runs job within new goroutine
	RunJob(ctx context.Context, job models.Job) error
}
