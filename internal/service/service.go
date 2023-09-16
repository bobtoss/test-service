package service

import "test-service/internal/domain/user"

type Configuration func(s *Service) error

type Service struct {
	userRepository user.Repository
	userCache      user.Cache
}

func New(configs ...Configuration) (s *Service, err error) {
	s = &Service{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}
	return
}

func WithUserRepository(userRepository user.Repository) Configuration {
	return func(s *Service) error {
		s.userRepository = userRepository
		return nil
	}
}

func WithUserCache(userCache user.Cache) Configuration {
	return func(s *Service) error {
		s.userCache = userCache
		return nil
	}
}
