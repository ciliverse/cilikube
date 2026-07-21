/**
 * 统一的 HTTP 客户端封装
 * 整合了之前分散的 axios 配置，提供统一的请求接口
 */
import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from "axios"
import { useUserStoreHook } from "@/store/modules/user"
import { ElMessage } from "element-plus"
import { get, merge } from "lodash-es"
import { getToken } from "./cache/cookies"
import { getCurrentClusterId as getGlobalClusterId } from "./cluster-context"

// ============================== 类型定义 ==============================

export interface HttpClientConfig extends AxiosRequestConfig {
  /** 是否显示错误消息，默认 true */
  showErrorMessage?: boolean
  /** 是否显示加载状态，默认 false */
  showLoading?: boolean
  /** 自定义错误处理 */
  customErrorHandler?: (error: any) => void
}

export interface ApiResponse<T = any> {
  code: number
  data: T
  message: string
}

export interface KubernetesResource {
  apiVersion?: string
  kind?: string
  metadata?: {
    name?: string
    namespace?: string
    labels?: Record<string, string>
    annotations?: Record<string, string>
  }
}

// ============================== 枚举定义 ==============================

export enum ResponseType {
  /** 带 code 字段的标准响应格式 */
  STANDARD = 'standard',
  /** 直接返回数据的格式 */
  DIRECT = 'direct',
  /** 原始响应 */
  RAW = 'raw'
}

export enum RequestType {
  /** 标准 API 请求 */
  API = 'api',
  /** Kubernetes 资源请求 */
  KUBERNETES = 'kubernetes',
  /** 文件上传 */
  UPLOAD = 'upload',
  /** 文件下载 */
  DOWNLOAD = 'download'
}

// ============================== 工具函数 ==============================

/** 获取当前选择的集群ID */
function getCurrentClusterId(): string | null {
  const globalId = getGlobalClusterId()
  if (globalId) return globalId
  return null
}

/** 判断是否为需要集群上下文的资源API */
function isResourceApiUrl(url: string): boolean {
  const resourceApiPatterns = [
    // Core API resources
    /\/api\/v1\/nodes/,
    /\/api\/v1\/pods/,
    /\/api\/v1\/services/,
    /\/api\/v1\/namespaces/,
    /\/api\/v1\/persistentvolumes/,
    /\/api\/v1\/persistentvolumeclaims/,
    /\/api\/v1\/configmaps/,
    /\/api\/v1\/secrets/,
    /\/api\/v1\/endpoints/,
    /\/api\/v1\/replicationcontrollers/,
    // Apps API resources
    /\/apis\/apps\/v1\/deployments/,
    /\/apis\/apps\/v1\/replicasets/,
    /\/apis\/apps\/v1\/statefulsets/,
    /\/apis\/apps\/v1\/daemonsets/,
    // Networking API resources
    /\/apis\/networking\.k8s\.io\/v1\/ingresses/,
    /\/apis\/networking\.k8s\.io\/v1\/networkpolicies/,
    // Storage API resources
    /\/apis\/storage\.k8s\.io\/v1\/storageclasses/,
    // RBAC API resources
    /\/apis\/rbac\.authorization\.k8s\.io\/v1\//,
    // Extensions (deprecated but might still be used)
    /\/apis\/extensions\/v1beta1\//
  ]

  return resourceApiPatterns.some(pattern => pattern.test(url))
}

/** 退出登录并强制刷新页面 */
function logout() {
  useUserStoreHook().logout()
  location.reload()
}

// ============================== HTTP 客户端类 ==============================

export class HttpClient {
  private instance: AxiosInstance
  private responseType: ResponseType

  constructor(
    responseType: ResponseType = ResponseType.STANDARD,
    baseConfig: AxiosRequestConfig = {}
  ) {
    this.responseType = responseType
    this.instance = axios.create({
      timeout: 10000,
      baseURL: import.meta.env.VITE_BASE_API,
      ...baseConfig
    })

    this.setupInterceptors()
  }

  private setupInterceptors() {
    // 请求拦截器
    this.instance.interceptors.request.use(
      (config) => {
        // 添加认证头
        const token = getToken()
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }

        // 设置默认 Content-Type
        if (!config.headers['Content-Type']) {
          config.headers['Content-Type'] = 'application/json'
        }

        // 为资源API请求自动添加clusterId参数
        if (config.url && isResourceApiUrl(config.url)) {
          const clusterId = getCurrentClusterId()
          if (clusterId) {
            config.params = config.params || {}
            config.params.clusterId = clusterId
          }
        }

        return config
      },
      (error) => Promise.reject(error)
    )

    // 响应拦截器
    this.instance.interceptors.response.use(
      (response) => this.handleSuccess(response),
      (error) => this.handleError(error)
    )
  }

  private handleSuccess(response: AxiosResponse) {
    const { data } = response
    const responseType = response.request?.responseType

    // 处理二进制数据
    if (responseType === "blob" || responseType === "arraybuffer") {
      return data
    }

    // 根据响应类型处理数据
    switch (this.responseType) {
      case ResponseType.STANDARD:
        return this.handleStandardResponse(data)
      case ResponseType.DIRECT:
        return data
      case ResponseType.RAW:
        return response
      default:
        return data
    }
  }

  private handleStandardResponse(data: any) {
    const code = data.code

    if (code === undefined) {
      ElMessage.error("响应格式错误：缺少 code 字段")
      return Promise.reject(new Error("响应格式错误"))
    }

    switch (code) {
      case 0:
      case 200:
        return data
      case 401:
        logout()
        return Promise.reject(new Error("未授权"))
      default:
        const message = data.message || "请求失败"
        ElMessage.error(message)
        return Promise.reject(new Error(message))
    }
  }

  private handleError(error: any) {
    const status = get(error, "response.status")
    const responseData = get(error, "response.data", {})
    let message = error.message

    // 尝试从后端响应中获取更具体的错误信息
    if (responseData.error) {
      message = responseData.error
    } else if (responseData.message) {
      message = responseData.message
    }

    // 根据 HTTP 状态码设置错误信息
    switch (status) {
      case 400:
        message = message || "请求参数错误"
        break
      case 401:
        logout()
        return Promise.reject(error)
      case 403:
        message = message || "拒绝访问"
        break
      case 404:
        message = message || "请求地址不存在"
        break
      case 408:
        message = message || "请求超时"
        break
      case 500:
        message = message || "服务器内部错误"
        break
      case 502:
        message = message || "网关错误"
        break
      case 503:
        message = message || "服务不可用"
        break
      case 504:
        message = message || "网关超时"
        break
      default:
        if (error.code === 'ECONNREFUSED') {
          message = "无法连接到服务器"
        }
        break
    }

    ElMessage.error(message)
    return Promise.reject(error)
  }

  /**
   * 通用请求方法
   */
  public request<T = any>(config: HttpClientConfig): Promise<T> {
    const mergedConfig = merge({}, config)
    return this.instance(mergedConfig)
  }

  /**
   * GET 请求
   */
  public get<T = any>(url: string, config?: HttpClientConfig): Promise<T> {
    return this.request<T>({ ...config, method: 'GET', url })
  }

  /**
   * POST 请求
   */
  public post<T = any>(url: string, data?: any, config?: HttpClientConfig): Promise<T> {
    return this.request<T>({ ...config, method: 'POST', url, data })
  }

  /**
   * PUT 请求
   */
  public put<T = any>(url: string, data?: any, config?: HttpClientConfig): Promise<T> {
    return this.request<T>({ ...config, method: 'PUT', url, data })
  }

  /**
   * PATCH 请求
   */
  public patch<T = any>(url: string, data?: any, config?: HttpClientConfig): Promise<T> {
    return this.request<T>({ ...config, method: 'PATCH', url, data })
  }

  /**
   * DELETE 请求
   */
  public delete<T = any>(url: string, config?: HttpClientConfig): Promise<T> {
    return this.request<T>({ ...config, method: 'DELETE', url })
  }

  /**
   * 上传文件
   */
  public upload<T = any>(url: string, formData: FormData, config?: HttpClientConfig): Promise<T> {
    return this.request<T>({
      ...config,
      method: 'POST',
      url,
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data',
        ...config?.headers
      }
    })
  }

  /**
   * 下载文件
   */
  public download(url: string, config?: HttpClientConfig): Promise<Blob> {
    return this.request<Blob>({
      ...config,
      method: 'GET',
      url,
      responseType: 'blob'
    })
  }
}

// ============================== 预配置的实例 ==============================

/** 标准 API 客户端（处理带 code 字段的响应） */
export const standardApiClient = new HttpClient(ResponseType.STANDARD, {
  timeout: 20000 // 增加到20秒，为添加集群操作预留更多时间
})

/** 直接 API 客户端（处理直接返回数据的响应） */
export const directApiClient = new HttpClient(ResponseType.DIRECT, {
  timeout: 10000
})

/** Kubernetes 资源客户端 */
export const kubernetesClient = new HttpClient(ResponseType.STANDARD, {
  timeout: 15000 // K8s 操作可能需要更长时间
})

// ============================== 便捷方法 ==============================

/**
 * 标准 API 请求（兼容旧的 request 函数）
 */
export function request<T = any>(config: HttpClientConfig): Promise<T> {
  return standardApiClient.request<T>(config)
}

/**
 * Go 后端 API 请求（兼容旧的 requestGo 函数）
 */
export function requestGo<T = any>(config: HttpClientConfig): Promise<T> {
  return directApiClient.request<T>(config)
}

/**
 * Kubernetes 资源请求
 */
export async function kubernetesRequest<T = any>(config: {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  data?: any
  params?: any
  headers?: any
}): Promise<T> {
  const clusterId = getCurrentClusterId()

  // 如果需要集群上下文的API，添加集群参数
  if (isResourceApiUrl(config.url) && clusterId) {
    config.params = {
      ...config.params,
      clusterId: clusterId
    }
  }

  return kubernetesClient.request<T>({
    ...config,
    method: config.method || 'GET'
  })
}

// ============================== 导出默认实例 ==============================

export default {
  standard: standardApiClient,
  direct: directApiClient,
  kubernetes: kubernetesClient,
  request,
  requestGo,
  kubernetesRequest
}