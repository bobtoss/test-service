package service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
	"test-service/internal/domain/user"
	"test-service/pkg/log"
	"test-service/pkg/store"
)

func (s *Service) ListUsers(ctx context.Context, req *http.Request) (res []user.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListAuthors")

	data, err := s.userRepository.List(ctx, req)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}
	res = user.ParseFromEntities(data)

	return
}

func (s *Service) AddUser(ctx context.Context, req user.Request) (res user.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("AddUser")

	data := user.Entity{
		ID:     primitive.NewObjectID(),
		Name:   &req.Name,
		Email:  &req.Email,
		Status: &req.Status,
	}

	data.ID, err = s.userRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}
	res = user.ParseFromEntity(data)
	return
}

func (s *Service) GetUser(ctx context.Context, id primitive.ObjectID) (res user.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetUser").With(zap.String("id", id.Hex()))
	data, err := s.userCache.Get(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}

	fmt.Println(data)
	res = user.ParseFromEntity(data)

	return
}

func (s *Service) UpdateStatus(ctx context.Context, id primitive.ObjectID, req user.Request) (res user.UpdateResponse, err error) {
	logger := log.LoggerFromContext(ctx).Named("UpdateUserStatus").With(zap.String("id", id.Hex()))

	data := user.Entity{
		Status: &req.Status,
	}
	dest, err := s.userRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to update by id", zap.Error(err))
		return
	}

	res.ID = id.Hex()
	res.New = *data.Status
	res.Previous = *dest.Status

	return
}
