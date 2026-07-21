// 节点类型定义，适配后端 NodeResponse 结构
export interface NodeAddress {
  type: string
  address: string
}

export interface NodeCondition {
  type: string
  status: string
  lastHeartbeatTime?: string
  lastTransitionTime?: string
  reason?: string
  message?: string
}

export interface NodeInfo {
  kubeletVersion: string
  containerRuntimeVersion: string
  osImage: string
  kernelVersion: string
  kubeProxyVersion: string
  [key: string]: any
}

export interface NodeStatus {
  capacity: {
    cpu: string
    memory: string
  }
  allocatable: {
    cpu: string
    memory: string
  }
  conditions: NodeCondition[]
  addresses: NodeAddress[]
  nodeInfo: NodeInfo
}

export interface Node {
  metadata: {
    name: string
    uid: string
    creationTimestamp: string
    labels?: Record<string, string>
  }
  status: NodeStatus
  spec: {
    podCIDR?: string
  }
}

// 通用 K8s 列表响应
export interface KubernetesListResponse<T> {
  items: T[]
  total: number
}
