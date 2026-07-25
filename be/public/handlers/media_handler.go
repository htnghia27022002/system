package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"be/internal/common/response"
	"be/internal/services/media"
)

type MediaHandler struct {
	media *media.Service
}

func NewMediaHandler(mediaSvc *media.Service) *MediaHandler {
	return &MediaHandler{media: mediaSvc}
}

// ServeAvatar streams an uploaded avatar by opaque filename (public GET).
func (h *MediaHandler) ServeAvatar(c *gin.Context) {
	filename := c.Param("filename")
	path := h.media.ResolveAvatarFile(filename)
	if path == "" {
		response.Error(c, http.StatusNotFound, "avatar not found")
		return
	}
	if _, err := os.Stat(path); err != nil {
		response.Error(c, http.StatusNotFound, "avatar not found")
		return
	}
	c.File(path)
}
