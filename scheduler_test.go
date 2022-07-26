package scheduler

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/p8bin/scheduler/models"

	"github.com/p8bin/dlocker"
	dmodels "github.com/p8bin/dlocker/models"
	"github.com/p8bin/dlocker/storageproviders/testsp"
)

// todo: add tests

func TestNewScheduler(t *testing.T) {
	sp := testsp.NewStorageProvider()
	locker := dlocker.NewLocker(sp)
	sc := NewScheduler(locker)

	lock, err := dmodels.NewLock(
		"unique-job-name",
		30,
		10,
	)
	require.NoError(t, err)

	job, err := models.NewJobEx(
		lock,
		func(ctx context.Context, job models.Job) error {
			fmt.Println("start job action")
			fmt.Println("end before ctx.Done()")
			return nil
		}, 0, 0, nil,
	)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = sc.RunJob(ctx, job)
	require.NoError(t, err)
}
