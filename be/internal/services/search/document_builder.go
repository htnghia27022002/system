package searchsvc

import (
	"context"
	"fmt"
	"strings"

	"be/pkg/query"
	"be/internal/common/rbac"
	searchpkg "be/internal/search"
	permissionmodel "be/internal/models/permission"
	rolemodel "be/internal/models/role"
	usermodel "be/internal/models/user"
	"be/internal/repository/interfaces"
)

type DocumentBuilder struct {
	users       interfaces.UserRepository
	roles       interfaces.RoleRepository
	permissions interfaces.PermissionRepository
}

func NewDocumentBuilder(
	users interfaces.UserRepository,
	roles interfaces.RoleRepository,
	permissions interfaces.PermissionRepository,
) *DocumentBuilder {
	return &DocumentBuilder{users: users, roles: roles, permissions: permissions}
}

func (b *DocumentBuilder) Build(ctx context.Context, entityType, entityID string) (*searchpkg.Document, error) {
	switch entityType {
	case searchpkg.EntityUser:
		return b.buildUser(ctx, entityID)
	case searchpkg.EntityRole:
		return b.buildRole(ctx, entityID)
	case searchpkg.EntityPermission:
		return b.buildPermission(ctx, entityID)
	default:
		return nil, fmt.Errorf("unsupported entity type: %s", entityType)
	}
}

func (b *DocumentBuilder) buildUser(ctx context.Context, id string) (*searchpkg.Document, error) {
	user, err := b.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Status != usermodel.StatusActive {
		return nil, nil
	}

	searchable := strings.ToLower(strings.Join([]string{
		user.FullName,
		user.Email,
		string(user.Status),
	}, " "))

	return &searchpkg.Document{
		EntityType:     searchpkg.EntityUser,
		EntityID:       user.ID,
		Title:          user.FullName,
		SearchableText: searchable,
		PermissionKeys: []string{rbac.Key("users", rbac.ActionView)},
		Metadata: map[string]string{
			"email":  user.Email,
			"status": string(user.Status),
		},
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (b *DocumentBuilder) buildRole(ctx context.Context, id string) (*searchpkg.Document, error) {
	role, err := b.roles.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	searchable := strings.ToLower(strings.Join([]string{
		role.Name,
		role.Slug,
		role.Description,
	}, " "))

	return &searchpkg.Document{
		EntityType:     searchpkg.EntityRole,
		EntityID:       role.ID,
		Title:          role.Name,
		SearchableText: searchable,
		PermissionKeys: []string{rbac.Key("roles", rbac.ActionView)},
		Metadata: map[string]string{
			"slug": role.Slug,
		},
		UpdatedAt: role.UpdatedAt,
	}, nil
}

func (b *DocumentBuilder) buildPermission(ctx context.Context, id string) (*searchpkg.Document, error) {
	permissions, err := b.permissions.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range permissions {
		if item.ID != id {
			continue
		}
		return permissionDoc(item), nil
	}
	return nil, nil
}

func permissionDoc(item permissionmodel.Permission) *searchpkg.Document {
	searchable := strings.ToLower(strings.Join([]string{
		item.Key,
		item.Name,
		item.Group,
		item.Description,
	}, " "))

	return &searchpkg.Document{
		EntityType:     searchpkg.EntityPermission,
		EntityID:       item.ID,
		Title:          item.Name,
		SearchableText: searchable,
		PermissionKeys: []string{rbac.Key("permissions", rbac.ActionView)},
		Metadata: map[string]string{
			"key":   item.Key,
			"group": item.Group,
		},
		UpdatedAt: item.UpdatedAt,
	}
}

func (b *DocumentBuilder) BuildPermissionByKey(ctx context.Context, key string) (*searchpkg.Document, error) {
	item, err := b.permissions.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	doc := permissionDoc(*item)
	return doc, nil
}

func (b *DocumentBuilder) AllRoles(ctx context.Context) ([]rolemodel.Role, error) {
	roles, _, err := b.roles.List(ctx, query.Unbounded().OrderBy("created_at ASC"))
	return roles, err
}

func (b *DocumentBuilder) AllPermissions(ctx context.Context) ([]permissionmodel.Permission, error) {
	return b.permissions.ListAll(ctx)
}

func (b *DocumentBuilder) AllUsers(ctx context.Context) ([]usermodel.User, error) {
	return b.users.ListAll(ctx)
}
