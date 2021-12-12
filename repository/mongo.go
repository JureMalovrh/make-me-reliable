package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository is a definition for Mongo backed repository
type MongoRepository struct {
	c *mongo.Collection
}

// NewMongoRepository returns new MongoRepository
func NewMongoRepository(c *mongo.Collection) *MongoRepository {
	return &MongoRepository{
		c,
	}
}

// GetJobByIdCtx returns Job by ID
func (m MongoRepository) GetJobByIdCtx(ctx context.Context, id primitive.ObjectID) (*Job, error) {
	var j *Job
	res := m.c.FindOne(ctx, bson.M{"_id": id})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return nil, &EmptyResponseError{}
	}
	if res.Err() != nil {
		return nil, res.Err()
	}

	err := res.Decode(&j)
	if err != nil {
		return nil, err
	}

	return j, nil
}

// GetUnfinishedJobsCtx will return job that is not finished and is the oldest in the database
func (m MongoRepository) GetUnfinishedJobsCtx(ctx context.Context) (*Job, error) {
	var j *Job

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: "lastTried", Value: 1}}) // sort by last tried ascending to get last tried one
	res := m.c.FindOne(ctx, bson.M{"finished": false}, findOptions)
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return nil, &EmptyResponseError{}
	}
	if res.Err() != nil {
		return nil, res.Err()
	}

	err := res.Decode(&j)
	if err != nil {
		return nil, err
	}

	return j, nil
}

// InsertJobCtx will insert a new Job into database
func (m MongoRepository) InsertJobCtx(ctx context.Context, job *Job) (*primitive.ObjectID, error) {
	job.LastTried = time.Now()
	res, err := m.c.InsertOne(ctx, job)
	if err != nil {
		return nil, err
	}
	x := res.InsertedID.(primitive.ObjectID) // cast to ObjectId as InsertOne driver returns ObjectId directly
	return &x, nil
}

// UpdateJobCtx will replace existing Job with new Job
func (m MongoRepository) UpdateJobCtx(ctx context.Context, job *Job) error {
	job.LastTried = time.Now()
	res, err := m.c.ReplaceOne(ctx, bson.M{"_id": job.ID}, job)
	if err != nil {
		return err
	}

	if res.ModifiedCount != 1 {
		return fmt.Errorf("more than 1 job updated")
	}

	return nil
}
