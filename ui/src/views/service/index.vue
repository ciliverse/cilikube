<template>
  <div class="service-page-container">
    <!-- Breadcrumbs -->
    <el-breadcrumb separator="/" class="page-breadcrumb">
      <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
      <el-breadcrumb-item>网络管理</el-breadcrumb-item>
      <el-breadcrumb-item>Services (SVC)</el-breadcrumb-item>
    </el-breadcrumb>

    <!-- Header: Title & Create Button -->
    <div class="page-header">
      <h1 class="page-title">Services (SVC) 列表</h1>
       <el-button
         type="primary"
         :icon="PlusIcon"
         @click="handleAddService"
         :loading="loading.page"
         :disabled="!selectedNamespace"
       >
         创建 Service (YAML)
       </el-button>
    </div>

     <!-- Cluster Knowledge Alert -->
     <el-alert
       title="关于 Services"
       type="info"
       show-icon
       :closable="true"
       class="info-alert"
       description="在 Kubernetes 中，Service 是一种抽象，它定义了 Pod 的逻辑集和访问它们的策略。Service 使依赖 Pod 集合的应用解耦。Service 通过标签选择器（Label Selector）匹配一组 Pod。例如，一个运行了 3 个副本的后端应用，这些副本是可互换的，前端不需要关心它们调用了哪个后端副本。 然而组成这一后端集合的 Pod 实际上可能会发生变动， 前端客户端不应该也没必要知道，而且也不需要跟踪这一后端集合的状态。 Service 提供的这种抽象能够解耦这种关联。"
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
            placeholder="搜索 Service 名称 / 类型 / ClusterIP"
            :prefix-icon="SearchIcon"
            clearable
            @input="handleSearchDebounced"
            class="filter-item search-input"
            :disabled="!selectedNamespace || loading.services"
        />

        <el-tooltip content="刷新列表" placement="top">
            <el-button
              :icon="RefreshIcon"
              circle
              @click="fetchServiceData"
              :loading="loading.services"
              :disabled="!selectedNamespace"
            />
        </el-tooltip>
    </div>

    <!-- Services Table -->
    <el-table
        :data="paginatedData"
        v-loading="loading.services"
        border
        stripe
        style="width: 100%"
        @sort-change="handleSortChange"
        class="service-table"
        :default-sort="{ prop: 'createdAt', order: 'descending' }"
        row-key="uid"
    >
        <el-table-column prop="name" label="名称" min-width="250" sortable="custom" fixed show-overflow-tooltip>
             <template #default="{ row }">
                <el-icon class="service-icon"><ServiceIcon /></el-icon>
                <span class="service-name">{{ row.name }}</span>
            </template>
        </el-table-column>
        <el-table-column prop="namespace" label="命名空间" min-width="150" sortable="custom" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" min-width="130" sortable="custom" align="center">
             <template #default="{ row }">
                <el-tag size="small" effect="light" :type="getTypeTagType(row.type)">{{ row.type }}</el-tag>
            </template>
        </el-table-column>
         <el-table-column prop="clusterIP" label="ClusterIP" min-width="150" sortable="custom" show-overflow-tooltip>
             <template #default="{ row }">
                 {{ row.clusterIP || '<none>' }}
             </template>
         </el-table-column>
         <el-table-column prop="externalIP" label="External IP" min-width="180" sortable="custom" show-overflow-tooltip>
             <template #default="{ row }">
                 {{ row.externalIP || '<none>' }}
             </template>
         </el-table-column>
        <el-table-column prop="ports" label="端口" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">
               <div v-if="row.ports && row.ports.length > 0">
                    <el-tag
                        v-for="(port, index) in row.ports" :key="index"
                        size="small" type="info" effect="plain"
                        style="margin-right: 5px; margin-bottom: 3px;"
                    >
                       {{ formatPort(port) }}
                   </el-tag>
               </div>
                <span v-else>-</span>
            </template>
        </el-table-column>
        <el-table-column prop="selector" label="选择器 (Selector)" min-width="200" show-overflow-tooltip>
             <template #default="{ row }">
                <div v-if="row.selector && Object.keys(row.selector).length > 0">
                     <el-tag
                        v-for="(value, key) in row.selector" :key="key"
                        size="small" type="primary" effect="light"
                        style="margin-right: 5px; margin-bottom: 3px;"
                    >
                        {{ key }}={{ value }}
                    </el-tag>
                </div>
                 <span v-else>-</span>
             </template>
         </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" min-width="180" sortable="custom" />
        <el-table-column label="操作" width="130" align="center" fixed="right">
            <template #default="{ row }">
                 <el-tooltip content="编辑 YAML" placement="top">
                    <el-button link type="primary" :icon="EditPenIcon" @click="editServiceYaml(row)" />
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                    <el-button link type="danger" :icon="DeleteIcon" @click="handleDeleteService(row)" />
                </el-tooltip>
            </template>
        </el-table-column>
         <template #empty>
          <el-empty v-if="!loading.services && !selectedNamespace" description="请先选择一个命名空间以加载 Services" image-size="100" />
          <el-empty v-else-if="!loading.services && paginatedData.length === 0" :description="`在命名空间 '${selectedNamespace}' 中未找到 Services`" image-size="100" />
         </template>
    </el-table>

    <!-- Pagination -->
    <div class="pagination-container" v-if="!loading.services && totalServices > 0">
        <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="totalServices"
            layout="total, sizes, prev, pager, next, jumper"
            background
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
            :disabled="loading.services"
        />
    </div>

    <!-- Create/Edit Dialog (YAML focus) -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="70%" :close-on-click-modal="false">
       <el-alert type="info" :closable="false" style="margin-bottom: 20px;">
         请在此处粘贴或编辑 Service 的 YAML 配置。确保 YAML 中的 `namespace` 与当前选定的命名空间 (`${selectedNamespace || '未选定'}`) 匹配或省略。
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

import {
    Plus as PlusIcon, Search as SearchIcon, Refresh as RefreshIcon, Service as ServiceIcon,
    EditPen as EditPenIcon, Delete as DeleteIcon, View as ViewIcon // Added View for potential future use
} from '@element-plus/icons-vue'

// --- Interfaces ---
// K8s Structures based on sample
interface K8sMetadata { name: string; namespace: string; uid: string; resourceVersion: string; creationTimestamp: string; labels?: { [key: string]: string }; annotations?: { [key: string]: string }; managedFields?: any[] }
interface K8sServicePort { name?: string; protocol?: string; port: number; targetPort?: number | string; nodePort?: number }
interface K8sLoadBalancerStatus { ingress?: { ip?: string; hostname?: string }[] }
interface K8sServiceSpec { ports?: K8sServicePort[]; selector?: { [key: string]: string }; clusterIP?: string; clusterIPs?: string[]; type?: string; sessionAffinity?: string; externalTrafficPolicy?: string; internalTrafficPolicy?: string; ipFamilies?: string[]; ipFamilyPolicy?: string; externalIPs?: string[]; loadBalancerIP?: string; } // Added more fields
interface K8sServiceStatus { loadBalancer?: K8sLoadBalancerStatus }

// API Response Item
interface ServiceApiItem {
  metadata: K8sMetadata;
  spec: K8sServiceSpec;
  status: K8sServiceStatus;
}
interface ServiceListApiResponseData { items: ServiceApiItem[]; total?: number; metadata?: { resourceVersion?: string } }
interface ServiceApiResponse { code: number; data: ServiceListApiResponseData; message: string }
interface NamespaceListResponse { code: number; data: string[]; message: string }

// Internal Display/Table Item
interface ServiceDisplayItem {
  uid: string
  name: string
  namespace: string
  type: string
  clusterIP: string
  externalIP: string // Combined from status.loadBalancer.ingress or spec.externalIPs
  ports: K8sServicePort[] // Keep original ports array
  selector: { [key: string]: string } | null
  createdAt: string
  rawData: ServiceApiItem // Store raw API item
}

// --- Reactive State ---
const allServices = ref<ServiceDisplayItem[]>([])
const namespaces = ref<string[]>([])
const selectedNamespace = ref<string>("")
const currentPage = ref(1)
const pageSize = ref(10)
const totalServices = ref(0)
const searchQuery = ref("")
const sortParams = reactive({ key: 'createdAt', order: 'descending' as ('ascending' | 'descending' | null) })

const loading = reactive({
    page: false, namespaces: false, services: false, dialogSave: false
})

// Dialog state (YAML focus)
const dialogVisible = ref(false)
const dialogTitle = ref("创建 Service (YAML)");
const currentEditService = ref<ServiceApiItem | null>(null);
const yamlContent = ref("");
const placeholderYaml = computed(() => `apiVersion: v1
kind: Service
metadata:
  name: my-new-service
  namespace: ${selectedNamespace.value || 'default'}
spec:
  selector:
    app: my-app # Must match the labels of the Pods you want to target
  ports:
    - protocol: TCP
      port: 80       # Port exposed by the Service
      targetPort: 80 # Port on the Pods to forward traffic to
  # type: ClusterIP # Default type. Others: NodePort, LoadBalancer, ExternalName
`);


// --- Computed Properties ---
const filteredData = computed(() => {
    const query = searchQuery.value.trim().toLowerCase()
    if (!query) return allServices.value;
    return allServices.value.filter(svc =>
        svc.name.toLowerCase().includes(query) ||
        (svc.type && svc.type.toLowerCase().includes(query)) ||
        (svc.clusterIP && svc.clusterIP.toLowerCase().includes(query)) ||
        (svc.externalIP && svc.externalIP.toLowerCase().includes(query))
    );
});

const sortedData = computed(() => { /* ... similar to PVC/Pod, handle string sort for type/IP etc. ... */
    const data = [...filteredData.value];
    const { key, order } = sortParams;
    if (!key || !order) return data;

    data.sort((a, b) => {
        let valA: any; let valB: any;
        if (key === 'createdAt') {
            const timeA = a.createdAt ? dayjs(a.createdAt, "YYYY-MM-DD HH:mm:ss").valueOf() : 0;
            const timeB = b.createdAt ? dayjs(b.createdAt, "YYYY-MM-DD HH:mm:ss").valueOf() : 0;
            valA = isNaN(timeA) ? 0 : timeA; valB = isNaN(timeB) ? 0 : timeB;
        } else {
            valA = a[key as keyof ServiceDisplayItem] ?? ''; valB = b[key as keyof ServiceDisplayItem] ?? '';
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

// Format service ports for display
const formatPort = (port: K8sServicePort): string => {
  let str = `${port.port}`;
  if (port.nodePort) {
    str += `:${port.nodePort}`;
  }
  str += `/${port.protocol || 'TCP'}`;
  if (port.targetPort && port.targetPort !== port.port) {
    str += ` -> ${port.targetPort}`;
  }
  if (port.name) {
      str += ` (${port.name})`;
  }
  return str;
};

// Get External IP (simplistic, handles LoadBalancer primarily)
const getExternalIP = (item: ServiceApiItem): string => {
    if (item.spec?.type === 'LoadBalancer' && item.status?.loadBalancer?.ingress?.length > 0) {
        // Prefer IP, fallback to hostname
        return item.status.loadBalancer.ingress.map(ing => ing.ip || ing.hostname).filter(Boolean).join(', ');
    }
    if (item.spec?.externalIPs && item.spec.externalIPs.length > 0) {
        return item.spec.externalIPs.join(', ');
    }
    return ''; // Return empty string if no external IP found
}

// Get tag type based on service type
const getTypeTagType = (type: string): '' | 'success' | 'warning' | 'info' | 'danger' => {
    switch(type) {
        case 'ClusterIP': return 'info';
        case 'NodePort': return ''; // Default/primary
        case 'LoadBalancer': return 'success';
        case 'ExternalName': return 'warning';
        default: return 'info';
    }
}

const VITE_API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://1992.168.1.100:8080'; // Ensure this is set in your environment


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

const fetchServiceData = async () => {
    if (!selectedNamespace.value) { allServices.value = []; totalServices.value = 0; return; }
    loading.services = true;
    try {
        const params: Record<string, any> = { /* Server-side params */ };
        const url = `/api/v1/namespaces/${selectedNamespace.value}/services`;
        const response = await request<ServiceApiResponse>({ url, method: "get", params, baseURL: VITE_API_BASE_URL });

        if (response.code === 200 && response.data?.items) {
            totalServices.value = response.data.total ?? response.data.items.length;
            allServices.value = response.data.items.map((item, index) => ({
                uid: item.metadata.uid || `${item.metadata.namespace}-${item.metadata.name}-${index}`,
                name: item.metadata.name,
                namespace: item.metadata.namespace,
                type: item.spec?.type || 'ClusterIP', // Default to ClusterIP
                clusterIP: item.spec?.clusterIP || '',
                externalIP: getExternalIP(item),
                ports: item.spec?.ports || [],
                selector: item.spec?.selector || null,
                createdAt: formatTimestamp(item.metadata.creationTimestamp),
                rawData: item,
            }));
             const totalPages = Math.ceil(totalServices.value / pageSize.value);
             if (currentPage.value > totalPages && totalPages > 0) currentPage.value = totalPages;
             else if (totalServices.value === 0) currentPage.value = 1;
        } else {
            ElMessage.error(`获取 Service 数据失败: ${response.message || '未知错误'}`);
            allServices.value = []; totalServices.value = 0;
        }
    } catch (error: any) {
        console.error("获取 Service 数据失败:", error);
        ElMessage.error(`获取 Service 数据出错: ${error.message || '网络请求失败'}`);
        allServices.value = []; totalServices.value = 0;
    } finally {
        loading.services = false;
    }
}

// --- Event Handlers ---
const handleNamespaceChange = () => { /* ... */ currentPage.value = 1; searchQuery.value = ''; sortParams.key = 'createdAt'; sortParams.order = 'descending'; fetchServiceData(); };
const handlePageChange = (page: number) => { currentPage.value = page; /* Fetch if server-side */ };
const handleSizeChange = (size: number) => { pageSize.value = size; currentPage.value = 1; /* Fetch if server-side */ };
const handleSearchDebounced = debounce(() => { currentPage.value = 1; /* Fetch if server-side */ }, 300);
const handleSortChange = ({ prop, order }: { prop: string | null; order: 'ascending' | 'descending' | null }) => { /* ... */ sortParams.key = prop || 'createdAt'; sortParams.order = order; currentPage.value = 1; };


// --- Dialog and CRUD Actions ---
const handleAddService = () => { /* ... */ if (!selectedNamespace.value) { ElMessage.warning("请先选择一个命名空间"); return; } currentEditService.value = null; yamlContent.value = placeholderYaml.value; dialogTitle.value = "创建 Service (YAML)"; dialogVisible.value = true; };
const editServiceYaml = async (service: ServiceDisplayItem) => { /* ... */
    ElMessage.info(`模拟: 获取 Service "${service.name}" 的 YAML`);
    currentEditService.value = service.rawData;
    yamlContent.value = service.rawData ? yaml.dump(service.rawData) : placeholderYaml.value;
    dialogTitle.value = `编辑 Service: ${service.name} (YAML)`;
    dialogVisible.value = true;
};

const handleSaveYaml = async () => { /* ... */
    if (!selectedNamespace.value && !currentEditService.value?.metadata.namespace) { ElMessage.error("无法确定目标命名空间。"); return; }
    loading.dialogSave = true;
    // --- Replace with actual YAML editor interaction and API call (POST/PUT) ---
     await new Promise(resolve => setTimeout(resolve, 500)); // Simulate
     loading.dialogSave = false; dialogVisible.value = false;
     const action = currentEditService.value ? '更新' : '创建';
     ElMessage.success(`模拟 Service ${action}成功`); fetchServiceData();
};

const handleDeleteService = (service: ServiceDisplayItem) => { /* ... */
    ElMessageBox.confirm(
        `确定要删除 Service "${service.name}" (命名空间: ${service.namespace}) 吗？`,
        '确认删除', { type: 'warning' }
    ).then(async () => {
        loading.services = true;
        try {
            const response = await request<{ code: number; message: string }>({
                url: `/api/v1/namespaces/${service.namespace}/services/${service.name}`,
                method: "delete",
                baseURL: VITE_API_BASE_URL,
            });
             if (response.code === 200 || response.code === 204 || response.code === 202) {
                ElMessage.success(`Service "${service.name}" 已删除`); await fetchServiceData();
            } else { ElMessage.error(`删除 Service 失败: ${response.message || '未知错误'}`); loading.services = false; }
        } catch (error: any) { console.error("删除 Service 失败:", error); ElMessage.error(`删除 Service 失败: ${error.response?.data?.message || error.message || '请求失败'}`); loading.services = false; }
    }).catch(() => ElMessage.info('删除操作已取消'));
};


// --- Lifecycle Hooks ---
onMounted(async () => { /* ... */
    loading.page = true;
    await fetchNamespaces();
    if (selectedNamespace.value) { await fetchServiceData(); }
    loading.page = false;
});

</script>

<style lang="scss" scoped>
/* Using fallback variables directly */
$page-padding: 20px; $spacing-md: 15px; $spacing-lg: 20px;
$font-size-base: 14px; $font-size-small: 12px; $font-size-large: 16px; $font-size-extra-large: 24px;
$border-radius-base: 4px; $kube-service-icon-color: #3498DB; // Example Blue

.service-page-container { padding: $page-padding; background-color: var(--el-bg-color-page); }
.page-breadcrumb { margin-bottom: $spacing-lg; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: $spacing-md; flex-wrap: wrap; gap: $spacing-md; }
.page-title { font-size: $font-size-extra-large; font-weight: 600; color: var(--el-text-color-primary); margin: 0; }
.info-alert { margin-bottom: $spacing-lg; background-color: var(--el-color-info-light-9); :deep(.el-alert__description) { font-size: $font-size-small; color: var(--el-text-color-regular); line-height: 1.6; } }
.filter-bar { display: flex; align-items: center; flex-wrap: wrap; gap: $spacing-md; margin-bottom: $spacing-lg; padding: $spacing-md; background-color: var(--el-bg-color); border-radius: $border-radius-base; border: 1px solid var(--el-border-color-lighter); }
.filter-item { }
.namespace-select { width: 240px; }
.search-input { width: 300px; }

.service-table {
   border-radius: $border-radius-base; border: 1px solid var(--el-border-color-lighter); overflow: hidden;
    :deep(th.el-table__cell) { background-color: var(--el-fill-color-lighter); color: var(--el-text-color-secondary); font-weight: 600; font-size: $font-size-small; }
    :deep(td.el-table__cell) { padding: 8px 0; font-size: $font-size-base; vertical-align: middle; }
   .service-icon { margin-right: 8px; color: $kube-service-icon-color; vertical-align: middle; font-size: 18px; position: relative; top: -1px; }
   .service-name { font-weight: 500; vertical-align: middle; color: var(--el-text-color-regular); }
}

.el-table .el-button.is-link { font-size: 14px; padding: 4px; margin: 0 3px; vertical-align: middle; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: $spacing-lg; }
.yaml-editor-placeholder { border: 1px dashed var(--el-border-color); padding: 20px; margin-top: 10px; min-height: 350px; max-height: 60vh; background-color: var(--el-fill-color-lighter); color: var(--el-text-color-secondary); font-family: monospace; white-space: pre-wrap; overflow: auto; font-size: 13px; line-height: 1.5; }
.dialog-footer { text-align: right; padding-top: 10px; }
</style>