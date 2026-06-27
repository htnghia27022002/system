package interfaces

import (
	"context"

	permissionmodel "be/internal/models/permission"
)

type PermissionRepository interface {
	ListAll(ctx context.Context) ([]permissionmodel.Permission, error)
	GetByKey(ctx context.Context, key string) (*permissionmodel.Permission, error)
}
