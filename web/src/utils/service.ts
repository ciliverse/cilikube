/**
 * HTTP 服务封装 - 使用统一的 HttpClient
 * 保持向后兼容性，同时提供新的统一接口
 */
import { 
  standardApiClient, 
  directApiClient, 
  kubernetesClient,
  request as newRequest,
  requestGo as newRequestGo,
  kubernetesRequest as newKubernetesRequest,
  type HttpClientConfig
} from './http-client'

// ============================== 向后兼容的导出 ==============================

/**
 * 标准 API 请求（兼容旧版本）
 * @deprecated 请使用 standardApiClient 或 httpClient.standard
 */
export const request = newRequest

/**
 * Go 后端 API 请求（兼容旧版本）
 * @deprecated 请使用 directApiClient 或 httpClient.direct
 */
export const requestGo = newRequestGo

/**
 * Kubernetes 资源请求（兼容旧版本）
 * @deprecated 请使用 kubernetesClient 或 httpClient.kubernetes
 */
export const kubernetesRequest = newKubernetesRequest

// ============================== 新的推荐导出 ==============================

/**
 * HTTP 客户端集合
 * 推荐使用这种方式来访问不同类型的客户端
 */
export const httpClient = {
  /** 标准 API 客户端（处理带 code 字段的响应） */
  standard: standardApiClient,
  /** 直接 API 客户端（处理直接返回数据的响应） */
  direct: directApiClient,
  /** Kubernetes 资源客户端 */
  kubernetes: kubernetesClient
}

// ============================== 便捷方法 ==============================

/**
 * 创建标准 API 请求
 */
export function createStandardRequest<T = any>(config: HttpClientConfig) {
  return standardApiClient.request<T>(config)
}

/**
 * 创建直接 API 请求
 */
export function createDirectRequest<T = any>(config: HttpClientConfig) {
  return directApiClient.request<T>(config)
}

/**
 * 创建 Kubernetes 资源请求
 */
export function createKubernetesRequest<T = any>(config: HttpClientConfig) {
  return kubernetesClient.request<T>(config)
}

// ============================== 默认导出 ==============================

export default httpClient