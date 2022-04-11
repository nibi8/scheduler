package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nibi8/scheduler"
	"github.com/nibi8/scheduler/models"

	"github.com/nibi8/dlocker"
	"github.com/nibi8/dlocker/storageproviders/testsp"
)

func main() {

	// create storage provider
	sp := testsp.NewStorageProvider()

	// create locker
	locker := dlocker.NewLocker(sp)

	// create scheduler
	schedulerSvc := scheduler.NewScheduler(locker)

	// create jobs

	jobName := "super_job"
	jobName2 := "another_job"

	// create job
	job1, err := newJob(jobName, "instace_1")
	if err != nil {
		log.Fatal("newJob")
	}

	// create job again (simulate another service instance)
	job1AnotherInstace, err := newJob(jobName, "instace_2")
	if err != nil {
		log.Fatal("newJob")
	}

	// create another job
	job2, err := newJob(jobName2, "")
	if err != nil {
		log.Fatal("newJob")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err = schedulerSvc.RunJob(ctx, job1)
	if err != nil {
		log.Fatal("RunJob")
	}

	err = schedulerSvc.RunJob(ctx, job1AnotherInstace)
	if err != nil {
		log.Fatal("RunJob")
	}

	err = schedulerSvc.RunJob(ctx, job2)
	if err != nil {
		log.Fatal("RunJob")
	}

	<-ctx.Done()
}

func newJob(jobName string, instanceName string) (job models.Job, err error) {
	jobPrintName := jobName
	if instanceName != "" {
		jobPrintName += " " + instanceName
	}
	job, err = models.NewJobEx(jobName, 10, 5, func(ctx context.Context, job models.Job) error {
		for i := 0; i < 5; i++ {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			now := time.Now()
			fmt.Printf(
				"[%v] %v: %v \n",
				now.Format("15:04:05"),
				jobPrintName,
				i+1,
			)

			time.Sleep(1 * time.Second)
		}
		return nil
	}, 5, 5, func(ctx context.Context, job models.Job, err error) {
		fmt.Println(job.Name, err)
	})
	if err != nil {
		return job, err
	}

	return job, nil
}
