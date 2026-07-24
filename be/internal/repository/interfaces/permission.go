package interfaces

import (
	"context"

	"be/pkg/query"
	permissionmodel "be/internal/models/permission"
)

type PermissionRepository interface {
	ListAll(ctx context.Context) ([]permissionmodel.Permission, error)
	List(ctx context.Context, q *query.Query) ([]permissionmodel.Permission, int64, error)
	GetByKey(ctx context.Context, key string) (*permissionmodel.Permission, error)
}
