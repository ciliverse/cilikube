import { standardApiClient } from '@/utils/http-client'

export interface PrometheusStatus {
  enabled: boolean
  url: string
  healthy: boolean
  error?: string
  data?: unknown
}

export interface PromQueryResult {
  status: string
  data: unknown
  error?: string
  errorType?: string
}

export function getPrometheusStatus() {
  return standardApiClient.get<{ code: number; data: PrometheusStatus; message: string }>(
    '/api/v1/prometheus/status'
  )
}

export function queryPrometheus(params: { query: string; time?: string }) {
  return standardApiClient.get<{ code: number; data: PromQueryResult; message: string }>(
    '/api/v1/prometheus/query',
    { params }
  )
}

export function queryPrometheusRange(params: {
  query: string
  start: string
  end: string
  step?: string
}) {
  return standardApiClient.get<{ code: number; data: PromQueryResult; message: string }>(
    '/api/v1/prometheus/query_range',
    { params }
  )
}
