import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'
import type { LookupItem } from '@/features/devices/types'

function unwrapEnvelope<T>(response: { data: HmdmEnvelope<T> }, message: string): T {
  return unwrapHmdmData(response.data, message)
}

export async function getGroups(): Promise<LookupItem[]> {
  const response = await apiClient.get<HmdmEnvelope<LookupItem[]>>('/private/groups/search')
  return unwrapEnvelope(response, 'Failed to load groups.')
}

export async function createGroup(name: string): Promise<void> {
  const response = await apiClient.put<HmdmEnvelope<unknown>>('/private/groups', { id: null, name })
  assertHmdmOk(response.data, 'Failed to create group.')
}

export async function updateGroup(group: LookupItem): Promise<void> {
  const response = await apiClient.put<HmdmEnvelope<unknown>>('/private/groups', group)
  assertHmdmOk(response.data, 'Failed to update group.')
}

export async function deleteGroup(id: number): Promise<void> {
  const response = await apiClient.delete<HmdmEnvelope<unknown>>(`/private/groups/${id}`)
  assertHmdmOk(response.data, 'Failed to delete group.')
}
