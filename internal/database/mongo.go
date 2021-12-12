package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"make-me-reliable/internal/config"
)

func getMongoUrl(uri string) string {
	return fmt.Sprintf("mongodb://%s:27017/", uri)
}

// MongoConnectionFromConfigCtx connects to Mongo database using config struct
func MongoConnectionFromConfigCtx(ctx context.Context, c config.Config) (*mongo.Client, error) {
	credential := options.Credential{
		Username:   c.DatabaseUser,
		Password:   c.DatabasePass,
		AuthSource: c.DatabaseName,
	}

	clientOpts := options.
		Client().
		ApplyURI(getMongoUrl(c.DatabaseURL)).
		SetAuth(credential)

	client, err := mongo.Connect(ctx, clientOpts)
	return client, err
}
