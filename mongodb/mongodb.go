package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// New creates a new connection to mongodb by 'url'
func New(url string, timeout int) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	ctx, _ = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}
