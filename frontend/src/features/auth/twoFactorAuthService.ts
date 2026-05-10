import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk } from '@/services/hmdmEnvelope'

export async function verifyTwoFactor(userId: number, code: string): Promise<void> {
  const response = await apiClient.get<HmdmEnvelope<unknown>>(`/private/twofactor/verify/${userId}/${code}`)
  assertHmdmOk(response.data, 'Invalid verification code.')
}

export async function enableTwoFactorAfterVerify(): Promise<void> {
  const response = await apiClient.get<HmdmEnvelope<unknown>>('/private/twofactor/set')
  assertHmdmOk(response.data, 'Could not finalize two-factor setup.')
}

export async function resetTwoFactor(): Promise<void> {
  const response = await apiClient.get<HmdmEnvelope<unknown>>('/private/twofactor/reset')
  assertHmdmOk(response.data, 'Could not reset two-factor authentication.')
}

/** PNG bytes for authenticated user (requires Bearer header). */
export async function fetchTwoFactorQrPngBlob(userId: number): Promise<Blob> {
  const response = await apiClient.get(`/private/twofactor/qr/${userId}`, {
    responseType: 'blob',
    params: { t: Date.now() },
  })
  return response.data as Blob
}
