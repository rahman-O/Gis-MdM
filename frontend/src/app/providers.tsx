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

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setTokenState] = useState<string | null>(() => getToken())
  const [username, setUsername] = useState<string | null>(null)
  const navigate = useNavigate()

  // Hydrate username from token on mount (token is opaque, so we store login name separately)
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
      const response = await authService.login(credentials)
      setToken(response.authToken)
      setTokenState(response.authToken)
      setUsername(credentials.login)
      localStorage.setItem('hmdm_username', credentials.login)
      navigate('/dashboard')
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
