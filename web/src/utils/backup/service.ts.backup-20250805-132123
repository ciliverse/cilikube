import axios, { type AxiosInstance, type AxiosRequestConfig } from "axios"
import { useUserStoreHook } from "@/store/modules/user"
import { ElMessage } from "element-plus"
import { get, merge } from "lodash-es"
import { getToken } from "./cache/cookies"
import { getCurrentClusterId as getGlobalClusterId } from "./cluster-context"

/** 退出登录并强制刷新页面（会重定向到登录页） */
function logout() {
  useUserStoreHook().logout()
  location.reload()
}

/** 获取当前选择的集群ID */
function getCurrentClusterId(): string | null {
  // 优先从全局上下文获取
  const globalId = getGlobalClusterId()
  if (globalId) return globalId
  
  // 如果全局上下文为空，可能是应用刚启动，尝试从localStorage获取作为fallback
  // 但这只能获取到名称，无法获取ID，所以这里暂时返回null
  // 实际使用中，应该确保在应用启动时就初始化集群store
  return null
}

/** 判断是否为需要集群上下文的资源API */
function isResourceApiUrl(url: string): boolean {
  // 匹配需要集群上下文的资源API路径
  const resourceApiPatterns = [
    /\/api\/v1\/nodes/,
    /\/api\/v1\/pods/,
    /\/api\/v1\/deployments/,
    /\/api\/v1\/services/,
    /\/api\/v1\/namespaces/,
    /\/api\/v1\/persistentvolumes/,
    /\/api\/v1\/persistentvolumeclaims/,
    /\/api\/v1\/configmaps/,
    /\/api\/v1\/secrets/,
    /\/api\/v1\/ingresses/,
    /\/api\/v1\/daemonsets/,
    /\/api\/v1\/statefulsets/
  ]
  
  return resourceApiPatterns.some(pattern => pattern.test(url))
}

// ==================================================================================
// ======== 第一套：用于 Mock/旧版 API (期望返回带 code 的数据) 的请求实例 =========
// ==================================================================================

/** 创建用于 Mock API 的请求实例 */
function createMockService() {
  const service = axios.create()
  // 请求拦截
  service.interceptors.request.use(
    (config) => {
      // 为资源API请求自动添加clusterId参数
      if (config.url && isResourceApiUrl(config.url)) {
        const clusterId = getCurrentClusterId()
        if (clusterId) {
          // 如果URL已有查询参数，添加到现有参数中
          if (config.params) {
            config.params.clusterId = clusterId
          } else {
            config.params = { clusterId }
          }
        }
      }
      return config
    },
    (error) => Promise.reject(error)
  )
  // 响应拦截（处理带 code 的返回数据）
  service.interceptors.response.use(
    (response) => {
      const apiData = response.data
      const responseType = response.request?.responseType
      if (responseType === "blob" || responseType === "arraybuffer") return apiData
      
      const code = apiData.code
      if (code === undefined) {
        ElMessage.error("非本系统的接口（缺少code字段）")
        return Promise.reject(new Error("非本系统的接口"))
      }
      switch (code) {
        case 0:
        case 200:
          return apiData
        case 401:
          return logout()
        default:
          ElMessage.error(apiData.message || "Error")
          return Promise.reject(new Error(apiData.message || "Error"))
      }
    },
    (error) => {
      const status = get(error, "response.status")
      switch (status) {
        case 400:
          error.message = "请求错误"
          break
        case 401:
          logout()
          break
        case 403:
          error.message = "拒绝访问"
          break
        case 404:
          error.message = "请求地址出错"
          break
        case 408:
          error.message = "请求超时"
          break
        case 500:
          error.message = "服务器内部错误"
          break
        case 501:
          error.message = "服务未实现"
          break
        case 502:
          error.message = "网关错误"
          break
        case 503:
          error.message = "服务不可用"
          break
        case 504:
          error.message = "网关超时"
          break
        case 505:
          error.message = "HTTP 版本不受支持"
          break
        default:
          break
      }
      ElMessage.error(error.message)
      return Promise.reject(error)
    }
  )
  return service
}

/** 创建用于 Mock API 的请求方法 */
function createMockRequest(service: AxiosInstance) {
  return function <T>(config: AxiosRequestConfig): Promise<T> {
    const token = getToken()
    const defaultConfig = {
      headers: {
        Authorization: token ? `Bearer ${token}` : undefined,
        "Content-Type": "application/json"
      },
      timeout: 5000,
      // 在开发环境使用代理，生产环境使用完整URL
      baseURL: import.meta.env.DEV ? '' : import.meta.env.VITE_BASE_API,
      data: {}
    }
    const mergeConfig = merge(defaultConfig, config)
    return service(mergeConfig)
  }
}

/** 用于 Mock/旧版 API 的实例和方法 (导出给登录等旧功能使用) */
const mockService = createMockService()
export const request = createMockRequest(mockService)


// ======================================================================================
// ======== 第二套：用于 Go 后端新 API (直接返回 {data: ...}) 的专用请求实例 =========
// ======================================================================================

/** 创建用于 Go 后端 API 的请求实例 */
function createGoApiService() {
  const service = axios.create()
  // 请求拦截
  service.interceptors.request.use(
    (config) => config,
    (error) => Promise.reject(error)
  )
  // 响应拦截（直接返回数据，不检查 code）
  service.interceptors.response.use(
    (response) => {
      // 对于Go后端接口，我们约定它直接返回数据 { data: [...] }
      // 或者在HTTP层面报错（由下面的错误拦截器处理）。
      return response.data
    },
    (error) => {
      // 错误处理逻辑
      const status = get(error, "response.status")
      const errData = get(error, "response.data", {})
      let message = error.message

      // 尝试从后端返回的JSON中获取更具体的错误信息
      if(errData.error) {
        message = errData.error
      } else if (errData.message) {
        message = errData.message
      }
      
      switch (status) {
        case 400:
          message = message || "请求错误"
          break
        case 401:
          logout()
          break
        case 403:
          message = message || "拒绝访问"
          break
        case 404:
          message = message || "请求地址出错"
          break
        case 500:
          message = message || "服务器内部错误"
          break
        default:
          break
      }
      ElMessage.error(message)
      return Promise.reject(error)
    }
  )
  return service
}

/** 创建用于 Go 后端 API 的请求方法 */
function createGoApiRequest(service: AxiosInstance) {
  return function <T>(config: AxiosRequestConfig): Promise<T> {
    const token = getToken()
    const defaultConfig = {
      headers: {
        Authorization: token ? `Bearer ${token}` : undefined,
        "Content-Type": "application/json"
      },
      timeout: 10000, // 为真实后端设置更长的超时时间
      // 在开发环境使用代理，生产环境使用完整URL
      baseURL: import.meta.env.DEV ? '' : import.meta.env.VITE_BASE_API,
      data: {}
    }
    const mergeConfig = merge(defaultConfig, config)
    return service(mergeConfig)
  }
}

/** 用于 Go 后端 API 的实例和方法 (导出给集群管理等新功能使用) */
const goApiService = createGoApiService()
export const requestGo = createGoApiRequest(goApiService)