package role

import "be/pkg/query"

type CreateRoleRequest struct {
	Name           string   `json:"name" binding:"required,min=2,max=50"`
	Slug           string   `json:"slug" binding:"required,min=2,max=50"`
	PermissionKeys []string `json:"permissionKeys" binding:"required"`
}

type UpdateRoleRequest struct {
	Name           *string  `json:"name" binding:"omitempty,min=2,max=50"`
	Slug           *string  `json:"slug" binding:"omitempty,min=2,max=50"`
	PermissionKeys []string `json:"permissionKeys"`
}

type RoleResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Slug           string   `json:"slug"`
	PermissionKeys []string `json:"permissionKeys"`
}

type ListRolesQuery struct {
	query.PageParams
	Search        string `form:"search" binding:"omitempty,max=200"`
	ID            string `form:"id" binding:"omitempty,uuid"`
	PermissionKey string `form:"permissionKey" binding:"omitempty,max=100"`
}

type PaginatedRolesResponse struct {
	Items    []RoleResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}
