package models

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dmodels "github.com/p8bin/dlocker/models"
)

func TestValidate(t *testing.T) {

	errHandler := func(context.Context, Job, error) {
	}

	lock, err := dmodels.NewLock(
		"unique-job-name",
		30,
		10,
	)
	require.NoError(t, err)

	_, err = NewJobEx(
		lock,
		func(ctx context.Context, job Job) error {
			fmt.Println("start job action")
			fmt.Println("end before ctx.Done()")
			return nil
		}, 1, 2, errHandler,
	)
	require.NoError(t, err)

	_, err = NewJobEx(
		dmodels.Lock{},
		func(ctx context.Context, job Job) error {
			fmt.Println("start job action")
			fmt.Println("end before ctx.Done()")
			return nil
		}, 1, 2, errHandler,
	)
	require.Error(t, err)

	_, err = NewJobEx(
		lock,
		nil, 1, 2, errHandler,
	)
	require.Error(t, err)

}

func TestNewJob(t *testing.T) {
	errHandler := func(context.Context, Job, error) {
	}

	lock, err := dmodels.NewLock(
		"unique-job-name",
		30,
		10,
	)
	require.NoError(t, err)

	job, err := NewJob(
		lock,
		func(ctx context.Context, job Job) error {
			fmt.Println("start job action")
			fmt.Println("end before ctx.Done()")
			return nil
		}, errHandler,
	)
	require.NoError(t, err)

	jobEx, err := NewJobEx(
		lock,
		job.Action, job.PeekTimeoutSec, job.ErrTimeoutSec, job.ErrHandler,
	)
	require.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("%+v", job), fmt.Sprintf("%+v", jobEx))

	assert.NotEmpty(t, job.Lock)
	assert.NotEmpty(t, job.Action)
	assert.NotEmpty(t, job.PeekTimeoutSec)
	assert.NotEmpty(t, job.ErrTimeoutSec)
	assert.NotEmpty(t, job.ErrHandler)
}

func TestNewJobPnc(t *testing.T) {
	errHandler := func(context.Context, Job, error) {
	}

	lock := dmodels.NewLockPnc(
		"unique-job-name",
		30,
		10,
	)

	job := NewJobPnc(
		lock,
		func(ctx context.Context, job Job) error {
			fmt.Println("start job action")
			fmt.Println("end before ctx.Done()")
			return nil
		}, errHandler,
	)

	jobEx := NewJobExPnc(
		lock,
		job.Action, job.PeekTimeoutSec, job.ErrTimeoutSec, job.ErrHandler,
	)

	assert.Equal(t, fmt.Sprintf("%+v", job), fmt.Sprintf("%+v", jobEx))

	assert.NotEmpty(t, job.Lock)
	assert.NotEmpty(t, job.Action)
	assert.NotEmpty(t, job.PeekTimeoutSec)
	assert.NotEmpty(t, job.ErrTimeoutSec)
	assert.NotEmpty(t, job.ErrHandler)
}
