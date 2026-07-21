<template>
  <div class="resource-summary-container">
    <!-- Header -->
    <div class="summary-header">
      <div class="header-content">
        <div class="title-section">
          <h3 class="main-title">{{ $t('clusterOverview.resourceOverview') }}</h3>
          <el-tag type="success" size="small" class="status-tag">{{ $t('clusterOverview.running') }}</el-tag>
        </div>
        <div class="header-actions">
          <el-button :icon="Refresh" text :loading="loading" @click="fetchSummary" class="refresh-btn">
            {{ $t('actions.refresh') }}
          </el-button>
        </div>
      </div>
    </div>

    <!-- Statistics Overview -->
    <div v-if="summaryData && displayItems.length > 0" class="stats-overview">
      <div class="stat-item">
        <div class="stat-number">{{ getTotalResources() }}</div>
        <div class="stat-label">{{ $t('clusterOverview.totalResources') }}</div>
      </div>
      <div class="stat-item">
        <div class="stat-number">{{ getActiveResources() }}</div>
        <div class="stat-label">{{ $t('clusterOverview.activeTypes') }}</div>
      </div>
      <div class="stat-item">
        <div class="stat-number">{{ summaryData.namespaces || 0 }}</div>
        <div class="stat-label">{{ $t('clusterOverview.namespaces') }}</div>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="error" class="error-state">
      <el-icon class="error-icon">
        <Warning />
      </el-icon>
      <div class="error-text">{{ error }}</div>
      <el-button type="primary" size="small" @click="fetchSummary">{{ $t('common.retry') }}</el-button>
    </div>

    <!-- Loading State -->
    <div v-else-if="loading && !summaryData" class="loading-state">
      <el-icon class="loading-icon">
        <Loading />
      </el-icon>
      <span>{{ $t('common.loading') }}</span>
    </div>

    <!-- Resource Grid -->
    <div v-else-if="summaryData" class="resources-grid">
      <TransitionGroup name="resource" tag="div" class="grid-container">
        <div v-for="(item, index) in displayItems" :key="item.key" class="resource-card"
          :style="{ '--delay': index * 0.02 + 's' }" @click="handleResourceClick(item)">
          <div class="card-content">
            <div class="status-indicator" :class="getResourceStatusClass(summaryData[item.key])">
              <div class="indicator-dot"></div>
            </div>
            <div class="resource-info">
              <div class="resource-name">{{ item.label }}</div>
              <div class="resource-count">{{ formatCount(summaryData[item.key]) }}</div>
            </div>
            <el-tag :type="getResourceStatus(summaryData[item.key]).type" size="small" effect="plain"
              class="resource-status">
              {{ getResourceStatus(summaryData[item.key]).text }}
            </el-tag>
          </div>
        </div>
      </TransitionGroup>

      <!-- Empty State -->
      <div v-if="displayItems.length === 0" class="empty-state">
        <el-icon class="empty-icon">
          <Files />
        </el-icon>
        <div class="empty-text">{{ $t('clusterOverview.noResourceData') }}</div>
        <el-button type="primary" size="small" @click="fetchSummary">{{ $t('actions.refresh') }}</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { ElMessage } from 'element-plus';
import { request } from '@/utils/service'; // Adjust path
import {
  DataAnalysis, Refresh, Platform, CollectionTag, Files, TakeawayBox, Service,
  Coin, MessageBox, SetUp, Notebook, Key, Lock, Connection, Warning, // Added more icons
  Loading, Clock, CircleCheck, View, TrendCharts, SuccessFilled, InfoFilled,
  CaretTop, Minus
} from '@element-plus/icons-vue';

interface ResourceSummaryData {
  nodes?: number | null;
  namespaces?: number | null;
  pods?: number | null;
  deployments?: number | null;
  services?: number | null;
  persistentVolumes?: number | null;
  pvcs?: number | null;
  statefulSets?: number | null;
  daemonSets?: number | null;
  configMaps?: number | null;
  secrets?: number | null;
  ingresses?: number | null;
  // Add keys corresponding to the Go model
}

interface DisplayConfig {
  key: keyof ResourceSummaryData;
  label: string;
  icon: any; // Vue component type
  color: string;
}

const loading = ref(false);
const summaryData = ref<ResourceSummaryData | null>(null);
const error = ref<string | null>(null);
const lastUpdateTime = ref<string>('');

// Update last update time
const updateLastUpdateTime = () => {
  const now = new Date();
  const month = String(now.getMonth() + 1).padStart(2, '0');
  const day = String(now.getDate()).padStart(2, '0');
  let hour = now.getHours();
  const ampm = hour >= 12 ? 'PM' : 'AM';
  if (hour > 12) hour -= 12;
  if (hour === 0) hour = 12;
  const minute = String(now.getMinutes()).padStart(2, '0');
  lastUpdateTime.value = `${month}/${day} ${String(hour).padStart(2, '0')}:${minute} ${ampm}`;
};

// Define the display order, labels, icons, and colors
const displayConfig: DisplayConfig[] = [
  { key: 'nodes', label: 'Nodes', icon: Platform, color: '#409EFF' },
  { key: 'namespaces', label: 'Namespaces', icon: CollectionTag, color: '#67C23A' },
  { key: 'pods', label: 'Pods', icon: Files, color: '#E6A23C' },
  { key: 'deployments', label: 'Deployments', icon: TakeawayBox, color: '#F56C6C' },
  { key: 'statefulSets', label: 'StatefulSets', icon: SetUp, color: '#a774d1' }, // Purple
  { key: 'daemonSets', label: 'DaemonSets', icon: Notebook, color: '#7f8c8d' }, // Grey
  { key: 'services', label: 'Services', icon: Service, color: '#3498DB' }, // Blue
  { key: 'ingresses', label: 'Ingresses', icon: Connection, color: '#1ABC9C' }, // Turquoise
  { key: 'persistentVolumes', label: 'PVs', icon: Coin, color: '#f39c12' }, // Orange
  { key: 'pvcs', label: 'PVCs', icon: MessageBox, color: '#f1c40f' }, // Yellow
  { key: 'configMaps', label: 'ConfigMaps', icon: Key, color: '#2ecc71' }, // Emerald
  { key: 'secrets', label: 'Secrets', icon: Lock, color: '#e74c3c' }, // Red
  // Add more resources here
];

// Filter config based on available data keys from backend
const displayItems = computed(() => {
  if (!summaryData.value) return [];
  return displayConfig.filter(item => summaryData.value?.[item.key] !== undefined && summaryData.value?.[item.key] !== null);
});

const fetchSummary = async () => {
  loading.value = true;
  error.value = null;
  try {
    const response = await request<{ code: number; data: ResourceSummaryData; message: string }>({
      url: '/api/v1/summary/resources',
      method: 'get',
      timeout: 15000 // 增加超时时间，因为需要从K8s集群获取数据
      // 移除硬编码的baseURL，让它使用service.ts中的配置
    });
    if (response.code === 200 && response.data) {
      summaryData.value = response.data;
      updateLastUpdateTime();
    } else {
      throw new Error(response.message || 'Data format error');
    }
  } catch (err: any) {
    console.error("Failed to fetch resource summary:", err);
    error.value = err.message || 'Network request failed';
  } finally {
    loading.value = false;
  }
};

// Get resource status
const getResourceStatus = (count: number | null | undefined) => {
  if (count === null || count === undefined) {
    return { type: 'info', text: 'Unknown' };
  }
  if (count === 0) {
    return { type: 'info', text: 'Idle' };
  }
  if (count > 20) {
    return { type: 'danger', text: 'High' };
  }
  if (count > 10) {
    return { type: 'warning', text: 'Medium' };
  }
  return { type: 'success', text: 'Normal' };
};

// Get resource status class for indicator
const getResourceStatusClass = (count: number | null | undefined) => {
  if (count === null || count === undefined) {
    return 'status-unknown';
  }
  if (count === 0) {
    return 'status-idle';
  }
  if (count > 20) {
    return 'status-danger';
  }
  if (count > 10) {
    return 'status-warning';
  }
  return 'status-normal';
};

// Get resource progress (simulated usage)
const getResourceProgress = (key: string, count: number | null | undefined): number => {
  if (count === null || count === undefined || count === 0) return 0;

  // Set different progress calculation logic for different resource types
  const progressMap: Record<string, number> = {
    'nodes': Math.min((count / 10) * 100, 100),
    'pods': Math.min((count / 100) * 100, 100),
    'namespaces': Math.min((count / 20) * 100, 100),
    'deployments': Math.min((count / 50) * 100, 100),
    'services': Math.min((count / 30) * 100, 100),
  };

  return Math.round(progressMap[key] || Math.min((count / 20) * 100, 100));
};

// Handle resource card click
const handleResourceClick = (item: DisplayConfig) => {
  ElMessage.info(`Clicked ${item.label}, feature under development...`);
};

// Get total resources count
const getTotalResources = (): number => {
  if (!summaryData.value) return 0;
  return Object.values(summaryData.value).reduce((total, count) => {
    return total + (count || 0);
  }, 0);
};

// Get active resources count (non-zero resources)
const getActiveResources = (): number => {
  if (!summaryData.value) return 0;
  return Object.values(summaryData.value).filter(count => count && count > 0).length;
};

const formatCount = (count: number | null | undefined): string | number => {
  // Backend uses pointers, so check for null explicitly
  if (count === null || count === undefined) {
    return '?'; // Indicate data wasn't fetched or error occurred for this type
  }
  return count;
};

onMounted(() => {
  fetchSummary();
});

// Expose fetchSummary if you want to call it from parent
// defineExpose({ fetchSummary });

</script>

<style lang="scss" scoped>
.resource-summary-container {
  background: var(--el-bg-color);
  border-radius: 8px;
  border: 1px solid var(--el-border-color-lighter);
  margin-bottom: 20px;
  overflow: hidden;
  transition: all 0.2s ease;

  &:hover {
    border-color: var(--el-border-color);
    box-shadow: 0 4px 12px var(--el-box-shadow-light);
  }
}

// 简洁头部
.summary-header {
  padding: 16px 20px 12px;
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

    .main-title {
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

// 统计概览
.stats-overview {
  display: flex;
  padding: 12px 20px;
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-lighter);

  .stat-item {
    flex: 1;
    text-align: center;

    &:not(:last-child) {
      border-right: 1px solid var(--el-border-color-lighter);
    }

    .stat-number {
      font-size: 20px;
      font-weight: 700;
      color: var(--el-text-color-primary);
      line-height: 1;
      margin-bottom: 3px;
    }

    .stat-label {
      font-size: 11px;
      color: var(--el-text-color-regular);
      font-weight: 500;
    }
  }
}

// 错误状态
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 32px 24px;
  text-align: center;

  .error-icon {
    font-size: 32px;
    color: var(--el-color-danger);
  }

  .error-text {
    font-size: 14px;
    color: var(--el-text-color-regular);
    max-width: 300px;
  }
}

// 加载状态
.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 32px 24px;
  color: var(--el-text-color-regular);
  font-size: 14px;

  .loading-icon {
    font-size: 16px;
    animation: spin 1s linear infinite;
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

// 资源网格
.resources-grid {
  padding: 8px 12px 12px;

  .grid-container {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
    gap: 6px;
  }

  .resource-card {
    background: var(--el-bg-color);
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 4px;
    transition: all 0.2s ease;
    cursor: pointer;
    animation: slideIn 0.3s ease-out;
    animation-delay: var(--delay);
    animation-fill-mode: both;

    &:hover {
      border-color: var(--el-border-color);
      box-shadow: 0 2px 8px var(--el-box-shadow-light);
    }

    .card-content {
      padding: 6px;
      display: flex;
      flex-direction: column;
      gap: 4px;
      min-height: 60px;
    }

    .status-indicator {
      width: 20px;
      height: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
      align-self: flex-start;

      .indicator-dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        transition: all 0.2s ease;
      }

      &.status-normal .indicator-dot {
        background: #10b981;
        box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
      }

      &.status-warning .indicator-dot {
        background: #f59e0b;
        box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.2);
      }

      &.status-danger .indicator-dot {
        background: #ef4444;
        box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.2);
      }

      &.status-idle .indicator-dot {
        background: #64748b;
        box-shadow: 0 0 0 2px rgba(100, 116, 139, 0.2);
      }

      &.status-unknown .indicator-dot {
        background: #64748b;
        box-shadow: 0 0 0 2px rgba(100, 116, 139, 0.2);
      }
    }

    .resource-info {
      flex: 1;

      .resource-name {
        font-size: 11px;
        font-weight: 500;
        color: var(--el-text-color-regular);
        margin-bottom: 2px;
        line-height: 1.2;
      }

      .resource-count {
        font-size: 16px;
        font-weight: 700;
        color: var(--el-text-color-primary);
        line-height: 1;
      }
    }

    .resource-status {
      align-self: flex-end;

      .el-tag {
        font-size: 9px;
        font-weight: 500;
        padding: 1px 4px;
        border-radius: 2px;
      }
    }
  }
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 空状态
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 32px 24px;
  text-align: center;

  .empty-icon {
    font-size: 32px;
    color: var(--el-text-color-placeholder);
  }

  .empty-text {
    font-size: 14px;
    color: var(--el-text-color-regular);
  }
}

// 过渡动画
.resource-enter-active,
.resource-leave-active {
  transition: all 0.2s ease;
}

.resource-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.resource-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

// 响应式设计
@media (max-width: 768px) {
  .resource-summary-container {
    margin: 0 -16px 20px;
    border-radius: 0;
    border-left: none;
    border-right: none;
  }

  .summary-header {
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

  .stats-overview {
    padding: 12px 20px;

    .stat-item .stat-number {
      font-size: 20px;
    }
  }

  .resources-grid {
    padding: 10px 20px 16px;

    .grid-container {
      grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
      gap: 8px;
    }

    .resource-card {
      .card-content {
        padding: 10px;
        min-height: 80px;
      }

      .resource-icon {
        width: 28px;
        height: 28px;

        .el-icon {
          font-size: 14px;
        }
      }

      .resource-info {
        .resource-name {
          font-size: 12px;
        }

        .resource-count {
          font-size: 18px;
        }
      }
    }
  }

  .error-state,
  .empty-state {
    padding: 24px 20px;
  }

  .loading-state {
    padding: 24px 20px;
  }
}

@media (max-width: 480px) {
  .stats-overview {
    flex-direction: column;
    gap: 12px;

    .stat-item {
      &:not(:last-child) {
        border-right: none;
        border-bottom: 1px solid var(--el-border-color-lighter);
        padding-bottom: 12px;
      }
    }
  }

  .resources-grid {
    .grid-container {
      grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
      gap: 6px;
    }

    .resource-card {
      .card-content {
        padding: 8px;
        min-height: 70px;
        gap: 6px;
      }

      .resource-icon {
        width: 24px;
        height: 24px;

        .el-icon {
          font-size: 12px;
        }
      }

      .resource-info {
        .resource-name {
          font-size: 11px;
        }

        .resource-count {
          font-size: 16px;
        }
      }

      .resource-status .el-tag {
        font-size: 9px;
        padding: 1px 4px;
      }
    }
  }
}
</style>