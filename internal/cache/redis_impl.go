package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// RedisCacheClient is an implementation of the CacheClient interface for Redis
type RedisCacheClient struct {
	client *redis.Client
}

// NewRedisCacheClient creates a new RedisCacheClient instance1
func NewRedisCacheClient(addr string) *RedisCacheClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, // example: localhost:6379
	})
	return &RedisCacheClient{client: rdb}
}

// Set stores a key-value pair in Redis
func (r *RedisCacheClient) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

// Get retrieves a value from Redis by key
func (r *RedisCacheClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del deletes a key from Redis
func (r *RedisCacheClient) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (r *RedisCacheClient) Exists(ctx context.Context, key string) (bool, error) {
	res, err := r.client.Exists(ctx, key).Result()
	return res > 0, err
}
