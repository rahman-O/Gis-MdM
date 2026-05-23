import { encodeLoginPassword } from '@/features/auth/loginPasswordEncode'
import type { AuthLandingOptions, LoginOutcome, LoginRequest, LoginUserPayload } from '@/features/auth/types'
import { applySessionFromUserPayload, type SessionUserPayload } from '@/features/auth/session'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmData } from '@/services/hmdmEnvelope'
import axios from 'axios'
import apiClient from './apiClient'
import { clearToken } from '@/shared/utils/tokenStorage'

interface AuthOptionsPayload extends AuthLandingOptions {}

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
      return { publicKey: null, signup: false, recover: false }
    }
    throw err
  }
}

export async function getAuthLandingOptions(): Promise<AuthLandingOptions> {
  const o = await fetchAuthOptions()
  return {
    signup: Boolean(o.signup),
    recover: Boolean(o.recover),
    publicKey: o.publicKey,
  }
}

function toSessionPayload(login: LoginUserPayload): SessionUserPayload {
  return {
    superAdmin: login.superAdmin,
    singleCustomer: login.singleCustomer,
    userRole: login.userRole ?? undefined,
  }
}

/** Same envelope as login (`User`), used to hydrate stale tokens. */
function toSessionPayloadFromCurrentUser(body: LoginUserPayload): SessionUserPayload {
  return {
    superAdmin: body.superAdmin,
    singleCustomer: body.singleCustomer,
    userRole: body.userRole ?? undefined,
  }
}

function redirectAfterLogin(user: LoginUserPayload): string {
  if (user.passwordReset && user.passwordResetToken) {
    return `/password-reset/${encodeURIComponent(user.passwordResetToken)}`
  }
  if (user.twoFactor === true && user.twoFactorAccepted !== true) {
    return '/twofactor'
  }
  return '/dashboard'
}

/**
 * Go backend expects a signed JWT in `Authorization`, not the DB `authToken` from session login.
 * Prefer `/public/jwt/login`; fall back to session login for legacy Java-only stacks.
 */
export async function login(credentials: LoginRequest): Promise<LoginOutcome> {
  clearToken()
  const options = await fetchAuthOptions()
  const encodedPassword = encodeLoginPassword(credentials.password, options.publicKey)
  const body = { login: credentials.login, password: encodedPassword }

  try {
    const jwtRes = await apiClient.post<{ id_token?: string }>('/public/jwt/login', body)
    const headerAuth = jwtRes.headers?.authorization ?? jwtRes.headers?.Authorization
    const fromHeader =
      typeof headerAuth === 'string' && headerAuth.startsWith('Bearer ')
        ? headerAuth.slice(7).trim()
        : ''
    const idToken = (jwtRes.data?.id_token ?? fromHeader).trim()
    if (idToken) {
      applySessionFromUserPayload({}, idToken)
      const current = await fetchCurrentUserAfterLogin()
      applySessionFromUserPayload(toSessionPayloadFromCurrentUser(current), idToken)
      return { authToken: idToken, redirectPath: redirectAfterLogin(current) }
    }
  } catch (err: unknown) {
    if (!axios.isAxiosError(err) || err.response?.status !== 404) {
      if (axios.isAxiosError(err) && (err.response?.status === 401 || err.response?.status === 403)) {
        throw new Error('Invalid username or password.')
      }
      throw err
    }
  }

  const response = await apiClient.post<HmdmEnvelope<LoginUserPayload>>('/public/auth/login', body)
  const data = unwrapHmdmData(response.data, 'Login failed.')
  // Session cookie is set by the server; do not send DB authToken as Bearer (Go JWT middleware rejects it).
  applySessionFromUserPayload(toSessionPayload(data))
  await refreshSessionFromCurrentUser()
  const token = getSessionAuthMarker()
  return { authToken: token, redirectPath: redirectAfterLogin(data) }
}

async function fetchCurrentUserAfterLogin(): Promise<LoginUserPayload> {
  const response = await apiClient.get<HmdmEnvelope<LoginUserPayload>>('/private/users/current')
  return unwrapHmdmData(response.data, 'Failed to load user profile after login.')
}

/** Marker stored so AuthGuard treats session-only login as authenticated without a fake Bearer token. */
function getSessionAuthMarker(): string {
  return 'session'
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
