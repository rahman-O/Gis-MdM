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

/** Mark a hint key as shown (JSON string body per Go handler). */
export async function markHintShown(hintKey: string): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/private/hints/history', JSON.stringify(hintKey), {
    headers: { 'Content-Type': 'application/json' },
  })
  assertHmdmOk(response.data, 'Failed to record hint.')
}
