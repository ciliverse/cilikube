<template>
    <div class="pv-page-container">
      <!-- Header: Title & Create Button -->
      <div class="page-header">
        <h1 class="page-title">{{ $t('persistentVolumeManagement.title') }}</h1>
         <el-button type="primary" :icon="PlusIcon" @click="handleAddPV" :loading="loading.page">
           {{ $t('persistentVolumeManagement.actions.create') }}
         </el-button>
      </div>

      <!-- Filter Bar: Search, Refresh -->
      <div class="filter-bar">
          <!-- No Namespace Select for PVs -->
          <el-input
              v-model="searchQuery"
              :placeholder="$t('persistentVolumeManagement.searchPlaceholder')"
              :prefix-icon="SearchIcon"
              clearable
              @input="handleSearchDebounced"
              class="filter-item search-input"
              :disabled="loading.pvs"
          />

          <el-tooltip :content="$t('actions.refresh')" placement="top">
              <el-button
                :icon="RefreshIcon"
                circle
                @click="fetchPVData"
                :loading="loading.pvs"
              />
          </el-tooltip>
      </div>

      <!-- PVs Table Card -->
      <el-card class="pv-table-card" shadow="hover">
        <el-table
            :data="paginatedData"
            v-loading="loading.pvs"
            border
            stripe
            style="width: 100%"
            @sort-change="handleSortChange"
            class="pv-table"
            :default-sort="{ prop: 'createdAt', order: 'descending' }"
            row-key="uid"
        >
          <el-table-column prop="name" :label="$t('persistentVolumeManagement.columns.name')" min-width="250" sortable="custom" fixed show-overflow-tooltip>
               <template #default="{ row }">
                  <el-icon class="pv-icon"><Coin /></el-icon>
                  <span class="pv-name">{{ row.name }}</span>
              </template>
          </el-table-column>
          <el-table-column prop="status" :label="$t('persistentVolumeManagement.columns.status')" min-width="120" sortable="custom" align="center">
              <template #default="{ row }">
                  <el-tag :type="getStatusTagType(row.status)" size="small" effect="light">
                       <el-icon class="status-icon" v-if="getStatusIcon(row.status)">
                           <component :is="getStatusIcon(row.status)" />
                       </el-icon>
                      {{ row.status }}
                  </el-tag>
              </template>
          </el-table-column>
          <el-table-column prop="capacity" :label="$t('persistentVolumeManagement.columns.capacity')" min-width="100" sortable="custom" align="right">
               <template #default="{ row }">
                  {{ formatCapacity(row.capacity) }}
              </template>
          </el-table-column>
          <el-table-column prop="accessModes" :label="$t('persistentVolumeManagement.columns.accessModes')" min-width="150">
               <template #default="{ row }">
                  <div v-for="mode in row.accessModes" :key="mode">
                      <el-tag size="small" type="info" effect="plain" style="margin-right: 4px;">{{ mode }}</el-tag>
                  </div>
                   <span v-if="!row.accessModes || row.accessModes.length === 0">N/A</span>
              </template>
          </el-table-column>
           <el-table-column prop="reclaimPolicy" :label="$t('persistentVolumeManagement.columns.reclaimPolicy')" min-width="120" align="center">
              <template #default="{ row }">
                  <el-tag :type="getReclaimPolicyTagType(row.reclaimPolicy)" size="small">{{ row.reclaimPolicy }}</el-tag>
              </template>
           </el-table-column>
            <el-table-column prop="volumeMode" :label="$t('persistentVolumeManagement.columns.volumeMode')" min-width="120" align="center">
               <template #default="{ row }">
                  <el-tag :type="getVolumeModeTagType(row.volumeMode)" size="small" effect="light">{{ row.volumeMode }}</el-tag>
              </template>
            </el-table-column>
           <el-table-column prop="storageClass" :label="$t('persistentVolumeManagement.columns.storageClass')" min-width="150" show-overflow-tooltip>
                <template #default="{ row }">
                   {{ row.storageClass || '-' }}
               </template>
           </el-table-column>
           <el-table-column prop="claim" :label="$t('persistentVolumeManagement.columns.boundClaim')" min-width="200" show-overflow-tooltip>
               <template #default="{ row }">
                   {{ row.claim || '-' }}
               </template>
           </el-table-column>
          <el-table-column prop="createdAt" :label="$t('persistentVolumeManagement.columns.createdAt')" min-width="180" sortable="custom" />
          <el-table-column :label="$t('persistentVolumeManagement.columns.actions')" width="130" align="center" fixed="right">
              <template #default="{ row }">
                   <el-tooltip :content="$t('persistentVolumeManagement.actions.editYaml')" placement="top">
                      <el-button link type="primary" :icon="EditPenIcon" @click="editPVYaml(row)" />
                  </el-tooltip>
                  <el-tooltip :content="$t('persistentVolumeManagement.actions.delete')" placement="top">
                      <el-button link type="danger" :icon="DeleteIcon" @click="handleDeletePV(row)" />
                  </el-tooltip>
              </template>
          </el-table-column>
           <template #empty>
            <el-empty :description="$t('persistentVolumeManagement.messages.noPVs')" :image-size="100" />
           </template>
      </el-table>

        <!-- Pagination -->
        <div class="pagination-container" v-if="!loading.pvs && totalPVs > 0">
            <el-pagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :page-sizes="[10, 20, 50, 100]"
                :total="totalPVs"
                layout="total, sizes, prev, pager, next, jumper"
                background
                @size-change="handleSizeChange"
                @current-change="handlePageChange"
                :disabled="loading.pvs"
            />
        </div>
      </el-card>
  
      <!-- Create/Edit Dialog (YAML focus) -->
      <el-dialog :title="dialogTitle" v-model="dialogVisible" width="70%" :close-on-click-modal="false">
         <el-alert type="info" :closable="false" style="margin-bottom: 20px;">
           {{ $t('persistentVolumeManagement.messages.yamlEditTip') }}
         </el-alert>
         <!-- Integrate a YAML editor component here -->
         <div class="yaml-editor-placeholder">
              {{ $t('persistentVolumeManagement.messages.yamlEditorPlaceholder') }}
              <pre>{{ yamlContent || placeholderYaml }}</pre>
         </div>
        <template #footer>
          <div class="dialog-footer">
              <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
              <el-button type="primary" @click="handleSaveYaml" :loading="loading.dialogSave">{{ $t('persistentVolumeManagement.actions.applyYaml') }}</el-button>
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
//   import yaml from 'js-yaml'; // Ensure installed
//   import Qty from 'js-quantities'; // For capacity formatting (npm install js-quantities)
  
  
  import {
      Plus as PlusIcon, Search as SearchIcon, Refresh as RefreshIcon, Coin, // Using Coin icon for PV
      EditPen as EditPenIcon, Delete as DeleteIcon,
      CircleCheckFilled, WarningFilled, CloseBold, Loading as LoadingIcon,
      QuestionFilled, InfoFilled, Link as LinkIcon // For Bound status
  } from '@element-plus/icons-vue'
  
  // --- Interfaces ---
  // Matching backend K8s native PersistentVolumeList response
  interface K8sPVListResponse {
      metadata: { resourceVersion: string };
      items: K8sPersistentVolume[];
  }

  interface K8sPersistentVolume {
      metadata: {
          name: string;
          uid: string;
          resourceVersion: string;
          creationTimestamp: string;
          annotations?: { [key: string]: string };
          finalizers?: string[];
          managedFields?: any[];
      };
      spec: {
          capacity: { storage: string };
          accessModes: string[];
          persistentVolumeReclaimPolicy: string;
          storageClassName: string;
          volumeMode: string;
          hostPath?: { path: string; type: string };
          claimRef?: {
              kind: string;
              namespace: string;
              name: string;
              uid: string;
              apiVersion: string;
              resourceVersion: string;
          };
      };
      status: {
          phase: string;
          lastPhaseTransitionTime?: string;
      };
  }

  interface PVApiResponse { 
      code: number; 
      data: K8sPVListResponse; 
      message: string;
  }
  
  // Internal Display/Table Item
  interface PVDisplayItem {
    uid: string
    name: string
    status: string
    capacity: string // Keep original string for sorting maybe? Or parse?
    capacityBytes: number // Parsed capacity for sorting
    accessModes: string[]
    reclaimPolicy: string
    storageClass: string
    volumeMode: string
    claim: string // namespace/pvcName
    createdAt: string
    // Raw API data for editing
    rawData?: K8sPersistentVolume // Store raw API response if needed for PUT/YAML edit
  }
  
  
  // --- Reactive State ---
  const allPVs = ref<PVDisplayItem[]>([])
  const currentPage = ref(1)
  const pageSize = ref(10)
  const totalPVs = ref(0) // Use API total
  const searchQuery = ref("")
  const sortParams = reactive({ key: 'createdAt', order: 'descending' as ('ascending' | 'descending' | null) })
  
  const loading = reactive({
      page: false, pvs: false, dialogSave: false
  })
  
  // Get i18n instance
  const { t } = useI18n()
  
  // Dialog state (YAML focus)
  const dialogVisible = ref(false)
  const dialogTitle = computed(() => currentEditPV.value ? `${t('persistentVolumeManagement.actions.edit')}: ${currentEditPV.value.metadata.name} (YAML)` : t('persistentVolumeManagement.messages.createPVTitle'))
  const currentEditPV = ref<K8sPersistentVolume | null>(null); // Store raw data for editing
  const yamlContent = ref("");
  const placeholderYaml = ref(`apiVersion: v1
  kind: PersistentVolume
  metadata:
    name: my-pv-name # Replace with unique name
  spec:
    capacity:
      storage: 5Gi # Example: 5 Gigabytes
    volumeMode: Filesystem # Or Block
    accessModes:
      - ReadWriteOnce # Options: ReadWriteOnce, ReadOnlyMany, ReadWriteMany, ReadWriteOncePod
    persistentVolumeReclaimPolicy: Retain # Options: Retain, Recycle, Delete
    storageClassName: standard # Match existing StorageClass or omit for manual binding
    # --- Configure specific volume source based on your storage ---
    # Example: HostPath (for testing/local dev ONLY)
    hostPath:
      path: "/mnt/data/my-pv-name" # Replace with actual path on node
    # Example: NFS
    # nfs:
    #   path: /path/on/nfs/server
    #   server: nfs-server.example.com
    # Example: CSI
    # csi:
    #   driver: ...
    #   volumeHandle: ...
    #   fsType: ext4
  `);
  
  // --- Computed Properties ---
  const filteredData = computed(() => {
      const query = searchQuery.value.trim().toLowerCase()
      if (!query) return allPVs.value;
      return allPVs.value.filter(pv =>
          pv.name.toLowerCase().includes(query) ||
          pv.storageClass.toLowerCase().includes(query)
          // Add more fields to search if needed
      );
  });
  
  const sortedData = computed(() => {
      const data = [...filteredData.value];
      const { key, order } = sortParams;
      if (!key || !order) return data;
  
      data.sort((a, b) => {
          let valA: any;
          let valB: any;
  
          if (key === 'capacity') {
               valA = a.capacityBytes ?? 0;
               valB = b.capacityBytes ?? 0;
          } else if (key === 'createdAt') {
              const timeA = a.createdAt ? dayjs(a.createdAt, "YYYY-MM-DD HH:mm:ss").valueOf() : 0;
              const timeB = b.createdAt ? dayjs(b.createdAt, "YYYY-MM-DD HH:mm:ss").valueOf() : 0;
              valA = isNaN(timeA) ? 0 : timeA;
              valB = isNaN(timeB) ? 0 : timeB;
          } else {
              valA = a[key as keyof PVDisplayItem] ?? '';
              valB = b[key as keyof PVDisplayItem] ?? '';
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
      if (!timestamp) return 'N/A'; return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss");
  }
  
  // Format capacity string (e.g., "5Gi", "100Mi") to human-readable format
  const formatCapacity = (capacity: string): string => {
      if (!capacity) return 'N/A';
      try {
          // 简单的格式化，保持原有格式但加空格
          return capacity.replace(/([0-9])([A-Za-z])/g, '$1 $2');
      } catch (e) {
          console.warn(`Could not parse capacity string "${capacity}":`, e);
          return capacity; // Return original string if parsing fails
      }
  }
  
  // Parse capacity string to bytes for sorting - simplified version
  const parseCapacityToBytes = (capacity: string): number => {
      if (!capacity) return 0;
      try {
          // 简单的解析，支持常见的单位
          const match = capacity.match(/^(\d+(?:\.\d+)?)([A-Za-z]*)$/);
          if (!match) return 0;
          
          const value = parseFloat(match[1]);
          const unit = match[2].toLowerCase();
          
          const multipliers: Record<string, number> = {
              'ki': 1024,
              'mi': 1024 * 1024,
              'gi': 1024 * 1024 * 1024,
              'ti': 1024 * 1024 * 1024 * 1024,
              'k': 1000,
              'm': 1000 * 1000,
              'g': 1000 * 1000 * 1000,
              't': 1000 * 1000 * 1000 * 1000,
              '': 1
          };
          
          return value * (multipliers[unit] || 1);
      } catch (e) {
          return 0; // Return 0 if parsing fails
      }
  }
  
  
  const getStatusTagType = (status: string): 'success' | 'warning' | 'danger' | 'info' => {
      const lowerStatus = status?.toLowerCase();
      if (lowerStatus === 'available') return 'success';
      if (lowerStatus === 'bound') return 'success'; // Bound is also a 'good' state
      if (lowerStatus === 'pending') return 'warning'; // Waiting for something
      if (lowerStatus === 'released') return 'info'; // Needs manual reclaim or deletion usually
      if (lowerStatus === 'failed') return 'danger';
      return 'info'; // Unknown
  }
  const getStatusIcon = (status: string) => {
      const lowerStatus = status?.toLowerCase();
      if (lowerStatus === 'available') return CircleCheckFilled;
      if (lowerStatus === 'bound') return LinkIcon; // Link icon for bound state
      if (lowerStatus === 'pending') return LoadingIcon;
      if (lowerStatus === 'released') return WarningFilled; // Released might need attention
      if (lowerStatus === 'failed') return CloseBold;
      return QuestionFilled;
   }
  
  const getReclaimPolicyTagType = (policy: string): 'success' | 'warning' | 'danger' | 'info' => {
      const lowerPolicy = policy?.toLowerCase();
      if (lowerPolicy === 'retain') return 'success'; // Safest option
      if (lowerPolicy === 'recycle') return 'warning'; // Deprecated and potentially unsafe
      if (lowerPolicy === 'delete') return 'danger'; // Data deleted when PVC is deleted
      return 'info';
  }
  
  const getVolumeModeTagType = (mode: string): 'info' | 'primary' => {
       return mode === 'Block' ? 'primary' : 'info';
  }
  
  // --- API Interaction ---
  const fetchPVData = async () => {
  loading.pvs = true;
  try {
    const params: Record<string, any> = {};
    const url = `/api/v1/persistentvolumes`;
    const response = await kubernetesRequest<PVApiResponse>({ 
      url, 
      method: "get", 
      params
    });

    if (response.code === 200 && response.data?.items) {
      // 适配K8s原生PersistentVolumeList格式
      const pvList = response.data as any; // K8s原生格式
      totalPVs.value = pvList.items?.length || 0; // 使用 items 数组长度作为 total
      allPVs.value = (pvList.items || []).map((item: any) => ({
        uid: item.metadata?.uid || '',
        name: item.metadata?.name || '',
        status: item.status?.phase || 'Unknown',
        capacity: item.spec?.capacity?.storage || 'N/A',
        capacityBytes: parseCapacityToBytes(item.spec?.capacity?.storage || ''),
        accessModes: item.spec?.accessModes || [],
        reclaimPolicy: item.spec?.persistentVolumeReclaimPolicy || 'N/A',
        storageClass: item.spec?.storageClassName || '',
        volumeMode: item.spec?.volumeMode || 'N/A',
        claim: item.spec?.claimRef ? `${item.spec.claimRef.namespace}/${item.spec.claimRef.name}` : '', // 格式化为 namespace/name
        createdAt: formatTimestamp(item.metadata?.creationTimestamp),
        rawData: item, // 存储原始数据
      }));
      const totalPages = Math.ceil(totalPVs.value / pageSize.value);
      if (currentPage.value > totalPages && totalPages > 0) currentPage.value = totalPages;
      else if (totalPVs.value === 0) currentPage.value = 1;
    } else {
      ElMessage.error(`获取 PV 数据失败: ${response.message || '未知错误'}`);
      allPVs.value = []; 
      totalPVs.value = 0;
    }
  } catch (error: any) {
    console.error("获取 PV 数据失败:", error);
    ElMessage.error(`获取 PV 数据出错: ${error.message || '网络请求失败'}`);
    allPVs.value = []; 
    totalPVs.value = 0;
  } finally {
    loading.pvs = false;
  }
};
  
  // --- Event Handlers ---
  const handlePageChange = (page: number) => { currentPage.value = page; /* Fetch only if server-side */ };
  const handleSizeChange = (size: number) => { pageSize.value = size; currentPage.value = 1; /* Fetch only if server-side */ };
  const handleSearchDebounced = debounce(() => { currentPage.value = 1; /* Fetch only if server-side */ }, 300);
  const handleSortChange = ({ prop, order }: { prop: string | null; order: 'ascending' | 'descending' | null }) => {
      sortParams.key = prop || 'createdAt';
      sortParams.order = order;
      currentPage.value = 1;
  };
  
  
  // --- Dialog and CRUD Actions ---
  const handleAddPV = () => {
      currentEditPV.value = null;
      yamlContent.value = placeholderYaml.value; // Use placeholder for new
      dialogVisible.value = true;
  };
  
  // Fetch full YAML for editing
  const editPVYaml = async (pv: PVDisplayItem) => {
       ElMessage.info(t('persistentVolumeManagement.messages.fetchingYaml').replace('{name}', pv.name));
       // --- Replace with actual API call to get YAML ---
       // try {
       //    loading.pvs = true;
       //    const response = await request<PVResponse>({ // Assuming backend returns our structured response on GET
       //       url: `/api/v1/persistentvolumes/${pv.name}`,
       //       method: 'get',
       //       baseURL: "VITE_API_BASE_URL"
       //    });
       //    if (response.code === 200 && response.data) {
       //        currentEditPV.value = response.data; // Store raw data
       //        // Fetch the *actual* K8s resource YAML if backend GET doesn't return it
       //        // Or convert the response object back to a simplified editable YAML
       //        yamlContent.value = yaml.dump(cleanK8sDataForEdit(response.data)); // Convert object to YAML string
       //        dialogTitle.value = `编辑 PV: ${pv.name} (YAML)`;
       //        dialogVisible.value = true;
       //    } else { // ... error handling ...}
       // } catch(e { // ... error handling ...}
       // finally { loading.pvs = false; }
  
       // Simulate fetching and opening editor
       currentEditPV.value = pv.rawData || null; // Use stored raw data if available
       yamlContent.value = JSON.stringify(pv.rawData, null, 2); // Use JSON instead of yaml for now
       dialogVisible.value = true;
  };
  
  const handleSaveYaml = async () => {
      loading.dialogSave = true;
      // --- Replace with actual YAML editor interaction and API call ---
      // const currentYaml = yamlEditorRef.value.getContent(); // Get from editor
      // try {
      //     let response;
      //     let parsedYaml = yaml.load(currentYaml);
      //     // Basic validation
      //     if (typeof parsedYaml !== 'object' || parsedYaml === null || !parsedYaml.metadata?.name) {
      //         throw new Error("无效的 YAML 或缺少 metadata.name");
      //     }
      //     const name = parsedYaml.metadata.name;
      //     const method = currentEditPV.value ? 'put' : 'post';
      //     const url = currentEditPV.value
      //                 ? `/api/v1/persistentvolumes/${name}`
      //                 : `/api/v1/persistentvolumes`;
      //
      //     // If PUT, ensure resourceVersion might be needed depending on backend/K8s API server behavior
      //     // if (method === 'put' && currentEditPV.value) {
      //     //     parsedYaml.metadata.resourceVersion = currentEditPV.value.resourceVersion;
      //     // }
      //
      //     response = await request({
      //          url: url,
      //          method: method,
      //          headers: { 'Content-Type': 'application/json' }, // Send JSON if backend expects JSON after parsing YAML
      //          // data: parsedYaml // Send parsed object
      //          // OR send raw YAML if backend handles it:
      //          // headers: { 'Content-Type': 'application/yaml' },
      //          // data: currentYaml
      //          data: parsedYaml, // Assuming backend handler binds JSON to corev1.PersistentVolume
      //          baseURL: "VITE_API_BASE_URL",
      //     });
      //
      //     if (response.code === 200 || response.code === 201) {
      //        ElMessage.success(`PV ${currentEditPV.value ? '更新' : '创建'}成功`);
      //        dialogVisible.value = false; fetchPVData();
      //     } else { // ... error handling ... }
      // } catch (error: any) { // ... error handling ... }
      // finally { loading.dialogSave = false; }
  
       // Simulate success
       await new Promise(resolve => setTimeout(resolve, 500));
       loading.dialogSave = false;
       dialogVisible.value = false;
       const action = currentEditPV.value ? '更新' : '创建';
       ElMessage.success(`模拟 PV ${action}成功`);
       fetchPVData();
  };
  
  
  const handleDeletePV = (pv: PVDisplayItem) => {
       ElMessageBox.confirm(
          `确定要删除 PersistentVolume "${pv.name}" 吗？请确保它没有被 PVC 绑定，否则可能需要先解除绑定。此操作可能不可恢复。`,
          '确认删除', { type: 'warning' }
      ).then(async () => {
          loading.pvs = true;
          try {
              const response = await kubernetesRequest<{ code: number; message: string }>({
                  url: `/api/v1/persistentvolumes/${pv.name}`,
                  method: "delete"
              });
               if (response.code === 200 || response.code === 204 || response.code === 202) { // Check for success codes
                  ElMessage.success(`PV "${pv.name}" 已删除`);
                  await fetchPVData(); // Refresh list fully
              } else {
                   ElMessage.error(`删除 PV 失败: ${response.message || '未知错误'}`);
                   loading.pvs = false;
              }
          } catch (error: any) { console.error("删除 PV 失败:", error); ElMessage.error(`删除 PV 失败: ${error.response?.data?.message || error.message || '请求失败'}`); loading.pvs = false; }
      }).catch(() => ElMessage.info('删除操作已取消'));
  };
  
  
  // --- Lifecycle Hooks ---
  onMounted(async () => {
      loading.page = true;
      // No namespaces to fetch for PVs
      await fetchPVData();
      loading.page = false;
  });
  
  </script>
  
<style lang="scss" scoped>
// Define fallbacks or import variables
$page-padding: 24px; $spacing-md: 16px; $spacing-lg: 24px;
$font-size-base: 14px; $font-size-small: 12px; $font-size-large: 16px; $font-size-extra-large: 28px;
$border-radius-base: 4px; $kube-pv-icon-color: #f59e0b; // Example Amber color

.pv-page-container { 
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
.search-input { width: 300px; }

.pv-table-card {
  margin-bottom: $spacing-lg;
  border-radius: $border-radius-base;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
  border: 1px solid var(--el-border-color-lighter);
  
  :deep(.el-card__body) {
    padding: 0;
  }
}
  
.pv-table {
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
     .pv-icon { 
       margin-right: 8px; 
       color: $kube-pv-icon-color; 
       vertical-align: middle; 
       font-size: 18px; 
       position: relative; 
       top: -1px; 
     }
     .pv-name { 
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
  border-color: var(--el-color-primary);
  color: var(--el-color-primary);
  transition: all 0.2s;
}

:deep(.el-button--primary:hover) {
  background-color: var(--el-color-primary);
  border-color: var(--el-color-primary);
  color: #ffffff;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(59, 130, 246, 0.2);
}

:deep(.el-button--primary:active),
:deep(.el-button--primary.is-active) {
  background-color: var(--el-color-primary);
  border-color: var(--el-color-primary);
  color: #ffffff;
  transform: translateY(0);
}

:deep(.el-button--success) {
  background-color: transparent;
  border-color: var(--el-color-success);
  color: var(--el-color-success);
  transition: all 0.2s;
}

:deep(.el-button--success:hover) {
  background-color: var(--el-color-success);
  border-color: var(--el-color-success);
  color: #ffffff;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(16, 185, 129, 0.2);
}

:deep(.el-button--success:active),
:deep(.el-button--success.is-active) {
  background-color: var(--el-color-success);
  border-color: var(--el-color-success);
  color: #ffffff;
  transform: translateY(0);
}

:deep(.el-button--danger) {
  background-color: transparent;
  border-color: var(--el-color-danger);
  color: var(--el-color-danger);
  transition: all 0.2s;
}

:deep(.el-button--danger:hover) {
  background-color: var(--el-color-danger);
  border-color: var(--el-color-danger);
  color: #ffffff;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(239, 68, 68, 0.2);
}

:deep(.el-button--danger:active),
:deep(.el-button--danger.is-active) {
  background-color: var(--el-color-danger);
  border-color: var(--el-color-danger);
  color: #ffffff;
  transform: translateY(0);
}
  </style>