<template>
  <div class="cluster-overview">
    <div class="page-header">
      <h1 class="page-title">{{ $t('menu.clusterOverview') }}</h1>
    </div>

    <div class="overview-content">
      <!-- 上半部分：资源概览 -->
      <div class="top-section">
        <ResourceSummaryCard v-if="resourceCardExists" />
        <el-card v-else class="resource-placeholder">
          <div class="placeholder-content">
            <el-icon class="placeholder-icon">
              <DataAnalysis />
            </el-icon>
            <h3>{{ $t('clusterOverview.resourceOverview') }}</h3>
            <p>{{ $t('clusterOverview.resourceDescription') }}</p>
          </div>
        </el-card>
      </div>

      <!-- 下半部分：监控和事件 -->
      <div class="bottom-section">
        <el-row :gutter="20">
          <el-col :span="12">
            <NodeMetricsCard v-if="nodeMetricsCardExists" />
            <el-card v-else class="metrics-placeholder">
              <div class="placeholder-content">
                <el-icon class="placeholder-icon">
                  <Monitor />
                </el-icon>
                <h3>{{ $t('clusterOverview.nodeMetrics') }}</h3>
                <p>节点监控数据加载中...</p>
              </div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card class="events-card">
              <template #header>
                <div class="card-header-content">
                  <div class="header-left">
                    <el-icon class="header-icon">
                      <Document />
                    </el-icon>
                    <span class="header-title">{{ $t('clusterEvents.title') }}</span>
                  </div>
                  <div class="header-actions">
                    <el-button :icon="Refresh" text @click="refreshEvents" class="refresh-btn">
                      {{ $t('clusterEvents.refresh') }}
                    </el-button>
                  </div>
                </div>
              </template>
              <div class="events-content">
                <div v-if="eventsLoading" class="loading-events">
                  <el-icon class="is-loading">
                    <Refresh />
                  </el-icon>
                  <span>{{ $t('common.loading') }}</span>
                </div>
                <div v-else-if="!hasEvents" class="no-events">
                  <el-icon class="no-events-icon">
                    <Document />
                  </el-icon>
                  <span class="no-events-text">{{ $t('clusterEvents.noEvents') }}</span>
                </div>
                <template v-else>
                  <div class="event-item" v-for="event in events" :key="event.id">
                    <div class="event-message">{{ event.message }}</div>
                    <el-tag :type="getEventTagType(event.type)" size="small">
                      {{ $t(`clusterEvents.eventTypes.${event.type.toLowerCase()}`) }}
                    </el-tag>
                  </div>
                  <div class="events-footer">
                    <el-button text size="small" @click="goToEventsPage" class="view-all-btn">{{ $t('clusterEvents.viewAll') }}</el-button>
                  </div>
                </template>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { DataAnalysis, Monitor, Document, Refresh } from '@element-plus/icons-vue'
import ResourceSummaryCard from '@/components/ClusterSummary/ResourceSummaryCard.vue'
import NodeMetricsCard from '@/components/ClusterSummary/NodeMetricsCard.vue'
import { useEvents } from '@/composables/useEvents'

const { t } = useI18n()
const router = useRouter()

// Check if components exist
const resourceCardExists = ref(true)
const nodeMetricsCardExists = ref(true)

// 使用事件组合式函数
const {
  events,
  loading: eventsLoading,
  hasEvents,
  getEventTagType,
  fetchRecentEvents,
  refreshEvents
} = useEvents()

// 自动刷新定时器
let refreshTimer: NodeJS.Timeout | null = null

// 组件挂载时初始化
onMounted(async () => {
  // 获取初始事件数据（只显示5条）
  await fetchRecentEvents(5)
  
  // 设置自动刷新（每30秒）
  refreshTimer = setInterval(() => {
    fetchRecentEvents(5)
  }, 30000)
})

// 跳转到事件页面
const goToEventsPage = () => {
  router.push('/cluster/events')
}

// 组件卸载时清理定时器
onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
})
</script>

<style scoped>
.cluster-overview {
  padding: 24px;
  background-color: var(--el-bg-color-page);
  box-sizing: border-box;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0;
}

.overview-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.top-section {
  flex: 0 0 auto;
}

.bottom-section {
  flex: 0 0 auto;
}

.bottom-section .el-row {
  align-items: stretch;
}

.bottom-section .el-col {
  display: flex;
}

.bottom-section .el-col > * {
  flex: 1;
}

.resource-placeholder,
.metrics-placeholder,
.events-placeholder {
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.placeholder-content {
  text-align: center;
  color: var(--el-text-color-regular);
}

.placeholder-icon {
  font-size: 48px;
  color: var(--el-color-primary);
  margin-bottom: 16px;
}

.placeholder-content h3 {
  font-size: 18px;
  font-weight: 500;
  margin: 0 0 8px 0;
  color: var(--el-text-color-primary);
}

.placeholder-content p {
  font-size: 14px;
  margin: 0;
}

/* 事件卡片样式 */
.events-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.events-card :deep(.el-card__body) {
  flex: 1;
  padding: 0;
  display: flex;
  flex-direction: column;
}

.card-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-icon {
  font-size: 16px;
  color: var(--el-color-primary);
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
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

.events-content {
  flex: 1;
  padding: 16px;
  display: flex;
  flex-direction: column;
}

.event-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.event-item:last-of-type {
  border-bottom: none;
}



.event-message {
  flex: 1;
  font-size: 13px;
  color: var(--el-text-color-primary);
  line-height: 1.4;
}

.events-footer {
  margin-top: auto;
  padding-top: 12px;
  border-top: 1px solid var(--el-border-color-lighter);
  text-align: center;
}

.no-events {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
  text-align: center;
  flex: 1;
}

.no-events-icon {
  font-size: 32px;
  color: var(--el-text-color-placeholder);
  margin-bottom: 12px;
}

.no-events-text {
  font-size: 14px;
  color: var(--el-text-color-regular);
}

.loading-events {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
  text-align: center;
  flex: 1;
  gap: 12px;
}

.loading-events .el-icon {
  font-size: 24px;
  color: var(--el-color-primary);
}

.loading-events span {
  font-size: 14px;
  color: var(--el-text-color-regular);
}

/* 查看所有事件按钮样式 */
.view-all-btn {
  font-weight: 700 !important;
  color: var(--el-color-primary) !important;
  font-size: 13px !important;
}

.view-all-btn:hover {
  color: var(--el-color-primary-light-3) !important;
}

/* Responsive Design */
@media (max-width: 768px) {
  .cluster-overview {
    padding: 16px 16px 60px 16px;
  }

  .page-title {
    font-size: 24px;
  }

  .page-header {
    margin-bottom: 20px;
  }

  .overview-content {
    gap: 16px;
  }

  /* 在移动端让两个卡片垂直排列 */
  .bottom-section .el-col {
    margin-bottom: 16px;
  }

  .bottom-section .el-col:last-child {
    margin-bottom: 0;
  }

  .bottom-section .el-row {
    flex-direction: column;
  }

  .bottom-section .el-col {
    width: 100% !important;
    max-width: 100% !important;
  }
}

@media (max-width: 480px) {
  .cluster-overview {
    padding: 12px 12px 50px 12px;
  }

  .page-title {
    font-size: 20px;
  }

  .page-header {
    margin-bottom: 16px;
  }

  .overview-content {
    gap: 12px;
  }
}
</style>