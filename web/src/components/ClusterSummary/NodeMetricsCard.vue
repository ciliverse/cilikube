<template>
  <div class="node-metrics-card">
    <!-- Header -->
    <div class="card-header">
      <div class="header-content">
        <div class="title-section">
          <h3 class="card-title">{{ $t('clusterOverview.nodeMetrics') }}</h3>
        </div>
        <div class="header-actions">
          <el-button :icon="Refresh" text :loading="loading" @click="fetchNodeMetrics" class="refresh-btn">
            {{ $t('actions.refresh') }}
          </el-button>
        </div>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="error" class="error-state">
      <el-icon class="error-icon">
        <Warning />
      </el-icon>
      <div class="error-text">{{ error }}</div>
      <el-button type="primary" size="small" @click="fetchNodeMetrics">
        {{ $t('common.retry') }}
      </el-button>
    </div>

    <!-- Loading State -->
    <div v-else-if="loading && !nodeMetrics" class="loading-state">
      <el-icon class="loading-icon">
        <Loading />
      </el-icon>
      <span>{{ $t('common.loading') }}</span>
    </div>

    <!-- Metrics Content -->
    <div v-else-if="nodeMetrics" class="metrics-content">
      <!-- Overview Stats -->
      <div class="overview-stats">
        <div class="stat-item">
          <div class="stat-icon cpu" :class="cpuStatusClass">
            <div class="status-indicator"></div>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ nodeMetrics.avgCpuUsagePercent }}%</div>
            <div class="stat-label">{{ $t('nodeMetrics.avgCpuUsage') }}</div>
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-icon memory" :class="memoryStatusClass">
            <div class="status-indicator"></div>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ nodeMetrics.avgMemoryUsagePercent }}%</div>
            <div class="stat-label">{{ $t('nodeMetrics.avgMemoryUsage') }}</div>
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-icon nodes" :class="nodesStatusClass">
            <div class="status-indicator"></div>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ nodeMetrics.healthyNodes }}/{{ nodeMetrics.totalNodes }}</div>
            <div class="stat-label">{{ $t('nodeMetrics.healthyNodes') }}</div>
          </div>
        </div>
      </div>

      <!-- Resource Summary -->
      <div class="resource-summary">
        <div class="resource-item">
          <div class="resource-label">{{ $t('nodeMetrics.totalCpuCores') }}</div>
          <div class="resource-value">{{ nodeMetrics.totalCpuCores }} {{ $t('nodeMetrics.cores') }}</div>
        </div>
        <div class="resource-item">
          <div class="resource-label">{{ $t('nodeMetrics.totalMemory') }}</div>
          <div class="resource-value">{{ nodeMetrics.totalMemoryGB }} GB</div>
        </div>
      </div>

      <!-- Resource Requests -->
      <div class="requests-section">
        <div class="section-title">{{ $t('nodeMetrics.resourceRequests') }}</div>
        <div class="requests-grid">
          <div class="request-item">
            <div class="request-header">
              <span class="request-label">{{ $t('nodeMetrics.cpuRequests') }}</span>
              <div class="request-value-section">
                <div class="percentage-badge cpu">
                  <span class="percentage-number">{{ nodeMetrics.totalCpuRequestsPercent }}</span>
                  <span class="percentage-symbol">%</span>
                </div>
                <span class="request-detail">{{ nodeMetrics.totalCpuRequests }} {{ $t('nodeMetrics.cores') }}</span>
              </div>
            </div>
            <div class="request-progress">
              <div class="progress-bar">
                <div class="progress-fill cpu" :style="{ width: nodeMetrics.totalCpuRequestsPercent + '%' }"></div>
              </div>
            </div>
          </div>
          <div class="request-item">
            <div class="request-header">
              <span class="request-label">{{ $t('nodeMetrics.memoryRequests') }}</span>
              <div class="request-value-section">
                <div class="percentage-badge memory">
                  <span class="percentage-number">{{ nodeMetrics.totalMemoryRequestsPercent }}</span>
                  <span class="percentage-symbol">%</span>
                </div>
                <span class="request-detail">{{ nodeMetrics.totalMemoryRequests }} GB</span>
              </div>
            </div>
            <div class="request-progress">
              <div class="progress-bar">
                <div class="progress-fill memory" :style="{ width: nodeMetrics.totalMemoryRequestsPercent + '%' }">
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="empty-state">
      <el-icon class="empty-icon">
        <Monitor />
      </el-icon>
      <div class="empty-text">{{ $t('nodeMetrics.noData') }}</div>
      <el-button type="primary" size="small" @click="fetchNodeMetrics">
        {{ $t('actions.refresh') }}
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getNodeMetrics, parseNodeMetricsSummary } from '@/api/node-metrics'
import type { NodeMetricsSummary } from '@/types/node-metrics'
import {
  Monitor, Refresh, Warning, Loading
} from '@element-plus/icons-vue'

const loading = ref(false)
const nodeMetrics = ref<NodeMetricsSummary | null>(null)
const error = ref<string | null>(null)
const lastUpdateTime = ref<string>('')

// Update last update time
const updateLastUpdateTime = () => {
  const now = new Date()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  let hour = now.getHours()
  const ampm = hour >= 12 ? 'PM' : 'AM'
  if (hour > 12) hour -= 12
  if (hour === 0) hour = 12
  const minute = String(now.getMinutes()).padStart(2, '0')
  lastUpdateTime.value = `${month}/${day} ${String(hour).padStart(2, '0')}:${minute} ${ampm}`
}

// Get overall health status
const getOverallHealthStatus = () => {
  if (!nodeMetrics.value) {
    return { type: 'info', text: 'Unknown' }
  }

  const { healthyNodes, totalNodes, avgCpuUsagePercent, avgMemoryUsagePercent } = nodeMetrics.value

  if (healthyNodes < totalNodes) {
    return { type: 'danger', text: 'Error' }
  }

  if (avgCpuUsagePercent > 80 || avgMemoryUsagePercent > 80) {
    return { type: 'warning', text: 'High Load' }
  }

  if (avgCpuUsagePercent > 60 || avgMemoryUsagePercent > 60) {
    return { type: 'warning', text: 'Medium Load' }
  }

  return { type: 'success', text: 'Normal' }
}

// CPU status indicator class
const cpuStatusClass = computed(() => {
  if (!nodeMetrics.value) return 'blue'
  const percent = nodeMetrics.value.avgCpuUsagePercent
  if (percent < 30) return 'blue'
  if (percent < 80) return 'green'
  return 'red'
})

// Memory status indicator class
const memoryStatusClass = computed(() => {
  if (!nodeMetrics.value) return 'blue'
  const percent = nodeMetrics.value.avgMemoryUsagePercent
  if (percent < 30) return 'blue'
  if (percent < 80) return 'green'
  return 'red'
})

// Nodes status indicator class (green if all healthy, red otherwise)
const nodesStatusClass = computed(() => {
  if (!nodeMetrics.value) return 'gray'
  const { healthyNodes, totalNodes } = nodeMetrics.value
  return healthyNodes === totalNodes ? 'green' : 'red'
})

// Fetch node metrics
const fetchNodeMetrics = async () => {
  loading.value = true
  error.value = null

  try {
    const response = await getNodeMetrics()
    if (response.code === 200 && response.data) {
      nodeMetrics.value = parseNodeMetricsSummary(response)
      updateLastUpdateTime()
    } else {
      throw new Error(response.message || '获取节点监控数据失败')
    }
  } catch (err: any) {
    console.error('Failed to fetch node metrics:', err)
    error.value = err.message || '网络请求失败'
    ElMessage.error(error.value || '网络请求失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchNodeMetrics()
})

// Expose fetchNodeMetrics for parent component
defineExpose({ fetchNodeMetrics })
</script>

<style lang="scss" scoped>
.node-metrics-card {
  background: var(--el-bg-color);
  border-radius: 8px;
  border: 1px solid var(--el-border-color-lighter);
  margin-bottom: 0;
  overflow: hidden;
  transition: all 0.2s ease;
  min-height: 400px;

  &:hover {
    border-color: var(--el-border-color);
    box-shadow: 0 4px 12px var(--el-box-shadow-light);
  }
}

.card-header {
  padding: 12px 16px 8px;
  border-bottom: 1px solid var(--el-border-color-lighter);

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .title-section {
    display: flex;
    align-items: center;
    gap: 12px;

    .title-icon {
      font-size: 18px;
      color: var(--el-color-primary);
    }

    .card-title {
      font-size: 18px;
      font-weight: 600;
      color: var(--el-text-color-primary);
      margin: 0;
    }

    .status-tag {
      font-size: 12px;
      font-weight: 500;
      padding: 2px 8px;
      border-radius: 4px;
    }
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 12px;

    .last-update {
      font-size: 12px;
      color: var(--el-text-color-regular);
    }

    .refresh-btn {
      font-size: 12px;
      font-weight: 500;
      color: var(--el-color-primary);
    }
  }
}

.error-state,
.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 32px 24px;
  text-align: center;

  .error-icon,
  .empty-icon {
    font-size: 32px;
    color: var(--el-color-danger);
  }

  .loading-icon {
    font-size: 16px;
    animation: spin 1s linear infinite;
  }

  .error-text,
  .empty-text {
    font-size: 14px;
    color: var(--el-text-color-regular);
    max-width: 300px;
  }
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}

.metrics-content {
  padding: 20px 20px 24px;
}

.overview-stats {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 10px;
  margin-bottom: 14px;

  .stat-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 12px;
    border-radius: 6px;
    transition: all 0.2s ease;

    &:hover {
      background: var(--el-bg-color-page);
    }

    .stat-icon {
      width: 32px;
      height: 32px;
      border-radius: 6px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 15px;
      flex-shrink: 0;

      .status-indicator {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        transition: all 0.2s ease;
      }

      &.cpu {
        background: rgba(64, 158, 255, 0.1);
        color: #409EFF;

        &.blue .status-indicator {
          background: #3b82f6;
          box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
        }

        &.green .status-indicator {
          background: #10b981;
          box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2);
        }

        &.red .status-indicator {
          background: #ef4444;
          box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.2);
        }
      }

      &.memory {
        background: rgba(103, 194, 58, 0.1);
        color: #67C23A;

        &.blue .status-indicator {
          background: #3b82f6;
          box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
        }

        &.green .status-indicator {
          background: #10b981;
          box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2);
        }

        &.red .status-indicator {
          background: #ef4444;
          box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.2);
        }

        .memory-text {
          font-size: 12px;
          font-weight: 600;
          letter-spacing: 0.5px;
        }
      }

      &.nodes {
        background: rgba(230, 162, 60, 0.1);
        color: #E6A23C;

        &.gray .status-indicator {
          background: #64748b;
          box-shadow: 0 0 0 3px rgba(100, 116, 139, 0.2);
        }

        &.green .status-indicator {
          background: #10b981;
          box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2);
        }

        &.red .status-indicator {
          background: #ef4444;
          box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.2);
        }
      }
    }

    .stat-info {
      flex: 1;
      min-width: 0;

      .stat-value {
        font-size: 16px;
        font-weight: 700;
        color: var(--el-text-color-primary);
        line-height: 1.2;
        margin-bottom: 3px;
      }

      .stat-label {
        font-size: 10px;
        color: var(--el-text-color-regular);
        font-weight: 500;
        line-height: 1;
      }
    }
  }
}

.resource-summary {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-bottom: 16px;

  .resource-item {
    padding: 10px 12px;
    border-radius: 6px;
    transition: all 0.2s ease;

    &:hover {
      background: var(--el-bg-color-page);
    }

    .resource-label {
      font-size: 11px;
      color: var(--el-text-color-regular);
      margin-bottom: 3px;
      font-weight: 500;
    }

    .resource-value {
      font-size: 15px;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }
  }
}

.requests-section {
  .section-title {
    font-size: 12px;
    font-weight: 600;
    color: var(--el-text-color-primary);
    margin-bottom: 10px;
    padding-bottom: 4px;
    border-bottom: 1px solid var(--el-border-color-lighter);
  }

  .requests-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 10px;

    .request-item {
      padding: 10px 12px;
      border-radius: 6px;
      border: 1px solid var(--el-border-color-lighter);
      transition: all 0.2s ease;

      &:hover {
        border-color: var(--el-border-color);
        box-shadow: 0 2px 8px var(--el-box-shadow-light);
      }

      .request-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 8px;

        .request-label {
          font-size: 11px;
          color: var(--el-text-color-primary);
          font-weight: 700;
        }

        .request-value-section {
          display: flex;
          align-items: center;
          gap: 8px;

          .percentage-badge {
            display: flex;
            align-items: baseline;
            padding: 4px 8px;
            border-radius: 6px;
            font-weight: 700;
            min-width: 50px;
            justify-content: center;

            &.cpu {
              background: var(--el-bg-color-page);
              border: 1px solid var(--el-border-color);
              color: var(--el-text-color-primary);
            }

            &.memory {
              background: var(--el-bg-color-page);
              border: 1px solid var(--el-border-color);
              color: var(--el-text-color-primary);
            }

            .percentage-number {
              font-size: 16px;
              font-weight: 700;
            }

            .percentage-symbol {
              font-size: 12px;
              font-weight: 500;
              margin-left: 1px;
            }
          }

          .request-detail {
            font-size: 10px;
            color: var(--el-text-color-regular);
            font-weight: 500;
          }
        }
      }

      .request-progress {
        .progress-bar {
          width: 100%;
          height: 5px;
          background: var(--el-bg-color);
          border-radius: 3px;
          overflow: hidden;
          border: 1px solid var(--el-border-color-lighter);

          .progress-fill {
            height: 100%;
            border-radius: 3px;
            transition: width 0.3s ease;

            &.cpu {
              background: linear-gradient(90deg, #409EFF 0%, #66B1FF 100%);
            }

            &.memory {
              background: linear-gradient(90deg, #67C23A 0%, #85CE61 100%);
            }
          }
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .node-metrics-card {
    margin: 0 -16px 20px;
    border-radius: 0;
    border-left: none;
    border-right: none;
  }

  .card-header {
    padding: 16px 20px 12px;

    .header-content {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .title-section {
      align-self: stretch;
      justify-content: space-between;
    }

    .header-actions {
      align-self: stretch;
      justify-content: space-between;
    }
  }

  .metrics-content {
    padding: 16px 20px;
  }

  .overview-stats {
    grid-template-columns: 1fr;
    gap: 12px;
    margin-bottom: 20px;

    .stat-item {
      padding: 12px;

      .stat-icon {
        width: 36px;
        height: 36px;
        font-size: 16px;
      }

      .stat-info .stat-value {
        font-size: 18px;
      }
    }
  }

  .resource-summary {
    grid-template-columns: 1fr;
    gap: 12px;
    margin-bottom: 20px;
  }
}
</style>