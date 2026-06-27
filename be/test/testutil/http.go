package testutil

import (
	"testing"

	"be/internal/app"
	"be/internal/config"
	"be/internal/database"
	"be/public/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ConnectPostgres connects using config.Load() env vars. Skips the test when DB is unavailable.
func ConnectPostgres(t *testing.T) *gorm.DB {
	t.Helper()
	SkipIfShort(t)

	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		t.Skipf("postgres not available: %v", err)
	}
	return db
}

// MigrateTestSchema applies SQL migrations for integration and e2e tests.
func MigrateTestSchema(t *testing.T, db *gorm.DB) {
	t.Helper()
	cfg := config.Load()
	if err := database.RunMigrations(cfg); err != nil {
		t.Fatalf("migrate test schema: %v", err)
	}
}

// NewTestContainer wires the DI container against a test database.
func NewTestContainer(t *testing.T, db *gorm.DB) *app.Container {
	t.Helper()
	cfg := config.Load()
	cfg.JWTSecret = "integration-test-secret"
	return app.NewContainer(cfg, db)
}

// NewTestRouter returns a Gin engine with /api auth routes for HTTP-level tests.
func NewTestRouter(t *testing.T, container *app.Container) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	routes.RegisterAuthRoutes(r.Group("/api"), container)
	return r
}
