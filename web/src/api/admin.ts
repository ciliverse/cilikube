import { apiDelete, apiGet, apiPost, apiPut } from '@/lib/api'

export type AdminUser = {
  id: number
  username: string
  email?: string
  display_name?: string
  is_active?: boolean
  roles?: string[]
  created_at?: string
  last_login?: string
}

export type AdminRole = {
  id: number
  name: string
  display_name?: string
  description?: string
  is_system?: boolean
  user_count?: number
  permission_count?: number
  main_permissions?: string[]
}

export type PermissionItem = {
  name: string
  display_name?: string
  description?: string
}

export type PermissionCategory = {
  name: string
  display_name?: string
  description?: string
  permissions: PermissionItem[]
}

export async function listAdminUsers(page = 1, pageSize = 100) {
  const data = await apiGet<{ data: AdminUser[]; pagination?: { total?: number } }>(
    '/api/v1/admin/users',
    { page, page_size: pageSize },
  )
  return Array.isArray(data?.data) ? data.data : []
}

export async function createAdminUser(body: {
  username: string
  email: string
  password: string
  confirmPassword: string
  display_name?: string
  roles?: string
}) {
  return apiPost('/api/v1/admin/users', body)
}

export async function updateAdminUser(
  id: number,
  body: {
    email: string
    display_name?: string
    avatar_url?: string
    is_active?: boolean
    roles?: string[]
  },
) {
  return apiPut(`/api/v1/admin/users/${id}`, body)
}

export async function updateAdminUserStatus(id: number, isActive: boolean) {
  return apiPut(`/api/v1/admin/users/${id}/status`, { is_active: isActive })
}

export async function deleteAdminUser(id: number) {
  return apiDelete(`/api/v1/admin/users/${id}`)
}

export async function listAdminRoles() {
  const data = await apiGet<{ roles: AdminRole[]; total: number }>('/api/v1/admin/roles')
  return data?.roles || []
}

export async function createAdminRole(body: {
  name: string
  display_name: string
  description?: string
  permissions?: string[]
}) {
  return apiPost('/api/v1/admin/roles', body)
}

export async function updateAdminRole(
  id: number,
  body: { display_name: string; description?: string; permissions?: string[] },
) {
  return apiPut(`/api/v1/admin/roles/${id}`, body)
}

export async function deleteAdminRole(id: number) {
  return apiDelete(`/api/v1/admin/roles/${id}`)
}

export async function getAvailablePermissions() {
  const data = await apiGet<{ categories?: PermissionCategory[] } | PermissionCategory[]>(
    '/api/v1/admin/permissions',
  )
  if (Array.isArray(data)) return data
  return data?.categories || []
}

export async function getRolePermissions(roleId: number) {
  const data = await apiGet<{ permissions: Array<string | { name: string }> }>(
    `/api/v1/admin/roles/${roleId}/permissions`,
  )
  const list = data?.permissions || []
  return list.map((p) => (typeof p === 'string' ? p : p.name))
}

export async function setRolePermissions(roleId: number, permissions: string[]) {
  return apiPut(`/api/v1/admin/roles/${roleId}/permissions`, { permissions })
}
