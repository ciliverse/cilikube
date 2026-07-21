import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import { login as loginApi, logout as logoutApi, type UserInfo } from '@/api/auth'
import { clearToken, getToken } from '@/lib/api'

type AuthContextValue = {
  token: string | null
  user: UserInfo | null
  isAuthenticated: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setTokenState] = useState<string | null>(() => getToken())
  const [user, setUser] = useState<UserInfo | null>(() => {
    const raw = localStorage.getItem('cilikube_user')
    return raw ? (JSON.parse(raw) as UserInfo) : null
  })

  const login = useCallback(async (username: string, password: string) => {
    const result = await loginApi(username, password)
    setTokenState(result.token)
    setUser(result.user)
    localStorage.setItem('cilikube_user', JSON.stringify(result.user))
  }, [])

  const logout = useCallback(async () => {
    await logoutApi()
    clearToken()
    localStorage.removeItem('cilikube_user')
    setTokenState(null)
    setUser(null)
  }, [])

  const value = useMemo(
    () => ({
      token,
      user,
      isAuthenticated: Boolean(token),
      login,
      logout,
    }),
    [token, user, login, logout],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
