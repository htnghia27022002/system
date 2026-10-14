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

func HandleError(c *gin.Context, err error) {
	switch {
	case apperrors.IsNotFound(err):
		Error(c, http.StatusNotFound, err.Error())
	case apperrors.IsUnauthorized(err):
		Error(c, http.StatusUnauthorized, err.Error())
	case apperrors.IsForbidden(err):
		Error(c, http.StatusForbidden, err.Error())
	case apperrors.IsConflict(err):
		Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, apperrors.ErrBadRequest):
		Error(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, apperrors.ErrTooManyRequests):
		Error(c, http.StatusTooManyRequests, err.Error())
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
