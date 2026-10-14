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

func Error(c *gin.Context, status int, message string) {
	JSON(c, status, gin.H{"error": message})
}

// clientError writes {"error": message} plus {"code": ...} when err carries one.
func clientError(c *gin.Context, status int, err error) {
	body := gin.H{"error": err.Error()}
	if code := apperrors.Code(err); code != "" {
		body["code"] = code
	}
	JSON(c, status, body)
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
