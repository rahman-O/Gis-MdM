import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'
import type { LookupItem } from '@/features/devices/types'
import type { Role, User, UserPayload } from '@/features/users/types'
import { encodePasswordForUserSave } from '@/features/users/userPasswordEncode'

type BackendUser = {
  id: number
  login: string
  name: string
  email: string
  userRole?: Role | null
  allDevicesAvailable?: boolean | null
  allConfigAvailable?: boolean | null
  groups?: LookupItem[] | null
  configurations?: LookupItem[] | null
}

function normalizeUser(user: BackendUser): User {
  return {
    id: user.id,
    login: user.login,
    name: user.name,
    email: user.email,
    role: user.userRole ?? null,
    allDevicesAvailable: Boolean(user.allDevicesAvailable),
    allConfigAvailable: user.allConfigAvailable == null ? true : Boolean(user.allConfigAvailable),
    groups: user.groups ?? [],
    configurations: user.configurations ?? [],
  }
}

function unwrapEnvelope<T>(response: { data: HmdmEnvelope<T> }, message: string): T {
  return unwrapHmdmData(response.data, message)
}

function buildUserPutBody(id: number | undefined, data: UserPayload): Record<string, unknown> {
  const body: Record<string, unknown> = {
    login: data.login,
    name: data.name,
    email: data.email,
    userRole: { id: data.roleId },
    allDevicesAvailable: data.allDevicesAvailable,
    allConfigAvailable: data.allConfigAvailable,
  }
  if (id != null) {
    body.id = id
  }
  if (data.password && data.password.trim().length > 0) {
    body.newPassword = encodePasswordForUserSave(data.password.trim())
  }
  if (data.allDevicesAvailable) {
    body.groups = null
  } else {
    body.groups = data.groups.map((g) => ({ id: g.id }))
  }
  if (data.allConfigAvailable) {
    body.configurations = null
  } else {
    body.configurations = data.configurations.map((c) => ({ id: c.id }))
  }
  return body
}

export async function getUsers(): Promise<User[]> {
  const response = await apiClient.get<HmdmEnvelope<BackendUser[]>>('/private/users/all')
  const users = unwrapEnvelope(response, 'Failed to load users.')
  return users.map((user) => normalizeUser(user))
}

export async function createUser(data: UserPayload): Promise<void> {
  const response = await apiClient.put<HmdmEnvelope<unknown>>('/private/users', buildUserPutBody(undefined, data))
  assertHmdmOk(response.data, 'Failed to create user.')
}

export async function updateUser(id: number, data: UserPayload): Promise<void> {
  const response = await apiClient.put<HmdmEnvelope<unknown>>('/private/users', buildUserPutBody(id, data))
  assertHmdmOk(response.data, 'Failed to update user.')
}

export async function deleteUser(id: number): Promise<void> {
  const response = await apiClient.delete<HmdmEnvelope<unknown>>(`/private/users/other/${id}`)
  assertHmdmOk(response.data, 'Failed to delete user.')
}

export async function getRoles(): Promise<Role[]> {
  const response = await apiClient.get<HmdmEnvelope<Role[]>>('/private/users/roles')
  return unwrapEnvelope(response, 'Failed to load roles.')
}
