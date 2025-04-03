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
           :loading="loading.page"
           :disabled="!selectedNamespace"
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
              <!-- No "All Namespaces" option -->
              <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
               <template #empty>
                  <div style="padding: 10px; text-align: center; color: #999;">
                      {{ loading.namespaces ? '正在加载...' : '无可用命名空间' }}
                  </div>
              </template>
          </el-select>
  
          <el-input
              v-model="searchQuery"
              placeholder="搜索 Pod 名称 / IP / 节点"
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
          <el-table-column prop="name" label="名称" min-width="250" sortable="custom" fixed show-overflow-tooltip>
               <template #default="{ row }">
                  <el-icon class="pod-icon"><Tickets /></el-icon>
                  <span class="pod-name">{{ row.name }}</span>
              </template>
          </el-table-column>
          <el-table-column prop="namespace" label="命名空间" min-width="150" sortable="custom" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" min-width="130" sortable="custom" align="center">
              <template #default="{ row }">
                  <el-tag :type="getStatusTagType(row.status)" size="small" effect="light" class="status-tag">
                      <el-icon class="status-icon" :class="getSpinClass(row.status)">
                          <component :is="getStatusIcon(row.status)" />
                      </el-icon>
                      {{ row.status }}
                  </el-tag>
              </template>
          </el-table-column>
          <el-table-column prop="ip" label="Pod IP" min-width="150" sortable="custom" show-overflow-tooltip />
          <el-table-column prop="node" label="所在节点" min-width="180" sortable="custom" show-overflow-tooltip />
          <el-table-column prop="createdAt" label="创建时间" min-width="180" sortable="custom" />
          <el-table-column label="操作" width="180" align="center" fixed="right">
              <template #default="{ row }">
                  <el-tooltip content="查看日志" placement="top">
                      <el-button link type="primary" :icon="DocumentIcon" @click="viewPodLogs(row)" />
                  </el-tooltip>
                   <el-tooltip content="进入容器" placement="top">
                      <el-button link type="primary" :icon="MonitorIcon" @click="execIntoPod(row)" />
                  </el-tooltip>
                   <el-tooltip content="编辑 YAML" placement="top">
                      <el-button link type="primary" :icon="EditPenIcon" @click="editPodYaml(row)" />
                  </el-tooltip>
                  <el-tooltip content="删除" placement="top">
                      <el-button link type="danger" :icon="DeleteIcon" @click="handleDeletePod(row)" />
                  </el-tooltip>
              </template>
          </el-table-column>
           <template #empty>
            <el-empty v-if="!loading.pods && !selectedNamespace" description="请先选择一个命名空间以加载 Pods" image-size="100" />
            <el-empty v-else-if="!loading.pods && paginatedData.length === 0" :description="`在命名空间 '${selectedNamespace}' 中未找到 Pods`" image-size="100" />
           </template>
      </el-table>
  
      <!-- Pagination -->
      <div class="pagination-container" v-if="!loading.pods && totalPods > 0">
          <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="totalPods"
              layout="total, sizes, prev, pager, next, jumper"
              background
              @size-change="handleSizeChange"
              @current-change="handlePageChange"
              :disabled="loading.pods"
          />
      </div>
  
      <!-- Create Pod Dialog (YAML focus) -->
      <el-dialog title="创建 Pod (YAML)" v-model="dialogVisible" width="70%" :close-on-click-modal="false">
         <el-alert type="info" :closable="false" style="margin-bottom: 20px;">
           建议使用 Deployment、StatefulSet 等控制器管理 Pods。直接创建 Pod 通常用于调试或特殊场景。请在此处粘贴或编辑 Pod 的 YAML 配置。确保 YAML 中的 `namespace` 与当前选定的命名空间 (`${selectedNamespace || '未选定'}`) 匹配或省略（将使用选定命名空间）。
         </el-alert>
         <!-- Integrate a YAML editor component here -->
         <div class="yaml-editor-placeholder">
              YAML 编辑器占位符 (例如: 使用 Monaco Editor 或 Codemirror)
              <pre>{{ placeholderYaml }}</pre>
              <!-- Replace above pre with actual editor component -->
              <!-- Example: <YamlEditor v-model="yamlContent" ref="yamlEditorRef" /> -->
         </div>
        <template #footer>
          <div class="dialog-footer">
              <el-button @click="dialogVisible = false">取 消</el-button>
              <el-button type="primary" @click="handleSaveYaml" :loading="loading.dialogSave">应用 YAML</el-button>
          </div>
        </template>
      </el-dialog>
  
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, reactive, computed, onMounted } from "vue"
  import { ElMessage, ElMessageBox } from "element-plus"
  // Removed: type { FormInstance } from 'element-plus'
  import { request } from "@/utils/service" // Ensure correct path
  import dayjs from "dayjs"
  import { debounce } from 'lodash-es'
  
  import {
      Plus as PlusIcon,
      Search as SearchIcon,
      Refresh as RefreshIcon,
      Tickets, // Icon for Pod
      CircleCheckFilled,
      WarningFilled,
      CloseBold,
      Loading as LoadingIcon,
      QuestionFilled,
      Document as DocumentIcon,
      Monitor as MonitorIcon,
      EditPen as EditPenIcon,
      Delete as DeleteIcon
  } from '@element-plus/icons-vue'
  
  // --- Interfaces ---
  interface PodApiItem {
    name: string
    namespace: string
    labels?: { [key: string]: string }
    annotations?: { [key: string]: string } | null
    status: string
    ip?: string
    node?: string
    createdAt: string
    uid?: string // Add if API provides it
  }
  interface PodListApiResponseData { items: PodApiItem[]; total: number }
  interface PodApiResponse { code: number; data: PodListApiResponseData; message: string }
  interface PodDisplayItem {
    uid: string
    name: string
    namespace: string
    status: string
    ip: string
    node: string
    createdAt: string
    labels?: { [key: string]: string }
  }
  interface NamespaceListResponse { code: number; data: string[]; message: string }
  
  // --- Reactive State ---
  const allPods = ref<PodDisplayItem[]>([])
  const namespaces = ref<string[]>([])
  const selectedNamespace = ref<string>("") // ** Crucial: Start empty **
  const currentPage = ref(1)
  const pageSize = ref(10)
  const totalPods = ref(0)
  const searchQuery = ref("")
  const sortParams = reactive({ key: 'createdAt', order: 'descending' as ('ascending' | 'descending' | null) })
  
  const loading = reactive({
      page: false,
      namespaces: false,
      pods: false,
      dialogSave: false,
  })
  
  // Dialog state (YAML focus)
  const dialogVisible = ref(false)
  // const podFormRef = ref<FormInstance>() // Removed simple form ref
  const yamlContent = ref(""); // State for the actual YAML editor model
  const placeholderYaml = computed(() => `apiVersion: v1
  kind: Pod
  metadata:
    name: my-new-pod
    namespace: ${selectedNamespace.value || 'default'} # Reflect selected NS
    labels:
      app: myapp
  spec:
    containers:
    - name: my-container
      image: nginx:latest
      ports:
      - containerPort: 80
  `);
  
  
  // --- Computed Properties ---
  const filteredData = computed(() => {
      const query = searchQuery.value.trim().toLowerCase()
      if (!query) return allPods.value;
      return allPods.value.filter(pod =>
          pod.name.toLowerCase().includes(query) ||
          (pod.ip && pod.ip.toLowerCase().includes(query)) ||
          (pod.node && pod.node.toLowerCase().includes(query))
      );
  });
  
  const sortedData = computed(() => {
      const data = [...filteredData.value];
      const { key, order } = sortParams;
      if (!key || !order) return data;
  
      data.sort((a, b) => {
          let valA = a[key as keyof PodDisplayItem];
          let valB = b[key as keyof PodDisplayItem];
          if (key === 'createdAt') {
              const timeA = valA ? dayjs(valA, "YYYY-MM-DD HH:mm:ss").valueOf() : 0;
              const timeB = valB ? dayjs(valB, "YYYY-MM-DD HH:mm:ss").valueOf() : 0;
              valA = isNaN(timeA) ? 0 : timeA;
              valB = isNaN(timeB) ? 0 : timeB;
          } else {
               valA = valA ?? ''; valB = valB ?? '';
          }
          let comparison = 0;
          if (valA < valB) comparison = -1;
          else if (valA > valB) comparison = 1;
          return order === 'ascending' ? comparison : -comparison;
      });
      return data;
  });
  
  const paginatedData = computed(() => {
      const start = (currentPage.value - 1) * pageSize.value;
      const end = start + pageSize.value;
      return sortedData.value.slice(start, end);
  });
  
  
  // --- Helper Functions ---
  const formatTimestamp = (timestamp: string): string => {
      if (!timestamp) return 'N/A';
      return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss");
  }
  const getStatusTagType = (status: string): 'success' | 'warning' | 'danger' | 'info' => { /* ... as before ... */
      const lowerStatus = status?.toLowerCase();
      if (lowerStatus === 'running') return 'success';
      if (lowerStatus === 'succeeded') return 'success';
      if (lowerStatus === 'pending' || lowerStatus === 'containercreating') return 'warning';
      if (lowerStatus === 'failed' || lowerStatus === 'error') return 'danger';
      if (lowerStatus === 'terminating') return 'info';
      return 'info';
  }
  const getStatusIcon = (status: string) => { /* ... as before ... */
      const lowerStatus = status?.toLowerCase();
      if (lowerStatus === 'running') return CircleCheckFilled;
      if (lowerStatus === 'succeeded') return CircleCheckFilled;
      if (lowerStatus === 'pending' || lowerStatus === 'containercreating') return LoadingIcon;
      if (lowerStatus === 'failed' || lowerStatus === 'error') return CloseBold;
      if (lowerStatus === 'terminating') return LoadingIcon;
      return QuestionFilled;
   }
  const getSpinClass = (status: string) => { /* ... as before ... */
      const lowerStatus = status?.toLowerCase();
      return (lowerStatus === 'pending' || lowerStatus === 'containercreating' || lowerStatus === 'terminating') ? 'is-loading' : '';
  }
  
  // --- API Interaction ---
  const fetchNamespaces = async () => {
      loading.namespaces = true;
      try {
          const response = await request<NamespaceListResponse>({ url: "/api/v1/namespaces", method: "get", baseURL: "http://192.168.1.100:8080" });
          if (response.code === 200 && Array.isArray(response.data)) {
              namespaces.value = response.data;
              // ** Important: Select a default only if namespaces are found **
              if (namespaces.value.length > 0 && !selectedNamespace.value) {
                   // Try to find 'default', otherwise pick the first one
                   selectedNamespace.value = namespaces.value.find(ns => ns === 'default') || namespaces.value[0];
                   // console.log('Default namespace selected:', selectedNamespace.value);
                   // ** Trigger pod fetch for the newly selected default namespace **
                   // await fetchPodData(); // Moved to onMounted after both fetches
              } else if (namespaces.value.length === 0) {
                  ElMessage.warning("未找到任何命名空间。");
              }
          } else {
              ElMessage.error(`获取命名空间失败: ${response.message || '数据格式错误'}`);
              namespaces.value = []; // Ensure empty on error
          }
      } catch (error: any) {
          console.error("获取命名空间失败:", error);
          ElMessage.error(`获取命名空间失败: ${error.message || '网络请求失败'}`);
          namespaces.value = []; // Ensure empty on error
      } finally {
          loading.namespaces = false;
      }
  }
  
  const fetchPodData = async () => {
      // ** Guard: Only fetch if a namespace is actually selected **
      if (!selectedNamespace.value) {
          // console.log("Skipping pod fetch: No namespace selected.");
          allPods.value = [];
          totalPods.value = 0;
          return; // Do not proceed without a namespace
      }
      // console.log(`Fetching pods for namespace: ${selectedNamespace.value}`);
      loading.pods = true;
      try {
          const params: Record<string, any> = { /* Server-side params if needed */ };
          // ** Corrected: Always use the namespaced endpoint **
          const url = `/api/v1/namespaces/${selectedNamespace.value}/pods`;
  
          const response = await request<PodApiResponse>({ url, method: "get", params, baseURL: "http://192.168.1.100:8080" });
  
          if (response.code === 200 && response.data?.items) {
              // console.log(`Received ${response.data.items.length} pods, total: ${response.data.total}`);
              allPods.value = response.data.items.map((item, index) => ({
                  uid: item.uid || `${item.namespace}-${item.name}-${index}`,
                  name: item.name,
                  namespace: item.namespace,
                  status: item.status || 'Unknown',
                  ip: item.ip || 'N/A',
                  node: item.node || 'N/A',
                  createdAt: formatTimestamp(item.createdAt),
                  labels: item.labels,
              }));
              totalPods.value = response.data.total;
               const totalPages = Math.ceil(totalPods.value / pageSize.value);
               if (currentPage.value > totalPages && totalPages > 0) currentPage.value = totalPages;
               else if (totalPods.value === 0) currentPage.value = 1;
          } else {
              ElMessage.error(`获取 Pod 数据失败 (ns: ${selectedNamespace.value}): ${response.message || '未知错误'}`);
              allPods.value = []; totalPods.value = 0;
          }
      } catch (error: any) {
          console.error("获取 Pod 数据失败:", error);
          ElMessage.error(`获取 Pod 数据出错 (ns: ${selectedNamespace.value}): ${error.message || '网络请求失败'}`);
          allPods.value = []; totalPods.value = 0;
      } finally {
          loading.pods = false;
      }
  }
  
  // --- Event Handlers ---
  const handleNamespaceChange = () => {
      currentPage.value = 1;
      searchQuery.value = '';
      sortParams.key = 'createdAt'; sortParams.order = 'descending';
      // ** Crucial: Fetch data AFTER the namespace has been changed **
      fetchPodData();
  };
  
  const handlePageChange = (page: number) => { currentPage.value = page; /* Fetch only if server-side */ };
  const handleSizeChange = (size: number) => { pageSize.value = size; currentPage.value = 1; /* Fetch only if server-side */ };
  const handleSearchDebounced = debounce(() => { currentPage.value = 1; /* Fetch only if server-side */ }, 300);
  const handleSortChange = ({ prop, order }: { prop: string | null; order: 'ascending' | 'descending' | null }) => {
      sortParams.key = prop || 'createdAt';
      sortParams.order = order;
      currentPage.value = 1;
      /* Fetch only if server-side */
  };
  
  
  // --- Dialog and CRUD Actions ---
  const handleAddPod = () => {
      if (!selectedNamespace.value) {
          ElMessage.warning("请先选择一个命名空间");
          return;
      }
      // Update placeholder YAML with current namespace before showing
      yamlContent.value = placeholderYaml.value; // Reset editor content
      dialogVisible.value = true;
  };
  
  // Placeholder for saving YAML
  const handleSaveYaml = async () => {
      if (!selectedNamespace.value) {
           ElMessage.error("无法创建 Pod，未选择命名空间。");
           return;
      }
      loading.dialogSave = true;
      ElMessage.info("模拟应用 YAML...");
      // --- Replace with actual YAML editor interaction and API call ---
      // const currentYaml = yamlEditorRef.value.getContent(); // Get from editor
      // const namespaceToUse = parseNamespaceFromYaml(currentYaml) || selectedNamespace.value; // Determine NS
      // try {
      //    const response = await request({... POST request with YAML ...});
      //    ... handle response ...
      // } catch(e) {...}
  
      await new Promise(resolve => setTimeout(resolve, 500)); // Simulate API call
      loading.dialogSave = false;
      dialogVisible.value = false;
      ElMessage.success(`模拟 Pod 应用成功 (ns: ${selectedNamespace.value})`);
      fetchPodData(); // Refresh list
  };
  
  
  const handleDeletePod = (pod: PodDisplayItem) => { /* ... as before ... */
      ElMessageBox.confirm(
          `确定要删除 Pod "${pod.name}" (命名空间: ${pod.namespace}) 吗？`,
          '确认删除', { type: 'warning' }
      ).then(async () => {
          // const loadingInstance = ElLoading.service({ lock: true, text: `正在删除 Pod ${pod.name}...` }); // Optional overlay
          loading.pods = true; // Use table loading indicator
          try {
              const response = await request<{ code: number; message: string }>({
                  url: `/api/v1/namespaces/${pod.namespace}/pods/${pod.name}`,
                  method: "delete",
                  baseURL: "http://192.168.1.100:8080", // Keep baseURL if needed
              });
              // Use status code from Gin handler (likely 200 or 204 if status() was used)
              if (response.code === 200 || response.code === 204 || response.code === 202) {
                  ElMessage.success(`Pod "${pod.name}" 已删除`);
                  // Refresh might be needed if total count changed significantly or for consistency
                   await fetchPodData();
                   // Or optimistic update:
                   // allPods.value = allPods.value.filter(p => p.uid !== pod.uid);
                   // totalPods.value--;
              } else {
                   ElMessage.error(`删除 Pod 失败: ${response.message || '未知错误'}`);
                   loading.pods = false; // Stop loading on error
              }
          } catch (error: any) {
              console.error("删除 Pod 失败:", error);
              const errMsg = error.response?.data?.message || error.message || '请求失败';
              ElMessage.error(`删除 Pod 失败: ${errMsg}`);
              loading.pods = false; // Stop loading on error
          }
          // finally { loadingInstance.close(); } // If using overlay
      }).catch(() => ElMessage.info('删除操作已取消'));
  };
  
  // --- Other Actions (Simulations) ---
  const viewPodLogs = (pod: PodDisplayItem) => { ElMessage.info(`模拟: 查看 Pod "${pod.name}" 日志 (ns: ${pod.namespace})`); };
  const execIntoPod = (pod: PodDisplayItem) => { ElMessage.info(`模拟: 进入 Pod "${pod.name}" 终端 (ns: ${pod.namespace})`); };
  const editPodYaml = (pod: PodDisplayItem) => { ElMessage.info(`模拟: 打开 Pod "${pod.name}" 的 YAML 编辑器 (ns: ${pod.namespace})`); };
  
  
  // --- Lifecycle Hooks ---
  onMounted(async () => {
      loading.page = true;
      await fetchNamespaces();
      // ** Crucial: Fetch pods only AFTER namespaces are loaded AND a default is selected **
      if (selectedNamespace.value) {
          await fetchPodData();
      } else {
          // Handle the case where no namespaces were found or default couldn't be set
          console.log("No namespace selected after initial load, skipping initial pod fetch.");
      }
      loading.page = false;
  });
  
  </script>
  
  <style lang="scss" scoped>
  // Import variables if the file exists and path is correct
  // @import "@/styles/variables.scss";
  
  // Define fallbacks if variables are not used
  $page-padding: 20px;
  $spacing-md: 15px;
  $spacing-lg: 20px;
  $font-size-base: 14px;
  $font-size-small: 12px;
  $font-size-large: 16px;
  $font-size-extra-large: 24px; // Matched h1
  $border-radius-base: 4px;
  $kube-pod-icon-color: #326ce5; // Example color
  
  
  .pod-page-container {
    padding: $page-padding;
    background-color: var(--el-bg-color-page);
  }
  
  .page-breadcrumb {
    margin-bottom: $spacing-lg;
  }
  
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: $spacing-md;
    flex-wrap: wrap;
    gap: $spacing-md;
  }
  
  .page-title {
    font-size: $font-size-extra-large;
    font-weight: 600;
    color: var(--el-text-color-primary);
    margin: 0;
  }
  
  .info-alert {
    margin-bottom: $spacing-lg;
    background-color: var(--el-color-info-light-9);
     :deep(.el-alert__description) {
         font-size: $font-size-small;
         color: var(--el-text-color-regular);
         line-height: 1.6;
     }
  }
  
  .filter-bar {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: $spacing-md;
    margin-bottom: $spacing-lg;
    padding: $spacing-md;
    background-color: var(--el-bg-color);
    border-radius: $border-radius-base;
    border: 1px solid var(--el-border-color-lighter);
  }
  
  .filter-item {
    // Base styles if needed
  }
  
  .namespace-select {
    width: 240px; // Adjust width
  }
  
  .search-input {
    width: 300px; // Adjust width
  }
  
  .pod-table {
     border-radius: $border-radius-base;
     border: 1px solid var(--el-border-color-lighter);
     overflow: hidden; // Ensures border radius is applied to table content
  
      :deep(th.el-table__cell) {
         background-color: var(--el-fill-color-lighter);
         color: var(--el-text-color-secondary);
         font-weight: 600;
         font-size: $font-size-small;
     }
      :deep(td.el-table__cell) {
          padding: 8px 0; // Adjust cell padding
          font-size: $font-size-base;
          vertical-align: middle; // Align cell content vertically
      }
  
  
     .pod-icon {
         margin-right: 8px; // Increased space
         color: $kube-pod-icon-color;
         vertical-align: middle;
         font-size: 18px; // Slightly larger icon
         position: relative; // Allow fine-tuning
         top: -1px;
     }
     .pod-name {
         font-weight: 500;
         vertical-align: middle;
         color: var(--el-text-color-regular);
     }
  
     .status-tag {
       display: inline-flex;
       align-items: center;
       gap: 4px;
       padding: 0 6px; // Adjust tag padding
       height: 22px;
       line-height: 20px;
       font-size: $font-size-small; // Ensure tag text size is small
     }
     .status-icon {
         font-size: 12px; // Adjust icon size in tag
     }
  
     // Loading spin animation
     .is-loading {
         animation: rotating 1.5s linear infinite;
     }
      @keyframes rotating {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
      }
  }
  
  .el-table .el-button.is-link {
     font-size: 14px;
     padding: 4px;
     margin: 0 3px;
     vertical-align: middle;
  }
  
  
  .pagination-container {
    display: flex;
    justify-content: flex-end;
    margin-top: $spacing-lg;
  }
  
  .yaml-editor-placeholder {
      border: 1px dashed var(--el-border-color);
      padding: 20px;
      margin-top: 10px; // Space from alert
      min-height: 350px; // Ensure enough space
      max-height: 60vh; // Prevent excessive height
      background-color: var(--el-fill-color-lighter);
      color: var(--el-text-color-secondary);
      font-family: monospace;
      white-space: pre-wrap;
      overflow: auto;
      font-size: 13px;
      line-height: 1.5;
  }
  
  
  .dialog-footer {
    text-align: right;
    padding-top: 10px; // Add some space above buttons
  }
  </style>
  