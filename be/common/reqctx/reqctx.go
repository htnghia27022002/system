// Package reqctx holds per-request values shared by middleware, handlers and the
// response helpers. It imports only gin, so any package can depend on it.
package reqctx

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// RequestIDKey is the gin context key holding the request ID.
const RequestIDKey = "request_id"

// RequestID returns the current request ID, or "" when none was set.
func RequestID(c *gin.Context) string {
	return c.GetString(RequestIDKey)
}

// IsHTTPS reports whether the request reached us over HTTPS, directly or via a
// TLS-terminating proxy that sets X-Forwarded-Proto.
func IsHTTPS(c *gin.Context) bool {
	return c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}
