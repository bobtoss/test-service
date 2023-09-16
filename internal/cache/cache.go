package cache

import (
	"test-service/internal/cache/redis"
	"test-service/internal/domain/user"
	"test-service/pkg/store"
)

type Dependencies struct {
	UserRepository user.Repository
}

type Configuration func(r *Cache) error

type Cache struct {
	dependencies Dependencies
	redis        store.Redis

	User user.Cache
}

func New(d Dependencies, configs ...Configuration) (s *Cache, err error) {
	s = &Cache{
		dependencies: d,
	}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func (r *Cache) Close() {
	if r.redis.Connection != nil {
		r.redis.Connection.Close()
	}
}

func WithRedisStore(url string) Configuration {
	return func(s *Cache) (err error) {
		s.redis, err = store.NewRedis(url)
		if err != nil {
			return
		}

		s.User = redis.NewAuthorCache(s.redis.Connection, s.dependencies.UserRepository)

		return
	}
}
