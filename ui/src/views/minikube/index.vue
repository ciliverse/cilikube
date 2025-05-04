<template>
  <el-card class="installer-card" shadow="hover">
    <template #header>
      <div class="card-header">
        <span>
          <el-icon><Platform /></el-icon>
          Minikube 安装器
        </span>
        <!-- 优化状态显示 -->
        <el-tooltip
            v-if="currentStep === 'idle' && !backendReachable && !isCheckingBackend"
            content="无法连接到后端服务。"
            placement="top"
          >
             <el-tag type="danger" effect="plain" size="small" round>
                <el-icon><CircleCloseFilled /></el-icon> 离线
             </el-tag>
        </el-tooltip>
         <el-tag v-if="currentStep === 'idle' && backendReachable && !isCheckingBackend" type="success" effect="plain" size="small" round>
             <el-icon><SuccessFilled /></el-icon> 在线
         </el-tag>
         <el-tag v-if="isCheckingBackend" type="info" effect="light" size="small" round>
             <el-icon class="is-loading"><Loading /></el-icon> 检查中
         </el-tag>
         <el-tag v-if="currentStep === 'connecting'" type="warning" effect="light" size="small" round>
             <el-icon class="is-loading"><Loading /></el-icon> 连接中
         </el-tag>
         <el-tag v-if="currentStep === 'download'" type="primary" effect="light" size="small" round>
             <el-icon><Download /></el-icon> 下载中
         </el-tag>
         <!-- 移除 "(模拟)" -->
         <el-tag v-if="currentStep === 'install'" color="#E6A23C" effect="light" size="small" round>
              <el-icon><Setting /></el-icon> 安装
         </el-tag>
         <el-tag v-if="currentStep === 'start'" color="#67C23A" effect="light" size="small" round>
             <el-icon class="is-loading"><Loading /></el-icon> 启动中
         </el-tag>
         <el-tag v-if="currentStep === 'finished'" type="success" effect="dark" size="small" round>
              <el-icon><Select /></el-icon> 成功
         </el-tag>
         <el-tag v-if="currentStep === 'error'" type="danger" effect="dark" size="small" round>
             <el-icon><CloseBold /></el-icon> 失败
         </el-tag>
      </div>
    </template>

    <!-- Initial State -->
    <div v-if="currentStep === 'idle'" class="installer-content initial-state">
        <div v-if="backendReachable">
            <p>此工具将尝试通过后端服务安装并启动 Minikube。</p>
            <el-steps direction="vertical" :active="3" finish-status="process" style="margin-top: 15px; max-width: 400px;">
                <el-step title="下载 Minikube" icon="Download" description="从官方源获取最新版本。" />
                <!-- 移除 "(模拟)" -->
                <el-step title="安装" icon="Setting" description="安装到系统路径 (可能需要免密 sudo)。" />
                <el-step title="启动 Minikube" icon="VideoPlay" description="使用 'minikube start --force' 启动。" />
            </el-steps>
        </div>
        <div v-else class="offline-state">
             <el-alert title="后端服务离线" type="error" :description="`无法连接到后端服务 (${backendBaseUrl})，请确保后端服务正在运行并检查网络连接。`" :closable="false" show-icon center/>
              <div style="text-align: center; margin-top: 15px;">
                <!-- 添加重试连接按钮 -->
                <el-button type="primary" @click="checkBackend" :icon="Refresh" :loading="isCheckingBackend">
                    {{ isCheckingBackend ? '正在重试...' : '重试连接' }}
                </el-button>
             </div>
        </div>
    </div>

    <!-- Installation Progress State (Connecting, Downloading, Installing, Starting) -->
    <div v-if="isInstalling || currentStep === 'connecting'" class="installer-content progress-state">
       <!-- Step Indicator -->
       <el-steps :active="activeStepIndex" :process-status="stepStatus" :finish-status="finishStatus" align-center style="margin-bottom: 25px;">
            <el-step title="下载" :icon="Download" />
            <!-- 移除 "(模拟)" -->
            <el-step title="安装" :icon="Setting" />
            <el-step title="启动" :icon="VideoPlay" />
       </el-steps>

       <!-- **单一进度条** -->
       <el-progress
            v-if="currentStep !== 'connecting'"
            :text-inside="false"
            :stroke-width="14"
            :percentage="overallProgress"
            :status="overallProgressStatus"
            striped
            striped-flow
            :duration="10"
            style="margin: 0 10px 15px 10px;"
          >
             <span>总进度: {{ overallProgress }}%</span>
          </el-progress>
        <div v-if="currentStep === 'connecting'" style="text-align: center; margin-bottom: 15px; color: #909399; font-size: 0.9em;">
           <el-icon class="is-loading"><Loading /></el-icon> 正在建立连接...
        </div>

        <!-- Status Message & Log -->
        <div class="status-box">
            <p v-if="statusMessage">
              <el-text :type="messageType" class="status-message-text">
                  <el-icon v-if="messageType === 'warning'" style="margin-right: 3px;"><WarningFilled /></el-icon>
                  <strong>{{ stepTitle }}:</strong> {{ statusMessage }}
              </el-text>
            </p>
             <!-- 注意: 如果后端仍然发送模拟警告，这里会显示 -->
            <el-alert v-if="currentStep === 'install' && showInstallWarning" title="安装提示" type="warning" :description="installWarningMessage" :closable="false" style="margin-top: 10px; font-size: 0.85em;" />
        </div>

         <el-collapse v-model="activeLog" style="margin-top: 15px;">
            <el-collapse-item name="1">
                 <template #title>
                    <el-icon><Memo /></el-icon> <span style="margin-left: 5px;">详细日志</span>
                    <!-- 错误和警告徽章 -->
                    <el-badge :value="logErrorCount" type="danger" :hidden="logErrorCount === 0" style="margin-left: 10px;"/>
                    <el-badge :value="logWarningCount" type="warning" :hidden="logWarningCount === 0" style="margin-left: 10px;" />
                 </template>
                 <!-- 使用 v-html 渲染带样式的日志 -->
                <pre class="raw-output" ref="logContainer" v-html="filteredLogOutput"></pre>
            </el-collapse-item>
        </el-collapse>
    </div>

    <!-- Success State -->
    <div v-if="installComplete && !error" class="installer-content success-state">
       <el-result icon="success" title="Minikube 安装/启动成功" :sub-title="statusMessage">
        <template #extra>
          <el-button type="primary" @click="resetState" :icon="RefreshLeft">重新开始</el-button>
           <el-button @click="toggleSuccessLog" :icon="Document">
             {{ showSuccessLog ? '隐藏日志' : '显示日志' }}
           </el-button>
        </template>
      </el-result>
       <el-collapse-transition>
         <div v-show="showSuccessLog">
              <!-- 使用 v-html -->
             <pre class="raw-output success-log" v-html="filteredLogOutput"></pre>
         </div>
       </el-collapse-transition>
    </div>

    <!-- Error State -->
     <div v-if="error && currentStep === 'error'" class="installer-content error-state">
       <el-steps :active="activeStepIndex" :process-status="stepStatus" :finish-status="finishStatus" align-center style="margin-bottom: 25px;">
            <el-step title="下载" :icon="Download" />
            <!-- 移除 "(模拟)" -->
            <el-step title="安装" :icon="Setting" />
            <el-step title="启动" :icon="VideoPlay" />
       </el-steps>
       <el-result icon="error" :title="`在步骤 '${stepTitleOnError}' 发生错误`" :sub-title="error">
         <template #extra>
            <el-button type="primary" @click="resetState" :icon="Refresh">重试</el-button>
            <el-button @click="toggleErrorLog" :icon="Document">
              {{ showErrorLog ? '隐藏日志' : '显示日志' }}
            </el-button>
         </template>
       </el-result>
        <el-collapse-transition>
           <div v-show="showErrorLog">
                <!-- 使用 v-html -->
             <pre class="raw-output error-log" v-html="filteredLogOutput"></pre>
           </div>
        </el-collapse-transition>
    </div>

    <!-- Footer Actions -->
    <template #footer>
      <div class="card-footer">
          <el-button
             v-if="currentStep === 'idle'"
             type="primary"
             :icon="CaretRight"
             @click="startInstallation"
             :loading="isCheckingBackend"
             :disabled="isCheckingBackend || !backendReachable"
           >
             {{ isCheckingBackend ? '检查连接...' : '开始安装' }}
           </el-button>
           <el-button
              v-if="isInstalling || currentStep === 'connecting'"
              type="danger"
              :icon="CircleClose"
              @click="cancelInstallation"
              plain
           >
             取消
           </el-button>
           <el-button
              v-if="currentStep === 'error'"
              type="primary"
              :icon="Refresh"
              @click="resetState"
           >
             重试
           </el-button>
            <el-button
              v-if="currentStep === 'finished'"
              type="success"
              :icon="Check"
              plain
              disabled
           >
             已完成
           </el-button>
      </div>
    </template>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue';
import { ElNotification, ElMessage, ElProgress, ElAlert, ElCard, ElButton, ElIcon, ElTag, ElText, ElCollapse, ElCollapseItem, ElResult, ElCollapseTransition, ElTooltip, ElSteps, ElStep, ElBadge } from 'element-plus';
import {
    Platform, CaretRight, CircleClose, Refresh, Check, Download, Warning, VideoPlay,
    WarningFilled, Loading, Select, CloseBold, Memo, Document, RefreshLeft,
    SuccessFilled, CircleCloseFilled, Setting // 确保 Setting 已导入
} from '@element-plus/icons-vue';

// --- Types ---
type StepType = 'idle' | 'connecting' | 'download' | 'install' | 'start' | 'finished' | 'error';
const VITE_API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://192.168.1.100:8080";
// --- Configuration ---
const backendBaseUrl = ref<string>(VITE_API_BASE_URL);
const healthCheckUrl = computed(() => `${backendBaseUrl.value}/healthz`);
const installUrl = computed(() => `${backendBaseUrl.value}/api/v1/system/install-minikube`);

// --- State ---
const currentStep = ref<StepType>('idle');
const overallProgress = ref<number>(0);
// 移除 stepProgress
// const stepProgress = ref<number>(0);
const statusMessage = ref<string>('正在检查后端状态...');
const error = ref<string | null>(null);
const eventSource = ref<EventSource | null>(null);
const rawOutputLog = ref<string>('');
const backendReachable = ref<boolean>(false);
const isCheckingBackend = ref<boolean>(true);
const activeLog = ref<string[]>([]);
const showSuccessLog = ref<boolean>(false);
const showErrorLog = ref<boolean>(false);
const showInstallWarning = ref<boolean>(false);
const installWarningMessage = ref<string>('');
const stepIndexOnError = ref<number | null>(null);
const parsingErrorCount = ref<number>(0);
const logErrorCount = ref<number>(0);
const logWarningCount = ref<number>(0); // 添加警告计数
const logContainer = ref<HTMLPreElement | null>(null);

// --- Computed ---
const isInstalling = computed(() => ['download', 'install', 'start'].includes(currentStep.value));
const installComplete = computed(() => currentStep.value === 'finished');

// ** 优化：过滤并高亮日志 **
const filteredLogOutput = computed(() => {
    const lines = rawOutputLog.value.split('\n');
    const filteredLines = [];
    for (const line of lines) {
        const trimmed = line.trim();
        // 过滤掉空行和看起来像 curl 进度条的行
        if (trimmed === '') continue;
        if (trimmed.startsWith('#') && (trimmed.includes('%') || trimmed.includes('O'))) continue;
        if (/^[#=O\- \t\r]+$/.test(trimmed)) continue; // 过滤掉只包含这些字符的行
        if (/^\d+\.\d+%$/.test(trimmed)) continue; // 过滤掉纯百分比的行 (来自curl stderr)

        // 添加样式标记
        let styledLine = line; // 保留原始缩进和时间戳
        if (trimmed.includes('[ERROR]')) {
             styledLine = `<span class="log-error">${line}</span>`;
        } else if (trimmed.includes('[WARN]')) {
             styledLine = `<span class="log-warning">${line}</span>`;
        }
        filteredLines.push(styledLine);
    }
    return filteredLines.join('\n');
});

const activeStepIndex = computed(() => { /* ... (保持不变) ... */
    switch (currentStep.value) { case 'connecting': return 0; case 'download': return 0; case 'install': return 1; case 'start': return 2; case 'finished': return 3; case 'error': return stepIndexOnError.value !== null ? stepIndexOnError.value : 0; default: return -1; } });
const stepStatus = computed(() => { /* ... (保持不变) ... */
    if (currentStep.value === 'error') return 'error'; if (currentStep.value === 'connecting' || isInstalling.value) return 'process'; return 'wait'; });
const finishStatus = computed(() => { /* ... (保持不变) ... */
    if (currentStep.value === 'finished') return 'success'; if (currentStep.value === 'error') return 'error'; return 'process'; });
const overallProgressStatus = computed(() => { /* ... (保持不变) ... */
  if (currentStep.value === 'error') return 'exception'; if (currentStep.value === 'finished') return 'success'; return undefined; });

// ** 修改：移除 "(模拟)" **
const stepTitle = computed(() => {
     switch (currentStep.value) {
        case 'connecting': return '连接中';
        case 'download': return '下载';
        case 'install': return '安装'; // 修改这里
        case 'start': return '启动';
        case 'finished': return '完成';
        case 'error': return '错误';
        default: return '状态';
    }
});
// ** 修改：移除 "(模拟)" **
const stepTitleOnError = computed(() => {
     switch (stepIndexOnError.value) {
        case 0: return '下载';
        case 1: return '安装'; // 修改这里
        case 2: return '启动';
        default: return '未知步骤';
     }
});

const messageType = computed(() => { /* ... (保持不变) ... */
    switch (currentStep.value) { case 'error': return 'danger'; case 'finished': return 'success'; case 'install': return 'warning'; default: return 'primary'; } });

// --- Methods ---
const resetState = () => { /* ... (移除 stepProgress 重置) ... */
    console.log('Resetting state...'); currentStep.value = 'idle'; overallProgress.value = 0; /*stepProgress.value = 0;*/ statusMessage.value = '正在检查后端状态...'; error.value = null; rawOutputLog.value = ''; activeLog.value = []; showSuccessLog.value = false; showErrorLog.value = false; showInstallWarning.value = false; installWarningMessage.value = ''; stepIndexOnError.value = null; parsingErrorCount.value = 0; logErrorCount.value = 0; logWarningCount.value = 0; closeEventSource(); checkBackend(); };
const closeEventSource = () => { /* ... (保持不变) ... */
    if (eventSource.value) { eventSource.value.close(); eventSource.value = null; console.log('EventSource connection explicitly closed.'); } };
const cancelInstallation = () => { /* ... (保持不变) ... */
    if (isInstalling.value || currentStep.value === 'connecting') { ElMessage.warning('已手动取消连接。后端过程可能仍在继续。'); const lastStep = currentStep.value; error.value = "用户手动取消"; currentStep.value = 'error'; stepIndexOnError.value = ['connecting', 'download'].includes(lastStep) ? 0 : lastStep === 'install' ? 1 : lastStep === 'start' ? 2 : activeStepIndex.value; closeEventSource(); } };

// ** 修改：appendLog 添加错误/警告计数 **
const appendLog = (line: string | undefined) => {
    if (line && typeof line === 'string') {
        const timestamp = new Date().toLocaleTimeString();
        const trimmedLine = line.trim();
        if (trimmedLine === '') return; // 跳过完全空行

        // 标记错误和警告
        let logType = '';
        let processedLine = trimmedLine; // 用于检查关键字
        if (trimmedLine.toLowerCase().includes('[error]') || trimmedLine.toLowerCase().includes('fail') || trimmedLine.toLowerCase().startsWith('error:')) {
             logType = 'error';
             logErrorCount.value++;
        } else if (trimmedLine.toLowerCase().includes('[warn]') || trimmedLine.toLowerCase().includes('warning') || trimmedLine.startsWith('!')) {
             logType = 'warn';
             logWarningCount.value++;
             // 移除 Minikube 的 '!' 前缀，如果存在
             if (trimmedLine.startsWith('!')) {
                 processedLine = trimmedLine.substring(1).trim();
             }
        }

        // 添加时间戳和标记到原始日志记录
        rawOutputLog.value += `[${timestamp}]${logType ? ` [${logType.toUpperCase()}]` : ''} ${processedLine}\n`;

        // 滚动逻辑
        nextTick(() => { if (logContainer.value) { logContainer.value.scrollTo({ top: logContainer.value.scrollHeight, behavior: 'smooth' }); } });
    }
};

const startInstallation = () => { /* ... (移除 stepProgress 重置) ... */
    if (!backendReachable.value) { ElMessage.error("无法连接到后端服务，请检查。"); return; }
    resetState(); currentStep.value = 'connecting'; statusMessage.value = '正在建立 SSE 连接...'; rawOutputLog.value = '--- 请求安装 ---\n'; activeLog.value = ['1']; parsingErrorCount.value = 0; logErrorCount.value = 0; logWarningCount.value = 0; closeEventSource();
    console.log(`Connecting to SSE endpoint: ${installUrl.value}`);
    eventSource.value = new EventSource(installUrl.value, { withCredentials: false });

    eventSource.value.addEventListener('message', (event: MessageEvent) => { /* ... (移除 stepProgress 更新) ... */
        console.log("Raw SSE data received:", event.data); try { const data = JSON.parse(event.data); console.log("Parsed SSE data:", data); parsingErrorCount.value = 0; if (data.step) { currentStep.value = data.step as StepType; } if (typeof data.progress === 'number' && data.progress >= 0) { overallProgress.value = Math.min(data.progress, 100); } /* 移除 stepProgress 更新 */ if (typeof data.message === 'string') { if (!data.rawLine || data.message !== data.rawLine.trim()) { statusMessage.value = data.message; } else { statusMessage.value = ''; } if (currentStep.value === 'install' && data.message.includes("模拟")) { showInstallWarning.value = true; installWarningMessage.value = data.message; } else if (currentStep.value !== 'install') { showInstallWarning.value = false; } } if (typeof data.rawLine === 'string') { appendLog(data.rawLine); } if (data.error) { console.error('Backend reported error:', data.error); error.value = data.error; currentStep.value = 'error'; stepIndexOnError.value = activeStepIndex.value; ElNotification({ title: '安装失败', message: error.value, type: 'error', duration: 0 }); closeEventSource(); } else if (data.done && data.step === 'finished') { console.log('Backend reported finished.'); currentStep.value = 'finished'; overallProgress.value = 100; ElNotification({ title: '安装成功', message: statusMessage.value || '操作成功完成！', type: 'success', duration: 5000 }); closeEventSource(); } else if (data.done && data.step !== 'finished') { console.warn("Backend reported done=true but step is not 'finished':", data.step); if (!data.error) { error.value = `流程意外终止于步骤 '${data.step}'`; currentStep.value = 'error'; stepIndexOnError.value = activeStepIndex.value; ElNotification({ title: '流程错误', message: error.value, type: 'error', duration: 0 }); closeEventSource(); } } } catch (e) { parsingErrorCount.value++; console.error(`解析 SSE 消息 #${parsingErrorCount.value} 时出错:`, e, "原始数据:", event.data); appendLog(`*** 前端错误: 无法解析消息 (第 ${parsingErrorCount.value} 次): ${event.data} - 错误: ${e} ***`); ElMessage.warning(`收到格式错误的消息 (错误 #${parsingErrorCount.value})，已跳过。`); logErrorCount.value++; } });
    eventSource.value.onopen = () => { /* ... (保持不变) ... */ console.log('SSE connection opened successfully.'); statusMessage.value = '已连接，等待后端开始...'; appendLog('--- SSE 连接已建立 ---'); };
    eventSource.value.onerror = (event) => { /* ... (保持不变) ... */ console.error('EventSource 连接失败:', event); let readyState = event.target ? (event.target as EventSource).readyState : 'unknown'; console.error(`EventSource readyState: ${readyState} (0=CONNECTING, 1=OPEN, 2=CLOSED)`); appendLog(`--- SSE 连接错误 (State: ${readyState}) ---`); let errorMessage = `无法连接到后端 (${installUrl.value}) 或连接意外断开。`; if (currentStep.value === 'connecting' || isInstalling.value) { error.value = errorMessage; statusMessage.value = "连接错误"; currentStep.value = 'error'; stepIndexOnError.value = activeStepIndex.value > 0 ? activeStepIndex.value : 0; ElNotification({ title: '连接错误', message: errorMessage, type: 'error', duration: 0 }); } else { console.warn("SSE error occurred outside of active installation/connection phase."); if(currentStep.value === 'idle') { statusMessage.value = '与后端连接失败。'; backendReachable.value = false; } } closeEventSource(); };
};

// --- Backend Check & Lifecycle (保持不变) ---
const checkBackend = async () => { /* ... (保持不变) ... */
    console.log(`[Debug] 正在尝试连接健康检查 URL: ${healthCheckUrl.value}`); isCheckingBackend.value = true; try { const response = await fetch(healthCheckUrl.value, { method: 'GET', mode: 'cors', signal: AbortSignal.timeout(5000) }); backendReachable.value = response.ok; if (!response.ok) console.warn(`Backend health check non-OK status: ${response.status}`); else console.log("Backend health check successful."); } catch (err: any) { console.error(`[Debug] Fetch 请求失败，错误对象:`, err); console.error('Backend health check failed:', err.name === 'TimeoutError' ? 'Request timed out' : err); backendReachable.value = false; } finally { isCheckingBackend.value = false; if (currentStep.value === 'idle') { statusMessage.value = backendReachable.value ? '后端服务在线，准备就绪。' : '后端服务离线。'; } } };
let healthCheckInterval: number | null = null;
onMounted(() => { checkBackend(); healthCheckInterval = window.setInterval(checkBackend, 30000); });
onUnmounted(() => { closeEventSource(); if (healthCheckInterval) clearInterval(healthCheckInterval); });

// --- Log Toggling (保持不变) ---
const toggleSuccessLog = () => showSuccessLog.value = !showSuccessLog.value;
const toggleErrorLog = () => showErrorLog.value = !showErrorLog.value;

</script>

<style scoped>
.installer-card { margin-bottom: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; font-weight: bold; }
.card-header .el-icon { margin-right: 8px; vertical-align: -1px; }
.card-header .el-tag { margin-left: 10px; display: inline-flex; align-items: center; }
.card-header .el-tag .el-icon { margin-right: 4px; }

.installer-content { padding: 20px 15px; min-height: 150px; } /* Adjusted padding */

.initial-state .el-steps { margin-left: 10px; margin-top: 20px;} /* Added margin top */
.initial-state .offline-state .el-alert { margin-top: 20px;}
.initial-state p { margin-bottom: 15px; } /* More space */

.progress-state .status-box {
    text-align: center;
    margin-top: 10px; /* Reduced margin */
    margin-bottom: 20px; /* Added bottom margin */
    padding: 12px 15px;
    background-color: var(--el-fill-color-lighter);
    border-radius: 4px;
    border: 1px solid var(--el-border-color-lighter);
}
.progress-state .status-message-text { font-size: 1em; display: inline-flex; align-items: center; }

/* Style for the single progress bar text */
.el-progress--line { width: calc(100% - 20px); margin-left: 10px; margin-right: 10px; } /* Give some horizontal margin */
.el-progress__text { font-size: 12px !important; color: var(--el-text-color-primary); }

.raw-output {
    background-color: var(--el-bg-color-page);
    border: 1px solid var(--el-border-color);
    border-radius: 4px;
    padding: 10px 15px;
    max-height: 400px;
    overflow-y: auto;
    font-family: 'Consolas', 'Monaco', 'Menlo', 'Courier New', monospace;
    font-size: 12px;
    line-height: 1.6;
    color: var(--el-text-color-regular);
    margin-top: 10px;
    white-space: pre-wrap;
    word-break: break-word;
}
/* Styles for log highlighting */
.raw-output .log-error { color: var(--el-color-danger); font-weight: bold; }
.raw-output .log-warning { color: var(--el-color-warning); }

.success-state .raw-output,
.error-state .raw-output { margin-top: 10px; }
.success-log { border-color: var(--el-color-success-light-3); background-color: var(--el-color-success-light-9); }
.error-log { border-color: var(--el-color-danger-light-3); background-color: var(--el-color-danger-light-9); }

.card-footer { display: flex; justify-content: flex-end; align-items: center; padding-top: 15px; }
.el-result { padding: 15px 0; }
:deep(.el-step.is-vertical .el-step__line) { background-color: var(--el-border-color-lighter); }
</style>