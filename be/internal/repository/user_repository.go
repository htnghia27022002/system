package repository

import (
	"context"

	usermodel "be/internal/models/user"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type UserRepository struct {
	*repo.Repository[usermodel.User]
}

var _ interfaces.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *postgres.Postgres) *UserRepository {
	return &UserRepository{
		Repository: repo.New[usermodel.User](db, repo.Opts{
			Table:            "users",
			PK:               "id",
			SoftDeleteColumn: "deleted_at",
		}),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *usermodel.User) error {
	return r.Insert(ctx, user)
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*usermodel.User, error) {
	return r.FindByID(ctx, id)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	return r.FindOne(ctx, query.New(1, 1).WhereEqual("email", email))
}

func (r *UserRepository) List(ctx context.Context, q *query.Query) ([]usermodel.User, int64, error) {
	return r.Paginate(ctx, q)
}

func (r *UserRepository) ListAll(ctx context.Context) ([]usermodel.User, error) {
	return r.Find(ctx, query.Unbounded().
		WhereEqual("status", string(usermodel.StatusActive)).
		OrderBy("created_at ASC"))
}

func (r *UserRepository) Update(ctx context.Context, user *usermodel.User) error {
	return r.Repository.Update(ctx, user)
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.DeleteByID(ctx, id)
}

func (r *UserRepository) CountSuperAdmins(ctx context.Context) (int64, error) {
	_, total, err := r.Paginate(ctx, query.New(1, 1).WhereRaw("is_super_admin = TRUE"))
	return total, err
}
