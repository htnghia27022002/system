package role

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
