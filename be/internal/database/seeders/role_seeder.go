package seeders

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	rolemodel "be/internal/models/role"
)

// RoleSeeder seeds default roles.
type RoleSeeder struct{}

func NewRoleSeeder() *RoleSeeder {
	return &RoleSeeder{}
}

func (s *RoleSeeder) Name() string {
	return "RoleSeeder"
}

func (s *RoleSeeder) Run(ctx context.Context, db *gorm.DB) error {
	roles := []rolemodel.Role{
		{
			ID:          RoleAdminID,
			Name:        "Administrator",
			Slug:        "admin",
			Description: "Full system access",
		},
		{
			ID:          RoleUserID,
			Name:        "Member",
			Slug:        "user",
			Description: "Standard member access",
		},
	}

	for _, role := range roles {
		if err := db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "slug"}},
			DoNothing: true,
		}).Create(&role).Error; err != nil {
			return err
		}
	}

	return nil
}
