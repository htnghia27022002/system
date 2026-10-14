package storage

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

// Local stores objects under Root and serves them at PublicPrefix (the API
// route that streams them, e.g. "/api/media").
type Local struct {
	Root         string
	PublicPrefix string
}

var _ Storage = (*Local)(nil)

func NewLocal(root, publicPrefix string) *Local {
	return &Local{Root: root, PublicPrefix: strings.TrimRight(publicPrefix, "/")}
}

func (l *Local) path(key string) (string, error) {
	key, err := CleanKey(key)
	if err != nil {
		return "", err
	}
	return filepath.Join(l.Root, filepath.FromSlash(key)), nil
}

func (l *Local) Put(_ context.Context, key string, r io.Reader, _ int64, _ string) error {
	dest, err := l.path(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	// Write to a temp file and rename so readers never see a partial object.
	tmp, err := os.CreateTemp(filepath.Dir(dest), ".upload-*")
	if err != nil {
		return err
	}
	if _, err := io.Copy(tmp, r); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), dest)
}

func (l *Local) Open(_ context.Context, key string) (io.ReadCloser, Object, error) {
	p, err := l.path(key)
	if err != nil {
		return nil, Object{}, ErrNotFound
	}
	f, err := os.Open(p)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, Object{}, ErrNotFound
	}
	if err != nil {
		return nil, Object{}, err
	}
	info, err := f.Stat()
	if err != nil || info.IsDir() {
		f.Close()
		return nil, Object{}, ErrNotFound
	}
	return f, Object{Key: key, Size: info.Size(), ContentType: mime.TypeByExtension(filepath.Ext(p))}, nil
}

func (l *Local) Delete(_ context.Context, key string) error {
	p, err := l.path(key)
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

func (l *Local) URL(key string) string { return l.PublicPrefix + "/" + key }

func (l *Local) KeyFromURL(url string) (string, bool) { return keyFromPrefixedURL(l.PublicPrefix, url) }
