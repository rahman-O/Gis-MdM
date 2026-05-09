const PERMISSIONS_KEY = 'hmdm_permissions'

export function getPermissions(): string[] {
  if (typeof window === 'undefined') return []
  const raw = window.localStorage.getItem(PERMISSIONS_KEY)
  if (!raw) return []
  try {
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed.map((item) => String(item)) : []
  } catch {
    return []
  }
}

/**
 * Fail-open to preserve existing UX until backend permissions
 * are fully wired into auth/session bootstrap.
 */
export function hasPermission(permission: string): boolean {
  const permissions = getPermissions()
  if (permissions.length === 0) return true
  return permissions.includes(permission)
}
