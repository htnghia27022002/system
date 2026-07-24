package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type fileEntry struct {
	ExpiresAt time.Time `json:"expiresAt"`
	Value     []byte    `json:"value"`
}

type fileStore struct {
	dir        string
	defaultTTL time.Duration
	mu         sync.RWMutex
	closed     bool
}

func newFileStore(dir string, defaultTTL time.Duration) (*fileStore, error) {
	if dir == "" {
		dir = "storage/cache"
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("cache file dir: %w", err)
	}
	return &fileStore{dir: dir, defaultTTL: defaultTTL}, nil
}

func (s *fileStore) Get(_ context.Context, key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return nil, ErrClosed
	}

	path, err := s.entryPath(key)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var entry fileEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		_ = os.Remove(path)
		return nil, ErrNotFound
	}

	if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		_ = os.Remove(path)
		return nil, ErrNotFound
	}

	return entry.Value, nil
}

func (s *fileStore) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}

	if ttl <= 0 {
		ttl = s.defaultTTL
	}

	entry := fileEntry{Value: append([]byte(nil), value...)}
	if ttl > 0 {
		entry.ExpiresAt = time.Now().Add(ttl)
	}

	payload, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	path, err := s.entryPath(key)
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, payload, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (s *fileStore) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}

	path, err := s.entryPath(key)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return nil
}

func (s *fileStore) Purge(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return err
	}

	for _, item := range entries {
		if item.IsDir() {
			continue
		}
		name := item.Name()
		if filepath.Ext(name) != ".cache" {
			continue
		}
		if err := os.Remove(filepath.Join(s.dir, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func (s *fileStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

func (s *fileStore) entryPath(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("cache: empty key")
	}
	sum := sha256.Sum256([]byte(key))
	return filepath.Join(s.dir, hex.EncodeToString(sum[:])+".cache"), nil
}
