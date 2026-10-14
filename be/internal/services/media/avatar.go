package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	apperrors "be/common/errors"
	"be/pkg/storage"
)

const (
	MaxAvatarBytes = 2 * 1024 * 1024 // 2 MB
	avatarsPrefix  = "avatars"
	uploadTimeout  = 30 * time.Second
)

var allowedAvatarMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// Service stores uploads through a storage.Storage (local disk or S3) and
// returns the public URL saved in avatar_url.
type Service struct {
	store storage.Storage
}

func NewService(store storage.Storage) *Service {
	return &Service{store: store}
}

// Store exposes the underlying object store (media route, modules).
func (s *Service) Store() storage.Storage { return s.store }

// SaveAvatar validates type/size, stores avatars/<uuid>.<ext> and returns its public URL.
func (s *Service) SaveAvatar(file multipart.File, header *multipart.FileHeader) (string, error) {
	if header == nil {
		return "", fmt.Errorf("%w: avatar file is required", apperrors.ErrBadRequest)
	}
	if header.Size > MaxAvatarBytes {
		return "", fmt.Errorf("%w: avatar must be at most 2 MB", apperrors.ErrBadRequest)
	}

	limited := io.LimitReader(file, MaxAvatarBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return "", err
	}
	if int64(len(data)) > MaxAvatarBytes {
		return "", fmt.Errorf("%w: avatar must be at most 2 MB", apperrors.ErrBadRequest)
	}

	mime := http.DetectContentType(data)
	ext, ok := allowedAvatarMIME[mime]
	if !ok {
		// Some JPEG detectors return image/jpeg; also check Content-Type hint for webp edge cases.
		ct := strings.ToLower(strings.TrimSpace(header.Header.Get("Content-Type")))
		if e, okCT := allowedAvatarMIME[ct]; okCT && (mime == "application/octet-stream" || mime == "image/webp") {
			ext = e
			mime = ct
			ok = true
		}
	}
	if !ok {
		return "", fmt.Errorf("%w: avatar must be JPEG, PNG, or WebP", apperrors.ErrBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), uploadTimeout)
	defer cancel()
	key := avatarsPrefix + "/" + uuid.NewString() + ext
	if err := s.store.Put(ctx, key, bytes.NewReader(data), int64(len(data)), mime); err != nil {
		return "", err
	}
	return s.store.URL(key), nil
}

// DeleteByPublicPath best-effort removes an avatar this store issued. URLs from
// elsewhere (OAuth provider photos, other buckets) are ignored.
func (s *Service) DeleteByPublicPath(publicURL string) {
	key, ok := s.store.KeyFromURL(strings.TrimSpace(publicURL))
	if !ok || !strings.HasPrefix(key, avatarsPrefix+"/") {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), uploadTimeout)
	defer cancel()
	if err := s.store.Delete(ctx, key); err != nil {
		slog.Warn("delete avatar", "key", key, "err", err)
	}
}
