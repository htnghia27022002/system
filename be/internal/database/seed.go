package database

import (
	"context"

	"be/internal/database/seeders"
	"be/pkg/postgres"
)

// Seed runs the DatabaseSeeder: core seeders plus each module contribution.
// It refuses to run when two permissions share a key or an ID.
func Seed(ctx context.Context, db *postgres.Postgres, contribs ...seeders.Contribution) error {
	if _, err := seeders.MergePermissions(contribs...); err != nil {
		return err
	}
	return seeders.NewDatabaseSeeder(contribs...).Run(ctx, db)
}

// SeedRBAC seeds core RBAC and users only (used by tests).
func SeedRBAC(ctx context.Context, db *postgres.Postgres) error {
	return Seed(ctx, db)
}

// SeedClass runs a single registered seeder by class name.
func SeedClass(ctx context.Context, db *postgres.Postgres, className string, contribs ...seeders.Contribution) error {
	if _, err := seeders.MergePermissions(contribs...); err != nil {
		return err
	}
	return seeders.RunByName(ctx, db, className, contribs...)
}
