package cache

import (
	"context"
	"errors"
	"time"

	goredis "github.com/redis/go-redis/v9"
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
