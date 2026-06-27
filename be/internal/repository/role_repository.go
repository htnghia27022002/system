package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	permissionmodel "be/internal/models/permission"
	rolemodel "be/internal/models/role"
	"be/internal/repository/interfaces"
)

type RoleRepository struct {
	db *gorm.DB
}

var _ interfaces.RoleRepository = (*RoleRepository)(nil)

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *rolemodel.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *RoleRepository) GetByID(ctx context.Context, id string) (*rolemodel.Role, error) {
	var role rolemodel.Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetBySlug(ctx context.Context, slug string) (*rolemodel.Role, error) {
	var role rolemodel.Role
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) List(ctx context.Context) ([]rolemodel.Role, error) {
	var roles []rolemodel.Role
	if err := r.db.WithContext(ctx).Order("created_at ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *rolemodel.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *RoleRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&rolemodel.RolePermission{}).Error; err != nil {
			return err
		}
		return tx.Delete(&rolemodel.Role{}, "id = ?", id).Error
	})
}

func (r *RoleRepository) AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&rolemodel.RolePermission{}).Error; err != nil {
			return err
		}
		for _, permissionID := range permissionIDs {
			link := rolemodel.RolePermission{
				RoleID:       roleID,
				PermissionID: permissionID,
			}
			if err := tx.Create(&link).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RoleRepository) GetPermissionKeysByRoleID(ctx context.Context, roleID string) ([]string, error) {
	var keys []string
	err := r.db.WithContext(ctx).
		Table("role_permissions").
		Select("permissions.key").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Order("permissions.key ASC").
		Scan(&keys).Error
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (r *RoleRepository) GetPermissionIDsByKeys(ctx context.Context, keys []string) ([]string, error) {
	if len(keys) == 0 {
		return []string{}, nil
	}
	var ids []string
	if err := r.db.WithContext(ctx).
		Model(&permissionmodel.Permission{}).
		Where("key IN ?", keys).
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}
