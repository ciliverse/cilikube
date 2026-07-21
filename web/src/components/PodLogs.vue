<template>
  <div class="pod-logs">
    <div class="logs-header">
      <div class="logs-controls">
        <el-form :inline="true" size="small">
          <el-form-item :label="$t('podManagement.form.container')">
            <el-select v-model="selectedContainer" :placeholder="$t('podManagement.form.container')"
              style="width: 200px;" @change="reconnectLogs" :disabled="connecting">
              <el-option v-for="container in containers" :key="container.name" 
                :label="container.name" :value="container.name" />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('podManagement.form.tailLines')">
            <el-input-number v-model="tailLines" :min="10" :max="10000" 
              style="width: 120px;" @change="reconnectLogs" :disabled="connecting" />
          </el-form-item>
          <el-form-item>
            <el-checkbox v-model="timestamps" @change="reconnectLogs" :disabled="connecting">
              时间戳
            </el-checkbox>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="reconnectLogs" :loading="connecting" :disabled="!selectedContainer">
              {{ $t('common.refresh') }}
            </el-button>
            <el-button @click="clearLogs" :disabled="connecting">
              清空
            </el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="connection-status">
        <el-tag :type="connectionStatus.type" size="small">
          {{ connectionStatus.text }}
        </el-tag>
      </div>
    </div>
    
    <div class="logs-content-container" ref="logsContainerRef">
      <pre class="logs-content" v-loading="connecting">{{ logsContent || (connecting ? '正在连接日志流...' : '请选择容器并点击刷新') }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'

const { t } = useI18n()

// Props
interface Props {
  namespace: string
  podName: string
  containers: Array<{ name: string; image: string }>
}

const props = defineProps<Props>()

// 响应式状态
const selectedContainer = ref('')
const tailLines = ref(100)
const timestamps = ref(true)
const logsContent = ref('')
const connecting = ref(false)
const logsContainerRef = ref<HTMLDivElement | null>(null)

// WebSocket 连接
let logWebSocket: WebSocket | null = null
let shouldAutoScroll = ref(true)

// 连接状态
const connectionStatus = computed(() => {
  if (connecting.value) {
    return { type: 'warning', text: '连接中...' }
  }
  if (logWebSocket && logWebSocket.readyState === WebSocket.OPEN) {
    return { type: 'success', text: '已连接' }
  }
  return { type: 'info', text: '未连接' }
})

// 初始化容器选择
watch(() => props.containers, (newContainers) => {
  if (newContainers.length > 0 && !selectedContainer.value) {
    selectedContainer.value = newContainers[0].name
    nextTick(() => {
      connectLogs()
    })
  }
}, { immediate: true })

// 连接日志流
const connectLogs = async () => {
  if (!selectedContainer.value || !props.namespace || !props.podName) {
    return
  }

  // 关闭之前的连接
  if (logWebSocket) {
    logWebSocket.close()
    logWebSocket = null
  }

  connecting.value = true
  logsContent.value = ''

  try {
    // 构建 WebSocket URL
    const wsProtocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://'
    const wsHost = window.location.host
    const params = new URLSearchParams({
      container: selectedContainer.value,
      tailLines: tailLines.value.toString(),
      timestamps: timestamps.value.toString(),
      follow: 'true'
    })
    
    const wsUrl = `${wsProtocol}${wsHost}/api/v1/namespaces/${props.namespace}/pods/${props.podName}/logs?${params.toString()}`

    // 建立 WebSocket 连接
    logWebSocket = new WebSocket(wsUrl)

    logWebSocket.onopen = () => {
      console.log('日志 WebSocket 连接成功')
      connecting.value = false
    }

    logWebSocket.onmessage = (event) => {
      if (typeof event.data === 'string') {
        logsContent.value += event.data + '\n'
        // 自动滚动到底部
        if (shouldAutoScroll.value) {
          nextTick(() => {
            const container = logsContainerRef.value?.querySelector('.logs-content')
            if (container) {
              container.scrollTop = container.scrollHeight
            }
          })
        }
      }
    }

    logWebSocket.onerror = (error) => {
      console.error('日志 WebSocket 连接出错:', error)
      connecting.value = false
      logsContent.value += '\n# WebSocket 连接出错\n'
      ElMessage.error('日志连接出错，请重试')
    }

    logWebSocket.onclose = (event) => {
      console.log('日志 WebSocket 关闭:', event.code, event.reason)
      connecting.value = false
      if (event.code !== 1000) {
        logsContent.value += `\n# WebSocket 连接已断开: ${event.reason || event.code}\n`
      }
    }

  } catch (error: any) {
    console.error("建立日志连接失败:", error)
    connecting.value = false
    logsContent.value = `建立日志连接出错: ${error.message || '网络错误'}`
    ElMessage.error(`建立日志连接出错: ${error.message || '网络错误'}`)
  }
}

// 重新连接日志
const reconnectLogs = () => {
  connectLogs()
}

// 清空日志
const clearLogs = () => {
  logsContent.value = ''
}

// 监听滚动事件，判断是否需要自动滚动
const setupScrollListener = () => {
  nextTick(() => {
    const container = logsContainerRef.value?.querySelector('.logs-content')
    if (container) {
      container.addEventListener('scroll', () => {
        const { scrollTop, scrollHeight, clientHeight } = container
        // 判断是否滚动到最底部
        shouldAutoScroll.value = scrollTop + clientHeight >= scrollHeight - 10
      })
    }
  })
}

// 组件卸载时清理
onBeforeUnmount(() => {
  if (logWebSocket) {
    logWebSocket.close(1000, 'Component unmounted')
    logWebSocket = null
  }
})

// 暴露方法给父组件
defineExpose({
  connectLogs,
  clearLogs
})

// 设置滚动监听
setupScrollListener()
</script>

<style scoped>
.pod-logs {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.logs-controls {
  flex: 1;
}

.connection-status {
  margin-left: 16px;
}

.logs-content-container {
  flex: 1;
  overflow: hidden;
  background: #1e293b;
  border-radius: 0 0 8px 8px;
}

.logs-content {
  height: 100%;
  overflow: auto;
  padding: 16px;
  margin: 0;
  font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  line-height: 1.5;
  color: #e2e8f0;
  white-space: pre-wrap;
  word-break: break-all;
  background: #1e293b;
}

/* 滚动条样式 */
.logs-content::-webkit-scrollbar {
  width: 8px;
}

.logs-content::-webkit-scrollbar-track {
  background: #374151;
}

.logs-content::-webkit-scrollbar-thumb {
  background: #6b7280;
  border-radius: 4px;
}

.logs-content::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}
</style>