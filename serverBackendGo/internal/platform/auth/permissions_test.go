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

func TestPrincipal_CanEditDevices(t *testing.T) {
	if !(&Principal{SuperAdmin: true}).CanEditDevices() {
		t.Fatal("super admin can edit devices")
	}
	if !(&Principal{Permissions: []string{PermEditDevices}}).CanEditDevices() {
		t.Fatal("named permission")
	}
	if (&Principal{Permissions: []string{"settings"}}).CanEditDevices() {
		t.Fatal("settings alone is not edit_devices")
	}
}

func TestPrincipal_CanManageApplications(t *testing.T) {
	if !(&Principal{SuperAdmin: true}).CanManageApplications() {
		t.Fatal("super admin")
	}
	if !(&Principal{Permissions: []string{PermApplications}}).CanManageApplications() {
		t.Fatal("named permission")
	}
	if (&Principal{Permissions: []string{"settings"}}).CanManageApplications() {
		t.Fatal("settings alone is not applications")
	}
}

func TestPrincipal_CanManageConfigurations(t *testing.T) {
	if !(&Principal{SuperAdmin: true}).CanManageConfigurations() {
		t.Fatal("super admin")
	}
	if !(&Principal{Permissions: []string{PermConfigurations}}).CanManageConfigurations() {
		t.Fatal("named permission")
	}
}

func TestPrincipal_CanBrowseAndEditFiles(t *testing.T) {
	if !(&Principal{SuperAdmin: true}).CanBrowseFiles() || !(&Principal{SuperAdmin: true}).CanEditFiles() {
		t.Fatal("super admin")
	}
	if !(&Principal{Permissions: []string{PermFiles}}).CanBrowseFiles() {
		t.Fatal("files permission")
	}
	if !(&Principal{Permissions: []string{PermEditFiles}}).CanEditFiles() {
		t.Fatal("edit_files permission")
	}
	if (&Principal{Permissions: []string{PermFiles}}).CanEditFiles() {
		t.Fatal("files alone cannot edit")
	}
}

func TestPrincipal_CanUsePushAPI(t *testing.T) {
	if !(&Principal{SuperAdmin: true}).CanUsePushAPI() {
		t.Fatal("super admin")
	}
	if !(&Principal{Permissions: []string{PermPushAPI}}).CanUsePushAPI() {
		t.Fatal("push_api permission")
	}
}

func TestPrincipal_CanPluginPush(t *testing.T) {
	if !(&Principal{Permissions: []string{PermPluginPushSend}}).CanPluginPushSend() {
		t.Fatal("plugin_push_send")
	}
	if !(&Principal{Permissions: []string{PermPluginPushDelete}}).CanPluginPushDelete() {
		t.Fatal("plugin_push_delete")
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
