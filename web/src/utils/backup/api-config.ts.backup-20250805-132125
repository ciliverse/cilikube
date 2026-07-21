/**
 * API配置和通用请求方法
 */
import { request } from "@/utils/service"
import { useClusterStore } from "@/store/modules/clusterStore"

// API基础配置
export const API_CONFIG = {
  BASE_URL: import.meta.env.VITE_BASE_API || 'http://127.0.0.1:8080',
  TIMEOUT: 10000
}

/**
 * 获取当前选中的集群ID
 */
export function getCurrentClusterId(): string | null {
  const clusterStore = useClusterStore()
  return clusterStore.selectedClusterId || null
}

/**
 * 获取当前选中的集群名称（向后兼容）
 */
export function getCurrentClusterName(): string | null {
  const clusterStore = useClusterStore()
  return clusterStore.selectedClusterName || null
}

/**
 * 通用的Kubernetes资源API请求方法
 */
export async function kubernetesRequest<T>(config: {
  url: string
  method?: 'get' | 'post' | 'put' | 'delete' | 'patch'
  data?: any
  params?: any
  headers?: any
}): Promise<T> {
  const clusterId = getCurrentClusterId()

  // 如果需要集群上下文的API，添加集群参数
  if (isKubernetesResourceUrl(config.url) && clusterId) {
    config.params = {
      ...config.params,
      clusterId: clusterId
    }
  }

  return request<T>({
    ...config,
    baseURL: API_CONFIG.BASE_URL,
    timeout: API_CONFIG.TIMEOUT
  })
}

/**
 * 判断是否是需要集群上下文的Kubernetes资源URL
 */
function isKubernetesResourceUrl(url: string): boolean {
  const resourcePatterns = [
    /\/api\/v1\/namespaces/,
    /\/api\/v1\/pods/,
    /\/api\/v1\/deployments/,
    /\/api\/v1\/services/,
    /\/api\/v1\/nodes/,
    /\/api\/v1\/persistentvolumes/,
    /\/api\/v1\/persistentvolumeclaims/,
    /\/api\/v1\/configmaps/,
    /\/api\/v1\/secrets/,
    /\/api\/v1\/ingresses/,
    /\/apis\/apps\/v1\/deployments/,
    /\/apis\/networking\.k8s\.io\/v1\/ingresses/
  ]

  return resourcePatterns.some(pattern => pattern.test(url))
}

/**
 * 标准化Kubernetes API响应
 */
export interface KubernetesApiResponse<T = any> {
  code: number
  data: T
  message: string
}

/**
 * 标准化的命名空间列表响应
 */
export interface NamespaceListResponse {
  code: number
  data: {
    items: Array<{
      metadata: {
        name: string
        uid: string
        creationTimestamp: string
        labels?: { [key: string]: string }
      }
    }>
  }
  message: string
}

/**
 * 获取命名空间列表的通用方法
 */
export async function fetchNamespaces(): Promise<string[]> {
  try {
    const response = await kubernetesRequest<NamespaceListResponse>({
      url: '/api/v1/namespaces',
      method: 'get'
    })

    if (response.code === 200 && response.data?.items) {
      return response.data.items.map(ns => ns.metadata.name).sort()
    } else {
      throw new Error(response.message || '获取命名空间失败')
    }
  } catch (error: any) {
    console.error('获取命名空间失败:', error)
    throw new Error(error.message || '网络请求失败')
  }
}