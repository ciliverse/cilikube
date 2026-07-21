<template>
    <div class="pvc-page-container">
      <!-- Header: Title & Create Button -->
      <div class="page-header">
        <h1 class="page-title">{{ $t('persistentVolumeClaimManagement.title') }}</h1>
         <el-button
           type="primary"
           :icon="PlusIcon"
           @click="handleAddPVC"
           :loading="loading.page"
           :disabled="!selectedNamespace"
         >
           {{ $t('persistentVolumeClaimManagement.actions.create') }}
         </el-button>
      </div>
  
      <!-- Filter Bar: Namespace, Search, Refresh -->
      <div class="filter-bar">
          <el-select
              v-model="selectedNamespace"
              :placeholder="$t('persistentVolumeClaimManagement.filterNamespace')"
              @change="handleNamespaceChange"
              filterable
              :loading="loading.namespaces"
              class="filter-item namespace-select"
              style="width: 280px;"
          >
              <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
               <template #empty>
                  <div style="padding: 10px; text-align: center; color: #999;">
                      {{ loading.namespaces ? $t('common.loading') : $t('persistentVolumeClaimManagement.messages.noNamespaces') }}
                  </div>
              </template>
          </el-select>
  
          <el-input
              v-model="searchQuery"
              :placeholder="$t('persistentVolumeClaimManagement.searchPlaceholder')"
              :prefix-icon="SearchIcon"
              clearable
              @input="handleSearchDebounced"
              class="filter-item search-input"
              :disabled="!selectedNamespace || loading.pvcs"
          />
  
          <el-tooltip :content="$t('actions.refresh')" placement="top">
              <el-button
                :icon="RefreshIcon"
                circle
                @click="fetchPvcData"
                :loading="loading.pvcs"
                :disabled="!selectedNamespace"
              />
          </el-tooltip>
      </div>
  
      <!-- PVCs Table Card -->
      <el-card class="pvc-table-card" shadow="hover">
        <el-table
            :data="paginatedData"
            v-loading="loading.pvcs"
            border
            stripe
            style="width: 100%"
            @sort-change="handleSortChange"
            class="pvc-table"
            :default-sort="{ prop: 'createdAt', order: 'descending' }"
            row-key="uid"
        >
          <el-table-column prop="name" :label="$t('persistentVolumeClaimManagement.columns.name')" min-width="250" sortable="custom" fixed show-overflow-tooltip>
               <template #default="{ row }">
                  <el-icon class="pvc-icon"><MessageBox /></el-icon>
                  <span class="pvc-name">{{ row.name }}</span>
              </template>
          </el-table-column>
          <el-table-column prop="namespace" :label="$t('persistentVolumeClaimManagement.columns.namespace')" min-width="150" sortable="custom" show-overflow-tooltip />
          <el-table-column prop="status" :label="$t('persistentVolumeClaimManagement.columns.status')" min-width="120" sortable="custom" align="center">
              <template #default="{ row }">
                  <el-tag :type="getStatusTagType(row.status)" size="small" effect="light">
                      <el-icon class="status-icon" v-if="getStatusIcon(row.status)" :class="getSpinClass(row.status)">
                           <component :is="getStatusIcon(row.status)" />
                       </el-icon>
                      {{ row.status }}
                  </el-tag>
              </template>
          </el-table-column>
          <el-table-column prop="volumeName" :label="$t('persistentVolumeClaimManagement.columns.boundVolume')" min-width="200" sortable="custom" show-overflow-tooltip>
              <template #default="{ row }">
                  {{ row.volumeName || '-' }}
              </template>
          </el-table-column>
          <el-table-column prop="capacity" :label="$t('persistentVolumeClaimManagement.columns.requestCapacity')" min-width="110" sortable="custom" align="right">
               <template #default="{ row }">
                  {{ formatCapacity(row.capacity) }}
              </template>
          </el-table-column>
           <el-table-column prop="actualCapacity" :label="$t('persistentVolumeClaimManagement.columns.actualCapacity')" min-width="110" sortable="custom" align="right">
               <template #default="{ row }">
                  {{ formatCapacity(row.actualCapacity) }}
              </template>
          </el-table-column>
          <el-table-column prop="accessModes" :label="$t('persistentVolumeClaimManagement.columns.accessModes')" min-width="150">
               <template #default="{ row }">
                  <div v-for="mode in row.accessModes" :key="mode">
                      <el-tag size="small" type="info" effect="plain" style="margin-right: 4px;">{{ mode }}</el-tag>
                  </div>
                   <span v-if="!row.accessModes || row.accessModes.length === 0">N/A</span>
              </template>
          </el-table-column>
           <el-table-column prop="volumeMode" :label="$t('persistentVolumeClaimManagement.columns.volumeMode')" min-width="100" align="center">
               <template #default="{ row }">
                  <el-tag type="info" size="small" effect="light">{{ row.volumeMode }}</el-tag>
              </template>
            </el-table-column>
           <el-table-column prop="storageClass" :label="$t('persistentVolumeClaimManagement.columns.storageClass')" min-width="150" show-overflow-tooltip>
               <template #default="{ row }">
                   {{ row.storageClass || '<none>' }}
               </template>
           </el-table-column>
          <el-table-column prop="createdAt" :label="$t('persistentVolumeClaimManagement.columns.createdAt')" min-width="180" sortable="custom" />
          <el-table-column :label="$t('persistentVolumeClaimManagement.columns.actions')" width="130" align="center" fixed="right">
              <template #default="{ row }">
                   <el-tooltip :content="$t('persistentVolumeClaimManagement.actions.editYaml')" placement="top">
                      <el-button link type="primary" :icon="EditPenIcon" @click="editPvcYaml(row)" />
                  </el-tooltip>
                  <el-tooltip :content="$t('persistentVolumeClaimManagement.actions.delete')" placement="top">
                      <el-button link type="danger" :icon="DeleteIcon" @click="handleDeletePVC(row)" />
                  </el-tooltip>
              </template>
          </el-table-column>
           <template #empty>
            <el-empty v-if="!loading.pvcs && !selectedNamespace" :description="$t('persistentVolumeClaimManagement.messages.selectNamespace')" image-size="100" />
            <el-empty v-else-if="!loading.pvcs && paginatedData.length === 0" :description="$t('persistentVolumeClaimManagement.messages.noPVCs')" image-size="100" />
           </template>
      </el-table>

        <!-- Pagination -->
        <div class="pagination-container" v-if="!loading.pvcs && totalPvcs > 0">
            <el-pagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :page-sizes="[10, 20, 50, 100]"
                :total="totalPvcs"
                layout="total, sizes, prev, pager, next, jumper"
                background
                @size-change="handleSizeChange"
                @current-change="handlePageChange"
                :disabled="loading.pvcs"
            />
        </div>
      </el-card>
  
      <!-- Create/Edit Dialog (YAML focus) -->
      <el-dialog :title="dialogTitle" v-model="dialogVisible" width="70%" :close-on-click-modal="false">
         <el-alert type="info" :closable="false" style="margin-bottom: 20px;">
           {{ $t('persistentVolumeClaimManagement.messages.yamlEditTip') }}
         </el-alert>
         <!-- Integrate a YAML editor component here -->
         <div class="yaml-editor-placeholder">
              {{ $t('persistentVolumeClaimManagement.messages.yamlEditorPlaceholder') }}
              <pre>{{ yamlContent || placeholderYaml }}</pre>
         </div>
        <template #footer>
          <div class="dialog-footer">
              <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
              <el-button type="primary" @click="handleSaveYaml" :loading="loading.dialogSave">{{ $t('persistentVolumeClaimManagement.actions.applyYaml') }}</el-button>
          </div>
        </template>
      </el-dialog>
  
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, reactive, computed, onMounted } from "vue"
  import { useI18n } from 'vue-i18n'
  import { ElMessage, ElMessageBox } from "element-plus"
  import { kubernetesRequest, fetchNamespaces, KubernetesApiResponse } from "@/utils/api-config" // Ensure correct path
  import dayjs from "dayjs"
  import { debounce } from 'lodash-es'
  import yaml from 'js-yaml'; // Ensure installed
  import Qty from 'js-quantities'; // Ensure installed
  
  import {
      Plus as PlusIcon, Search as SearchIcon, Refresh as RefreshIcon, MessageBox, // Icon for PVC
      EditPen as EditPenIcon, Delete as DeleteIcon,
      CircleCheckFilled, WarningFilled, CloseBold, Loading as LoadingIcon,
      QuestionFilled, InfoFilled, Link as LinkIcon
  } from '@element-plus/icons-vue'
  
  // --- Interfaces ---
  // ** Interface matching the ACTUAL API response item structure **
  interface PVCApiItem {
    name: string;
    namespace: string;
    uid: string;
    status: string; // Direct status string
    volumeName: string; // Direct volume name
    storageClass: string | null; // Can be null
    accessModes: string[];
    requestedStorage: string; // Direct requested storage string
    actualCapacity?: string | null; // Direct actual capacity string (optional)
    volumeMode: string;
    createdAt: string;
    annotations?: { [key: string]: string };
    labels?: { [key: string]: string }; // Added labels based on sample
    resourceVersion: string;
  }
  
  interface PVCListApiResponseData { items: PVCApiItem[]; total: number } // Assuming total is always present now
  interface PVCApiResponse { code: number; data: PVCListApiResponseData; message: string }
  interface NamespaceListResponse { code: number; data: { items: Array<{ metadata: { name: string } }> }; message: string }
  
  // Internal Display/Table Item - adjusted to use direct fields from PVCApiItem
  interface PVCDisplayItem {
    uid: string
    name: string
    namespace: string
    status: string
    volumeName: string
    capacity: string      // Renamed to capacity, stores requestedStorage
    actualCapacity: string // Stores actualCapacity
    capacityBytes: number // Parsed requested capacity
    accessModes: string[]
    storageClass: string
    volumeMode: string
    createdAt: string
    rawData: PVCApiItem // Store raw API item
  }
  
  // --- Reactive State ---
  const allPvcs = ref<PVCDisplayItem[]>([])
  const namespaces = ref<string[]>([])
  const selectedNamespace = ref<string>("")
  const currentPage = ref(1)
  const pageSize = ref(10)
  const totalPvcs = ref(0)
  const searchQuery = ref("")
  const sortParams = reactive({ key: 'createdAt', order: 'descending' as ('ascending' | 'descending' | null) })
  
  const loading = reactive({
      page: false, namespaces: false, pvcs: false, dialogSave: false
  })
  
  // Get i18n instance
  const { t } = useI18n()
  
  // Dialog state (YAML focus)
  const dialogVisible = ref(false)
  const dialogTitle = computed(() => currentEditPvc.value ? `${t('persistentVolumeClaimManagement.actions.edit')}: ${currentEditPvc.value.metadata.name} (YAML)` : t('persistentVolumeClaimManagement.messages.createPVCTitle'))
  const currentEditPvc = ref<PVCApiItem | null>(null); // Store raw API item
  const yamlContent = ref("");
  const placeholderYaml = computed(() => `apiVersion: v1
  kind: PersistentVolumeClaim
  metadata:
    name: my-new-pvc
    namespace: ${selectedNamespace.value || 'default'}
  spec:
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
    storageClassName: standard
  `);
  
  // --- Computed Properties ---
  const filteredData = computed(() => { /* ... as before ... */
      const query = searchQuery.value.trim().toLowerCase()
      if (!query) return allPvcs.value;
      return allPvcs.value.filter(pvc =>
          pvc.name.toLowerCase().includes(query) ||
          (pvc.volumeName && pvc.volumeName.toLowerCase().includes(query))
      );
  });
  
  const sortedData = computed(() => {
      const data = [...filteredData.value];
      const { key, order } = sortParams;
      if (!key || !order) return data;
  
      data.sort((a, b) => {
          let valA: any; let valB: any;
          // ** Sort by 'capacity' uses the parsed 'capacityBytes' **
          if (key === 'capacity' || key === 'requestedStorage') {
               valA = a.capacityBytes ?? 0;
               valB = b.capacityBytes ?? 0;
          } else if (key === 'actualCapacity') {
               // Parse actual capacity for sorting if needed, otherwise sort by string
               valA = parseCapacityToBytes(a.actualCapacity) ?? 0;
               valB = parseCapacityToBytes(b.actualCapacity) ?? 0;
          }
          else if (key === 'createdAt') {
              const timeA = a.createdAt ? dayjs(a.createdAt, "YYYY-MM-DD HH:mm:ss").valueOf() : 0;
              const timeB = b.createdAt ? dayjs(b.createdAt, "YYYY-MM-DD HH:mm:ss").valueOf() : 0;
              valA = isNaN(timeA) ? 0 : timeA; valB = isNaN(timeB) ? 0 : timeB;
          } else {
              // Ensure type safety for other keys
              valA = a[key as keyof Omit<PVCDisplayItem, 'capacityBytes' | 'rawData'>] ?? '';
              valB = b[key as keyof Omit<PVCDisplayItem, 'capacityBytes' | 'rawData'>] ?? '';
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
  
  
  const paginatedData = computed(() => { /* ... as before ... */
      const start = (currentPage.value - 1) * pageSize.value;
      const end = start + pageSize.value;
      return sortedData.value.slice(start, end);
  });
  
  
  // --- Helper Functions ---
  const formatTimestamp = (timestamp: string): string => { /* ... */ if (!timestamp) return 'N/A'; return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss"); }
  const formatCapacity = (capacity: string | undefined | null): string => { /* ... Added null check */ if (!capacity) return '-'; try { const qty = Qty(capacity); return qty.format('gib').replace(/([KMG])iB/, ' $1iB'); } catch (e) { console.warn(`Failed to format capacity: ${capacity}`); return capacity; } }
  const parseCapacityToBytes = (capacity: string | undefined | null): number => { /* ... Added null check */ if (!capacity) return 0; try { const qty = Qty(capacity); return qty.toBase().scalar; } catch (e) { return 0; } }
  const getStatusTagType = (status: string): 'success' | 'warning' | 'danger' | 'info' => { /* ... */ const lowerStatus = status?.toLowerCase(); if (lowerStatus === 'bound') return 'success'; if (lowerStatus === 'pending') return 'warning'; if (lowerStatus === 'lost') return 'danger'; return 'info'; }
  const getStatusIcon = (status: string) => { /* ... */ const lowerStatus = status?.toLowerCase(); if (lowerStatus === 'bound') return LinkIcon; if (lowerStatus === 'pending') return LoadingIcon; if (lowerStatus === 'lost') return CloseBold; return QuestionFilled; }
  const getSpinClass = (status: string) => { /* ... */ return status?.toLowerCase() === 'pending' ? 'is-loading' : ''; }
  
  // --- API Interaction ---
  const fetchNamespaces = async () => { /* ... same as before ... */
      loading.namespaces = true;
      try {
          const response = await kubernetesRequest<NamespaceListResponse>({ url: "/api/v1/namespaces", method: "get" });
          if (response.code === 200 && response.data && Array.isArray(response.data.items)) {
              namespaces.value = response.data.items.map(ns => ns.metadata.name);
              if (namespaces.value.length > 0 && !selectedNamespace.value) {
                   selectedNamespace.value = namespaces.value.find(ns => ns === 'default') || namespaces.value[0];
              } else if (namespaces.value.length === 0) { ElMessage.warning(t('persistentVolumeClaimManagement.messages.noNamespaces')); }
          } else { ElMessage.error(t('persistentVolumeClaimManagement.messages.fetchNamespacesError')); namespaces.value = []; }
      } catch (error: any) { console.error("获取命名空间失败:", error); ElMessage.error(t('persistentVolumeClaimManagement.messages.fetchNamespacesError') + `: ${error.message || t('common.error')}`); namespaces.value = []; }
      finally { loading.namespaces = false; }
  }
  
  const fetchPvcData = async () => {
      if (!selectedNamespace.value) { allPvcs.value = []; totalPvcs.value = 0; return; }
      loading.pvcs = true;
      allPvcs.value = [];
      totalPvcs.value = 0;
      try {
          const params: Record<string, any> = { /* Server-side params if needed */ };
          // ** Use correct endpoint name from backend routing **
          const url = `/api/v1/namespaces/${selectedNamespace.value}/persistentvolumeclaims`; // Make sure this matches Go route
          const response = await kubernetesRequest<PVCApiResponse>({ url, method: "get", params });
  
          if (response.code === 200 && response.data?.items && Array.isArray(response.data.items)) {
              totalPvcs.value = response.data.total; // Use total from API
  
              allPvcs.value = response.data.items
                  .filter(item => item && item.name && item.namespace) // Basic check
                  .map((item, index) => {
                      // ** Direct mapping based on the NEW API response structure **
                      const requestedCapacityStr = item.requestedStorage || 'N/A';
                      const actualCapacityStr = item.actualCapacity || '';
  
                      return {
                          // Use item.uid directly if provided by backend, otherwise fallback
                          uid: item.uid || `${item.namespace}-${item.name}-${index}`,
                          name: item.name,
                          namespace: item.namespace,
                          status: item.status || 'Unknown',
                          volumeName: item.volumeName || '',
                          capacity: requestedCapacityStr, // Store requested capacity here
                          actualCapacity: actualCapacityStr,
                          capacityBytes: parseCapacityToBytes(item.requestedStorage), // Parse requested for sorting
                          accessModes: item.accessModes || [],
                          storageClass: item.storageClass || '',
                          volumeMode: item.volumeMode || 'Filesystem', // Default if missing
                          createdAt: formatTimestamp(item.createdAt),
                          rawData: item, // Store the raw item from API
                      };
              });
  
              const totalPages = Math.ceil(totalPvcs.value / pageSize.value);
               if (currentPage.value > totalPages && totalPages > 0) currentPage.value = totalPages;
               else if (totalPvcs.value === 0) currentPage.value = 1;
  
          } else if (response.code === 200 && response.data?.items === null) {
              console.log(`No PVCs found in namespace '${selectedNamespace.value}' (items is null).`);
              allPvcs.value = []; totalPvcs.value = 0; currentPage.value = 1;
          } else {
              ElMessage.error(t('persistentVolumeClaimManagement.messages.fetchPVCsError'));
              allPvcs.value = []; totalPvcs.value = 0;
          }
      } catch (error: any) {
          console.error("获取 PVC 数据失败:", error);
          ElMessage.error(t('persistentVolumeClaimManagement.messages.fetchPVCsError'));
          allPvcs.value = []; totalPvcs.value = 0;
      } finally {
          loading.pvcs = false;
      }
  }
  
  
  // --- Event Handlers ---
  const handleNamespaceChange = () => { /* ... */ currentPage.value = 1; searchQuery.value = ''; sortParams.key = 'createdAt'; sortParams.order = 'descending'; fetchPvcData(); };
  const handlePageChange = (page: number) => { currentPage.value = page; /* Fetch if server-side */ };
  const handleSizeChange = (size: number) => { pageSize.value = size; currentPage.value = 1; /* Fetch if server-side */ };
  const handleSearchDebounced = debounce(() => { currentPage.value = 1; /* Fetch if server-side */ }, 300);
  const handleSortChange = ({ prop, order }: { prop: string | null; order: 'ascending' | 'descending' | null }) => {
      // ** Remap sort key if needed to match internal display item property **
      let sortKey = prop;
      if (prop === 'requestedStorage') {
          sortKey = 'capacity'; // Sort by the parsed capacityBytes via 'capacity' key in sortedData computed
      }
       // Add similar mapping for other columns if prop name differs from display item key
      sortParams.key = sortKey || 'createdAt';
      sortParams.order = order;
      currentPage.value = 1;
  };
  
  
  // --- Dialog and CRUD Actions ---
  const handleAddPVC = () => { /* ... */ if (!selectedNamespace.value) { ElMessage.warning(t('persistentVolumeClaimManagement.messages.selectNamespaceWarning')); return; } currentEditPvc.value = null; yamlContent.value = placeholderYaml.value; dialogVisible.value = true; };
  const editPvcYaml = async (pvc: PVCDisplayItem) => { /* ... */
      ElMessage.info(t('persistentVolumeClaimManagement.messages.fetchingYaml').replace('{name}', pvc.name));
      if (!pvc.rawData) { ElMessage.error(t('persistentVolumeClaimManagement.messages.editError')); return; }
      currentEditPvc.value = pvc.rawData;
      // Convert raw API structure *back* to standard K8s PVC for YAML editing
      // This is needed if your API mapping function (ToPVCResponse) significantly changed the structure
      const k8sObjectForYaml = {
          apiVersion: "v1",
          kind: "PersistentVolumeClaim",
          metadata: {
              name: pvc.rawData.name,
              namespace: pvc.rawData.namespace,
              labels: pvc.rawData.labels,
              annotations: pvc.rawData.annotations,
              resourceVersion: pvc.rawData.resourceVersion, // Keep RV for PUT
              // uid and creationTimestamp should usually be omitted for edits
          },
          spec: {
              accessModes: pvc.rawData.accessModes,
              resources: { requests: { storage: pvc.rawData.requestedStorage } },
              volumeName: pvc.rawData.volumeName || undefined, // Only include if present
              storageClassName: pvc.rawData.storageClass,
              volumeMode: pvc.rawData.volumeMode
          },
          // status is omitted for editing
      };
      yamlContent.value = yaml.dump(k8sObjectForYaml, { skipInvalid: true }); // skipInvalid helps with undefined fields
      dialogVisible.value = true;
  };
  
  const handleSaveYaml = async () => { /* ... */
      if (!selectedNamespace.value && !currentEditPvc.value?.namespace) { ElMessage.error("无法确定目标命名空间。"); return; }
      loading.dialogSave = true;
      // --- Replace with actual YAML editor interaction and API call ---
      // const currentYaml = yamlEditorRef.value.getContent();
      // try {
      //     let parsedYaml = yaml.load(currentYaml); ... validate ...
      //     const name = parsedYaml.metadata.name;
      //     const namespaceToUse = parsedYaml.metadata.namespace || currentEditPvc.value?.namespace || selectedNamespace.value;
      //     const method = currentEditPvc.value ? 'put' : 'post';
      //     const url = currentEditPvc.value ? `/api/v1/namespaces/${namespaceToUse}/persistentvolumeclaims/${name}` : `/api/v1/namespaces/${namespaceToUse}/persistentvolumeclaims`;
      //     // Backend expects the standard K8s structure, so send parsedYaml
      //     const response = await request({ url, method, data: parsedYaml, baseURL:"..." });
      //     ... handle response ...
      // } catch (e) { ... }
  
       await new Promise(resolve => setTimeout(resolve, 500)); // Simulate
       loading.dialogSave = false; dialogVisible.value = false;
       const action = currentEditPvc.value ? t('persistentVolumeClaimManagement.actions.edit') : t('persistentVolumeClaimManagement.actions.create');
       ElMessage.success(t('persistentVolumeClaimManagement.messages.operationSuccess').replace('{action}', action)); fetchPvcData();
  };
  
  const handleDeletePVC = (pvc: PVCDisplayItem) => { /* ... */
      ElMessageBox.confirm(
          t('persistentVolumeClaimManagement.messages.deleteConfirmDetail').replace('{name}', pvc.name).replace('{namespace}', pvc.namespace),
          t('common.confirm'), { type: 'warning' }
      ).then(async () => {
          loading.pvcs = true;
          try {
              const response = await kubernetesRequest<{ code: number; message: string }>({
                  url: `/api/v1/namespaces/${pvc.namespace}/persistentvolumeclaims/${pvc.name}`,
                  method: "delete"
              });
               if (response.code === 200 || response.code === 204 || response.code === 202) {
                  ElMessage.success(t('persistentVolumeClaimManagement.messages.deleteSuccess')); await fetchPvcData();
              } else { ElMessage.error(t('persistentVolumeClaimManagement.messages.deleteFailed')); loading.pvcs = false; }
          } catch (error: any) { console.error("删除 PVC 失败:", error); ElMessage.error(t('persistentVolumeClaimManagement.messages.deleteFailed')); loading.pvcs = false; }
      }).catch(() => ElMessage.info(t('persistentVolumeClaimManagement.messages.operationCancelled')));
  };
  
  // --- Lifecycle Hooks ---
  onMounted(async () => { /* ... */
      loading.page = true;
      await fetchNamespaces();
      if (selectedNamespace.value) { await fetchPvcData(); }
      loading.page = false;
  });
  
  </script>
  
  <style lang="scss" scoped>
  /* Using fallback variables directly */
  $page-padding: 24px; $spacing-md: 16px; $spacing-lg: 24px;
  $font-size-base: 14px; $font-size-small: 12px; $font-size-large: 16px; $font-size-extra-large: 28px;
  $border-radius-base: 4px; $kube-pvc-icon-color: #9b59b6;
  
  .pvc-page-container { 
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
  .namespace-select { width: 240px; }
  .search-input { width: 300px; }
  
  .pvc-table-card {
    margin-bottom: $spacing-lg;
    border-radius: $border-radius-base;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
    border: 1px solid var(--el-border-color-lighter);
    
    :deep(.el-card__body) {
      padding: 0;
    }
  }
  
  .pvc-table {
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
     .pvc-icon { 
       margin-right: 8px; 
       color: $kube-pvc-icon-color; 
       vertical-align: middle; 
       font-size: 18px; 
       position: relative; 
       top: -1px; 
     }
     .pvc-name { 
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

  /* 按钮样式 - 透明背景 + 彩色边框 */
  :deep(.el-button--primary) {
    background-color: transparent;
    border-color: #3b82f6;
    color: #3b82f6;
  }

  :deep(.el-button--primary:hover) {
    background-color: #3b82f6;
    border-color: #3b82f6;
    color: #ffffff;
  }

  :deep(.el-button--primary:active),
  :deep(.el-button--primary.is-active) {
    background-color: #3b82f6;
    border-color: #3b82f6;
    color: #ffffff;
  }

  :deep(.el-button--success) {
    background-color: transparent;
    border-color: #10b981;
    color: #10b981;
  }

  :deep(.el-button--success:hover) {
    background-color: #10b981;
    border-color: #10b981;
    color: #ffffff;
  }

  :deep(.el-button--success:active),
  :deep(.el-button--success.is-active) {
    background-color: #10b981;
    border-color: #10b981;
    color: #ffffff;
  }

  :deep(.el-button--danger) {
    background-color: transparent;
    border-color: #ef4444;
    color: #ef4444;
  }

  :deep(.el-button--danger:hover) {
    background-color: #ef4444;
    border-color: #ef4444;
    color: #ffffff;
  }

  :deep(.el-button--danger:active),
  :deep(.el-button--danger.is-active) {
    background-color: #ef4444;
    border-color: #ef4444;
    color: #ffffff;
  }
  </style>