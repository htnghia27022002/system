package repository

import (
	"context"

	permissionmodel "be/internal/models/permission"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type PermissionRepository struct {
	*repo.Repository[permissionmodel.Permission]
}

var _ interfaces.PermissionRepository = (*PermissionRepository)(nil)

func NewPermissionRepository(db *postgres.Postgres) *PermissionRepository {
	return &PermissionRepository{
		Repository: repo.New[permissionmodel.Permission](db, repo.Opts{
			Table: "permissions",
			PK:    "id",
		}),
	}
}

func (r *PermissionRepository) ListAll(ctx context.Context) ([]permissionmodel.Permission, error) {
	return r.Find(ctx, query.Unbounded().OrderBy(`"group" ASC, key ASC`))
}

func (r *PermissionRepository) List(ctx context.Context, q *query.Query) ([]permissionmodel.Permission, int64, error) {
	return r.Paginate(ctx, q)
}

func (r *PermissionRepository) GetByKey(ctx context.Context, key string) (*permissionmodel.Permission, error) {
	return r.FindOne(ctx, query.New(1, 1).WhereEqual("key", key))
}
