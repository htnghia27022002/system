package routes

import (
	"github.com/gin-gonic/gin"

	"be/internal/app"
	"be/public/handlers"
)

// RegisterAppRoutes mounts public deployment info under /api/app (no auth).
func RegisterAppRoutes(r *gin.RouterGroup, c *app.Container) {
	h := handlers.NewAppHandler(c.Modules.Names())
	r.GET("/app/info", h.Info)
}
