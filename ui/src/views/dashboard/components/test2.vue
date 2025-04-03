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
              :disabled="refreshing"
            />
          </div>
          <div class="control-group">
            <el-select
              v-model="selectedNamespace"
              placeholder="全部命名空间"
              size="small"
              clearable
              class="namespace-select"
              @change="handleNamespaceChange"
              :disabled="refreshing"
            >
              <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
            </el-select>
            <el-button
              type="primary"
              size="small"
              :icon="RefreshIcon"
              @click="fetchDashboardData"
              :loading="refreshing"
              class="refresh-btn"
            >
              刷新数据
            </el-button>
            <el-tooltip content="最后更新时间">
              <div class="last-update">
                <el-icon><ClockIcon /></el-icon>
                <span>{{ lastUpdateTime }}</span>
              </div>
            </el-tooltip>
          </div>
        </div>
      </div>
  
      <!-- 健康状况指示灯 -->
      <el-skeleton :rows="1" animated :loading="refreshing && !healthStatus.length">
        <template #default>
          <div class="health-indicator">
            <div class="indicator-item" v-for="(item, index) in healthStatus" :key="index">
              <div class="indicator-dot" :style="{ backgroundColor: item.color }"></div>
              <div class="indicator-text">
                <span class="label">{{ item.label }}</span>
                <span class="value">{{ item.value }}</span>
              </div>
            </div>
          </div>
        </template>
      </el-skeleton>
  
  
      <!-- 告警统计卡片 -->
       <el-skeleton :rows="1" animated :loading="refreshing && !alertSummary.length">
         <template #default>
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
        </template>
      </el-skeleton>
  
      <!-- 集群概览卡片 -->
      <div class="dashboard-card">
        <div class="card-header">
          <div class="card-title">
            <el-icon><DataAnalysisIcon /></el-icon>
            <span>集群概览</span>
          </div>
        </div>
        <div class="card-body">
           <el-skeleton :rows="3" animated :loading="refreshing && !overviewData.length">
             <template #default>
              <el-row :gutter="20">
                <el-col :xs="24" :sm="12" :md="6" v-for="(item, index) in overviewData" :key="index">
                  <div class="overview-card-item">
                    <div class="card-icon" :style="{ backgroundColor: item.color + '1a' }">
                       <!-- Method 1: Using dynamic component based on string name -->
                       <component :is="iconComponents[item.icon]" :style="{ color: item.color }" v-if="iconComponents[item.icon]"/>
                       <!-- Method 2: If you passed the component directly (less flexible with API) -->
                       <!-- <component :is="item.icon" :style="{ color: item.color }" /> -->
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
            </template>
          </el-skeleton>
        </div>
      </div>
  
      <!-- 资源使用率图表 -->
      <div class="dashboard-card">
        <div class="card-header">
          <div class="card-title">
            <el-icon><LoadingIcon /></el-icon> <!-- Consider a different icon like BarChart -->
            <span>资源使用率</span>
          </div>
        </div>
        <div class="card-body">
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
               <el-skeleton :rows="5" animated :loading="refreshing">
                 <template #default>
                  <div class="chart-container">
                    <div class="chart-title">CPU使用情况</div>
                    <div class="chart-wrapper">
                      <v-chart :option="cpuUsageOption" autoresize />
                    </div>
                  </div>
                </template>
              </el-skeleton>
            </el-col>
            <el-col :xs="24" :sm="12">
               <el-skeleton :rows="5" animated :loading="refreshing">
                 <template #default>
                  <div class="chart-container">
                    <div class="chart-title">内存使用情况</div>
                    <div class="chart-wrapper">
                      <v-chart :option="memoryUsageOption" autoresize />
                    </div>
                  </div>
                </template>
              </el-skeleton>
            </el-col>
            <el-col :xs="24" :sm="12">
               <el-skeleton :rows="5" animated :loading="refreshing">
                 <template #default>
                  <div class="chart-container">
                    <div class="chart-title">存储使用情况</div>
                    <div class="chart-wrapper">
                      <v-chart :option="storageUsageOption" autoresize />
                    </div>
                  </div>
                </template>
              </el-skeleton>
            </el-col>
            <el-col :xs="24" :sm="12">
               <el-skeleton :rows="5" animated :loading="refreshing">
                 <template #default>
                  <div class="chart-container">
                    <div class="chart-title">网络流量</div>
                    <div class="chart-wrapper">
                      <v-chart :option="networkUsageOption" autoresize />
                    </div>
                  </div>
                </template>
              </el-skeleton>
            </el-col>
          </el-row>
        </div>
      </div>
  
      <!-- 节点状态表格 -->
      <div class="dashboard-card">
        <div class="card-header">
          <div class="card-title">
            <!-- Using a more specific icon might be better, e.g., Grid or List -->
            <el-icon><DataBoardIcon /></el-icon>
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
            <!-- Columns definition remains mostly the same -->
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
              <template #default="{ row }"> <!-- Pass row to handler -->
                <el-button type="primary" size="small" link @click.stop="handleNodeDetail(row)">详情</el-button>
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
          <el-tabs v-model="activeEventTab" class="event-tabs" @tab-change="handleEventTabChange">
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
               <el-skeleton :rows="5" animated :loading="loadingEventStats">
                 <template #default>
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
                </template>
              </el-skeleton>
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
           <el-skeleton :rows="3" animated :loading="refreshing && !resourceAlerts.length">
             <template #default>
              <div class="resource-alerts">
                <div v-if="!resourceAlerts.length && !refreshing" class="no-alerts">
                  <el-empty description="暂无资源警报" />
                </div>
                <div v-for="(item, index) in resourceAlerts" :key="index" class="alert-item">
                  <div class="alert-level" :class="'level-' + item.level">
                    <span>{{ item.levelText }}</span>
                  </div>
                  <div class="alert-content">
                    <div class="alert-name">{{ item.name }}</div>
                    <div class="alert-message">{{ item.message }}</div>
                  </div>
                  <div class="alert-time">{{ item.time }}</div>
                  <el-button size="small" type="text" @click="handleAlertAction(item)">处理</el-button>
                </div>
              </div>
            </template>
          </el-skeleton>
        </div>
      </div>
  
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, onMounted, reactive, watch, shallowRef } from 'vue'
  import { use } from 'echarts/core'
  import { CanvasRenderer } from 'echarts/renderers'
  import { PieChart, BarChart, LineChart } from 'echarts/charts'
  import {
    TitleComponent, TooltipComponent, LegendComponent, GridComponent, MarkLineComponent
  } from 'echarts/components'
  import VChart from 'vue-echarts'
  import { ElMessage, ElMessageBox } from 'element-plus' // Import message/confirm if needed
  import {
    Refresh as RefreshIcon, Clock as ClockIcon, DataAnalysis as DataAnalysisIcon,
    // PieChart as PieChartIcon, // Not used directly
    Loading as LoadingIcon, Document, ArrowUpBold as ArrowUpBoldIcon, ArrowDownBold as ArrowDownBoldIcon,
    DataBoard as DataBoardIcon, Collection as CollectionIcon, Box as BoxIcon,
    Connection as ConnectionIcon, CircleCheck, CircleClose, Top, Bottom, Warning
    // Import other icons if needed, e.g., Grid, List, BarChart
  } from '@element-plus/icons-vue'
  import dayjs from 'dayjs'
  import { request } from '@/utils/service' // Adjust path as needed
  import type { AxiosRequestConfig } from 'axios'; // If using the provided request structure
  
  // Import Types
  import type {
      HealthStatusItem, AlertSummaryItem, OverviewDataItem, NodeData,
      EventItem, EventStatisticItem, ResourceAlertItem, ChartData, DashboardDataResponse
  } from '@/types/kubernetes'; // Adjust path as needed
  
  // Import Mock Data functions
  import {
      getMockNamespaces, getMockHealthStatus, getMockAlertSummary, getMockOverviewData,
      getMockCpuUsageChartData, getMockMemoryUsageChartData, getMockStorageUsageChartData,
      getMockNetworkUsageChartData, getMockNodes, getMockRecentEvents, getMockEventStatistics,
      getMockEventStatisticsChartData, getMockResourceAlerts
  } from '@/mockData/kubernetesMock'; // Adjust path as needed
  
  
  // ECharts setup
  use([
    CanvasRenderer, PieChart, BarChart, LineChart, TitleComponent,
    TooltipComponent, LegendComponent, GridComponent, MarkLineComponent
  ]);
  
  // --- State ---
  
  const selectedNamespace = ref('')
  const refreshing = ref(false) // Overall refresh state
  const loadingNodes = ref(false)
  const loadingEvents = ref(false)
  const loadingEventStats = ref(false); // Separate loading for stats tab
  const activeEventTab = ref('recent')
  const lastUpdateTime = ref('--') // Default value
  const timeRange = ref<[Date, Date]>([
    dayjs().subtract(6, 'hour').toDate(),
    dayjs().toDate()
  ]);
  
  // Map icon string names to actual components (if using string names from API/mock)
  const iconComponents = shallowRef({ // Use shallowRef for performance if icons don't change
      DataBoardIcon,
      CollectionIcon,
      BoxIcon,
      ConnectionIcon,
      // Add other icons used in overviewData here if needed
  });
  
  
  // --- Data Refs (Initialized with Mock Data initially) ---
  const namespaces = ref<string[]>(getMockNamespaces());
  const healthStatus = ref<HealthStatusItem[]>([]);
  const alertSummary = ref<AlertSummaryItem[]>([]);
  const overviewData = ref<OverviewDataItem[]>([]);
  const nodes = ref<NodeData[]>([]);
  const recentEvents = ref<EventItem[]>([]);
  const eventStatistics = ref<EventStatisticItem[]>([]);
  const resourceAlerts = ref<ResourceAlertItem[]>([]);
  
  // --- Chart Options (Structure defined, data will be populated) ---
  // Keep the reactive structure but initialize data arrays empty or with mock placeholders
  const cpuUsageOption = ref({
    backgroundColor: 'transparent',
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' }, formatter: '{b}<br/>{a0}: {c0}%' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: [] as string[], axisLine: { lineStyle: { color: '#e0e0e0' } }, axisLabel: { color: '#999' } },
    yAxis: { type: 'value', max: 100, min: 0, interval: 20, axisLine: { show: false }, axisTick: { show: false }, splitLine: { lineStyle: { color: 'rgba(0, 0, 0, 0.05)' } }, axisLabel: { color: '#999' } },
    series: [{ name: 'CPU使用率', type: 'bar', barWidth: '30%', data: [] as number[], itemStyle: { color: (params: { data: number }) => getProgressColor(params.data), borderRadius: [4, 4, 0, 0] }, label: { show: true, position: 'top', formatter: '{c}%', color: '#666' } }]
  });
  
  const memoryUsageOption = ref({
    backgroundColor: 'transparent',
    tooltip: { trigger: 'axis', axisPointer: { type: 'cross', crossStyle: { color: '#999' } } },
    legend: { data: ['已用内存', '缓存', '可用内存'], right: 10, top: 10, textStyle: { color: '#666' } },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: [] as string[], axisLine: { lineStyle: { color: '#e0e0e0' } }, axisLabel: { color: '#999' } },
    yAxis: { type: 'value', name: '内存 (GB)', min: 0, // Max/Interval can be dynamic
      axisLine: { lineStyle: { color: '#e0e0e0' } }, axisLabel: { color: '#999' }, splitLine: { lineStyle: { color: 'rgba(0, 0, 0, 0.05)' } } },
    series: [
      { name: '已用内存', type: 'bar', stack: 'total', data: [] as number[], itemStyle: { color: '#F56C6C', borderRadius: [4, 4, 0, 0] }, emphasis: { focus: 'series' } },
      { name: '缓存', type: 'bar', stack: 'total', data: [] as number[], itemStyle: { color: '#E6A23C' }, emphasis: { focus: 'series' } }, // No border radius for middle
      { name: '可用内存', type: 'bar', stack: 'total', data: [] as number[], itemStyle: { color: '#67C23A' }, emphasis: { focus: 'series' } } // No border radius for bottom in stacked bar usually
    ]
  });
  
  const storageUsageOption = ref({
      backgroundColor: 'transparent',
      tooltip: { trigger: 'axis', axisPointer: { type: 'cross', label: { backgroundColor: '#6a7985' } } },
      legend: { data: ['已使用', '总量'], right: 10, top: 10, textStyle: { color: '#666' } },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: [{ type: 'category', boundaryGap: false, data: [] as string[], axisLine: { lineStyle: { color: '#e0e0e0' } }, axisLabel: { color: '#999' } }],
      yAxis: [{ type: 'value', name: '存储 (GB)', min: 0, // Max/Interval can be dynamic
          axisLine: { lineStyle: { color: '#e0e0e0' } }, axisLabel: { color: '#999' }, splitLine: { lineStyle: { color: 'rgba(0, 0, 0, 0.05)' } } }],
      series: [
          { name: '已使用', type: 'line', stack: '总量', areaStyle: { color: '#F56C6C', opacity: 0.8 }, emphasis: { focus: 'series' }, lineStyle: { width: 2, color: '#F56C6C' }, itemStyle: { color: '#F56C6C' }, data: [] as number[] },
          { name: '总量', type: 'line', stack: '总量', areaStyle: { color: '#E4E7ED', opacity: 0.3 }, emphasis: { focus: 'series' }, lineStyle: { width: 2, color: '#909399' }, itemStyle: { color: '#909399' }, data: [] as number[] } // This might represent the *remaining* if stacked, or the total line itself
      ]
  });
  
  const networkUsageOption = ref({
      backgroundColor: 'transparent',
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      legend: { data: ['入站', '出站'], right: 10, top: 10, textStyle: { color: '#666' } },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: [] as string[], axisLine: { lineStyle: { color: '#e0e0e0' } }, axisLabel: { color: '#999' } },
      yAxis: { type: 'value', name: '网络流量 (Mbps)', axisLine: { lineStyle: { color: '#e0e0e0' } }, axisLabel: { color: '#999' }, splitLine: { lineStyle: { color: 'rgba(0, 0, 0, 0.05)' } } },
      series: [
          { name: '入站', type: 'bar', stack: '流量', barWidth: '40%', data: [] as number[], itemStyle: { color: '#409EFF' } },
          { name: '出站', type: 'bar', stack: '流量', barWidth: '40%', data: [] as number[], itemStyle: { color: '#67C23A' } }
      ]
  });
  
  const eventStatisticsOption = ref({
      backgroundColor: 'transparent',
      tooltip: { trigger: 'item' },
      legend: { top: '5%', left: 'center', textStyle: { color: '#666' } },
      series: [{
          name: '事件统计', type: 'pie', radius: ['50%', '70%'], center: ['50%', '60%'],
          avoidLabelOverlap: false, itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
          label: { show: false, position: 'center' },
          emphasis: { label: { show: true, fontSize: '18', fontWeight: 'bold', color: '#333' } },
          labelLine: { show: false },
          data: [] as any[] // Initialize empty, will be populated
      }]
  });
  
  // --- Time Shortcuts ---
  const timeShortcuts = [ /* ... as before ... */
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
  ];
  
  // --- API Endpoints (Placeholder - Replace with your actual endpoints) ---
  const API_ENDPOINTS = {
    // Option 1: Consolidated Endpoint (Recommended)
    // dashboardData: '/api/kubernetes/dashboard-data',
  
    // Option 2: Separate Endpoints
    namespaces: '/api/kubernetes/namespaces',
    health: '/api/kubernetes/health',
    alertSummary: '/api/kubernetes/alerts/summary',
    overview: '/api/kubernetes/overview',
    cpuUsage: '/api/kubernetes/metrics/cpu',
    memoryUsage: '/api/kubernetes/metrics/memory',
    storageUsage: '/api/kubernetes/metrics/storage',
    networkUsage: '/api/kubernetes/metrics/network',
    nodes: '/api/kubernetes/nodes',
    recentEvents: '/api/kubernetes/events/recent',
    eventStatistics: '/api/kubernetes/events/statistics',
    resourceAlerts: '/api/kubernetes/alerts/resource',
  };
  
  // --- Helper Functions ---
  const formatDate = (date: string | Date): string => {
    return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
  }
  
  const getProgressColor = (percentage: number): string => {
    if (percentage > 80) return '#f56c6c';
    if (percentage > 60) return '#e6a23c';
    return '#67c23a';
  }
  
  // --- API Fetching Logic ---
  const fetchDashboardData = async () => {
    if (refreshing.value) return; // Prevent concurrent refreshes
  
    refreshing.value = true;
    loadingNodes.value = true; // Assume nodes/events might take longer
    loadingEvents.value = true;
    // Reset data arrays before fetching new data, prevents stale data flashing
    healthStatus.value = [];
    alertSummary.value = [];
    overviewData.value = [];
    nodes.value = [];
    recentEvents.value = [];
    resourceAlerts.value = [];
    // Keep stats data until specifically requested or fetched
    // Don't reset chart options completely, just their data arrays
  
    const params: AxiosRequestConfig['params'] = {
      // Convert dates to ISO string or timestamp as required by your API
      startTime: timeRange.value[0]?.toISOString(),
      endTime: timeRange.value[1]?.toISOString(),
      namespace: selectedNamespace.value || undefined, // Send undefined or omit if empty/null
    };
  
    console.log("Fetching data with params:", params);
  
    // --- Option 1: Fetch all data from a single endpoint ---
    /*
    try {
        const response = await request<DashboardDataResponse>({
            url: API_ENDPOINTS.dashboardData,
            method: 'GET',
            params: params
        });
  
        namespaces.value = response.namespaces;
        healthStatus.value = response.healthStatus;
        alertSummary.value = response.alertSummary;
        overviewData.value = response.overviewData;
  
        // Update chart data
        cpuUsageOption.value.xAxis.data = response.cpuUsageData.labels;
        cpuUsageOption.value.series[0].data = response.cpuUsageData.values as number[];
  
        memoryUsageOption.value.xAxis.data = response.memoryUsageData.labels;
        memoryUsageOption.value.series[0].data = (response.memoryUsageData.values as number[][])[0]; // Used
        memoryUsageOption.value.series[1].data = (response.memoryUsageData.values as number[][])[1]; // Cached
        memoryUsageOption.value.series[2].data = (response.memoryUsageData.values as number[][])[2]; // Available
  
        storageUsageOption.value.xAxis.data = response.storageUsageData.labels;
        storageUsageOption.value.series[0].data = (response.storageUsageData.values as number[][])[0]; // Used
        storageUsageOption.value.series[1].data = (response.storageUsageData.values as number[][])[1]; // Total/Remaining
  
        networkUsageOption.value.xAxis.data = response.networkUsageData.labels;
        networkUsageOption.value.series[0].data = (response.networkUsageData.values as number[][])[0]; // Inbound
        networkUsageOption.value.series[1].data = (response.networkUsageData.values as number[][])[1]; // Outbound
  
        nodes.value = response.nodes;
        recentEvents.value = response.recentEvents;
        eventStatistics.value = response.eventStatistics;
        // Update event stats chart if needed (or fetch separately on tab click)
        eventStatisticsOption.value.series[0].data = getMockEventStatisticsChartData(response.eventStatistics); // Process data for chart
  
        resourceAlerts.value = response.resourceAlerts;
  
        lastUpdateTime.value = dayjs().format('YYYY-MM-DD HH:mm:ss');
  
    } catch (error) {
        console.error("Failed to fetch dashboard data:", error);
        ElMessage.error('获取仪表盘数据失败，将显示模拟数据。');
        // Fallback to all mock data
        namespaces.value = getMockNamespaces();
        healthStatus.value = getMockHealthStatus();
        alertSummary.value = getMockAlertSummary();
        overviewData.value = getMockOverviewData();
        // Fallback chart data
        const mockCpu = getMockCpuUsageChartData();
        cpuUsageOption.value.xAxis.data = mockCpu.labels;
        cpuUsageOption.value.series[0].data = mockCpu.values;
        // ... fallback for other charts ...
        nodes.value = getMockNodes();
        recentEvents.value = getMockRecentEvents();
        eventStatistics.value = getMockEventStatistics();
        eventStatisticsOption.value.series[0].data = getMockEventStatisticsChartData(getMockEventStatistics());
        resourceAlerts.value = getMockResourceAlerts();
        lastUpdateTime.value = `模拟数据 ${dayjs().format('HH:mm:ss')}`;
    } finally {
        refreshing.value = false;
        loadingNodes.value = false;
        loadingEvents.value = false;
    }
    */
  
    // --- Option 2: Fetch data from separate endpoints ---
    // Use Promise.allSettled to fetch concurrently and handle individual failures
    const results = await Promise.allSettled([
      request<string[]>({ url: API_ENDPOINTS.namespaces, method: 'GET' }), // Namespaces might not need time/ns params
      request<HealthStatusItem[]>({ url: API_ENDPOINTS.health, method: 'GET', params }),
      request<AlertSummaryItem[]>({ url: API_ENDPOINTS.alertSummary, method: 'GET', params }),
      request<OverviewDataItem[]>({ url: API_ENDPOINTS.overview, method: 'GET', params }),
      request<ChartData>({ url: API_ENDPOINTS.cpuUsage, method: 'GET', params }),
      request<ChartData>({ url: API_ENDPOINTS.memoryUsage, method: 'GET', params }),
      request<ChartData>({ url: API_ENDPOINTS.storageUsage, method: 'GET', params }),
      request<ChartData>({ url: API_ENDPOINTS.networkUsage, method: 'GET', params }),
      request<NodeData[]>({ url: API_ENDPOINTS.nodes, method: 'GET', params }), // Nodes might only need namespace filter
      request<EventItem[]>({ url: API_ENDPOINTS.recentEvents, method: 'GET', params }), // Events need time/ns filter
      // Event stats often fetched separately on tab click
      // request<EventStatisticItem[]>({ url: API_ENDPOINTS.eventStatistics, method: 'GET', params }),
      request<ResourceAlertItem[]>({ url: API_ENDPOINTS.resourceAlerts, method: 'GET', params }), // Alerts might only need ns filter
    ]);
  
    // Process results, falling back to mock data on failure
    if (results[0].status === 'fulfilled') {
      namespaces.value = results[0].value;
    } else {
      console.error("Failed to fetch namespaces:", results[0].reason);
      namespaces.value = getMockNamespaces(); // Keep mock namespaces if fetch fails
    }
  
    if (results[1].status === 'fulfilled') {
      healthStatus.value = results[1].value;
    } else {
      console.error("Failed to fetch health status:", results[1].reason);
      healthStatus.value = getMockHealthStatus();
    }
  
    if (results[2].status === 'fulfilled') {
      alertSummary.value = results[2].value;
    } else {
      console.error("Failed to fetch alert summary:", results[2].reason);
      alertSummary.value = getMockAlertSummary();
    }
  
      if (results[3].status === 'fulfilled') {
          overviewData.value = results[3].value;
      } else {
          console.error("Failed to fetch overview data:", results[3].reason);
          overviewData.value = getMockOverviewData();
      }
  
      // CPU Chart
      if (results[4].status === 'fulfilled') {
          const data = results[4].value;
          cpuUsageOption.value.xAxis.data = data.labels;
          cpuUsageOption.value.series[0].data = data.values as number[];
      } else {
          console.error("Failed to fetch CPU data:", results[4].reason);
          const mockData = getMockCpuUsageChartData();
          cpuUsageOption.value.xAxis.data = mockData.labels;
          cpuUsageOption.value.series[0].data = mockData.values;
      }
  
      // Memory Chart
      if (results[5].status === 'fulfilled') {
          const data = results[5].value;
          memoryUsageOption.value.xAxis.data = data.labels;
          // Assuming API returns [[used], [cached], [available]]
          if (Array.isArray(data.values) && data.values.length >= 3 && Array.isArray(data.values[0])) {
              memoryUsageOption.value.series[0].data = (data.values as number[][])[0];
              memoryUsageOption.value.series[1].data = (data.values as number[][])[1];
              memoryUsageOption.value.series[2].data = (data.values as number[][])[2];
          } else { // Fallback if structure is wrong
               const mockData = getMockMemoryUsageChartData();
               memoryUsageOption.value.xAxis.data = mockData.labels;
               memoryUsageOption.value.series[0].data = mockData.values[0];
               memoryUsageOption.value.series[1].data = mockData.values[1];
               memoryUsageOption.value.series[2].data = mockData.values[2];
          }
      } else {
          console.error("Failed to fetch Memory data:", results[5].reason);
          const mockData = getMockMemoryUsageChartData();
          memoryUsageOption.value.xAxis.data = mockData.labels;
          memoryUsageOption.value.series[0].data = mockData.values[0];
          memoryUsageOption.value.series[1].data = mockData.values[1];
          memoryUsageOption.value.series[2].data = mockData.values[2];
      }
  
      // Storage Chart
      if (results[6].status === 'fulfilled') {
          const data = results[6].value;
          storageUsageOption.value.xAxis.data = data.labels;
           // Assuming API returns [[used], [total]]
          if (Array.isArray(data.values) && data.values.length >= 2 && Array.isArray(data.values[0])) {
              storageUsageOption.value.series[0].data = (data.values as number[][])[0]; // Used
              storageUsageOption.value.series[1].data = (data.values as number[][])[1]; // Total
          } else { // Fallback
               const mockData = getMockStorageUsageChartData();
               storageUsageOption.value.xAxis.data = mockData.labels;
               storageUsageOption.value.series[0].data = mockData.values[0];
               storageUsageOption.value.series[1].data = mockData.values[1];
          }
      } else {
          console.error("Failed to fetch Storage data:", results[6].reason);
          const mockData = getMockStorageUsageChartData();
          storageUsageOption.value.xAxis.data = mockData.labels;
          storageUsageOption.value.series[0].data = mockData.values[0];
          storageUsageOption.value.series[1].data = mockData.values[1];
      }
  
      // Network Chart
      if (results[7].status === 'fulfilled') {
          const data = results[7].value;
          networkUsageOption.value.xAxis.data = data.labels;
          // Assuming API returns [[in], [out]]
          if (Array.isArray(data.values) && data.values.length >= 2 && Array.isArray(data.values[0])) {
              networkUsageOption.value.series[0].data = (data.values as number[][])[0]; // In
              networkUsageOption.value.series[1].data = (data.values as number[][])[1]; // Out
          } else { // Fallback
              const mockData = getMockNetworkUsageChartData();
              networkUsageOption.value.xAxis.data = mockData.labels;
              networkUsageOption.value.series[0].data = mockData.values[0];
              networkUsageOption.value.series[1].data = mockData.values[1];
          }
      } else {
          console.error("Failed to fetch Network data:", results[7].reason);
          const mockData = getMockNetworkUsageChartData();
          networkUsageOption.value.xAxis.data = mockData.labels;
          networkUsageOption.value.series[0].data = mockData.values[0];
          networkUsageOption.value.series[1].data = mockData.values[1];
      }
  
  
      // Nodes
      if (results[8].status === 'fulfilled') {
          nodes.value = results[8].value;
      } else {
          console.error("Failed to fetch nodes:", results[8].reason);
          nodes.value = getMockNodes();
      }
      loadingNodes.value = false; // Set node loading false here
  
      // Recent Events
      if (results[9].status === 'fulfilled') {
          recentEvents.value = results[9].value;
      } else {
          console.error("Failed to fetch recent events:", results[9].reason);
          recentEvents.value = getMockRecentEvents();
      }
       loadingEvents.value = false; // Set event loading false here
  
  
      // Resource Alerts
      if (results[10].status === 'fulfilled') {
          resourceAlerts.value = results[10].value;
      } else {
          console.error("Failed to fetch resource alerts:", results[10].reason);
          resourceAlerts.value = getMockResourceAlerts();
      }
  
      // Check if any request failed to decide the update time message
      const hasFailedRequest = results.some(r => r.status === 'rejected');
      if (hasFailedRequest) {
           ElMessage.warning('部分数据获取失败，已使用模拟数据替代。');
           lastUpdateTime.value = `部分模拟数据 ${dayjs().format('HH:mm:ss')}`;
      } else {
           lastUpdateTime.value = dayjs().format('YYYY-MM-DD HH:mm:ss');
      }
  
      refreshing.value = false; // Set overall refreshing false at the end
  
      // Fetch event statistics if that tab is active (or on initial load maybe)
      if (activeEventTab.value === 'statistics' && !loadingEventStats.value) {
          fetchEventStatistics();
      }
  };
  
  // Separate function to fetch event statistics
  const fetchEventStatistics = async () => {
      if (loadingEventStats.value) return;
      loadingEventStats.value = true;
      eventStatistics.value = []; // Clear previous stats
  
      const params: AxiosRequestConfig['params'] = {
          startTime: timeRange.value[0]?.toISOString(),
          endTime: timeRange.value[1]?.toISOString(),
          namespace: selectedNamespace.value || undefined,
      };
  
      try {
          const stats = await request<EventStatisticItem[]>({
              url: API_ENDPOINTS.eventStatistics,
              method: 'GET',
              params
          });
          eventStatistics.value = stats;
          // Process data for the pie chart
          eventStatisticsOption.value.series[0].data = getMockEventStatisticsChartData(stats); // Replace with actual processing if needed
      } catch (error) {
          console.error("Failed to fetch event statistics:", error);
          const mockStats = getMockEventStatistics();
          eventStatistics.value = mockStats;
          // Use mock data for chart as well
          eventStatisticsOption.value.series[0].data = getMockEventStatisticsChartData(mockStats);
      } finally {
          loadingEventStats.value = false;
      }
  };
  
  
  // --- Event Handlers ---
  const handleTimeRangeChange = (val: [Date, Date] | null) => {
    if (val) {
      console.log('时间范围变更:', val);
      fetchDashboardData(); // Fetch data on time change
    }
  }
  
  const handleNamespaceChange = (ns: string | undefined) => {
      console.log('命名空间变更:', ns);
      fetchDashboardData(); // Fetch data on namespace change
  }
  
  const handleNodeClick = (row: NodeData) => {
    console.log('点击节点:', row);
    // Navigate to node detail page or show modal
    // Example: router.push(`/kubernetes/nodes/${row.name}`)
    ElMessage.info(`查看节点 ${row.name} 详情 (模拟)`);
  }
  
  const handleNodeDetail = (row: NodeData) => {
      console.log('点击节点详情按钮:', row);
      ElMessageBox.alert(`显示节点 ${row.name} 的详细信息面板。`, '节点详情 (模拟)', {
          confirmButtonText: '确定',
      });
  }
  
  const handleAlertAction = (alertItem: ResourceAlertItem) => {
      console.log('处理警报:', alertItem);
       ElMessageBox.confirm(`确定要处理警报 "${alertItem.name}" 吗？`, '确认处理', {
          confirmButtonText: '确认',
          cancelButtonText: '取消',
          type: 'warning',
      }).then(() => {
          ElMessage.success(`警报 "${alertItem.name}" 已标记为处理中 (模拟)`);
          // Optionally remove the alert from the list or call an API to acknowledge it
          // resourceAlerts.value = resourceAlerts.value.filter(a => a !== alertItem);
      }).catch(() => {
          ElMessage.info('取消处理');
      });
  }
  
  const handleEventTabChange = (tabName: string | number) => {
      console.log('切换事件Tab:', tabName);
      if (tabName === 'statistics' && !eventStatistics.value.length && !loadingEventStats.value) {
          // Fetch stats only when switching to the tab if they haven't been loaded yet
          fetchEventStatistics();
      }
  }
  
  // --- Lifecycle Hooks ---
  onMounted(() => {
    fetchDashboardData(); // Initial data load
  });
  
  // Optional: Watch for namespace changes if the select doesn't trigger refresh automatically
  // watch(selectedNamespace, () => {
  //   fetchDashboardData();
  // });
  
  </script>
  
  <style lang="scss" scoped>
  // --- Styles (Keep your existing styles) ---
  .dashboard-container {
    padding: 24px;
    background: #f5f7fa;
    min-height: calc(100vh - 64px); /* Adjust based on your layout's header height */
    font-family: 'Arial', sans-serif;
  
    .breadcrumb-section {
      margin-bottom: 20px;
      padding: 10px 15px; // Slightly more padding
      background: #fff;
      border-radius: 8px;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
      .el-breadcrumb {
          font-size: 14px;
      }
    }
  
    .dashboard-header {
      display: flex;
      flex-direction: column; // Stack title and controls vertically
      // align-items: flex-start; // Align items to the start
      gap: 16px; // Space between title and controls row
      margin-bottom: 24px;
      padding: 20px;
      background: #fff;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  
      h2 {
        margin: 0; // Remove default margins
        font-size: 22px; // Slightly smaller H2
        color: #212529;
        font-weight: 600;
        line-height: 1.2;
      }
  
      .header-controls {
        display: flex;
        flex-wrap: wrap; // Allow wrapping on smaller screens
        justify-content: space-between;
        align-items: center; // Align items vertically
        gap: 16px; // Space between control groups
        width: 100%; // Take full width
  
        .time-range-selector {
          // flex-grow: 1; // Allow it to grow
          min-width: 320px; // Minimum width before wrapping
           max-width: 400px; // Max width
          .el-date-editor {
              width: 100%; // Make picker take full width of its container
          }
        }
  
        .control-group {
          display: flex;
          align-items: center;
          gap: 12px; // Slightly less gap
  
          .namespace-select {
            width: 180px;
          }
  
          .refresh-btn {
            transition: all 0.3s;
            .el-icon {
              transition: transform 0.5s ease; // Smoother rotation
            }
            &:hover .el-icon {
              transform: rotate(360deg);
            }
             &.is-loading .el-icon {
                 animation: rotating 2s linear infinite; // Ensure loading icon spins
             }
          }
  
          .last-update {
            display: flex;
            align-items: center;
            gap: 6px;
            font-size: 13px; // Slightly smaller
            color: #6c757d;
            padding: 6px 10px;
            background: #f8f9fa;
            border: 1px solid #e9ecef; // Subtle border
            border-radius: 4px;
            // box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05); // Shadow might be too much here
            .el-icon {
              font-size: 15px; // Slightly smaller icon
              color: #495057; // Darker icon color
            }
          }
        }
      }
    }
  
      /* Add styles for el-skeleton */
      .el-skeleton {
          margin-bottom: 24px; /* Match card margin */
          padding: 16px;
          background: #fff;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0,0,0,0.05);
      }
  
      .no-alerts {
          padding: 20px;
          text-align: center;
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
        background: #f8f9fa; // Lighter background
        border: 1px solid #e9ecef; // Subtle border
        // background: linear-gradient(135deg, #ffffff, #f8f9fa); // Gradient can be nice too
        // box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04); // Lighter shadow
        min-width: 160px; // Slightly smaller min-width
        transition: all 0.2s ease-in-out; // Faster transition
  
        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 3px 6px rgba(0, 0, 0, 0.08); // Slightly stronger hover shadow
          border-color: #dee2e6;
        }
  
        .indicator-dot {
          width: 10px; // Smaller dot
          height: 10px;
          border-radius: 50%;
          margin-right: 10px; // Less margin
          flex-shrink: 0; // Prevent dot shrinking
        }
  
        .indicator-text {
          display: flex;
          flex-direction: column;
          line-height: 1.4; // Improve line spacing
  
          .label {
            font-size: 12px; // Smaller label
            color: #6c757d;
            margin-bottom: 2px; // Space between label and value
          }
  
          .value {
            font-size: 15px; // Slightly smaller value
            font-weight: 600; // Bolder value
            color: #343a40;
          }
        }
      }
    }
  
    .alert-summary-card {
      display: grid; // Use grid for better alignment
      grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); // Responsive columns
      gap: 16px;
      margin-bottom: 24px;
      padding: 16px;
      background: #fff;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  
      .alert-item {
        // flex: 1; // Not needed with grid
        // min-width: 180px; // Set in grid-template-columns
        padding: 16px;
        border-radius: 6px;
        background: #f8f9fa;
        border: 1px solid #e9ecef;
        // background: linear-gradient(135deg, #ffffff, #f8f9fa);
        // box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
        transition: all 0.2s ease-in-out;
  
        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 3px 6px rgba(0, 0, 0, 0.08);
           border-color: #dee2e6;
        }
  
        .alert-count {
          font-size: 26px; // Slightly smaller count
          font-weight: 700; // Bolder
          line-height: 1.1;
          margin-bottom: 8px; // Space below count
        }
  
        .alert-label {
          display: flex;
          justify-content: space-between;
          align-items: center;
          font-size: 13px; // Smaller label
          color: #495057;
          // margin-top: 8px; // Removed, using margin-bottom on count
  
          .el-tag {
            font-size: 11px; // Smaller tag
            height: 20px; // Explicit height
            padding: 0 6px; // Adjust padding
            .el-icon {
                font-size: 10px; // Smaller icon in tag
            }
          }
        }
      }
    }
  
    .dashboard-card {
      background: #fff;
      border-radius: 8px;
      border: 1px solid #e9ecef; // Add subtle border to cards
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03); // Lighter shadow
      margin-bottom: 24px;
      transition: box-shadow 0.3s ease; // Only transition shadow
  
      &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08); // Softer hover shadow
      }
  
      .card-header {
        padding: 14px 20px; // Adjust padding
        border-bottom: 1px solid #e9ecef;
        background: #fcfdff; // Very light background for header
        // background: linear-gradient(90deg, #f8f9fa, #ffffff); // Keep gradient if preferred
  
        .card-title {
          font-size: 17px; // Adjust size
          font-weight: 600;
          color: #343a40;
          display: flex;
          align-items: center;
          gap: 8px; // Space between icon and text
  
          .el-icon {
            // margin-right: 10px; // Use gap instead
            font-size: 18px; // Adjust size
            color: #4a5568; // More subtle icon color
            // color: #007bff; // Keep primary color if preferred
          }
        }
      }
  
      .card-body {
        padding: 20px;
      }
    }
  
      .overview-card-item {
          display: flex;
          align-items: center; // Vertically align icon and content
          height: 120px; // Slightly smaller height
          padding: 16px;
          background: #fff; // Keep white or use f8f9fa
          border-radius: 6px;
          border: 1px solid #eef2f7; // Very light border
          // box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
          transition: all 0.2s ease-in-out;
          margin-bottom: 20px; // Add margin between items on smaller screens if needed
  
          &:hover {
              transform: translateY(-3px); // More lift
              box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08);
              border-color: #e5e9f2;
          }
  
          .card-icon {
              display: flex;
              align-items: center;
              justify-content: center;
              width: 48px; // Smaller icon circle
              height: 48px;
              border-radius: 50%;
              margin-right: 16px;
              flex-shrink: 0;
              font-size: 22px; // Smaller icon size
              background-color: var(--icon-bg-color, #e0f2fe); // Use CSS variable
              color: var(--icon-color, #0ea5e9); // Use CSS variable
  
              // Example setting variables based on item color (can be done inline style too)
              &[style*="--icon-color:#409EFF"] { --icon-bg-color: #e0f2fe; --icon-color: #0ea5e9; }
              &[style*="--icon-color:#67C23A"] { --icon-bg-color: #dcfce7; --icon-color: #22c55e; }
              &[style*="--icon-color:#E6A23C"] { --icon-bg-color: #fef9c3; --icon-color: #eab308; }
              &[style*="--icon-color:#F56C6C"] { --icon-bg-color: #fee2e2; --icon-color: #ef4444; }
  
               // Use the inline style from the template for simplicity if preferred
              // :style="{ backgroundColor: item.color + '1a' }"
               .el-icon {
                  // color: var(--icon-color); // Set color via CSS variable if needed
               }
          }
  
          .card-content {
              flex: 1;
              display: flex;
              flex-direction: column;
              overflow: hidden; // Prevent text overflow issues
  
              .card-title {
                  font-size: 13px; // Smaller title
                  color: #6c757d;
                  margin-bottom: 4px; // Less space
                  white-space: nowrap;
                  overflow: hidden;
                  text-overflow: ellipsis;
              }
  
              .card-value {
                  font-size: 24px; // Adjust value size
                  font-weight: 600; // Less bold
                  color: #212529;
                  margin-bottom: 10px; // Adjust space
                  line-height: 1.2;
              }
  
              .el-progress {
                  margin: 0 0 6px; // Adjust space
                  height: 6px; // Thinner progress bar
                  border-radius: 3px;
                  .el-progress-bar__outer {
                      background-color: #e9ecef; // Lighter background for progress
                  }
              }
  
              .card-description {
                  display: flex;
                  justify-content: space-between;
                  font-size: 11px; // Smaller description
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
      background: #fff; // Keep white or use a very light grey
      border-radius: 6px;
      border: 1px solid #eef2f7; // Consistent light border
      // box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  
      .chart-title {
        font-size: 15px; // Adjust size
        font-weight: 600; // Bolder
        color: #343a40;
        margin-bottom: 16px; // More space below title
      }
  
      .chart-wrapper {
        width: 100%;
        height: 250px; // Keep height or adjust as needed
      }
    }
  
      .event-tabs {
          :deep(.el-tabs__header) {
              margin: 0 0 20px; // More space below tabs nav
              border-bottom: 1px solid #e9ecef; // Lighter border
          }
  
          :deep(.el-tabs__nav-wrap::after) {
               height: 1px; // Make the bottom line thinner
          }
  
  
          :deep(.el-tabs__item) {
              font-size: 15px; // Adjust size
              color: #6c757d;
              padding: 0 20px; // Adjust padding
               height: 45px; // Taller tabs
               line-height: 45px;
              font-weight: 500;
  
              &.is-active {
                  color: #007bff; // Keep primary color
                   font-weight: 600;
                  // border-bottom: 2px solid #007bff; // Handled by el-tabs active-bar
              }
              &:hover {
                  color: #343a40;
              }
          }
           :deep(.el-tabs__active-bar) {
              height: 3px; // Thicker active bar
              background-color: #007bff;
           }
      }
  
  
      .event-statistics-container {
          display: flex;
          flex-wrap: wrap; // Allow wrap on small screens
          gap: 20px;
  
          .statistics-chart {
              flex: 1; // Take available space
              min-width: 250px; // Minimum width for chart
              height: 300px;
              background: #fff;
              border-radius: 6px;
               border: 1px solid #eef2f7;
              // box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
              padding: 16px;
          }
  
          .statistics-table {
              flex: 1; // Take available space
              min-width: 300px; // Minimum width for table
  
              .percentage-bar {
                  position: relative;
                  height: 18px; // Slightly smaller
                  background: #f1f3f5; // Lighter grey background
                  border-radius: 9px;
                  overflow: hidden;
  
                  .percentage {
                      position: absolute;
                      left: 8px;
                      top: 0;
                      z-index: 1;
                      font-size: 10px; // Smaller percentage text
                      line-height: 18px;
                      color: #fff;
                      font-weight: 500;
                       text-shadow: 1px 1px 1px rgba(0,0,0,0.2); // Add shadow for readability
                  }
  
                  .percentage-progress {
                      height: 100%;
                      background: linear-gradient(90deg, #3b82f6, #2563eb); // Nicer blue gradient
                      // background: linear-gradient(45deg, #007bff, #0056b3); // Original gradient
                      border-radius: 9px;
                      transition: width 0.4s ease; // Smoother transition
                  }
              }
               .el-table { // Style table within stats
                  border-radius: 6px;
                  border: 1px solid #e9ecef;
                   th.el-table__cell {
                      background: #f8f9fa;
                      color: #495057;
                      font-weight: 600;
                      font-size: 12px;
                   }
                   td.el-table__cell {
                       font-size: 12px;
                       padding: 8px 0; // Adjust cell padding
                   }
               }
          }
      }
  
      .resource-alerts {
          .alert-item {
              display: flex;
              align-items: center;
              padding: 12px 0;
              border-bottom: 1px solid #f1f3f5; // Lighter separator
               transition: background-color 0.2s;
  
               &:hover {
                  background-color: #f8f9fa; // Subtle hover background
               }
  
              &:last-child {
                  border-bottom: none;
              }
  
              .alert-level {
                  width: 55px; // Slightly narrower
                  text-align: center;
                  padding: 4px 6px; // Adjust padding
                  border-radius: 4px;
                  font-size: 11px; // Smaller text
                  font-weight: 600; // Bolder
                  color: #fff;
                  flex-shrink: 0;
                  margin-right: 12px; // Space after level indicator
  
                  &.level-critical { background: linear-gradient(135deg, #f56c6c, #dc3545); }
                  &.level-warning { background: linear-gradient(135deg, #e6a23c, #fd7e14); }
                  &.level-notice { background: linear-gradient(135deg, #409eff, #0d6efd); } // Adjusted blue
                  // Add styles for other potential levels
              }
  
              .alert-content {
                  flex: 1;
                  padding: 0 12px 0 0; // Space before time/button
                  overflow: hidden; // Prevent content overflow
  
                  .alert-name {
                      font-weight: 600; // Bolder name
                      color: #343a40;
                      margin-bottom: 3px; // Less space
                      font-size: 14px; // Slightly larger name
                       white-space: nowrap;
                       overflow: hidden;
                       text-overflow: ellipsis;
                  }
  
                  .alert-message {
                      font-size: 12px;
                      color: #6c757d;
                      line-height: 1.4;
                       white-space: nowrap;
                       overflow: hidden;
                       text-overflow: ellipsis;
                  }
              }
  
              .alert-time {
                  width: 80px; // Less width for time
                  font-size: 11px; // Smaller time text
                  color: #adb5bd; // Lighter time color
                  text-align: right;
                  flex-shrink: 0;
                  margin-right: 10px; // Space before button
              }
  
              .el-button {
                  width: auto; // Let button size naturally
                  padding: 4px 8px; // Adjust button padding
                   height: 24px; // Explicit height
                  font-size: 12px;
                  color: #007bff;
                  transition: color 0.2s, background-color 0.2s;
                  flex-shrink: 0; // Prevent button shrinking
  
                  &:hover {
                      color: #0056b3;
                      background-color: rgba(0, 123, 255, 0.05); // Slight background on hover
                  }
              }
          }
      }
  
  
      :deep(.el-table) {
          border-radius: 8px;
          overflow: hidden; // Keep overflow hidden for radius
          border: 1px solid #e9ecef; // Ensure table border matches card border
  
          th.el-table__cell {
              background: #f8f9fa; // Lighter header background
              color: #495057; // Darker grey header text
              font-weight: 600;
              font-size: 13px; // Adjust header font size
               padding: 10px 0; // Adjust header padding
          }
  
           /* Remove default element plus striped style if you prefer no stripes */
          /* .el-table__row--striped {
              background-color: transparent !important;
          } */
           /* Or customize stripe color */
           .el-table__row--striped td.el-table__cell{
               background-color: #fcfdff !important; // Very subtle stripe
           }
  
           td.el-table__cell {
               font-size: 13px; // Consistent row font size
               padding: 10px 0; // Adjust cell padding
                border-color: #f1f3f5; // Lighter cell borders
           }
            .el-table__row {
               transition: background-color 0.2s;
                &:hover > td.el-table__cell { // Row hover effect
                     background-color: #f5f7fa !important;
                 }
            }
  
  
          .node-name {
              display: flex;
              align-items: center;
              gap: 6px; // Less gap
  
              .el-icon {
                  font-size: 15px; // Slightly smaller icon
              }
              span {
                  white-space: nowrap;
                  overflow: hidden;
                  text-overflow: ellipsis;
              }
  
              .new-tag {
                  margin-left: 5px;
                  background-color: #198754; // Bootstrap success green
                  border-color: #198754;
                  color: #fff;
                  height: 18px; // Smaller tag
                   padding: 0 5px;
                   line-height: 16px;
                   font-size: 10px;
              }
          }
  
          .progress-container {
              display: flex;
              align-items: center;
              gap: 8px;
  
              .el-progress {
                  flex: 1;
                   height: 12px; // Slightly thicker progress bar in table
                  .el-progress-bar__outer {
                      background: #e9ecef;
                      border-radius: 6px;
                  }
                   .el-progress-bar__inner {
                      border-radius: 6px;
                   }
              }
  
              .progress-text {
                  width: 35px; // Less width for text
                  font-size: 11px;
                  color: #666;
                  text-align: right;
              }
          }
  
          .pod-count {
            font-size: 12px;
            .running {
              color: #212529;
              font-weight: 600;
            }
            .separator { margin: 0 3px; color: #adb5bd; }
            .total { color: #6c757d; }
          }
  
          // Style for action buttons in table
           .el-table__cell .el-button--small.is-link {
              font-size: 12px;
              padding: 2px 4px; // Smaller padding for link buttons
           }
      }
  
    // Responsive adjustments
    @media (max-width: 1200px) {
       .overview-card-item { height: auto; min-height: 120px;} // Allow wrap inside
    }
  
    @media (max-width: 992px) {
         .dashboard-header .header-controls {
             flex-direction: column;
             align-items: stretch; // Make controls take full width
             .time-range-selector,
             .control-group {
                 width: 100%;
                 justify-content: flex-start; // Align controls left
             }
             .control-group {
                  justify-content: space-between; // Space out items in the group again
             }
         }
         .alert-summary-card {
             grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); // Adjust grid for smaller screens
         }
    }
  
  
    @media (max-width: 768px) {
      .dashboard-container { padding: 16px; }
      .dashboard-header h2 { font-size: 20px; }
      .health-indicator, .alert-summary-card {
        // Keep grid/flex as is, they handle wrapping
      }
       .overview-card-item {
           height: auto;
           flex-direction: column; // Stack icon and content vertically
           align-items: flex-start;
           .card-icon { margin-bottom: 12px; }
       }
      .event-statistics-container {
          flex-direction: column; // Stack chart and table
          .statistics-chart, .statistics-table {
              width: 100%; // Make full width
          }
           .statistics-chart { height: 250px; } // Reduce chart height
      }
       .resource-alerts .alert-item {
          flex-wrap: wrap; // Allow wrapping inside alert item
           .alert-content { padding-right: 0; margin-bottom: 5px;} // Adjust padding
           .alert-time { width: 100%; text-align: left; margin-bottom: 5px; }
       }
       // Hide less critical columns in tables if needed
       .el-table {
           // Example: Hide Role column on small screens
            .el-table__column[prop='role'] { display: none; }
       }
    }
  
     @media (max-width: 576px) {
        .dashboard-header .control-group {
            flex-direction: column;
             align-items: stretch;
              .last-update { justify-content: center;}
        }
        .health-indicator .indicator-item { min-width: 100%; } // Full width indicators
         .alert-summary-card .alert-item { min-width: 100%; } // Full width alerts
          .el-table {
               // Example: Hide Pods column too
               .el-table__column[prop='pods'] { display: none; }
          }
     }
  }
  </style>