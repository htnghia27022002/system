package database

import (
	"be/internal/config"
	pkgredis "be/pkg/redis"

	goredis "github.com/redis/go-redis/v9"
)

// ConnectRedis dials Redis when cfg.RedisURL is set. Returns (nil, nil) when URL is empty.
func ConnectRedis(cfg config.Config) (*goredis.Client, error) {
	return ConnectRedisURL(cfg.RedisURL)
}

// ConnectRedisURL dials Redis from a URL. Returns (nil, nil) when url is empty.
func ConnectRedisURL(url string) (*goredis.Client, error) {
	return pkgredis.Connect(url)
}
