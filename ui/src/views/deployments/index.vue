<template>
    <div class="deployment-page-container">
      <!-- Breadcrumbs -->
      <el-breadcrumb separator="/" class="page-breadcrumb">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item>应用管理</el-breadcrumb-item>
        <el-breadcrumb-item>Deployments</el-breadcrumb-item>
      </el-breadcrumb>
  
      <!-- Header: Title & Create Button -->
      <div class="page-header">
        <h1 class="page-title">Deployments 列表</h1>
         <el-button
           type="primary"
           :icon="PlusIcon"
           @click="handleAddDeployment"
           :loading="loading.page"
           :disabled="!selectedNamespace"
         >
           创建 Deployment (YAML)
         </el-button>
      </div>
  
       <!-- Cluster Knowledge Alert -->
       <el-alert
         title="关于 Deployments"
         type="info"
         show-icon
         :closable="true"
         class="info-alert"
         description="Deployment 为 Pod 和 ReplicaSet 提供声明式的更新能力。您负责描述 Deployment 中的期望状态，Deployment 控制器（Controller）以受控速率更改实际状态，使其变为期望状态。您可以定义 Deployment 以创建新的 ReplicaSet，或删除现有 Deployment，并通过新的 Deployment 采用其所有资源。"
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
              placeholder="搜索 Deployment 名称"
              :prefix-icon="SearchIcon"
              clearable
              @input="handleSearchDebounced"
              class="filter-item search-input"
              :disabled="!selectedNamespace || loading.deployments"
          />
  
          <el-tooltip content="刷新列表" placement="top">
              <el-button
                :icon="RefreshIcon"
                circle
                @click="fetchDeploymentData"
                :loading="loading.deployments"
                :disabled="!selectedNamespace"
              />
          </el-tooltip>
      </div>
  
      <!-- Deployments Table -->
      <el-table
          :data="paginatedData"
          v-loading="loading.deployments"
          border
          stripe
          style="width: 100%"
          @sort-change="handleSortChange"
          class="deployment-table"
          :default-sort="{ prop: 'createdAt', order: 'descending' }"
          row-key="uid"
      >
          <el-table-column prop="name" label="名称" min-width="250" sortable="custom" fixed show-overflow-tooltip>
               <template #default="{ row }">
                  <el-icon class="deployment-icon"><TakeawayBox /></el-icon>
                  <span class="deployment-name">{{ row.name }}</span>
              </template>
          </el-table-column>
          <el-table-column prop="namespace" label="命名空间" min-width="150" sortable="custom" show-overflow-tooltip />
          <el-table-column label="状态" min-width="130" sortable="custom" align="center" :sort-by="['statusSort', 'name']">
               <template #default="{ row }">
                  <el-tag :type="row.statusInfo.tagType" size="small" effect="light" class="status-tag">
                      <el-icon v-if="row.statusInfo.icon" class="status-icon" :class="row.statusInfo.spin ? 'is-loading' : ''">
                          <component :is="row.statusInfo.icon" />
                      </el-icon>
                      {{ row.statusInfo.text }}
                  </el-tag>
              </template>
          </el-table-column>
          <el-table-column label="副本 (Ready/Desired)" min-width="160" sortable="custom" align="center" :sort-by="['readyReplicas', 'desiredReplicas', 'name']">
              <template #default="{ row }">
                 <span :class="getReplicasClass(row)">{{ row.readyReplicas }} / {{ row.desiredReplicas }}</span>
              </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" min-width="180" sortable="custom" />
           <el-table-column label="镜像" min-width="250" show-overflow-tooltip>
               <template #default="{ row }">
                  <div v-for="(image, index) in row.images" :key="index" class="image-tag">
                      <el-tag size="small" type="info" effect="plain">{{ image }}</el-tag>
                  </div>
                  <span v-if="!row.images || row.images.length === 0">N/A</span>
              </template>
           </el-table-column>
          <el-table-column label="操作" width="180" align="center" fixed="right">
              <template #default="{ row }">
                   <el-tooltip content="调整副本数" placement="top">
                      <el-button link type="primary" :icon="ScaleToOriginalIcon" @click="scaleDeployment(row)" />
                   </el-tooltip>
                  <el-tooltip content="查看关联 Pods" placement="top">
                      <el-button link type="primary" :icon="ViewIcon" @click="viewRelatedPods(row)" />
                  </el-tooltip>
                   <el-tooltip content="编辑 YAML" placement="top">
                      <el-button link type="primary" :icon="EditPenIcon" @click="editDeploymentYaml(row)" />
                  </el-tooltip>
                  <el-tooltip content="删除" placement="top">
                      <el-button link type="danger" :icon="DeleteIcon" @click="handleDeleteDeployment(row)" />
                  </el-tooltip>
              </template>
          </el-table-column>
           <template #empty>
            <el-empty v-if="!loading.deployments && !selectedNamespace" description="请先选择一个命名空间以加载 Deployments" image-size="100" />
            <el-empty v-else-if="!loading.deployments && paginatedData.length === 0" :description="`在命名空间 '${selectedNamespace}' 中未找到 Deployments`" image-size="100" />
           </template>
      </el-table>
  
      <!-- Pagination -->
      <div class="pagination-container" v-if="!loading.deployments && totalDeployments > 0">
          <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="totalDeployments"
              layout="total, sizes, prev, pager, next, jumper"
              background
              @size-change="handleSizeChange"
              @current-change="handlePageChange"
              :disabled="loading.deployments"
          />
      </div>
  
      <!-- Create/Edit Dialog (YAML focus) -->
      <el-dialog title="创建/编辑 Deployment (YAML)" v-model="dialogVisible" width="70%" :close-on-click-modal="false">
         <el-alert type="info" :closable="false" style="margin-bottom: 20px;">
           请在此处粘贴或编辑 Deployment 的 YAML 配置。确保 YAML 中的 `namespace` 与当前选定的命名空间 (`${selectedNamespace || '未选定'}`) 匹配或省略。
         </el-alert>
         <!-- Integrate a YAML editor component here -->
         <div class="yaml-editor-placeholder">
              YAML 编辑器占位符 (例如: 使用 Monaco Editor 或 Codemirror)
              <pre>{{ yamlContent || placeholderYaml }}</pre>
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
  
      <!-- Scale Dialog -->
       <el-dialog title="调整副本数" v-model="scaleDialogVisible" width="400px" :close-on-click-modal="false">
           <el-form :model="scaleForm" label-width="80px" ref="scaleFormRef">
               <el-form-item label="名称">
                   <el-input :model-value="scaleForm.name" disabled />
               </el-form-item>
                <el-form-item label="命名空间">
                   <el-input :model-value="scaleForm.namespace" disabled />
               </el-form-item>
               <el-form-item label="副本数" prop="replicas" :rules="[{ required: true, message: '副本数不能为空'}, { type: 'integer', min: 0, message: '副本数必须是非负整数', trigger: 'blur'}]">
                   <el-input-number v-model.number="scaleForm.replicas" :min="0" controls-position="right" style="width: 100%" />
               </el-form-item>
           </el-form>
          <template #footer>
              <div class="dialog-footer">
                  <el-button @click="scaleDialogVisible = false">取 消</el-button>
                  <el-button type="primary" @click="handleConfirmScale" :loading="loading.scaleSave">确 定</el-button>
              </div>
          </template>
       </el-dialog>
  
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, reactive, computed, onMounted } from "vue"
  import { ElMessage, ElMessageBox } from "element-plus"
  import type { FormInstance } from 'element-plus'
  import { request } from "@/utils/service" // Ensure correct path
  import dayjs from "dayjs"
  import { debounce } from 'lodash-es'
  import yaml from 'js-yaml'; // For YAML parsing/handling (npm install js-yaml @types/js-yaml)
  
  import {
      Plus as PlusIcon,
      Search as SearchIcon,
      Refresh as RefreshIcon,
      TakeawayBox, // Icon for Deployment
      ScaleToOriginal as ScaleToOriginalIcon,
      View as ViewIcon,
      EditPen as EditPenIcon,
      Delete as DeleteIcon,
      CircleCheckFilled,
      WarningFilled,
      CloseBold,
      Loading as LoadingIcon,
      QuestionFilled,
      InfoFilled // For unknown/info status
  } from '@element-plus/icons-vue'
  
  // --- Interfaces ---
  // K8s API structures based on sample
  interface K8sMetadata { name: string; namespace: string; uid: string; resourceVersion: string; generation?: number; creationTimestamp: string; labels?: { [key: string]: string }; annotations?: { [key: string]: string }; managedFields?: any[] }
  interface K8sLabelSelector { matchLabels: { [key: string]: string } }
  interface K8sContainer { name: string; image: string; ports?: { containerPort: number, protocol: string }[]; resources?: any; terminationMessagePath?: string; terminationMessagePolicy?: string; imagePullPolicy?: string }
  interface K8sPodTemplateSpec { metadata: { creationTimestamp: null, labels: { [key: string]: string } }; spec: { containers: K8sContainer[]; restartPolicy?: string; terminationGracePeriodSeconds?: number; dnsPolicy?: string; securityContext?: any; schedulerName?: string } }
  interface K8sRollingUpdateDeployment { maxUnavailable?: string | number; maxSurge?: string | number }
  interface K8sDeploymentStrategy { type?: string; rollingUpdate?: K8sRollingUpdateDeployment }
  interface K8sDeploymentSpec { replicas: number; selector: K8sLabelSelector; template: K8sPodTemplateSpec; strategy?: K8sDeploymentStrategy; revisionHistoryLimit?: number; progressDeadlineSeconds?: number }
  interface K8sDeploymentCondition { type: string; status: string; lastUpdateTime?: string; lastTransitionTime?: string; reason?: string; message?: string }
  interface K8sDeploymentStatus { observedGeneration?: number; replicas?: number; updatedReplicas?: number; readyReplicas?: number; availableReplicas?: number; unavailableReplicas?: number; conditions?: K8sDeploymentCondition[]; collisionCount?: number }
  
  // API Response Item (Matches Go backend sample)
  interface DeploymentApiItem { metadata: K8sMetadata; spec: K8sDeploymentSpec; status: K8sDeploymentStatus }
  interface DeploymentListApiResponseData { items: DeploymentApiItem[]; total?: number } // Make total optional as backend might not always provide it accurately
  interface DeploymentApiResponse { code: number; data: DeploymentListApiResponseData; message: string }
  interface NamespaceListResponse { code: number; data: string[]; message: string }
  
  // Status info derived for display
  interface StatusInfo { text: string; tagType: 'success' | 'warning' | 'danger' | 'info'; icon: any; spin?: boolean, sort: number }
  
  // Internal Display/Table Item
  interface DeploymentDisplayItem {
    uid: string
    name: string
    namespace: string
    desiredReplicas: number
    availableReplicas: number
    readyReplicas: number
    updatedReplicas: number // Added for status calculation
    statusInfo: StatusInfo
    statusSort: number // Add a sortable numeric status value
    createdAt: string
    images: string[]
    // Raw API data for editing/details
    rawData: DeploymentApiItem
  }
  
  // --- Reactive State ---
  const allDeployments = ref<DeploymentDisplayItem[]>([])
  const namespaces = ref<string[]>([])
  const selectedNamespace = ref<string>("")
  const currentPage = ref(1)
  const pageSize = ref(10)
  const totalDeployments = ref(0) // Use API total if available, otherwise calculate
  const apiTotalCount = ref(0) // Store total from API response
  const searchQuery = ref("")
  const sortParams = reactive({ key: 'createdAt', order: 'descending' as ('ascending' | 'descending' | null) })
  
  const loading = reactive({
      page: false, namespaces: false, deployments: false, dialogSave: false, scaleSave: false
  })
  
  // Dialog state (YAML focus)
  const dialogVisible = ref(false)
  const currentEditDeployment = ref<DeploymentApiItem | null>(null); // Store raw data for editing
  const yamlContent = ref("");
  const placeholderYaml = computed(() => `apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: my-new-deployment
    namespace: ${selectedNamespace.value || 'default'}
    labels:
      app: myapp
  spec:
    replicas: 1
    selector:
      matchLabels:
        app: myapp
    template:
      metadata:
        labels:
          app: myapp
      spec:
        containers:
        - name: my-container
          image: nginx:latest
          ports:
          - containerPort: 80
  `);
  
  // Scale Dialog State
  const scaleDialogVisible = ref(false);
  const scaleFormRef = ref<FormInstance>();
  const scaleForm = reactive({ name: '', namespace: '', replicas: 1 });
  
  
  // --- Computed Properties ---
  const filteredData = computed(() => {
      const query = searchQuery.value.trim().toLowerCase()
      if (!query) return allDeployments.value;
      return allDeployments.value.filter(dep =>
          dep.name.toLowerCase().includes(query)
          // Add more fields to search if needed, e.g., labels
      );
  });
  
  const sortedData = computed(() => {
      const data = [...filteredData.value];
      const { key, order } = sortParams;
      if (!key || !order) return data;
  
      data.sort((a, b) => {
          let valA: any;
          let valB: any;
  
          // Handle specific sort keys
          if (key === 'statusSort' || key === 'readyReplicas' || key === 'desiredReplicas') {
              valA = a[key as keyof DeploymentDisplayItem] ?? 0;
              valB = b[key as keyof DeploymentDisplayItem] ?? 0;
          } else if (key === 'createdAt') {
              const timeA = a.createdAt ? dayjs(a.createdAt, "YYYY-MM-DD HH:mm:ss").valueOf() : 0;
              const timeB = b.createdAt ? dayjs(b.createdAt, "YYYY-MM-DD HH:mm:ss").valueOf() : 0;
              valA = isNaN(timeA) ? 0 : timeA;
              valB = isNaN(timeB) ? 0 : timeB;
          } else {
              // Default string comparison for name, namespace etc.
              valA = a[key as keyof DeploymentDisplayItem] ?? '';
              valB = b[key as keyof DeploymentDisplayItem] ?? '';
              // Ensure case-insensitive sort for strings if needed
              if (typeof valA === 'string') valA = valA.toLowerCase();
              if (typeof valB === 'string') valB = valB.toLowerCase();
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
  const formatTimestamp = (timestamp: string): string => { /* ... as before ... */
      if (!timestamp) return 'N/A';
      return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss");
  }
  
  // Enhanced status calculation for Deployments
  const getDeploymentStatusInfo = (dep: DeploymentApiItem): StatusInfo => {
      const desired = dep.spec.replicas ?? 0;
      const current = dep.status?.replicas ?? 0;
      const available = dep.status?.availableReplicas ?? 0;
      const ready = dep.status?.readyReplicas ?? 0;
      const updated = dep.status?.updatedReplicas ?? 0;
  
      // Check conditions for Progressing and Available
      const progressing = dep.status?.conditions?.find(c => c.type === 'Progressing');
      const availableCond = dep.status?.conditions?.find(c => c.type === 'Available');
  
      // Default status
      let status: StatusInfo = { text: 'Unknown', tagType: 'info', icon: QuestionFilled, sort: 4 };
  
      if (availableCond?.status === 'False' && availableCond.reason === 'MinimumReplicasUnavailable') {
          status = { text: 'Unavailable', tagType: 'danger', icon: CloseBold, sort: 0 };
      } else if (progressing?.status === 'True' && progressing.reason === 'NewReplicaSetAvailable') {
           // Check if update is happening or scaling
           if (updated < desired || current < desired || ready < desired || available < desired) {
               status = { text: 'Progressing', tagType: 'warning', icon: LoadingIcon, spin: true, sort: 1 };
           } else if (ready === desired && available === desired && updated === desired) {
               status = { text: 'Available', tagType: 'success', icon: CircleCheckFilled, sort: 3 };
           }
      } else if (progressing?.status === 'False') {
           status = { text: 'Failed', tagType: 'danger', icon: CloseBold, sort: 0 }; // Deployment failed to progress
      } else if (ready === desired && available === desired && updated === desired && current === desired) {
          status = { text: 'Available', tagType: 'success', icon: CircleCheckFilled, sort: 3 };
      } else if (desired > 0 && available === 0) {
          // Might be progressing but condition not updated yet, or truly unavailable
          status = { text: 'Initializing', tagType: 'warning', icon: LoadingIcon, spin: true, sort: 1 };
      } else if (available < desired) {
           status = { text: 'Degraded', tagType: 'warning', icon: WarningFilled, sort: 2 };
      }
  
      // Override if replicas set to 0 and none are running
      if (desired === 0 && current === 0) {
         status = { text: 'Scaled to 0', tagType: 'info', icon: InfoFilled, sort: 4 };
      }
  
      return status;
  }
  
  const getReplicasClass = (row: DeploymentDisplayItem): string => {
      if (row.readyReplicas === row.desiredReplicas && row.desiredReplicas > 0) return 'replicas-ok';
      if (row.desiredReplicas > 0 && row.readyReplicas < row.desiredReplicas) return 'replicas-warning';
      if (row.desiredReplicas === 0) return 'replicas-scaled-down';
      return '';
  };
  const VITE_API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://192.168.1.100:8080";
  // Extract image names from containers spec
  const extractImages = (spec: K8sDeploymentSpec | undefined): string[] => {
      if (!spec?.template?.spec?.containers) {
          return [];
      }
      return spec.template.spec.containers.map(c => c.image);
  };
  
  // --- API Interaction ---
  const fetchNamespaces = async () => { /* ... same as in Pod component ... */
      loading.namespaces = true;
      try {
          const response = await request<NamespaceListResponse>({ url: "/api/v1/namespaces", method: "get", baseURL: VITE_API_BASE_URL });
          if (response.code === 200 && Array.isArray(response.data)) {
              namespaces.value = response.data;
              if (namespaces.value.length > 0 && !selectedNamespace.value) {
                   selectedNamespace.value = namespaces.value.find(ns => ns === 'default') || namespaces.value[0];
              } else if (namespaces.value.length === 0) {
                   ElMessage.warning("未找到任何命名空间。");
              }
          } else { ElMessage.error(`获取命名空间失败: ${response.message || '数据格式错误'}`); namespaces.value = []; }
      } catch (error: any) { console.error("获取命名空间失败:", error); ElMessage.error(`获取命名空间失败: ${error.message || '网络请求失败'}`); namespaces.value = []; }
      finally { loading.namespaces = false; }
  }
  
  const fetchDeploymentData = async () => {
      if (!selectedNamespace.value) { allDeployments.value = []; totalDeployments.value = 0; apiTotalCount.value = 0; return; }
      loading.deployments = true;
      try {
          const params: Record<string, any> = { /* Server-side params */ };
          const url = `/api/v1/namespaces/${selectedNamespace.value}/deployments`;
          const response = await request<DeploymentApiResponse>({ url, method: "get", params, baseURL: VITE_API_BASE_URL });
  
          if (response.code === 200 && response.data?.items) {
              apiTotalCount.value = response.data.total ?? response.data.items.length; // Use API total or length
              allDeployments.value = response.data.items.map((item, index) => {
                   const statusInfo = getDeploymentStatusInfo(item);
                   return {
                      uid: item.metadata.uid || `${item.metadata.namespace}-${item.metadata.name}-${index}`,
                      name: item.metadata.name,
                      namespace: item.metadata.namespace,
                      desiredReplicas: item.spec?.replicas ?? 0,
                      availableReplicas: item.status?.availableReplicas ?? 0,
                      readyReplicas: item.status?.readyReplicas ?? 0,
                      updatedReplicas: item.status?.updatedReplicas ?? 0,
                      statusInfo: statusInfo,
                      statusSort: statusInfo.sort, // Add sortable status value
                      createdAt: formatTimestamp(item.metadata.creationTimestamp),
                      images: extractImages(item.spec),
                      rawData: item // Store raw data
                  };
              });
              totalDeployments.value = apiTotalCount.value; // Update total count for pagination
  
              // Adjust pagination if needed
               const totalPages = Math.ceil(totalDeployments.value / pageSize.value);
               if (currentPage.value > totalPages && totalPages > 0) currentPage.value = totalPages;
               else if (totalDeployments.value === 0) currentPage.value = 1;
  
          } else {
              ElMessage.error(`获取 Deployment 数据失败: ${response.message || '未知错误'}`);
              allDeployments.value = []; totalDeployments.value = 0; apiTotalCount.value = 0;
          }
      } catch (error: any) {
          console.error("获取 Deployment 数据失败:", error);
          ElMessage.error(`获取 Deployment 数据出错: ${error.message || '网络请求失败'}`);
          allDeployments.value = []; totalDeployments.value = 0; apiTotalCount.value = 0;
      } finally {
          loading.deployments = false;
      }
  }
  
  // --- Event Handlers ---
  const handleNamespaceChange = () => { /* ... same as Pod component ... */
      currentPage.value = 1; searchQuery.value = ''; sortParams.key = 'createdAt'; sortParams.order = 'descending';
      fetchDeploymentData();
  };
  const handlePageChange = (page: number) => { currentPage.value = page; /* Fetch only if server-side */ };
  const handleSizeChange = (size: number) => { pageSize.value = size; currentPage.value = 1; /* Fetch only if server-side */ };
  const handleSearchDebounced = debounce(() => { currentPage.value = 1; /* Fetch only if server-side */ }, 300);
  const handleSortChange = ({ prop, order }: { prop: string | null; order: 'ascending' | 'descending' | null }) => { /* ... same as Pod component ... */
      sortParams.key = prop || 'createdAt';
      sortParams.order = order;
      currentPage.value = 1;
  };
  
  
  // --- Dialog and CRUD Actions ---
  const handleAddDeployment = () => {
      if (!selectedNamespace.value) { ElMessage.warning("请先选择一个命名空间"); return; }
      currentEditDeployment.value = null; // Indicate adding
      yamlContent.value = placeholderYaml.value; // Use placeholder for new deployment
      dialogVisible.value = true;
  };
  
  // Fetch full YAML for editing
  const editDeploymentYaml = async (deployment: DeploymentDisplayItem) => {
       currentEditDeployment.value = deployment.rawData; // Store raw data if needed later
       ElMessage.info(`模拟: 获取 Deployment "${deployment.name}" 的 YAML`);
       // --- Replace with actual API call to get YAML ---
       // try {
       //    loading.deployments = true; // Show loading maybe
       //    const response = await request<DeploymentApiItem>({ // Assuming GET returns full object
       //       url: `/api/v1/namespaces/${deployment.namespace}/deployments/${deployment.name}`,
       //       method: 'get',
       //       baseURL: "VITE_API_BASE_URL"
       //    });
       //    if (response.code === 200 && response.data) {
       //        currentEditDeployment.value = response.data; // Store raw data
       //        // Remove fields typically managed by Kubernetes before displaying
       //        const displayData = cleanK8sMetadataForEdit(response.data);
       //        yamlContent.value = yaml.dump(displayData); // Convert object to YAML string
       //        dialogVisible.value = true;
       //    } else { // ... error handling ...}
       // } catch(e) { // ... error handling ...}
       // finally { loading.deployments = false; }
  
       // Simulate fetching and opening editor
       yamlContent.value = yaml.dump(deployment.rawData); // Simulate with stored raw data
       dialogVisible.value = true;
  };
  
  // Helper to remove server-managed fields before editing YAML
  // const cleanK8sMetadataForEdit = (resource: DeploymentApiItem) => {
  //    const cleaned = JSON.parse(JSON.stringify(resource)); // Deep copy
  //    delete cleaned.metadata.uid;
  //    delete cleaned.metadata.resourceVersion;
  //    delete cleaned.metadata.creationTimestamp;
  //    delete cleaned.metadata.managedFields;
  //    delete cleaned.metadata.generation; // Often managed
  //    delete cleaned.status; // Remove status for edits
  //    return cleaned;
  // }
  
  
  const handleSaveYaml = async () => {
      if (!selectedNamespace.value && !currentEditDeployment.value) { ElMessage.error("无法应用，未选择命名空间。"); return; }
      loading.dialogSave = true;
      // --- Replace with actual YAML editor interaction and API call ---
      // const currentYaml = yamlEditorRef.value.getContent(); // Get from editor
      // try {
      //     let response;
      //     let parsedYaml = yaml.load(currentYaml); // Parse to check validity and get metadata
      //     if (typeof parsedYaml !== 'object' || parsedYaml === null || !parsedYaml.metadata) {
      //         throw new Error("无效的 YAML 格式或缺少 metadata");
      //     }
      //     const name = parsedYaml.metadata.name;
      //     // Use namespace from YAML if present, otherwise fallback to selected or edited object's NS
      //     const namespaceToUse = parsedYaml.metadata.namespace || currentEditDeployment.value?.metadata.namespace || selectedNamespace.value;
      //
      //     if (!namespaceToUse) { throw new Error("无法确定目标命名空间"); }
      //
      //     if (currentEditDeployment.value) { // Editing existing
      //         response = await request({
      //              url: `/api/v1/namespaces/${namespaceToUse}/deployments/${name}`,
      //              method: 'put', // Use PUT for replace semantics
      //              headers: { 'Content-Type': 'application/yaml' },
      //              data: currentYaml
      //         });
      //     } else { // Creating new
      //          response = await request({
      //              url: `/api/v1/namespaces/${namespaceToUse}/deployments`,
      //              method: 'post',
      //              headers: { 'Content-Type': 'application/yaml' },
      //              data: currentYaml
      //         });
      //     }
      //     if (response.code === 200 || response.code === 201) {
      //        ElMessage.success(`Deployment ${currentEditDeployment.value ? '更新' : '创建'}成功`);
      //        dialogVisible.value = false; fetchDeploymentData();
      //     } else { // ... error handling ... }
      // } catch (error: any) {
      //     console.error("应用 YAML 失败:", error);
      //     const errMsg = error.message || error.response?.data?.message || '请求失败';
      //     ElMessage.error(`应用 YAML 失败: ${errMsg}`);
      // } finally { loading.dialogSave = false; }
  
       // Simulate success
       await new Promise(resolve => setTimeout(resolve, 500));
       loading.dialogSave = false;
       dialogVisible.value = false;
       const action = currentEditDeployment.value ? '更新' : '创建';
       ElMessage.success(`模拟 Deployment ${action}成功`);
       fetchDeploymentData();
  };
  
  
  const handleDeleteDeployment = (deployment: DeploymentDisplayItem) => { /* ... similar to Pod delete ... */
       ElMessageBox.confirm(
          `确定要删除 Deployment "${deployment.name}" (命名空间: ${deployment.namespace}) 吗？关联的 Pods 也会被删除。`,
          '确认删除', { type: 'warning' }
      ).then(async () => {
          loading.deployments = true;
          try {
              const response = await request<{ code: number; message: string }>({
                  url: `/api/v1/namespaces/${deployment.namespace}/deployments/${deployment.name}`,
                  method: "delete",
                  baseURL: VITE_API_BASE_URL,
              });
               if (response.code === 200 || response.code === 202 || response.code === 204) {
                  ElMessage.success(`Deployment "${deployment.name}" 已删除`);
                  await fetchDeploymentData(); // Refresh
              } else {
                   ElMessage.error(`删除 Deployment 失败: ${response.message || '未知错误'}`);
                   loading.deployments = false;
              }
          } catch (error: any) { console.error("删除 Deployment 失败:", error); ElMessage.error(`删除 Deployment 失败: ${error.response?.data?.message || error.message || '请求失败'}`); loading.deployments = false; }
      }).catch(() => ElMessage.info('删除操作已取消'));
  };
  
  // --- Scale Action ---
  const scaleDeployment = (deployment: DeploymentDisplayItem) => {
      scaleForm.name = deployment.name;
      scaleForm.namespace = deployment.namespace;
      scaleForm.replicas = deployment.desiredReplicas; // Start with current desired count
      scaleDialogVisible.value = true;
      scaleFormRef.value?.clearValidate(); // Clear potential previous validation errors
  };
  
  const handleConfirmScale = async () => {
      if (!scaleFormRef.value) return;
      await scaleFormRef.value.validate(async (valid) => {
          if (valid) {
              loading.scaleSave = true;
              ElMessage.info(`模拟调整 Deployment "${scaleForm.name}" 副本数为 ${scaleForm.replicas}`);
               // --- Actual API call for scaling ---
               // try {
               //    // K8s typically uses PATCH with Scale subresource or strategic merge patch
               //    const scalePayload = {
               //        apiVersion: "apps/v1", // Or autoscaling/v1 depending on method
               //        kind: "Scale",
               //        metadata: { name: scaleForm.name, namespace: scaleForm.namespace },
               //        spec: { replicas: scaleForm.replicas }
               //    };
               //    const response = await request({
               //         url: `/apis/apps/v1/namespaces/${scaleForm.namespace}/deployments/${scaleForm.name}/scale`,
               //         method: 'put', // PUT often used for scale subresource
               //         data: scalePayload
               //    });
               //    if (response.code === 200) {
               //        ElMessage.success("副本数调整成功");
               //        scaleDialogVisible.value = false;
               //        // Find and update the specific deployment in the list optimistically, or fetch all
               //        fetchDeploymentData();
               //    } else { //... error handling ...}
               // } catch (error) { // ... error handling ...}
               // finally { loading.scaleSave = false; }
  
                // Simulate success
               await new Promise(resolve => setTimeout(resolve, 500));
               loading.scaleSave = false;
               scaleDialogVisible.value = false;
               ElMessage.success("模拟副本数调整成功");
               fetchDeploymentData(); // Refresh list
  
          } else {
              console.log('Scale form validation failed');
              return false;
          }
      });
  };
  
  
  // --- Other Actions ---
  const viewRelatedPods = (deployment: DeploymentDisplayItem) => {
      ElMessage.info(`模拟: 查看 Deployment "${deployment.name}" 关联的 Pods`);
      // Navigate to Pods list, pre-filtered by deployment's label selector
      // const selector = deployment.rawData.spec.selector.matchLabels; // Get selector
      // const selectorString = Object.entries(selector).map(([k, v]) => `${k}=${v}`).join(',');
      // router.push(`/pods?namespace=${deployment.namespace}&labelSelector=${encodeURIComponent(selectorString)}`);
  };
  
  
  // --- Lifecycle Hooks ---
  onMounted(async () => {
      loading.page = true;
      await fetchNamespaces();
      if (selectedNamespace.value) {
          await fetchDeploymentData();
      }
      loading.page = false;
  });
  
  </script>
  
  <style lang="scss" scoped>
  // Define fallbacks or import variables
  $page-padding: 20px; $spacing-md: 15px; $spacing-lg: 20px;
  $font-size-base: 14px; $font-size-small: 12px; $font-size-large: 16px; $font-size-extra-large: 24px;
  $border-radius-base: 4px; $kube-deployment-icon-color: #16a34a; // Example Green
  
  .deployment-page-container { padding: $page-padding; background-color: var(--el-bg-color-page); }
  .page-breadcrumb { margin-bottom: $spacing-lg; }
  .page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: $spacing-md; flex-wrap: wrap; gap: $spacing-md; }
  .page-title { font-size: $font-size-extra-large; font-weight: 600; color: var(--el-text-color-primary); margin: 0; }
  .info-alert { margin-bottom: $spacing-lg; background-color: var(--el-color-info-light-9); :deep(.el-alert__description) { font-size: $font-size-small; color: var(--el-text-color-regular); line-height: 1.6; } }
  .filter-bar { display: flex; align-items: center; flex-wrap: wrap; gap: $spacing-md; margin-bottom: $spacing-lg; padding: $spacing-md; background-color: var(--el-bg-color); border-radius: $border-radius-base; border: 1px solid var(--el-border-color-lighter); }
  .filter-item { }
  .namespace-select { width: 240px; }
  .search-input { width: 300px; }
  
  .deployment-table {
     border-radius: $border-radius-base; border: 1px solid var(--el-border-color-lighter); overflow: hidden;
      :deep(th.el-table__cell) { background-color: var(--el-fill-color-lighter); color: var(--el-text-color-secondary); font-weight: 600; font-size: $font-size-small; }
      :deep(td.el-table__cell) { padding: 8px 0; font-size: $font-size-base; vertical-align: middle; }
     .deployment-icon { margin-right: 8px; color: $kube-deployment-icon-color; vertical-align: middle; font-size: 18px; position: relative; top: -1px; }
     .deployment-name { font-weight: 500; vertical-align: middle; color: var(--el-text-color-regular); }
     .status-tag { display: inline-flex; align-items: center; gap: 4px; padding: 0 6px; height: 22px; line-height: 20px; font-size: $font-size-small; }
     .status-icon { font-size: 12px; }
     .is-loading { animation: rotating 1.5s linear infinite; }
     @keyframes rotating { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  
      .replicas-ok { color: var(--el-color-success); font-weight: 500; }
      .replicas-warning { color: var(--el-color-warning); font-weight: 500; }
      .replicas-scaled-down { color: var(--el-text-color-secondary); }
  
      .image-tag {
          margin-bottom: 3px;
          display: block; // Each image on a new line if multiple
          max-width: 100%;
          .el-tag {
              max-width: 100%;
              white-space: nowrap;
              overflow: hidden;
              text-overflow: ellipsis;
          }
      }
  }
  
  .el-table .el-button.is-link { font-size: 14px; padding: 4px; margin: 0 3px; vertical-align: middle; }
  .pagination-container { display: flex; justify-content: flex-end; margin-top: $spacing-lg; }
  .yaml-editor-placeholder { border: 1px dashed var(--el-border-color); padding: 20px; margin-top: 10px; min-height: 350px; max-height: 60vh; background-color: var(--el-fill-color-lighter); color: var(--el-text-color-secondary); font-family: monospace; white-space: pre-wrap; overflow: auto; font-size: 13px; line-height: 1.5; }
  .dialog-footer { text-align: right; padding-top: 10px; }
  </style>