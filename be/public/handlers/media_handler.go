package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"be/common/response"
	"be/pkg/storage"
)

// MediaHandler streams objects from the configured store. With the local
// driver this is how uploads are served; with S3 clients normally fetch the
// bucket/CDN URL directly.
type MediaHandler struct {
	store storage.Storage
}

func NewMediaHandler(store storage.Storage) *MediaHandler {
	return &MediaHandler{store: store}
}

// Serve streams GET /api/media/*key (public).
func (h *MediaHandler) Serve(c *gin.Context) {
	key := strings.TrimPrefix(c.Param("key"), "/")
	body, obj, err := h.store.Open(c.Request.Context(), key)
	if errors.Is(err, storage.ErrNotFound) {
		response.Error(c, http.StatusNotFound, "file not found")
		return
	}
	if err != nil {
		response.HandleError(c, err)
		return
	}
	defer body.Close()

	contentType := obj.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	// Uploaded content is untrusted: never let a browser run it as a page.
	c.Header("Content-Security-Policy", "default-src 'none'; img-src 'self' data:; style-src 'unsafe-inline'; sandbox")
	c.Header("Cache-Control", "public, max-age=86400")
	c.Header("Content-Type", contentType)
	if obj.Size > 0 {
		c.Header("Content-Length", strconv.FormatInt(obj.Size, 10))
	}
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, body)
}
