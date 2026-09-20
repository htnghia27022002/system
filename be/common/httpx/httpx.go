package httpx

import (
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	apperrors "be/common/errors"
	"be/common/response"
	"be/internal/middleware"
)

// BindJSON binds the request body. Writes 400 and returns false on failure.
func BindJSON(c *gin.Context, dest any) bool {
	if err := c.ShouldBindJSON(dest); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

// BindQuery binds URL query parameters. Writes 400 and returns false on failure.
func BindQuery(c *gin.Context, dest any) bool {
	if err := c.ShouldBindQuery(dest); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

// RequireUserID returns the JWT user id or writes 401 and false.
func RequireUserID(c *gin.Context) (string, bool) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.HandleError(c, apperrors.ErrUnauthorized)
		return "", false
	}
	return userID, true
}

// OK writes payload at 200, or maps err through HandleError.
func OK(c *gin.Context, payload any, err error) {
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, payload)
}

// Created writes payload at 201, or maps err through HandleError.
func Created(c *gin.Context, payload any, err error) {
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, payload)
}

// NoContent writes 204, or maps err through HandleError.
func NoContent(c *gin.Context, err error) {
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// FormFile opens the first matching multipart field name.
func FormFile(c *gin.Context, errMessage string, names ...string) (multipart.File, *multipart.FileHeader, bool) {
	for _, name := range names {
		file, header, err := c.Request.FormFile(name)
		if err == nil {
			return file, header, true
		}
	}
	if strings.TrimSpace(errMessage) == "" {
		errMessage = "file is required"
	}
	response.Error(c, http.StatusBadRequest, errMessage)
	return nil, nil, false
}

// ValuesMap converts url.Values / headers into a JSON-friendly map.
func ValuesMap(values map[string][]string) map[string]any {
	out := make(map[string]any, len(values))
	for key, items := range values {
		if len(items) == 1 {
			out[key] = items[0]
			continue
		}
		copied := make([]string, len(items))
		copy(copied, items)
		out[key] = copied
	}
	return out
}

// HeaderMap copies an HTTP header into a JSON-friendly map.
func HeaderMap(h http.Header) map[string]any {
	return ValuesMap(h)
}

// QueryMap copies URL query values into a JSON-friendly map.
func QueryMap(v url.Values) map[string]any {
	return ValuesMap(v)
}

// ClientIP returns a best-effort client address behind reverse proxies.
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
