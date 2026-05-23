import type React from 'react'

export interface NavItem {
  label: string
  path: string
  icon: React.ComponentType<{ className?: string }>
  permission?: string
  /** Hidden unless `isSuperAdmin()` (Angular control-panel parity). */
  requiresSuperAdmin?: boolean
  /** Hidden unless user can manage roles (single-tenant or super-admin), matching legacy settings submenu. */
  requiresCanManageRoles?: boolean
  /** Requires one of these permissions when set (OR). */
  anyPermission?: string[]
}
