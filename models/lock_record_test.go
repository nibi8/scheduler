package models

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// todo: add tests

func TestNewLockRecord(t *testing.T) {

	job, err := NewJobEx(
		"unique-job-name",
		int((30 * time.Second).Seconds()),
		int((10 * time.Second).Seconds()),
		func(ctx context.Context, job Job) error {
			fmt.Println("start job action")
			fmt.Println("end before ctx.Done()")
			return nil
		}, 0, 0, nil,
	)

	if err != nil {
		t.Error(err)
	}

	lr := NewLockRecord(job)

	fmt.Println("Lock record new unique version:", lr.Version)

}
