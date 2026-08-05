import { apiGet } from '@/lib/api'

export type TimelineSegment = {
  from: string
  to: string
  status: string
  reason?: string
}

export type TimelineEventMarker = {
  at: string
  type: string
  reason: string
  message: string
  marker: string
}

export type TimelineRow = {
  kind: string
  namespace: string
  name: string
  uid?: string
  appGroup: string
  href: string
  segments: TimelineSegment[]
  events: TimelineEventMarker[]
  eventBadges: Record<string, number>
}

export type TimelineGroup = {
  name: string
  rows: TimelineRow[]
}

export type TimelineResponse = {
  clusterId: string
  from: string
  to: string
  groupBy: string
  groups: TimelineGroup[]
  sampling: {
    enabled: boolean
    heartbeatSeconds: number
    retentionDays: number
    provisional: boolean
    scanSeconds: number
  }
}

export type TimelineMeta = {
  enabled: boolean
  scanSeconds: number
  heartbeatSeconds: number
  retentionDays: number
  sampleCount: number
  lastOk: Record<string, string>
  lastError: Record<string, string>
}

export function getTimeline(params: {
  namespace?: string
  window?: '15m' | '1h' | '6h'
  groupBy?: 'app' | 'none'
  q?: string
  kinds?: string[]
}) {
  return apiGet<TimelineResponse>('/api/v1/timeline', {
    namespace: params.namespace || undefined,
    window: params.window || '15m',
    groupBy: params.groupBy || 'app',
    q: params.q || undefined,
    kinds: params.kinds?.length ? params.kinds.join(',') : undefined,
  })
}

export function getTimelineMeta() {
  return apiGet<TimelineMeta>('/api/v1/timeline/meta')
}
