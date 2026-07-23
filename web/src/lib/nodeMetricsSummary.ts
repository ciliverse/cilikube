export type NodeMetricsRow = {
  nodeName?: string
  cpuPercent?: string
  memoryPercent?: string
  cpuCapacity?: string
  memoryCapacity?: string
  cpuRequests?: string
  cpuRequestsPercent?: string
  memoryRequests?: string
  memoryRequestsPercent?: string
  cpuCores?: string
  memoryBytes?: string
}

export type NodeMetricsSummary = {
  totalNodes: number
  totalCpuCores: number
  totalMemoryGB: number
  avgCpuUsagePercent: number
  avgMemoryUsagePercent: number
  totalCpuRequests: number
  totalCpuRequestsPercent: number
  totalMemoryRequests: number
  totalMemoryRequestsPercent: number
}

function parseNum(v: unknown) {
  if (v == null) return 0
  const n = parseFloat(String(v).replace(/%/g, '').replace(/Gi/gi, '').replace(/[^\d.-]/g, ''))
  return Number.isFinite(n) ? n : 0
}

function parseMemoryGi(v: unknown) {
  const s = String(v || '')
  if (!s) return 0
  if (/Gi/i.test(s)) return parseNum(s)
  if (/Mi/i.test(s)) return parseNum(s) / 1024
  if (/Ti/i.test(s)) return parseNum(s) * 1024
  // raw bytes
  const bytes = parseFloat(s)
  if (Number.isFinite(bytes) && bytes > 1024) return bytes / (1024 ** 3)
  return parseNum(s)
}

/** Aggregate per-node metrics-server rows into cluster rollup (cilikube-web style). */
export function parseNodeMetricsSummary(nodes: NodeMetricsRow[]): NodeMetricsSummary {
  if (!nodes.length) {
    return {
      totalNodes: 0,
      totalCpuCores: 0,
      totalMemoryGB: 0,
      avgCpuUsagePercent: 0,
      avgMemoryUsagePercent: 0,
      totalCpuRequests: 0,
      totalCpuRequestsPercent: 0,
      totalMemoryRequests: 0,
      totalMemoryRequestsPercent: 0,
    }
  }

  let totalCpuCapacity = 0
  let totalMemoryGi = 0
  let totalCpuUsagePercent = 0
  let totalMemoryUsagePercent = 0
  let totalCpuRequests = 0
  let totalCpuRequestsPercent = 0
  let totalMemoryRequests = 0
  let totalMemoryRequestsPercent = 0

  for (const node of nodes) {
    totalCpuCapacity += parseNum(node.cpuCapacity || node.cpuCores)
    totalMemoryGi += parseMemoryGi(node.memoryCapacity || node.memoryBytes)
    totalCpuUsagePercent += parseNum(node.cpuPercent)
    totalMemoryUsagePercent += parseNum(node.memoryPercent)
    totalCpuRequests += parseNum(node.cpuRequests)
    totalCpuRequestsPercent += parseNum(node.cpuRequestsPercent)
    totalMemoryRequests += parseMemoryGi(node.memoryRequests)
    totalMemoryRequestsPercent += parseNum(node.memoryRequestsPercent)
  }

  const n = nodes.length
  const round2 = (v: number) => Math.round(v * 100) / 100

  return {
    totalNodes: n,
    totalCpuCores: round2(totalCpuCapacity),
    totalMemoryGB: round2(totalMemoryGi),
    avgCpuUsagePercent: round2(totalCpuUsagePercent / n),
    avgMemoryUsagePercent: round2(totalMemoryUsagePercent / n),
    totalCpuRequests: round2(totalCpuRequests),
    totalCpuRequestsPercent: round2(totalCpuRequestsPercent / n),
    totalMemoryRequests: round2(totalMemoryRequests),
    totalMemoryRequestsPercent: round2(totalMemoryRequestsPercent / n),
  }
}
