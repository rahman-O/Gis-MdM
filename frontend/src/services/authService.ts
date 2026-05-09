import { encodeLoginPassword } from '@/features/auth/loginPasswordEncode'
import type { LoginRequest, LoginResponse } from '@/features/auth/types'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmData } from '@/services/hmdmEnvelope'
import apiClient from './apiClient'

interface LoginUserPayload {
  authToken?: string | null
}

interface AuthOptionsPayload {
  publicKey?: string | null
}

async function fetchAuthOptions(): Promise<AuthOptionsPayload> {
  const response = await apiClient.get<HmdmEnvelope<AuthOptionsPayload>>('/public/auth/options')
  return unwrapHmdmData(response.data, 'Failed to load login options.')
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
  return { authToken: token }
}

export async function logout(): Promise<void> {
  await apiClient.post('/public/auth/logout')
}
