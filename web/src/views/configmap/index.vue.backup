<template>
    <div class="configmap-page-container">
      <!-- Header: Title & Create Button -->
      <div class="page-header">
        <h1 class="page-title">{{ $t('configmapManagement.title') }}</h1>
         <el-button
           type="primary"
           :icon="PlusIcon"
           @click="handleAddCM"
           :loading="loading.page"
           :disabled="!selectedNamespace"
         >
           {{ $t('configmapManagement.actions.create') }}
         </el-button>
      </div>
  
      <!-- Filter Bar: Namespace, Search, Refresh -->
      <div class="filter-bar">
          <el-select
              v-model="selectedNamespace"
              :placeholder="$t('configmapManagement.filterNamespace')"
              @change="handleNamespaceChange"
              filterable
              :loading="loading.namespaces"
              class="filter-item namespace-select"
              style="width: 280px;"
          >
              <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
               <template #empty>
                  <div style="padding: 10px; text-align: center; color: #999;">
                      {{ loading.namespaces ? $t('common.loading') : $t('configmapManagement.messages.noNamespaces') }}
                  </div>
              </template>
          </el-select>
  
          <el-input
              v-model="searchQuery"
              :placeholder="$t('configmapManagement.searchPlaceholder')"
              :prefix-icon="SearchIcon"
              clearable
              @input="handleSearchDebounced"
              class="filter-item search-input"
              :disabled="!selectedNamespace || loading.configmaps"
          />
  
          <el-tooltip :content="$t('common.refresh')" placement="top">
              <el-button
                :icon="RefreshIcon"
                circle
                @click="fetchConfigMapData"
                :loading="loading.configmaps"
                :disabled="!selectedNamespace"
              />
          </el-tooltip>
      </div>
  
      <!-- ConfigMaps Table Card -->
      <el-card class="configmap-table-card" shadow="hover">
        <el-table
            :data="paginatedData"
            v-loading="loading.configmaps"
            border
            stripe
            style="width: 100%"
            @sort-change="handleSortChange"
            class="configmap-table"
            :default-sort="{ prop: 'createdAt', order: 'descending' }"
            row-key="uid"
        >
          <el-table-column :label="$t('configmapManagement.columns.name')" prop="name" min-width="300" sortable="custom" fixed show-overflow-tooltip>
               <template #default="{ row }">
                  <el-icon class="cm-icon"><TicketsIcon /></el-icon>
                  <span class="cm-name">{{ row.name }}</span>
              </template>
          </el-table-column>
          <el-table-column :label="$t('configmapManagement.columns.namespace')" prop="namespace" min-width="150" sortable="custom" show-overflow-tooltip />
          <el-table-column :label="$t('configmapManagement.columns.dataCount')" prop="dataCount" min-width="150" sortable="custom" align="center">
              <template #default="{ row }">
                  <el-tag size="small">{{ row.dataCount }}</el-tag>
              </template>
          </el-table-column>
          <el-table-column :label="$t('configmapManagement.columns.createdAt')" prop="createdAt" min-width="180" sortable="custom" />
          <el-table-column :label="$t('configmapManagement.columns.actions')" width="130" align="center" fixed="right">
              <template #default="{ row }">
                   <el-tooltip :content="$t('configmapManagement.actions.editYaml')" placement="top">
                      <el-button link type="primary" :icon="EditPenIcon" @click="editCmYaml(row)" />
                  </el-tooltip>
                  <el-tooltip :content="$t('common.delete')" placement="top">
                      <el-button link type="danger" :icon="DeleteIcon" @click="handleDeleteCM(row)" />
                  </el-tooltip>
              </template>
          </el-table-column>
           <template #empty>
            <el-empty v-if="!loading.configmaps && !selectedNamespace" :description="$t('configmapManagement.messages.selectNamespace')" image-size="100" />
            <el-empty v-else-if="!loading.configmaps && paginatedData.length === 0" :description="$t('configmapManagement.messages.noConfigMaps')" image-size="100" />
           </template>
      </el-table>

        <!-- Pagination -->
        <div class="pagination-container" v-if="!loading.configmaps && totalConfigMaps > 0">
          <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="totalConfigMaps"
              layout="total, sizes, prev, pager, next, jumper"
              background
              @size-change="handleSizeChange"
              @current-change="handlePageChange"
              :disabled="loading.configmaps"
          />
        </div>
      </el-card>
  
      <!-- Create/Edit Dialog (YAML focus) -->
      <el-dialog :title="dialogTitle" v-model="dialogVisible" width="70%" :close-on-click-modal="false">
         <el-alert type="info" :closable="false" style="margin-bottom: 20px;">
           {{ $t('configmapManagement.messages.yamlEditTip') }}
         </el-alert>
         <!-- Integrate a YAML editor component here -->
         <div class="yaml-editor-placeholder">
              {{ $t('configmapManagement.messages.yamlEditorPlaceholder') }}
              <pre>{{ yamlContent || placeholderYaml }}</pre>
         </div>
        <template #footer>
          <div class="dialog-footer">
              <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
              <el-button type="primary" @click="handleSaveYaml" :loading="loading.dialogSave">{{ $t('configmapManagement.actions.applyYaml') }}</el-button>
          </div>
        </template>
      </el-dialog>
  
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, reactive, computed, onMounted } from "vue"
  import { ElMessage, ElMessageBox } from "element-plus"
  import { useI18n } from 'vue-i18n'
  import { kubernetesRequest, fetchNamespaces, KubernetesApiResponse } from "@/utils/api-config" // Ensure correct path
  import dayjs from "dayjs"
  import { debounce } from 'lodash-es'
  import yaml from 'js-yaml'; // Ensure installed
  
  import {
      Plus as PlusIcon, Search as SearchIcon, Refresh as RefreshIcon, Tickets as TicketsIcon, // Icon for CM
      EditPen as EditPenIcon, Delete as DeleteIcon
  } from '@element-plus/icons-vue'

  const { t } = useI18n()
  // --- Interfaces ---
  // ** Adjusted to match a simplified API response (like the corrected PVC) **
  interface ConfigMapApiItem {
    name: string;
    namespace: string;
    uid: string;
    dataCount: number; // Expect count directly from API list view
    createdAt: string;
    labels?: { [key: string]: string };
    annotations?: { [key: string]: string };
    resourceVersion: string;
    // Include data/binaryData only if the LIST API actually returns them
    // data?: { [key: string]: string };
    // binaryData?: { [key: string]: string };
  }
  
  // Detail response still needed if GET /name returns full data
  interface ConfigMapDetailApiItem {
      apiVersion?: string; // Added for reconstructing K8s object
      kind?: string;       // Added for reconstructing K8s object
      metadata: K8sMetadata; // Reusing K8sMetadata from previous examples
      data?: { [key: string]: string };
      binaryData?: { [key: string]: string }; // Base64 encoded string
  }
  
  interface ConfigMapListApiResponseData { items: ConfigMapApiItem[]; total: number }
  interface ConfigMapApiResponse { code: number; data: ConfigMapListApiResponseData; message: string }
  // For GET /:name endpoint returning full K8s object (or detail structure)
  interface ConfigMapDetailApiResponse { code: number; data: ConfigMapDetailApiItem; message: string }
  interface NamespaceListResponse { code: number; data: { items: Array<{ metadata: { name: string } }> }; message: string }
  
  // Internal Display/Table Item - Simplified to match expected API list item
  interface ConfigMapDisplayItem {
    uid: string
    name: string
    namespace: string
    dataCount: number
    createdAt: string
    // Store raw API list item. For edit, we'll need to fetch details.
    rawData: ConfigMapApiItem
  }
  
  // Interface for K8s metadata (reuse if defined elsewhere)
  interface K8sMetadata { name: string; namespace: string; uid?: string; resourceVersion?: string; creationTimestamp?: string; labels?: { [key: string]: string }; annotations?: { [key: string]: string }; }
  
  
  // --- Reactive State ---
  const allConfigMaps = ref<ConfigMapDisplayItem[]>([])
  const namespaces = ref<string[]>([])
  const selectedNamespace = ref<string>("")
  const currentPage = ref(1)
  const pageSize = ref(10)
  const totalConfigMaps = ref(0)
  const searchQuery = ref("")
  const sortParams = reactive({ key: 'createdAt', order: 'descending' as ('ascending' | 'descending' | null) })
  
  const loading = reactive({
      page: false, namespaces: false, configmaps: false, dialogSave: false
  })
  
  // Dialog state (YAML focus)
  const dialogVisible = ref(false)
  const currentEditCm = ref<ConfigMapApiItem | null>(null); // Store basic info from list
  const yamlContent = ref("");

  const dialogTitle = computed(() => 
    currentEditCm.value 
      ? `${t('configmapManagement.actions.edit')}: ${currentEditCm.value.name} (YAML)`
      : t('configmapManagement.messages.createConfigMapTitle')
  );

  const placeholderYaml = computed(() => `apiVersion: v1
kind: ConfigMap
metadata:
  name: my-new-configmap
  namespace: ${selectedNamespace.value || 'default'}
data:
  key1: value1
  key2: |
    multi-line
    value
  `);
  
  // --- Computed Properties ---
  const filteredData = computed(() => { /* ... as before ... */
      const query = searchQuery.value.trim().toLowerCase()
      if (!query) return allConfigMaps.value;
      return allConfigMaps.value.filter(cm => cm.name.toLowerCase().includes(query));
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
              valA = a[key as keyof ConfigMapDisplayItem] ?? ''; valB = b[key as keyof ConfigMapDisplayItem] ?? '';
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
  
  // --- API Interaction ---
  const fetchNamespaces = async () => { /* ... same as before ... */
      loading.namespaces = true;
      try {
          const response = await kubernetesRequest<NamespaceListResponse>({ url: "/api/v1/namespaces", method: "get"});
          if (response.code === 200 && response.data && Array.isArray(response.data.items)) {
              namespaces.value = response.data.items.map(ns => ns.metadata.name);
              if (namespaces.value.length > 0 && !selectedNamespace.value) {
                   selectedNamespace.value = namespaces.value.find(ns => ns === 'default') || namespaces.value[0];
              } else if (namespaces.value.length === 0) { ElMessage.warning(t('configmapManagement.messages.noNamespaces')); }
          } else { ElMessage.error(`${t('common.error')}: ${response.message || '数据格式错误'}`); namespaces.value = []; }
      } catch (error: any) { console.error("获取命名空间失败:", error); ElMessage.error(`${t('common.error')}: ${error.message || '网络请求失败'}`); namespaces.value = []; }
      finally { loading.namespaces = false; }
  }
  
  const fetchConfigMapData = async () => {
      if (!selectedNamespace.value) { allConfigMaps.value = []; totalConfigMaps.value = 0; return; }
      loading.configmaps = true;
      allConfigMaps.value = [];
      totalConfigMaps.value = 0;
      try {
          const params: Record<string, any> = { /* Server-side params */ };
          const url = `/api/v1/namespaces/${selectedNamespace.value}/configmaps`;
          // ** Use the interface matching the SIMPLIFIED API list response **
          const response = await kubernetesRequest<ConfigMapApiResponse>({ url, method: "get", params });
  
          if (response.code === 200 && response.data?.items && Array.isArray(response.data.items)) {
              totalConfigMaps.value = response.data.total ?? response.data.items.length;
  
              // ** Map directly from the simplified ConfigMapApiItem structure **
              allConfigMaps.value = response.data.items
                  .filter(item => item && item.name && item.namespace && item.uid) // Ensure basic fields exist
                  .map((item, index) => ({
                      uid: item.uid, // UID is expected directly
                      name: item.name,
                      namespace: item.namespace,
                      dataCount: item.dataCount, // dataCount is expected directly
                      createdAt: formatTimestamp(item.createdAt), // createdAt is expected directly
                      rawData: item, // Store the simplified list item data
              }));
  
              const totalPages = Math.ceil(totalConfigMaps.value / pageSize.value);
               if (currentPage.value > totalPages && totalPages > 0) currentPage.value = totalPages;
               else if (totalConfigMaps.value === 0) currentPage.value = 1;
          } else if (response.code === 200 && response.data?.items === null) {
               console.log(`No ConfigMaps found in namespace '${selectedNamespace.value}' (items is null).`);
               allConfigMaps.value = []; totalConfigMaps.value = 0; currentPage.value = 1;
          } else {
              ElMessage.error(`${t('configmapManagement.messages.fetchConfigMapsError')}: ${response.message || '无效的响应数据'}`);
              allConfigMaps.value = []; totalConfigMaps.value = 0;
          }
      } catch (error: any) {
          console.error("获取 ConfigMap 数据失败:", error);
          ElMessage.error(`${t('configmapManagement.messages.fetchConfigMapsError')}: ${error.message || '网络请求失败'}`);
          allConfigMaps.value = []; totalConfigMaps.value = 0;
      } finally {
          loading.configmaps = false;
      }
  }
  
  // --- Event Handlers ---
  const handleNamespaceChange = () => { /* ... */ currentPage.value = 1; searchQuery.value = ''; sortParams.key = 'createdAt'; sortParams.order = 'descending'; fetchConfigMapData(); };
  const handlePageChange = (page: number) => { currentPage.value = page; /* Fetch if server-side */ };
  const handleSizeChange = (size: number) => { pageSize.value = size; currentPage.value = 1; /* Fetch if server-side */ };
  const handleSearchDebounced = debounce(() => { currentPage.value = 1; /* Fetch if server-side */ }, 300);
  const handleSortChange = ({ prop, order }: { prop: string | null; order: 'ascending' | 'descending' | null }) => { /* ... */ sortParams.key = prop || 'createdAt'; sortParams.order = order; currentPage.value = 1; };
  
  
  // --- Dialog and CRUD Actions ---
  const handleAddCM = () => { /* ... */ if (!selectedNamespace.value) { ElMessage.warning("请先选择一个命名空间"); return; } currentEditCm.value = null; yamlContent.value = placeholderYaml.value; dialogTitle.value = "创建 ConfigMap (YAML)"; dialogVisible.value = true; };
  
  const editCmYaml = async (cm: ConfigMapDisplayItem) => {
      ElMessage.info(`获取 ConfigMap "${cm.name}" 的详细信息...`);
      loading.configmaps = true; // Indicate loading for the fetch
      try {
         // ** Fetch the FULL details using the GET endpoint **
         const response = await kubernetesRequest<ConfigMapDetailApiResponse>({
             url: `/api/v1/namespaces/${cm.namespace}/configmaps/${cm.name}`,
             method: 'get'
         });
         if (response.code === 200 && response.data) {
              currentEditCm.value = response.data; // Store full data from GET request
              // ** Reconstruct the standard K8s object for YAML editing **
              const k8sObjectForYaml = {
                   apiVersion: "v1", // Assuming v1
                   kind: "ConfigMap", // Assuming Kind
                   metadata: {
                       name: response.data.name,
                       namespace: response.data.namespace,
                       labels: response.data.labels,
                       annotations: response.data.annotations,
                       resourceVersion: response.data.resourceVersion, // Include RV for potential update
                       // Exclude uid, creationTimestamp, managedFields etc. for editing
                   },
                   data: response.data.data,
                   binaryData: response.data.binaryData, // Include binaryData if present
               };
              yamlContent.value = yaml.dump(k8sObjectForYaml, { skipInvalid: true }); // Dump the reconstructed object
              dialogTitle.value = `编辑 ConfigMap: ${cm.name} (YAML)`;
              dialogVisible.value = true;
         } else {
             ElMessage.error(`获取 ConfigMap 详情失败: ${response.message || '未知错误'}`);
         }
      } catch (error: any) {
          console.error("获取 ConfigMap 详情失败:", error);
          ElMessage.error(`获取 ConfigMap 详情出错: ${error.message || '网络请求失败'}`);
      } finally {
         loading.configmaps = false;
      }
  };
  
  
  const handleSaveYaml = async () => { /* ... */
      if (!selectedNamespace.value && !currentEditCm.value?.namespace) { ElMessage.error("无法确定目标命名空间。"); return; }
      loading.dialogSave = true;
      // --- Replace with actual YAML editor interaction and API call ---
      // const currentYaml = yamlEditorRef.value.getContent();
      // try {
      //     let parsedYaml = yaml.load(currentYaml); ... validate ...
      //     const name = parsedYaml.metadata.name;
      //     const namespaceToUse = parsedYaml.metadata.namespace || currentEditCm.value?.namespace || selectedNamespace.value;
      //     const method = currentEditCm.value ? 'put' : 'post';
      //     const url = currentEditCm.value ? `/api/v1/namespaces/${namespaceToUse}/configmaps/${name}` : `/api/v1/namespaces/${namespaceToUse}/configmaps`;
      //     // Backend expects corev1.ConfigMap structure
      //     const response = await request({ url, method, data: parsedYaml, baseURL:"..." });
      //     ... handle response ...
      // } catch (e) { ... }
  
       await new Promise(resolve => setTimeout(resolve, 500)); // Simulate
       loading.dialogSave = false; dialogVisible.value = false;
       const action = currentEditCm.value ? '更新' : '创建';
       ElMessage.success(`模拟 ConfigMap ${action}成功`); fetchConfigMapData();
  };
  
  const handleDeleteCM = (cm: ConfigMapDisplayItem) => { /* ... */
      ElMessageBox.confirm(
          `确定要删除 ConfigMap "${cm.name}" (命名空间: ${cm.namespace}) 吗？`,
          '确认删除', { type: 'warning' }
      ).then(async () => {
          loading.configmaps = true;
          try {
              const response = await kubernetesRequest<{ code: number; message: string }>({
                  url: `/api/v1/namespaces/${cm.namespace}/configmaps/${cm.name}`,
                  method: "delete"
              });
               if (response.code === 200 || response.code === 204 || response.code === 202) {
                  ElMessage.success(`ConfigMap "${cm.name}" 已删除`); await fetchConfigMapData();
              } else { ElMessage.error(`删除 ConfigMap 失败: ${response.message || '未知错误'}`); loading.configmaps = false; }
          } catch (error: any) { console.error("删除 ConfigMap 失败:", error); ElMessage.error(`删除 ConfigMap 失败: ${error.response?.data?.message || error.message || '请求失败'}`); loading.configmaps = false; }
      }).catch(() => ElMessage.info('删除操作已取消'));
  };
  
  // --- Lifecycle Hooks ---
  onMounted(async () => { /* ... */
      loading.page = true;
      await fetchNamespaces();
      if (selectedNamespace.value) { await fetchConfigMapData(); }
      loading.page = false;
  });
  
  </script>
  
<style lang="scss" scoped>
/* Using fallback variables directly */
$page-padding: 24px; $spacing-md: 16px; $spacing-lg: 24px;
$font-size-base: 14px; $font-size-small: 12px; $font-size-large: 16px; $font-size-extra-large: 28px;
$border-radius-base: 4px; $kube-cm-icon-color: #2ecc71;

.configmap-page-container { 
  padding: $page-padding; 
  background-color: var(--el-bg-color-page);
  box-sizing: border-box;
}
.page-header { 
  display: flex; 
  justify-content: space-between; 
  align-items: center; 
  margin-bottom: $spacing-lg; 
  flex-wrap: wrap; 
  gap: $spacing-md;
}
.page-title { 
  font-size: $font-size-extra-large; 
  font-weight: 600; 
  color: var(--el-text-color-primary); 
  margin: 0;
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
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
}
.filter-item { }
.namespace-select { width: 240px; }
.search-input { width: 300px; }

.configmap-table-card {
  margin-bottom: $spacing-lg;
  border-radius: $border-radius-base;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
  border: 1px solid var(--el-border-color-lighter);
  
  :deep(.el-card__body) {
    padding: 0;
  }
}

.configmap-table {
  border: none;
  border-radius: 0;
  overflow: hidden;
  :deep(th.el-table__cell) { 
    background-color: var(--el-fill-color-lighter); 
    color: var(--el-text-color-secondary); 
    font-weight: 600; 
    font-size: $font-size-small;
    border-bottom: 1px solid var(--el-border-color-light);
  }
  :deep(td.el-table__cell) { 
    padding: 12px 0; 
    font-size: $font-size-base; 
    vertical-align: middle;
    border-bottom: 1px solid var(--el-border-color-lighter);
  }
  :deep(.el-table__inner-wrapper::before) {
    height: 0;
  }
  .cm-icon { 
    margin-right: 8px; 
    color: $kube-cm-icon-color; 
    vertical-align: middle; 
    font-size: 18px; 
    position: relative; 
    top: -1px; 
  }
  .cm-name { 
    font-weight: 500; 
    vertical-align: middle; 
    color: var(--el-text-color-primary); 
  }
  .status-tag { 
    display: inline-flex; 
    align-items: center; 
    gap: 4px; 
    padding: 2px 8px; 
    height: 24px; 
    line-height: 20px; 
    font-size: $font-size-small;
    border-radius: 12px;
  }
  .status-icon { font-size: 12px; }
  .is-loading { animation: rotating 1.5s linear infinite; }
  @keyframes rotating { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

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
    padding: 0 16px 16px 16px;
  }
  .yaml-editor-placeholder { 
    border: 1px dashed var(--el-border-color); 
    padding: 20px; 
    margin-top: 10px; 
    min-height: 350px; 
    max-height: 60vh; 
    background-color: var(--el-fill-color-lighter); 
    color: var(--el-text-color-secondary); 
    font-family: monospace; 
    white-space: pre-wrap; 
    overflow: auto; 
    font-size: 13px; 
    line-height: 1.5; 
    border-radius: $border-radius-base;
  }
  .dialog-footer { 
    text-align: right; 
    padding-top: 16px;
  }
</style>