package scheduler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/p8bin/dlocker"
	lockmodels "github.com/p8bin/dlocker/models"

	"github.com/p8bin/scheduler/models"
)

type SchedulerImp struct {
	locker dlocker.Locker
}

func NewScheduler(
	locker dlocker.Locker,
) *SchedulerImp {
	svc := SchedulerImp{
		locker: locker,
	}
	return &svc
}

func (s *SchedulerImp) RunJob(
	ctx context.Context,
	job models.Job,
) (err error) {
	lock := job.Lock

	// todo: ? catch panic
	go func() {
		for {
			lockCtx, _, err := s.locker.LockWithWait(ctx, lock)
			if err != nil {
				if errors.Is(err, lockmodels.ErrNoLuck) {
					// no luck
					// continue
				} else {
					job.ErrHandler(ctx, job, err)
					time.Sleep(time.Duration(job.ErrTimeoutSec) * time.Second)
				}
			} else {
				err := s.runJobAction(lockCtx, job)
				if err != nil {
					job.ErrHandler(ctx, job, err)
					time.Sleep(time.Duration(job.ErrTimeoutSec) * time.Second)
				}
			}
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Duration(job.PeekTimeoutSec) * time.Second)
			}
		}
	}()

	return nil
}

func (s *SchedulerImp) runJobAction(
	ctx context.Context,
	job models.Job,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recover panic %v", r)
		}
	}()

	err = job.Action(ctx, job)
	if err != nil {
		return err
	}

	return nil
}
