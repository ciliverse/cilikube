import axios from 'axios'

const TOKEN_KEY = 'cilikube_token'
const CLUSTER_KEY = 'cilikube_cluster'

export const getToken = () => localStorage.getItem(TOKEN_KEY)
export const setToken = (token: string) => localStorage.setItem(TOKEN_KEY, token)
export const clearToken = () => localStorage.removeItem(TOKEN_KEY)

export const getClusterId = () => localStorage.getItem(CLUSTER_KEY) || ''
export const setClusterId = (id: string) => localStorage.setItem(CLUSTER_KEY, id)

export const api = axios.create({
  baseURL: import.meta.env.VITE_BASE_API || '',
  timeout: 15000,
})

api.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  const clusterId = getClusterId()
  const url = config.url || ''
  const needsCluster =
    url.includes('/api/v1/') &&
    !url.includes('/auth/') &&
    !url.includes('/clusters')

  if (clusterId && needsCluster) {
    config.params = { ...config.params, clusterId }
  }

  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      clearToken()
      if (!window.location.pathname.startsWith('/login')) {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  },
)

export type ApiEnvelope<T> = {
  code: number
  data: T
  message: string
}

export async function apiGet<T>(url: string, params?: Record<string, unknown>) {
  const res = await api.get<ApiEnvelope<T>>(url, { params })
  if (res.data.code !== 0 && res.data.code !== 200) {
    throw new Error(res.data.message || 'Request failed')
  }
  return res.data.data
}

export async function apiPost<T>(url: string, data?: unknown) {
  const res = await api.post<ApiEnvelope<T>>(url, data)
  if (res.data.code !== 0 && res.data.code !== 200) {
    throw new Error(res.data.message || 'Request failed')
  }
  return res.data.data
}
