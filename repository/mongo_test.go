package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/go-cmp/cmp/cmpopts"

	"make-me-reliable/internal/config"
	"make-me-reliable/internal/database"

	"github.com/google/go-cmp/cmp"
)

func getDBCollection() (*mongo.Collection, error) {
	c := config.ParseFromEnv()
	db, err := database.MongoConnectionFromConfigCtx(context.Background(), c)
	if err != nil {
		return nil, err
	}

	col := db.Database(c.DatabaseName).Collection("test-1")
	return col, nil
}

func TestMongoRepository_GetUnfinishedJobsCtx(t *testing.T) {
	col, err := getDBCollection()
	if err != nil {
		t.Fatalf("Can't connect to DB %+v", err)
	}

	j := NewJob()
	j.Message = "test message"
	j.Finished = false

	j2 := NewJob()
	j2.Message = "test message"
	j2.Finished = true

	x, err := col.InsertOne(context.TODO(), j)
	if err != nil {
		t.Fatal(err)
	}

	x, err = col.InsertOne(context.TODO(), j2)
	if err != nil {
		t.Fatal(err)
	}

	r := NewMongoRepository(col)
	y, err := r.GetUnfinishedJobsCtx(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(j.ID, y.ID) {
		t.Errorf("expected %s, got %s", y.ID, j.ID)
	}

	err = col.Drop(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}

func TestMongoRepository_GetJobByIdCtx(t *testing.T) {
	col, err := getDBCollection()
	if err != nil {
		t.Fatalf("Can't connect to DB %+v", err)
	}

	j := NewJob()
	j.Message = "test message"
	j.Finished = false
	j.LastTried = time.Now()

	_, err = col.InsertOne(context.TODO(), j)
	if err != nil {
		t.Fatal(err)
	}

	r := NewMongoRepository(col)
	y, err := r.GetJobByIdCtx(context.TODO(), j.ID)
	if err != nil {
		t.Fatal(err)
	}
	opt := cmpopts.IgnoreFields(Job{}, "LastTried")
	if !cmp.Equal(&j, y, opt) {
		fmt.Println(cmp.Diff(j, y))
		t.Errorf("expected %s, got %s", y.ID, j.ID)
	}

	err = col.Drop(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}

func TestMongoRepository_InsertJobCtx(t *testing.T) {
	col, err := getDBCollection()
	if err != nil {
		t.Fatalf("Can't connect to DB %+v", err)
	}
	r := NewMongoRepository(col)

	j := NewJob()
	j.Message = "test message"
	j.Finished = false
	j.LastTried = time.Now()

	id, err := r.InsertJobCtx(context.TODO(), &j)
	if err != nil {
		t.Fatal(err)
	}

	jobFromDB, err := r.GetJobByIdCtx(context.TODO(), *id)
	if err != nil {
		t.Fatal(err)
	}

	opt := cmpopts.IgnoreFields(Job{}, "LastTried")
	if !cmp.Equal(&j, jobFromDB, opt) {
		fmt.Println(cmp.Diff(j, jobFromDB))
		t.Errorf("expected %s, got %s", jobFromDB.ID, j.ID)
	}

	err = col.Drop(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}

func TestMongoRepository_UpdateJobCtx(t *testing.T) {
	col, err := getDBCollection()
	if err != nil {
		t.Fatalf("Can't connect to DB %+v", err)
	}

	r := NewMongoRepository(col)

	j := NewJob()
	fmt.Println(j)
	j.Message = "test message"
	j.Finished = false
	j.LastTried = time.Now()

	id, err := r.InsertJobCtx(context.TODO(), &j)
	if err != nil {
		t.Fatal(err)
	}

	jobFromDB, err := r.GetJobByIdCtx(context.TODO(), *id)
	if err != nil {
		t.Fatal(err)
	}

	jobFromDB.Result = "done"
	jobFromDB.Finished = true

	err = r.UpdateJobCtx(context.TODO(), jobFromDB)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := r.GetJobByIdCtx(context.TODO(), jobFromDB.ID)
	if err != nil {
		t.Fatal(err)
	}

	if updated.Finished != true {
		t.Error("expected job to be finished")
	}

	opt := cmpopts.IgnoreFields(Job{}, "LastTried")
	if !cmp.Equal(updated, jobFromDB, opt) {
		fmt.Println(cmp.Diff(j, jobFromDB))
		t.Errorf("expected %s, got %s", jobFromDB.ID, j.ID)
	}

	err = col.Drop(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}
