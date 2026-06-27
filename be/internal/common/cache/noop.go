package cache

import (
	"context"
	"time"
)

type noopStore struct{}

func (noopStore) Get(context.Context, string) ([]byte, error) {
	return nil, ErrNotFound
}

func (noopStore) Set(context.Context, string, []byte, time.Duration) error {
	return nil
}

func (noopStore) Delete(context.Context, string) error {
	return nil
}

func (noopStore) Purge(context.Context) error {
	return nil
}

func (noopStore) Close() error {
	return nil
}
