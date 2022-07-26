package scheduler

import (
	"context"

	"github.com/p8bin/scheduler/models"
)

// Runs distributed jobs.
type Scheduler interface {
	// Runs job within new goroutine
	RunJob(ctx context.Context, job models.Job) error
}
