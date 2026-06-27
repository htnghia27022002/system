package cache

import (
	"context"
	"time"
)

// Store is the cache backend contract (file, Redis, or noop when disabled).
type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Purge(ctx context.Context) error
	Close() error
}
