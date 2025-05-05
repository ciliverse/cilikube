<template>
    <div class="secret-page-container">
      <!-- Breadcrumbs -->
      <el-breadcrumb separator="/" class="page-breadcrumb">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item>配置管理</el-breadcrumb-item>
        <el-breadcrumb-item>Secrets</el-breadcrumb-item>
      </el-breadcrumb>
  
      <!-- Header: Title & Create Button -->
      <div class="page-header">
        <h1 class="page-title">Secrets 列表</h1>
         <el-button
           type="primary"
           :icon="PlusIcon"
           @click="handleAddSecret"
           :loading="loading.page"
           :disabled="!selectedNamespace"
         >
           创建 Secret (YAML)
         </el-button>
      </div>
  
       <!-- Cluster Knowledge Alert -->
       <el-alert
         title="关于 Secrets"
         type="warning" 
         show-icon
         :closable="true"
         class="info-alert"
         description="Secret 用于存储敏感信息，例如密码、OAuth 令牌和 SSH 密钥。 将此类信息放在 Secret 中比放在 Pod 定义或容器镜像中更安全、更灵活。 Pod 可以通过卷挂载或环境变量的方式来使用 Secret。"
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
              placeholder="搜索 Secret 名称 / 类型"
              :prefix-icon="SearchIcon"
              clearable
              @input="handleSearchDebounced"
              class="filter-item search-input"
              :disabled="!selectedNamespace || loading.secrets"
          />
  
          <el-tooltip content="刷新列表" placement="top">
              <el-button
                :icon="RefreshIcon"
                circle
                @click="fetchSecretData"
                :loading="loading.secrets"
                :disabled="!selectedNamespace"
              />
          </el-tooltip>
      </div>
  
      <!-- Secrets Table -->
      <el-table
          :data="paginatedData"
          v-loading="loading.secrets"
          border
          stripe
          style="width: 100%"
          @sort-change="handleSortChange"
          class="secret-table"
          :default-sort="{ prop: 'createdAt', order: 'descending' }"
          row-key="uid"
      >
          <el-table-column prop="name" label="名称" min-width="300" sortable="custom" fixed show-overflow-tooltip>
               <template #default="{ row }">
                  <el-icon class="secret-icon"><LockIcon /></el-icon>
                  <span class="secret-name">{{ row.name }}</span>
              </template>
          </el-table-column>
          <el-table-column prop="namespace" label="命名空间" min-width="150" sortable="custom" show-overflow-tooltip />
          <el-table-column prop="type" label="类型" min-width="200" sortable="custom">
               <template #default="{ row }">
                   <el-tooltip :content="row.type" placement="top">
                      <el-tag size="small" type="info" effect="plain" class="type-tag">{{ formatSecretType(row.type) }}</el-tag>
                   </el-tooltip>
              </template>
          </el-table-column>
          <el-table-column prop="dataCount" label="数据条目 (Keys)" min-width="150" sortable="custom" align="center">
              <template #default="{ row }">
                  <el-tag size="small">{{ row.dataCount }}</el-tag>
              </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" min-width="180" sortable="custom" />
          <el-table-column label="操作" width="130" align="center" fixed="right">
              <template #default="{ row }">
                   <el-tooltip content="编辑 YAML" placement="top">
                      <el-button link type="primary" :icon="EditPenIcon" @click="editSecretYaml(row)" />
                  </el-tooltip>
                  <el-tooltip content="删除" placement="top">
                      <el-button link type="danger" :icon="DeleteIcon" @click="handleDeleteSecret(row)" />
                  </el-tooltip>
              </template>
          </el-table-column>
           <template #empty>
            <el-empty v-if="!loading.secrets && !selectedNamespace" description="请先选择一个命名空间以加载 Secrets" image-size="100" />
            <el-empty v-else-if="!loading.secrets && paginatedData.length === 0" :description="`在命名空间 '${selectedNamespace}' 中未找到 Secrets`" image-size="100" />
           </template>
      </el-table>
  
      <!-- Pagination -->
      <div class="pagination-container" v-if="!loading.secrets && totalSecrets > 0">
          <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="totalSecrets"
              layout="total, sizes, prev, pager, next, jumper"
              background
              @size-change="handleSizeChange"
              @current-change="handlePageChange"
              :disabled="loading.secrets"
          />
      </div>
  
      <!-- Create/Edit Dialog (YAML focus) -->
      <el-dialog :title="dialogTitle" v-model="dialogVisible" width="70%" :close-on-click-modal="false">
         <el-alert type="warning" :closable="false" style="margin-bottom: 20px;">
           <strong>安全警告:</strong> 编辑 Secret 时，`data` 字段中的值应为 Base64 编码格式。 `stringData` 字段允许使用明文字符串（Kubernetes 会自动编码）。请谨慎处理敏感信息。
         </el-alert>
         <!-- Integrate a YAML editor component here -->
         <div class="yaml-editor-placeholder">
              YAML 编辑器占位符
              <pre>{{ yamlContent || placeholderYaml }}</pre>
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
  import { request } from "@/utils/service" // Ensure correct path
  import dayjs from "dayjs"
  import { debounce } from 'lodash-es'
  import yaml from 'js-yaml'; // Ensure installed
  import { Base64 } from 'js-base64'; // For placeholder
  
  import {
      Plus as PlusIcon, Search as SearchIcon, Refresh as RefreshIcon, Lock as LockIcon, // Icon for Secret
      EditPen as EditPenIcon, Delete as DeleteIcon
  } from '@element-plus/icons-vue'
  
  // --- Interfaces ---
  // ** Adjusted to match a simplified/flattened API list response **
  interface SecretApiItem {
    name: string;
    namespace: string;
    uid: string;
    type: string; // K8s type is string enum
    dataCount: number; // Expect count directly from API list view
    createdAt: string;
    labels?: { [key: string]: string };
    annotations?: { [key: string]: string };
    resourceVersion: string;
    // ** List response likely DOES NOT contain data/stringData **
  }
  
  // Interface for the GET /name response (needs full data)
  // Assuming it returns a structure compatible with SecretDetailBackend structure
  interface K8sMetadata { name: string; namespace: string; uid?: string; resourceVersion?: string; creationTimestamp?: string; labels?: { [key: string]: string }; annotations?: { [key: string]: string }; }
  interface SecretDetailBackend {
      apiVersion?: string;
      kind?: string;
      metadata: K8sMetadata;
      type?: string;
      data?: { [key: string]: string }; // Base64 encoded strings
      stringData?: { [key: string]: string };
  }
  
  interface SecretListApiResponseData { items: SecretApiItem[]; total: number }
  interface SecretApiResponse { code: number; data: SecretListApiResponseData; message: string }
  // For GET /:name endpoint
  interface SecretDetailApiResponse { code: number; data: SecretDetailBackend; message: string } // Expecting detail structure
  interface NamespaceListResponse { code: number; data: string[]; message: string }
  
  // Internal Display/Table Item
  interface SecretDisplayItem {
    uid: string
    name: string
    namespace: string
    type: string
    dataCount: number
    createdAt: string
    // Store the simple list item data. Full data fetched on demand for edit.
    rawData: SecretApiItem
  }
  
  // --- Reactive State ---
  const allSecrets = ref<SecretDisplayItem[]>([])
  const namespaces = ref<string[]>([])
  const selectedNamespace = ref<string>("")
  const currentPage = ref(1)
  const pageSize = ref(10)
  const totalSecrets = ref(0)
  const searchQuery = ref("")
  const sortParams = reactive({ key: 'createdAt', order: 'descending' as ('ascending' | 'descending' | null) })
  
  const loading = reactive({
      page: false, namespaces: false, secrets: false, dialogSave: false
  })
  
  // Dialog state (YAML focus)
  const dialogVisible = ref(false)
  const dialogTitle = ref("创建 Secret (YAML)");
  // ** Store the DETAILED structure fetched via GET for editing **
  const currentEditSecretDetail = ref<SecretDetailBackend | null>(null);
  const yamlContent = ref("");
  const placeholderYaml = computed(() => `apiVersion: v1
  kind: Secret
  metadata:
    name: my-new-secret
    namespace: ${selectedNamespace.value || 'default'}
  # type: Opaque
  # stringData:
  #   app.conf: |-
  #     key1=value1
  #     key2=value2
  data:
    # Values must be base64 encoded
    user: ${Base64.encode('admin')}
    pass: ${Base64.encode('password123')}
  `);
  
  // --- Computed Properties ---
  const filteredData = computed(() => { /* ... as before ... */
      const query = searchQuery.value.trim().toLowerCase()
      if (!query) return allSecrets.value;
      return allSecrets.value.filter(s =>
          s.name.toLowerCase().includes(query) ||
          s.type.toLowerCase().includes(query)
      );
  });
  
  const sortedData = computed(() => { /* ... as before ... */
      const data = [...filteredData.value];
      const { key, order } = sortParams;
      if (!key || !order) return data;
  
      data.sort((a, b) => {
          let valA: any; let valB: any;
          if (key === 'dataCount') {
              valA = a.dataCount ?? 0; valB = b.dataCount ?? 0;
          } else if (key === 'createdAt') {
              const timeA = a.createdAt ? dayjs(a.createdAt, "YYYY-MM-DD HH:mm:ss").valueOf() : 0;
              const timeB = b.createdAt ? dayjs(b.createdAt, "YYYY-MM-DD HH:mm:ss").valueOf() : 0;
              valA = isNaN(timeA) ? 0 : timeA; valB = isNaN(timeB) ? 0 : timeB;
          } else {
              valA = a[key as keyof SecretDisplayItem] ?? ''; valB = b[key as keyof SecretDisplayItem] ?? '';
              if (typeof valA === 'string') valA = valA.toLowerCase();
              if (typeof valB === 'string') valB = valB.toLowerCase();
          }
          let comparison = 0;
          if (valA < valB) comparison = -1; else if (valA > valB) comparison = 1;
          return order === 'ascending' ? comparison : -comparison;
      });
      return data;
  });
  
  
  const paginatedData = computed(() => { /* ... as before ... */
      const start = (currentPage.value - 1) * pageSize.value;
      const end = start + pageSize.value;
      return sortedData.value.slice(start, end);
  });
  
  
  // --- Helper Functions ---
  const formatTimestamp = (timestamp: string): string => { /* ... */ if (!timestamp) return 'N/A'; return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss"); }
  const formatSecretType = (type: string | undefined): string => { /* ... */ if (!type) return 'Opaque'; if (type.startsWith('kubernetes.io/')) { return type.substring('kubernetes.io/'.length); } return type; };
  const VITE_API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://192.168.1.100:8080"; // Ensure this is set in your environment variables
  // --- API Interaction ---
  const fetchNamespaces = async () => { /* ... same as before ... */
      loading.namespaces = true;
      try {
          const response = await request<NamespaceListResponse>({ url: "/api/v1/namespaces", method: "get", baseURL: VITE_API_BASE_URL });
          if (response.code === 200 && Array.isArray(response.data)) {
              namespaces.value = response.data;
              if (namespaces.value.length > 0 && !selectedNamespace.value) {
                   selectedNamespace.value = namespaces.value.find(ns => ns === 'default') || namespaces.value[0];
              } else if (namespaces.value.length === 0) { ElMessage.warning("未找到任何命名空间。"); }
          } else { ElMessage.error(`获取命名空间失败: ${response.message || '数据格式错误'}`); namespaces.value = []; }
      } catch (error: any) { console.error("获取命名空间失败:", error); ElMessage.error(`获取命名空间失败: ${error.message || '网络请求失败'}`); namespaces.value = []; }
      finally { loading.namespaces = false; }
  }
  
  const fetchSecretData = async () => {
      if (!selectedNamespace.value) { allSecrets.value = []; totalSecrets.value = 0; return; }
      loading.secrets = true;
      allSecrets.value = [];
      totalSecrets.value = 0;
      try {
          const params: Record<string, any> = { /* Server-side params */ };
          const url = `/api/v1/namespaces/${selectedNamespace.value}/secrets`;
          // ** Expect the SIMPLIFIED list response structure **
          const response = await request<SecretApiResponse>({ url, method: "get", params, baseURL: VITE_API_BASE_URL });
  
          if (response.code === 200 && response.data?.items && Array.isArray(response.data.items)) {
              totalSecrets.value = response.data.total ?? response.data.items.length;
  
              // ** Map from the simplified SecretApiItem (list view doesn't have data) **
              allSecrets.value = response.data.items
                  .filter(item => item && item.name && item.namespace && item.uid)
                  .map((item, index) => ({
                      uid: item.uid,
                      name: item.name,
                      namespace: item.namespace,
                      type: item.type || 'Opaque',
                      dataCount: item.dataCount, // Directly use dataCount from list API
                      createdAt: formatTimestamp(item.createdAt),
                      rawData: item, // Store the basic list item data
              }));
  
              const totalPages = Math.ceil(totalSecrets.value / pageSize.value);
               if (currentPage.value > totalPages && totalPages > 0) currentPage.value = totalPages;
               else if (totalSecrets.value === 0) currentPage.value = 1;
  
          } else if (response.code === 200 && response.data?.items === null) {
               console.log(`No Secrets found in namespace '${selectedNamespace.value}' (items is null).`);
               allSecrets.value = []; totalSecrets.value = 0; currentPage.value = 1;
          } else {
              ElMessage.error(`获取 Secret 数据失败: ${response.message || '无效的响应数据'}`);
              allSecrets.value = []; totalSecrets.value = 0;
          }
      } catch (error: any) {
          console.error("获取 Secret 数据失败:", error);
          ElMessage.error(`获取 Secret 数据出错: ${error.message || '网络请求失败'}`);
          allSecrets.value = []; totalSecrets.value = 0;
      } finally {
          loading.secrets = false;
      }
  }
  
  // --- Event Handlers ---
  const handleNamespaceChange = () => { /* ... */ currentPage.value = 1; searchQuery.value = ''; sortParams.key = 'createdAt'; sortParams.order = 'descending'; fetchSecretData(); };
  const handlePageChange = (page: number) => { currentPage.value = page; /* Fetch if server-side */ };
  const handleSizeChange = (size: number) => { pageSize.value = size; currentPage.value = 1; /* Fetch if server-side */ };
  const handleSearchDebounced = debounce(() => { currentPage.value = 1; /* Fetch if server-side */ }, 300);
  const handleSortChange = ({ prop, order }: { prop: string | null; order: 'ascending' | 'descending' | null }) => { /* ... */ sortParams.key = prop || 'createdAt'; sortParams.order = order; currentPage.value = 1; };
  
  
  // --- Dialog and CRUD Actions ---
  const handleAddSecret = () => { /* ... */ if (!selectedNamespace.value) { ElMessage.warning("请先选择一个命名空间"); return; } currentEditSecret.value = null; yamlContent.value = placeholderYaml.value; dialogTitle.value = "创建 Secret (YAML)"; dialogVisible.value = true; };
  
  const editSecretYaml = async (secret: SecretDisplayItem) => {
      ElMessage.info(`获取 Secret "${secret.name}" 的详细信息...`);
      loading.secrets = true; // Indicate loading
      currentEditSecret.value = null; // Clear previous edit data
      yamlContent.value = ""; // Clear previous yaml
      try {
         // ** Fetch the FULL details using the GET endpoint **
         const response = await request<SecretDetailApiResponse>({
             url: `/api/v1/namespaces/${secret.namespace}/secrets/${secret.name}`,
             method: 'get',
             baseURL: VITE_API_BASE_URL,
         });
         if (response.code === 200 && response.data) {
              currentEditSecret.value = response.data; // Store full data from GET request
  
              // ** Reconstruct standard K8s object for YAML editor **
              //    (Ensure SecretDetailBackend matches K8s structure or adapt this reconstruction)
              const k8sObjectForYaml = {
                   apiVersion: response.data.apiVersion || "v1",
                   kind: response.data.kind || "Secret",
                   metadata: {
                       name: response.data.metadata.name,
                       namespace: response.data.metadata.namespace,
                       labels: response.data.metadata.labels,
                       annotations: response.data.metadata.annotations,
                       resourceVersion: response.data.metadata.resourceVersion,
                   },
                   type: response.data.type || 'Opaque',
                   // Use stringData if present (usually preferred for text), else use base64 data
                   stringData: response.data.stringData,
                   data: response.data.data,
               };
               // Remove empty fields before dumping for cleaner YAML
               if (!k8sObjectForYaml.stringData || Object.keys(k8sObjectForYaml.stringData).length === 0) delete k8sObjectForYaml.stringData;
               if (!k8sObjectForYaml.data || Object.keys(k8sObjectForYaml.data).length === 0) delete k8sObjectForYaml.data;
  
  
              yamlContent.value = yaml.dump(k8sObjectForYaml, { skipInvalid: true, sortKeys: false }); // Don't sort keys for K8s objects
              dialogTitle.value = `编辑 Secret: ${secret.name} (YAML)`;
              dialogVisible.value = true;
         } else {
             ElMessage.error(`获取 Secret 详情失败: ${response.message || '未知错误'}`);
         }
      } catch (error: any) {
          console.error("获取 Secret 详情失败:", error);
          ElMessage.error(`获取 Secret 详情出错: ${error.message || '网络请求失败'}`);
      } finally {
         loading.secrets = false;
      }
  };
  
  const handleSaveYaml = async () => { /* ... */
      // Determine namespace
      const namespaceToUse = currentEditSecret.value?.metadata?.namespace || selectedNamespace.value;
      if (!namespaceToUse) { ElMessage.error("无法确定目标命名空间。"); return; }
  
      loading.dialogSave = true;
      // --- Replace with actual YAML editor interaction and API call ---
      // const currentYaml = yamlEditorRef.value.getContent();
      // try {
      //     let parsedYaml = yaml.load(currentYaml); ... validate ...
      //     const name = parsedYaml.metadata.name;
      //     const method = currentEditSecret.value ? 'put' : 'post';
      //     const url = currentEditSecret.value ? `/api/v1/namespaces/${namespaceToUse}/secrets/${name}` : `/api/v1/namespaces/${namespaceToUse}/secrets`;
      //     // Backend needs to handle the parsed JSON object (corev1.Secret)
      //     // Ensure frontend YAML editor content parses correctly to required Go struct
      //     const response = await request({ url, method, data: parsedYaml, baseURL:"..." });
      //     ... handle response ...
      // } catch (e) { ... }
  
       await new Promise(resolve => setTimeout(resolve, 500)); // Simulate
       loading.dialogSave = false; dialogVisible.value = false;
       const action = currentEditSecret.value ? '更新' : '创建';
       ElMessage.success(`模拟 Secret ${action}成功`); fetchSecretData();
  };
  
  const handleDeleteSecret = (secret: SecretDisplayItem) => { /* ... */
      ElMessageBox.confirm(
          `确定要删除 Secret "${secret.name}" (命名空间: ${secret.namespace}) 吗？`,
          '确认删除', { type: 'warning' }
      ).then(async () => {
          loading.secrets = true;
          try {
              const response = await request<{ code: number; message: string }>({
                  url: `/api/v1/namespaces/${secret.namespace}/secrets/${secret.name}`,
                  method: "delete",
                  baseURL: VITE_API_BASE_URL,
              });
               if (response.code === 200 || response.code === 204 || response.code === 202) {
                  ElMessage.success(`Secret "${secret.name}" 已删除`); await fetchSecretData();
              } else { ElMessage.error(`删除 Secret 失败: ${response.message || '未知错误'}`); loading.secrets = false; }
          } catch (error: any) { console.error("删除 Secret 失败:", error); ElMessage.error(`删除 Secret 失败: ${error.response?.data?.message || error.message || '请求失败'}`); loading.secrets = false; }
      }).catch(() => ElMessage.info('删除操作已取消'));
  };
  
  // --- Lifecycle Hooks ---
  onMounted(async () => { /* ... */
      loading.page = true;
      await fetchNamespaces();
      if (selectedNamespace.value) { await fetchSecretData(); }
      loading.page = false;
  });
  
  </script>
  
  <style lang="scss" scoped>
  /* Using fallback variables directly */
  $page-padding: 20px; $spacing-md: 15px; $spacing-lg: 20px;
  $font-size-base: 14px; $font-size-small: 12px; $font-size-large: 16px; $font-size-extra-large: 24px;
  $border-radius-base: 4px; $kube-secret-icon-color: #e74c3c;
  
  .secret-page-container { padding: $page-padding; background-color: var(--el-bg-color-page); }
  .page-breadcrumb { margin-bottom: $spacing-lg; }
  .page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: $spacing-md; flex-wrap: wrap; gap: $spacing-md; }
  .page-title { font-size: $font-size-extra-large; font-weight: 600; color: var(--el-text-color-primary); margin: 0; }
  .info-alert { margin-bottom: $spacing-lg; background-color: var(--el-color-warning-light-9); border-color: var(--el-color-warning-light-8); :deep(.el-alert__title) { color: var(--el-color-warning-dark-2); } :deep(.el-alert__description) { font-size: $font-size-small; color: var(--el-text-color-regular); line-height: 1.6; } }
  .filter-bar { display: flex; align-items: center; flex-wrap: wrap; gap: $spacing-md; margin-bottom: $spacing-lg; padding: $spacing-md; background-color: var(--el-bg-color); border-radius: $border-radius-base; border: 1px solid var(--el-border-color-lighter); }
  .filter-item { }
  .namespace-select { width: 240px; }
  .search-input { width: 300px; }
  
  .secret-table {
     border-radius: $border-radius-base; border: 1px solid var(--el-border-color-lighter); overflow: hidden;
      :deep(th.el-table__cell) { background-color: var(--el-fill-color-lighter); color: var(--el-text-color-secondary); font-weight: 600; font-size: $font-size-small; }
      :deep(td.el-table__cell) { padding: 8px 0; font-size: $font-size-base; vertical-align: middle; }
     .secret-icon { margin-right: 8px; color: $kube-secret-icon-color; vertical-align: middle; font-size: 18px; position: relative; top: -1px; }
     .secret-name { font-weight: 500; vertical-align: middle; color: var(--el-text-color-regular); }
     .type-tag { max-width: 180px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  }
  
  .el-table .el-button.is-link { font-size: 14px; padding: 4px; margin: 0 3px; vertical-align: middle; }
  .pagination-container { display: flex; justify-content: flex-end; margin-top: $spacing-lg; }
  .yaml-editor-placeholder { border: 1px dashed var(--el-border-color); padding: 20px; margin-top: 10px; min-height: 350px; max-height: 60vh; background-color: var(--el-fill-color-lighter); color: var(--el-text-color-secondary); font-family: monospace; white-space: pre-wrap; overflow: auto; font-size: 13px; line-height: 1.5; }
  .dialog-footer { text-align: right; padding-top: 10px; }
  </style>