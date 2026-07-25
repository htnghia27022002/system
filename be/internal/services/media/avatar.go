package media

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	apperrors "be/internal/common/errors"
)

const (
	MaxAvatarBytes = 2 * 1024 * 1024 // 2 MB
	avatarsSubdir  = "avatars"
)

var allowedAvatarMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// Service stores avatar binaries under UPLOAD_DIR and returns public API paths.
type Service struct {
	uploadDir string
}

func NewService(uploadDir string) *Service {
	return &Service{uploadDir: uploadDir}
}

func (s *Service) AvatarsDir() string {
	return filepath.Join(s.uploadDir, avatarsSubdir)
}

// PublicPath returns the API-relative path stored in avatar_url.
func PublicPath(filename string) string {
	return "/api/media/avatars/" + filename
}

// SaveAvatar validates type/size, writes under UPLOAD_DIR/avatars/, returns public path.
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
			ok = true
		}
	}
	if !ok {
		return "", fmt.Errorf("%w: avatar must be JPEG, PNG, or WebP", apperrors.ErrBadRequest)
	}

	if err := os.MkdirAll(s.AvatarsDir(), 0o755); err != nil {
		return "", err
	}

	filename := uuid.NewString() + ext
	dest := filepath.Join(s.AvatarsDir(), filename)
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return "", err
	}
	return PublicPath(filename), nil
}

// ResolveAvatarFile maps a public path or bare filename to an absolute disk path.
// Returns empty string if invalid / traversal attempt.
func (s *Service) ResolveAvatarFile(nameOrPath string) string {
	name := filepath.Base(strings.TrimSpace(nameOrPath))
	if name == "." || name == ".." || name == "" || strings.Contains(name, "..") {
		return ""
	}
	// Only allow UUID-like names with known extensions.
	ext := strings.ToLower(filepath.Ext(name))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return ""
	}
	full := filepath.Join(s.AvatarsDir(), name)
	// Ensure resolved path stays under avatars dir.
	avatarsAbs, err := filepath.Abs(s.AvatarsDir())
	if err != nil {
		return ""
	}
	fullAbs, err := filepath.Abs(full)
	if err != nil {
		return ""
	}
	if !strings.HasPrefix(fullAbs, avatarsAbs+string(os.PathSeparator)) && fullAbs != avatarsAbs {
		return ""
	}
	return fullAbs
}

// DeleteByPublicPath best-effort removes a previously stored avatar file.
func (s *Service) DeleteByPublicPath(publicPath string) {
	if publicPath == "" {
		return
	}
	path := s.ResolveAvatarFile(publicPath)
	if path == "" {
		return
	}
	_ = os.Remove(path)
}
