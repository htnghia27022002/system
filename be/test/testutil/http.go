package testutil

import (
	"context"
	"testing"

	"be/internal/app"
	"be/internal/config"
	"be/internal/database"
	"be/internal/middleware"
	"be/pkg/postgres"
	"be/public/handlers"
	"be/public/routes"

	"github.com/gin-gonic/gin"
)

// ConnectPostgres connects using config.Load() env vars. Skips the test when DB is unavailable.
func ConnectPostgres(t *testing.T) *postgres.Postgres {
	t.Helper()
	SkipIfShort(t)

	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		t.Skipf("postgres not available: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// MigrateTestSchema applies the core SQL migrations for integration and e2e tests.
func MigrateTestSchema(t *testing.T, _ *postgres.Postgres) {
	t.Helper()
	cfg := config.Load()
	if err := database.RunMigrations(cfg); err != nil {
		t.Fatalf("migrate test schema: %v", err)
	}
}

// TruncateAuthTables clears auth-related tables between tests.
func TruncateAuthTables(t *testing.T, db *postgres.Postgres) {
	t.Helper()
	_, err := db.Pool.Exec(context.Background(),
		`TRUNCATE users, roles, role_permissions, permissions, refresh_tokens, oauth_accounts, user_tokens, audit_logs, app_settings RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

// NewTestContainer wires the DI container against a test database. Options adjust
// the loaded config (e.g. turn the refresh-token cookie off).
func NewTestContainer(t *testing.T, db *postgres.Postgres, opts ...func(*config.Config)) *app.Container {
	t.Helper()
	return NewTestContainerWith(t, db, nil, opts...)
}

// NewTestContainerWith is NewTestContainer plus app options (e.g. app.WithMailer).
func NewTestContainerWith(t *testing.T, db *postgres.Postgres, appOpts []app.Option, opts ...func(*config.Config)) *app.Container {
	t.Helper()
	cfg := config.Load()
	cfg.JWTSecret = "integration-test-secret"
	for _, opt := range opts {
		opt(&cfg)
	}
	container, err := app.NewContainer(cfg, db, appOpts...)
	if err != nil {
		t.Fatalf("new container: %v", err)
	}
	t.Cleanup(func() { container.Close() })
	return container
}

// WithRefreshInBody returns the refresh token in JSON instead of an httpOnly cookie.
func WithRefreshInBody(cfg *config.Config) { cfg.Auth.RefreshCookie = false }

// NewTestRouter returns a Gin engine with /api auth routes for HTTP-level tests.
func NewTestRouter(t *testing.T, container *app.Container) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	routes.RegisterAuthRoutes(r.Group("/api"), container)
	return r
}

// NewTestAPIRouter mounts auth and core admin routes.
func NewTestAPIRouter(t *testing.T, container *app.Container) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.AuditContext())
	api := r.Group("/api")
	routes.RegisterAppRoutes(api, container, handlers.NewHealthHandler())
	routes.RegisterAuthRoutes(api, container)
	routes.RegisterAdminRoutes(api, container)
	return r
}
