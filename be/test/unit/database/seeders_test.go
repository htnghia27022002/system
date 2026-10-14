package database_test

import (
	"strings"
	"testing"

	"be/common/rbac"
	"be/internal/database/seeders"
)

func TestDefaultPermissionsAreCoreViewModifyKeys(t *testing.T) {
	t.Parallel()

	perms := seeders.DefaultPermissions()
	want := []string{
		rbac.Key("dashboard", rbac.ActionView),
		rbac.Key("users", rbac.ActionView),
		rbac.Key("users", rbac.ActionModify),
		rbac.Key("roles", rbac.ActionView),
		rbac.Key("roles", rbac.ActionModify),
		rbac.Key("permissions", rbac.ActionView),
		rbac.Key("audit", rbac.ActionView),
		rbac.Key("settings", rbac.ActionView),
		rbac.Key("settings", rbac.ActionModify),
	}
	if len(perms) != len(want) {
		t.Fatalf("expected %d core permissions, got %d", len(want), len(perms))
	}

	keys := map[string]bool{}
	for _, p := range perms {
		keys[p.Key] = true
	}
	for _, key := range want {
		if !keys[key] {
			t.Fatalf("missing permission key %q", key)
		}
	}
}

func TestMergePermissionsAppendsModulesAndRejectsDuplicates(t *testing.T) {
	t.Parallel()

	extra := seeders.PermissionDef{ID: "10000001-0000-4000-8000-000000000099", Key: "widgets:view"}
	all, err := seeders.MergePermissions(seeders.Contribution{Permissions: []seeders.PermissionDef{extra}})
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	if got := all[len(all)-1].Key; got != "widgets:view" {
		t.Fatalf("expected module permission last, got %q", got)
	}

	dupKey := seeders.PermissionDef{ID: "10000001-0000-4000-8000-000000000098", Key: "users:view"}
	if _, err := seeders.MergePermissions(seeders.Contribution{Permissions: []seeders.PermissionDef{dupKey}}); err == nil ||
		!strings.Contains(err.Error(), "duplicate permission key") {
		t.Fatalf("expected duplicate key error, got %v", err)
	}

	dupID := seeders.PermissionDef{ID: seeders.DefaultPermissions()[0].ID, Key: "widgets:modify"}
	if _, err := seeders.MergePermissions(seeders.Contribution{Permissions: []seeders.PermissionDef{dupID}}); err == nil ||
		!strings.Contains(err.Error(), "duplicate permission id") {
		t.Fatalf("expected duplicate id error, got %v", err)
	}
}

type namedSeeder struct{ seeders.Seeder }

func (namedSeeder) Name() string { return "WidgetSeeder" }

func TestRegistryContainsCoreAndModuleSeeders(t *testing.T) {
	t.Parallel()

	registry := seeders.Registry(seeders.Contribution{Seeders: []seeders.Seeder{namedSeeder{}}})
	for _, name := range []string{"DatabaseSeeder", "PermissionSeeder", "RolePermissionSeeder", "WidgetSeeder"} {
		if _, ok := registry[name]; !ok {
			t.Fatalf("%s not registered", name)
		}
	}
}
