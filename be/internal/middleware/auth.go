package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	apperrors "be/internal/common/errors"
	jwtmanager "be/internal/common/jwt"
	"be/internal/common/rbac"
	"be/internal/common/response"
	"be/internal/repository/interfaces"
)

const (
	ContextUserIDKey      = "userID"
	ContextUserEmailKey   = "userEmail"
	ContextUserRoleKey    = "userRole"
	ContextUserRoleIDKey  = "userRoleID"
	ContextPermissionsKey = "permissions"
)

func Auth(jwtManager *jwtmanager.Manager, roleRepo interfaces.RoleRepository) gin.HandlerFunc {
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

		permissions := claims.Permissions
		if roleRepo != nil && claims.RoleID != "" {
			if fresh, err := roleRepo.GetPermissionKeysByRoleID(c.Request.Context(), claims.RoleID); err == nil {
				permissions = fresh
			}
		}
		if permissions == nil {
			permissions = []string{}
		}

		c.Set(ContextUserIDKey, claims.Subject)
		c.Set(ContextUserEmailKey, claims.Email)
		c.Set(ContextUserRoleKey, claims.Role)
		c.Set(ContextUserRoleIDKey, claims.RoleID)
		c.Set(ContextPermissionsKey, permissions)
		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, ok := permissionsFromContext(c)
		if !ok {
			response.HandleError(c, apperrors.ErrForbidden)
			c.Abort()
			return
		}
		if !rbac.Allowed(permissions, permission) {
			response.HandleError(c, apperrors.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireView(resource string) gin.HandlerFunc {
	return RequirePermission(rbac.Key(resource, rbac.ActionView))
}

func RequireModify(resource string) gin.HandlerFunc {
	return RequirePermission(rbac.Key(resource, rbac.ActionModify))
}

func GetUserID(c *gin.Context) string {
	value, _ := c.Get(ContextUserIDKey)
	id, _ := value.(string)
	return id
}

func permissionsFromContext(c *gin.Context) ([]string, bool) {
	raw, ok := c.Get(ContextPermissionsKey)
	if !ok {
		return nil, false
	}
	permissions, ok := raw.([]string)
	if !ok {
		return nil, false
	}
	return permissions, true
}
