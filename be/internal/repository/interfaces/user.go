package interfaces

import (
	"context"

	usermodel "be/internal/models/user"
)

type UserListFilter struct {
	Search string
	RoleID string
	Offset int
	Limit  int
}

type UserRepository interface {
	Create(ctx context.Context, user *usermodel.User) error
	GetByID(ctx context.Context, id string) (*usermodel.User, error)
	GetByEmail(ctx context.Context, email string) (*usermodel.User, error)
	List(ctx context.Context, filter UserListFilter) ([]usermodel.User, int64, error)
	Update(ctx context.Context, user *usermodel.User) error
	Delete(ctx context.Context, id string) error
}
