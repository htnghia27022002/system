package role

import (
	"context"
	"fmt"

	apperrors "be/internal/common/errors"
	roledto "be/internal/dto/role"
	rolemodel "be/internal/models/role"
	"be/internal/repository/interfaces"
)

type Service struct {
	repo interfaces.RoleRepository
}

func NewService(repo interfaces.RoleRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req roledto.CreateRoleRequest) (*roledto.RoleResponse, error) {
	permissionIDs, err := s.repo.GetPermissionIDsByKeys(ctx, req.PermissionKeys)
	if err != nil {
		return nil, err
	}
	if len(permissionIDs) != len(req.PermissionKeys) {
		return nil, fmt.Errorf("%w: one or more permission keys are invalid", apperrors.ErrBadRequest)
	}

	role := &rolemodel.Role{
		Name: req.Name,
		Slug: req.Slug,
	}
	if err := s.repo.Create(ctx, role); err != nil {
		return nil, err
	}
	if err := s.repo.AssignPermissions(ctx, role.ID, permissionIDs); err != nil {
		return nil, err
	}
	return s.toResponse(ctx, role)
}

func (s *Service) GetByID(ctx context.Context, id string) (*roledto.RoleResponse, error) {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, apperrors.ErrNotFound
	}
	return s.toResponse(ctx, role)
}

func (s *Service) List(ctx context.Context) ([]roledto.RoleResponse, error) {
	roles, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]roledto.RoleResponse, 0, len(roles))
	for i := range roles {
		item, err := s.toResponse(ctx, &roles[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *item)
	}
	return out, nil
}

func (s *Service) Update(ctx context.Context, id string, req roledto.UpdateRoleRequest) (*roledto.RoleResponse, error) {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, apperrors.ErrNotFound
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Slug != nil {
		role.Slug = *req.Slug
	}
	if err := s.repo.Update(ctx, role); err != nil {
		return nil, err
	}

	if req.PermissionKeys != nil {
		permissionIDs, err := s.repo.GetPermissionIDsByKeys(ctx, req.PermissionKeys)
		if err != nil {
			return nil, err
		}
		if len(permissionIDs) != len(req.PermissionKeys) {
			return nil, fmt.Errorf("%w: one or more permission keys are invalid", apperrors.ErrBadRequest)
		}
		if err := s.repo.AssignPermissions(ctx, role.ID, permissionIDs); err != nil {
			return nil, err
		}
	}

	return s.toResponse(ctx, role)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return apperrors.ErrNotFound
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) toResponse(ctx context.Context, role *rolemodel.Role) (*roledto.RoleResponse, error) {
	keys, err := s.repo.GetPermissionKeysByRoleID(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	if keys == nil {
		keys = []string{}
	}
	return &roledto.RoleResponse{
		ID:             role.ID,
		Name:           role.Name,
		Slug:           role.Slug,
		PermissionKeys: keys,
	}, nil
}
