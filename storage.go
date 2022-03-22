package scheduler

import (
	"context"
	"errors"
	"time"

	"github.com/nibi8/scheduler/models"
)

type Storage struct {
	sp StorageProvider
}

func NewStorage(
	sp StorageProvider,
) *Storage {
	s := Storage{
		sp: sp,
	}
	return &s
}

// Set a lock on the job or return ErrNoLuck if no luck or an unexpected error occurs
func (s *Storage) SetLock(ctx context.Context, job models.Job) (err error) {

	lrFound := false
	lr, err := s.sp.GetLockRecord(ctx, job.Name)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			return err
		}
		// not found
		// continue
		err = nil
		lrFound = false
	} else {
		lrFound = true
	}

	if lrFound && lr.DurationSec > 0 {
		time.Sleep(time.Duration(lr.DurationSec) * time.Second)
	}

	if !lrFound {
		lr = models.NewLockRecord(job)
		err = s.sp.CreateLockRecord(ctx, lr)
		if err != nil {
			if errors.Is(err, models.ErrDuplicate) {
				return models.ErrNoLuck
			}
			return err
		}
		return nil
	}

	lrPatch := models.NewLockRecordPatch(job.GetDurationSec())
	err = s.sp.UpdateLockRecord(ctx, job.Name, lr.Version, lrPatch)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return models.ErrNoLuck
		}
		return err
	}

	return nil
}
