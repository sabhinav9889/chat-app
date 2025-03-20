package cache

import "context"

// CacheClient defines methods for interacting with a generic cache
type CacheClient interface {
    Set(ctx context.Context, key string, value interface{}) error
    Get(ctx context.Context, key string) (string, error)
    Del(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
}
