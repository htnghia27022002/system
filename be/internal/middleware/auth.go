package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	apperrors "be/common/errors"
	jwtmanager "be/common/jwt"
	"be/common/rbac"
	"be/common/response"
	"be/internal/repository/interfaces"
)

const (
	ContextUserIDKey      = "userID"
	ContextUserEmailKey   = "userEmail"
	ContextUserRoleKey    = "userRole"
	ContextUserRoleIDKey  = "userRoleID"
	ContextPermissionsKey = "permissions"
	ContextSuperAdminKey  = "superAdmin"
)

func Auth(jwtManager *jwtmanager.Manager, roleRepo interfaces.RoleRepository, userRepo interfaces.UserRepository) gin.HandlerFunc {
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

		superAdmin := claims.SuperAdmin
		if userRepo != nil && claims.Subject != "" {
			if user, err := userRepo.GetByID(c.Request.Context(), claims.Subject); err == nil && user != nil {
				superAdmin = user.IsSuperAdmin
			}
		}

		c.Set(ContextUserIDKey, claims.Subject)
		c.Set(ContextUserEmailKey, claims.Email)
		c.Set(ContextUserRoleKey, claims.Role)
		c.Set(ContextUserRoleIDKey, claims.RoleID)
		c.Set(ContextPermissionsKey, permissions)
		c.Set(ContextSuperAdminKey, superAdmin)
		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsSuperAdmin(c) {
			c.Next()
			return
		}
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

func IsSuperAdmin(c *gin.Context) bool {
	value, ok := c.Get(ContextSuperAdminKey)
	if !ok {
		return false
	}
	superAdmin, _ := value.(bool)
	return superAdmin
}

func GetPermissions(c *gin.Context) ([]string, bool) {
	return permissionsFromContext(c)
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
