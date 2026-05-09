import { encodeLoginPassword } from '@/features/auth/loginPasswordEncode'
import type { LoginRequest, LoginResponse, LoginUserPayload } from '@/features/auth/types'
import { applySessionFromUserPayload, type SessionUserPayload } from '@/features/auth/session'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmData } from '@/services/hmdmEnvelope'
import axios from 'axios'
import apiClient from './apiClient'

interface AuthOptionsPayload {
  publicKey?: string | null
}

/**
 * Legacy Angular calls GET `/rest/public/auth/options`. Some proxies or older WARs may return 404;
 * login still works with MD5 password when `transmit.password` is off (default).
 */
async function fetchAuthOptions(): Promise<AuthOptionsPayload> {
  try {
    const response = await apiClient.get<HmdmEnvelope<AuthOptionsPayload>>('/public/auth/options')
    return unwrapHmdmData(response.data, 'Failed to load login options.')
  } catch (err: unknown) {
    if (axios.isAxiosError(err) && err.response?.status === 404) {
      return { publicKey: null }
    }
    throw err
  }
}

function toSessionPayload(login: LoginUserPayload): SessionUserPayload {
  return {
    superAdmin: login.superAdmin,
    userRole: login.userRole ?? undefined,
  }
}

/** Same envelope as login (`User`), used to hydrate stale tokens. */
function toSessionPayloadFromCurrentUser(body: LoginUserPayload): SessionUserPayload {
  return {
    superAdmin: body.superAdmin,
    userRole: body.userRole ?? undefined,
  }
}

export async function login(credentials: LoginRequest): Promise<LoginResponse> {
  const options = await fetchAuthOptions()
  const encodedPassword = encodeLoginPassword(credentials.password, options.publicKey)
  const response = await apiClient.post<HmdmEnvelope<LoginUserPayload>>('/public/auth/login', {
    login: credentials.login,
    password: encodedPassword,
  })
  const data = unwrapHmdmData(response.data, 'Login failed.')
  const token = data.authToken?.trim()
  if (!token) {
    throw new Error('Login failed.')
  }
  applySessionFromUserPayload(toSessionPayload(data), token)
  return { authToken: token }
}

/**
 * Refreshes `hmdm_permissions` and `hmdm_super_admin` from `/private/users/current`.
 * Keeps HTTP-only-visible token in localStorage untouched unless server sends `authToken` (typically unchanged).
 */
export async function refreshSessionFromCurrentUser(): Promise<void> {
  const response = await apiClient.get<HmdmEnvelope<LoginUserPayload>>('/private/users/current')
  const envelope = response.data
  if (envelope.status !== 'OK' || envelope.data == null) {
    return
  }
  applySessionFromUserPayload(toSessionPayloadFromCurrentUser(envelope.data))
}

export async function logout(): Promise<void> {
  await apiClient.post('/public/auth/logout')
}
