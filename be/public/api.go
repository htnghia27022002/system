package public

import (
	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"be/internal/app"
	"be/internal/common/cache"
	"be/internal/config"
	"be/internal/database"
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
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	routes.RegisterAuthRoutes(api, container)
	routes.RegisterAdminRoutes(api, container)

	return r.Run(":" + cfg.Port)
}
