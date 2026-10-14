package response

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "be/common/errors"
	"be/common/reqctx"
)

func JSON(c *gin.Context, status int, payload any) {
	c.JSON(status, payload)
}

// Error writes {"error": message, "code": <default code for status>}.
func Error(c *gin.Context, status int, message string) {
	JSON(c, status, gin.H{"error": message, "code": DefaultCode(status)})
}

// Stable codes clients can branch on or translate (fe: errors.<code>). Specific
// errors add their own via apperrors.WithCode (e.g. "email_not_verified").
const (
	CodeBadRequest       = "bad_request"
	CodeValidationFailed = "validation_failed"
	CodeInvalidJSON      = "invalid_json"
	CodeUnauthorized     = "unauthorized"
	CodeForbidden        = "forbidden"
	CodeNotFound         = "not_found"
	CodeConflict         = "conflict"
	CodePayloadTooLarge  = "payload_too_large"
	CodeTooManyRequests  = "too_many_requests"
	CodeInternal         = "internal"
)

// DefaultCode maps an HTTP status to its generic error code.
func DefaultCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return CodeBadRequest
	case http.StatusUnauthorized:
		return CodeUnauthorized
	case http.StatusForbidden:
		return CodeForbidden
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusConflict:
		return CodeConflict
	case http.StatusRequestEntityTooLarge:
		return CodePayloadTooLarge
	case http.StatusTooManyRequests:
		return CodeTooManyRequests
	default:
		if status >= 500 {
			return CodeInternal
		}
		return CodeBadRequest
	}
}

// clientError writes {"error": message, "code": ...}; the code is err's own
// (apperrors.WithCode) or the status default.
func clientError(c *gin.Context, status int, err error) {
	code := apperrors.Code(err)
	if code == "" {
		code = DefaultCode(status)
	}
	JSON(c, status, gin.H{"error": err.Error(), "code": code})
}

// ValidationError writes 400 {"error", "code": "validation_failed", "fields":
// {jsonField: rule}} so clients can highlight fields and translate rules.
func ValidationError(c *gin.Context, message string, fields map[string]string) {
	JSON(c, http.StatusBadRequest, gin.H{"error": message, "code": CodeValidationFailed, "fields": fields})
}

func HandleError(c *gin.Context, err error) {
	switch {
	case apperrors.IsNotFound(err):
		clientError(c, http.StatusNotFound, err)
	case apperrors.IsUnauthorized(err):
		clientError(c, http.StatusUnauthorized, err)
	case apperrors.IsForbidden(err):
		clientError(c, http.StatusForbidden, err)
	case apperrors.IsConflict(err):
		clientError(c, http.StatusConflict, err)
	case errors.Is(err, apperrors.ErrBadRequest):
		clientError(c, http.StatusBadRequest, err)
	case errors.Is(err, apperrors.ErrTooManyRequests):
		clientError(c, http.StatusTooManyRequests, err)
	default:
		// Never echo unexpected errors (SQL, driver, upstream) to clients.
		slog.ErrorContext(c.Request.Context(), "internal error",
			slog.String("request_id", reqctx.RequestID(c)),
			slog.String("method", c.Request.Method),
			slog.String("route", c.FullPath()),
			slog.String("error", err.Error()),
		)
		Error(c, http.StatusInternalServerError, "internal server error")
	}
}
