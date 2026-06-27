package seeders

import (
	"context"

	"be/internal/common/rbac"
	rolemodel "be/internal/models/role"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RolePermissionSeeder assigns permissions to roles.
type RolePermissionSeeder struct{}

func NewRolePermissionSeeder() *RolePermissionSeeder {
	return &RolePermissionSeeder{}
}

func (s *RolePermissionSeeder) Name() string {
	return "RolePermissionSeeder"
}

func (s *RolePermissionSeeder) Run(ctx context.Context, db *gorm.DB) error {
	for _, item := range DefaultPermissions() {
		row := rolemodel.RolePermission{
			RoleID:       RoleAdminID,
			PermissionID: item.ID,
		}
		if err := db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error; err != nil {
			return err
		}
	}

	memberDashboard := rolemodel.RolePermission{
		RoleID:       RoleUserID,
		PermissionID: PermissionIDByKey(rbac.Key("dashboard", rbac.ActionView)),
	}
	if err := db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&memberDashboard).Error; err != nil {
		return err
	}

	return nil
}
