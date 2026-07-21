// 用户资料相关类型定义

export interface UserProfile {
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
  bio?: string
  location?: string
  website?: string
}

export interface UpdateUserProfileData {
  display_name?: string
  email?: string
  bio?: string
  location?: string
  website?: string
}

export type UserProfileResponseData = ApiResponseData<UserProfile>

export interface ChangePasswordData {
  current_password: string
  new_password: string
}

// 用户会话相关
export interface UserSession {
  id: string
  device_type: string
  device_name: string
  ip_address: string
  location?: string
  user_agent: string
  last_activity: string
  is_current: boolean
}

export type UserSessionsResponseData = ApiResponseData<UserSession[]>

// 安全日志相关
export interface SecurityLog {
  id: string
  action: string
  description: string
  ip_address: string
  location?: string
  created_at: string
}

export type SecurityLogsResponseData = ApiResponseData<SecurityLog[]>

// OAuth账户相关
export interface OAuthAccount {
  provider: string
  account_info: {
    username: string
    display_name?: string
    email: string
    avatar_url?: string
  }
  connected_at: string
  last_sync: string
}

export type OAuthAccountsResponseData = ApiResponseData<OAuthAccount[]>

export interface ConnectOAuthData {
  code: string
  state: string
}

export interface SyncOAuthAccountResponseData {
  user_updated: boolean
  user_info?: Partial<UserProfile>
}

// 用户偏好设置相关
export interface UserPreferences {
  theme: string
  language: string
  sidebarCollapsed: boolean
  showBreadcrumb: boolean
  showTabs: boolean
  notifications: {
    email: {
      security: boolean
      system: boolean
      marketing: boolean
    }
    browser: {
      enabled: boolean
      sound: boolean
    }
    frequency: string
  }
  privacy: {
    analytics: boolean
    errorReporting: boolean
    sessionDuration: string
  }
  advanced: {
    developerMode: boolean
    experimentalFeatures: boolean
    apiDebug: boolean
  }
}

export type UserPreferencesResponseData = ApiResponseData<UserPreferences>

export type UpdateUserPreferencesData = Partial<UserPreferences>

// 头像相关
export interface UpdateAvatarData {
  avatar_url: string
  avatar_type?: 'upload' | 'color' | 'preset'
  avatar_config?: {
    background?: string
    text?: string
    initials?: string
  }
}

export interface UpdateAvatarResponseData {
  avatar_url: string
}

// 数据导出相关
export interface ExportUserDataResponseData {
  user_profile: UserProfile
  preferences: UserPreferences
  oauth_accounts: OAuthAccount[]
  security_logs: SecurityLog[]
  sessions: UserSession[]
  export_date: string
}