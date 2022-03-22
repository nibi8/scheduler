package testsp

import (
	"context"
	"fmt"
	"testing"

	"github.com/nibi8/scheduler/models"
)

// todo: add tests

func TestStorageProvider(t *testing.T) {

	ctx := context.Background()

	sp := NewStorageProvider()

	jobName := "job1"

	job, err := models.NewJob(jobName, 20, 10, func(ctx context.Context, job models.Job) error {
		fmt.Println(jobName, "action1")
		return nil
	}, nil)
	if err != nil {
		t.Error(err)
		return
	}

	lr := models.NewLockRecord(job)

	err = sp.CreateLockRecord(ctx, lr)
	if err != nil {
		t.Error(err)
		return
	}

	lrResp, err := sp.GetLockRecord(ctx, job.Name)
	if err != nil {
		t.Error(err)
		return
	}

	if lr != lrResp {
		t.Errorf("GetLockRecord result differs")
		return
	}

	lrPatch := models.NewLockRecordPatch(lr.DurationSec)

	err = sp.UpdateLockRecord(ctx, job.Name, lr.Version, lrPatch)
	if err != nil {
		t.Error(err)
		return
	}

	lrUpdated, err := sp.GetLockRecord(ctx, job.Name)
	if err != nil {
		t.Error(err)
		return
	}

	lr.ApplyPatch(lrPatch)

	if lr != lrUpdated {
		t.Errorf("GetLockRecord (after UpdateLockRecord) result differs")
		return
	}

}
