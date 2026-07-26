import { apiGet, apiPost, setToken, clearToken } from '@/lib/api'

export type UserInfo = {
  id: number
  username: string
  email?: string
  role?: string
  roles?: string[]
  display_name?: string
  must_change_password?: boolean
}

export type LoginResult = {
  token: string
  user: UserInfo
  expires_at?: string
  is_new_user?: boolean
}

export type OAuthProviderInfo = {
  name: string
  display_name?: string
  enabled?: boolean
  configured?: boolean
  login_ready?: boolean
  auth_url?: string
}

export type OAuthProvidersResponse = {
  providers: OAuthProviderInfo[]
  allow_registration?: boolean
  auto_link_accounts?: boolean
  default_oauth_role?: string
}

export type LinkedOAuthProvider = {
  provider: string
  provider_user_id?: string
  connected_at?: string
}

const OAUTH_NOTICE_KEY = 'cilikube_oauth_notice'

export type OAuthNotice = {
  kind: 'new_viewer' | 'viewer'
  username: string
}

export function setOAuthNotice(notice: OAuthNotice) {
  sessionStorage.setItem(OAUTH_NOTICE_KEY, JSON.stringify(notice))
}

export function peekOAuthNotice(): OAuthNotice | null {
  const raw = sessionStorage.getItem(OAUTH_NOTICE_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw) as OAuthNotice
  } catch {
    return null
  }
}

export function clearOAuthNotice() {
  sessionStorage.removeItem(OAUTH_NOTICE_KEY)
}

export async function login(username: string, password: string) {
  const data = await apiPost<LoginResult>('/api/v1/auth/login', { username, password })
  setToken(data.token)
  return data
}

export async function register(username: string, email: string, password: string) {
  return apiPost<UserInfo>('/api/v1/auth/register', { username, email, password })
}

export async function fetchOAuthProviders() {
  return apiGet<OAuthProvidersResponse>('/api/v1/auth/oauth/providers')
}

export async function completeOAuthLogin(provider: string, code: string, state?: string) {
  const data = await apiPost<LoginResult>('/api/v1/auth/oauth/callback', {
    provider,
    code,
    state: state || '',
  })
  setToken(data.token)
  return data
}

export async function fetchOAuthAuthURL(provider: string, state = 'cilikube_link') {
  return apiGet<{ auth_url: string; state: string }>(
    `/api/v1/auth/oauth/${encodeURIComponent(provider)}/auth`,
    { state },
  )
}

export async function linkOAuthAccount(provider: string, code: string) {
  return apiPost('/api/v1/auth/oauth/link', { provider, code })
}

export async function unlinkOAuthAccount(provider: string) {
  return apiPost('/api/v1/auth/oauth/unlink', { provider })
}

export async function logout() {
  try {
    await apiPost('/api/v1/auth/logout')
  } catch {
    // ignore network errors on logout
  } finally {
    clearToken()
  }
}
