package testutil

import (
	"context"
	"testing"

	"be/internal/app"
	"be/internal/config"
	"be/internal/database"
	"be/pkg/postgres"
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

// MigrateTestSchema applies SQL migrations for integration and e2e tests.
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
		`TRUNCATE users, roles, role_permissions, permissions, refresh_tokens, oauth_accounts RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

// NewTestContainer wires the DI container against a test database.
func NewTestContainer(t *testing.T, db *postgres.Postgres) *app.Container {
	t.Helper()
	cfg := config.Load()
	cfg.JWTSecret = "integration-test-secret"
	container := app.NewContainer(cfg, db)
	t.Cleanup(func() { container.Close() })
	return container
}

// NewTestRouter returns a Gin engine with /api auth routes for HTTP-level tests.
func NewTestRouter(t *testing.T, container *app.Container) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	routes.RegisterAuthRoutes(r.Group("/api"), container)
	return r
}
