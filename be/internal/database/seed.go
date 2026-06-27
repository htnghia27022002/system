package database

import (
	"context"

	"gorm.io/gorm"

	"be/internal/database/seeders"
)

// Seed runs the default DatabaseSeeder (Laravel-style).
func Seed(ctx context.Context, db *gorm.DB) error {
	return seeders.NewDatabaseSeeder().Run(ctx, db)
}

// SeedRBAC is kept for backward compatibility with existing callers.
func SeedRBAC(ctx context.Context, db *gorm.DB) error {
	return Seed(ctx, db)
}

// SeedClass runs a single registered seeder by class name.
func SeedClass(ctx context.Context, db *gorm.DB, className string) error {
	return seeders.RunByName(ctx, db, className)
}
