package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"be/internal/config"
)

const (
	DriverFile  = "file"
	DriverRedis = "redis"
)

var (
	mu       sync.RWMutex
	instance Store = noopStore{}
)

type redisStore struct {
	client     *goredis.Client
	defaultTTL time.Duration
}

func newRedisStore(client *goredis.Client, defaultTTL time.Duration) Store {
	return &redisStore{client: client, defaultTTL: defaultTTL}
}

func (s *redisStore) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return value, nil
}

func (s *redisStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = s.defaultTTL
	}
	if ttl <= 0 {
		return s.client.Set(ctx, key, value, 0).Err()
	}
	return s.client.Set(ctx, key, value, ttl).Err()
}

func (s *redisStore) Delete(ctx context.Context, key string) error {
	err := s.client.Del(ctx, key).Err()
	if errors.Is(err, goredis.Nil) {
		return nil
	}
	return err
}

func (s *redisStore) Purge(ctx context.Context) error {
	return s.client.FlushDB(ctx).Err()
}

func (s *redisStore) Close() error {
	return nil
}

// New builds a cache store from config. Returns noop store when cache is disabled.
// When driver is redis, pass the shared client from database.ConnectRedis when available.
func New(cfg config.CacheConfig, redis *goredis.Client) (Store, error) {
	if !cfg.Enabled {
		return noopStore{}, nil
	}

	switch cfg.Driver {
	case DriverFile:
		return newFileStore(cfg.FileDir, cfg.DefaultTTL)
	case DriverRedis:
		if redis == nil {
			return nil, fmt.Errorf("cache redis: connect via database.ConnectRedis and pass the client")
		}
		return newRedisStore(redis, cfg.DefaultTTL), nil
	default:
		return nil, fmt.Errorf("cache: unsupported driver %q", cfg.Driver)
	}
}

// Init replaces the package-default store. Call once at startup (e.g. from api.go).
func Init(cfg config.CacheConfig, redis *goredis.Client) error {
	store, err := New(cfg, redis)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	_ = instance.Close()
	instance = store
	return nil
}

// Close shuts down the active store and resets to noop.
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if err := instance.Close(); err != nil {
		return err
	}
	instance = noopStore{}
	return nil
}

// Default returns the active store (noop when cache is disabled).
func Default() Store {
	mu.RLock()
	defer mu.RUnlock()
	return instance
}

// Get reads a value from the default store.
func Get(ctx context.Context, key string) ([]byte, error) {
	return Default().Get(ctx, key)
}

// Set writes a value to the default store. ttl=0 uses configured default TTL.
func Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return Default().Set(ctx, key, value, ttl)
}

// Delete removes a key from the default store.
func Delete(ctx context.Context, key string) error {
	return Default().Delete(ctx, key)
}

// Purge clears all entries in the default store.
func Purge(ctx context.Context) error {
	return Default().Purge(ctx)
}
