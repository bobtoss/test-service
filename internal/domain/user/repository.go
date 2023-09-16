package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Repository interface {
	List(ctx context.Context, req *http.Request) (dest []Entity, err error)
	Add(ctx context.Context, data Entity) (id primitive.ObjectID, err error)
	Get(ctx context.Context, id primitive.ObjectID) (dest Entity, err error)
	Update(ctx context.Context, id primitive.ObjectID, data Entity) (dest Entity, err error)
}
