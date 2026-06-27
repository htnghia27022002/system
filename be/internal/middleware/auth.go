package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	apperrors "be/internal/common/errors"
	jwtmanager "be/internal/common/jwt"
	"be/internal/common/response"
)

const (
	ContextUserIDKey      = "userID"
	ContextUserEmailKey   = "userEmail"
	ContextUserRoleKey    = "userRole"
	ContextUserRoleIDKey  = "userRoleID"
	ContextPermissionsKey = "permissions"
)

func Auth(jwtManager *jwtmanager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.HandleError(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.HandleError(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := jwtManager.VerifyAccessToken(parts[1])
		if err != nil {
			response.HandleError(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.Subject)
		c.Set(ContextUserEmailKey, claims.Email)
		c.Set(ContextUserRoleKey, claims.Role)
		c.Set(ContextUserRoleIDKey, claims.RoleID)
		c.Set(ContextPermissionsKey, claims.Permissions)
		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, ok := c.Get(ContextPermissionsKey)
		if !ok {
			response.HandleError(c, apperrors.ErrForbidden)
			c.Abort()
			return
		}
		permissions, ok := raw.([]string)
		if !ok {
			response.HandleError(c, apperrors.ErrForbidden)
			c.Abort()
			return
		}
		for _, item := range permissions {
			if item == permission {
				c.Next()
				return
			}
		}
		response.HandleError(c, apperrors.ErrForbidden)
		c.Abort()
	}
}

func GetUserID(c *gin.Context) string {
	value, _ := c.Get(ContextUserIDKey)
	id, _ := value.(string)
	return id
}
