package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"test-service/internal/domain/user"
	"time"
)

type AuthorCache struct {
	cache      *redis.Client
	repository user.Repository
}

func NewAuthorCache(c *redis.Client, r user.Repository) *AuthorCache {
	return &AuthorCache{
		cache:      c,
		repository: r,
	}
}

func (c *AuthorCache) Get(ctx context.Context, id primitive.ObjectID) (dest user.Entity, err error) {
	data, err := c.cache.Get(ctx, id.Hex()).Result()
	fmt.Println("cache", err)
	if err == nil {
		// Data found in cache, unmarshal JSON into struct
		if err = json.Unmarshal([]byte(data), &dest); err != nil {
			return
		}
		fmt.Println(dest, "found")
		return
	}

	// Data not found in cache, retrieve it from the data source
	dest, err = c.repository.Get(ctx, id)
	fmt.Println(err, "not found")
	if err != nil {
		return
	}

	// Marshal struct data into JSON and database it in Redis cache
	payload, err := json.Marshal(dest)
	if err != nil {
		fmt.Println(err, "error in marshal")
		return
	}

	if err = c.cache.Set(ctx, id.Hex(), payload, 5*time.Minute).Err(); err != nil {
		fmt.Println(err, "error in set")
		return
	}
	fmt.Println(dest)
	return
}
