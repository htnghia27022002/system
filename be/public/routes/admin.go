package routes

import (
	"github.com/gin-gonic/gin"

	"be/internal/app"
	"be/internal/middleware"
)

func RegisterAdminRoutes(r *gin.RouterGroup, c *app.Container) {
	admin := r.Group("/admin")
	admin.Use(middleware.Auth(c.JWT, c.RoleRepo))
	{
		permissions := admin.Group("/permissions")
		permissions.GET("", middleware.RequireView("permissions"), c.PermissionHandler.List)

		roles := admin.Group("/roles")
		roles.GET("", middleware.RequireView("roles"), c.RoleHandler.List)
		roles.POST("", middleware.RequireModify("roles"), c.RoleHandler.Create)
		roles.GET("/:id", middleware.RequireView("roles"), c.RoleHandler.Get)
		roles.PATCH("/:id", middleware.RequireModify("roles"), c.RoleHandler.Update)
		roles.DELETE("/:id", middleware.RequireModify("roles"), c.RoleHandler.Delete)

		users := admin.Group("/users")
		users.GET("", middleware.RequireView("users"), c.UserHandler.List)
		users.POST("", middleware.RequireModify("users"), c.UserHandler.Create)
		users.GET("/:id", middleware.RequireView("users"), c.UserHandler.Get)
		users.PATCH("/:id", middleware.RequireModify("users"), c.UserHandler.Update)
		users.DELETE("/:id", middleware.RequireModify("users"), c.UserHandler.Delete)
	}
}
