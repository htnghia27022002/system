package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Options holds connection parameters for PostgreSQL via GORM.
type Options struct {
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	SSLMode string
}

// Connect opens a GORM PostgreSQL connection from Options.
func Connect(opts Options) (*gorm.DB, error) {
	sslMode := opts.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		opts.Host, opts.User, opts.Pass, opts.Name, opts.Port, sslMode,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// ConnectDSN opens a GORM PostgreSQL connection from a DSN string.
func ConnectDSN(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
