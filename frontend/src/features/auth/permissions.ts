import { readStoredSuperAdmin, getStoredPermissions } from '@/features/auth/session'

export function isSuperAdmin(): boolean {
  return readStoredSuperAdmin()
}

/** Same list legacy Angular persisted from `user.userRole.permissions`. */
export function getPermissions(): string[] {
  return getStoredPermissions()
}

/**
 * Super-admin matches Angular: all UI permissions allowed.
 * Non-super-admin: explicit permission names when list is populated.
 * Fail-open when list empty (sessions not warmed) to preserve prior SPA behavior.
 */
export function hasPermission(permission: string): boolean {
  if (readStoredSuperAdmin()) return true
  const permissions = getStoredPermissions()
  if (permissions.length === 0) return true
  return permissions.includes(permission)
}

/** Legacy permission name from Angular ({@code enroll_devices}); super-admin bypasses checks. */
export function canEnrollDevicesViaQr(): boolean {
  return hasPermission('enroll_devices')
}
