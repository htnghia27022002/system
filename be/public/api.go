package public

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"

	"be/common/cache"
	"be/internal/app"
	"be/internal/config"
	"be/internal/database"
	"be/internal/module"
	"be/pkg/postgres"
	"be/public/routes"
)

const shutdownTimeout = 15 * time.Second

func Run(cfg config.Config, db *postgres.Postgres, redis *goredis.Client) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container, err := app.NewContainer(cfg, db)
	if err != nil {
		return err
	}
	defer container.Close()
	log.Printf("modules enabled: %v", container.Modules.Names())

	if cfg.AutoMigrate {
		if err := database.RunMigrations(cfg, container.Modules.MigrationSets()...); err != nil {
			return err
		}
		if err := database.Seed(ctx, db, container.Modules.SeedContributions()...); err != nil {
			return err
		}
	}

	if err := cache.Init(cfg.Cache, redis); err != nil {
		return err
	}
	defer func() { _ = cache.Close() }()

	r := gin.Default()
	// Module middleware (e.g. permissive CORS for public webhook capture) runs before
	// the credentialed core CORS so its browser preflight succeeds.
	for _, mw := range container.Modules.GlobalMiddleware() {
		r.Use(mw)
	}
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
	admin := routes.RegisterAdminRoutes(api, container)
	routes.RegisterMediaRoutes(api, container)
	container.Modules.RegisterRoutes(module.Routes{API: api, Admin: admin})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	serveErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
		close(serveErr)
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	log.Printf("shutting down HTTP server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
