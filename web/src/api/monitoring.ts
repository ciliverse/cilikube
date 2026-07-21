import { standardApiClient } from '@/utils/http-client'

export interface RealTimeMetrics {
  login_attempts_per_minute: number
  failed_logins_per_minute: number
  active_sessions: number
  locked_accounts: number
  security_violations_per_hour: number
  suspicious_activities: number
  active_threats: number
  total_users: number
  active_users: number
  api_requests_per_minute: number
  resource_access_per_minute: number
  permission_denials_per_hour: number
  last_updated: string
  update_interval: number
}

export interface HealthIssue {
  type: string
  severity: string
  description: string
  value: number
  threshold: number
}

export interface SystemHealth {
  status: string
  timestamp: string
  metrics: RealTimeMetrics
  issues: HealthIssue[]
}

export interface MonitoringDashboard {
  metrics: RealTimeMetrics
  health: SystemHealth
  summary: {
    status: string
    total_users: number
    active_users: number
    active_sessions: number
    active_threats: number
    failed_logins_rate: number
    security_violations: number
    last_updated: string
  }
  alerts: {
    critical: number
    warning: number
    info: number
  }
}

export interface MonitoringAlert {
  id: string
  level: string
  type: string
  title: string
  description: string
  timestamp: string
  resolved: boolean
}

export interface MetricsHistory {
  period: string
  interval: string
  timestamps: string[]
  series: {
    login_attempts: number[]
    failed_logins: number[]
    api_requests: number[]
  }
}

export function getMonitoringDashboard() {
  return standardApiClient.get<{ code: number; data: MonitoringDashboard; message: string }>(
    '/api/v1/monitoring/dashboard'
  )
}

export function getSystemHealth() {
  return standardApiClient.get<{ code: number; data: SystemHealth; message: string }>(
    '/api/v1/monitoring/health'
  )
}

export function getMonitoringAlerts(params?: { severity?: string; limit?: number }) {
  return standardApiClient.get<{
    code: number
    data: { alerts: MonitoringAlert[]; count: number; filter: string }
    message: string
  }>('/api/v1/monitoring/alerts', { params })
}

export function getMetricsHistory(params?: { period?: string; interval?: string }) {
  return standardApiClient.get<{ code: number; data: MetricsHistory; message: string }>(
    '/api/v1/monitoring/metrics/history',
    { params }
  )
}
