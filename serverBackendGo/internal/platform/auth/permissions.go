package auth

import (
	"strings"
)

const (
	// OrgAdminRoleID matches role.orgadmin.id from Java context (default 2).
	OrgAdminRoleID = 2
	PermSettings         = "settings"
	PermEditDevices        = "edit_devices"
	PermEditDeviceDesc     = "edit_device_desc"
	PermApplications       = "applications"
	PermConfigurations     = "configurations"
	PermFiles              = "files"
	PermEditFiles          = "edit_files"
	PermPushAPI            = "push_api"
	PermPluginPushSend     = "plugin_push_send"
	PermPluginPushDelete   = "plugin_push_delete"
	PermPluginsCustomerAccess = "plugins_customer_access_management"
	PermPluginAuditAccess     = "plugin_audit_access"
	PermPluginMessagingSend   = "plugin_messaging_send"
	PermPluginMessagingDelete = "plugin_messaging_delete"
	PermPluginDeviceinfoAccess = "plugin_deviceinfo_access"
	PermPluginDevicelogAccess  = "plugin_devicelog_access"
)

// HasPermission returns true if principal is super admin or has the named permission.
func (p *Principal) HasPermission(name string) bool {
	if p == nil {
		return false
	}
	if p.SuperAdmin {
		return true
	}
	want := strings.ToLower(strings.TrimSpace(name))
	for _, n := range p.Permissions {
		if strings.EqualFold(n, want) {
			return true
		}
	}
	return false
}

// IsOrgAdmin matches UserDAO.isOrgAdmin (role id 2).
func (p *Principal) IsOrgAdmin() bool {
	if p == nil {
		return false
	}
	return p.RoleID == OrgAdminRoleID
}

// CanManageUsers returns true for user create/update/delete (super admin or org admin).
func (p *Principal) CanManageUsers() bool {
	if p == nil {
		return false
	}
	return p.SuperAdmin || p.IsOrgAdmin()
}

// CanListUsers requires settings permission (super admin satisfies via HasPermission).
func (p *Principal) CanListUsers() bool {
	return p.HasPermission(PermSettings)
}

// CanManageRoles matches UserRoleDAO.hasAccess for multi-tenant (super admin only)
// or single-customer (super admin or org admin).
func (p *Principal) CanManageRoles(singleCustomer bool) bool {
	if p == nil {
		return false
	}
	if singleCustomer {
		return p.SuperAdmin || p.IsOrgAdmin()
	}
	return p.SuperAdmin
}

// CanEditDevices matches DeviceResource mutation permission.
func (p *Principal) CanEditDevices() bool {
	return p.HasPermission(PermEditDevices)
}

// CanEditDeviceDescription matches description update permission.
func (p *Principal) CanEditDeviceDescription() bool {
	return p.HasPermission(PermEditDeviceDesc)
}

// CanManageApplications matches ApplicationResource mutations.
func (p *Principal) CanManageApplications() bool {
	return p.HasPermission(PermApplications)
}

// CanManageConfigurations matches ConfigurationResource mutations and browse.
func (p *Principal) CanManageConfigurations() bool {
	return p.HasPermission(PermConfigurations)
}

// CanBrowseFiles matches FilesResource list/search.
func (p *Principal) CanBrowseFiles() bool {
	return p.HasPermission(PermFiles)
}

// CanEditFiles matches FilesResource upload/remove/update.
func (p *Principal) CanEditFiles() bool {
	return p.HasPermission(PermEditFiles)
}

// CanUsePushAPI matches PushApiResource.
func (p *Principal) CanUsePushAPI() bool {
	return p.HasPermission(PermPushAPI)
}

// CanPluginPushSend matches PushResource send.
func (p *Principal) CanPluginPushSend() bool {
	return p.HasPermission(PermPluginPushSend)
}

// CanPluginPushDelete matches PushResource delete/purge.
func (p *Principal) CanPluginPushDelete() bool {
	return p.HasPermission(PermPluginPushDelete)
}

func (p *Principal) CanManagePluginsCustomer() bool {
	return p.HasPermission(PermPluginsCustomerAccess)
}

func (p *Principal) CanPluginAuditAccess() bool {
	return p.HasPermission(PermPluginAuditAccess)
}

func (p *Principal) CanPluginMessagingSend() bool {
	return p.HasPermission(PermPluginMessagingSend)
}

func (p *Principal) CanPluginMessagingDelete() bool {
	return p.HasPermission(PermPluginMessagingDelete)
}

func (p *Principal) CanPluginDeviceinfoAccess() bool {
	return p.HasPermission(PermPluginDeviceinfoAccess)
}

func (p *Principal) CanPluginDevicelogAccess() bool {
	return p.HasPermission(PermPluginDevicelogAccess)
}
