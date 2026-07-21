import { standardApiClient } from '@/utils/http-client'
import type { NodeMetricsResponse, NodeMetricsSummary } from '@/types/node-metrics'

/**
 * 获取节点监控数据
 */
export function getNodeMetrics(): Promise<NodeMetricsResponse> {
  return standardApiClient.get<NodeMetricsResponse>('/api/v1/nodes/metrics')
}

/**
 * 解析节点监控数据为摘要信息
 */
export function parseNodeMetricsSummary(response: NodeMetricsResponse): NodeMetricsSummary {
  if (!response.data?.nodes || response.data.nodes.length === 0) {
    return {
      totalNodes: 0,
      totalCpuCores: 0,
      totalMemoryGB: 0,
      avgCpuUsagePercent: 0,
      avgMemoryUsagePercent: 0,
      healthyNodes: 0,
      totalCpuRequests: 0,
      totalCpuRequestsPercent: 0,
      totalMemoryRequests: 0,
      totalMemoryRequestsPercent: 0
    }
  }

  const nodes = response.data.nodes
  let totalCpuCapacity = 0
  let totalMemoryCapacityGB = 0
  let totalCpuUsagePercent = 0
  let totalMemoryUsagePercent = 0
  let totalCpuRequests = 0
  let totalCpuRequestsPercent = 0
  let totalMemoryRequests = 0
  let totalMemoryRequestsPercent = 0

  nodes.forEach(node => {
    // 解析 CPU 容量 (cores)
    const cpuCapacity = parseFloat(node.cpuCapacity)
    totalCpuCapacity += cpuCapacity

    // 解析内存容量 (Gi to GB)
    const memoryCapacityGi = parseFloat(node.memoryCapacity.replace('Gi', ''))
    totalMemoryCapacityGB += memoryCapacityGi * 1.073741824 // Gi to GB

    // 解析 CPU 使用百分比
    const cpuPercent = parseFloat(node.cpuPercent.replace('%', ''))
    totalCpuUsagePercent += cpuPercent

    // 解析内存使用百分比
    const memoryPercent = parseFloat(node.memoryPercent.replace('%', ''))
    totalMemoryUsagePercent += memoryPercent

    // 解析 CPU 请求 (cores)
    const cpuRequests = parseFloat(node.cpuRequests)
    totalCpuRequests += cpuRequests

    // 解析 CPU 请求百分比
    const cpuRequestsPercent = parseFloat(node.cpuRequestsPercent.replace('%', ''))
    totalCpuRequestsPercent += cpuRequestsPercent

    // 解析内存请求 (Gi to GB)
    const memoryRequestsGi = parseFloat(node.memoryRequests.replace('Gi', ''))
    totalMemoryRequests += memoryRequestsGi * 1.073741824 // Gi to GB

    // 解析内存请求百分比
    const memoryRequestsPercent = parseFloat(node.memoryRequestsPercent.replace('%', ''))
    totalMemoryRequestsPercent += memoryRequestsPercent
  })

  // 计算平均值
  const avgCpuUsagePercent = nodes.length > 0 ? totalCpuUsagePercent / nodes.length : 0
  const avgMemoryUsagePercent = nodes.length > 0 ? totalMemoryUsagePercent / nodes.length : 0
  const avgCpuRequestsPercent = nodes.length > 0 ? totalCpuRequestsPercent / nodes.length : 0
  const avgMemoryRequestsPercent = nodes.length > 0 ? totalMemoryRequestsPercent / nodes.length : 0

  return {
    totalNodes: nodes.length,
    totalCpuCores: Math.round(totalCpuCapacity * 100) / 100,
    totalMemoryGB: Math.round(totalMemoryCapacityGB * 100) / 100,
    avgCpuUsagePercent: Math.round(avgCpuUsagePercent * 100) / 100,
    avgMemoryUsagePercent: Math.round(avgMemoryUsagePercent * 100) / 100,
    healthyNodes: nodes.length, // 假设所有返回的节点都是健康的
    totalCpuRequests: Math.round(totalCpuRequests * 100) / 100,
    totalCpuRequestsPercent: Math.round(avgCpuRequestsPercent * 100) / 100,
    totalMemoryRequests: Math.round(totalMemoryRequests * 100) / 100,
    totalMemoryRequestsPercent: Math.round(avgMemoryRequestsPercent * 100) / 100
  }
}