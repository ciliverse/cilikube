export interface NodeMetrics {
  nodeName: string
  cpuCores: string
  cpuPercent: string
  memoryBytes: string
  memoryPercent: string
  cpuCapacity: string
  memoryCapacity: string
  cpuRequests: string
  cpuRequestsPercent: string
  memoryRequests: string
  memoryRequestsPercent: string
  cpuLimits: string
  cpuLimitsPercent: string
  memoryLimits: string
  memoryLimitsPercent: string
  timestamp: string
}

export interface NodeMetricsResponse {
  code: number
  data: {
    nodes: NodeMetrics[]
    total: number
  }
  message: string
}

export interface NodeMetricsSummary {
  totalNodes: number
  totalCpuCores: number
  totalMemoryGB: number
  avgCpuUsagePercent: number
  avgMemoryUsagePercent: number
  healthyNodes: number
  totalCpuRequests: number
  totalCpuRequestsPercent: number
  totalMemoryRequests: number
  totalMemoryRequestsPercent: number
}