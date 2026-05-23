import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'

export async function signupVerifyEmail(email: string, language: string): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/public/signup/verifyEmail', {
    email: email.trim().toLowerCase(),
    language: language.trim().toLowerCase().slice(0, 5),
  })
  assertHmdmOk(response.data, 'Signup request failed.')
}

export async function signupVerifyTokenOk(token: string): Promise<boolean> {
  const safe = encodeURIComponent(token)
  const response = await apiClient.get<HmdmEnvelope<unknown>>(`/public/signup/verifyToken/${safe}`)
  return response.data.status === 'OK'
}

export interface SignupCompleteBody {
  token: string
  name: string
  firstName: string
  lastName: string
  company?: string
  description?: string
  passwd: string
}

export async function signupComplete(body: SignupCompleteBody): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/public/signup/complete', body)
  if (response.data.status !== 'OK') {
    throw new Error(response.data.message ?? 'Signup failed.')
  }
}

export async function fetchPendingSignup(token: string): Promise<Record<string, unknown> | null> {
  const safe = encodeURIComponent(token)
  const response = await apiClient.get<HmdmEnvelope<Record<string, unknown>>>(
    `/public/signup/verifyToken/${safe}`
  )
  if (response.data.status !== 'OK') return null
  try {
    return unwrapHmdmData(response.data, 'Invalid signup token.') as Record<string, unknown>
  } catch {
    return null
  }
}
