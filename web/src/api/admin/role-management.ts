import { request } from "@/utils/service"

/** 角色接口类型定义 */
export interface Role {
  id: number
  name: string
  display_name: string
  description: string
  is_system: boolean
  is_active: boolean
  user_count: number
  permission_count: number
  key_permissions: string[]
  created_at: string
  updated_at: string
}

export interface CreateRoleRequest {
  name: string
  display_name: string
  description: string
  is_active: boolean
  permissions?: string[]
}

export interface UpdateRoleRequest {
  name: string
  display_name: string
  description: string
  is_active: boolean
}

export interface RolePermissions {
  role_id: number
  permissions: string[]
}

export interface UpdateRolePermissionsRequest {
  permissions: string[]
}

export interface RoleUser {
  id: number
  username: string
  display_name: string
  email: string
  avatar_url: string
  is_active: boolean
  is_admin: boolean
  last_login_at: string
  created_at: string
}

export interface Permission {
  name: string
  display_name: string
  description: string
  system_required?: boolean
}

export interface PermissionCategory {
  name: string
  display_name: string
  description: string
  permissions: Permission[]
}

/** API响应类型 */
export interface GetRolesResponse {
  data: Role[]
  total: number
}

export interface GetRoleResponse {
  data: Role
}

export interface CreateRoleResponse {
  data: Role
}

export interface UpdateRoleResponse {
  data: Role
}

export interface GetRolePermissionsResponse {
  data: RolePermissions
}

export interface GetRoleUsersResponse {
  data: {
    users: RoleUser[]
    total: number
  }
}

export interface GetAvailablePermissionsResponse {
  data: {
    categories: PermissionCategory[]
  }
}

/** 角色管理API */

/** 获取角色列表 */
export function getRolesApi(params?: {
  page?: number
  page_size?: number
  search?: string
  is_active?: boolean
}) {
  return request<GetRolesResponse>({
    url: "/api/v1/admin/roles",
    method: "get",
    params
  })
}

/** 获取角色详情 */
export function getRoleApi(id: number) {
  return request<GetRoleResponse>({
    url: `/api/v1/admin/roles/${id}`,
    method: "get"
  })
}

/** 创建角色 */
export function createRoleApi(data: CreateRoleRequest) {
  return request<CreateRoleResponse>({
    url: "/api/v1/admin/roles",
    method: "post",
    data
  })
}

/** 更新角色 */
export function updateRoleApi(id: number, data: UpdateRoleRequest) {
  return request<UpdateRoleResponse>({
    url: `/api/v1/admin/roles/${id}`,
    method: "put",
    data
  })
}

/** 删除角色 */
export function deleteRoleApi(id: number) {
  return request({
    url: `/api/v1/admin/roles/${id}`,
    method: "delete"
  })
}

/** 获取角色权限 */
export function getRolePermissionsApi(roleId: number) {
  return request<GetRolePermissionsResponse>({
    url: `/api/v1/admin/roles/${roleId}/permissions`,
    method: "get"
  })
}

/** 更新角色权限 */
export function updateRolePermissionsApi(roleId: number, data: UpdateRolePermissionsRequest) {
  return request({
    url: `/api/v1/admin/roles/${roleId}/permissions`,
    method: "put",
    data
  })
}

/** 获取角色用户列表 */
export function getRoleUsersApi(roleId: number, params?: {
  page?: number
  page_size?: number
  search?: string
}) {
  return request<GetRoleUsersResponse>({
    url: `/api/v1/admin/roles/${roleId}/users`,
    method: "get",
    params
  })
}

/** 从角色中移除用户 */
export function removeUserFromRoleApi(roleId: number, userId: number) {
  return request({
    url: `/api/v1/admin/roles/${roleId}/users/${userId}`,
    method: "delete"
  })
}

/** 为角色添加用户 */
export function addUserToRoleApi(roleId: number, userId: number) {
  return request({
    url: `/api/v1/admin/roles/${roleId}/users`,
    method: "post",
    data: { user_id: userId }
  })
}

/** 批量为角色添加用户 */
export function addUsersToRoleApi(roleId: number, userIds: number[]) {
  return request({
    url: `/api/v1/admin/roles/${roleId}/users/batch`,
    method: "post",
    data: { user_ids: userIds }
  })
}

/** 获取可用权限列表 */
export function getAvailablePermissionsApi() {
  return request<GetAvailablePermissionsResponse>({
    url: "/api/v1/admin/permissions",
    method: "get"
  })
}

/** 角色状态管理 */
export function toggleRoleStatusApi(id: number, is_active: boolean) {
  return request({
    url: `/api/v1/admin/roles/${id}/status`,
    method: "patch",
    data: { is_active }
  })
}

/** 复制角色 */
export function duplicateRoleApi(id: number, data: { name: string; display_name: string }) {
  return request<CreateRoleResponse>({
    url: `/api/v1/admin/roles/${id}/duplicate`,
    method: "post",
    data
  })
}

/** 获取角色统计信息 */
export function getRoleStatsApi() {
  return request<{
    data: {
      total_roles: number
      active_roles: number
      system_roles: number
      custom_roles: number
    }
  }>({
    url: "/api/v1/admin/roles/stats",
    method: "get"
  })
}