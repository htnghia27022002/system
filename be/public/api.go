package public

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"be/internal/app"
	"be/internal/common/cache"
	"be/internal/config"
	"be/internal/database"
	"be/internal/middleware"
	"be/public/routes"
)

func Run(cfg config.Config, db *gorm.DB, redis *goredis.Client) error {
	if err := database.RunMigrations(cfg); err != nil {
		return err
	}
	if err := database.Seed(context.Background(), db); err != nil {
		return err
	}

	if err := cache.Init(cfg.Cache, redis); err != nil {
		return err
	}
	defer func() { _ = cache.Close() }()

	container := app.NewContainer(cfg, db)
	defer container.Close()
	defer container.Close()

	ctx := context.Background()
	if container.SearchClient != nil && container.SearchClient.Enabled() {
		if err := container.SearchService.EnsureIndex(ctx); err != nil {
			log.Printf("search index ensure failed: %v", err)
		}
	}

	r := gin.Default()
	// Permissive CORS for public webhook capture (no credentials) must run before
	// credentialed owner CORS so browser preflight to /api/webhooks/capture/* succeeds.
	r.Use(middleware.WebhookCaptureCORS())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	routes.RegisterAuthRoutes(api, container)
	routes.RegisterAdminRoutes(api, container)
	routes.RegisterMediaRoutes(api, container)
	routes.RegisterWebhookRoutes(api, container)

	return r.Run(":" + cfg.Port)
}
