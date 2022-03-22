package models

import (
	"time"

	"github.com/google/uuid"
)

type LockRecord struct {
	JobName     string
	Version     string
	DurationSec int

	Dt time.Time
}

func NewLockRecord(
	job Job,
) LockRecord {
	lr := LockRecord{
		JobName:     job.Name,
		Version:     uuid.New().String(),
		DurationSec: job.GetDurationSec(),

		Dt: time.Now(),
	}
	return lr
}

type LockRecordPatch struct {
	Version     string
	DurationSec int

	Dt time.Time
}

func NewLockRecordPatch(
	durationSec int,
) LockRecordPatch {
	lr := LockRecordPatch{
		Version:     uuid.New().String(),
		DurationSec: durationSec,

		Dt: time.Now(),
	}
	return lr
}

func (lr *LockRecord) ApplyPatch(patch LockRecordPatch) {
	lr.Version = patch.Version
	lr.DurationSec = patch.DurationSec
	lr.Dt = patch.Dt
}
