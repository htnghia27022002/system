package permission

import (
	"context"

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
