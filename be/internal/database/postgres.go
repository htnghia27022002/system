package database

import (
	"context"

	"be/internal/config"
	"be/pkg/postgres"
)

// Connect opens PostgreSQL via pgx using app config.
// When cfg.DBURL is set (from DB_URL), it is used; otherwise discrete DB_* fields.
func Connect(cfg config.Config) (*postgres.Postgres, error) {
	ctx := context.Background()
	if cfg.DBURL != "" {
		return postgres.ConnectDSN(ctx, cfg.DBURL)
	}
	return postgres.Connect(ctx, postgres.Options{
		Host:    cfg.DBHost,
		Port:    cfg.DBPort,
		User:    cfg.DBUser,
		Pass:    cfg.DBPass,
		Name:    cfg.DBName,
		SSLMode: cfg.DBSSLMode,
	})
}
