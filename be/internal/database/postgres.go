package database

import (
	"be/internal/config"
	"be/pkg/postgres"

	"gorm.io/gorm"
)

// Connect opens PostgreSQL via GORM using app config.
// When cfg.DBURL is set (from DB_URL), it is used; otherwise discrete DB_* fields.
func Connect(cfg config.Config) (*gorm.DB, error) {
	if cfg.DBURL != "" {
		return postgres.ConnectDSN(cfg.DBURL)
	}
	return postgres.Connect(postgres.Options{
		Host:    cfg.DBHost,
		Port:    cfg.DBPort,
		User:    cfg.DBUser,
		Pass:    cfg.DBPass,
		Name:    cfg.DBName,
		SSLMode: cfg.DBSSLMode,
	})
}
