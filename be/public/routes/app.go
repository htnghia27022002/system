package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"be/common/response"
	"be/internal/apidoc"
	"be/internal/app"
	"be/internal/module"
	"be/public/handlers"
)

// RegisterAppRoutes mounts public, unauthenticated endpoints: deployment info and
// the liveness/readiness probes.
func RegisterAppRoutes(r *gin.RouterGroup, c *app.Container, health *handlers.HealthHandler) {
	h := handlers.NewAppHandler(c.Modules.Names(), c.Settings)
	r.GET("/app/info", h.Info)
	r.GET("/health", health.Health)
	r.GET("/ready", health.Ready)

	spec, err := json.Marshal(OpenAPI(c.Modules))
	r.GET("/openapi.json", func(ctx *gin.Context) {
		if err != nil {
			response.HandleError(ctx, err)
			return
		}
		ctx.Data(http.StatusOK, "application/json", spec)
	})
}

// OpenAPI is the API document for the core plus the enabled modules.
func OpenAPI(mods module.Set) map[string]any {
	return apidoc.Build(apidoc.Info{
		Title:       "System API",
		Version:     "1.0.0",
		Description: "Core API (auth, users, roles, settings, audit) plus enabled modules. Errors share the Error schema; branch on its code.",
	}, append(CoreOperations(), mods.APIDoc()...))
}
