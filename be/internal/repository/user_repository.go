package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	usermodel "be/internal/models/user"
	"be/internal/common/query"
	"be/internal/repository/interfaces"
)

type UserRepository struct {
	db *gorm.DB
}

var _ interfaces.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *usermodel.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*usermodel.User, error) {
	var user usermodel.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	var user usermodel.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context, q *query.Query) ([]usermodel.User, int64, error) {
	var users []usermodel.User
	total, err := query.Paginate[usermodel.User](ctx, r.db, q, &users)
	return users, total, err
}

func (r *UserRepository) ListAll(ctx context.Context) ([]usermodel.User, error) {
	var users []usermodel.User
	if err := r.db.WithContext(ctx).
		Where("status = ?", usermodel.StatusActive).
		Order("created_at ASC").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *usermodel.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&usermodel.User{}, "id = ?", id).Error
}
