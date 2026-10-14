package seeders

import (
	"context"
	"fmt"

	"be/pkg/postgres"
)

// DatabaseSeeder runs all application seeders (Laravel DatabaseSeeder equivalent):
// core RBAC + users first, then each module's seeders.
type DatabaseSeeder struct {
	seeders []Seeder
}

// NewDatabaseSeeder builds the core seed run plus the given module contributions.
func NewDatabaseSeeder(contribs ...Contribution) *DatabaseSeeder {
	extra := modulePermissions(contribs)
	seeders := []Seeder{
		NewPermissionSeeder(extra...),
		NewRoleSeeder(),
		NewRolePermissionSeeder(extra...),
		NewUserSeeder(),
	}
	for _, c := range contribs {
		seeders = append(seeders, c.Seeders...)
	}
	return &DatabaseSeeder{seeders: seeders}
}

func modulePermissions(contribs []Contribution) []PermissionDef {
	var extra []PermissionDef
	for _, c := range contribs {
		extra = append(extra, c.Permissions...)
	}
	return extra
}

func (s *DatabaseSeeder) Name() string {
	return "DatabaseSeeder"
}

func (s *DatabaseSeeder) Run(ctx context.Context, db *postgres.Postgres) error {
	return Run(ctx, db, s.seeders...)
}

// Registry maps seeder class names to constructors for CLI --class usage: the core
// seeders plus every module seeder (by its Name()).
func Registry(contribs ...Contribution) map[string]func() Seeder {
	extra := modulePermissions(contribs)
	registry := map[string]func() Seeder{
		"DatabaseSeeder":       func() Seeder { return NewDatabaseSeeder(contribs...) },
		"PermissionSeeder":     func() Seeder { return NewPermissionSeeder(extra...) },
		"RoleSeeder":           func() Seeder { return NewRoleSeeder() },
		"RolePermissionSeeder": func() Seeder { return NewRolePermissionSeeder(extra...) },
		"UserSeeder":           func() Seeder { return NewUserSeeder() },
	}
	for _, c := range contribs {
		for _, seeder := range c.Seeders {
			seeder := seeder
			registry[seeder.Name()] = func() Seeder { return seeder }
		}
	}
	return registry
}

// RunByName executes a registered seeder by class name.
func RunByName(ctx context.Context, db *postgres.Postgres, className string, contribs ...Contribution) error {
	constructor, ok := Registry(contribs...)[className]
	if !ok {
		return fmt.Errorf("unknown seeder %q", className)
	}

	seeder := constructor()
	if className == "DatabaseSeeder" {
		return seeder.Run(ctx, db)
	}

	return db.WithTx(ctx, func(ctx context.Context) error {
		return seeder.Run(ctx, db)
	})
}
