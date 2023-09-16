package repository

import (
	"test-service/internal/domain/user"
	"test-service/internal/repository/mongo"
	"test-service/pkg/store"
)

type Configuration func(r *Repository) error

type Repository struct {
	mongo store.Mongo

	User user.Repository
}

func New(configs ...Configuration) (s *Repository, err error) {
	// Create the repository
	s = &Repository{}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		// Pass the repository into the configuration function
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func (r *Repository) Close() {
	if r.mongo.Client != nil {
		r.mongo.Client.Disconnect(nil)
	}
}

func WithMongoStore(uri, name string) Configuration {
	return func(s *Repository) (err error) {
		// Create the mongo store, if we needed parameters, such as connection strings they could be inputted here
		s.mongo, err = store.NewMongo(uri)
		if err != nil {
			return
		}
		database := s.mongo.Client.Database(name)

		s.User = mongo.NewUserRepository(database)

		return
	}
}
