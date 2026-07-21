import { apiPost, setToken, clearToken } from '@/lib/api'

export type UserInfo = {
  id: number
  username: string
  email?: string
  role?: string
  roles?: string[]
}

export type LoginResult = {
  token: string
  user: UserInfo
  expires_at?: string
}

export async function login(username: string, password: string) {
  const data = await apiPost<LoginResult>('/api/v1/auth/login', { username, password })
  setToken(data.token)
  return data
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
