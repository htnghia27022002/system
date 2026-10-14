package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"be/common/reqctx"
)

// RequestIDHeader carries the request ID in and out (echoed on every response).
const RequestIDHeader = "X-Request-ID"

var validRequestID = regexp.MustCompile(`^[A-Za-z0-9._:-]{1,128}$`)

// RequestID reuses a well-formed inbound X-Request-ID (e.g. from nginx or a caller)
// or generates one, stores it in the context and echoes it on the response.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if !validRequestID.MatchString(id) {
			id = uuid.NewString()
		}
		c.Set(reqctx.RequestIDKey, id)
		c.Header(RequestIDHeader, id)
		c.Next()
	}
}

// AccessLog writes one structured line per request. Query strings are not logged
// (they may carry tokens); the matched route template is logged alongside the path.
func AccessLog(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		level := slog.LevelInfo
		switch {
		case status >= http.StatusInternalServerError:
			level = slog.LevelError
		case status >= http.StatusBadRequest:
			level = slog.LevelWarn
		}

		attrs := []slog.Attr{
			slog.String("request_id", reqctx.RequestID(c)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("route", c.FullPath()),
			slog.Int("status", status),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
			slog.String("ip", c.ClientIP()),
			slog.Int("bytes", c.Writer.Size()),
		}
		if userID := GetUserID(c); userID != "" {
			attrs = append(attrs, slog.String("user_id", userID))
		}
		logger.LogAttrs(c.Request.Context(), level, "http request", attrs...)
	}
}

// Recovery turns a panic into a logged 500 with a generic body.
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic recovered",
					slog.String("request_id", reqctx.RequestID(c)),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
					slog.String("panic", fmt.Sprint(rec)),
					slog.String("stack", string(debug.Stack())),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}

// SecurityHeaders sets conservative headers for API responses. HSTS is only sent
// when the request arrived over HTTPS (directly or via X-Forwarded-Proto).
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		if reqctx.IsHTTPS(c) {
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}
