package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type Entity struct {
	ID     primitive.ObjectID `bson:"_id" db:"id"`
	Name   *string            `db:"name"`
	Email  *string            `db:"email"`
	Status *string            `db:"status"`
}
