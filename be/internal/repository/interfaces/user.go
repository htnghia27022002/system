package interfaces

import (
	"context"

	"be/pkg/query"
	usermodel "be/internal/models/user"
)

type UserRepository interface {
	Create(ctx context.Context, user *usermodel.User) error
	GetByID(ctx context.Context, id string) (*usermodel.User, error)
	GetByEmail(ctx context.Context, email string) (*usermodel.User, error)
	List(ctx context.Context, q *query.Query) ([]usermodel.User, int64, error)
	ListAll(ctx context.Context) ([]usermodel.User, error)
	Update(ctx context.Context, user *usermodel.User) error
	Delete(ctx context.Context, id string) error
}
