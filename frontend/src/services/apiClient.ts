import axios from 'axios'
import { clearToken, getToken } from '@/shared/utils/tokenStorage'

const apiClient = axios.create({
  baseURL: '/rest',
  withCredentials: true,
})

apiClient.interceptors.request.use((config) => {
  const isFormData = typeof FormData !== 'undefined' && config.data instanceof FormData
  if (isFormData) {
    // Let browser/axios set multipart boundary automatically.
    if (config.headers) {
      delete (config.headers as Record<string, unknown>)['Content-Type']
      delete (config.headers as Record<string, unknown>)['content-type']
    }
  } else if (config.headers && !('Content-Type' in config.headers) && !('content-type' in config.headers)) {
    ;(config.headers as Record<string, unknown>)['Content-Type'] = 'application/json'
  }
  // Attach Bearer only for JWTs; session login uses the HTTP-only cookie (withCredentials).
  try {
    const token = getToken()
    if (token && shouldAttachBearer(token) && config.headers) {
      if (!('Authorization' in config.headers) && !('authorization' in config.headers)) {
        ;(config.headers as Record<string, unknown>)['Authorization'] = `Bearer ${token}`
      }
    }
  } catch {
    // ignore storage errors
  }
  return config
})

// Response interceptor: handle 401 by clearing token and redirecting to login
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status
    const url: string = error.config?.url ?? ''
    const isPrivate = typeof url === 'string' && url.includes('/private/')
    const tokenPresent = !!getToken()
    if (status === 401 && (isPrivate || tokenPresent)) {
      clearToken()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

/** Session-only marker in localStorage; real auth is the cookie. */
function shouldAttachBearer(token: string): boolean {
  if (token === 'session') return false
  const parts = token.split('.')
  return parts.length === 3 && parts.every((p) => p.length > 0)
}

export default apiClient
