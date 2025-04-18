<template>
    <el-card class="installer-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>
            <el-icon><Platform /></el-icon>
            Minikube 安装器 (多步骤)
          </span>
           <!-- Status Tags -->
          <el-tag v-if="currentStep === 'idle'" type="info" effect="light" size="small">空闲</el-tag>
          <el-tag v-if="isCheckingBackend" type="info" effect="light" size="small">检查中...</el-tag>
          <el-tooltip
            v-if="!backendReachable && !isCheckingBackend && currentStep === 'idle'"
            content="无法连接到后端服务，请检查其是否运行。"
            placement="top"
          >
             <el-tag type="danger" effect="dark" size="small">离线</el-tag>
          </el-tooltip>
          <el-tag v-if="currentStep === 'connecting'" type="warning" effect="light" size="small">连接中</el-tag>
          <el-tag v-if="currentStep === 'download'" type="primary" effect="light" size="small">下载中</el-tag>
          <el-tag v-if="currentStep === 'install'" color="#E6A23C" effect="light" size="small">安装 (模拟)</el-tag>
          <el-tag v-if="currentStep === 'start'" color="#67C23A" effect="light" size="small">启动中</el-tag>
          <el-tag v-if="currentStep === 'finished'" type="success" effect="dark" size="small">成功</el-tag>
          <el-tag v-if="currentStep === 'error'" type="danger" effect="dark" size="small">失败</el-tag>
        </div>
      </template>
  
      <!-- Initial State -->
      <div v-if="currentStep === 'idle'" class="installer-content initial-state">
          <p>此工具将按顺序执行以下步骤来安装/启动 Minikube:</p>
           <ol>
              <li><el-icon><Download /></el-icon> 从 GitHub 下载最新的 Minikube 二进制文件。</li>
              <li><el-icon><WarningFilled color="#E6A23C" /></el-icon> <el-text type="warning">模拟安装步骤。</el-text> (需要手动执行 <code>sudo install</code> 或预配置)</li>
              <li><el-icon><VideoPlay /></el-icon> 执行 <code>minikube start --force</code> 启动集群。</li>
           </ol>
          <p><el-text type="info" size="small">点击下方按钮开始。</el-text></p>
           <el-alert v-if="!backendReachable" title="后端连接失败" type="warning" description="请确保后端服务正在运行并且网络可达。" :closable="false" show-icon style="margin-top: 15px;"/>
      </div>
  
      <!-- Installation Progress State (Connecting, Downloading, Installing, Starting) -->
      <div v-if="isInstalling || currentStep === 'connecting'" class="installer-content progress-state">
         <!-- Step Indicator -->
         <el-steps :active="activeStepIndex" :process-status="stepStatus" :finish-status="finishStatus" align-center style="margin-bottom: 25px;">
              <el-step title="下载" :icon="Download" />
              <el-step title="安装 (模拟)" :icon="Warning" />
              <el-step title="启动" :icon="VideoPlay" />
         </el-steps>
  
         <!-- Overall Progress -->
         <el-progress
              v-if="currentStep !== 'connecting'"
              :text-inside="true"
              :stroke-width="18"
              :percentage="overallProgress"
              :status="overallProgressStatus"
              striped
              striped-flow
              :duration="10"
              style="margin-bottom: 15px;"
            >
             <span>总进度: {{ overallProgress }}%</span>
            </el-progress>
          <div v-else style="text-align: center; margin-bottom: 15px; color: #909399;">
             <el-icon class="is-loading"><Loading /></el-icon> 正在建立连接...
          </div>
  
  
          <!-- Status Message & Log -->
          <div class="status-box">
              <p>
                <el-text :type="messageType">
                    <strong>{{ stepTitle }}:</strong> {{ statusMessage }}
                </el-text>
              </p>
              <el-alert v-if="showInstallWarning" title="模拟安装提示" type="warning" :description="installWarningMessage" show-icon :closable="false" style="margin-top: 10px;"/>
          </div>
  
           <el-collapse v-model="activeLog" style="margin-top: 15px;">
              <el-collapse-item title="显示/隐藏详细日志" name="1">
                  <pre class="raw-output">{{ rawOutputLog }}</pre>
              </el-collapse-item>
          </el-collapse>
      </div>
  
      <!-- Success State -->
      <div v-if="installComplete && !error" class="installer-content success-state">
         <el-result
            icon="success"
            title="Minikube 安装/启动成功"
            :sub-title="statusMessage"
          >
          <template #extra>
            <el-button type="primary" @click="resetState">返回初始状态</el-button>
             <el-button @click="toggleSuccessLog">
               {{ showSuccessLog ? '隐藏日志' : '显示日志' }}
             </el-button>
          </template>
        </el-result>
         <el-collapse-transition>
           <div v-show="showSuccessLog">
               <pre class="raw-output success-log">{{ rawOutputLog }}</pre>
           </div>
         </el-collapse-transition>
      </div>
  
      <!-- Error State -->
       <div v-if="error && currentStep === 'error'" class="installer-content error-state">
         <el-steps :active="activeStepIndex" :process-status="stepStatus" :finish-status="finishStatus" align-center style="margin-bottom: 25px;">
              <el-step title="下载" :icon="Download" />
              <el-step title="安装 (模拟)" :icon="Warning" />
              <el-step title="启动" :icon="VideoPlay" />
         </el-steps>
  
         <el-result
            icon="error"
            :title="`在步骤 '${stepTitleOnError}' 发生错误`"
            :sub-title="error"
          >
           <template #extra>
              <el-button type="primary" @click="resetState">重试</el-button>
              <el-button @click="toggleErrorLog">
                {{ showErrorLog ? '隐藏日志' : '显示日志' }}
              </el-button>
           </template>
         </el-result>
          <el-collapse-transition>
             <div v-show="showErrorLog">
               <pre class="raw-output error-log">{{ rawOutputLog }}</pre>
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
               {{ isCheckingBackend ? '检查连接...' : '开始' }}
             </el-button>
             <el-button
                v-if="isInstalling || currentStep === 'connecting'"
                type="danger"
                :icon="CircleClose"
                @click="cancelInstallation"
                plain
             >
               取消连接
             </el-button>
             <el-button
                v-if="currentStep === 'error'"
                type="warning"
                :icon="Refresh"
                @click="resetState"
                plain
             >
               重置状态
             </el-button>
              <el-button
                v-if="currentStep === 'finished'"
                type="success"
                :icon="Check"
                plain
                disabled
             >
               完成
             </el-button>
        </div>
      </template>
  
    </el-card>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
  import { ElNotification, ElMessage, ElProgress, ElAlert, ElCard, ElButton, ElIcon, ElTag, ElText, ElCollapse, ElCollapseItem, ElResult, ElCollapseTransition, ElTooltip, ElSteps, ElStep } from 'element-plus';
  import { Platform, CaretRight, CircleClose, Refresh, Check, Download, Warning, VideoPlay, WarningFilled, Loading } from '@element-plus/icons-vue'; // Import Loading icon
  
  // --- Types ---
  type StepType = 'idle' | 'connecting' | 'download' | 'install' | 'start' | 'finished' | 'error'; // Added 'connecting'
  
  // --- Configuration ---
  const backendBaseUrl = ref<string>('http://192.168.1.100:8080'); // Hardcoded as requested
  const healthCheckUrl = computed(() => `${backendBaseUrl.value}/health`);
  const installUrl = computed(() => `${backendBaseUrl.value}/api/v1/system/install-minikube`);
  
  // --- State ---
  const currentStep = ref<StepType>('idle');
  const overallProgress = ref<number>(0);
  const statusMessage = ref<string>('等待连接后端...');
  const error = ref<string | null>(null);
  const eventSource = ref<EventSource | null>(null);
  const rawOutputLog = ref<string>('');
  const backendReachable = ref<boolean>(false);
  const isCheckingBackend = ref<boolean>(false);
  const activeLog = ref<string[]>([]);
  const showSuccessLog = ref<boolean>(false);
  const showErrorLog = ref<boolean>(false);
  const showInstallWarning = ref<boolean>(false);
  const installWarningMessage = ref<string>('');
  const stepIndexOnError = ref<number | null>(null); // Store step index when error occurs
  const parsingErrorCount = ref<number>(0); // Count parsing errors
  
  // --- Computed ---
  const isInstalling = computed(() => ['download', 'install', 'start'].includes(currentStep.value));
  const installComplete = computed(() => currentStep.value === 'finished');
  
  const activeStepIndex = computed(() => {
      switch (currentStep.value) {
          case 'connecting': return 0; // Visually on first step while connecting
          case 'download': return 0;
          case 'install': return 1;
          case 'start': return 2;
          case 'finished': return 3;
          case 'error':
               return stepIndexOnError.value !== null ? stepIndexOnError.value : 0;
          default: return -1; // Idle state, no active step shown usually
      }
  });
  
  const stepStatus = computed(() => {
      if (currentStep.value === 'error') return 'error';
      if (currentStep.value === 'connecting' || isInstalling.value) return 'process';
      return 'wait';
  });
  
  const finishStatus = computed(() => {
      if (currentStep.value === 'finished') return 'success';
      if (currentStep.value === 'error') return 'error';
      return 'process';
  });
  
  const overallProgressStatus = computed(() => {
    if (currentStep.value === 'error') return 'exception';
    if (currentStep.value === 'finished') return 'success';
    return undefined; // Default 'processing' blue
  });
  
  const stepTitle = computed(() => {
       switch (currentStep.value) {
          case 'connecting': return '连接中';
          case 'download': return '下载';
          case 'install': return '安装 (模拟)';
          case 'start': return '启动';
          case 'finished': return '完成';
          case 'error': return '错误';
          default: return '状态';
      }
  });
  
  const stepTitleOnError = computed(() => { // For displaying error title
       switch (stepIndexOnError.value) {
          case 0: return '下载';
          case 1: return '安装 (模拟)';
          case 2: return '启动';
          default: return '未知步骤';
       }
  });
  
  const messageType = computed(() => {
      switch (currentStep.value) {
          case 'error': return 'danger';
          case 'finished': return 'success';
          case 'install': return 'warning';
          default: return 'primary';
      }
  });
  
  // --- Methods ---
  
  const resetState = () => {
    console.log('Resetting state...');
    currentStep.value = 'idle';
    overallProgress.value = 0;
    statusMessage.value = '准备就绪';
    error.value = null;
    rawOutputLog.value = '';
    activeLog.value = [];
    showSuccessLog.value = false;
    showErrorLog.value = false;
    showInstallWarning.value = false;
    installWarningMessage.value = '';
    stepIndexOnError.value = null;
    parsingErrorCount.value = 0; // Reset parsing error count
    closeEventSource();
    checkBackend();
  };
  
  const closeEventSource = () => {
    if (eventSource.value) {
      eventSource.value.close();
      eventSource.value = null;
      console.log('EventSource connection explicitly closed.');
    }
  };
  
  const cancelInstallation = () => {
      if (isInstalling.value || currentStep.value === 'connecting') {
          ElMessage.warning('已手动取消连接。后端过程可能仍在继续。');
          const lastStep = currentStep.value; // Remember step before setting error
          error.value = "用户手动取消";
          currentStep.value = 'error';
          // Determine step index based on when cancelled
          stepIndexOnError.value = ['connecting', 'download'].includes(lastStep) ? 0 :
                                   lastStep === 'install' ? 1 :
                                   lastStep === 'start' ? 2 : activeStepIndex.value;
          closeEventSource();
      }
  };
  
  const appendLog = (line: string | undefined) => {
      if (line && typeof line === 'string') {
          const timestamp = new Date().toLocaleTimeString();
          rawOutputLog.value += `[${timestamp}] ${line.trim()}\n`;
      }
  }
  
  const startInstallation = () => {
    if (!backendReachable.value) {
        ElMessage.error("无法连接到后端服务，请检查。");
        return;
    }
    resetState();
    currentStep.value = 'connecting'; // Set initial step to connecting
    statusMessage.value = '正在建立 SSE 连接...';
    overallProgress.value = 0;
    rawOutputLog.value = '--- 请求安装 ---\n';
    activeLog.value = ['1'];
    parsingErrorCount.value = 0; // Reset counter
  
    closeEventSource();
  
    console.log(`Connecting to SSE endpoint: ${installUrl.value}`);
    eventSource.value = new EventSource(installUrl.value, { withCredentials: false });
  
    // --- onmessage: Handle incoming messages ---
    // **关键改动:** 调整 catch 块，使其不关闭连接
    eventSource.value.addEventListener('message', (event: MessageEvent) => {
      console.log("Raw SSE data received:", event.data); // Log raw data first
  
      try {
        // Attempt to parse potentially garbled data
        const data = JSON.parse(event.data);
        console.log("Parsed SSE data:", data);
        parsingErrorCount.value = 0; // Reset error count on successful parse
  
        // Update state based on parsed data
        if (data.step) { currentStep.value = data.step as StepType; }
        if (typeof data.progress === 'number' && data.progress >= 0) { overallProgress.value = Math.min(data.progress, 100); }
        if (typeof data.message === 'string') {
            statusMessage.value = data.message;
            if (currentStep.value === 'install' && data.message.includes("模拟")) {
                showInstallWarning.value = true; installWarningMessage.value = data.message;
            } else if (currentStep.value !== 'install') { showInstallWarning.value = false; }
        }
        if (data.rawLine) { appendLog(data.rawLine); }
  
        // Handle final states (Error or Done)
        if (data.error) {
          console.error('Backend reported error:', data.error);
          error.value = data.error;
          currentStep.value = 'error';
          stepIndexOnError.value = activeStepIndex.value; // Use computed value at time of error
          ElNotification({ title: '安装失败', message: error.value, type: 'error', duration: 0 });
          closeEventSource(); // Close connection on explicit backend error
        } else if (data.done && data.step === 'finished') {
          console.log('Backend reported finished.');
          currentStep.value = 'finished';
          overallProgress.value = 100;
          ElNotification({ title: '安装成功', message: statusMessage.value, type: 'success', duration: 5000 });
          closeEventSource(); // Close connection on explicit backend success
        } else if (data.done && data.step !== 'finished') {
            // Handle unexpected 'done' signal
             console.warn("Backend reported done=true but step is not 'finished':", data.step);
             if (!data.error) {
                  error.value = `流程意外终止于步骤 '${data.step}'`;
                  currentStep.value = 'error';
                  stepIndexOnError.value = activeStepIndex.value;
                  ElNotification({ title: '流程错误', message: error.value, type: 'error', duration: 0 });
                  closeEventSource();
             }
        }
        // If it's just a progress update, do nothing more here
  
      } catch (e) {
        // --- Graceful handling of JSON parse errors ---
        parsingErrorCount.value++;
        console.error(`解析 SSE 消息 #${parsingErrorCount.value} 时出错:`, e, "原始数据:", event.data);
        appendLog(`*** 前端错误: 无法解析消息 (第 ${parsingErrorCount.value} 次): ${event.data} - 错误: ${e} ***`);
  
        // Show a warning message to the user, but *do not* close the connection
        ElMessage.warning(`收到格式错误的消息 (错误 #${parsingErrorCount.value})，已跳过。`);
  
        // Optional: Implement logic to give up after too many consecutive errors
        // if (parsingErrorCount.value > 10) { // Example threshold
        //    console.error("连续收到过多解析错误，关闭连接。");
        //    error.value = "与后端通信时连续发生错误，连接已关闭。";
        //    currentStep.value = 'error';
        //    stepIndexOnError.value = activeStepIndex.value;
        //    ElNotification({ title: '通信错误', message: error.value, type: 'error', duration: 0 });
        //    closeEventSource();
        // }
      }
    });
  
    // --- onopen: Connection established ---
    eventSource.value.onopen = () => {
      console.log('SSE connection opened successfully.');
      // No need to change currentStep here, message handler will set it to 'download' etc.
      statusMessage.value = '已连接，等待后端开始...';
      appendLog('--- SSE 连接已建立 ---');
      // Set step back from 'connecting' if needed, or let message handler do it
      // currentStep.value = 'download'; // Or wait for first message
    };
  
    // --- onerror: Handle connection errors ---
    eventSource.value.onerror = (event) => {
      console.error('EventSource 连接失败:', event);
      let readyState = event.target ? (event.target as EventSource).readyState : 'unknown';
      console.error(`EventSource readyState: ${readyState} (0=CONNECTING, 1=OPEN, 2=CLOSED)`);
      appendLog(`--- SSE 连接错误 (State: ${readyState}) ---`);
  
      let errorMessage = `无法连接到后端 (${installUrl.value}) 或连接意外断开。`;
  
      // Only set fatal error state if we were actively trying to install/connect
      if (currentStep.value === 'connecting' || isInstalling.value) {
          error.value = errorMessage;
          statusMessage.value = "连接错误";
          currentStep.value = 'error';
          // Try to determine step index based on when error occurred
          stepIndexOnError.value = activeStepIndex.value > 0 ? activeStepIndex.value : 0; // Default to first step if connecting
          ElNotification({ title: '连接错误', message: errorMessage, type: 'error', duration: 0 });
      } else {
          console.warn("SSE error occurred outside of active installation/connection phase.");
          if(currentStep.value === 'idle') {
              statusMessage.value = '与后端连接失败。';
              backendReachable.value = false;
          }
      }
      closeEventSource(); // Close on any error
    };
  };
  
  // --- Backend Check & Lifecycle (Unchanged) ---
  const checkBackend = async () => {
      console.log(`[Debug] 正在尝试连接健康检查 URL: ${healthCheckUrl.value}`);
      isCheckingBackend.value = true;
      try {
          const response = await fetch(healthCheckUrl.value, { method: 'GET', mode: 'cors', signal: AbortSignal.timeout(10000) });
          backendReachable.value = response.ok;
          if (!response.ok) console.warn(`Backend health check non-OK status: ${response.status}`);
          else console.log("Backend health check successful.");
      } catch (err: any) {
          console.error(`[Debug] Fetch 请求失败，错误对象:`, err);
          console.error('Backend health check failed:', err.name === 'TimeoutError' ? 'Request timed out' : err);
          backendReachable.value = false;
      } finally {
          isCheckingBackend.value = false;
      }
  };
  let healthCheckInterval: number | null = null;
  onMounted(() => { checkBackend(); healthCheckInterval = window.setInterval(checkBackend, 30000); });
  onUnmounted(() => { closeEventSource(); if (healthCheckInterval) clearInterval(healthCheckInterval); });
  watch(backendReachable, (newVal) => { if (currentStep.value === 'idle') { statusMessage.value = newVal ? '准备就绪' : '等待连接后端...'; } });
  
  // --- Log Toggling (Unchanged) ---
  const toggleSuccessLog = () => showSuccessLog.value = !showSuccessLog.value;
  const toggleErrorLog = () => showErrorLog.value = !showErrorLog.value;
  
  </script>
  
  <style scoped> /* Styles remain the same */
  .installer-card { margin-bottom: 20px; }
  .card-header { display: flex; justify-content: space-between; align-items: center; font-weight: bold; }
  .card-header .el-icon { margin-right: 8px; vertical-align: middle; }
  .card-header .el-tag { margin-left: 10px; }
  .installer-content { padding: 15px 5px; min-height: 150px; }
  .initial-state ol { list-style: none; padding-left: 5px; margin: 15px 0; }
  .initial-state li { margin-bottom: 10px; display: flex; align-items: center; }
  .initial-state li .el-icon { margin-right: 8px; font-size: 1.1em; }
  .initial-state li .el-text { margin-left: 5px;}
  .progress-state .status-box { text-align: center; margin-top: 20px; padding: 10px; background-color: #f4f4f5; border-radius: 4px; }
  .progress-state p { font-size: 0.95em; }
  .raw-output { background-color: #f8f9fa; border: 1px solid #e9ecef; border-radius: 4px; padding: 15px; max-height: 350px; overflow-y: auto; font-family: 'Consolas', 'Monaco', 'monospace'; font-size: 0.8em; line-height: 1.5; color: #495057; margin-top: 10px; white-space: pre-wrap; word-break: break-all; }
  .success-state .raw-output, .error-state .raw-output { margin-top: 0; }
  .success-log { border-color: var(--el-color-success-light-5); background-color: var(--el-color-success-light-9); }
  .error-log { border-color: var(--el-color-danger-light-5); background-color: var(--el-color-danger-light-9); }
  .card-footer { display: flex; justify-content: flex-end; align-items: center; padding-top: 10px; }
  .el-result { padding: 10px 0; }
  :deep(.el-step__head.is-process) { color: var(--el-color-primary); border-color: var(--el-color-primary); }
  :deep(.el-step__head.is-error) { color: var(--el-color-danger); border-color: var(--el-color-danger); }
  :deep(.el-step__title.is-error) { color: var(--el-color-danger); }
  :deep(.el-step__description.is-error) { color: var(--el-color-danger); }
  </style>