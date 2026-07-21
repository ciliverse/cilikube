import { request } from '@/utils/http-client'

export interface ClusterEvent {
  id: string
  name: string
  namespace: string
  reason: string
  message: string
  type: string
  source: string
  object: string
  objectKind: string
  count: number
  firstTime: string
  lastTime: string
  createdAt: string
}

export interface EventListParams {
  namespace?: string
  type?: string
  limit?: number
  since?: string
}

export interface EventListResponse {
  events: ClusterEvent[]
  total: number
}

export interface RecentEventsResponse {
  events: ClusterEvent[]
  total: number
}

// 获取集群事件列表
export function getClusterEvents(params?: EventListParams) {
  return request<EventListResponse>({
    url: '/api/v1/events',
    method: 'get',
    params
  })
}

// 获取最近的集群事件（用于仪表板）
export function getRecentEvents(limit?: number) {
  return request<RecentEventsResponse>({
    url: '/api/v1/events/recent',
    method: 'get',
    params: { limit }
  })
}

// 获取特定对象的相关事件
export function getEventsByObject(kind: string, name: string, namespace?: string) {
  return request<{
    events: ClusterEvent[]
    total: number
    object: {
      kind: string
      name: string
      namespace?: string
    }
  }>({
    url: `/api/v1/events/object/${kind}/${name}`,
    method: 'get',
    params: { namespace }
  })
}

// 获取所有命名空间
export function getAllNamespaces() {
  return request<{
    namespaces: string[]
    total: number
  }>({
    url: '/api/v1/namespaces',
    method: 'get'
  })
}