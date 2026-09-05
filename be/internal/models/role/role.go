package role

import "time"

type Role struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type RolePermission struct {
	RoleID       string    `json:"roleId" db:"role_id"`
	PermissionID string    `json:"permissionId" db:"permission_id"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
