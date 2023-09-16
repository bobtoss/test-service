package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"test-service/internal/domain/user"
	"test-service/pkg/store"
)

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		db: db.Collection("users"),
	}
}

func (r *UserRepository) List(ctx context.Context, req *http.Request) (dest []user.Entity, err error) {
	filters := r.prepareFilters(req)
	cur, err := r.db.Find(ctx, filters)
	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &dest); err != nil {
		return nil, err
	}

	return
}

func (r *UserRepository) Add(ctx context.Context, data user.Entity) (id primitive.ObjectID, err error) {
	res, err := r.db.InsertOne(ctx, data)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *UserRepository) Get(ctx context.Context, id primitive.ObjectID) (dest user.Entity, err error) {
	err = r.db.FindOne(ctx, bson.D{{"_id", id}}).Decode(&dest)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = store.ErrorNotFound
		}
	}
	fmt.Println(dest, "mongo")
	return dest, err
}

func (r *UserRepository) Update(ctx context.Context, id primitive.ObjectID, data user.Entity) (dest user.Entity, err error) {
	args := r.prepareArgs(data)
	if len(args) > 0 {
		if err = r.db.FindOne(ctx, bson.D{{"_id", id}}).Decode(&dest); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				err = store.ErrorNotFound
			}
		}
		_, err = r.db.UpdateOne(ctx, bson.D{{"_id", id}}, bson.M{"$set": args})
		if err != nil {
			return
		}
	}

	return
}

func (r *UserRepository) prepareFilters(req *http.Request) bson.M {
	filters := bson.M{}

	name := req.URL.Query().Get("name")
	if name != "" {
		filters["name"] = name
	}

	email := req.URL.Query().Get("email")
	if email != "" {
		filters["email"] = email
	}

	status := req.URL.Query().Get("status")
	if status != "" {
		filters["status"] = status
	}

	return filters
}

func (r *UserRepository) prepareArgs(data user.Entity) bson.M {
	args := bson.M{}
	if data.Status != nil {
		args["status"] = data.Status
	}

	return args
}
