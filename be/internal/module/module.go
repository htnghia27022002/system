// Package module defines how optional feature modules plug into the core app.
//
// The core (auth, users, roles, permissions, profile, media) always runs. A module
// is a vertical slice under internal/modules/<name> that owns its handlers,
// services, repositories, models, DTOs, migrations, seeders and queue consumers,
// and is wired only through this package. Removing a module from a cloned project
// means deleting its folder and its line in internal/modules/registry.go.
package module

import (
	"context"
	"io/fs"

	"github.com/gin-gonic/gin"

	jwtmanager "be/common/jwt"
	"be/internal/config"
	"be/internal/database/seeders"
	"be/internal/handlers/publisher"
	"be/internal/queue"
	"be/pkg/postgres"
)

// Kernel is the shared core infrastructure a module may depend on.
type Kernel struct {
	Config    config.Config
	DB        *postgres.Postgres
	JWT       *jwtmanager.Manager
	Publisher *publisher.Publisher
	// Auth authenticates the request (JWT) and loads the caller's permissions;
	// use it on routes outside Routes.Admin. Routes.Admin already applies it.
	Auth gin.HandlerFunc
}

// Routes are the router groups a module mounts on.
type Routes struct {
	// API is /api (no auth).
	API *gin.RouterGroup
	// Admin is /api/admin with Kernel.Auth applied; add RequireView/RequireModify per route.
	Admin *gin.RouterGroup
}

// Module is the required surface of every module. Its constructor
// (func New(*Kernel) Module) must only wire objects — no I/O — because CLIs such
// as cmd/migrate build modules from a Kernel that has only Config set.
type Module interface {
	Name() string
	RegisterRoutes(r Routes)
}

// Optional capabilities — implement only what the module needs.

// PermissionProvider contributes {resource}:view / {resource}:modify keys to the
// RBAC catalog. The admin role receives them automatically.
type PermissionProvider interface {
	Permissions() []seeders.PermissionDef
}

// SeederProvider contributes seeders that run after the core seeders.
type SeederProvider interface {
	Seeders() []seeders.Seeder
}

// MigrationProvider contributes SQL migrations, versioned in their own table
// (schema_migrations_<name>) so modules never share a version sequence.
type MigrationProvider interface {
	Migrations() fs.FS
}

// MiddlewareProvider contributes engine-level middleware that must run before
// the core CORS middleware (e.g. permissive CORS for a public capture endpoint).
type MiddlewareProvider interface {
	GlobalMiddleware() []gin.HandlerFunc
}

// StreamProvider declares JetStream streams/consumers. cmd/queue ensures and
// consumes them when NATS is enabled.
type StreamProvider interface {
	Streams() []queue.Stream
}

// Poller processes queued work when NATS is disabled; cmd/queue calls it on a ticker.
type Poller interface {
	Poll(ctx context.Context) (processed int, err error)
}

// Closer releases module resources on shutdown.
type Closer interface {
	Close()
}
