package mongosp

import (
	"time"

	"github.com/nibi8/scheduler/models"
)

type LockRecordDB struct {
	JobName     string `bson:"jobname"`
	Version     string `bson:"version"`
	DurationSec int    `bson:"durationsec"`

	Dt time.Time `bson:"dt"`
}

func FromLockRecord(in models.LockRecord) LockRecordDB {
	out := LockRecordDB{
		JobName:     in.JobName,
		Version:     in.Version,
		DurationSec: in.DurationSec,
		Dt:          in.Dt,
	}
	return out
}

func ToLockRecord(in LockRecordDB) models.LockRecord {
	out := models.LockRecord{
		JobName:     in.JobName,
		Version:     in.Version,
		DurationSec: in.DurationSec,
		Dt:          in.Dt,
	}
	return out
}
