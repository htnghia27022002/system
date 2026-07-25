package routes

import (
	"github.com/gin-gonic/gin"

	"be/internal/app"
	"be/internal/middleware"
)

// RegisterWebhookRoutes mounts public capture and JWT owner inbox APIs under /api/webhooks.
func RegisterWebhookRoutes(r *gin.RouterGroup, c *app.Container) {
	webhooks := r.Group("/webhooks")
	{
		// Public capture — no JWT. Product URL /tools/webhooks/{uuid} rewrites here via nginx/Next.
		webhooks.Any("/capture/:uuid", c.WebhookHandler.Capture)

		owner := webhooks.Group("")
		owner.Use(middleware.Auth(c.JWT, c.RoleRepo))
		{
			owner.GET("/inbox", middleware.RequireView("webhooks"), c.WebhookHandler.GetInbox)
			owner.POST("/inbox/regenerate", middleware.RequireModify("webhooks"), c.WebhookHandler.Regenerate)
			owner.GET("/inbox/requests", middleware.RequireView("webhooks"), c.WebhookHandler.ListRequests)
			owner.GET("/inbox/requests/:id", middleware.RequireView("webhooks"), c.WebhookHandler.GetRequest)
			owner.PATCH("/inbox/requests/:id/read", middleware.RequireView("webhooks"), c.WebhookHandler.SetRequestRead)
			// Clear-all must be registered before /:id so Gin does not capture "requests" as an id segment conflict.
			owner.DELETE("/inbox/requests", middleware.RequireModify("webhooks"), c.WebhookHandler.SoftDeleteAllRequests)
			owner.DELETE("/inbox/requests/:id", middleware.RequireModify("webhooks"), c.WebhookHandler.SoftDeleteRequest)
		}
	}
}
