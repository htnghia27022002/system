package user

import (
	"context"
	"fmt"
	"strings"

	apperrors "be/internal/common/errors"
	"be/internal/common/hash"
	"be/internal/common/pagination"
	userdto "be/internal/dto/user"
	usermodel "be/internal/models/user"
	"be/internal/repository/interfaces"
)

type Service struct {
	repo interfaces.UserRepository
}

func NewService(repo interfaces.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req userdto.CreateUserRequest) (*usermodel.User, error) {
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: email already exists", apperrors.ErrConflict)
	}

	hashed, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	status := usermodel.StatusActive
	if req.Status != "" {
		status = usermodel.Status(req.Status)
	}

	user := &usermodel.User{
		Email:        req.Email,
		PasswordHash: hashed,
		FullName:     req.Name,
		RoleID:       req.RoleID,
		Status:       status,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*usermodel.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}
	return user, nil
}

func (s *Service) List(ctx context.Context, query userdto.ListUsersQuery) ([]usermodel.User, int64, int, int, error) {
	page, pageSize := pagination.NormalizePage(query.Page, query.PageSize)
	filter := interfaces.UserListFilter{
		Search: strings.TrimSpace(query.Search),
		RoleID: query.RoleID,
		Offset: pagination.Offset(page, pageSize),
		Limit:  pageSize,
	}
	users, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return users, total, page, pageSize, nil
}

func (s *Service) Update(ctx context.Context, id string, req userdto.UpdateUserRequest, sessionUserID string) (*usermodel.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}

	if req.Email != nil && *req.Email != user.Email {
		existing, err := s.repo.GetByEmail(ctx, *req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, fmt.Errorf("%w: email already exists", apperrors.ErrConflict)
		}
		user.Email = *req.Email
	}
	if req.Name != nil {
		user.FullName = *req.Name
	}
	if req.RoleID != nil {
		user.RoleID = *req.RoleID
	}
	if req.Status != nil {
		if id == sessionUserID && *req.Status == string(usermodel.StatusInactive) {
			return nil, fmt.Errorf("%w: cannot deactivate your own account", apperrors.ErrForbidden)
		}
		user.Status = usermodel.Status(*req.Status)
	}
	if req.Password != nil && *req.Password != "" {
		hashed, err := hash.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hashed
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Delete(ctx context.Context, id, sessionUserID string) error {
	if id == sessionUserID {
		return fmt.Errorf("%w: cannot delete your own account", apperrors.ErrForbidden)
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return apperrors.ErrNotFound
	}
	return s.repo.Delete(ctx, id)
}

func ToResponse(user *usermodel.User) userdto.UserResponse {
	return userdto.UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.FullName,
		RoleID: user.RoleID,
		Status: string(user.Status),
	}
}
