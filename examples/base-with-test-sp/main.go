package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nibi8/scheduler"
	"github.com/nibi8/scheduler/models"
	"github.com/nibi8/scheduler/storageproviders/testsp"
)

func main() {

	// create storage provider
	sp := testsp.NewStorageProvider()

	// create scheduler
	schedulerSvc := scheduler.NewScheduler(sp)

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

	schedulerSvc.RunJob(ctx, job1)
	schedulerSvc.RunJob(ctx, job1AnotherInstace)
	schedulerSvc.RunJob(ctx, job2)

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
