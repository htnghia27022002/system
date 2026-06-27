package seeders

import "be/internal/common/rbac"

// Stable IDs for dev/test references.
const (
	RoleAdminID = "11111111-1111-4111-8111-111111111111"
	RoleUserID  = "22222222-2222-4222-8222-222222222222"
	AdminUserID = "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	DemoUserID  = "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
)

// PermissionDef describes one permission catalog row.
type PermissionDef struct {
	ID          string
	Key         string
	Name        string
	Group       string
	Description string
}

// DefaultPermissions returns the RBAC permission catalog (view/modify model).
func DefaultPermissions() []PermissionDef {
	return []PermissionDef{
		{
			ID:          "10000001-0000-4000-8000-000000000011",
			Key:         rbac.Key("dashboard", rbac.ActionView),
			Name:        "View dashboard",
			Group:       "dashboard",
			Description: "View the admin dashboard overview",
		},
		{
			ID:          "10000001-0000-4000-8000-000000000012",
			Key:         rbac.Key("users", rbac.ActionView),
			Name:        "View users",
			Group:       "users",
			Description: "List and view user accounts",
		},
		{
			ID:          "10000001-0000-4000-8000-000000000013",
			Key:         rbac.Key("users", rbac.ActionModify),
			Name:        "Modify users",
			Group:       "users",
			Description: "Create, update, and delete user accounts",
		},
		{
			ID:          "10000001-0000-4000-8000-000000000014",
			Key:         rbac.Key("roles", rbac.ActionView),
			Name:        "View roles",
			Group:       "roles",
			Description: "List and view roles",
		},
		{
			ID:          "10000001-0000-4000-8000-000000000015",
			Key:         rbac.Key("roles", rbac.ActionModify),
			Name:        "Modify roles",
			Group:       "roles",
			Description: "Create, update, and delete roles",
		},
		{
			ID:          "10000001-0000-4000-8000-000000000016",
			Key:         rbac.Key("permissions", rbac.ActionView),
			Name:        "View permissions",
			Group:       "permissions",
			Description: "View the permissions catalog",
		},
	}
}

// PermissionIDByKey maps permission key to stable UUID for role assignment.
func PermissionIDByKey(key string) string {
	for _, item := range DefaultPermissions() {
		if item.Key == key {
			return item.ID
		}
	}
	return ""
}
