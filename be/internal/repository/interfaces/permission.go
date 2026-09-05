package interfaces

import (
	"context"

	permissionmodel "be/internal/models/permission"
	"be/pkg/query"
)

type PermissionRepository interface {
	ListAll(ctx context.Context) ([]permissionmodel.Permission, error)
	List(ctx context.Context, q *query.Query) ([]permissionmodel.Permission, int64, error)
	GetByKey(ctx context.Context, key string) (*permissionmodel.Permission, error)
}
