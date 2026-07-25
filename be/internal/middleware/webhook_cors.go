package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const webhookCapturePathPrefix = "/api/webhooks/capture/"

// WebhookCaptureCORS allows unauthenticated cross-origin capture without credentials.
// It short-circuits OPTIONS preflight so credentialed owner CORS does not reject unknown origins.
func WebhookCaptureCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, webhookCapturePathPrefix) {
			c.Next()
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Max-Age", "86400")
		// Explicitly omit Allow-Credentials for public capture.

		// Only short-circuit browser CORS preflight. Bare OPTIONS (e.g. curl)
		// continues to Capture so FR-004 can record the method.
		if c.Request.Method == http.MethodOptions &&
			c.GetHeader("Access-Control-Request-Method") != "" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
