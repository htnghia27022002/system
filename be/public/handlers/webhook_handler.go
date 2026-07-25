package handlers

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	apperrors "be/internal/common/errors"
	"be/internal/common/response"
	webhookdto "be/internal/dto/webhook"
	"be/internal/middleware"
	webhooksvc "be/internal/services/webhook"
)

type WebhookHandler struct {
	svc *webhooksvc.Service
}

func NewWebhookHandler(svc *webhooksvc.Service) *WebhookHandler {
	return &WebhookHandler{svc: svc}
}

func (h *WebhookHandler) GetInbox(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.HandleError(c, apperrors.ErrUnauthorized)
		return
	}
	result, err := h.svc.GetOrCreateInbox(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
}

func (h *WebhookHandler) Regenerate(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.HandleError(c, apperrors.ErrUnauthorized)
		return
	}
	result, err := h.svc.RegenerateUUID(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
}

func (h *WebhookHandler) ListRequests(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.HandleError(c, apperrors.ErrUnauthorized)
		return
	}
	var query webhookdto.ListRequestsQuery
	_ = c.ShouldBindQuery(&query)
	result, err := h.svc.ListRequests(c.Request.Context(), userID, query)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
}

func (h *WebhookHandler) GetRequest(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.HandleError(c, apperrors.ErrUnauthorized)
		return
	}
	result, err := h.svc.GetRequest(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
}

func (h *WebhookHandler) SoftDeleteAllRequests(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.HandleError(c, apperrors.ErrUnauthorized)
		return
	}
	result, err := h.svc.SoftDeleteAllActive(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
}

func (h *WebhookHandler) SoftDeleteRequest(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.HandleError(c, apperrors.ErrUnauthorized)
		return
	}
	result, err := h.svc.SoftDeleteRequest(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
}

func (h *WebhookHandler) SetRequestRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.HandleError(c, apperrors.ErrUnauthorized)
		return
	}
	var body webhookdto.SetReadRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := h.svc.SetRequestRead(c.Request.Context(), userID, c.Param("id"), body.IsRead)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
}

// Capture handles unauthenticated ingest for ANY /api/webhooks/capture/:uuid.
func (h *WebhookHandler) Capture(c *gin.Context) {
	publicUUID := c.Param("uuid")

	body, oversized, err := webhooksvc.ReadBodyLimited(c.Request.Body)
	if err != nil && !errors.Is(err, io.EOF) {
		response.Error(c, http.StatusBadRequest, "failed to read request body")
		return
	}
	_ = c.Request.Body.Close()

	headers := map[string]any{}
	for key, values := range c.Request.Header {
		if len(values) == 1 {
			headers[key] = values[0]
		} else {
			copied := make([]string, len(values))
			copy(copied, values)
			headers[key] = copied
		}
	}

	query := map[string]any{}
	for key, values := range c.Request.URL.Query() {
		if len(values) == 1 {
			query[key] = values[0]
		} else {
			copied := make([]string, len(values))
			copy(copied, values)
			query[key] = copied
		}
	}

	form := map[string]any{}
	ct := c.ContentType()
	if strings.Contains(ct, "application/x-www-form-urlencoded") || strings.Contains(ct, "multipart/form-data") {
		_ = c.Request.ParseMultipartForm(webhooksvc.MaxBodyBytes)
		_ = c.Request.ParseForm()
		for key, values := range c.Request.PostForm {
			if len(values) == 1 {
				form[key] = values[0]
			} else {
				copied := make([]string, len(values))
				copy(copied, values)
				form[key] = copied
			}
		}
	}

	fullURL := c.Request.URL.RequestURI()
	if fullURL == "" {
		fullURL = c.Request.URL.Path
	}

	captureErr := h.svc.Capture(c.Request.Context(), publicUUID, webhooksvc.CaptureInput{
		Method:      c.Request.Method,
		URL:         fullURL,
		ClientIP:    ClientIP(c),
		Headers:     headers,
		Query:       query,
		Form:        form,
		Body:        body,
		ContentType: ct,
		Oversized:   oversized,
	})

	if captureErr != nil {
		if webhooksvc.IsBodyTooLarge(captureErr) {
			response.Error(c, http.StatusRequestEntityTooLarge, webhooksvc.FormatCaptureError(captureErr))
			return
		}
		if apperrors.IsNotFound(captureErr) {
			response.Error(c, http.StatusNotFound, "webhook URL not found")
			return
		}
		response.HandleError(c, captureErr)
		return
	}

	ack := webhookdto.CaptureAckResponse{OK: true, Message: "Request received"}
	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
		return
	}
	response.JSON(c, http.StatusOK, ack)
}

// ClientIP returns best-effort client IP behind nginx proxies.
func ClientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		for _, part := range parts {
			ip := strings.TrimSpace(part)
			if ip != "" {
				return ip
			}
		}
	}
	if xri := strings.TrimSpace(c.GetHeader("X-Real-IP")); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return c.Request.RemoteAddr
}
