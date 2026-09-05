package repository

import (
	"context"

	"github.com/Masterminds/squirrel"

	rolemodel "be/internal/models/role"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type RoleRepository struct {
	*repo.Repository[rolemodel.Role]
	links *repo.Repository[rolemodel.RolePermission]
}

var _ interfaces.RoleRepository = (*RoleRepository)(nil)

func NewRoleRepository(db *postgres.Postgres) *RoleRepository {
	return &RoleRepository{
		Repository: repo.New[rolemodel.Role](db, repo.Opts{Table: "roles", PK: "id"}),
		links:      repo.New[rolemodel.RolePermission](db, repo.Opts{Table: "role_permissions", PK: "role_id"}),
	}
}

func (r *RoleRepository) Create(ctx context.Context, role *rolemodel.Role) error {
	return r.Insert(ctx, role)
}

func (r *RoleRepository) GetByID(ctx context.Context, id string) (*rolemodel.Role, error) {
	return r.FindByID(ctx, id)
}

func (r *RoleRepository) GetBySlug(ctx context.Context, slug string) (*rolemodel.Role, error) {
	return r.FindOne(ctx, query.New(1, 1).WhereEqual("slug", slug))
}

func (r *RoleRepository) ListAll(ctx context.Context) ([]rolemodel.Role, error) {
	return r.Find(ctx, query.Unbounded().OrderBy("created_at ASC"))
}

func (r *RoleRepository) List(ctx context.Context, q *query.Query) ([]rolemodel.Role, int64, error) {
	return r.Paginate(ctx, q)
}

func (r *RoleRepository) Update(ctx context.Context, role *rolemodel.Role) error {
	return r.Repository.Update(ctx, role)
}

func (r *RoleRepository) Delete(ctx context.Context, id string) error {
	return r.DB().WithTx(ctx, func(ctx context.Context) error {
		if _, err := r.links.ExecBuilder(ctx, r.DB().Builder.
			Delete(postgres.QuoteIdent("role_permissions")).
			Where(squirrel.Eq{postgres.QuoteIdent("role_id"): id})); err != nil {
			return err
		}
		return r.DeleteByID(ctx, id)
	})
}

func (r *RoleRepository) AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return r.DB().WithTx(ctx, func(ctx context.Context) error {
		if _, err := r.links.ExecBuilder(ctx, r.DB().Builder.
			Delete(postgres.QuoteIdent("role_permissions")).
			Where(squirrel.Eq{postgres.QuoteIdent("role_id"): roleID})); err != nil {
			return err
		}
		for _, permissionID := range permissionIDs {
			link := rolemodel.RolePermission{RoleID: roleID, PermissionID: permissionID}
			if err := r.links.Insert(ctx, &link); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RoleRepository) GetPermissionKeysByRoleID(ctx context.Context, roleID string) ([]string, error) {
	sql, args, err := r.DB().Builder.
		Select("permissions.key").
		From("role_permissions").
		Join("permissions ON permissions.id = role_permissions.permission_id").
		Where(squirrel.Eq{"role_permissions.role_id": roleID}).
		OrderBy("permissions.key ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	return postgres.QueryStrings(ctx, r.DB().Querier(ctx), sql, args...)
}

func (r *RoleRepository) GetPermissionIDsByKeys(ctx context.Context, keys []string) ([]string, error) {
	if len(keys) == 0 {
		return []string{}, nil
	}
	sql, args, err := r.DB().Builder.
		Select("id").
		From(postgres.QuoteIdent("permissions")).
		Where(squirrel.Eq{postgres.QuoteIdent("key"): keys}).
		ToSql()
	if err != nil {
		return nil, err
	}

	return postgres.QueryStrings(ctx, r.DB().Querier(ctx), sql, args...)
}
