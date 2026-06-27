package role

import "time"

type Role struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"type:varchar(50);not null;uniqueIndex"`
	Slug        string    `json:"slug" gorm:"type:varchar(50);not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type RolePermission struct {
	RoleID       string    `json:"roleId" gorm:"type:uuid;primaryKey"`
	PermissionID string    `json:"permissionId" gorm:"type:uuid;primaryKey"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
