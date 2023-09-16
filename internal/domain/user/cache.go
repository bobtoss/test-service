package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cache interface {
	Get(ctx context.Context, id primitive.ObjectID) (dest Entity, err error)
}
