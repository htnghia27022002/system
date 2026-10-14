package public

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
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
	"be/internal/module"
	"be/pkg/postgres"
	"be/public/handlers"
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

	logger := slog.Default()
	r := gin.New()
	// Only these peers may set X-Forwarded-For / X-Real-IP; everyone else is
	// identified by the socket address (rate limits and session IPs rely on it).
	if err := r.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		return fmt.Errorf("TRUSTED_PROXIES: %w", err)
	}
	r.Use(
		middleware.RequestID(),
		middleware.AuditContext(),
		middleware.Recovery(logger),
		middleware.AccessLog(logger),
		middleware.SecurityHeaders(),
	)
	// Module middleware (e.g. permissive CORS for public webhook capture) runs before
	// the credentialed core CORS so its browser preflight succeeds.
	for _, mw := range container.Modules.GlobalMiddleware() {
		r.Use(mw)
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization",
			"X-Requested-With", "X-Session-Id", middleware.RequestIDHeader,
		},
		ExposeHeaders:    []string{"Content-Length", middleware.RequestIDHeader},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	health := handlers.NewHealthHandler(readinessChecks(container, db, redis)...)

	api := r.Group("/api")
	routes.RegisterAppRoutes(api, container, health)
	routes.RegisterAuthRoutes(api, container)
	admin := routes.RegisterAdminRoutes(api, container)
	routes.RegisterMediaRoutes(api, container)
	container.Modules.RegisterRoutes(module.Routes{API: api, Admin: admin})

	go container.Janitor.Run(ctx)

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

// readinessChecks lists the dependencies /api/ready verifies: Postgres always,
// Redis and NATS only when configured.
func readinessChecks(c *app.Container, db *postgres.Postgres, redis *goredis.Client) []handlers.ReadinessCheck {
	checks := []handlers.ReadinessCheck{
		{Name: "database", Check: func(ctx context.Context) error { return db.Pool.Ping(ctx) }},
	}
	if redis != nil {
		checks = append(checks, handlers.ReadinessCheck{
			Name:  "redis",
			Check: func(ctx context.Context) error { return redis.Ping(ctx).Err() },
		})
	}
	if c.Queue.Connection.Enabled {
		checks = append(checks, handlers.ReadinessCheck{
			Name: "queue",
			Check: func(context.Context) error {
				if c.QueueClient == nil || !c.QueueClient.Enabled() {
					return errors.New("nats not connected")
				}
				return nil
			},
		})
	}
	// Remote object stores (S3) expose Ping; the local disk needs no probe.
	if pinger, ok := c.Storage.(interface{ Ping(context.Context) error }); ok {
		checks = append(checks, handlers.ReadinessCheck{Name: "storage", Check: pinger.Ping})
	}
	return checks
}
