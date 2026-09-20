package handlers

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	apperrors "be/common/errors"
	"be/common/httpx"
	"be/common/response"
	webhookdto "be/internal/dto/webhook"
	webhooksvc "be/internal/services/webhook"
)

type WebhookHandler struct {
	svc *webhooksvc.Service
}

func NewWebhookHandler(svc *webhooksvc.Service) *WebhookHandler {
	return &WebhookHandler{svc: svc}
}

func (h *WebhookHandler) GetInbox(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	result, err := h.svc.GetOrCreateInbox(c.Request.Context(), userID)
	httpx.OK(c, result, err)
}

func (h *WebhookHandler) Regenerate(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	result, err := h.svc.RegenerateUUID(c.Request.Context(), userID)
	httpx.OK(c, result, err)
}

func (h *WebhookHandler) ListRequests(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	var form webhookdto.ListRequestsQuery
	_ = c.ShouldBindQuery(&form)
	result, err := h.svc.ListRequests(c.Request.Context(), userID, form)
	httpx.OK(c, result, err)
}

func (h *WebhookHandler) GetRequest(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	result, err := h.svc.GetRequest(c.Request.Context(), userID, c.Param("id"))
	httpx.OK(c, result, err)
}

func (h *WebhookHandler) SoftDeleteAllRequests(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	result, err := h.svc.SoftDeleteAllActive(c.Request.Context(), userID)
	httpx.OK(c, result, err)
}

func (h *WebhookHandler) SoftDeleteRequest(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	result, err := h.svc.SoftDeleteRequest(c.Request.Context(), userID, c.Param("id"))
	httpx.OK(c, result, err)
}

func (h *WebhookHandler) SetRequestRead(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	var body webhookdto.SetReadRequest
	if !httpx.BindJSON(c, &body) {
		return
	}
	result, err := h.svc.SetRequestRead(c.Request.Context(), userID, c.Param("id"), body.IsRead)
	httpx.OK(c, result, err)
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

	form := map[string]any{}
	ct := c.ContentType()
	if strings.Contains(ct, "application/x-www-form-urlencoded") || strings.Contains(ct, "multipart/form-data") {
		_ = c.Request.ParseMultipartForm(webhooksvc.MaxBodyBytes)
		_ = c.Request.ParseForm()
		form = httpx.ValuesMap(c.Request.PostForm)
	}

	fullURL := c.Request.URL.RequestURI()
	if fullURL == "" {
		fullURL = c.Request.URL.Path
	}

	captureErr := h.svc.Capture(c.Request.Context(), publicUUID, webhooksvc.CaptureInput{
		Method:      c.Request.Method,
		URL:         fullURL,
		ClientIP:    httpx.ClientIP(c),
		Headers:     httpx.HeaderMap(c.Request.Header),
		Query:       httpx.QueryMap(c.Request.URL.Query()),
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
