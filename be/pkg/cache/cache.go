package cache

import (
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const (
	DriverFile  = "file"
	DriverRedis = "redis"
)

// Options configures a cache Store without depending on app config types.
type Options struct {
	Enabled    bool
	Driver     string
	DefaultTTL time.Duration
	FileDir    string
}

// NewNoop returns a no-op store that always misses.
func NewNoop() Store {
	return noopStore{}
}

// New builds a cache store from Options. Returns a noop store when disabled.
// When driver is redis, pass a shared client from be/pkg/redis (or equivalent).
func New(opts Options, redis *goredis.Client) (Store, error) {
	if !opts.Enabled {
		return NewNoop(), nil
	}

	switch opts.Driver {
	case DriverFile:
		return newFileStore(opts.FileDir, opts.DefaultTTL)
	case DriverRedis:
		if redis == nil {
			return nil, fmt.Errorf("cache redis: pass a connected redis client")
		}
		return newRedisStore(redis, opts.DefaultTTL), nil
	default:
		return nil, fmt.Errorf("cache: unsupported driver %q", opts.Driver)
	}
}
