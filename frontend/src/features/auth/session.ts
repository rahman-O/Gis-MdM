import { setToken } from '@/shared/utils/tokenStorage'

const PERMISSIONS_KEY = 'hmdm_permissions'
const SUPER_ADMIN_KEY = 'hmdm_super_admin'
const SINGLE_CUSTOMER_KEY = 'hmdm_single_customer'

/** Payload shape merged from login `UserView` and `/private/users/current` user JSON. */
export interface UserRoleJson {
  superAdmin?: boolean
  permissions?: Array<{ name?: string | null } | null>
}

export interface SessionUserPayload {
  superAdmin?: boolean
  singleCustomer?: boolean
  userRole?: UserRoleJson | null
}

export function permissionNamesFromUserRole(role?: UserRoleJson | null): string[] {
  if (!role?.permissions?.length) return []
  return role.permissions.map((p) => String(p?.name ?? '')).filter(Boolean)
}

export function inferSuperAdmin(payload: SessionUserPayload): boolean {
  if (payload.superAdmin === true) return true
  if (payload.userRole?.superAdmin === true) return true
  return false
}

/** Persist RBAC mirrors of legacy Angular `authService` cookie state + token. */
export function applySessionFromUserPayload(payload: SessionUserPayload, token?: string): void {
  if (typeof window === 'undefined') return
  if (typeof token === 'string' && token.trim()) {
    setToken(token.trim())
  }
  const names = permissionNamesFromUserRole(payload.userRole ?? null)
  const superAdmin = inferSuperAdmin(payload)
  window.localStorage.setItem(PERMISSIONS_KEY, JSON.stringify(names))
  window.localStorage.setItem(SUPER_ADMIN_KEY, superAdmin ? 'true' : 'false')
  if (typeof payload.singleCustomer === 'boolean') {
    window.localStorage.setItem(SINGLE_CUSTOMER_KEY, payload.singleCustomer ? 'true' : 'false')
  }
}

export function clearSessionExtras(): void {
  if (typeof window === 'undefined') return
  window.localStorage.removeItem(PERMISSIONS_KEY)
  window.localStorage.removeItem(SUPER_ADMIN_KEY)
  window.localStorage.removeItem(SINGLE_CUSTOMER_KEY)
}

export function readStoredSuperAdmin(): boolean {
  if (typeof window === 'undefined') return false
  return window.localStorage.getItem(SUPER_ADMIN_KEY) === 'true'
}

export function readStoredSingleCustomer(): boolean {
  if (typeof window === 'undefined') return false
  return window.localStorage.getItem(SINGLE_CUSTOMER_KEY) === 'true'
}

export function getStoredPermissions(): string[] {
  if (typeof window === 'undefined') return []
  const raw = window.localStorage.getItem(PERMISSIONS_KEY)
  if (!raw) return []
  try {
    const parsed = JSON.parse(raw) as unknown
    return Array.isArray(parsed) ? parsed.map((item) => String(item)) : []
  } catch {
    return []
  }
}

/** True until first login persists keys; avoids extra `/users/current` when already warmed. */
export function sessionBootstrapKeysMissing(): boolean {
  if (typeof window === 'undefined') return false
  const sa = window.localStorage.getItem(SUPER_ADMIN_KEY)
  const perm = window.localStorage.getItem(PERMISSIONS_KEY)
  const sc = window.localStorage.getItem(SINGLE_CUSTOMER_KEY)
  return sa === null || perm === null || sc === null
}
