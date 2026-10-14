package routes

import (
	"github.com/gin-gonic/gin"

	"be/internal/app"
	"be/public/handlers"
)

// RegisterAppRoutes mounts public, unauthenticated endpoints: deployment info and
// the liveness/readiness probes.
func RegisterAppRoutes(r *gin.RouterGroup, c *app.Container, health *handlers.HealthHandler) {
	h := handlers.NewAppHandler(c.Modules.Names())
	r.GET("/app/info", h.Info)
	r.GET("/health", health.Health)
	r.GET("/ready", health.Ready)
}
