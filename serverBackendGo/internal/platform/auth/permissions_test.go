package auth

import "testing"

func TestPrincipal_HasPermission_superAdmin(t *testing.T) {
	p := &Principal{SuperAdmin: true, AuthLoaded: true}
	if !p.HasPermission("settings") {
		t.Fatal("super admin should have settings")
	}
}

func TestPrincipal_HasPermission_named(t *testing.T) {
	p := &Principal{Permissions: []string{"settings"}, AuthLoaded: true}
	if !p.HasPermission("settings") {
		t.Fatal("expected settings permission")
	}
	if p.HasPermission("devices") {
		t.Fatal("should not have devices")
	}
}

func TestPrincipal_IsOrgAdmin(t *testing.T) {
	p := &Principal{RoleID: OrgAdminRoleID, AuthLoaded: true}
	if !p.IsOrgAdmin() {
		t.Fatal("role 2 should be org admin")
	}
}

func TestPrincipal_CanManageUsers(t *testing.T) {
	if !(&Principal{SuperAdmin: true}).CanManageUsers() {
		t.Fatal("super admin can manage users")
	}
	if !(&Principal{RoleID: OrgAdminRoleID}).CanManageUsers() {
		t.Fatal("org admin can manage users")
	}
	if (&Principal{RoleID: 3}).CanManageUsers() {
		t.Fatal("regular user cannot manage users")
	}
}
