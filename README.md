# Distributed scheduler

Tiny distributed scheduler.
For distributed locks, dlocker package is used.

## Scheme

Execute action during lock period.

## Usage

Create storage provider. Implement `StorageProvider` interface for you persistent storage or use default implementation (mongosp package has a mongodb implementation)

```go
sp, err := mongosp.NewStorageProvider(ctx, db, "scheduler")
```

Create locker

```go
locker := dlocker.NewLocker(sp)
```

Create scheduler

```go
schedulerSvc := scheduler.NewScheduler(locker)
```

Create job

```go
executionDurationSec := 60
spanDurationSec := 10
job, err := models.NewJob("super-job", executionDurationSec, spanDurationSec, func(ctx context.Context, job models.Job) error {
  if ctx.Err() != nil {
    return ctx.Err()
  } 
  fmt.Println("Run some code")
  return nil
}, func(ctx context.Context, job models.Job, err error) {
  fmt.Println(job.Name, err)
})
```

For more details see examples.
