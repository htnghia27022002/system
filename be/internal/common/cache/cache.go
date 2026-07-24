package cache

import (
	"context"
	"sync"
	"time"

	"be/internal/config"
	pkgcache "be/pkg/cache"

	goredis "github.com/redis/go-redis/v9"
)

const (
	DriverFile  = pkgcache.DriverFile
	DriverRedis = pkgcache.DriverRedis
)

var (
	ErrNotFound = pkgcache.ErrNotFound
	ErrClosed   = pkgcache.ErrClosed
)

// Store is the cache backend contract.
type Store = pkgcache.Store

var (
	mu       sync.RWMutex
	instance Store = pkgcache.NewNoop()
)

// OptionsFromConfig maps app cache config to pkg Options.
func OptionsFromConfig(cfg config.CacheConfig) pkgcache.Options {
	return pkgcache.Options{
		Enabled:    cfg.Enabled,
		Driver:     cfg.Driver,
		DefaultTTL: cfg.DefaultTTL,
		FileDir:    cfg.FileDir,
	}
}

// New builds a cache store from config. Returns noop store when cache is disabled.
// When driver is redis, pass the shared client from database.ConnectRedis when available.
func New(cfg config.CacheConfig, redis *goredis.Client) (Store, error) {
	return pkgcache.New(OptionsFromConfig(cfg), redis)
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
	instance = pkgcache.NewNoop()
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
