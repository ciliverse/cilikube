import { ref } from "vue"
import store from "@/store"
import { defineStore } from "pinia"
import { useTagsViewStore } from "./tags-view"
import { useSettingsStore } from "./settings"
import { getToken, removeToken, setToken } from "@/utils/cache/cookies"
import { resetRouter } from "@/router"
import { loginApi, getUserInfoApi, getOAuthAuthUrlApi, oauthLoginApi } from "@/api/login"
import { type LoginRequestData, type OAuthLoginData } from "@/api/login/types/login"
import routeSettings from "@/config/route"

export const useUserStore = defineStore("user", () => {
  const token = ref<string>(getToken() || "")
  const roles = ref<string[]>([])
  const username = ref<string>("")
  const userInfo = ref<{
    id?: number
    email?: string
    display_name?: string
    avatar_url?: string
    is_active?: boolean
    email_verified?: boolean
    last_login?: string | null
    created_at?: string
  }>({})

  const tagsViewStore = useTagsViewStore()
  const settingsStore = useSettingsStore()

  /** 登录 */
  const login = async ({ username, password }: LoginRequestData) => {
    const { data } = await loginApi({ username, password })
    setToken(data.token)
    token.value = data.token
    // 获取用户信息和角色
    await getInfo()
  }

  /** OAuth登录 */
  const oauthLogin = async (oauthData: OAuthLoginData) => {
    const { data } = await oauthLoginApi(oauthData)
    setToken(data.token)
    token.value = data.token
    // 获取用户信息和角色
    await getInfo()
  }

  /** 获取OAuth授权URL */
  const getOAuthAuthUrl = async (provider: string) => {
    // 生成更安全的随机state
    const state = crypto.getRandomValues(new Uint8Array(16))
      .reduce((acc, byte) => acc + byte.toString(16).padStart(2, '0'), '')
    
    const { data } = await getOAuthAuthUrlApi(provider, state)
    
    // 存储state用于验证，设置过期时间（10分钟）
    const stateData = {
      state,
      timestamp: Date.now(),
      provider
    }
    localStorage.setItem('oauth_state', JSON.stringify(stateData))
    
    return data
  }
  /** 获取用户详情 */
  const getInfo = async () => {
    const { data } = await getUserInfoApi()
    
    username.value = data.username
    // 验证返回的 roles 是否为一个非空数组，否则塞入一个没有任何作用的默认角色，防止路由守卫逻辑进入无限循环
    roles.value = data.roles?.length > 0 ? data.roles : routeSettings.defaultRoles
    
    // 存储完整的用户信息
    userInfo.value = {
      id: data.id,
      email: data.email,
      display_name: data.display_name,
      avatar_url: data.avatar_url,
      is_active: data.is_active,
      email_verified: data.email_verified,
      last_login: data.last_login,
      created_at: data.created_at
    }
  }
  /** 设置角色 */
  const setRoles = (newRoles: string[]) => {
    roles.value = newRoles
  }
  /** 模拟角色变化 */
  const changeRoles = async (role: string) => {
    const newToken = "token-" + role
    token.value = newToken
    setToken(newToken)
    // 用刷新页面代替重新登录
    window.location.reload()
  }
  /** 登出 */
  const logout = () => {
    removeToken()
    token.value = ""
    roles.value = []
    username.value = ""
    userInfo.value = {}
    resetRouter()
    _resetTagsView()
  }
  /** 重置 Token */
  const resetToken = () => {
    removeToken()
    token.value = ""
    roles.value = []
    username.value = ""
    userInfo.value = {}
  }
  /** 重置 Visited Views 和 Cached Views */
  const _resetTagsView = () => {
    if (!settingsStore.cacheTagsView) {
      tagsViewStore.delAllVisitedViews()
      tagsViewStore.delAllCachedViews()
    }
  }

  return { token, roles, username, userInfo, login, oauthLogin, getOAuthAuthUrl, getInfo, setRoles, changeRoles, logout, resetToken }
})

/** 在 setup 外使用 */
export function useUserStoreHook() {
  return useUserStore(store)
}
