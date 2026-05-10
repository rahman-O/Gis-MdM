import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmData } from '@/services/hmdmEnvelope'

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
