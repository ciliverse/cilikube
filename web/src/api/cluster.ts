// [重要] 导入我们的请求方法 request (Go backend returns {code, data, message} format)
import { request } from "@/utils/service"

// ... 接口定义部分保持不变 ...
export interface ClusterInfo {
  id: string
  name: string
  server: string
  version: string
  status: string
  source: string
  environment: string
}

export interface CreateClusterRequest {
  name: string
  kubeconfigData: string
  provider?: string
  description?: string
  environment?: string
  region?: string
}


/** 获取集群列表 */
export function getClusterList() {
  // Go 后端返回 {code: 200, data: [...], message: "..."} 格式
  return request<{ data: ClusterInfo[] }>({
    url: "/api/v1/clusters",
    method: "get",
    timeout: 10000 // 增加超时时间到10秒
  })
}

/** 获取当前活动集群名称 */
export function getActiveCluster() {
  // Go 后端返回 {code: 200, data: "cluster-id", message: "..."} 格式
  return request<{ data: string }>({
    url: "/api/v1/clusters/active",
    method: "get",
    timeout: 10000 // 增加超时时间到10秒
  })
}

/** 设置活动集群 */
export function setActiveCluster(idOrName: string) {
  // 优先使用 ID，向后兼容 name
  const isId = idOrName.length > 10 && idOrName.includes('-') // 简单的ID检测
  return request({
    url: "/api/v1/clusters/active",
    method: "post",
    data: isId ? { id: idOrName } : { name: idOrName },
    timeout: 10000 // 增加超时时间到10秒
  })
}

/** 创建新集群 */
export function createCluster(data: CreateClusterRequest) {
  // 使用更安全的 Base64 编码方式，支持 UTF-8 字符
  const kubeconfigBase64 = btoa(unescape(encodeURIComponent(data.kubeconfigData)))

  const payload = {
    ...data,
    kubeconfigData: kubeconfigBase64
  }

  console.log('发送的 kubeconfig 数据长度:', data.kubeconfigData.length)
  console.log('Base64 编码后长度:', kubeconfigBase64.length)

  return request({
    url: "/api/v1/clusters",
    method: "post",
    data: payload,
    timeout: 15000 // 创建集群可能需要更长时间
  })
}

/** 删除集群 */
export function deleteCluster(clusterId: string) {
  // 现在使用 ID 而不是 name
  return request({
    url: `/api/v1/clusters/${clusterId}`,
    method: "delete",
    timeout: 10000 // 增加超时时间到10秒
  })
}