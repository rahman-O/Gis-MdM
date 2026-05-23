export interface RolePermission {
  id: number
  name: string | null
  description?: string | null
  superAdmin?: boolean
}

/** Matches backend `UserRole` for REST responses. */
export interface ManagedRole {
  id: number | null
  name: string
  description?: string | null
  superAdmin?: boolean
  permissions?: RolePermission[] | null
}

export interface RoleSavePayload {
  id: number | null
  name: string
  description?: string | null
  permissions: Array<{ id: number }>
}
