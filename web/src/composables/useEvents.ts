import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { getRecentEvents, type ClusterEvent } from '@/api/events'
import { ElMessage } from 'element-plus'

export function useEvents() {
  const { t } = useI18n()
  
  const events = ref<ClusterEvent[]>([])
  const loading = ref(false)
  const lastUpdateTime = ref<string>('')

  // 格式化事件时间
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

  // 获取事件标签类型
  const getEventTagType = (eventType: string) => {
    const typeMap: Record<string, 'success' | 'warning' | 'danger' | 'info'> = {
      Normal: 'success',
      Warning: 'warning',
      Error: 'danger'
    }
    return typeMap[eventType] || 'info'
  }

  // 更新最后更新时间
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

  // 获取最近事件
  const fetchRecentEvents = async (limit = 10) => {
    loading.value = true
    try {
      const response = await getRecentEvents(limit)
      // 后端返回格式: { code: 200, data: { events: [...], total: N }, message: "..." }
      events.value = (response as any).data?.events || []
      updateLastUpdateTime()
    } catch (error) {
      console.error('Failed to fetch recent events:', error)
      ElMessage.error(t('clusterEvents.fetchError') || 'Failed to fetch events')
      events.value = []
    } finally {
      loading.value = false
    }
  }

  // 刷新事件
  const refreshEvents = async () => {
    await fetchRecentEvents()
  }

  // 计算属性：是否有事件
  const hasEvents = computed(() => events.value.length > 0)

  return {
    events,
    loading,
    lastUpdateTime,
    hasEvents,
    formatEventTime,
    getEventTagType,
    fetchRecentEvents,
    refreshEvents
  }
}