package permission

import (
	"context"

	"be/internal/common/query"
	permissiondto "be/internal/dto/permission"
	"be/internal/repository/interfaces"
)

type Service struct {
	repo interfaces.PermissionRepository
}

func NewService(repo interfaces.PermissionRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListAll(ctx context.Context) ([]permissiondto.PermissionResponse, error) {
	permissions, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]permissiondto.PermissionResponse, 0, len(permissions))
	for _, item := range permissions {
		out = append(out, permissiondto.PermissionResponse{
			Key:         item.Key,
			Name:        item.Name,
			Group:       item.Group,
			Description: item.Description,
		})
	}
	return out, nil
}

func (s *Service) List(ctx context.Context, form permissiondto.ListPermissionsQuery) ([]permissiondto.PermissionResponse, int64, int, int, error) {
	q := query.New(form.Page, form.PageSize).
		OrderBy(`"group" ASC, key ASC`).
		WhereEqual(`"group"`, form.Group).
		WhereLikeAny([]string{"key", "name", "description"}, form.Search)

	permissions, total, err := s.repo.List(ctx, q)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	out := make([]permissiondto.PermissionResponse, 0, len(permissions))
	for _, item := range permissions {
		out = append(out, permissiondto.PermissionResponse{
			Key:         item.Key,
			Name:        item.Name,
			Group:       item.Group,
			Description: item.Description,
		})
	}
	return out, total, q.Page, q.PageSize, nil
}
