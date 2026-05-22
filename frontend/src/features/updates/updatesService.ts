import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'

export interface UpdateEntryRow {
  pkg?: string
  version?: string
  description?: string
  url?: string
  outdated?: boolean
}

export async function checkUpdates(): Promise<UpdateEntryRow[]> {
  const response = await apiClient.get<HmdmEnvelope<UpdateEntryRow[]>>('/private/update/check')
  const data = unwrapHmdmData(response.data, 'Failed to check updates.')
  return Array.isArray(data) ? data : []
}

export async function applyUpdates(updates: UpdateEntryRow[]): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/private/update', {
    updates,
    update: true,
    sendStats: false,
  })
  assertHmdmOk(response.data, 'Failed to apply updates.')
}
