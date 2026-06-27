package rbac_test

import (
	"testing"

	"be/internal/common/rbac"
)

func TestAllowedViewAndModify(t *testing.T) {
	granted := []string{"users:view", "roles:modify"}

	if !rbac.Allowed(granted, rbac.Key("users", rbac.ActionView)) {
		t.Fatal("expected users:view")
	}
	if rbac.Allowed(granted, rbac.Key("users", rbac.ActionModify)) {
		t.Fatal("did not expect users:modify")
	}
	if !rbac.Allowed(granted, rbac.Key("roles", rbac.ActionModify)) {
		t.Fatal("expected roles:modify")
	}
}

func TestAllowedLegacyReadMapsToView(t *testing.T) {
	granted := []string{"dashboard:read"}

	if !rbac.Allowed(granted, rbac.Key("dashboard", rbac.ActionView)) {
		t.Fatal("expected legacy read to satisfy view")
	}
}

func TestAllowedLegacyCrudMapsToModify(t *testing.T) {
	granted := []string{"users:update"}

	if !rbac.Allowed(granted, rbac.Key("users", rbac.ActionModify)) {
		t.Fatal("expected legacy update to satisfy modify")
	}
}
