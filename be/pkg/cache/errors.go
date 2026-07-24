package cache

import "errors"

var (
	ErrNotFound = errors.New("cache: key not found")
	ErrClosed   = errors.New("cache: store is closed")
)
