package routes

import (
	"github.com/gin-gonic/gin"

	"be/internal/app"
	"be/internal/middleware"
)

func RegisterAdminRoutes(r *gin.RouterGroup, c *app.Container) {
	admin := r.Group("/admin")
	admin.Use(middleware.Auth(c.JWT))
	{
		permissions := admin.Group("/permissions")
		permissions.Use(middleware.RequirePermission("permissions:read"))
		permissions.GET("", c.PermissionHandler.List)

		roles := admin.Group("/roles")
		roles.GET("", middleware.RequirePermission("roles:read"), c.RoleHandler.List)
		roles.POST("", middleware.RequirePermission("roles:create"), c.RoleHandler.Create)
		roles.GET("/:id", middleware.RequirePermission("roles:read"), c.RoleHandler.Get)
		roles.PATCH("/:id", middleware.RequirePermission("roles:update"), c.RoleHandler.Update)
		roles.DELETE("/:id", middleware.RequirePermission("roles:delete"), c.RoleHandler.Delete)

		users := admin.Group("/users")
		users.GET("", middleware.RequirePermission("users:read"), c.UserHandler.List)
		users.POST("", middleware.RequirePermission("users:create"), c.UserHandler.Create)
		users.GET("/:id", middleware.RequirePermission("users:read"), c.UserHandler.Get)
		users.PATCH("/:id", middleware.RequirePermission("users:update"), c.UserHandler.Update)
		users.DELETE("/:id", middleware.RequirePermission("users:delete"), c.UserHandler.Delete)
	}
}
