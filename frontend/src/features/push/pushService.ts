import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk } from '@/services/hmdmEnvelope'

export interface PushPayload {
  messageType: string
  payload: string
  deviceNumbers?: string[]
  groups?: string[]
  broadcast?: boolean
}

export async function sendPush(body: PushPayload): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/private/push', body)
  assertHmdmOk(response.data, 'Failed to send push message.')
}
