import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'
import type { ManagedRole, RolePermission, RoleSavePayload } from '@/features/roles/types'

function unwrapEnvelope<T>(response: { data: HmdmEnvelope<T> }, message: string): T {
  return unwrapHmdmData(response.data, message)
}

export async function getRolePermissions(): Promise<RolePermission[]> {
  const response = await apiClient.get<HmdmEnvelope<RolePermission[]>>('/private/roles/permissions')
  return unwrapEnvelope(response, 'Failed to load permissions.')
}

export async function getManagedRoles(): Promise<ManagedRole[]> {
  const response = await apiClient.get<HmdmEnvelope<ManagedRole[]>>('/private/roles/all')
  return unwrapEnvelope(response, 'Failed to load roles.')
}

export async function saveRole(role: RoleSavePayload): Promise<void> {
  const response = await apiClient.put<HmdmEnvelope<unknown>>('/private/roles', {
    id: role.id,
    name: role.name.trim(),
    description: role.description?.trim() ? role.description.trim() : null,
    permissions: role.permissions.map((p) => ({ id: p.id })),
  })
  assertHmdmOk(response.data, 'Failed to save role.')
}

export async function deleteRole(id: number): Promise<void> {
  const response = await apiClient.delete<HmdmEnvelope<unknown>>(`/private/roles/${id}`)
  assertHmdmOk(response.data, 'Failed to delete role.')
}
