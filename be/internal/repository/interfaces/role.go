package interfaces

import (
	"context"

	rolemodel "be/internal/models/role"
)

type RoleRepository interface {
	Create(ctx context.Context, role *rolemodel.Role) error
	GetByID(ctx context.Context, id string) (*rolemodel.Role, error)
	GetBySlug(ctx context.Context, slug string) (*rolemodel.Role, error)
	List(ctx context.Context) ([]rolemodel.Role, error)
	Update(ctx context.Context, role *rolemodel.Role) error
	Delete(ctx context.Context, id string) error
	AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error
	GetPermissionKeysByRoleID(ctx context.Context, roleID string) ([]string, error)
	GetPermissionIDsByKeys(ctx context.Context, keys []string) ([]string, error)
}
