<template>
    <div class="pod-page-container">
      <!-- Breadcrumbs -->
      <el-breadcrumb separator="/" class="page-breadcrumb">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item>资源管理</el-breadcrumb-item>
        <el-breadcrumb-item>Pods</el-breadcrumb-item>
      </el-breadcrumb>
  
      <!-- Header: Title & Create Button -->
      <div class="page-header">
        <h1 class="page-title">Pods 列表</h1>
        <el-button
          type="primary"
          :icon="PlusIcon"
          @click="handleAddPod"
          :disabled="!selectedNamespace || loading.page"
        >
          创建 Pod (YAML)
        </el-button>
      </div>
  
      <!-- Cluster Knowledge Alert -->
      <el-alert
        title="关于 Pods"
        type="info"
        show-icon
        :closable="true"
        class="info-alert"
        description="在 Kubernetes 中，Pod 是最小的可部署单元。一个 Pod 表示一个运行中的进程，通常包含一个或多个容器。Pod 提供了容器之间共享的网络和存储资源，并规定了这些容器应如何运行。"
      />
  
      <!-- Filter Bar: Namespace, Search, Refresh -->
      <div class="filter-bar">
         <el-select
             v-model="selectedNamespace"
             placeholder="请选择命名空间"
             @change="handleNamespaceChange"
             filterable
             :loading="loading.namespaces"
             class="filter-item namespace-select"
             style="width: 280px;"
         >
             <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
              <template #empty>
                 <div style="padding: 10px; text-align: center; color: #999;">
                     {{ loading.namespaces ? '正在加载...' : '无可用命名空间' }}
                 </div>
             </template>
         </el-select>
  
         <el-input
             v-model="searchQuery"
             placeholder="搜索 Pod 名称 / IP / 节点 / 状态"
             :prefix-icon="SearchIcon"
             clearable
             @input="handleSearchDebounced"
             class="filter-item search-input"
             :disabled="!selectedNamespace || loading.pods"
         />
  
         <el-tooltip content="刷新列表" placement="top">
             <el-button
               :icon="RefreshIcon"
               circle
               @click="fetchPodData"
               :loading="loading.pods"
               :disabled="!selectedNamespace"
             />
         </el-tooltip>
      </div>
  
      <!-- Pods Table -->
      <el-table
        :data="paginatedData"
        v-loading="loading.pods"
        border
        stripe
        style="width: 100%"
        @sort-change="handleSortChange"
        class="pod-table"
        :default-sort="{ prop: 'createdAt', order: 'descending' }"
        row-key="uid"
      >
        <!-- Columns -->
         <el-table-column prop="name" label="名称" min-width="250" sortable="custom" fixed show-overflow-tooltip>
            <template #default="{ row }">
               <el-icon class="pod-icon"><Tickets /></el-icon>
               <span class="pod-name">{{ row.name }}</span>
           </template>
       </el-table-column>
       <el-table-column prop="namespace" label="命名空间" min-width="150" sortable="custom" show-overflow-tooltip />
       <el-table-column prop="status" label="状态" min-width="130" sortable="custom" align="center">
           <template #default="{ row }">
               <el-tooltip placement="top" :disabled="!row.reason && !row.message">
                   <template #content>
                       <div v-if="row.reason">Reason: {{ row.reason }}</div>
                       <div v-if="row.message">Message: {{ row.message }}</div>
                   </template>
                   <el-tag :type="getStatusTagType(row.status)" size="small" effect="light" class="status-tag">
                       <el-icon class="status-icon" :class="getSpinClass(row.status)">
                           <component :is="getStatusIcon(row.status)" />
                       </el-icon>
                       {{ row.status }}
                   </el-tag>
               </el-tooltip>
           </template>
       </el-table-column>
       <el-table-column prop="ip" label="Pod IP" min-width="150" sortable="custom" show-overflow-tooltip />
       <el-table-column prop="node" label="所在节点" min-width="180" sortable="custom" show-overflow-tooltip />
       <!-- FIXED: Apply formatting only for display -->
       <el-table-column prop="createdAt" label="创建时间" min-width="180" sortable="custom">
          <template #default="{ row }">
              {{ formatTimestamp(row.createdAt) }}
          </template>
       </el-table-column>
  
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-tooltip content="查看日志" placement="top">
              <el-button link type="primary" :icon="DocumentIcon" @click="handleViewLogs(row)" />
            </el-tooltip>
            <el-tooltip content="进入容器" placement="top">
              <el-button link type="primary" :icon="MonitorIcon" @click="handleExec(row)" />
            </el-tooltip>
            <el-tooltip content="编辑 YAML" placement="top">
              <el-button link type="primary" :icon="EditPenIcon" @click="handleEditPod(row)" />
            </el-tooltip>
            <el-tooltip content="删除" placement="top">
              <el-button link type="danger" :icon="DeleteIcon" @click="handleDeletePod(row)" />
            </el-tooltip>
          </template>
        </el-table-column>
        <template #empty>
          <!-- FIXED: Use v-bind for image-size -->
          <el-empty v-if="!loading.pods && !selectedNamespace" description="请先选择一个命名空间以加载 Pods" :image-size="100" />
          <el-empty v-else-if="!loading.pods && filteredData.length === 0 && searchQuery" description="未找到匹配的 Pods" :image-size="100" />
          <el-empty v-else-if="!loading.pods && allPods.length === 0" :description="`在命名空间 '${selectedNamespace}' 中未找到 Pods`" :image-size="100" />
         </template>
      </el-table>
  
      <!-- Pagination -->
      <!-- Use totalFilteredPods for client-side pagination total -->
      <div class="pagination-container" v-if="!loading.pods && totalFilteredPods > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="totalFilteredPods"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
          :disabled="loading.pods"
        />
      </div>
  
      <!-- Create/Edit Pod Dialog -->
      <el-dialog
          :title="isEditMode ? `编辑 Pod YAML: ${editingPodName}` : '创建 Pod (YAML)'"
          v-model="yamlDialogConfig.visible"
          width="70%"
          :close-on-click-modal="false"
          @closed="handleYamlDialogClose"
      >
        <el-alert type="info" :closable="false" style="margin-bottom: 15px;">
          建议通过 Deployment 等控制器管理 Pod。直接操作 Pod 常用于调试。确保 YAML 中的 `namespace` (如果提供) 与当前选定命名空间 (`${selectedNamespace || '未选定'}`) 匹配。
        </el-alert>
        <div class="yaml-editor-container" style="height: 50vh;">
          <!-- Monaco Editor Component -->
          <vue-monaco-editor
            v-if="yamlDialogConfig.visible"
            v-model:value="yamlDialogConfig.content"
            theme="vs-dark"
            language="yaml"
            :options="{ minimap: { enabled: false }, scrollBeyondLastLine: false, automaticLayout: true }"
            style="height: 100%; border: 1px solid #eee;"
          />
        </div>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="yamlDialogConfig.visible = false">取 消</el-button>
            <el-button type="primary" @click="handleSaveYaml" :loading="yamlDialogConfig.saving">
              {{ isEditMode ? '更新 YAML' : '应用 YAML' }}
            </el-button>
          </div>
        </template>
      </el-dialog>
  
       <!-- View Logs Dialog -->
       <el-dialog
           title="查看 Pod 日志"
           v-model="logDialogConfig.visible"
           width="75%"
           :close-on-click-modal="false"
           @closed="handleLogDialogClose"
           top="5vh"
           class="log-dialog"
           destroy-on-close
       >
           <div v-if="logDialogConfig.targetPod">
               <el-form :inline="true" size="small" class="log-options-form">
                   <el-form-item label="Pod">
                      <strong>{{ logDialogConfig.targetPod.name }}</strong> ({{ logDialogConfig.targetPod.namespace }})
                   </el-form-item>
                   <el-form-item label="容器" required>
                       <el-select
                           v-model="logDialogConfig.selectedContainer"
                           placeholder="请选择容器"
                           style="width: 200px;"
                           :disabled="logDialogConfig.loadingContainers"
                           clearable
                           @change="fetchLogs"
                           filterable
                       >
                           <el-option
                               v-for="container in logDialogConfig.containers"
                               :key="container.name"
                               :label="container.name"
                               :value="container.name"
                           />
                            <template #empty>
                                <div style="padding: 10px; text-align: center; color: #999;">
                                    {{ logDialogConfig.loadingContainers ? '加载中...' : '无可用容器' }}
                                </div>
                            </template>
                       </el-select>
                   </el-form-item>
                   <el-form-item label="显示行数">
                       <el-input-number
                           v-model="logDialogConfig.tailLines"
                           :min="10"
                           :max="10000"
                           :step="100"
                           controls-position="right"
                           style="width: 130px;"
                           placeholder="默认全部"
                           clearable
                           :disabled="!logDialogConfig.selectedContainer || logDialogConfig.follow"
                           @change="fetchLogs"
                       />
                   </el-form-item>
                   <el-form-item>
                      <el-button type="primary" :icon="RefreshIcon" @click="fetchLogs" :loading="logDialogConfig.loadingLogs" :disabled="!logDialogConfig.selectedContainer">刷新</el-button>
                   </el-form-item>
               </el-form>
  
               <div class="log-display-container">
                  <pre v-loading="logDialogConfig.loadingLogs" class="log-content">{{ logDialogConfig.content || (logDialogConfig.selectedContainer ? (logDialogConfig.loadingLogs ? '正在加载日志...' : '请点击刷新加载日志') : '请先选择一个容器') }}</pre>
               </div>
  
           </div>
           <template #footer>
               <span class="dialog-footer">
                   <el-button @click="logDialogConfig.visible = false">关 闭</el-button>
               </span>
           </template>
       </el-dialog>
  
       <!-- Exec into Container Dialog -->
      <el-dialog
          title="进入 Pod 容器终端"
          v-model="execDialogConfig.visible"
          width="75%"
          :close-on-click-modal="false"
          @opened="initTerminal"
          @closed="handleExecDialogClose"
          top="5vh"
          class="exec-dialog"
          destroy-on-close
      >
           <div v-if="execDialogConfig.targetPod">
               <el-form :inline="true" size="small" style="margin-bottom: 10px;">
                    <el-form-item label="Pod">
                       <strong>{{ execDialogConfig.targetPod.name }}</strong> ({{ execDialogConfig.targetPod.namespace }})
                    </el-form-item>
                   <el-form-item label="容器" required>
                       <el-select
                           v-model="execDialogConfig.selectedContainer"
                           placeholder="请选择容器"
                           style="width: 200px;"
                           :disabled="execDialogConfig.loadingContainers || execDialogConfig.connected"
                           clearable
                           filterable
                       >
                            <el-option
                                v-for="container in execDialogConfig.containers"
                                :key="container.name"
                                :label="container.name"
                                :value="container.name"
                            />
                             <template #empty>
                                 <div style="padding: 10px; text-align: center; color: #999;">
                                     {{ execDialogConfig.loadingContainers ? '加载中...' : '无可用容器' }}
                                 </div>
                             </template>
                       </el-select>
                   </el-form-item>
                    <el-form-item label="命令">
                         <el-input v-model="execDialogConfig.command" style="width: 150px;" placeholder="默认 sh" :disabled="execDialogConfig.connected"/>
                    </el-form-item>
                    <el-form-item>
                       <el-button
                          type="primary"
                          @click="connectWebSocket"
                          :loading="execDialogConfig.connecting"
                          :disabled="!execDialogConfig.selectedContainer || execDialogConfig.connected">
                           {{ execDialogConfig.connected ? '已连接' : '连 接' }}
                       </el-button>
                       <el-button @click="disconnectWebSocket" v-if="execDialogConfig.connected">断开连接</el-button>
                    </el-form-item>
                    <el-form-item v-if="execDialogConfig.statusText">
                         <span :class="`status-text status-${execDialogConfig.statusType}`">{{ execDialogConfig.statusText }}</span>
                    </el-form-item>
               </el-form>
  
              <!-- Terminal Container -->
              <div ref="terminalContainerRef" class="terminal-container"></div>
  
           </div>
           <template #footer>
               <span class="dialog-footer">
                   <el-button @click="execDialogConfig.visible = false">关 闭</el-button>
               </span>
           </template>
       </el-dialog>
  
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, reactive, computed, onMounted, watch, nextTick, onBeforeUnmount, markRaw } from "vue"; // Import markRaw
  import { ElMessage, ElMessageBox } from "element-plus";
  import { request } from "@/utils/service"; // Ensure correct path
  import dayjs from "dayjs";
  import { debounce } from 'lodash-es';
  // Standard import for vue-monaco-editor
//   import VueMonacoEditor from 'vue-monaco-editor';
  // Standard imports for xterm
  import { Terminal } from 'xterm';
  import { FitAddon } from 'xterm-addon-fit';
  import 'xterm/css/xterm.css';
  
  
  import {
      Plus as PlusIcon, Search as SearchIcon, Refresh as RefreshIcon, Tickets,
      CircleCheckFilled, WarningFilled, CloseBold, Loading as LoadingIcon,
      QuestionFilled, Document as DocumentIcon, Monitor as MonitorIcon,
      EditPen as EditPenIcon, Delete as DeleteIcon
  } from '@element-plus/icons-vue';
  
  // --- Constants ---
  const VITE_API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://192.168.1.100:8080";
  const wsProtocol = VITE_API_BASE_URL.startsWith('https://') ? 'wss://' : 'ws://';
  const wsHostPort = VITE_API_BASE_URL.replace(/^https?:\/\//, '');
  const WS_BASE_URL = `${wsProtocol}${wsHostPort}`;
  
  // --- Interfaces (Ensure they match backend and usage) ---
  interface K8sContainer { name: string; image: string; [key: string]: any; }
  interface PodSpec { containers: K8sContainer[]; initContainers?: K8sContainer[]; [key: string]: any; }
  interface PodDetail {
    name: string; namespace: string; uid: string; spec: PodSpec;
    status?: { phase?: string; containerStatuses?: Array<{ name: string; ready: boolean; state: any }>; initContainerStatuses?: Array<{ name: string; ready: boolean; state: any }>; };
    [key: string]: any;
  }
  interface PodDetailApiResponse { code: number; data: PodDetail; message: string; }
  interface PodApiItem {
    uid: string; name: string; namespace: string;
    labels?: { [key: string]: string }; annotations?: { [key: string]: string } | null;
    status: string; reason?: string; message?: string;
    ip?: string; node?: string; createdAt: string; // Expect ISO string
  }
  interface PodListApiResponseData { items: PodApiItem[]; total: number }
  interface PodApiResponse { code: number; data: PodListApiResponseData; message: string }
  interface PodDisplayItem extends PodApiItem {} // Extend if needed for frontend state
  interface NamespaceListResponse { code: number; data: string[]; message: string }
  interface YamlApiResponse { code: number; data: string; message: string; }
  
  // --- Reactive State ---
  const allPods = ref<PodDisplayItem[]>([]);
  const namespaces = ref<string[]>([]);
  const selectedNamespace = ref<string>("");
  const currentPage = ref(1);
  const pageSize = ref(10);
  // Removed totalPods as it's derived client-side now via totalFilteredPods
  const searchQuery = ref("");
  const sortParams = reactive({ key: 'createdAt', order: 'descending' as ('ascending' | 'descending' | null) });
  
  const loading = reactive({ page: false, namespaces: false, pods: false });
  
  // --- Create/Edit Dialog State ---
  const isEditMode = ref(false);
  const editingPodName = ref<string | null>(null);
  const yamlDialogConfig = reactive({ visible: false, content: '', saving: false });
  const placeholderYaml = computed(() => `apiVersion: v1
  kind: Pod
  metadata:
    generateName: new-pod-from-ui-
    namespace: ${selectedNamespace.value || 'default'} # Will use selected NS
    labels:
      app: myapp-${Math.random().toString(36).substring(2, 6)} # Add random part to label
      created-by: cilikube-ui
  spec:
    containers:
    - name: main-container # Use a more descriptive name
      image: nginx:alpine
      ports:
      - containerPort: 80
        protocol: TCP
      resources: # Add basic resource requests/limits as good practice
        requests:
          memory: "64Mi"
          cpu: "100m"
        limits:
          memory: "128Mi"
          cpu: "200m"
    restartPolicy: Never # Default for Pod, use Deployment for Always/OnFailure
  `);
  
  // --- Log Dialog State ---
  const logDialogConfig = reactive({ visible: false, targetPod: null as PodDisplayItem | null, containers: [] as K8sContainer[], selectedContainer: '', content: '', follow: false, tailLines: 500 as number | undefined, loadingContainers: false, loadingLogs: false });
  
  // --- Exec Dialog State ---
  const execDialogConfig = reactive({ visible: false, targetPod: null as PodDisplayItem | null, containers: [] as K8sContainer[], selectedContainer: '', command: 'sh', loadingContainers: false, connecting: false, connected: false, statusText: '', statusType: 'info' as 'info' | 'success' | 'error' });
  const terminalContainerRef = ref<HTMLDivElement | null>(null);
  let terminal: Terminal | null = null;
  let fitAddon: FitAddon | null = null;
  let websocket: WebSocket | null = null;
  
  
  // --- Computed Properties ---
  const filteredData = computed(() => {
      const query = searchQuery.value.trim().toLowerCase();
      if (!query) return allPods.value;
      return allPods.value.filter(pod =>
          pod.name.toLowerCase().includes(query) ||
          (pod.ip && pod.ip.toLowerCase().includes(query)) ||
          (pod.node && pod.node.toLowerCase().includes(query)) ||
          pod.status.toLowerCase().includes(query)
      );
  });
  
  // FIXED: Corrected sorting logic for dates and other types
  const sortedData = computed(() => {
      const data = [...filteredData.value];
      const { key, order } = sortParams;
  
      if (!key || !order) return data; // No sorting needed
  
      data.sort((a, b) => {
          // Get values, handle potential undefined/null early
          const valA = a[key as keyof PodDisplayItem] ?? '';
          const valB = b[key as keyof PodDisplayItem] ?? '';
  
          let comparison = 0;
  
          // Specific handling for 'createdAt' using original ISO string
          if (key === 'createdAt') {
              // Ensure values are valid date strings before parsing
              const timeA = typeof valA === 'string' && valA ? dayjs(valA).valueOf() : 0;
              const timeB = typeof valB === 'string' && valB ? dayjs(valB).valueOf() : 0;
              const numA = isNaN(timeA) ? 0 : timeA;
              const numB = isNaN(timeB) ? 0 : timeB;
              if (numA < numB) comparison = -1;
              else if (numA > numB) comparison = 1;
          } else {
              // General comparison for other fields (treat as strings, case-insensitive)
              const strA = String(valA).toLowerCase();
              const strB = String(valB).toLowerCase();
              if (strA < strB) comparison = -1;
              else if (strA > strB) comparison = 1;
          }
  
          return order === 'ascending' ? comparison : -comparison;
      });
      return data;
  });
  
  
  const paginatedData = computed(() => {
      const start = (currentPage.value - 1) * pageSize.value;
      const end = start + pageSize.value;
      return sortedData.value.slice(start, end);
  });
  
  const totalFilteredPods = computed(() => filteredData.value.length);
  
  
  // --- Helper Functions ---
  // FIXED: Format timestamp only for display
  const formatTimestamp = (isoTimestamp: string): string => {
      if (!isoTimestamp) return 'N/A';
      const d = dayjs(isoTimestamp);
      return d.isValid() ? d.format("YYYY-MM-DD HH:mm:ss") : 'Invalid Date';
  }
  // Status helpers remain the same
  const getStatusTagType = (status: string): 'success' | 'warning' | 'danger' | 'info' => { /* ... same ... */
      if (!status) return 'info';
      const lowerStatus = status.toLowerCase();
      if (lowerStatus === 'running' || lowerStatus === 'succeeded') return 'success';
      if (lowerStatus.includes('pending') || lowerStatus.includes('creating') || lowerStatus.includes('initializing')) return 'warning';
      if (lowerStatus.includes('failed') || lowerStatus.includes('error') || lowerStatus.includes('crashloop') || lowerStatus.includes('imagepull') || lowerStatus === 'nodeaffinity' || lowerStatus === 'unschedulable') return 'danger';
      if (lowerStatus.includes('terminating')) return 'info';
      return 'info';
  }
  const getStatusIcon = (status: string) => { /* ... same ... */
      if (!status) return QuestionFilled;
      const lowerStatus = status.toLowerCase();
       if (lowerStatus === 'running' || lowerStatus === 'succeeded') return CircleCheckFilled;
       if (lowerStatus.includes('pending') || lowerStatus.includes('creating') || lowerStatus.includes('terminating') || lowerStatus.includes('initializing')) return LoadingIcon;
       if (lowerStatus.includes('failed') || lowerStatus.includes('error') || lowerStatus.includes('crashloop') || lowerStatus.includes('imagepull') || lowerStatus === 'nodeaffinity' || lowerStatus === 'unschedulable') return CloseBold;
      return QuestionFilled;
  }
  const getSpinClass = (status: string) => { /* ... same ... */
       if (!status) return '';
      const lowerStatus = status.toLowerCase();
      return (lowerStatus.includes('pending') || lowerStatus.includes('creating') || lowerStatus.includes('terminating') || lowerStatus.includes('initializing')) ? 'is-loading' : '';
  }
  
  
  // --- API Interaction ---
  const fetchNamespaces = async () => {
      loading.namespaces = true;
      try {
          const response = await request<NamespaceListResponse>({ url: "/api/v1/namespaces", method: "get", baseURL: VITE_API_BASE_URL });
          if (response.code === 200 && Array.isArray(response.data)) {
              namespaces.value = response.data.sort();
              if (namespaces.value.length > 0 && !selectedNamespace.value) {
                   selectedNamespace.value = namespaces.value.find(ns => ns === 'default') || namespaces.value[0];
              } else if (namespaces.value.length === 0) {
                  ElMessage.warning("未找到任何命名空间。");
                  selectedNamespace.value = ""; allPods.value = [];
              }
          } else {
              ElMessage.error(`获取命名空间失败: ${response.message || '格式错误'}`);
              namespaces.value = []; selectedNamespace.value = ""; allPods.value = [];
          }
      } catch (error: any) {
          console.error("获取命名空间失败:", error);
          ElMessage.error(`获取命名空间出错: ${error.message || '网络请求失败'}`);
          namespaces.value = []; selectedNamespace.value = ""; allPods.value = [];
      } finally {
          loading.namespaces = false;
      }
  }
  
  const fetchPodData = async () => {
      if (!selectedNamespace.value) {
          allPods.value = [];
          loading.pods = false;
          return;
      }
      loading.pods = true;
      // Ensure page number is valid before fetching
      if (currentPage.value < 1) { currentPage.value = 1; }
  
      try {
          const params = { /* Server-side params if needed */ };
          const url = `/api/v1/namespaces/${selectedNamespace.value}/pods`;
          const response = await request<PodApiResponse>({ url, method: "get", params, baseURL: VITE_API_BASE_URL });
  
          if (response.code === 200 && response.data?.items) {
              // FIXED: Store raw ISO date string for sorting
              allPods.value = response.data.items.map(item => ({
                  ...item,
                  // Keep createdAt as the original ISO string from backend
              }));
  
              // Adjust page number *after* data is loaded and filtering/sorting is applied
              // This logic runs implicitly because totalFilteredPods updates
               nextTick(() => {
                  const totalPages = Math.ceil(totalFilteredPods.value / pageSize.value);
                  if (currentPage.value > totalPages && totalPages > 0) {
                      currentPage.value = totalPages;
                  } else if (totalFilteredPods.value === 0 && currentPage.value !== 1) {
                      currentPage.value = 1;
                  }
              });
  
          } else {
              ElMessage.error(`获取 Pod 数据失败 (ns: ${selectedNamespace.value}): ${response.message || '未知错误'}`);
              allPods.value = [];
          }
      } catch (error: any) {
          console.error("获取 Pod 数据失败:", error);
          ElMessage.error(`获取 Pod 数据出错 (ns: ${selectedNamespace.value}): ${error.message || '网络请求失败'}`);
          allPods.value = [];
      } finally {
          // Ensure loading stops even if errors occur during data processing
          loading.pods = false;
      }
  }
  
  
  async function fetchPodDetails(namespace: string, name: string): Promise<PodDetail | null> {
       // Reuse table loading indicator or add specific loading state
       // loading.pods = true;
       try {
           const response = await request<PodDetailApiResponse>({ url: `/api/v1/namespaces/${namespace}/pods/${name}`, method: "get", baseURL: VITE_API_BASE_URL });
           if (response.code === 200 && response.data) {
               return response.data;
           } else {
               ElMessage.error(`获取 Pod 详情 '${name}' 失败: ${response.message || '未知错误'}`);
               return null;
           }
       } catch (error: any) {
           console.error(`获取 Pod 详情 '${name}' 失败:`, error);
           ElMessage.error(`获取 Pod 详情 '${name}' 出错: ${error.message || '网络错误'}`);
           return null;
       } finally {
          // loading.pods = false;
       }
   }
  
  // --- Event Handlers (Client-side based) ---
  const handleNamespaceChange = () => { currentPage.value = 1; searchQuery.value = ''; sortParams.key = 'createdAt'; sortParams.order = 'descending'; fetchPodData(); };
  const handlePageChange = (page: number) => { currentPage.value = page; };
  const handleSizeChange = (size: number) => { pageSize.value = size; currentPage.value = 1; };
  const handleSearchDebounced = debounce(() => { currentPage.value = 1; }, 300);
  const handleSortChange = ({ prop, order }: { prop: string | null; order: 'ascending' | 'descending' | null }) => { sortParams.key = prop || 'createdAt'; sortParams.order = order; currentPage.value = 1; };
  
  
  // --- Dialog and CRUD Actions ---
  const handleAddPod = () => {
      if (!selectedNamespace.value) { ElMessage.warning("请先选择一个命名空间"); return; }
      isEditMode.value = false; editingPodName.value = null;
      yamlDialogConfig.content = placeholderYaml.value; yamlDialogConfig.visible = true;
  };
  
  const handleEditPod = async (pod: PodDisplayItem) => {
       isEditMode.value = true; editingPodName.value = pod.name;
       yamlDialogConfig.saving = true; yamlDialogConfig.content = "# 正在加载 YAML..."; yamlDialogConfig.visible = true;
       try {
           const response = await request<YamlApiResponse>({ url: `/api/v1/namespaces/${pod.namespace}/pods/${pod.name}/yaml`, method: 'get', baseURL: VITE_API_BASE_URL });
           if (response.code === 200 && typeof response.data === 'string') {
               yamlDialogConfig.content = response.data;
           } else { /* ... error handling ... */
               yamlDialogConfig.content = `# 获取 YAML 失败: ${response.message || '未知错误'}`;
               ElMessage.error(`获取 Pod YAML 失败: ${response.message || '格式错误'}`);
           }
       } catch (error: any) { /* ... error handling ... */
           console.error("获取 Pod YAML 失败:", error);
           const errMsg = error.response?.data?.message || error.message || '网络错误';
           yamlDialogConfig.content = `# 获取 YAML 出错:\n${errMsg}`;
           ElMessage.error(`获取 Pod YAML 出错: ${errMsg}`);
       } finally { yamlDialogConfig.saving = false; }
   };
  
  const handleSaveYaml = async () => {
      const targetNamespace = selectedNamespace.value; // Get current namespace
      if (!targetNamespace && !isEditMode.value) { ElMessage.error("未选择命名空间。"); return; }
      yamlDialogConfig.saving = true;
      const currentYaml = yamlDialogConfig.content;
      try {
          let response: any; let successMsg = '';
          if (isEditMode.value && editingPodName.value) {
              response = await request({ url: `/api/v1/namespaces/${targetNamespace}/pods/${editingPodName.value}/yaml`, method: 'put', baseURL: VITE_API_BASE_URL, headers: { 'Content-Type': 'application/yaml' }, data: currentYaml });
              successMsg = `Pod "${editingPodName.value}" 更新成功！`;
          } else {
              response = await request({ url: `/api/v1/namespaces/${targetNamespace}/pods`, method: 'post', baseURL: VITE_API_BASE_URL, headers: { 'Content-Type': 'application/yaml' }, data: currentYaml });
              const createdName = response.data?.name || '新创建的 Pod';
              successMsg = `Pod "${createdName}" 创建成功！`;
          }
          if (response.code === 200 || response.code === 201) {
              ElMessage.success(successMsg); yamlDialogConfig.visible = false; fetchPodData();
          } else { ElMessage.error(`操作失败: ${response.message || '未知后端错误'}`); }
      } catch (error: any) { /* ... error handling ... */
          console.error("应用 YAML 失败:", error);
          const errMsg = error.response?.data?.message || error.message || '请求失败';
          ElMessage.error(`应用 YAML 出错: ${errMsg}`);
      } finally { yamlDialogConfig.saving = false; }
  };
  
  const handleYamlDialogClose = () => { yamlDialogConfig.content = ''; isEditMode.value = false; editingPodName.value = null; }
  
  const handleDeletePod = (pod: PodDisplayItem) => { /* ... same logic as before ... */
      ElMessageBox.confirm(
          `确定要删除 Pod "${pod.name}" (命名空间: ${pod.namespace}) 吗？`, '确认删除',
          { confirmButtonText: '删除', cancelButtonText: '取消', type: 'warning' }
      ).then(async () => {
          loading.pods = true;
          try {
              await request({ url: `/api/v1/namespaces/${pod.namespace}/pods/${pod.name}`, method: "delete", baseURL: VITE_API_BASE_URL });
              ElMessage.success(`Pod "${pod.name}" 已删除`);
              // Refresh list or optimistic update
              await fetchPodData();
          } catch (error: any) {
               console.error("删除 Pod 失败:", error);
               const errMsg = error.response?.data?.message || error.message || '请求失败';
               ElMessage.error(`删除 Pod "${pod.name}" 失败: ${errMsg}`);
               // Only stop loading on error if not refreshing
               // loading.pods = false; // Let refresh handle loading state
          }
          // finally { loading.pods = false; } // Let refresh handle loading state
      }).catch(() => ElMessage.info('删除操作已取消'));
  };
  
  // --- View Logs ---
  const handleViewLogs = async (pod: PodDisplayItem) => { /* ... same logic ... */
       logDialogConfig.targetPod = pod; logDialogConfig.visible = true; logDialogConfig.loadingContainers = true;
       logDialogConfig.containers = []; logDialogConfig.selectedContainer = ''; logDialogConfig.content = '正在加载容器列表...';
       const details = await fetchPodDetails(pod.namespace, pod.name);
       if (details?.spec) {
          const running = details.spec.containers || [];
          const init = (details.spec.initContainers || []).map(c => ({...c, name: `[init] ${c.name}`}));
          const all = [...running, ...init]; logDialogConfig.containers = all;
          if (all.length > 0) { logDialogConfig.selectedContainer = running[0]?.name || all[0].name; await fetchLogs(); }
          else { logDialogConfig.content = '此 Pod 没有找到容器。'; }
       } else { logDialogConfig.content = '获取 Pod 详情失败。'; }
       logDialogConfig.loadingContainers = false;
  };
  const fetchLogs = async () => { /* ... same logic ... */
       if (!logDialogConfig.targetPod || !logDialogConfig.selectedContainer) { logDialogConfig.content = '请选择容器。'; return; }
       logDialogConfig.loadingLogs = true; logDialogConfig.content = '正在加载日志...';
       try {
           const actualContainerName = logDialogConfig.selectedContainer.replace(/^\[init\]\s/, '');
           const response = await request<string>({ url: `/api/v1/namespaces/${logDialogConfig.targetPod.namespace}/pods/${logDialogConfig.targetPod.name}/logs`, method: 'get', baseURL: VITE_API_BASE_URL, params: { container: actualContainerName, tailLines: logDialogConfig.tailLines || undefined }, responseType: 'text' });
           logDialogConfig.content = response || '(日志内容为空)';
       } catch (error: any) { /* ... error handling ... */
           console.error("获取日志失败:", error);
           const errMsg = error.response?.data || error.message || '请求失败';
           logDialogConfig.content = `# 获取日志出错:\n${errMsg}`;
           ElMessage.error(`获取日志出错: ${errMsg}`);
       } finally { logDialogConfig.loadingLogs = false; }
  };
  const handleLogDialogClose = () => { /* ... same ... */ logDialogConfig.targetPod = null; logDialogConfig.containers = []; logDialogConfig.selectedContainer = ''; logDialogConfig.content = ''; logDialogConfig.follow = false; logDialogConfig.tailLines = 500; };
  
  // --- Exec into Container ---
  const handleExec = async (pod: PodDisplayItem) => { /* ... same logic ... */
      execDialogConfig.targetPod = pod; execDialogConfig.visible = true; execDialogConfig.loadingContainers = true;
      execDialogConfig.containers = []; execDialogConfig.selectedContainer = ''; execDialogConfig.connected = false; execDialogConfig.connecting = false;
      execDialogConfig.statusText = '加载容器列表...'; execDialogConfig.statusType = 'info';
      const details = await fetchPodDetails(pod.namespace, pod.name);
      const runningContainers = details?.status?.containerStatuses?.filter(cs => cs.state?.running)?.map(cs => details.spec.containers.find(c => c.name === cs.name)).filter(Boolean) as K8sContainer[] | undefined || [];
       if (runningContainers.length > 0) { execDialogConfig.containers = runningContainers; execDialogConfig.selectedContainer = runningContainers[0].name; execDialogConfig.statusText = '请确认容器并连接。'; }
       else { execDialogConfig.statusText = '未找到该 Pod 中正在运行的容器。'; execDialogConfig.containers = []; execDialogConfig.statusType = 'warning'; }
       execDialogConfig.loadingContainers = false;
  };
  const initTerminal = async () => { /* ... same logic, ensure markRaw is used ... */
       await nextTick(); if (!terminalContainerRef.value || terminal) return;
       terminal = markRaw(new Terminal({ cursorBlink: true, rows: 25, theme: { background: '#222', foreground: '#ddd' }, convertEol: true, fontFamily: 'Consolas, "Courier New", monospace', fontSize: 13, letterSpacing: 1 }));
       fitAddon = markRaw(new FitAddon()); terminal.loadAddon(fitAddon);
       terminal.open(terminalContainerRef.value);
       try { fitAddon.fit(); } catch (e) { console.error("Error fitting terminal:", e) }
       terminal.onData((data) => { if (websocket?.readyState === WebSocket.OPEN) { websocket.send(data); } });
       terminal.writeln('\r\n请选择容器并点击 "连接" 按钮...');
  };
  const connectWebSocket = () => { /* ... same logic ... */
       if (!execDialogConfig.targetPod || !execDialogConfig.selectedContainer || execDialogConfig.connected || execDialogConfig.connecting) return;
       execDialogConfig.connecting = true; execDialogConfig.statusText = '正在连接...'; execDialogConfig.statusType = 'info';
       terminal?.reset(); terminal?.writeln('尝试连接 WebSocket...');
       const { namespace, name } = execDialogConfig.targetPod; const container = execDialogConfig.selectedContainer; const command = execDialogConfig.command || 'sh';
       const params = new URLSearchParams({ container, command, tty: 'true', stdin: 'true', stdout: 'true', stderr: 'true' });
       const wsUrl = `${WS_BASE_URL}/api/v1/namespaces/${namespace}/pods/${name}/exec?${params.toString()}`;
       websocket = new WebSocket(wsUrl);
       websocket.onopen = () => { /* ... same ... */ execDialogConfig.connecting = false; execDialogConfig.connected = true; execDialogConfig.statusText = `已连接 (${container})`; execDialogConfig.statusType = 'success'; terminal?.writeln('\r\n\x1b[32m连接成功！\x1b[0m'); terminal?.focus(); if (fitAddon) fitAddon.fit(); };
       websocket.onmessage = (event) => { /* ... same ... */ if (event.data instanceof ArrayBuffer) { terminal?.write(new Uint8Array(event.data)); } else if (typeof event.data === 'string') { terminal?.write(event.data); } };
       websocket.onerror = (event) => { /* ... same ... */ console.error("WebSocket error:", event); execDialogConfig.connecting = false; execDialogConfig.connected = false; execDialogConfig.statusText = 'WebSocket 连接错误'; execDialogConfig.statusType = 'error'; terminal?.writeln(`\r\n\x1b[31mWebSocket 错误\x1b[0m`); };
       websocket.onclose = (event) => { /* ... same ... */ console.log("WebSocket closed:", event.code, event.reason); execDialogConfig.connecting = false; execDialogConfig.connected = false; const reason = event.reason || `Code ${event.code}`; execDialogConfig.statusText = `连接已断开 (${reason})`; execDialogConfig.statusType = 'info'; terminal?.writeln(`\r\n\x1b[33m连接已关闭: ${reason}\x1b[0m`); websocket = null; };
  };
  const disconnectWebSocket = () => { /* ... same logic ... */ if (websocket) { execDialogConfig.statusText = '正在断开连接...'; execDialogConfig.statusType = 'info'; websocket.close(1000, "User disconnected"); } };
  const handleExecDialogClose = () => { /* ... same logic ... */ disconnectWebSocket(); if (terminal) { terminal.dispose(); terminal = null; } fitAddon = null; execDialogConfig.targetPod = null; execDialogConfig.containers = []; execDialogConfig.selectedContainer = ''; execDialogConfig.connected = false; execDialogConfig.connecting = false; execDialogConfig.statusText = ''; };
  
  // --- Lifecycle Hooks ---
  onMounted(async () => {
      loading.page = true;
      await fetchNamespaces();
      if (selectedNamespace.value) {
          await fetchPodData();
      } else {
          console.log("挂载时未选择命名空间，跳过 Pod 数据获取。");
          // Explicitly stop loading if no initial fetch happens
          loading.pods = false;
      }
      loading.page = false;
  });
  onBeforeUnmount(() => { handleExecDialogClose(); });
  
  </script>
  
  <style lang="scss" scoped>
  /* Styles remain the same */
  $page-padding: 20px;
  $spacing-md: 15px;
  $spacing-lg: 20px;
  $font-size-base: 14px;
  $font-size-small: 12px;
  $font-size-large: 16px;
  $font-size-extra-large: 24px;
  $border-radius-base: 4px;
  $kube-pod-icon-color: #326ce5;
  
  .pod-page-container { padding: $page-padding; background-color: var(--el-bg-color-page, #f5f7fa); }
  .page-breadcrumb { margin-bottom: $spacing-lg; }
  .page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: $spacing-md; flex-wrap: wrap; gap: $spacing-md; }
  .page-title { font-size: $font-size-extra-large; font-weight: 600; color: var(--el-text-color-primary); margin: 0; }
  .info-alert { margin-bottom: $spacing-lg; background-color: var(--el-color-info-light-9); :deep(.el-alert__description) { font-size: $font-size-small; color: var(--el-text-color-regular); line-height: 1.6; } }
  .filter-bar { display: flex; align-items: center; flex-wrap: wrap; gap: $spacing-md; margin-bottom: $spacing-lg; padding: $spacing-md; background-color: var(--el-bg-color-overlay); border-radius: $border-radius-base; border: 1px solid var(--el-border-color-lighter); }
  .filter-item {}
  .namespace-select { width: 240px; }
  .search-input { width: 300px; }
  .pod-table { border-radius: $border-radius-base; border: 1px solid var(--el-border-color-lighter); overflow: hidden; :deep(th.el-table__cell) { background-color: var(--el-fill-color-lighter); color: var(--el-text-color-secondary); font-weight: 600; font-size: $font-size-small; } :deep(td.el-table__cell) { padding: 8px 10px; font-size: $font-size-base; vertical-align: middle; } .pod-icon { margin-right: 8px; color: $kube-pod-icon-color; vertical-align: middle; font-size: 18px; position: relative; top: -1px; } .pod-name { font-weight: 500; vertical-align: middle; color: var(--el-text-color-regular); } .status-tag { display: inline-flex; align-items: center; gap: 4px; padding: 0 6px; height: 22px; line-height: 20px; font-size: $font-size-small; cursor: default; } .status-icon { font-size: 12px; } .is-loading { animation: rotating 1.5s linear infinite; } }
  @keyframes rotating { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  .el-table .el-button.is-link { font-size: 14px; padding: 4px; margin: 0 3px; vertical-align: middle; }
  .pagination-container { display: flex; justify-content: flex-end; margin-top: $spacing-lg; }
  .yaml-editor-container { border: 1px solid var(--el-border-color); border-radius: $border-radius-base; overflow: hidden; }
  .dialog-footer { text-align: right; padding-top: 10px; }
  .log-options-form { margin-bottom: 15px; padding-bottom: 10px; border-bottom: 1px solid var(--el-border-color-lighter); }
  .log-display-container { height: 60vh; overflow: hidden; background-color: #222; /* Darker bg */ border-radius: $border-radius-base; border: 1px solid var(--el-border-color); }
  .log-content { height: 100%; overflow: auto; padding: 10px; margin: 0; font-family: Consolas, "Courier New", monospace; font-size: 13px; line-height: 1.5; color: #ddd; /* Lighter text */ white-space: pre-wrap; word-break: break-all; }
  .terminal-container { width: 100%; height: 60vh; background-color: #222; border-radius: $border-radius-base; overflow: hidden; border: 1px solid var(--el-border-color); }
  :deep(.exec-dialog .el-dialog__body) { padding: 15px 20px; }
  :deep(.log-dialog .el-dialog__body) { padding: 10px 20px 20px 20px; }
  .status-text { font-size: $font-size-small; margin-left: 10px; &.status-success { color: var(--el-color-success); } &.status-error { color: var(--el-color-error); } &.status-info { color: var(--el-text-color-secondary); } }
  </style>