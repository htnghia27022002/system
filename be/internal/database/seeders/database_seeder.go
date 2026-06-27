package seeders

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// DatabaseSeeder runs all application seeders (Laravel DatabaseSeeder equivalent).
type DatabaseSeeder struct {
	seeders []Seeder
}

func NewDatabaseSeeder() *DatabaseSeeder {
	return &DatabaseSeeder{
		seeders: []Seeder{
			NewPermissionSeeder(),
			NewRoleSeeder(),
			NewRolePermissionSeeder(),
			NewUserSeeder(),
		},
	}
}

func (s *DatabaseSeeder) Name() string {
	return "DatabaseSeeder"
}

func (s *DatabaseSeeder) Run(ctx context.Context, db *gorm.DB) error {
	return Run(ctx, db, s.seeders...)
}

// Registry maps seeder class names to constructors for CLI --class usage.
var Registry = map[string]func() Seeder{
	"DatabaseSeeder":        func() Seeder { return NewDatabaseSeeder() },
	"PermissionSeeder":      func() Seeder { return NewPermissionSeeder() },
	"RoleSeeder":            func() Seeder { return NewRoleSeeder() },
	"RolePermissionSeeder": func() Seeder { return NewRolePermissionSeeder() },
	"UserSeeder":            func() Seeder { return NewUserSeeder() },
}

// RunByName executes a registered seeder by class name.
func RunByName(ctx context.Context, db *gorm.DB, className string) error {
	constructor, ok := Registry[className]
	if !ok {
		return fmt.Errorf("unknown seeder %q", className)
	}

	seeder := constructor()
	if className == "DatabaseSeeder" {
		return seeder.Run(ctx, db)
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return seeder.Run(ctx, tx)
	})
}
