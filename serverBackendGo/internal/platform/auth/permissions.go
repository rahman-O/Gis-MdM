package auth

import (
	"strings"
)

const (
	// OrgAdminRoleID matches role.orgadmin.id from Java context (default 2).
	OrgAdminRoleID = 2
	PermSettings       = "settings"
	PermEditDevices    = "edit_devices"
	PermEditDeviceDesc = "edit_device_desc"
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
