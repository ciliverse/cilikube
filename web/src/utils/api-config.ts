/**
 * API 配置和通用请求方法
 * 基于新的统一 HttpClient 重构
 */
import { kubernetesRequest } from './http-client'
import { useClusterStore } from "@/store/modules/clusterStore"

// ============================== API 基础配置 ==============================

export const API_CONFIG = {
  BASE_URL: import.meta.env.VITE_BASE_API || 'http://127.0.0.1:8080',
  TIMEOUT: 10000,
  KUBERNETES_TIMEOUT: 15000
}

// ============================== 集群上下文工具 ==============================

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

// ============================== Kubernetes API 类型定义 ==============================

/**
 * 标准化Kubernetes API响应
 */
export interface KubernetesApiResponse<T = any> {
  code: number
  data: T
  message: string
}

/**
 * Kubernetes 资源元数据
 */
export interface KubernetesMetadata {
  name: string
  uid?: string
  namespace?: string
  creationTimestamp?: string
  labels?: { [key: string]: string }
  annotations?: { [key: string]: string }
}

/**
 * Kubernetes 资源列表响应
 */
export interface KubernetesResourceList<T = any> {
  apiVersion: string
  kind: string
  metadata: {
    resourceVersion: string
    continue?: string
  }
  items: T[]
}

/**
 * 标准化的命名空间列表响应
 */
export interface NamespaceListResponse {
  code: number
  data: {
    items: Array<{
      metadata: KubernetesMetadata
    }>
  }
  message: string
}

// ============================== Kubernetes API 请求方法 ==============================

/**
 * 通用的Kubernetes资源API请求方法
 * 使用新的统一 HttpClient
 */
export async function kubernetesApiRequest<T>(config: {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  data?: any
  params?: any
  headers?: any
}): Promise<T> {
  return kubernetesRequest<T>(config)
}

/**
 * 获取命名空间列表的通用方法
 */
export async function fetchNamespaces(): Promise<string[]> {
  try {
    const response = await kubernetesApiRequest<any>({
      url: '/api/v1/namespaces',
      method: 'GET'
    })

    console.debug('[fetchNamespaces] Full response:', response)
    
    // 后端返回 {code, data: KubernetesResourceList, message}
    const resourceList = response?.data || response
    
    if (resourceList?.items && Array.isArray(resourceList.items)) {
      const namespaceNames = resourceList.items
        .map((ns: any) => ns.metadata?.name || ns.name)
        .filter((name: any) => name) // 过滤掉空名称
        .sort()
      
      console.debug('[fetchNamespaces] Extracted namespaces:', namespaceNames)
      return namespaceNames
    } else {
      throw new Error('命名空间列表格式不正确: ' + JSON.stringify(resourceList))
    }
  } catch (error: any) {
    console.error('获取命名空间失败:', error)
    throw new Error(error.message || '网络请求失败')
  }
}

/**
 * 获取指定命名空间的资源列表
 */
export async function fetchNamespacedResources<T>(
  resourceType: string, 
  namespace?: string,
  params?: Record<string, any>
): Promise<KubernetesResourceList<T>> {
  const url = namespace 
    ? `/api/v1/namespaces/${namespace}/${resourceType}`
    : `/api/v1/${resourceType}`
    
  return kubernetesApiRequest<KubernetesResourceList<T>>({
    url,
    method: 'GET',
    params
  })
}

/**
 * 获取集群级别的资源列表
 */
export async function fetchClusterResources<T>(
  resourceType: string,
  apiVersion: string = 'v1',
  group?: string,
  params?: Record<string, any>
): Promise<KubernetesResourceList<T>> {
  const baseUrl = group ? `/apis/${group}/${apiVersion}` : `/api/${apiVersion}`
  const url = `${baseUrl}/${resourceType}`
    
  return kubernetesApiRequest<KubernetesResourceList<T>>({
    url,
    method: 'GET',
    params
  })
}

/**
 * 创建或更新 Kubernetes 资源
 */
export async function createOrUpdateResource<T>(
  resourceType: string,
  resource: T,
  namespace?: string,
  method: 'POST' | 'PUT' | 'PATCH' = 'POST'
): Promise<T> {
  const url = namespace 
    ? `/api/v1/namespaces/${namespace}/${resourceType}`
    : `/api/v1/${resourceType}`
    
  return kubernetesApiRequest<T>({
    url,
    method,
    data: resource
  })
}

/**
 * 删除 Kubernetes 资源
 */
export async function deleteResource(
  resourceType: string,
  name: string,
  namespace?: string
): Promise<any> {
  const url = namespace 
    ? `/api/v1/namespaces/${namespace}/${resourceType}/${name}`
    : `/api/v1/${resourceType}/${name}`
    
  return kubernetesApiRequest({
    url,
    method: 'DELETE'
  })
}

// ============================== 常用资源快捷方法 ==============================

/**
 * 获取 Pod 列表
 */
export const fetchPods = (namespace?: string, params?: Record<string, any>) =>
  fetchNamespacedResources('pods', namespace, params)

/**
 * 获取 Service 列表
 */
export const fetchServices = (namespace?: string, params?: Record<string, any>) =>
  fetchNamespacedResources('services', namespace, params)

/**
 * 获取 Deployment 列表
 */
export const fetchDeployments = async (namespace?: string, params?: Record<string, any>) => {
  if (namespace && namespace !== '__ALL__') {
    // 获取指定命名空间的 Deployment
    try {
      const response = await kubernetesApiRequest<any>({
        url: `/api/v1/namespaces/${namespace}/deployments`,
        method: 'GET',
        params
      })
      
      // 后端返回 {code, data: KubernetesResourceList, message}
      const deploymentList = response?.data || response
      console.debug(`[fetchDeployments] Namespace: ${namespace}`, deploymentList)
      
      // 确保返回的是 KubernetesResourceList 格式
      if (deploymentList && Array.isArray(deploymentList.items)) {
        return deploymentList
      } else if (Array.isArray(deploymentList)) {
        // 如果直接返回数组，则包装成 KubernetesResourceList
        return {
          apiVersion: 'v1',
          kind: 'DeploymentList',
          metadata: { resourceVersion: '' },
          items: deploymentList
        }
      } else {
        console.warn('[fetchDeployments] Unexpected response format:', deploymentList)
        return {
          apiVersion: 'v1',
          kind: 'DeploymentList',
          metadata: { resourceVersion: '' },
          items: []
        }
      }
    } catch (error) {
      console.error(`[fetchDeployments] Error fetching from ${namespace}:`, error)
      return {
        apiVersion: 'v1',
        kind: 'DeploymentList',
        metadata: { resourceVersion: '' },
        items: []
      }
    }
  } else {
    // 获取所有命名空间的 Deployment
    try {
      // 先获取所有命名空间
      const namespacesResponse = await kubernetesApiRequest<any>({
        url: `/api/v1/namespaces`,
        method: 'GET'
      })
      
      const namespacesData = namespacesResponse?.data || namespacesResponse
      console.debug('[fetchDeployments] Namespaces Response:', namespacesData)

      if (!namespacesData?.items || !Array.isArray(namespacesData.items)) {
        console.warn('[fetchDeployments] No namespaces found in response')
        return { apiVersion: 'v1', kind: 'DeploymentList', metadata: { resourceVersion: '' }, items: [] }
      }

      // 获取所有命名空间，然后并发查询每个命名空间的deployment
      const namespaceList = namespacesData.items as any[]
      console.debug(`[fetchDeployments] Found ${namespaceList.length} namespaces:`, namespaceList.map(n => n.metadata?.name))
      
      const deploymentPromises = namespaceList.map(ns => {
        const nsName = ns.metadata?.name
        if (!nsName) return Promise.resolve([])
        
        return kubernetesApiRequest<any>({
          url: `/api/v1/namespaces/${nsName}/deployments`,
          method: 'GET'
        })
          .then(res => {
            const deploymentList = res?.data || res
            const items = (deploymentList?.items || deploymentList || [])
            if (Array.isArray(items)) {
              console.debug(`[fetchDeployments] Namespace ${nsName}: ${items.length} deployments`)
              return items
            }
            return []
          })
          .catch(err => {
            console.warn(`[fetchDeployments] Error fetching deployments from ${nsName}:`, err)
            return []
          })
      })

      const allDeploymentsArrays = await Promise.all(deploymentPromises)
      const allDeployments = allDeploymentsArrays.flat()
      console.debug(`[fetchDeployments] Total deployments from all namespaces: ${allDeployments.length}`)

      return {
        apiVersion: 'v1',
        kind: 'DeploymentList',
        metadata: { resourceVersion: '' },
        items: allDeployments
      }
    } catch (error) {
      console.error('Error fetching all deployments:', error)
      return { apiVersion: 'v1', kind: 'DeploymentList', metadata: { resourceVersion: '' }, items: [] }
    }
  }
}

/**
 * 获取 Node 列表
 */
export const fetchNodes = (params?: Record<string, any>) =>
  fetchClusterResources('nodes', 'v1', undefined, params)

/**
 * 获取 ConfigMap 列表
 */
export const fetchConfigMaps = (namespace?: string, params?: Record<string, any>) =>
  fetchNamespacedResources('configmaps', namespace, params)

/**
 * 获取 Secret 列表
 */
export const fetchSecrets = (namespace?: string, params?: Record<string, any>) =>
  fetchNamespacedResources('secrets', namespace, params)

// ============================== 导出所有功能 ==============================

// 重新导出 kubernetesRequest 以保持兼容性
export { kubernetesRequest } from './http-client'

export default {
  API_CONFIG,
  getCurrentClusterId,
  getCurrentClusterName,
  kubernetesApiRequest,
  fetchNamespaces,
  fetchNamespacedResources,
  fetchClusterResources,
  createOrUpdateResource,
  deleteResource,
  // 快捷方法
  fetchPods,
  fetchServices,
  fetchDeployments,
  fetchNodes,
  fetchConfigMaps,
  fetchSecrets
}