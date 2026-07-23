import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import {
  login as loginApi,
  logout as logoutApi,
  type LoginResult,
  type UserInfo,
} from '@/api/auth'
import { clearToken, getToken, setToken } from '@/lib/api'
import {
  canDeleteResource,
  canEdit,
  canExecPods,
  canMutateResource,
  checkResourcePermission,
  isAdmin,
  isViewerOnly,
  userRoles,
  type ResourceAction,
} from '@/lib/permissions'

type AuthContextValue = {
  token: string | null
  user: UserInfo | null
  isAuthenticated: boolean
  roles: string[]
  isAdmin: boolean
  canEdit: boolean
  isViewerOnly: boolean
  canMutate: (resource: string) => boolean
  canDelete: (resource: string) => boolean
  canExec: boolean
  checkPermission: (resource: string, action: ResourceAction) => boolean
  login: (username: string, password: string) => Promise<void>
  applySession: (result: LoginResult) => void
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setTokenState] = useState<string | null>(() => getToken())
  const [user, setUser] = useState<UserInfo | null>(() => {
    const raw = localStorage.getItem('cilikube_user')
    return raw ? (JSON.parse(raw) as UserInfo) : null
  })

  const applySession = useCallback((result: LoginResult) => {
    setToken(result.token)
    setTokenState(result.token)
    setUser(result.user)
    localStorage.setItem('cilikube_user', JSON.stringify(result.user))
  }, [])

  const login = useCallback(
    async (username: string, password: string) => {
      const result = await loginApi(username, password)
      applySession(result)
    },
    [applySession],
  )

  const logout = useCallback(async () => {
    await logoutApi()
    clearToken()
    localStorage.removeItem('cilikube_user')
    setTokenState(null)
    setUser(null)
  }, [])

  const value = useMemo(() => {
    const roles = userRoles(user)
    return {
      token,
      user,
      isAuthenticated: Boolean(token),
      roles,
      isAdmin: isAdmin(user),
      canEdit: canEdit(user),
      isViewerOnly: isViewerOnly(user),
      canMutate: (resource: string) => canMutateResource(user, resource),
      canDelete: (resource: string) => canDeleteResource(user, resource),
      canExec: canExecPods(user),
      checkPermission: (resource: string, action: ResourceAction) =>
        checkResourcePermission(user, resource, action),
      login,
      applySession,
      logout,
    }
  }, [token, user, login, applySession, logout])

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
