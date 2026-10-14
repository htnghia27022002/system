package routes

import (
	"github.com/gin-gonic/gin"

	"be/internal/app"
	"be/internal/middleware"
)

// RegisterAdminRoutes mounts core admin APIs and returns the authenticated /admin
// group so modules can mount their own admin routes on it.
func RegisterAdminRoutes(r *gin.RouterGroup, c *app.Container) *gin.RouterGroup {
	admin := r.Group("/admin")
	admin.Use(c.Kernel.Auth)
	{
		permissions := admin.Group("/permissions")
		permissions.GET("", middleware.RequireView("permissions"), c.PermissionHandler.List)
		permissions.GET("/all", middleware.RequireView("permissions"), c.PermissionHandler.ListAll)

		roles := admin.Group("/roles")
		roles.GET("", middleware.RequireView("roles"), c.RoleHandler.List)
		roles.GET("/all", middleware.RequireView("roles"), c.RoleHandler.ListAll)
		roles.POST("", middleware.RequireModify("roles"), c.RoleHandler.Create)
		roles.GET("/:id", middleware.RequireView("roles"), c.RoleHandler.Get)
		roles.PATCH("/:id", middleware.RequireModify("roles"), c.RoleHandler.Update)
		roles.DELETE("/:id", middleware.RequireModify("roles"), c.RoleHandler.Delete)

		admin.GET("/dashboard/overview", middleware.RequireView("dashboard"), c.SystemHandler.DashboardOverview)
		admin.GET("/audit-logs", middleware.RequireView("audit"), c.SystemHandler.ListAuditLogs)
		admin.GET("/settings", middleware.RequireView("settings"), c.SystemHandler.ListSettings)
		admin.PATCH("/settings", middleware.RequireModify("settings"), c.SystemHandler.UpdateSettings)

		users := admin.Group("/users")
		users.GET("", middleware.RequireView("users"), c.UserHandler.List)
		users.POST("", middleware.RequireModify("users"), c.UserHandler.Create)
		users.POST("/invite", middleware.RequireModify("users"), c.AccountHandler.InviteUser)
		users.POST("/:id/invite", middleware.RequireModify("users"), c.AccountHandler.ResendInvite)
		users.GET("/:id", middleware.RequireView("users"), c.UserHandler.Get)
		users.PATCH("/:id", middleware.RequireModify("users"), c.UserHandler.Update)
		users.POST("/:id/avatar", middleware.RequireModify("users"), c.UserHandler.UploadAvatar)
		users.DELETE("/:id", middleware.RequireModify("users"), c.UserHandler.Delete)
	}
	return admin
}
