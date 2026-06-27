package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	usermodel "be/internal/models/user"
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

func (r *UserRepository) List(ctx context.Context, filter interfaces.UserListFilter) ([]usermodel.User, int64, error) {
	query := r.db.WithContext(ctx).Model(&usermodel.User{})
	if filter.Search != "" {
		term := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(email) LIKE ? OR LOWER(full_name) LIKE ?", term, term)
	}
	if filter.RoleID != "" {
		query = query.Where("role_id = ?", filter.RoleID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []usermodel.User
	if err := query.Order("created_at DESC").Offset(filter.Offset).Limit(filter.Limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) Update(ctx context.Context, user *usermodel.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&usermodel.User{}, "id = ?", id).Error
}
