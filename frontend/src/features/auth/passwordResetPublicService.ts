import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'

export async function requestPasswordRecovery(loginOrEmail: string): Promise<void> {
  const user = encodeURIComponent(loginOrEmail.trim())
  const response = await apiClient.get<HmdmEnvelope<unknown>>(`/public/passwordReset/recover/${user}`)
  assertHmdmOk(response.data, 'Recovery request failed.')
}

/** Settings snapshot for enforcing password complexity during reset (`passwordLength`, `passwordStrength`). */
export async function fetchResetSettings(token: string): Promise<Record<string, unknown>> {
  const safe = encodeURIComponent(token)
  const response = await apiClient.get<HmdmEnvelope<Record<string, unknown>>>(
    `/public/passwordReset/settings/${safe}`
  )
  const data = unwrapHmdmData(response.data, 'This reset link is invalid or expired.')
  return typeof data === 'object' && data ? data : {}
}

export async function submitPasswordReset(token: string, newPasswordMd5Upper: string): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/public/passwordReset/reset', {
    passwordResetToken: token,
    newPassword: newPasswordMd5Upper,
  })
  assertHmdmOk(response.data, 'Could not reset password.')
}
