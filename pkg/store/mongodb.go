package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongodb://username:password@localhost:27017/?retryWrites=true&w=majority&tls=false

const timeout = 10 * time.Second

type Mongo struct {
	Client *mongo.Client
}

func NewMongo(uri string) (store Mongo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	store.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return
	}
	if err = store.Client.Ping(context.Background(), nil); err != nil {
		return
	}

	return
}
