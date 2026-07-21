import { request } from "@/utils/service"

// 用户管理相关类型定义
export interface User {
  id: number
  username: string
  email: string
  display_name: string
  avatar_url: string
  role: string
  is_active: boolean
  email_verified: boolean
  last_login: string | null
  created_at: string
}

export interface UserStats {
  total: number
  active: number
  admins: number
  newThisMonth: number
}

export interface CreateUserData {
  username: string
  email: string
  display_name?: string
  password: string
  role: string
  is_active: boolean
  send_welcome_email: boolean
}

export interface UpdateUserData {
  email?: string
  display_name?: string
  role?: string
  is_active?: boolean
}

export interface GetUsersParams {
  page?: number
  page_size?: number
  search?: string
  status?: string
  role?: string
}

export interface GetUsersResponse {
  users: User[]
  total: number
  page: number
  page_size: number
}

export interface LoginRecord {
  id: string
  login_time: string
  ip_address: string
  location?: string
  user_agent: string
  status: "success" | "failed"
  failure_reason?: string
}

export interface Permission {
  resource: string
  actions: string[]
}

export interface UserPermissions {
  role_permissions: Permission[]
  custom_permissions: Permission[]
}

export interface ActivityRecord {
  id: string
  action: string
  description: string
  ip_address: string
  location?: string
  created_at: string
}

export interface ExportUsersParams {
  search?: string
  status?: string
  role?: string
}

/** 获取用户列表 */
export function getUsersApi(params: GetUsersParams) {
  return request<GetUsersResponse>({
    url: "/api/v1/admin/users",
    method: "get",
    params
  })
}

/** 获取用户统计 */
export function getUserStatsApi() {
  return request<UserStats>({
    url: "/api/v1/admin/users/stats",
    method: "get"
  })
}

/** 创建用户 */
export function createUserApi(data: CreateUserData) {
  return request<User>({
    url: "/api/v1/admin/users",
    method: "post",
    data
  })
}

/** 更新用户 */
export function updateUserApi(userId: number, data: UpdateUserData) {
  return request<User>({
    url: `/api/v1/admin/users/${userId}`,
    method: "put",
    data
  })
}

/** 删除用户 */
export function deleteUserApi(userId: number) {
  return request({
    url: `/api/v1/admin/users/${userId}`,
    method: "delete"
  })
}

/** 批量更新用户 */
export function batchUpdateUsersApi(userIds: number[], data: UpdateUserData) {
  return request({
    url: "/api/v1/admin/users/batch",
    method: "put",
    data: {
      user_ids: userIds,
      ...data
    }
  })
}

/** 重置用户密码 */
export function resetUserPasswordApi(userId: number) {
  return request({
    url: `/api/v1/admin/users/${userId}/reset-password`,
    method: "post"
  })
}

/** 发送用户验证邮件 */
export function sendUserVerificationApi(userId: number) {
  return request({
    url: `/api/v1/admin/users/${userId}/send-verification`,
    method: "post"
  })
}

/** 获取用户登录记录 */
export function getUserLoginHistoryApi(userId: number) {
  return request<LoginRecord[]>({
    url: `/api/v1/admin/users/${userId}/login-history`,
    method: "get"
  })
}

/** 获取用户权限信息 */
export function getUserPermissionsApi(userId: number) {
  return request<UserPermissions>({
    url: `/api/v1/admin/users/${userId}/permissions`,
    method: "get"
  })
}

/** 获取用户活动日志 */
export function getUserActivityLogApi(userId: number) {
  return request<ActivityRecord[]>({
    url: `/api/v1/admin/users/${userId}/activity-log`,
    method: "get"
  })
}

/** 导出用户数据 */
export function exportUsersApi(params: ExportUsersParams) {
  return request({
    url: "/api/v1/admin/users/export",
    method: "get",
    params,
    responseType: "blob"
  })
}