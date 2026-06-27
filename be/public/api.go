package public

import (
	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"be/internal/app"
	"be/internal/config"
	"be/internal/database"
	authmodel "be/internal/models/auth"
	permissionmodel "be/internal/models/permission"
	rolemodel "be/internal/models/role"
	usermodel "be/internal/models/user"
	"be/public/routes"
)

func Run(cfg config.Config, db *gorm.DB) error {
	if err := migrate(db); err != nil {
		return err
	}
	if err := database.SeedRBAC(context.Background(), db); err != nil {
		return err
	}

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

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&usermodel.User{},
		&rolemodel.Role{},
		&rolemodel.RolePermission{},
		&permissionmodel.Permission{},
		&authmodel.RefreshToken{},
		&authmodel.OAuthAccount{},
	)
}
