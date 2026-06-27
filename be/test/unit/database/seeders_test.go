package database_test

import (
	"testing"

	"be/internal/common/rbac"
	"be/internal/database/seeders"
)

func TestDefaultPermissionsUseViewModifyKeys(t *testing.T) {
	t.Parallel()

	perms := seeders.DefaultPermissions()
	if len(perms) != 6 {
		t.Fatalf("expected 6 permissions, got %d", len(perms))
	}

	keys := map[string]bool{}
	for _, p := range perms {
		keys[p.Key] = true
	}

	want := []string{
		rbac.Key("dashboard", rbac.ActionView),
		rbac.Key("users", rbac.ActionView),
		rbac.Key("users", rbac.ActionModify),
		rbac.Key("roles", rbac.ActionView),
		rbac.Key("roles", rbac.ActionModify),
		rbac.Key("permissions", rbac.ActionView),
	}

	for _, key := range want {
		if !keys[key] {
			t.Fatalf("missing permission key %q", key)
		}
	}
}

func TestRegistryContainsDatabaseSeeder(t *testing.T) {
	t.Parallel()

	if _, ok := seeders.Registry["PermissionSeeder"]; !ok {
		t.Fatal("PermissionSeeder not registered")
	}
	if _, ok := seeders.Registry["DatabaseSeeder"]; !ok {
		t.Fatal("DatabaseSeeder not registered")
	}
}
