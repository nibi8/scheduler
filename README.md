# Distributed scheduler

Tiny distributed scheduler.
You need to implement persistent lock storage in order to use it (or use the provided mongodb storage provider).

## Scheme

Job "super-job" with 5 min duration:

```
Instace 1: [set lock success] [running within the duration of the lock (5 minutes)] [try set lock] ...
```

Must complete execution before the lock expires.

```
Instace 2:       [set lock fail] [sleep during lock (5 min)                       ] [try set lock] ...
```

Execution of other jobs is not affected.

## Usage

Create storage provider. Implement `StorageProvider` interface for you persistent storage or use default implementation (mongosp package has a mongodb implementation)

```go
sp, err := mongosp.NewStorageProvider(ctx, db, "scheduler")
```

Create scheduler

```go
schedulerSvc := scheduler.NewScheduler(sp)
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
