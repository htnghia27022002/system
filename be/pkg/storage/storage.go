// Package storage is a small object-store abstraction for user uploads
// (avatars, attachments). Drivers: local filesystem and S3-compatible
// (AWS S3, Cloudflare R2, MinIO, …). Objects are addressed by slash-separated
// keys such as "avatars/<uuid>.png".
package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
)

// ErrNotFound is returned by Open for a missing object.
var ErrNotFound = errors.New("storage: object not found")

// Object describes a stored object.
type Object struct {
	Key         string
	Size        int64
	ContentType string
}

// Storage stores public objects. Do not put private files behind a public URL:
// use a separate bucket/prefix with signed URLs for those.
type Storage interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	// Open streams an object. Callers must close the reader.
	Open(ctx context.Context, key string) (io.ReadCloser, Object, error)
	Delete(ctx context.Context, key string) error
	// URL is the public URL clients use to fetch the object.
	URL(key string) string
	// KeyFromURL reverses URL; ok is false for URLs this store did not issue.
	KeyFromURL(url string) (key string, ok bool)
}

// CleanKey validates a key: relative, slash-separated, no "..", no empty parts.
func CleanKey(key string) (string, error) {
	key = strings.TrimSpace(key)
	if key == "" || strings.HasPrefix(key, "/") || strings.Contains(key, "\\") {
		return "", fmt.Errorf("storage: invalid key %q", key)
	}
	cleaned := path.Clean(key)
	if cleaned != key || cleaned == "." || strings.HasPrefix(cleaned, "../") || cleaned == ".." {
		return "", fmt.Errorf("storage: invalid key %q", key)
	}
	for _, part := range strings.Split(cleaned, "/") {
		if part == "" || part == "." || part == ".." || strings.HasPrefix(part, ".") {
			return "", fmt.Errorf("storage: invalid key %q", key)
		}
	}
	return cleaned, nil
}

func keyFromPrefixedURL(prefix, url string) (string, bool) {
	if prefix == "" || !strings.HasPrefix(url, prefix+"/") {
		return "", false
	}
	key, err := CleanKey(strings.TrimPrefix(url, prefix+"/"))
	if err != nil {
		return "", false
	}
	return key, true
}
