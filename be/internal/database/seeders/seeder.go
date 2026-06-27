package seeders

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Seeder is a single database seed class (Laravel-style).
type Seeder interface {
	Name() string
	Run(ctx context.Context, db *gorm.DB) error
}

// Run executes seeders in order inside one transaction.
func Run(ctx context.Context, db *gorm.DB, seeders ...Seeder) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, seeder := range seeders {
			if err := seeder.Run(ctx, tx); err != nil {
				return fmt.Errorf("%s: %w", seeder.Name(), err)
			}
		}
		return nil
	})
}
