package mongosp

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/nibi8/scheduler/models"
)

// todo: add tests

func TestStorageProvider(t *testing.T) {

	ctx := context.Background()

	// connect to mongodb
	constr := "mongodb://localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(constr))
	if err != nil {
		log.Fatal("mongo.Connect")
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("client.Ping")
	}

	db := client.Database("test")

	collectionName := "schedulerTest"

	_ = db.Collection(collectionName).Drop(ctx)

	// create storage provider
	sp, err := NewStorageProvider(ctx, db, collectionName)
	if err != nil {
		t.Error(err)
		return
	}

	jobName := "job1"

	job, err := models.NewJob(jobName, 20, 10, func(ctx context.Context, job models.Job) error {
		fmt.Println(jobName, "action1")
		return nil
	}, nil)
	if err != nil {
		t.Error(err)
		return
	}

	lr := models.NewLockRecord(job)
	lr.Dt = time.Unix(time.Now().Unix(), 0).UTC() // fix golang and db format for later compare in tests

	err = sp.CreateLockRecord(ctx, lr)
	if err != nil {
		t.Error(err)
		return
	}

	lrResp, err := sp.GetLockRecord(ctx, job.Name)
	if err != nil {
		t.Error(err)
		return
	}

	if lr != lrResp {
		t.Errorf("GetLockRecord result differs")
		return
	}

	lrPatch := models.NewLockRecordPatch(lr.DurationSec)
	lrPatch.Dt = time.Unix(time.Now().Unix(), 0).UTC() // fix golang and db format for later compare in tests

	err = sp.UpdateLockRecord(ctx, job.Name, lr.Version, lrPatch)
	if err != nil {
		t.Error(err)
		return
	}

	lrUpdated, err := sp.GetLockRecord(ctx, job.Name)
	if err != nil {
		t.Error(err)
		return
	}

	lr.ApplyPatch(lrPatch)

	if lr != lrUpdated {
		t.Errorf("GetLockRecord (after UpdateLockRecord) result differs")
		return
	}

}
