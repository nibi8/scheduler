package models

import (
	"context"
	"fmt"

	lockmodels "github.com/p8bin/dlocker/models"
)

type JobAction func(context.Context, Job) error

type JobErrorHandler func(context.Context, Job, error)

// Default spap in retry cycle (in addition to duration of job lock)
var DefaultPeekTimeoutSec = 10

// Default time period after error occurs
var DefaultErrTimeoutSec = 10

type Job struct {
	Lock lockmodels.Lock

	// Job action
	Action JobAction

	// Span in retry cycle (in addition to duration of job lock)
	PeekTimeoutSec int
	// Time period after error occurs
	ErrTimeoutSec int

	// Error handler
	ErrHandler JobErrorHandler
}

func (j Job) Validate() (err error) {

	err = j.Lock.Validate()
	if err != nil {
		return err
	}

	if j.Action == nil {
		return fmt.Errorf("Action == nil")
	}

	if j.PeekTimeoutSec < 1 {
		return fmt.Errorf("PeekTimeoutSec < 1")
	}

	if j.ErrTimeoutSec < 1 {
		return fmt.Errorf("ErrTimeoutSec < 1")
	}

	if j.ErrHandler == nil {
		return fmt.Errorf("ErrHandler == nil")
	}

	return nil
}

func NewJobPnc(
	lock lockmodels.Lock,
	action JobAction,
	errHandler JobErrorHandler,
) (job Job) {
	job, err := NewJob(
		lock,
		action,
		errHandler,
	)
	if err != nil {
		panic(err)
	}
	return job
}

func NewJob(
	lock lockmodels.Lock,
	action JobAction,
	errHandler JobErrorHandler,
) (job Job, err error) {
	return NewJobEx(
		lock,
		action,
		DefaultPeekTimeoutSec,
		DefaultErrTimeoutSec,
		errHandler,
	)
}

func NewJobExPnc(
	lock lockmodels.Lock,
	action JobAction,
	peekTimeoutSec int,
	errTimeoutSec int,
	errHandler JobErrorHandler,
) (job Job) {
	job, err := NewJobEx(
		lock,
		action,
		peekTimeoutSec,
		errTimeoutSec,
		errHandler,
	)
	if err != nil {
		panic(err)
	}
	return job
}

func NewJobEx(
	lock lockmodels.Lock,
	action JobAction,
	peekTimeoutSec int,
	errTimeoutSec int,
	errHandler JobErrorHandler,
) (job Job, err error) {

	if peekTimeoutSec == 0 {
		peekTimeoutSec = DefaultPeekTimeoutSec
	}

	if errTimeoutSec == 0 {
		errTimeoutSec = DefaultErrTimeoutSec
	}

	if errHandler == nil {
		errHandler = func(ctx context.Context, j Job, err error) {}
	}

	job = Job{
		Lock:           lock,
		Action:         action,
		PeekTimeoutSec: peekTimeoutSec,
		ErrTimeoutSec:  errTimeoutSec,
		ErrHandler:     errHandler,
	}

	err = job.Validate()
	if err != nil {
		return job, err
	}

	return job, nil
}
