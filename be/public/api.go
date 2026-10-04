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
	"be/internal/middleware"
	"be/pkg/postgres"
	"be/public/routes"
)

const shutdownTimeout = 15 * time.Second

func Run(cfg config.Config, db *postgres.Postgres, redis *goredis.Client) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.AutoMigrate {
		if err := database.RunMigrations(cfg); err != nil {
			return err
		}
		if err := database.Seed(ctx, db); err != nil {
			return err
		}
	}

	if err := cache.Init(cfg.Cache, redis); err != nil {
		return err
	}
	defer func() { _ = cache.Close() }()

	container := app.NewContainer(cfg, db)
	defer container.Close()

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
	routes.RegisterAddressRoutes(api, container)
	routes.RegisterAdminRoutes(api, container)
	routes.RegisterMediaRoutes(api, container)
	routes.RegisterWebhookRoutes(api, container)

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
