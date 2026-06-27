package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	permissionmodel "be/internal/models/permission"
	"be/internal/repository/interfaces"
)

type PermissionRepository struct {
	db *gorm.DB
}

var _ interfaces.PermissionRepository = (*PermissionRepository)(nil)

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) ListAll(ctx context.Context) ([]permissionmodel.Permission, error) {
	var permissions []permissionmodel.Permission
	if err := r.db.WithContext(ctx).Order(`"group" ASC, key ASC`).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *PermissionRepository) GetByKey(ctx context.Context, key string) (*permissionmodel.Permission, error) {
	var permission permissionmodel.Permission
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}
