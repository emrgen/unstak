package cache

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis() (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})

	return &Redis{client: client}, nil
}

func (r *Redis) Set(ctx context.Context, k string, v interface{}, ttl time.Duration) error {
	r.client.Set(ctx, k, v, ttl)

	return nil
}

func (r *Redis) Get(ctx context.Context, k string) (interface{}, error) {
	value, err := r.client.Get(ctx, k).Result()
	if err != nil {
		return nil, err
	}

	return value, nil
}
