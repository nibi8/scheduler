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
	// Unique job name
	Name string

	// Duration of job in seconds
	ExecutionDurationSec int

	// Interval between jobs execution in seconds
	SpanDurationSec int

	// Job action
	Action JobAction

	// Spap in retry cycle (in addition to duration of job lock)
	PeekTimeoutSec int
	// Time period after error occurs
	ErrTimeoutSec int

	// Error handler
	ErrHandler JobErrorHandler
}

func (job Job) ToLock() (lock lockmodels.Lock, err error) {
	return lockmodels.NewLock(
		job.Name, 
		job.ExecutionDurationSec,
		job.SpanDurationSec,
	)
}

func (j Job) Validate() (err error) {
	if j.Name == "" {
		return fmt.Errorf(`Name == ""`)
	}

	if j.ExecutionDurationSec < 1 {
		return fmt.Errorf("ExecutionDurationSec < 1")
	}

	if j.SpanDurationSec < 1 {
		return fmt.Errorf("SpanDurationSec < 1")
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

// Total lock period
func (j Job) GetDurationSec() int {
	return j.ExecutionDurationSec + j.SpanDurationSec
}

func NewJob(
	name string,
	executionDurationSec int,
	spanDurationSec int,
	action JobAction,
	errHandler JobErrorHandler,
) (job Job, err error) {
	return NewJobEx(
		name,
		executionDurationSec,
		spanDurationSec,
		action,
		DefaultPeekTimeoutSec,
		DefaultErrTimeoutSec,
		errHandler,
	)
}

func NewJobEx(
	name string,
	executionDurationSec int,
	spanDurationSec int,
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
		Name:                 name,
		ExecutionDurationSec: executionDurationSec,
		SpanDurationSec:      spanDurationSec,
		Action:               action,
		PeekTimeoutSec:       peekTimeoutSec,
		ErrTimeoutSec:        errTimeoutSec,
		ErrHandler:           errHandler,
	}

	err = job.Validate()
	if err != nil {
		return job, err
	}

	return job, nil
}
