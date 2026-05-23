import React, { createContext, useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import * as authService from '@/services/authService'
import { sessionBootstrapKeysMissing, clearSessionExtras } from '@/features/auth/session'
import { clearToken, getToken, setToken } from '@/shared/utils/tokenStorage'
import type { LoginRequest } from '@/features/auth/types'

interface AuthContextValue {
  token: string | null
  username: string | null
  login: (credentials: LoginRequest) => Promise<void>
  logout: () => Promise<void>
}

export const AuthContext = createContext<AuthContextValue | null>(null)

function isStoredAuthTokenValid(value: string | null): boolean {
  if (!value) return false
  if (value === 'session') return true
  const parts = value.split('.')
  return parts.length === 3 && parts.every((p) => p.length > 0)
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setTokenState] = useState<string | null>(() => {
    const stored = getToken()
    if (!isStoredAuthTokenValid(stored)) {
      if (stored) clearToken()
      return null
    }
    return stored
  })
  const [username, setUsername] = useState<string | null>(null)
  const navigate = useNavigate()

  // Hydrate username from token on mount
  useEffect(() => {
    const stored = localStorage.getItem('hmdm_username')
    if (stored) setUsername(stored)
  }, [])

  /** Legacy SPA users: token existed before RBAC mirrors were persisted */
  useEffect(() => {
    if (!token) return
    if (!sessionBootstrapKeysMissing()) return
    void authService.refreshSessionFromCurrentUser().catch(() => {})
  }, [token])

  const login = useCallback(
    async (credentials: LoginRequest) => {
      const outcome = await authService.login(credentials)
      setToken(outcome.authToken)
      setTokenState(outcome.authToken)
      setUsername(credentials.login)
      localStorage.setItem('hmdm_username', credentials.login)
      navigate(outcome.redirectPath)
    },
    [navigate]
  )

  const logout = useCallback(async () => {
    try {
      await authService.logout()
    } finally {
      clearToken()
      clearSessionExtras()
      localStorage.removeItem('hmdm_username')
      setTokenState(null)
      setUsername(null)
      navigate('/login')
    }
  }, [navigate])

  return (
    <AuthContext.Provider value={{ token, username, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}
