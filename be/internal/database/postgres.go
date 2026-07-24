package database

import (
	"be/internal/config"
	"be/pkg/postgres"

	"gorm.io/gorm"
)

// Connect opens PostgreSQL via GORM using app config.
func Connect(cfg config.Config) (*gorm.DB, error) {
	return postgres.Connect(postgres.Options{
		Host:    cfg.DBHost,
		Port:    cfg.DBPort,
		User:    cfg.DBUser,
		Pass:    cfg.DBPass,
		Name:    cfg.DBName,
		SSLMode: cfg.DBSSLMode,
	})
}
