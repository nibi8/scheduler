package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nibi8/scheduler/models"
	"github.com/nibi8/scheduler/storageproviders/testsp"
)

// todo: add tests

func TestNewStorage(t *testing.T) {
	sp := testsp.NewStorageProvider()
	storage := NewStorage(sp)

	job, err := models.NewJobEx(
		"unique-job-name",
		int((30 * time.Second).Seconds()),
		int((10 * time.Second).Seconds()),
		func(ctx context.Context, job models.Job) error {
			fmt.Println("start job action")
			fmt.Println("end before ctx.Done()")
			return nil
		}, 0, 0, nil,
	)
	if err != nil {
		t.Error(err)
	}

	err = storage.SetLock(context.Background(), job)
	if err != nil {
		t.Error(err)
	}

}
