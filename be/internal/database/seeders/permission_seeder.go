package seeders

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	permissionmodel "be/internal/models/permission"
)

// PermissionSeeder seeds the permissions catalog.
type PermissionSeeder struct{}

func NewPermissionSeeder() *PermissionSeeder {
	return &PermissionSeeder{}
}

func (s *PermissionSeeder) Name() string {
	return "PermissionSeeder"
}

func (s *PermissionSeeder) Run(ctx context.Context, db *gorm.DB) error {
	for _, item := range DefaultPermissions() {
		row := permissionmodel.Permission{
			ID:          item.ID,
			Key:         item.Key,
			Name:        item.Name,
			Group:       item.Group,
			Description: item.Description,
		}
		if err := db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoNothing: true,
		}).Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}
