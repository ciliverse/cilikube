import { request } from "@/utils/service"
import type * as Profile from "./types/profile"

/** 获取用户资料 */
export function getUserProfileApi() {
  return request<Profile.UserProfileResponseData>({
    url: "/api/v1/auth/profile/detailed",
    method: "get"
  })
}

/** 更新用户资料 */
export function updateUserProfileApi(data: Profile.UpdateUserProfileData) {
  return request<Profile.UserProfileResponseData>({
    url: "/api/v1/auth/profile",
    method: "put",
    data
  })
}

/** 发送邮箱验证邮件 */
export function sendEmailVerificationApi() {
  return request({
    url: "/api/v1/profile/email/verify",
    method: "post"
  })
}

/** 修改密码 */
export function changePasswordApi(data: Profile.ChangePasswordData) {
  return request({
    url: "/api/v1/auth/change-password",
    method: "post",
    data
  })
}

/** 获取用户会话列表 */
export function getUserSessionsApi() {
  return request<Profile.UserSessionsResponseData>({
    url: "/api/v1/profile/sessions",
    method: "get"
  })
}

/** 注销指定会话 */
export function logoutSessionApi(sessionId: string) {
  return request({
    url: `/api/v1/profile/sessions/${sessionId}`,
    method: "delete"
  })
}

/** 注销所有其他会话 */
export function logoutAllSessionsApi() {
  return request({
    url: "/api/v1/profile/sessions/all",
    method: "delete"
  })
}

/** 获取安全日志 */
export function getSecurityLogsApi() {
  return request<Profile.SecurityLogsResponseData>({
    url: "/api/v1/profile/security-logs",
    method: "get"
  })
}

/** 获取OAuth关联账户 */
export function getOAuthAccountsApi() {
  return request<Profile.OAuthAccountsResponseData>({
    url: "/api/v1/profile/oauth-accounts",
    method: "get"
  })
}

/** 关联OAuth账户 */
export function connectOAuthApi(provider: string, data: Profile.ConnectOAuthData) {
  return request({
    url: `/api/v1/profile/oauth-accounts/${provider}/connect`,
    method: "post",
    data
  })
}

/** 解除OAuth关联 */
export function disconnectOAuthApi(provider: string) {
  return request({
    url: `/api/v1/profile/oauth-accounts/${provider}/disconnect`,
    method: "delete"
  })
}

/** 同步OAuth账户信息 */
export function syncOAuthAccountApi(provider: string) {
  return request<Profile.SyncOAuthAccountResponseData>({
    url: `/api/v1/profile/oauth-accounts/${provider}/sync`,
    method: "post"
  })
}

/** 获取用户偏好设置 */
export function getUserPreferencesApi() {
  return request<Profile.UserPreferencesResponseData>({
    url: "/api/v1/profile/preferences",
    method: "get"
  })
}

/** 更新用户偏好设置 */
export function updateUserPreferencesApi(data: Profile.UpdateUserPreferencesData) {
  return request({
    url: "/api/v1/profile/preferences",
    method: "put",
    data
  })
}

/** 更新头像 */
export function updateAvatarApi(data: Profile.UpdateAvatarData) {
  return request<Profile.UserProfileResponseData>({
    url: "/api/v1/auth/profile",
    method: "put",
    data
  })
}

/** 导出用户数据 */
export function exportUserDataApi() {
  return request<Profile.ExportUserDataResponseData>({
    url: "/api/v1/profile/export",
    method: "get"
  })
}

/** 清除用户缓存 */
export function clearUserCacheApi() {
  return request({
    url: "/api/v1/profile/cache",
    method: "delete"
  })
}

/** 删除用户账户 */
export function deleteUserAccountApi() {
  return request({
    url: "/api/v1/profile/account",
    method: "delete"
  })
}