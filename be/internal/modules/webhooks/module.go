// Package webhooks is the Webhooks tool module: a per-user public capture URL
// (/api/webhooks/capture/:uuid, no auth) and an owner inbox under /api/webhooks/inbox.
package webhooks

import (
	"io/fs"

	"github.com/gin-gonic/gin"

	"be/common/rbac"
	"be/internal/database/seeders"
	"be/internal/middleware"
	"be/internal/module"
	"be/internal/modules/webhooks/handlers"
	"be/internal/modules/webhooks/migrations"
	"be/internal/modules/webhooks/repository"
	webhooksvc "be/internal/modules/webhooks/services/webhook"
)

const (
	Name     = "webhooks"
	resource = "webhooks"
)

type Module struct {
	kernel  *module.Kernel
	handler *handlers.WebhookHandler
}

var (
	_ module.Module             = (*Module)(nil)
	_ module.PermissionProvider = (*Module)(nil)
	_ module.MigrationProvider  = (*Module)(nil)
	_ module.MiddlewareProvider = (*Module)(nil)
)

func New(k *module.Kernel) module.Module {
	svc := webhooksvc.NewService(
		repository.NewWebhookInboxRepository(k.DB),
		repository.NewWebhookRequestRepository(k.DB),
	)
	return &Module{kernel: k, handler: handlers.NewWebhookHandler(svc)}
}

func (m *Module) Name() string { return Name }

func (m *Module) Migrations() fs.FS { return migrations.FS }

// Permission IDs are stable — they were part of the original catalog.
func (m *Module) Permissions() []seeders.PermissionDef {
	return []seeders.PermissionDef{
		{
			ID:          "10000001-0000-4000-8000-000000000017",
			Key:         rbac.Key(resource, rbac.ActionView),
			Name:        "View webhooks",
			Group:       resource,
			Description: "View webhook inbox and captured requests",
		},
		{
			ID:          "10000001-0000-4000-8000-000000000018",
			Key:         rbac.Key(resource, rbac.ActionModify),
			Name:        "Modify webhooks",
			Group:       resource,
			Description: "Regenerate webhook URL and soft-delete captured requests",
		},
	}
}

// Permissive CORS for public capture must run before the credentialed core CORS
// so browser preflight to /api/webhooks/capture/* succeeds.
func (m *Module) GlobalMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{CaptureCORS()}
}

// RegisterRoutes mounts public capture and JWT owner inbox APIs under /api/webhooks.
// Owner APIs live outside /api/admin by design (tool for any permitted user).
func (m *Module) RegisterRoutes(r module.Routes) {
	h := m.handler
	webhooks := r.API.Group("/webhooks")
	{
		// Public capture — no JWT. Product URL /tools/webhooks/{uuid} rewrites here via nginx/Next.
		webhooks.Any("/capture/:uuid", h.Capture)

		owner := webhooks.Group("")
		owner.Use(m.kernel.Auth)
		{
			owner.GET("/inbox", middleware.RequireView(resource), h.GetInbox)
			owner.POST("/inbox/regenerate", middleware.RequireModify(resource), h.Regenerate)
			owner.GET("/inbox/requests", middleware.RequireView(resource), h.ListRequests)
			owner.GET("/inbox/requests/:id", middleware.RequireView(resource), h.GetRequest)
			owner.PATCH("/inbox/requests/:id/read", middleware.RequireView(resource), h.SetRequestRead)
			// Clear-all must be registered before /:id so Gin does not capture "requests" as an id segment conflict.
			owner.DELETE("/inbox/requests", middleware.RequireModify(resource), h.SoftDeleteAllRequests)
			owner.DELETE("/inbox/requests/:id", middleware.RequireModify(resource), h.SoftDeleteRequest)
		}
	}
}
