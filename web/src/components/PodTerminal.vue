<template>
  <div class="pod-terminal">
    <div class="terminal-header">
      <div class="terminal-controls">
        <el-form :inline="true" size="small">
          <el-form-item :label="$t('podManagement.form.container')">
            <el-select v-model="selectedContainer" :placeholder="$t('podManagement.form.container')"
              style="width: 200px;" :disabled="connected || connecting">
              <el-option v-for="container in containers" :key="container.name" 
                :label="container.name" :value="container.name" />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('podManagement.form.command')">
            <el-input v-model="command" style="width: 150px;" placeholder="sh" 
              :disabled="connected || connecting" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="connectTerminal" 
              :disabled="!selectedContainer || connected" :loading="connecting">
              {{ connected ? '已连接' : '连接' }}
            </el-button>
            <el-button @click="disconnectTerminal" v-if="connected" type="danger">
              断开连接
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
    
    <div class="terminal-container" ref="terminalContainerRef">
      <div class="terminal-placeholder" v-if="!terminalReady">
        {{ connecting ? '正在连接终端...' : '请选择容器并点击连接按钮' }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onBeforeUnmount, onMounted } from 'vue'
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
const command = ref('sh')
const connecting = ref(false)
const connected = ref(false)
const terminalReady = ref(false)
const terminalContainerRef = ref<HTMLDivElement | null>(null)

// 终端相关变量
let terminal: any = null
let fitAddon: any = null
let execWebSocket: WebSocket | null = null

// 连接状态
const connectionStatus = computed(() => {
  if (connecting.value) {
    return { type: 'warning', text: '连接中...' }
  }
  if (connected.value) {
    return { type: 'success', text: '已连接' }
  }
  return { type: 'info', text: '未连接' }
})

// 初始化容器选择
watch(() => props.containers, (newContainers) => {
  if (newContainers.length > 0 && !selectedContainer.value) {
    selectedContainer.value = newContainers[0].name
  }
}, { immediate: true })

// 初始化终端
const initTerminal = async () => {
  if (!terminalContainerRef.value) {
    console.error('终端容器元素未找到')
    return false
  }

  // 如果已有终端，先清理
  if (terminal) {
    terminal.dispose()
    terminal = null
  }

  try {
    // 动态导入 xterm
    const { Terminal } = await import('xterm')
    const { FitAddon } = await import('xterm-addon-fit')
    
    terminal = new Terminal({
      cursorBlink: true,
      rows: 25,
      cols: 80,
      theme: { 
        background: '#1e293b', 
        foreground: '#e2e8f0',
        cursor: '#3b82f6',
        selection: '#3b82f6'
      },
      convertEol: true,
      fontFamily: 'Maple Mono, Monaco, Menlo, Ubuntu Mono, monospace',
      fontSize: 13,
      letterSpacing: 1
    })

    fitAddon = new FitAddon()
    terminal.loadAddon(fitAddon)
    terminal.open(terminalContainerRef.value)
    
    try {
      fitAddon.fit()
    } catch (e) {
      console.error("调整终端大小出错:", e)
    }

    // 监听用户输入
    terminal.onData((data: string) => {
      if (execWebSocket?.readyState === WebSocket.OPEN) {
        execWebSocket.send(data)
      }
    })

    // 监听窗口大小变化
    window.addEventListener('resize', () => {
      if (fitAddon) {
        try {
          fitAddon.fit()
        } catch (e) {
          console.error("调整终端大小出错:", e)
        }
      }
    })

    terminal.writeln('\r\n请点击 "连接" 按钮建立终端连接...')
    terminalReady.value = true
    console.log('终端初始化完成')
    return true
  } catch (error) {
    console.error('初始化终端失败:', error)
    ElMessage.error('初始化终端失败，请确保已安装 xterm 依赖')
    return false
  }
}

// 连接终端
const connectTerminal = async () => {
  if (!selectedContainer.value || !props.namespace || !props.podName) {
    ElMessage.warning('请选择容器')
    return
  }

  connecting.value = true

  try {
    // 确保终端已初始化
    if (!terminal) {
      const success = await initTerminal()
      if (!success) {
        connecting.value = false
        return
      }
    }

    // 构建 WebSocket URL
    const wsProtocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://'
    const wsHost = window.location.host
    const params = new URLSearchParams({
      container: selectedContainer.value,
      command: command.value || 'sh',
      tty: 'true',
      stdin: 'true',
      stdout: 'true',
      stderr: 'true'
    })
    
    const wsUrl = `${wsProtocol}${wsHost}/api/v1/namespaces/${props.namespace}/pods/${props.podName}/exec?${params.toString()}`

    // 建立 WebSocket 连接
    execWebSocket = new WebSocket(wsUrl)

    execWebSocket.onopen = () => {
      console.log('终端 WebSocket 连接成功')
      connecting.value = false
      connected.value = true
      terminal?.clear()
      terminal?.writeln('\r\n\x1b[32m终端连接成功！\x1b[0m')
      terminal?.focus()
      if (fitAddon) {
        try {
          fitAddon.fit()
        } catch (e) {
          console.error("调整终端大小出错:", e)
        }
      }
    }

    execWebSocket.onmessage = async (event) => {
      try {
        if (event.data instanceof ArrayBuffer) {
          terminal?.write(new Uint8Array(event.data))
        } else if (typeof event.data === 'string') {
          terminal?.write(event.data)
        } else if (event.data instanceof Blob) {
          const arrayBuffer = await event.data.arrayBuffer()
          terminal?.write(new Uint8Array(arrayBuffer))
        }
      } catch (error) {
        console.error('写入终端消息出错:', error)
      }
    }

    execWebSocket.onerror = (error) => {
      console.error('终端 WebSocket 连接出错:', error)
      connecting.value = false
      connected.value = false
      terminal?.writeln(`\r\n\x1b[31mWebSocket 连接错误\x1b[0m`)
      ElMessage.error('终端连接出错，请重试')
    }

    execWebSocket.onclose = (event) => {
      console.log('终端 WebSocket 关闭:', event.code, event.reason)
      connecting.value = false
      connected.value = false
      const reason = event.reason || `Code ${event.code}`
      terminal?.writeln(`\r\n\x1b[33m连接已关闭: ${reason}\x1b[0m`)
      execWebSocket = null
    }

  } catch (error: any) {
    console.error("建立终端连接失败:", error)
    connecting.value = false
    connected.value = false
    ElMessage.error(`建立终端连接出错: ${error.message || '网络错误'}`)
  }
}

// 断开终端连接
const disconnectTerminal = () => {
  if (execWebSocket) {
    execWebSocket.close(1000, 'User disconnected')
    execWebSocket = null
  }
  connected.value = false
  connecting.value = false
  terminal?.writeln('\r\n\x1b[33m用户主动断开连接\x1b[0m')
}

// 组件挂载时初始化终端
onMounted(() => {
  nextTick(() => {
    initTerminal()
  })
})

// 组件卸载时清理
onBeforeUnmount(() => {
  if (execWebSocket) {
    execWebSocket.close(1000, 'Component unmounted')
    execWebSocket = null
  }
  
  if (terminal) {
    terminal.dispose()
    terminal = null
  }
  
  if (fitAddon) {
    fitAddon = null
  }
})

// 暴露方法给父组件
defineExpose({
  connectTerminal,
  disconnectTerminal,
  initTerminal
})
</script>

<style scoped>
.pod-terminal {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.terminal-controls {
  flex: 1;
}

.connection-status {
  margin-left: 16px;
}

.terminal-container {
  flex: 1;
  background: #1e293b;
  border-radius: 0 0 8px 8px;
  position: relative;
  overflow: hidden;
}

.terminal-placeholder {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #64748b;
  font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  text-align: center;
}

/* xterm.js 样式覆盖 */
.terminal-container :deep(.xterm) {
  height: 100% !important;
  padding: 16px;
}

.terminal-container :deep(.xterm-viewport) {
  background: #1e293b !important;
}

.terminal-container :deep(.xterm-screen) {
  background: #1e293b !important;
}
</style>