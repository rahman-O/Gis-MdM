import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'

export async function fetchHintHistory(): Promise<string[]> {
  const response = await apiClient.get<HmdmEnvelope<string[]>>('/private/hints/history')
  const data = unwrapHmdmData(response.data, 'Failed to load hints.')
  return Array.isArray(data) ? data : []
}

export async function enableHints(): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/private/hints/enable')
  assertHmdmOk(response.data, 'Failed to enable hints.')
}

export async function disableHints(): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/private/hints/disable')
  assertHmdmOk(response.data, 'Failed to disable hints.')
}
