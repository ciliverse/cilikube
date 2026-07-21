export interface LoginRequestData {
  /** 用户名 */
  username: string
  /** 密码 */
  password: string
}

export interface OAuthLoginData {
  /** OAuth提供商 */
  provider: string
  /** 授权码 */
  code: string
  /** 状态参数 */
  state: string
}

export type LoginResponseData = ApiResponseData<{ 
  token: string
  expires_at: string
  user: {
    id: number
    username: string
    email: string
    display_name: string
    avatar_url: string
    roles: string[]
    is_active: boolean
    email_verified: boolean
    last_login: string | null
    created_at: string
  }
}>

export type UserInfoResponseData = ApiResponseData<{ 
  id: number
  username: string
  email: string
  display_name: string
  avatar_url: string
  roles: string[]
  is_active: boolean
  email_verified: boolean
  last_login: string | null
  created_at: string
}>

export type OAuthAuthUrlResponseData = ApiResponseData<{
  auth_url: string
  state: string
}>
