<template>
    <div class="secret-page-container">
      <!-- Breadcrumbs -->
      <el-breadcrumb separator="/" class="page-breadcrumb">
        <el-breadcrumb-item :to="{ path: '/' }">{{ t('nav.home') }}</el-breadcrumb-item>
        <el-breadcrumb-item>{{ t('menu.config') }}</el-breadcrumb-item>
        <el-breadcrumb-item>{{ t('menu.secret') }}</el-breadcrumb-item>
      </el-breadcrumb>
  
      <!-- Header: Title & Create Button -->
      <div class="page-header">
        <h1 class="page-title">{{ t('secretManagement.title') }}</h1>
         <el-button
           type="primary"
           :icon="PlusIcon"
           @click="handleAddSecret"
           :loading="loading.page"
           :disabled="!selectedNamespace"
         >
           {{ t('secretManagement.create') }}
         </el-button>
      </div>
  
      <!-- Filter Bar: Namespace, Search, Refresh -->
      <div class="filter-bar">
          <el-select
              v-model="selectedNamespace"
              :placeholder="t('secretManagement.selectNamespace')"
              @change="handleNamespaceChange"
              filterable
              :loading="loading.namespaces"
              class="filter-item namespace-select"
              style="width: 280px;"
          >
              <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
               <template #empty>
                  <div style="padding: 10px; text-align: center; color: #999;">
                      {{ loading.namespaces ? t('common.loading') : t('secretManagement.noNamespace') }}
                  </div>
              </template>
          </el-select>
  
          <el-input
              v-model="searchQuery"
              :placeholder="t('secretManagement.searchPlaceholder')"
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
          <el-table-column prop="name" :label="t('common.name')" min-width="300" sortable="custom" fixed show-overflow-tooltip>
               <template #default="{ row }">
                  <el-icon class="secret-icon"><LockIcon /></el-icon>
                  <span class="secret-name">{{ row.name }}</span>
              </template>
          </el-table-column>
          <el-table-column prop="namespace" :label="t('common.namespace')" min-width="150" sortable="custom" show-overflow-tooltip />
          <el-table-column prop="type" :label="t('secretManagement.type')" min-width="200" sortable="custom">
               <template #default="{ row }">
                   <el-tooltip :content="row.type" placement="top">
                      <el-tag size="small" type="info" effect="plain" class="type-tag">{{ formatSecretType(row.type) }}</el-tag>
                   </el-tooltip>
              </template>
          </el-table-column>
          <el-table-column prop="dataCount" :label="t('secretManagement.dataKeys')" min-width="150" sortable="custom" align="center">
              <template #default="{ row }">
                  <el-tag size="small">{{ row.dataCount }}</el-tag>
              </template>
          </el-table-column>
          <el-table-column prop="createdAt" :label="t('common.createdAt')" min-width="180" sortable="custom" />
          <el-table-column :label="t('common.actions')" width="130" align="center" fixed="right">
              <template #default="{ row }">
                   <el-tooltip :content="t('secretManagement.editYaml')" placement="top">
                      <el-button link type="primary" :icon="EditPenIcon" @click="editSecretYaml(row)" />
                  </el-tooltip>
                  <el-tooltip :content="t('common.delete')" placement="top">
                      <el-button link type="danger" :icon="DeleteIcon" @click="handleDeleteSecret(row)" />
                  </el-tooltip>
              </template>
          </el-table-column>
           <template #empty>
            <el-empty v-if="!loading.secrets && !selectedNamespace" :description="t('secretManagement.selectNamespace')" image-size="100" />
            <el-empty v-else-if="!loading.secrets && paginatedData.length === 0" :description="`${t('secretManagement.noSecretsInNamespace')} '${selectedNamespace}'`" image-size="100" />
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
           {{ t('secretManagement.securityWarning') }}
         </el-alert>
         <!-- Integrate a YAML editor component here -->
         <div class="yaml-editor-placeholder">
              {{ t('secretManagement.yamlEditorPlaceholder') }}
              <pre>{{ yamlContent || placeholderYaml }}</pre>
         </div>
        <template #footer>
          <div class="dialog-footer">
              <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
              <el-button type="primary" @click="handleSaveYaml" :loading="loading.dialogSave">{{ t('secretManagement.applyYaml') }}</el-button>
          </div>
        </template>
      </el-dialog>
  
    </div>
  </template>
  
<script setup lang="ts">
import { ref, reactive, computed, onMounted } from "vue"
import { useI18n } from "vue-i18n"
import { ElMessage, ElMessageBox } from "element-plus"
import { kubernetesRequest, fetchNamespaces, KubernetesApiResponse } from "@/utils/api-config" // Ensure correct path
import dayjs from "dayjs"
import { debounce } from 'lodash-es'
import yaml from 'js-yaml'; // Ensure installed
import { Base64 } from 'js-base64'; // For placeholder

import {
    Plus as PlusIcon, Search as SearchIcon, Refresh as RefreshIcon, Lock as LockIcon, // Icon for Secret
    EditPen as EditPenIcon, Delete as DeleteIcon
} from '@element-plus/icons-vue'

// 设置国际化
const { t } = useI18n()
  
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
  interface NamespaceListResponse { code: number; data: { items: Array<{ metadata: { name: string } }> }; message: string }
  
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
  const dialogTitle = ref(t('secretManagement.createSecretTitle'));
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
  // --- API Interaction ---
  const fetchNamespaces = async () => { /* ... same as before ... */
      loading.namespaces = true;
      try {
          const response = await kubernetesRequest<NamespaceListResponse>({ url: "/api/v1/namespaces", method: "get" });
          if (response.code === 200 && response.data && Array.isArray(response.data.items)) {
              namespaces.value = response.data.items.map(ns => ns.metadata.name);
              if (namespaces.value.length > 0 && !selectedNamespace.value) {
                   selectedNamespace.value = namespaces.value.find(ns => ns === 'default') || namespaces.value[0];
              } else if (namespaces.value.length === 0) { ElMessage.warning(t('secretManagement.noNamespaces')); }
          } else { ElMessage.error(`${t('secretManagement.fetchNamespacesFailed')}: ${response.message || t('secretManagement.dataFormatError')}`); namespaces.value = []; }
      } catch (error: any) { console.error("获取命名空间失败:", error); ElMessage.error(`${t('secretManagement.fetchNamespacesFailed')}: ${error.message || t('secretManagement.networkError')}`); namespaces.value = []; }
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
          const response = await kubernetesRequest<SecretApiResponse>({ url, method: "get", params });
  
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
              ElMessage.error(`${t('secretManagement.fetchSecretsError')}: ${response.message || t('secretManagement.invalidResponse')}`);
              allSecrets.value = []; totalSecrets.value = 0;
          }
      } catch (error: any) {
          console.error("获取 Secret 数据失败:", error);
          ElMessage.error(`${t('secretManagement.fetchSecretsError')}: ${error.message || t('secretManagement.networkError')}`);
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
  const handleAddSecret = () => { /* ... */ if (!selectedNamespace.value) { ElMessage.warning(t('secretManagement.selectNamespace')); return; } currentEditSecretDetail.value = null; yamlContent.value = placeholderYaml.value; dialogTitle.value = t('secretManagement.createSecretTitle'); dialogVisible.value = true; };
  
  const editSecretYaml = async (secret: SecretDisplayItem) => {
      ElMessage.info(t('secretManagement.fetchingYaml', { name: secret.name }));
      loading.secrets = true; // Indicate loading
      currentEditSecretDetail.value = null; // Clear previous edit data
      yamlContent.value = ""; // Clear previous yaml
      try {
         // ** Fetch the FULL details using the GET endpoint **
         const response = await kubernetesRequest<SecretDetailApiResponse>({
             url: `/api/v1/namespaces/${secret.namespace}/secrets/${secret.name}`,
             method: 'get'
         });
         if (response.code === 200 && response.data) {
              currentEditSecretDetail.value = response.data; // Store full data from GET request
  
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
              dialogTitle.value = `${t('secretManagement.editSecretTitle')}: ${secret.name} (YAML)`;
              dialogVisible.value = true;
         } else {
             ElMessage.error(`${t('secretManagement.fetchSecretDetailsFailed')}: ${response.message || t('secretManagement.unknownError')}`);
         }
      } catch (error: any) {
          console.error("获取 Secret 详情失败:", error);
          ElMessage.error(`${t('secretManagement.fetchSecretDetailsFailed')}: ${error.message || t('secretManagement.networkError')}`);
      } finally {
         loading.secrets = false;
      }
  };
  
  const handleSaveYaml = async () => { /* ... */
      // Determine namespace
      const namespaceToUse = currentEditSecretDetail.value?.metadata?.namespace || selectedNamespace.value;
      if (!namespaceToUse) { ElMessage.error(t('secretManagement.cannotDetermineNamespace')); return; }
  
      loading.dialogSave = true;
      // --- Replace with actual YAML editor interaction and API call ---
      // const currentYaml = yamlEditorRef.value.getContent();
      // try {
      //     let parsedYaml = yaml.load(currentYaml); ... validate ...
      //     const name = parsedYaml.metadata.name;
      //     const method = currentEditSecretDetail.value ? 'put' : 'post';
      //     const url = currentEditSecretDetail.value ? `/api/v1/namespaces/${namespaceToUse}/secrets/${name}` : `/api/v1/namespaces/${namespaceToUse}/secrets`;
      //     // Backend needs to handle the parsed JSON object (corev1.Secret)
      //     // Ensure frontend YAML editor content parses correctly to required Go struct
      //     const response = await request({ url, method, data: parsedYaml, baseURL:"..." });
      //     ... handle response ...
      // } catch (e) { ... }
  
       await new Promise(resolve => setTimeout(resolve, 500)); // Simulate
       loading.dialogSave = false; dialogVisible.value = false;
       const action = currentEditSecretDetail.value ? t('secretManagement.update') : t('secretManagement.create');
       ElMessage.success(t('secretManagement.operationSuccess', { action })); fetchSecretData();
  };
  
  const handleDeleteSecret = (secret: SecretDisplayItem) => { /* ... */
      ElMessageBox.confirm(
          t('secretManagement.deleteConfirmDetail', { name: secret.name, namespace: secret.namespace }),
          t('secretManagement.confirmDelete'), { type: 'warning' }
      ).then(async () => {
          loading.secrets = true;
          try {
              const response = await kubernetesRequest<{ code: number; message: string }>({
                  url: `/api/v1/namespaces/${secret.namespace}/secrets/${secret.name}`,
                  method: "delete"
              });
               if (response.code === 200 || response.code === 204 || response.code === 202) {
                  ElMessage.success(t('secretManagement.deleteSuccess', { name: secret.name })); await fetchSecretData();
              } else { ElMessage.error(`${t('secretManagement.deleteFailed')}: ${response.message || t('secretManagement.unknownError')}`); loading.secrets = false; }
          } catch (error: any) { console.error("删除 Secret 失败:", error); ElMessage.error(`${t('secretManagement.deleteFailed')}: ${error.response?.data?.message || error.message || t('secretManagement.requestFailed')}`); loading.secrets = false; }
      }).catch(() => ElMessage.info(t('secretManagement.deleteCanceled')));
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