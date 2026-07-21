<template>
  <div class="events-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('menu.clusterEvents') }}</h1>
      <div class="header-actions">
        <el-button :icon="Refresh" @click="refreshEvents" :loading="loading">
          {{ $t('clusterEvents.refresh') }}
        </el-button>
      </div>
    </div>

    <!-- 过滤器 -->
    <el-card class="filter-card" shadow="never">
      <div class="filters">
        <el-form :model="filters" inline>
          <el-form-item :label="$t('clusterEvents.filters.namespace')">
            <el-select
              v-model="filters.namespace"
              :placeholder="$t('clusterEvents.filters.allNamespaces')"
              clearable
              style="width: 200px"
              @change="handleFilterChange"
            >
              <el-option
                :label="$t('clusterEvents.filters.allNamespaces')"
                value=""
              />
              <el-option
                v-for="ns in namespaces"
                :key="ns"
                :label="ns"
                :value="ns"
              />
            </el-select>
          </el-form-item>
          
          <el-form-item :label="$t('clusterEvents.filters.type')">
            <el-select
              v-model="filters.type"
              :placeholder="$t('clusterEvents.filters.allTypes')"
              clearable
              style="width: 150px"
              @change="handleFilterChange"
            >
              <el-option label="Normal" value="Normal" />
              <el-option label="Warning" value="Warning" />
            </el-select>
          </el-form-item>

          <el-form-item :label="$t('clusterEvents.filters.limit')">
            <el-select
              v-model="filters.limit"
              style="width: 100px"
              @change="handleFilterChange"
            >
              <el-option label="50" :value="50" />
              <el-option label="100" :value="100" />
              <el-option label="200" :value="200" />
            </el-select>
          </el-form-item>
        </el-form>
      </div>
    </el-card>

    <!-- 事件列表 -->
    <el-card class="events-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clusterEvents.eventList') }}</span>
          <span class="event-count" v-if="!loading">
            {{ $t('clusterEvents.totalEvents', { count: totalEvents }) }}
          </span>
        </div>
      </template>

      <div v-loading="loading" class="events-content">
        <div v-if="!loading && events.length === 0" class="no-events">
          <el-empty :description="$t('clusterEvents.noEvents')" />
        </div>

        <div v-else class="events-list">
          <div
            v-for="event in events"
            :key="event.id"
            class="event-item"
          >
            <div class="event-header">
              <div class="event-info">
                <div class="event-object">
                  <el-tag size="small" type="info">{{ event.objectKind }}</el-tag>
                  <span class="object-name">{{ event.object || event.name }}</span>
                  <span class="namespace" v-if="event.namespace">
                    ({{ event.namespace }})
                  </span>
                </div>
                <div class="event-reason">{{ event.reason }}</div>
              </div>
              <div class="event-meta">
                <el-tag :type="getEventTagType(event.type)" size="small">
                  {{ event.type }}
                </el-tag>
                <span class="event-time">{{ formatEventTime(event.lastTime) }}</span>
              </div>
            </div>
            
            <div class="event-message">{{ event.message }}</div>
            
            <div class="event-details">
              <span class="detail-item">
                <el-icon><User /></el-icon>
                {{ event.source }}
              </span>
              <span class="detail-item" v-if="event.count > 1">
                <el-icon><RefreshRight /></el-icon>
                {{ $t('clusterEvents.count', { count: event.count }) }}
              </span>
              <span class="detail-item">
                <el-icon><Clock /></el-icon>
                {{ $t('clusterEvents.firstTime') }}: {{ formatAbsoluteTime(event.firstTime) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Refresh, User, RefreshRight, Clock } from '@element-plus/icons-vue'
import { getClusterEvents, getAllNamespaces, type ClusterEvent, type EventListParams } from '@/api/events'
import { ElMessage } from 'element-plus'

const { t } = useI18n()

// 响应式数据
const events = ref<ClusterEvent[]>([])
const loading = ref(false)
const totalEvents = ref(0)
const namespaces = ref<string[]>([])
const allNamespaces = ref<string[]>([]) // 存储所有命名空间

// 过滤器
const filters = reactive<EventListParams>({
  namespace: '',
  type: '',
  limit: 100
})

// 获取所有命名空间列表（不带过滤器）
const fetchAllNamespaces = async () => {
  try {
    // 获取所有事件来提取命名空间列表
    const response = await getClusterEvents({ limit: 500 }) // 获取更多事件来收集命名空间
    const nsSet = new Set(response.data?.events?.map(e => e.namespace).filter(Boolean) || [])
    allNamespaces.value = Array.from(nsSet).sort()
    
    // 更新命名空间选项（始终显示完整列表）
    namespaces.value = allNamespaces.value
  } catch (error) {
    console.error('Failed to fetch namespaces:', error)
    // 如果获取失败，使用默认命名空间
    allNamespaces.value = ['default', 'kube-system', 'kube-public']
    namespaces.value = allNamespaces.value
  }
}

// 获取事件列表
const fetchEvents = async () => {
  loading.value = true
  try {
    const params: EventListParams = {}
    if (filters.namespace) params.namespace = filters.namespace
    if (filters.type) params.type = filters.type
    if (filters.limit) params.limit = filters.limit

    const response = await getClusterEvents(params)
    events.value = response.data?.events || []
    totalEvents.value = response.data?.total || 0
  } catch (error) {
    console.error('Failed to fetch events:', error)
    ElMessage.error(t('clusterEvents.fetchError'))
    events.value = []
    totalEvents.value = 0
  } finally {
    loading.value = false
  }
}

// 刷新事件
const refreshEvents = () => {
  fetchEvents()
}

// 处理过滤器变化
const handleFilterChange = () => {
  fetchEvents()
}

// 格式化事件时间（相对时间）
const formatEventTime = (timestamp: string) => {
  const now = Date.now()
  const eventTime = new Date(timestamp).getTime()
  const diff = now - eventTime
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (minutes < 1) {
    return t('clusterEvents.timeAgo.justNow')
  } else if (minutes < 60) {
    return `${minutes} ${t('clusterEvents.timeAgo.minutesAgo')}`
  } else if (hours < 24) {
    return `${hours} ${t('clusterEvents.timeAgo.hoursAgo')}`
  } else {
    return `${days} ${t('clusterEvents.timeAgo.daysAgo')}`
  }
}

// 格式化绝对时间
const formatAbsoluteTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleString()
}

// 获取事件标签类型
const getEventTagType = (eventType: string) => {
  const typeMap: Record<string, 'success' | 'warning' | 'danger' | 'info'> = {
    Normal: 'success',
    Warning: 'warning',
    Error: 'danger'
  }
  return typeMap[eventType] || 'info'
}

// 自动刷新定时器
let refreshTimer: NodeJS.Timeout | null = null

onMounted(async () => {
  // 先获取命名空间列表
  await fetchAllNamespaces()
  // 然后获取事件
  await fetchEvents()
  // 设置自动刷新（每60秒）
  refreshTimer = setInterval(fetchEvents, 60000)
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
})
</script>

<style scoped>
.events-page {
  padding: 24px;
  background-color: var(--el-bg-color-page);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.filter-card {
  margin-bottom: 16px;
}

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.events-card {
  min-height: 400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.event-count {
  font-size: 14px;
  color: var(--el-text-color-regular);
}

.events-content {
  min-height: 300px;
}

.no-events {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}

.events-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.event-item {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  padding: 16px;
  background: var(--el-bg-color);
  transition: all 0.2s;
}

.event-item:hover {
  border-color: var(--el-color-primary);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.event-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.event-info {
  flex: 1;
}

.event-object {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.object-name {
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.namespace {
  color: var(--el-text-color-regular);
  font-size: 12px;
}

.event-reason {
  font-size: 14px;
  color: var(--el-color-primary);
  font-weight: 500;
}

.event-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.event-time {
  font-size: 12px;
  color: var(--el-text-color-regular);
}

.event-message {
  color: var(--el-text-color-primary);
  line-height: 1.5;
  margin-bottom: 12px;
  padding: 8px 12px;
  background: var(--el-fill-color-lighter);
  border-radius: 4px;
}

.event-details {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 12px;
  color: var(--el-text-color-regular);
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.detail-item .el-icon {
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .events-page {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .page-title {
    font-size: 24px;
  }

  .filters {
    flex-direction: column;
  }

  .event-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .event-meta {
    align-self: flex-end;
  }

  .event-details {
    flex-direction: column;
    gap: 8px;
  }
}
</style>