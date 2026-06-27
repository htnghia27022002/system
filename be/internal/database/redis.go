package database

import (
	"context"
	"fmt"
	"time"

	"be/internal/config"

	goredis "github.com/redis/go-redis/v9"
)

const redisPingTimeout = 5 * time.Second

// ConnectRedis dials Redis when cfg.RedisURL is set. Returns (nil, nil) when URL is empty.
func ConnectRedis(cfg config.Config) (*goredis.Client, error) {
	return ConnectRedisURL(cfg.RedisURL)
}

// ConnectRedisURL dials Redis from a URL. Returns (nil, nil) when url is empty.
func ConnectRedisURL(url string) (*goredis.Client, error) {
	if url == "" {
		return nil, nil
	}

	opts, err := goredis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("redis: parse url: %w", err)
	}

	client := goredis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), redisPingTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis: ping: %w", err)
	}

	return client, nil
}
