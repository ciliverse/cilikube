<template>
  <div class="monitoring-page" v-loading="loading">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ $t('monitoringPage.title') }}</h1>
        <p class="page-subtitle">{{ $t('monitoringPage.subtitle') }}</p>
      </div>
      <div class="header-actions">
        <el-switch
          v-model="autoRefresh"
          :active-text="$t('monitoringPage.autoRefresh')"
          style="margin-right: 12px"
        />
        <el-button :icon="RefreshIcon" @click="refreshAll" :loading="loading">
          {{ $t('actions.refresh') }}
        </el-button>
      </div>
    </div>

    <!-- Node metrics summary -->
    <el-row :gutter="16" class="summary-row">
      <el-col :xs="12" :sm="6" v-for="card in nodeCards" :key="card.label">
        <el-card shadow="hover" class="metric-card">
          <div class="metric-label">{{ card.label }}</div>
          <div class="metric-value">{{ card.value }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="section-row">
      <!-- Security / app monitoring -->
      <el-col :xs="24" :lg="14">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('monitoringPage.securityOverview') }}</span>
              <el-tag :type="healthTagType" size="small">{{ healthStatus }}</el-tag>
            </div>
          </template>
          <el-row :gutter="12">
            <el-col :span="8" v-for="item in securityCards" :key="item.label">
              <div class="stat-block">
                <div class="stat-label">{{ item.label }}</div>
                <div class="stat-value">{{ item.value }}</div>
              </div>
            </el-col>
          </el-row>
          <el-divider />
          <div v-if="healthIssues.length === 0" class="muted">
            {{ $t('monitoringPage.noIssues') }}
          </div>
          <el-alert
            v-for="(issue, idx) in healthIssues"
            :key="idx"
            :title="issue.description"
            :type="issue.severity === 'critical' ? 'error' : issue.severity === 'warning' ? 'warning' : 'info'"
            :closable="false"
            show-icon
            style="margin-bottom: 8px"
          />
        </el-card>
      </el-col>

      <!-- Prometheus status + query -->
      <el-col :xs="24" :lg="10">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('monitoringPage.prometheus') }}</span>
              <el-tag :type="promStatus?.healthy ? 'success' : promStatus?.enabled ? 'danger' : 'info'" size="small">
                {{ prometheusLabel }}
              </el-tag>
            </div>
          </template>
          <p class="muted" v-if="promStatus">
            {{ $t('monitoringPage.prometheusUrl') }}: {{ promStatus.url || '-' }}
          </p>
          <el-alert
            v-if="promStatus?.error"
            :title="promStatus.error"
            type="warning"
            :closable="false"
            show-icon
            style="margin-bottom: 12px"
          />
          <el-form @submit.prevent="runQuery">
            <el-form-item :label="$t('monitoringPage.promQuery')">
              <el-input
                v-model="promQuery"
                placeholder="up"
                clearable
                :disabled="!promStatus?.enabled"
              />
            </el-form-item>
            <el-button
              type="primary"
              :disabled="!promStatus?.enabled || !promQuery.trim()"
              :loading="queryLoading"
              @click="runQuery"
            >
              {{ $t('monitoringPage.runQuery') }}
            </el-button>
          </el-form>
          <pre v-if="queryResult" class="query-result">{{ queryResult }}</pre>
        </el-card>
      </el-col>
    </el-row>

    <!-- Node table + Alerts -->
    <el-row :gutter="16" class="section-row">
      <el-col :xs="24" :lg="14">
        <el-card shadow="never">
          <template #header>{{ $t('monitoringPage.nodeMetrics') }}</template>
          <el-table :data="nodes" stripe border max-height="420" empty-text="-">
            <el-table-column prop="nodeName" :label="$t('monitoringPage.columns.node')" min-width="160" show-overflow-tooltip />
            <el-table-column prop="cpuPercent" :label="$t('monitoringPage.columns.cpu')" width="100" />
            <el-table-column prop="memoryPercent" :label="$t('monitoringPage.columns.memory')" width="100" />
            <el-table-column prop="cpuRequestsPercent" :label="$t('monitoringPage.columns.cpuRequests')" width="120" />
            <el-table-column prop="memoryRequestsPercent" :label="$t('monitoringPage.columns.memoryRequests')" width="130" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="10">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('monitoringPage.alerts') }}</span>
              <span class="muted">{{ alerts.length }}</span>
            </div>
          </template>
          <div v-if="alerts.length === 0" class="muted">{{ $t('monitoringPage.noAlerts') }}</div>
          <div v-else class="alerts-list">
            <div v-for="alert in alerts" :key="alert.id" class="alert-item">
              <el-tag
                size="small"
                :type="alert.level === 'critical' ? 'danger' : alert.level === 'warning' ? 'warning' : 'info'"
              >
                {{ alert.level }}
              </el-tag>
              <div class="alert-body">
                <div class="alert-title">{{ alert.title }}</div>
                <div class="alert-desc">{{ alert.description }}</div>
                <div class="alert-time">{{ formatTime(alert.timestamp) }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { Refresh as RefreshIcon } from '@element-plus/icons-vue'
import {
  getMonitoringDashboard,
  getMonitoringAlerts,
  type MonitoringDashboard,
  type MonitoringAlert,
  type HealthIssue
} from '@/api/monitoring'
import { getPrometheusStatus, queryPrometheus, type PrometheusStatus } from '@/api/prometheus'
import { getNodeMetrics, parseNodeMetricsSummary } from '@/api/node-metrics'
import type { NodeMetrics } from '@/types/node-metrics'

const { t } = useI18n()
const loading = ref(false)
const queryLoading = ref(false)
const autoRefresh = ref(true)
const dashboard = ref<MonitoringDashboard | null>(null)
const alerts = ref<MonitoringAlert[]>([])
const healthIssues = ref<HealthIssue[]>([])
const promStatus = ref<PrometheusStatus | null>(null)
const promQuery = ref('up')
const queryResult = ref('')
const nodes = ref<NodeMetrics[]>([])
const nodeSummary = ref(parseNodeMetricsSummary({ code: 200, data: { nodes: [], total: 0 }, message: '' }))

let timer: ReturnType<typeof setInterval> | null = null

const healthStatus = computed(() => dashboard.value?.health?.status || dashboard.value?.summary?.status || '-')
const healthTagType = computed(() => {
  const status = (healthStatus.value || '').toLowerCase()
  if (status === 'healthy' || status === 'ok') return 'success'
  if (status === 'warning' || status === 'degraded') return 'warning'
  if (status === 'critical' || status === 'unhealthy') return 'danger'
  return 'info'
})

const prometheusLabel = computed(() => {
  if (!promStatus.value) return '-'
  if (!promStatus.value.enabled) return t('monitoringPage.promDisabled')
  return promStatus.value.healthy
    ? t('monitoringPage.promHealthy')
    : t('monitoringPage.promUnhealthy')
})

const nodeCards = computed(() => [
  { label: t('monitoringPage.cards.nodes'), value: String(nodeSummary.value.totalNodes) },
  { label: t('monitoringPage.cards.cpuCores'), value: String(nodeSummary.value.totalCpuCores) },
  { label: t('monitoringPage.cards.avgCpu'), value: `${nodeSummary.value.avgCpuUsagePercent}%` },
  { label: t('monitoringPage.cards.avgMemory'), value: `${nodeSummary.value.avgMemoryUsagePercent}%` }
])

const securityCards = computed(() => {
  const s = dashboard.value?.summary
  return [
    { label: t('monitoringPage.cards.activeSessions'), value: s?.active_sessions ?? '-' },
    { label: t('monitoringPage.cards.activeUsers'), value: s?.active_users ?? '-' },
    { label: t('monitoringPage.cards.activeThreats'), value: s?.active_threats ?? '-' },
    { label: t('monitoringPage.cards.failedLogins'), value: s?.failed_logins_rate ?? '-' },
    { label: t('monitoringPage.cards.violations'), value: s?.security_violations ?? '-' },
    {
      label: t('monitoringPage.cards.alerts'),
      value: dashboard.value
        ? dashboard.value.alerts.critical + dashboard.value.alerts.warning + dashboard.value.alerts.info
        : '-'
    }
  ]
})

const formatTime = (value?: string) => (value ? dayjs(value).format('YYYY-MM-DD HH:mm:ss') : '-')

const refreshAll = async () => {
  loading.value = true
  try {
    const [dashRes, alertRes, promRes, nodeRes] = await Promise.allSettled([
      getMonitoringDashboard(),
      getMonitoringAlerts({ limit: 20 }),
      getPrometheusStatus(),
      getNodeMetrics()
    ])

    if (dashRes.status === 'fulfilled' && dashRes.value.code === 200) {
      dashboard.value = dashRes.value.data
      healthIssues.value = dashRes.value.data.health?.issues || []
    }
    if (alertRes.status === 'fulfilled' && alertRes.value.code === 200) {
      alerts.value = alertRes.value.data.alerts || []
    }
    if (promRes.status === 'fulfilled' && promRes.value.code === 200) {
      promStatus.value = promRes.value.data
    }
    if (nodeRes.status === 'fulfilled' && nodeRes.value.code === 200) {
      nodes.value = nodeRes.value.data.nodes || []
      nodeSummary.value = parseNodeMetricsSummary(nodeRes.value)
    }
  } catch (error: any) {
    ElMessage.error(error.message || t('monitoringPage.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const runQuery = async () => {
  if (!promQuery.value.trim()) return
  queryLoading.value = true
  try {
    const res = await queryPrometheus({ query: promQuery.value.trim() })
    if (res.code === 200) {
      queryResult.value = JSON.stringify(res.data, null, 2)
    }
  } catch (error: any) {
    queryResult.value = error.message || t('monitoringPage.queryFailed')
  } finally {
    queryLoading.value = false
  }
}

const startTimer = () => {
  stopTimer()
  if (!autoRefresh.value) return
  timer = setInterval(refreshAll, 30000)
}

const stopTimer = () => {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}

watch(autoRefresh, startTimer)

onMounted(async () => {
  await refreshAll()
  startTimer()
})

onUnmounted(stopTimer)
</script>

<style lang="scss" scoped>
.monitoring-page {
  padding: 24px;
  background: var(--el-bg-color-page);
  min-height: 100%;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}
.page-title {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
}
.page-subtitle {
  margin: 6px 0 0;
  color: var(--el-text-color-secondary);
  font-size: 14px;
}
.header-actions {
  display: flex;
  align-items: center;
}
.summary-row,
.section-row {
  margin-bottom: 16px;
}
.metric-card {
  margin-bottom: 12px;
}
.metric-label {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  margin-bottom: 8px;
}
.metric-value {
  font-size: 28px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.stat-block {
  padding: 12px;
  margin-bottom: 8px;
  background: var(--el-fill-color-lighter);
  border-radius: 4px;
}
.stat-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-bottom: 4px;
}
.stat-value {
  font-size: 20px;
  font-weight: 600;
}
.muted {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
.query-result {
  margin-top: 12px;
  max-height: 240px;
  overflow: auto;
  padding: 12px;
  background: var(--el-fill-color-lighter);
  border-radius: 4px;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-word;
}
.alerts-list {
  max-height: 420px;
  overflow: auto;
}
.alert-item {
  display: flex;
  gap: 10px;
  padding: 10px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}
.alert-title {
  font-weight: 600;
  margin-bottom: 4px;
}
.alert-desc,
.alert-time {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
