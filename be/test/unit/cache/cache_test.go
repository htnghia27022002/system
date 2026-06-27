package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"be/internal/common/cache"
	"be/internal/config"
)

func TestFileStoreSetGetDelete(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	store, err := cache.New(config.CacheConfig{
		Enabled:    true,
		Driver:     cache.DriverFile,
		FileDir:    dir,
		DefaultTTL: time.Minute,
	}, nil)
	if err != nil {
		t.Fatalf("new file store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	ctx := context.Background()
	if err := store.Set(ctx, "user:1", []byte("alpha"), 0); err != nil {
		t.Fatalf("set: %v", err)
	}

	value, err := store.Get(ctx, "user:1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if string(value) != "alpha" {
		t.Fatalf("unexpected value: %q", value)
	}

	if err := store.Delete(ctx, "user:1"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = store.Get(ctx, "user:1")
	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFileStoreExpiry(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	store, err := cache.New(config.CacheConfig{
		Enabled:    true,
		Driver:     cache.DriverFile,
		FileDir:    dir,
		DefaultTTL: time.Minute,
	}, nil)
	if err != nil {
		t.Fatalf("new file store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	ctx := context.Background()
	if err := store.Set(ctx, "expiring", []byte("x"), time.Millisecond); err != nil {
		t.Fatalf("set: %v", err)
	}

	time.Sleep(2 * time.Millisecond)

	_, err = store.Get(ctx, "expiring")
	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected expired key to be missing, got %v", err)
	}
}

func TestPackageAPIWhenDisabled(t *testing.T) {
	t.Parallel()

	if err := cache.Init(config.CacheConfig{Enabled: false}, nil); err != nil {
		t.Fatalf("init: %v", err)
	}
	t.Cleanup(func() { _ = cache.Close() })

	ctx := context.Background()
	if err := cache.Set(ctx, "k", []byte("v"), time.Minute); err != nil {
		t.Fatalf("set: %v", err)
	}

	_, err := cache.Get(ctx, "k")
	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound from noop cache, got %v", err)
	}
}

func TestPurgeClearsFileStore(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	store, err := cache.New(config.CacheConfig{
		Enabled:    true,
		Driver:     cache.DriverFile,
		FileDir:    dir,
		DefaultTTL: time.Minute,
	}, nil)
	if err != nil {
		t.Fatalf("new file store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	ctx := context.Background()
	if err := store.Set(ctx, "a", []byte("1"), time.Minute); err != nil {
		t.Fatalf("set: %v", err)
	}
	if err := store.Purge(ctx); err != nil {
		t.Fatalf("purge: %v", err)
	}

	_, err = store.Get(ctx, "a")
	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected purge to clear entries, got %v", err)
	}
}
