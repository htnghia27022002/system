package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "be/internal/common/errors"
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
	default:
		Error(c, http.StatusInternalServerError, err.Error())
	}
}
