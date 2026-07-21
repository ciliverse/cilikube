import { request } from "@/utils/service"
import type * as Login from "./types/login"

/** 登录并返回 Token */
export function loginApi(data: Login.LoginRequestData) {
  return request<Login.LoginResponseData>({
    url: "/api/v1/auth/login",
    method: "post",
    data
  })
}

/** 获取OAuth授权URL */
export function getOAuthAuthUrlApi(provider: string, state?: string) {
  return request<Login.OAuthAuthUrlResponseData>({
    url: `/api/v1/auth/oauth/${provider}/auth`,
    method: "get",
    params: { state }
  })
}

/** OAuth登录回调处理 */
export function oauthLoginApi(data: Login.OAuthLoginData) {
  return request<Login.LoginResponseData>({
    url: "/api/v1/auth/oauth/callback",
    method: "post",
    data
  })
}

/** 获取用户详情 */
export function getUserInfoApi() {
  return request<Login.UserInfoResponseData>({
    url: "/api/v1/auth/profile/detailed",
    method: "get"
  })
}
