<template>
  <div class="dashboard-container">
    <!-- 面包屑导航 -->
    <div class="breadcrumb-section">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>首页</el-breadcrumb-item>
        <el-breadcrumb-item>集群监控</el-breadcrumb-item>
        <el-breadcrumb-item>Kubernetes</el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <!-- 标题与控制器 -->
    <div class="dashboard-header">
      <h2>Kubernetes集群监控仪表盘</h2>
      <div class="header-controls">
        <div class="time-range-selector">
          <el-date-picker
            v-model="timeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            :shortcuts="timeShortcuts"
            @change="handleTimeRangeChange"
          />
        </div>
        <div class="control-group">
          <el-select
            v-model="selectedNamespace"
            placeholder="全部命名空间"
            size="small"
            clearable
            class="namespace-select"
          >
            <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
          </el-select>
          <el-button
            type="primary"
            size="small"
            :icon="RefreshIcon"
            @click="refreshData"
            :loading="refreshing"
            class="refresh-btn"
          >
            刷新数据
          </el-button>
          <el-tooltip content="最后更新时间">
            <div class="last-update">
              <el-icon><Clock /></el-icon>
              <span>{{ lastUpdateTime }}</span>
            </div>
          </el-tooltip>
        </div>
      </div>
    </div>

    <!-- 资源概览卡片 --> 

    
    <!-- 健康状况指示灯 -->
    <div class="health-indicator">
      <div class="indicator-item" v-for="(item, index) in healthStatus" :key="index">
        <div class="indicator-dot" :style="{ backgroundColor: item.color }"></div>
        <div class="indicator-text">
          <span class="label">{{ item.label }}</span>
          <span class="value">{{ item.value }}</span>
        </div>
      </div>
    </div>


     

    <!-- 告警统计卡片 -->
    <div class="alert-summary-card">
      <div class="alert-item" v-for="(item, index) in alertSummary" :key="index">
        <div class="alert-count" :style="{ color: item.color }">
          {{ item.count }}
        </div>
        <div class="alert-label">
          {{ item.label }}
          <el-tag :type="item.trend > 0 ? 'danger' : 'success'" size="small">
            <el-icon v-if="item.trend > 0"><Top /></el-icon>
            <el-icon v-else><Bottom /></el-icon>
            {{ Math.abs(item.trend) }}%
          </el-tag>
        </div>
      </div>
    </div>



    <!-- 集群概览卡片 -->
    <div class="dashboard-card">
      <div class="card-header">
        <div class="card-title">
          <el-icon><DataAnalysis /></el-icon>
          <span>集群概览</span>
        </div>
      </div>
      <div class="card-body">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="6" v-for="(item, index) in overviewData" :key="index">
            <div class="overview-card-item">
              <div class="card-icon" :style="{ backgroundColor: item.color + '1a' }">
                <component :is="item.icon" :style="{ color: item.color }" />
              </div>
              <div class="card-content">
                <div class="card-title">{{ item.title }}</div>
                <div class="card-value">{{ item.value }}</div>
                <el-progress
                  :percentage="item.percent"
                  :color="item.color"
                  :stroke-width="8"
                  :show-text="false"
                />
                <div class="card-description">
                  <span>总{{ item.total || 'N/A' }}</span>
                  <span class="usage">使用率 {{ item.percent }}%</span>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- 资源使用率图表 -->
    <div class="dashboard-card">
      <div class="card-header">
        <div class="card-title">
          <el-icon><LoadingIcon /></el-icon>
          <span>资源使用率</span>
        </div>
      </div>
      <div class="card-body">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <div class="chart-container">
              <div class="chart-title">CPU使用情况</div>
              <div class="chart-wrapper">
                <v-chart :option="cpuUsageOption" autoresize />
              </div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="12">
            <div class="chart-container">
              <div class="chart-title">内存使用情况</div>
              <div class="chart-wrapper">
                <v-chart :option="memoryUsageOption" autoresize />
              </div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="12">
            <div class="chart-container">
              <div class="chart-title">存储使用情况</div>
              <div class="chart-wrapper">
                <v-chart :option="storageUsageOption" autoresize />
              </div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="12">
            <div class="chart-container">
              <div class="chart-title">网络流量</div>
              <div class="chart-wrapper">
                <v-chart :option="networkUsageOption" autoresize />
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <!-- 节点状态表格 -->
    <div class="dashboard-card">
      <div class="card-header">
        <div class="card-title">
          <el-icon><Loading /></el-icon>
          <span>节点状态</span>
        </div>
      </div>
      <div class="card-body">
        <el-table
          :data="nodes"
          stripe
          border
          style="width: 100%"
          v-loading="loadingNodes"
          :header-cell-style="{ background: '#f5f7fa', color: '#666' }"
          highlight-current-row
          @row-click="handleNodeClick"
        >
          <el-table-column prop="name" label="节点名称" width="180" sortable>
            <template #default="{ row }">
              <div class="node-name">
                <el-icon :color="row.status === 'Ready' ? '#67C23A' : '#F56C6C'">
                  <CircleCheck v-if="row.status === 'Ready'" />
                  <CircleClose v-else />
                </el-icon>
                <span>{{ row.name }}</span>
                <el-tag v-if="row.isNew" size="small" type="success" class="new-tag">
                  新节点
                </el-tag>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="role" label="角色" width="120" sortable>
            <template #default="{ row }">
              <el-tag :type="row.role === 'master' ? 'info' : 'warning'" effect="light" size="small">
                <el-icon v-if="row.role === 'master'"><DataBoard /></el-icon>
                {{ row.role === 'master' ? '控制节点' : '工作节点' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100" sortable>
            <template #default="{ row }">
              <el-tag :type="row.status === 'Ready' ? 'success' : 'danger'" effect="light" size="small">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="cpuUsage" label="CPU使用率" sortable>
            <template #default="{ row }">
              <div class="progress-container">
                <el-progress
                  :percentage="row.cpuUsage"
                  :stroke-width="16"
                  :color="getProgressColor(row.cpuUsage)"
                  :show-text="false"
                />
                <span class="progress-text">{{ row.cpuUsage }}%</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="memoryUsage" label="内存使用率" sortable>
            <template #default="{ row }">
              <div class="progress-container">
                <el-progress
                  :percentage="row.memoryUsage"
                  :stroke-width="16"
                  :color="getProgressColor(row.memoryUsage)"
                  :show-text="false"
                />
                <span class="progress-text">{{ row.memoryUsage }}%</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="pods" label="Pods" sortable>
            <template #default="{ row }">
              <div class="pod-count">
                <span class="running">{{ row.runningPods }}</span>
                <span class="separator">/</span>
                <span class="total">{{ row.totalPods }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100">
            <template #default>
              <el-button type="primary" size="small" link>详情</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- 集群事件日志 -->
    <div class="dashboard-card">
      <div class="card-header">
        <div class="card-title">
          <el-icon><Document /></el-icon>
          <span>集群事件</span>
        </div>
      </div>
      <div class="card-body">
        <el-tabs v-model="activeEventTab" class="event-tabs">
          <el-tab-pane label="最新事件" name="recent">
            <el-table
              :data="recentEvents"
              style="width: 100%"
              height="300"
              v-loading="loadingEvents"
            >
              <el-table-column prop="timestamp" label="时间" width="160" sortable>
                <template #default="{ row }">
                  {{ formatDate(row.timestamp) }}
                </template>
              </el-table-column>
              <el-table-column prop="type" label="类型" width="100" sortable>
                <template #default="{ row }">
                  <el-tag :type="row.type === 'Warning' ? 'warning' : 'success'" effect="light" size="small">
                    {{ row.type }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="object" label="对象" />
              <el-table-column prop="namespace" label="命名空间" width="120" />
              <el-table-column prop="reason" label="原因" width="150" />
              <el-table-column prop="message" label="消息" show-overflow-tooltip />
            </el-table>
          </el-tab-pane>
          <el-tab-pane label="事件统计" name="statistics">
            <div class="event-statistics-container">
              <div class="statistics-chart">
                <v-chart :option="eventStatisticsOption" autoresize />
              </div>
              <div class="statistics-table">
                <el-table
                  :data="eventStatistics"
                  border
                  style="width: 100%"
                  height="280"
                >
                  <el-table-column prop="type" label="事件类型" width="120" />
                  <el-table-column prop="count" label="数量" width="80" />
                  <el-table-column prop="percentage" label="百分比">
                    <template #default="{ row }">
                      <div class="percentage-bar">
                        <span class="percentage">{{ row.percentage }}%</span>
                        <div class="percentage-progress" :style="{ width: row.percentage + '%' }"></div>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="trend" label="趋势">
                    <template #default="{ row }">
                      <el-tag
                        :type="row.trend === 'up' ? 'danger' : 'success'"
                        size="small"
                        :icon="row.trend === 'up' ? ArrowUpBoldIcon : ArrowDownBoldIcon"
                      >
                        {{ row.trend === 'up' ? '上升' : '下降' }}
                      </el-tag>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>

    <!-- 资源水位警报 -->
    <div class="dashboard-card">
      <div class="card-header">
        <div class="card-title">
          <el-icon><Warning /></el-icon>
          <span>资源水位警报</span>
        </div>
      </div>
      <div class="card-body">
        <div class="resource-alerts">
          <div v-for="(item, index) in resourceAlerts" :key="index" class="alert-item">
            <div class="alert-level" :class="'level-' + item.level">
              <span>{{ item.levelText }}</span>
            </div>
            <div class="alert-content">
              <div class="alert-name">{{ item.name }}</div>
              <div class="alert-message">{{ item.message }}</div>
            </div>
            <div class="alert-time">{{ item.time }}</div>
            <el-button size="small" type="text">处理</el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { PieChart, BarChart, LineChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'
import VChart from 'vue-echarts'
import {
  Refresh as RefreshIcon,
  Clock as ClockIcon,
  DataAnalysis as DataAnalysisIcon,
  PieChart as PieChartIcon,
  Loading as LoadingIcon,
  Document as Document,
  ArrowUpBold as ArrowUpBoldIcon,
  ArrowDownBold as ArrowDownBoldIcon,
  DataBoard as DataBoardIcon,
  Collection as CollectionIcon,
  Box as BoxIcon,
  Connection as Connection,
  CircleCheck,
  CircleClose,
  Top,
  Bottom,
  Warning
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import ResourceSummaryCard from '@/components/ClusterSummary/ResourceSummaryCard.vue'

// 注册 ECharts 组件
use([
  CanvasRenderer,
  PieChart,
  BarChart,
  LineChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

// 状态数据
const selectedNamespace = ref('')
const refreshing = ref(false)
const loadingNodes = ref(false)
const loadingEvents = ref(false)
const activeEventTab = ref('recent')
const lastUpdateTime = ref(dayjs().format('YYYY-MM-DD HH:mm:ss'))
const timeRange = ref<[Date, Date]>([
  new Date(dayjs().subtract(6, 'hour').toDate()),
  new Date(dayjs().toDate())
])

// 时间快捷选项
const timeShortcuts = [
  {
    text: '最近6小时',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 6)
      return [start, end]
    }
  },
  {
    text: '最近12小时',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 12)
      return [start, end]
    }
  },
  {
    text: '最近24小时',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24)
      return [start, end]
    }
  },
  {
    text: '最近7天',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start, end]
    }
  }
]

// 模拟数据
const namespaces = ref(['default', 'kube-system', 'monitoring', 'logging', 'dev', 'prod'])

const healthStatus = ref([
  { label: '集群状态', value: '健康', color: '#67C23A' },
  { label: '节点在线', value: '4/5', color: '#E6A23C' },
  { label: 'Pod运行', value: '148/180', color: '#409EFF' },
  { label: '警报', value: '3', color: '#F56C6C' }
])

const alertSummary = ref([
  { label: '严重警报', count: 1, color: '#F56C6C', trend: 12 },
  { label: '警告', count: 3, color: '#E6A23C', trend: -5 },
  { label: '通知', count: 8, color: '#909399', trend: 3 }
])

const overviewData = ref([
  { title: '集群节点', value: '5', percent: 95, total: '', icon: DataBoardIcon, color: '#409EFF' },
  { title: '命名空间', value: '6', percent: 75, total: '8', icon: CollectionIcon, color: '#67C23A' },
  { title: '运行Pods', value: '148', percent: 82, total: '180', icon: BoxIcon, color: '#E6A23C' },
  { title: '服务数量', value: '23', percent: 65, total: '35', icon: Connection, color: '#F56C6C' }
])

const cpuUsageOption = ref({
  backgroundColor: 'transparent',
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' }, formatter: '{b}<br/>{a0}: {c0}%' },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: {
    type: 'category',
    data: ['Master-1', 'Node-1', 'Node-2', 'Node-3', 'Node-4'],
    axisLine: { lineStyle: { color: '#e0e0e0' } },
    axisLabel: { color: '#999' }
  },
  yAxis: {
    type: 'value',
    max: 100,
    min: 0,
    interval: 20,
    axisLine: { show: false },
    axisTick: { show: false },
    splitLine: { lineStyle: { color: 'rgba(0, 0, 0, 0.05)' } },
    axisLabel: { color: '#999' }
  },
  series: [{
    name: 'CPU使用率',
    type: 'bar',
    barWidth: '30%',
    data: [65, 45, 30, 85, 60],
    itemStyle: {
      color: (params: { data: number }) => getProgressColor(params.data),
      borderRadius: [4, 4, 0, 0]
    },
    label: { show: true, position: 'top', formatter: '{c}%', color: '#666' }
  }]
})

const memoryUsageOption = ref({
  backgroundColor: 'transparent',
  tooltip: { trigger: 'axis', axisPointer: { type: 'cross', crossStyle: { color: '#999' } } },
  legend: { data: ['已用内存', '缓存', '可用内存'], right: 10, top: 10, textStyle: { color: '#666' } },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: {
    type: 'category',
    data: ['Master-1', 'Node-1', 'Node-2', 'Node-3', 'Node-4'],
    axisLine: { lineStyle: { color: '#e0e0e0' } },
    axisLabel: { color: '#999' }
  },
  yAxis: {
    type: 'value',
    name: '内存 (GB)',
    min: 0,
    max: 24,
    interval: 4,
    axisLine: { lineStyle: { color: '#e0e0e0' } },
    axisLabel: { color: '#999' },
    splitLine: { lineStyle: { color: 'rgba(0, 0, 0, 0.05)' } }
  },
  series: [
    { name: '已用内存', type: 'bar', stack: 'total', data: [6, 10, 8, 15, 12], itemStyle: { color: '#F56C6C', borderRadius: [4, 4, 0, 0] }, emphasis: { focus: 'series' } },
    { name: '缓存', type: 'bar', stack: 'total', data: [3, 4, 2, 5, 3], itemStyle: { color: '#E6A23C', borderRadius: [4, 4, 0, 0] }, emphasis: { focus: 'series' } },
    { name: '可用内存', type: 'bar', stack: 'total', data: [15, 10, 14, 4, 9], itemStyle: { color: '#67C23A', borderRadius: [4, 4, 0, 0] }, emphasis: { focus: 'series' } }
  ]
})

const storageUsageOption = ref({
  backgroundColor: 'transparent',
  tooltip: { trigger: 'axis', axisPointer: { type: 'cross', label: { backgroundColor: '#6a7985' } } },
  legend: { data: ['已使用', '总量'], right: 10, top: 10, textStyle: { color: '#666' } },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: [{ type: 'category', boundaryGap: false, data: ['Master-1', 'Node-1', 'Node-2', 'Node-3', 'Node-4'], axisLine: { lineStyle: { color: '#e0e0e0' } }, axisLabel: { color: '#999' } }],
  yAxis: [{ type: 'value', name: '存储 (GB)', min: 0, max: 300, interval: 50, axisLine: { lineStyle: { color: '#e0e0e0' } }, axisLabel: { color: '#999' }, splitLine: { lineStyle: { color: 'rgba(0, 0, 0, 0.05)' } } }],
  series: [
    { name: '已使用', type: 'line', stack: '总量', areaStyle: { color: '#F56C6C', opacity: 0.8 }, emphasis: { focus: 'series' }, lineStyle: { width: 2, color: '#F56C6C' }, itemStyle: { color: '#F56C6C' }, data: [120, 80, 90, 150, 100] },
    { name: '总量', type: 'line', stack: '总量', areaStyle: { color: '#E4E7ED', opacity: 0.3 }, emphasis: { focus: 'series' }, lineStyle: { width: 2, color: '#909399' }, itemStyle: { color: '#909399' }, data: [200, 200, 200, 200, 200] }
  ]
})

const networkUsageOption = ref({
  backgroundColor: 'transparent',
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  legend: { data: ['入站', '出站'], right: 10, top: 10, textStyle: { color: '#666' } },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: ['08:00', '10:00', '12:00', '14:00', '16:00', '18:00', '20:00'], axisLine: { lineStyle: { color: '#e0e0e0' } }, axisLabel: { color: '#999' } },
  yAxis: { type: 'value', name: '网络流量 (Mbps)', axisLine: { lineStyle: { color: '#e0e0e0' } }, axisLabel: { color: '#999' }, splitLine: { lineStyle: { color: 'rgba(0, 0, 0, 0.05)' } } },
  series: [
    { name: '入站', type: 'bar', stack: '流量', barWidth: '40%', data: [120, 200, 150, 180, 210, 190, 230], itemStyle: { color: '#409EFF' } },
    { name: '出站', type: 'bar', stack: '流量', barWidth: '40%', data: [100, 170, 130, 150, 180, 160, 190], itemStyle: { color: '#67C23A' } }
  ]
})

const nodes = ref([
  { name: 'master-1', role: 'master', status: 'Ready', cpuUsage: 65, memoryUsage: 75, totalPods: 64, runningPods: 60, isNew: false },
  { name: 'node-1', role: 'worker', status: 'Ready', cpuUsage: 45, memoryUsage: 60, totalPods: 128, runningPods: 110, isNew: true },
  { name: 'node-2', role: 'worker', status: 'Ready', cpuUsage: 30, memoryUsage: 55, totalPods: 128, runningPods: 98, isNew: false },
  { name: 'node-3', role: 'worker', status: 'NotReady', cpuUsage: 85, memoryUsage: 90, totalPods: 128, runningPods: 128, isNew: false },
  { name: 'node-4', role: 'worker', status: 'Ready', cpuUsage: 60, memoryUsage: 65, totalPods: 128, runningPods: 105, isNew: false }
])

const recentEvents = ref([
  { timestamp: '2023-10-15T09:23:17Z', type: 'Warning', object: 'pod/nginx-deployment-75675f5897-58xj7', namespace: 'default', reason: 'FailedScheduling', message: '0/4 nodes available: 3 Insufficient cpu, 1 Insufficient memory.' },
  { timestamp: '2023-10-15T09:20:45Z', type: 'Normal', object: 'deployment/nginx-deployment', namespace: 'default', reason: 'ScalingReplicaSet', message: 'Scaled up replica set nginx-deployment-75675f5897 to 1' },
  { timestamp: '2023-10-15T09:18:32Z', type: 'Warning', object: 'node/node-3', namespace: '', reason: 'NodeNotReady', message: 'Node node-3 status is now: NodeNotReady' },
  { timestamp: '2023-10-15T09:15:12Z', type: 'Normal', object: 'service/redis-master', namespace: 'prod', reason: 'UpdatedLoadBalancer', message: 'Updated load balancer with new hosts' },
  { timestamp: '2023-10-15T09:12:58Z', type: 'Normal', object: 'pod/mysql-56f4d7f65d-f5z2n', namespace: 'prod', reason: 'Pulled', message: 'Successfully pulled image "mysql:5.7"' },
  { timestamp: '2023-10-15T09:10:31Z', type: 'Warning', object: 'pod/cron-job-1625868000-abcde', namespace: 'dev', reason: 'Failed', message: 'Job has reached the specified backoff limit' }
])

const eventStatistics = ref([
  { type: 'Normal', count: 42, percentage: 70, trend: 'up' },
  { type: 'Warning', count: 18, percentage: 30, trend: 'down' },
  { type: 'Critical', count: 5, percentage: 10, trend: 'up' },
  { type: 'Info', count: 10, percentage: 20, trend: 'up' },
  { type: 'Notice', count: 8, percentage: 15, trend: 'down' },
  { type: 'Error', count: 3, percentage: 5, trend: 'up' },
  { type: 'Debug', count: 2, percentage: 5, trend: 'down' },
  { type: 'Alert', count: 1, percentage: 2, trend: 'up' },
  { type: 'Emergency', count: 0, percentage: 0, trend: 'down' }
])

const eventStatisticsOption = ref({
  backgroundColor: 'transparent',
  tooltip: { trigger: 'item' },
  legend: { top: '5%', left: 'center', textStyle: { color: '#666' } },
  series: [{
    name: '事件统计',
    type: 'pie',
    radius: ['50%', '70%'],
    center: ['50%', '60%'],
    avoidLabelOverlap: false,
    itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
    label: { show: false, position: 'center' },
    emphasis: { label: { show: true, fontSize: '18', fontWeight: 'bold', color: '#333' } },
    labelLine: { show: false },
    data: [
      { value: 42, name: '正常事件 (70%)', itemStyle: { color: '#67C23A' } },
      { value: 5, name: '严重事件 (10%)', itemStyle: { color: '#F56C6C' } },
      { value: 10, name: '通知事件 (20%)', itemStyle: { color: '#909399' } },
      { value: 18, name: '警告事件 (30%)', itemStyle: { color: '#E6A23C' } }
    ]
  }]
})

const resourceAlerts = ref([
  { name: 'node-3 CPU过载', message: '节点node-3 CPU使用率达85%，超过阈值80%', level: 'critical', levelText: '严重', time: '5分钟前' },
  { name: 'node-3 内存不足', message: '节点node-3 内存使用率达90%，超过阈值85%', level: 'warning', levelText: '警告', time: '10分钟前' },
  { name: 'default命名空间Pod数量', message: 'default命名空间Pod数量接近上限(80/90)', level: 'notice', levelText: '注意', time: '30分钟前' }
])

// 计算方法
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const getProgressColor = (percentage: number) => {
  if (percentage > 80) return '#f56c6c'
  if (percentage > 60) return '#e6a23c'
  return '#67c23a'
}

// 操作方法
const refreshData = () => {
  refreshing.value = true
  loadingNodes.value = true
  loadingEvents.value = true
  setTimeout(() => {
    lastUpdateTime.value = dayjs().format('YYYY-MM-DD HH:mm:ss')
    refreshing.value = false
    loadingNodes.value = false
    loadingEvents.value = false
  }, 1000)
}

const handleTimeRangeChange = (val: [Date, Date] | null) => {
  if (val) {
    console.log('时间范围变更:', val)
    refreshData()
  }
}

const handleNodeClick = (row: any) => {
  console.log('点击节点:', row)
}

// 组件挂载时加载数据
onMounted(() => {
  refreshData()
})
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 24px;
  background: #f5f7fa;
  min-height: calc(100vh - 64px);
  font-family: 'Arial', sans-serif;

  .breadcrumb-section {
    margin-bottom: 20px;
    padding: 10px 0;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  }

  .dashboard-header {
    display: flex;
    flex-direction: column;
    margin-bottom: 24px;
    padding: 20px;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);

    h2 {
      margin: 0 0 16px 0;
      font-size: 24px;
      color: #212529;
      font-weight: 600;
    }

    .header-controls {
      display: flex;
      flex-wrap: wrap;
      justify-content: space-between;
      gap: 16px;

      .time-range-selector {
        flex: 1;
        min-width: 300px;
      }

      .control-group {
        display: flex;
        align-items: center;
        gap: 16px;

        .namespace-select {
          width: 180px;
        }

        .refresh-btn {
          transition: all 0.3s;
          .el-icon {
            transition: transform 0.3s ease;
          }
          &:hover .el-icon {
            transform: rotate(360deg);
          }
        }

        .last-update {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 14px;
          color: #6c757d;
          padding: 6px 12px;
          background: #f8f9fa;
          border-radius: 4px;
          box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
          .el-icon {
            font-size: 16px;
          }
        }
      }
    }
  }

  .health-indicator {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    margin-bottom: 24px;
    padding: 16px;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);

    .indicator-item {
      display: flex;
      align-items: center;
      padding: 10px 16px;
      border-radius: 6px;
      background: linear-gradient(135deg, #ffffff, #f8f9fa);
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
      min-width: 180px;
      transition: all 0.3s;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
      }

      .indicator-dot {
        width: 12px;
        height: 12px;
        border-radius: 50%;
        margin-right: 12px;
      }

      .indicator-text {
        display: flex;
        flex-direction: column;

        .label {
          font-size: 13px;
          color: #6c757d;
        }

        .value {
          font-size: 16px;
          font-weight: 500;
          color: #343a40;
        }
      }
    }
  }

  .alert-summary-card {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    margin-bottom: 24px;
    padding: 16px;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);

    .alert-item {
      flex: 1;
      min-width: 200px;
      padding: 16px;
      border-radius: 6px;
      background: linear-gradient(135deg, #ffffff, #f8f9fa);
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
      transition: all 0.3s;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
      }

      .alert-count {
        font-size: 28px;
        font-weight: bold;
      }

      .alert-label {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-size: 14px;
        color: #495057;
        margin-top: 8px;

        .el-tag {
          font-size: 12px;
          padding: 0 8px;
        }
      }
    }
  }

  .dashboard-card {
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    margin-bottom: 24px;
    transition: all 0.3s ease;

    &:hover {
      box-shadow: 0 6px 16px rgba(0, 0, 0, 0.1);
    }

    .card-header {
      padding: 16px 20px;
      border-bottom: 1px solid #e9ecef;
      background: linear-gradient(90deg, #f8f9fa, #ffffff);

      .card-title {
        font-size: 18px;
        font-weight: 600;
        color: #343a40;
        display: flex;
        align-items: center;

        .el-icon {
          margin-right: 10px;
          font-size: 20px;
          color: #007bff;
        }
      }
    }

    .card-body {
      padding: 20px;
    }
  }

  .overview-card-item {
    display: flex;
    height: 130px;
    padding: 16px;
    background: #fff;
    border-radius: 6px;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
    transition: all 0.3s;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }

    .card-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 50px;
      height: 50px;
      border-radius: 50%;
      margin-right: 16px;
      flex-shrink: 0;
      font-size: 24px;
    }

    .card-content {
      flex: 1;
      display: flex;
      flex-direction: column;

      .card-title {
        font-size: 14px;
        color: #6c757d;
        margin-bottom: 6px;
      }

      .card-value {
        font-size: 26px;
        font-weight: bold;
        color: #212529;
        margin-bottom: 12px;
      }

      .el-progress {
        margin: 0 0 8px;
        height: 8px;
        border-radius: 4px;
      }

      .card-description {
        display: flex;
        justify-content: space-between;
        font-size: 12px;
        color: #6c757d;

        .usage {
          color: #343a40;
          font-weight: 500;
        }
      }
    }
  }

  .chart-container {
    margin-bottom: 20px;
    padding: 16px;
    background: #fff;
    border-radius: 6px;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);

    .chart-title {
      font-size: 16px;
      font-weight: 500;
      color: #343a40;
      margin-bottom: 12px;
    }

    .chart-wrapper {
      width: 100%;
      height: 250px;
    }
  }

  .event-tabs {
    :deep(.el-tabs__header) {
      margin: 0 0 16px;
      border-bottom: 2px solid #e9ecef;
    }

    :deep(.el-tabs__item) {
      font-size: 16px;
      color: #6c757d;
      padding: 12px 20px;
      font-weight: 500;

      &.is-active {
        color: #007bff;
        border-bottom: 2px solid #007bff;
      }
    }
  }

  .event-statistics-container {
    display: flex;
    gap: 20px;

    .statistics-chart {
      flex: 1;
      height: 300px;
      background: #fff;
      border-radius: 6px;
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
      padding: 16px;
    }

    .statistics-table {
      width: 320px;

      .percentage-bar {
        position: relative;
        height: 20px;
        background: #f8f9fa;
        border-radius: 10px;
        overflow: hidden;

        .percentage {
          position: absolute;
          left: 8px;
          z-index: 1;
          font-size: 12px;
          line-height: 20px;
          color: #fff;
        }

        .percentage-progress {
          height: 100%;
          background: linear-gradient(45deg, #007bff, #0056b3);
          border-radius: 10px;
          transition: width 0.3s ease;
        }
      }
    }
  }

  .resource-alerts {
    .alert-item {
      display: flex;
      align-items: center;
      padding: 12px 0;
      border-bottom: 1px solid #e9ecef;

      &:last-child {
        border-bottom: none;
      }

      .alert-level {
        width: 60px;
        text-align: center;
        padding: 6px;
        border-radius: 4px;
        font-size: 12px;
        font-weight: bold;
        color: #fff;

        &.level-critical {
          background: linear-gradient(45deg, #f56c6c, #dc3545);
        }

        &.level-warning {
          background: linear-gradient(45deg, #e6a23c, #fd7e14);
        }

        &.level-notice {
          background: linear-gradient(45deg, #409eff, #007bff);
        }
      }

      .alert-content {
        flex: 1;
        padding: 0 16px;

        .alert-name {
          font-weight: 500;
          color: #343a40;
          margin-bottom: 4px;
        }

        .alert-message {
          font-size: 12px;
          color: #6c757d;
        }
      }

      .alert-time {
        width: 100px;
        font-size: 12px;
        color: #999;
        text-align: right;
      }

      .el-button {
        width: 60px;
        color: #007bff;
        transition: color 0.3s;

        &:hover {
          color: #0056b3;
        }
      }
    }
  }

  :deep(.el-table) {
    border-radius: 8px;
    overflow: hidden;

    th.el-table__cell {
      background: #f5f7fa;
      color: #343a40;
      font-weight: 600;
    }

    .el-table__row--striped {
      background-color: #f8f9fa !important;
    }

    .node-name {
      display: flex;
      align-items: center;
      gap: 8px;

      .el-icon {
        font-size: 16px;
      }

      .new-tag {
        margin-left: 5px;
        background: #28a745;
        color: #fff;
      }
    }

    .progress-container {
      display: flex;
      align-items: center;
      gap: 8px;

      .el-progress {
        flex: 1;
        .el-progress-bar__outer {
          background: #e9ecef;
          border-radius: 4px;
        }
      }

      .progress-text {
        width: 40px;
        font-size: 12px;
        color: #666;
      }
    }

    .pod-count {
      .running {
        color: #343a40;
        font-weight: 500;
      }

      .separator {
        margin: 0 4px;
        color: #999;
      }

      .total {
        color: #999;
      }
    }
  }

  @media (max-width: 768px) {
    .dashboard-header .header-controls {
      flex-direction: column;
      .time-range-selector,
      .control-group {
        width: 100%;
      }
    }

    .health-indicator,
    .alert-summary-card {
      flex-direction: column;
      .indicator-item,
      .alert-item {
        width: 100%;
      }
    }

    .dashboard-container {
  padding: 20px;
}
/* Add spacing if needed */
.el-row {
    margin-bottom: 20px;
}


    .event-statistics-container {
      flex-direction: column;
      .statistics-table {
        width: 100%;
      }
    }
  }
}
</style>