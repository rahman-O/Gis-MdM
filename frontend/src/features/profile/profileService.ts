import { md5UpperHex } from '@/features/auth/loginPasswordEncode'
import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmData } from '@/services/hmdmEnvelope'
import type { LoginUserPayload } from '@/features/auth/types'

function unwrap<T>(response: { data: HmdmEnvelope<T> }, message: string): T {
  return unwrapHmdmData(response.data, message)
}

export async function fetchCurrentUser(): Promise<LoginUserPayload & Record<string, unknown>> {
  const response = await apiClient.get<HmdmEnvelope<LoginUserPayload & Record<string, unknown>>>('/private/users/current')
  return unwrap(response, 'Failed to load profile.')
}

export async function updateUserDetails(payload: {
  id: number
  name?: string | null
  email?: string | null
}): Promise<unknown> {
  const response = await apiClient.put<HmdmEnvelope<unknown>>('/private/users/details', payload)
  return unwrap(response, 'Failed to update profile.')
}

export async function updateCurrentPassword(payload: {
  id: number
  login?: string | null
  oldPasswordPlain: string
  newPasswordPlain: string
}): Promise<void> {
  const oldMd = md5UpperHex(payload.oldPasswordPlain)
  const newMd = md5UpperHex(payload.newPasswordPlain)
  const response = await apiClient.put<HmdmEnvelope<unknown>>('/private/users/current', {
    id: payload.id,
    login: payload.login ?? undefined,
    oldPassword: oldMd,
    newPassword: newMd,
  })
  unwrap(response, 'Password update failed.')
}
