package permission

import "be/internal/common/query"

type PermissionResponse struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

type ListPermissionsQuery struct {
	query.PageParams
	Search string `form:"search" binding:"omitempty,max=200"`
	Group  string `form:"group" binding:"omitempty,max=50"`
}

type PaginatedPermissionsResponse struct {
	Items    []PermissionResponse `json:"items"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
}
