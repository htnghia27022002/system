package seeders

import (
	"context"
	"fmt"

	"be/pkg/postgres"
)

// Seeder is a single database seed class (Laravel-style).
type Seeder interface {
	Name() string
	Run(ctx context.Context, db *postgres.Postgres) error
}

// Run executes seeders in order inside one transaction.
func Run(ctx context.Context, db *postgres.Postgres, seeders ...Seeder) error {
	return db.WithTx(ctx, func(ctx context.Context) error {
		for _, seeder := range seeders {
			if err := seeder.Run(ctx, db); err != nil {
				return fmt.Errorf("%s: %w", seeder.Name(), err)
			}
		}
		return nil
	})
}
