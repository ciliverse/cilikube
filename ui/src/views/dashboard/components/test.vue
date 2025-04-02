<template>
    <div class="cluster-monitoring-dashboard">
      <!-- 标题与刷新控制 -->
      <div class="dashboard-header">
        <h2>Kubernetes集群监控仪表盘</h2>
        <div class="header-actions">
          <el-select v-model="selectedNamespace" placeholder="全部命名空间" size="small">
            <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
          </el-select>
          <el-button 
            type="primary" 
            size="small" 
            :icon="Refresh" 
            @click="refreshData"
            :loading="refreshing"
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
  
      <!-- 集群概览卡片 -->
      <div class="overview-section">
        <div class="section-title">
          <el-icon><DataAnalysis /></el-icon>
          <span>集群概览</span>
        </div>
        <el-row :gutter="16" class="overview-cards">
          <el-col :span="6" v-for="(item, index) in overviewData" :key="index">
            <div class="overview-card">
              <div class="card-icon" :style="{ backgroundColor: item.color }">
                <el-icon><component :is="item.icon" /></el-icon>
              </div>
              <div class="card-content">
                <div class="card-title">{{ item.title }}</div>
                <div class="card-value">{{ item.value }}</div>
                <div class="card-description">
                  <el-progress :percentage="item.percent" :color="item.color" :show-text="false" />
                </div>
                <div class="card-footer">
                  <span>总{{ item.total }} · 使用率{{ item.percent }}%</span>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
  
      <!-- 资源使用率图表 -->
      <div class="resource-usage-section">
        <div class="section-title">
          <el-icon><PieChartIcon /></el-icon>
          <span>资源使用率</span>
        </div>
        <el-row :gutter="16" class="resource-usage-charts">
          <el-col :span="12">
            <div class="chart-card">
              <div class="chart-title">CPU使用情况</div>
              <div class="chart-container">
                <VChart :option="cpuUsageOption" autoresize />
              </div>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="chart-card">
              <div class="chart-title">内存使用情况</div>
              <div class="chart-container">
                <VChart :option="memoryUsageOption" autoresize />
              </div>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="chart-card">
              <div class="chart-title">存储使用情况</div>
              <div class="chart-container">
                <VChart :option="storageUsageOption" autoresize />
              </div>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="chart-card">
              <div class="chart-title">网络流量</div>
              <div class="chart-container">
                <VChart :option="networkUsageOption" autoresize />
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
  
      <!-- 节点状态表格 -->
      <div class="node-status-section">
        <div class="section-title">
          <el-icon><Loading /></el-icon>
          <span>节点状态</span>
        </div>
        <el-table :data="nodes" stripe style="width: 100%" v-loading="loadingNodes">
          <el-table-column prop="name" label="节点名称" width="180" />
          <el-table-column prop="role" label="角色" width="120">
            <template #default="{ row }">
              <el-tag :type="row.role === 'master' ? '' : 'info'" size="small">
                {{ row.role === 'master' ? '控制节点' : '工作节点' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="row.status === 'Ready' ? 'success' : 'danger'" size="small">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="cpuUsage" label="CPU使用率">
            <template #default="{ row }">
              <el-progress 
                :percentage="row.cpuUsage" 
                :stroke-width="15" 
                :color="getProgressColor(row.cpuUsage)"
              />
            </template>
          </el-table-column>
          <el-table-column prop="memoryUsage" label="内存使用率">
            <template #default="{ row }">
              <el-progress 
                :percentage="row.memoryUsage" 
                :stroke-width="15" 
                :color="getProgressColor(row.memoryUsage)"
              />
            </template>
          </el-table-column>
          <el-table-column prop="pods" label="Pods">
            <template #default="{ row }">
              {{ row.runningPods }}/{{ row.totalPods }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="80">
            <template #default>
              <el-button size="small" type="text">详情</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
  
      <!-- 集群事件日志 -->
      <div class="event-log-section">
        <div class="section-title">
          <el-icon><Document /></el-icon>
          <span>集群事件</span>
        </div>
        <el-tabs v-model="activeEventTab">
          <el-tab-pane label="最新事件" name="recent">
            <el-table :data="recentEvents" style="width: 100%" height="300" v-loading="loadingEvents">
              <el-table-column prop="timestamp" label="时间" width="160" sortable>
                <template #default="{ row }">
                  {{ formatDate(row.timestamp) }}
                </template>
              </el-table-column>
              <el-table-column prop="type" label="类型" width="120">
                <template #default="{ row }">
                  <el-tag :type="row.type === 'Warning' ? 'warning' : ''" size="small">
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
                <VChart :option="eventStatisticsOption" autoresize />
              </div>
              <div class="statistics-table">
                <el-table :data="eventStatistics" border style="width: 100%" height="280">
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
                        :icon="row.trend === 'up' ? ArrowUpBold : ArrowDownBold"
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
    Refresh,
    Clock,
    DataAnalysis,
    PieChart as PieChartIcon,
    Loading,
    Document,
    ArrowUpBold,
    ArrowDownBold,
    DataBoard,
    Collection,
    Box,
    Connection
  } from '@element-plus/icons-vue'
  import dayjs from 'dayjs'
  
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
  
  // 模拟数据
  const namespaces = ['default', 'kube-system', 'monitoring', 'logging', 'dev', 'prod']
  
  const overviewData = ref([
    {
      title: '集群节点',
      value: '5',
      percent: 95,
      total: '',
      icon: DataBoard,
      color: '#409EFF'
    },
    {
      title: '命名空间',
      value: '6',
      percent: 75,
      total: '8',
      icon: Collection,
      color: '#67C23A'
    },
    {
      title: '运行Pods',
      value: '148',
      percent: 82,
      total: '180',
      icon: Box,
      color: '#E6A23C'
    },
    {
      title: '服务数量',
      value: '23',
      percent: 65,
      total: '35',
      icon: Connection,
      color: '#F56C6C'
    }
  ])
  
  const cpuUsageOption = ref({
    tooltip: {
      trigger: 'item'
    },
    legend: {
      top: '5%',
      left: 'center'
    },
    series: [
      {
        name: 'CPU使用率',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: 70, name: '已使用 (70%)' },
          { value: 30, name: '剩余 (30%)' }
        ],
        color: ['#F56C6C', '#E4E7ED']
      }
    ]
  })
  
  const memoryUsageOption = ref({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'value',
      boundaryGap: [0, 0.01],
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      }
    },
    yAxis: {
      type: 'category',
      data: ['总内存', '已使用', '缓存', '可用'],
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      }
    },
    series: [
      {
        name: '内存分配',
        type: 'bar',
        data: [
          { value: 64, itemStyle: { color: '#E4E7ED' } },
          { value: 45, itemStyle: { color: '#F56C6C' } },
          { value: 12, itemStyle: { color: '#E6A23C' } },
          { value: 19, itemStyle: { color: '#67C23A' } }
        ],
        label: {
          show: true,
          position: 'right',
          formatter: '{@value} GB'
        }
      }
    ]
  })
  
  const storageUsageOption = ref({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'line',
        label: {
          backgroundColor: '#6a7985'
        }
      }
    },
    legend: {
      data: ['已使用', '总量']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: [
      {
        type: 'category',
        boundaryGap: false,
        data: ['Master-1', 'Node-1', 'Node-2', 'Node-3']
      }
    ],
    yAxis: [
      {
        type: 'value',
        name: '存储 (GB)'
      }
    ],
    series: [
      {
        name: '已使用',
        type: 'line',
        stack: '总量',
        areaStyle: {
          color: '#F56C6C'
        },
        emphasis: {
          focus: 'series'
        },
        data: [120, 80, 90, 100]
      },
      {
        name: '总量',
        type: 'line',
        stack: '总量',
        areaStyle: {
          color: '#E4E7ED'
        },
        emphasis: {
          focus: 'series'
        },
        data: [200, 200, 200, 200]
      }
    ]
  })
  
  const networkUsageOption = ref({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: ['08:00', '10:00', '12:00', '14:00', '16:00', '18:00', '20:00'],
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      }
    },
    yAxis: {
      type: 'value',
      name: '网络流量 (Mbps)',
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      }
    },
    series: [
      {
        name: '入站',
        type: 'bar',
        stack: '流量',
        data: [120, 200, 150, 180, 210, 190, 230],
        itemStyle: {
          color: '#409EFF'
        }
      },
      {
        name: '出站',
        type: 'bar',
        stack: '流量',
        data: [100, 170, 130, 150, 180, 160, 190],
        itemStyle: {
          color: '#67C23A'
        }
      }
    ]
  })
  
  const nodes = ref([
    {
      name: 'master-1',
      role: 'master',
      status: 'Ready',
      cpuUsage: 65,
      memoryUsage: 75,
      totalPods: 64,
      runningPods: 60
    },
    {
      name: 'node-1',
      role: 'worker',
      status: 'Ready',
      cpuUsage: 45,
      memoryUsage: 60,
      totalPods: 128,
      runningPods: 110
    },
    {
      name: 'node-2',
      role: 'worker',
      status: 'Ready',
      cpuUsage: 30,
      memoryUsage: 55,
      totalPods: 128,
      runningPods: 98
    },
    {
      name: 'node-3',
      role: 'worker',
      status: 'NotReady',
      cpuUsage: 85,
      memoryUsage: 90,
      totalPods: 128,
      runningPods: 128
    },
    {
      name: 'node-4',
      role: 'worker',
      status: 'Ready',
      cpuUsage: 60,
      memoryUsage: 65,
      totalPods: 128,
      runningPods: 105
    }
  ])
  
  const recentEvents = ref([
    {
      timestamp: '2023-10-15T09:23:17Z',
      type: 'Warning',
      object: 'pod/nginx-deployment-75675f5897-58xj7',
      namespace: 'default',
      reason: 'FailedScheduling',
      message: '0/4 nodes available: 3 Insufficient cpu, 1 Insufficient memory.'
    },
    {
      timestamp: '2023-10-15T09:20:45Z',
      type: 'Normal',
      object: 'deployment/nginx-deployment',
      namespace: 'default',
      reason: 'ScalingReplicaSet',
      message: 'Scaled up replica set nginx-deployment-75675f5897 to 1'
    },
    {
      timestamp: '2023-10-15T09:18:32Z',
      type: 'Warning',
      object: 'node/node-3',
      namespace: '',
      reason: 'NodeNotReady',
      message: 'Node node-3 status is now: NodeNotReady'
    },
    {
      timestamp: '2023-10-15T09:15:12Z',
      type: 'Normal',
      object: 'service/redis-master',
      namespace: 'prod',
      reason: 'UpdatedLoadBalancer',
      message: 'Updated load balancer with new hosts'
    },
    {
      timestamp: '2023-10-15T09:12:58Z',
      type: 'Normal',
      object: 'pod/mysql-56f4d7f65d-f5z2n',
      namespace: 'prod',
      reason: 'Pulled',
      message: 'Successfully pulled image "mysql:5.7"'
    },
    {
      timestamp: '2023-10-15T09:10:31Z',
      type: 'Warning',
      object: 'pod/cron-job-1625868000-abcde',
      namespace: 'dev',
      reason: 'Failed',
      message: 'Job has reached the specified backoff limit'
    }
  ])
  
  const eventStatistics = ref([
    { type: 'Normal', count: 42, percentage: 70, trend: 'up' },
    { type: 'Warning', count: 18, percentage: 30, trend: 'down' }
  ])
  
  const eventStatisticsOption = ref({
    tooltip: {
      trigger: 'item'
    },
    legend: {
      top: '5%',
      left: 'center'
    },
    series: [
      {
        name: '事件统计',
        type: 'pie',
        radius: ['50%', '70%'],
        center: ['50%', '60%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: 42, name: '正常事件 (70%)', itemStyle: { color: '#67C23A' } },
          { value: 18, name: '警告事件 (30%)', itemStyle: { color: '#E6A23C' } }
        ]
      }
    ]
  })
  
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
    // 模拟API延迟
    setTimeout(() => {
      lastUpdateTime.value = dayjs().format('YYYY-MM-DD HH:mm:ss')
      refreshing.value = false
    }, 1000)
  }
  
  // 组件挂载时加载数据
  onMounted(() => {
    loadingNodes.value = true
    loadingEvents.value = true
    
    // 模拟API调用
    setTimeout(() => {
      loadingNodes.value = false
      loadingEvents.value = false
    }, 1500)
  })
  </script>
  
  <style lang="scss" scoped>
  .cluster-monitoring-dashboard {
    padding: 20px;
    background-color: #f5f7fa;
    min-height: calc(100vh - 64px);
  }
  
  .dashboard-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 15px;
    border-bottom: 1px solid #ebeef5;
    
    h2 {
      margin: 0;
      font-size: 22px;
      color: #303133;
    }
    
    .header-actions {
      display: flex;
      align-items: center;
      gap: 15px;
      
      .last-update {
        display: flex;
        align-items: center;
        gap: 5px;
        font-size: 12px;
        color: #909399;
        
        .el-icon {
          font-size: 14px;
        }
      }
    }
  }
  
  .section-title {
    display: flex;
    align-items: center;
    margin: 25px 0 15px;
    font-size: 16px;
    font-weight: 500;
    color: #606266;
    
    .el-icon {
      margin-right: 8px;
      font-size: 18px;
    }
  }
  
  .overview-cards {
    margin-bottom: 20px;
    
    .overview-card {
      display: flex;
      background: #fff;
      border-radius: 4px;
      box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
      padding: 15px;
      height: 120px;
      
      .card-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 50px;
        height: 50px;
        border-radius: 50%;
        color: white;
        font-size: 22px;
        margin-right: 15px;
        flex-shrink: 0;
        
        .el-icon {
          display: flex;
        }
      }
      
      .card-content {
        flex: 1;
        display: flex;
        flex-direction: column;
        
        .card-title {
          font-size: 14px;
          color: #909399;
          margin-bottom: 5px;
        }
        
        .card-value {
          font-size: 24px;
          font-weight: bold;
          color: #303133;
          margin-bottom: 10px;
        }
        
        .card-description {
          flex: 1;
          
          .el-progress {
            height: 6px;
            margin-top: 5px;
          }
        }
        
        .card-footer {
          font-size: 12px;
          color: #909399;
        }
      }
    }
  }
  
  .resource-usage-charts {
    margin-bottom: 20px;
    
    .chart-card {
      background: #fff;
      border-radius: 4px;
      box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
      padding: 15px;
      margin-bottom: 16px;
      
      .chart-title {
        font-size: 15px;
        font-weight: 500;
        color: #606266;
        margin-bottom: 10px;
      }
      
      .chart-container {
        width: 100%;
        height: 250px;
      }
    }
  }
  
  .node-status-section,
  .event-log-section {
    background: #fff;
    border-radius: 4px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    padding: 15px;
    margin-bottom: 20px;
  }
  
  .event-statistics-container {
    display: flex;
    gap: 20px;
    
    .statistics-chart {
      flex: 1;
      height: 300px;
    }
    
    .statistics-table {
      width: 320px;
      
      .percentage-bar {
        position: relative;
        height: 18px;
        
        .percentage {
          position: absolute;
          left: 5px;
          z-index: 1;
          font-size: 12px;
          line-height: 18px;
          color: #fff;
        }
        
        .percentage-progress {
          position: absolute;
          height: 100%;
          background-color: #409eff;
          border-radius: 3px;
        }
      }
    }
  }
  
  :deep(.el-progress-bar) {
    padding-right: 0;
  }
  
  :deep(.el-progress__text) {
    min-width: 40px;
    font-size: 12px !important;
  }
  </style>