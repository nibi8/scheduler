//
// !!! For test use only
//
package testsp

import (
	"context"
	"fmt"
	"sync"

	"github.com/nibi8/scheduler/models"
)

//
// !!! For test use only
//

type StorageProvider struct {
	jobLocks sync.Map
}

func NewStorageProvider() *StorageProvider {
	res := StorageProvider{}
	return &res
}

type JobLock struct {
	Lock *sync.Mutex
	Lr   models.LockRecord
}

func NewJobLock(lr models.LockRecord) *JobLock {
	return &JobLock{
		Lock: &sync.Mutex{},
		Lr:   lr,
	}
}

func (sp *StorageProvider) GetLockRecord(
	ctx context.Context,
	jobName string,
) (lr models.LockRecord, err error) {

	value, ok := sp.jobLocks.Load(jobName)
	if !ok {
		return lr, models.ErrNotFound
	}

	jl, ok := value.(*JobLock)
	if !ok {
		return lr, fmt.Errorf("cast error")
	}

	return jl.Lr, nil
}

func (sp *StorageProvider) CreateLockRecord(
	ctx context.Context,
	lr models.LockRecord,
) (err error) {

	jl := NewJobLock(lr)
	_, loaded := sp.jobLocks.LoadOrStore(lr.JobName, jl)
	if loaded {
		return models.ErrDuplicate
	}

	return nil
}

func (sp *StorageProvider) UpdateLockRecord(
	ctx context.Context,
	jobName string,
	version string,
	patch models.LockRecordPatch,
) (err error) {

	value, ok := sp.jobLocks.Load(jobName)
	if !ok {
		return models.ErrNotFound
	}

	jl, ok := value.(*JobLock)
	if !ok {
		return fmt.Errorf("cast error")
	}

	jl.Lock.Lock()
	defer jl.Lock.Unlock()

	if jl.Lr.Version != version {
		return models.ErrNoLuck
	}

	jl.Lr.ApplyPatch(patch)

	return nil
}
