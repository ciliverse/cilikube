<template>
  <div class="ingress-page-container">
    <!-- Header: Title & Create Button -->
    <div class="page-header">
      <h1 class="page-title">{{ $t('ingressManagement.title') }}</h1>
       <el-button
         type="primary"
         :icon="PlusIcon"
         @click="handleAddIngress"
         :loading="loading.page"
         :disabled="!selectedNamespace"
       >
         {{ $t('ingressManagement.actions.create') }}
       </el-button>
    </div>

    <!-- Filter Bar: Namespace, Search, Refresh -->
    <div class="filter-bar">
        <el-select
            v-model="selectedNamespace"
            :placeholder="$t('ingressManagement.filterNamespace')"
            @change="handleNamespaceChange"
            filterable
            :loading="loading.namespaces"
            class="filter-item namespace-select"
            style="width: 280px;"
        >
            <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
             <template #empty>
                <div style="padding: 10px; text-align: center; color: #999;">
                    {{ loading.namespaces ? $t('common.loading') : $t('ingressManagement.messages.noNamespaces') }}
                </div>
            </template>
        </el-select>

        <el-input
            v-model="searchQuery"
            :placeholder="$t('ingressManagement.searchPlaceholder')"
            :prefix-icon="SearchIcon"
            clearable
            @input="handleSearchDebounced"
            class="filter-item search-input"
            :disabled="!selectedNamespace || loading.ingresses"
        />

        <el-tooltip :content="$t('common.refresh')" placement="top">
            <el-button
              :icon="RefreshIcon"
              circle
              @click="fetchIngressData"
              :loading="loading.ingresses"
              :disabled="!selectedNamespace"
            />
        </el-tooltip>
    </div>

    <!-- Ingresses Table Card -->
    <el-card class="ingress-table-card" shadow="hover">
      <el-table
          :data="paginatedData"
          v-loading="loading.ingresses"
          border
          stripe
          style="width: 100%"
          @sort-change="handleSortChange"
          class="ingress-table"
          :default-sort="{ prop: 'createdAt', order: 'descending' }"
          row-key="uid"
      >
        <el-table-column :label="$t('ingressManagement.columns.name')" prop="name" min-width="250" sortable="custom" fixed show-overflow-tooltip>
             <template #default="{ row }">
                <el-icon class="ingress-icon"><ConnectionIcon /></el-icon>
                <span class="ingress-name">{{ row.name }}</span>
            </template>
        </el-table-column>
        <el-table-column :label="$t('ingressManagement.columns.namespace')" prop="namespace" min-width="150" sortable="custom" show-overflow-tooltip />
        <el-table-column :label="$t('ingressManagement.columns.class')" prop="ingressClassName" min-width="150" sortable="custom" show-overflow-tooltip>
             <template #default="{ row }">
                <el-tag v-if="row.ingressClassName" size="small" type="info">{{ row.ingressClassName }}</el-tag>
                <span v-else>-</span>
            </template>
        </el-table-column>
         <el-table-column :label="$t('ingressManagement.columns.hosts')" prop="hosts" min-width="250" show-overflow-tooltip>
              <template #default="{ row }">
                 <div v-if="row.hosts && row.hosts.length > 0">
                     <span v-for="(host, index) in row.hosts" :key="index">
                         <el-link type="primary" :href="'http://' + host" target="_blank">{{ host }}</el-link>
                         <span v-if="index < row.hosts.length - 1">, </span>
                     </span>
                 </div>
                 <span v-else>*</span> <!-- Indicate default backend or rules without host -->
             </template>
         </el-table-column>
         <el-table-column :label="$t('ingressManagement.columns.address')" prop="address" min-width="180" show-overflow-tooltip>
              <template #default="{ row }">
                 {{ row.address || '-' }}
             </template>
         </el-table-column>
        <el-table-column :label="$t('ingressManagement.columns.ports')" prop="ports" min-width="100" align="center">
            <template #default>
               <el-tag size="small" type="info" effect="light">80, 443</el-tag> <!-- Typically 80/443 for Ingress -->
            </template>
        </el-table-column>
        <el-table-column :label="$t('ingressManagement.columns.rules')" prop="rules" min-width="300" show-overflow-tooltip>
             <template #default="{ row }">
                <div v-if="row.rules && row.rules.length > 0" class="rules-column">
                   <div v-for="(rule, index) in row.rules" :key="index" class="rule-item">
                        <span class="rule-path">{{ rule.path }}</span>
                        <el-icon><Right /></el-icon>
                        <span class="rule-backend">{{ rule.backendService }}</span>:
                        <span class="rule-port">{{ rule.backendPort }}</span>
                         <span v-if="rule.host" class="rule-host"> (Host: {{ rule.host }})</span>
                   </div>
                </div>
                 <span v-else-if="row.defaultBackend">Default Backend: {{ row.defaultBackend }}</span>
                 <span v-else>-</span>
             </template>
        </el-table-column>
        <el-table-column :label="$t('ingressManagement.columns.createdAt')" prop="createdAt" min-width="180" sortable="custom" />
        <el-table-column :label="$t('ingressManagement.columns.actions')" width="130" align="center" fixed="right">
            <template #default="{ row }">
                 <el-tooltip :content="$t('ingressManagement.actions.editYaml')" placement="top">
                    <el-button link type="primary" :icon="EditPenIcon" @click="editIngressYaml(row)" />
                </el-tooltip>
                <el-tooltip :content="$t('common.delete')" placement="top">
                    <el-button link type="danger" :icon="DeleteIcon" @click="handleDeleteIngress(row)" />
                </el-tooltip>
            </template>
        </el-table-column>
         <template #empty>
          <el-empty v-if="!loading.ingresses && !selectedNamespace" :description="$t('ingressManagement.messages.selectNamespace')" image-size="100" />
          <el-empty v-else-if="!loading.ingresses && paginatedData.length === 0" :description="`${$t('ingressManagement.messages.noIngresses')}`" image-size="100" />
         </template>
    </el-table>

      <!-- Pagination -->
      <div class="pagination-container" v-if="!loading.ingresses && totalIngresses > 0">
        <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="totalIngresses"
            layout="total, sizes, prev, pager, next, jumper"
            background
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
            :disabled="loading.ingresses"
        />
      </div>
    </el-card>

    <!-- Create/Edit Dialog (YAML focus) -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="70%" :close-on-click-modal="false">
       <el-alert type="info" :closable="false" style="margin-bottom: 20px;">
         {{ $t('ingressManagement.messages.yamlEditTip') }}
       </el-alert>
       <!-- Integrate a YAML editor component here -->
       <div class="yaml-editor-placeholder">
            {{ $t('ingressManagement.messages.yamlEditorPlaceholder') }}
            <pre>{{ yamlContent || placeholderYaml }}</pre>
       </div>
      <template #footer>
        <div class="dialog-footer">
            <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
            <el-button type="primary" @click="handleSaveYaml" :loading="loading.dialogSave">{{ $t('ingressManagement.actions.applyYaml') }}</el-button>
        </div>
      </template>
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from "vue"
import { ElMessage, ElMessageBox } from "element-plus"
import { useI18n } from 'vue-i18n'
import { kubernetesRequest, fetchNamespaces, KubernetesApiResponse } from "@/utils/api-config"
import dayjs from "dayjs"
import { debounce } from 'lodash-es'
import yaml from 'js-yaml'; // Ensure installed

import {
    Plus as PlusIcon, Search as SearchIcon, Refresh as RefreshIcon, Connection as ConnectionIcon,
    EditPen as EditPenIcon, Delete as DeleteIcon, Right
} from '@element-plus/icons-vue'

const { t } = useI18n()

// --- Interfaces ---
// K8s Structures based on API sample (networking.k8s.io/v1)
interface K8sMetadata { name: string; namespace: string; uid: string; resourceVersion: string; generation?: number; creationTimestamp: string; labels?: { [key: string]: string }; annotations?: { [key: string]: string }; managedFields?: any[] }
interface K8sIngressBackend { service?: { name: string; port?: { name?: string; number?: number } }; resource?: { apiGroup?: string; kind: string; name: string } }
interface K8sHTTPIngressPath { path?: string; pathType?: string; backend: K8sIngressBackend }
interface K8sHTTPIngressRuleValue { paths: K8sHTTPIngressPath[] }
interface K8sIngressRule { host?: string; http?: K8sHTTPIngressRuleValue }
interface K8sIngressTLS { hosts?: string[]; secretName?: string }
interface K8sIngressSpec { ingressClassName?: string | null; defaultBackend?: K8sIngressBackend; tls?: K8sIngressTLS[]; rules?: K8sIngressRule[] }
interface K8sLoadBalancerIngress { ip?: string; hostname?: string; ports?: { port: number; protocol: string; error?: string }[] }
interface K8sLoadBalancerStatus { ingress?: K8sLoadBalancerIngress[] }
interface K8sIngressStatus { loadBalancer?: K8sLoadBalancerStatus }

// API Response Item
interface IngressApiItem {
  metadata: K8sMetadata;
  spec: K8sIngressSpec;
  status: K8sIngressStatus;
}
interface IngressListApiResponseData { items: IngressApiItem[]; total?: number; metadata?: { resourceVersion?: string } }
interface IngressApiResponse { code: number; data: IngressListApiResponseData; message: string }
interface NamespaceListResponse { code: number; data: { items: Array<{ metadata: { name: string } }> }; message: string }

// Simple Rule structure for display
interface SimpleRule { host?: string; path?: string; backendService?: string; backendPort?: string | number; }

// Internal Display/Table Item
interface IngressDisplayItem {
  uid: string
  name: string
  namespace: string
  ingressClassName: string
  hosts: string[]
  address: string // LoadBalancer IP(s) or Hostname(s)
  ports: string // Usually "80, 443"
  rules: SimpleRule[]
  defaultBackend: string // Formatted default backend if rules are empty
  createdAt: string
  rawData: IngressApiItem // Store raw API item
}

// --- Reactive State ---
const allIngresses = ref<IngressDisplayItem[]>([])
const namespaces = ref<string[]>([])
const selectedNamespace = ref<string>("")
const currentPage = ref(1)
const pageSize = ref(10)
const totalIngresses = ref(0)
const searchQuery = ref("")
const sortParams = reactive({ key: 'createdAt', order: 'descending' as ('ascending' | 'descending' | null) })

const loading = reactive({
    page: false, namespaces: false, ingresses: false, dialogSave: false
})

// Dialog state (YAML focus)
const dialogVisible = ref(false)
const currentEditIngress = ref<IngressApiItem | null>(null);
const yamlContent = ref("");

const dialogTitle = computed(() => 
  currentEditIngress.value 
    ? `${t('ingressManagement.actions.edit')}: ${currentEditIngress.value.metadata.name} (YAML)`
    : t('ingressManagement.messages.createIngressTitle')
);

const placeholderYaml = computed(() => `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-new-ingress
  namespace: ${selectedNamespace.value || 'default'}
  # annotations:
  #   nginx.ingress.kubernetes.io/rewrite-target: / # Example for Nginx Ingress
spec:
  ingressClassName: nginx # Or your specific Ingress Class name
  rules:
  - host: myapp.example.com # Optional: host for the rule
    http:
      paths:
      - path: / # Path for this rule
        pathType: Prefix # Or Exact, ImplementationSpecific
        backend:
          service:
            name: my-service # Name of the backend Service
            port:
              number: 80 # Port number of the backend Service
  # tls: # Optional TLS configuration
  # - hosts:
  #   - myapp.example.com
  #   secretName: myapp-tls-secret # Secret containing TLS cert and key
`);


// --- Computed Properties ---
const filteredData = computed(() => {
    const query = searchQuery.value.trim().toLowerCase()
    if (!query) return allIngresses.value;
    return allIngresses.value.filter(ing =>
        ing.name.toLowerCase().includes(query) ||
        ing.ingressClassName.toLowerCase().includes(query) ||
        ing.hosts.some(h => h.toLowerCase().includes(query)) ||
        (ing.address && ing.address.toLowerCase().includes(query))
    );
});

const sortedData = computed(() => { /* ... similar sort logic, handle string/date ... */
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
            valA = a[key as keyof IngressDisplayItem] ?? ''; valB = b[key as keyof IngressDisplayItem] ?? '';
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


const paginatedData = computed(() => { /* ... */
    const start = (currentPage.value - 1) * pageSize.value;
    const end = start + pageSize.value;
    return sortedData.value.slice(start, end);
});


// --- Helper Functions ---
const formatTimestamp = (timestamp: string): string => { /* ... */ if (!timestamp) return 'N/A'; return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss"); }

// Extract Ingress LoadBalancer Addresses
const getAddress = (status: K8sIngressStatus | undefined): string => {
    if (!status?.loadBalancer?.ingress || status.loadBalancer.ingress.length === 0) {
        return '';
    }
    return status.loadBalancer.ingress.map(ing => ing.ip || ing.hostname).filter(Boolean).join(', ');
}

// Extract Hosts from rules
const getHosts = (rules: K8sIngressRule[] | undefined): string[] => {
    if (!rules) return [];
    return rules.map(rule => rule.host || '').filter(Boolean); // Get non-empty hosts
}

// Format backend service/port string
const formatBackend = (backend: K8sIngressBackend): string => {
    if (backend.service) {
        const port = backend.service.port?.number || backend.service.port?.name || '?';
        return `${backend.service.name}:${port}`;
    }
    if (backend.resource) {
        return `${backend.resource.kind}:${backend.resource.name}`;
    }
    return 'N/A';
}

// Simplify rules for display
const simplifyRules = (rules: K8sIngressRule[] | undefined): SimpleRule[] => {
    if (!rules) return [];
    const simpleRules: SimpleRule[] = [];
    rules.forEach(rule => {
        rule.http?.paths.forEach(path => {
            simpleRules.push({
                host: rule.host,
                path: path.path || '*',
                backendService: path.backend.service?.name || path.backend.resource?.name || '?',
                backendPort: path.backend.service?.port?.number || path.backend.service?.port?.name || path.backend.resource?.kind || '?' // Less ideal for resource, but shows something
            });
        });
    });
    return simpleRules;
}

// --- API Interaction ---
const fetchNamespacesList = async () => {
    loading.namespaces = true;
    try {
        const namespaceList = await fetchNamespaces();
        namespaces.value = namespaceList;
        if (namespaceList.length > 0 && !selectedNamespace.value) {
             selectedNamespace.value = namespaceList.find(ns => ns === 'default') || namespaceList[0];
        } else if (namespaceList.length === 0) { 
            ElMessage.warning(t('ingressManagement.messages.noNamespaces')); 
        }
    } catch (error: any) { 
        ElMessage.error(error.message); 
        namespaces.value = []; 
    }
    finally { loading.namespaces = false; }
}

const fetchIngressData = async () => {
    if (!selectedNamespace.value) { allIngresses.value = []; totalIngresses.value = 0; return; }
    loading.ingresses = true;
    try {
        const response = await kubernetesRequest<IngressApiResponse>({ 
            url: `/api/v1/namespaces/${selectedNamespace.value}/ingresses`,
            method: "get"
        });

        if (response.code === 200) {
            if (response.data?.items) {
                totalIngresses.value = response.data.total ?? response.data.items.length;
                allIngresses.value = response.data.items.map((item, index) => ({
                    uid: item.metadata.uid || `${item.metadata.namespace}-${item.metadata.name}-${index}`,
                    name: item.metadata.name,
                    namespace: item.metadata.namespace,
                    ingressClassName: item.spec?.ingressClassName || '',
                    hosts: getHosts(item.spec?.rules),
                    address: getAddress(item.status),
                    ports: '80, 443', // Fixed for typical Ingress
                    rules: simplifyRules(item.spec?.rules),
                    defaultBackend: item.spec?.defaultBackend ? formatBackend(item.spec.defaultBackend) : '',
                    createdAt: formatTimestamp(item.metadata.creationTimestamp),
                    rawData: item,
                }));
                const totalPages = Math.ceil(totalIngresses.value / pageSize.value);
                if (currentPage.value > totalPages && totalPages > 0) currentPage.value = totalPages;
                else if (totalIngresses.value === 0) currentPage.value = 1;
            } else {
                // Data is null/empty but HTTP status is 200
                allIngresses.value = [];
                totalIngresses.value = 0;
                // ElMessage.warning('获取到空的 Ingress 数据');
            }
        } else {
            ElMessage.error(`${t('ingressManagement.messages.fetchIngressesError')}: ${response.message || t('common.error')}`);
            allIngresses.value = [];
            totalIngresses.value = 0;
        }
    } catch (error: any) {
        console.error("获取 Ingress 数据失败:", error);
        // Handle 404 specifically for Ingress if the API group isn't enabled
        if (error.response?.status === 404) {
            ElMessage.warning(`${t('ingressManagement.messages.fetchIngressesError')}: 请确保 Ingress API (networking.k8s.io/v1) 在集群中已启用，并且您有读取权限。`);
        } else {
            ElMessage.error(`${t('ingressManagement.messages.fetchIngressesError')}: ${error.message || t('common.error')}`);
        }
        allIngresses.value = []; totalIngresses.value = 0;
    } finally {
        loading.ingresses = false;
    }
}

// --- Event Handlers ---
const handleNamespaceChange = () => { /* ... */ currentPage.value = 1; searchQuery.value = ''; sortParams.key = 'createdAt'; sortParams.order = 'descending'; fetchIngressData(); };
const handlePageChange = (page: number) => { currentPage.value = page; /* Fetch if server-side */ };
const handleSizeChange = (size: number) => { pageSize.value = size; currentPage.value = 1; /* Fetch if server-side */ };
const handleSearchDebounced = debounce(() => { currentPage.value = 1; /* Fetch if server-side */ }, 300);
const handleSortChange = ({ prop, order }: { prop: string | null; order: 'ascending' | 'descending' | null }) => { /* ... */ sortParams.key = prop || 'createdAt'; sortParams.order = order; currentPage.value = 1; };


// --- Dialog and CRUD Actions ---
const handleAddIngress = () => { /* ... */ if (!selectedNamespace.value) { ElMessage.warning(t('ingressManagement.messages.selectNamespaceWarning')); return; } currentEditIngress.value = null; yamlContent.value = placeholderYaml.value; dialogVisible.value = true; };
const editIngressYaml = async (ingress: IngressDisplayItem) => { /* ... */
    ElMessage.info(`${t('ingressManagement.messages.fetchingYaml').replace('{name}', ingress.name)}`);
    currentEditIngress.value = ingress.rawData;
    yamlContent.value = ingress.rawData ? yaml.dump(ingress.rawData) : placeholderYaml.value;
    dialogVisible.value = true;
};

const handleSaveYaml = async () => { /* ... */
    if (!selectedNamespace.value && !currentEditIngress.value?.metadata.namespace) { ElMessage.error(t('ingressManagement.messages.noNamespaces')); return; }
    loading.dialogSave = true;
    // --- Replace with actual YAML editor interaction and API call (POST/PUT) ---
     await new Promise(resolve => setTimeout(resolve, 500)); // Simulate
     loading.dialogSave = false; dialogVisible.value = false;
     const action = currentEditIngress.value ? t('common.edit') : t('common.add');
     ElMessage.success(t('ingressManagement.messages.operationSuccess').replace('{action}', action)); fetchIngressData();
};

const handleDeleteIngress = (ingress: IngressDisplayItem) => { /* ... */
    ElMessageBox.confirm(
        t('ingressManagement.messages.deleteConfirmDetail').replace('{name}', ingress.name).replace('{namespace}', ingress.namespace),
        t('common.confirm'), { type: 'warning' }
    ).then(async () => {
        loading.ingresses = true;
        try {
            const response = await kubernetesRequest<{ code: number; message: string }>({
                url: `/api/v1/namespaces/${ingress.namespace}/ingresses/${ingress.name}`,
                method: "delete"
            });
             if (response.code === 200 || response.code === 204 || response.code === 202) {
                ElMessage.success(`Ingress "${ingress.name}" ${t('common.delete')}${t('common.success')}`); await fetchIngressData();
            } else { ElMessage.error(`${t('ingressManagement.messages.deleteFailed')}: ${response.message || t('common.error')}`); loading.ingresses = false; }
        } catch (error: any) { console.error("删除 Ingress 失败:", error); ElMessage.error(`${t('ingressManagement.messages.deleteFailed')}: ${error.response?.data?.message || error.message || '请求失败'}`); loading.ingresses = false; }
    }).catch(() => ElMessage.info(t('ingressManagement.messages.operationCancelled')));
};

// --- Lifecycle Hooks ---
onMounted(async () => { /* ... */
    loading.page = true;
    await fetchNamespacesList();
    if (selectedNamespace.value) { await fetchIngressData(); }
    loading.page = false;
});

</script>

<style lang="scss" scoped>
/* Using fallback variables directly */
$page-padding: 24px; $spacing-md: 16px; $spacing-lg: 24px;
$font-size-base: 14px; $font-size-small: 12px; $font-size-large: 16px; $font-size-extra-large: 28px;
$border-radius-base: 4px; $kube-ingress-icon-color: #1ABC9C; // Example Turquoise

.ingress-page-container { 
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

.ingress-table-card {
  margin-bottom: $spacing-lg;
  border-radius: $border-radius-base;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
  border: 1px solid var(--el-border-color-lighter);
  
  :deep(.el-card__body) {
    padding: 0;
  }
}

.ingress-table {
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
  .ingress-icon { 
    margin-right: 8px; 
    color: $kube-ingress-icon-color; 
    vertical-align: middle; 
    font-size: 18px; 
    position: relative; 
    top: -1px; 
  }
  .ingress-name { 
    font-weight: 500; 
    vertical-align: middle; 
    color: var(--el-text-color-primary); 
  }

    .rules-column {
        font-size: $font-size-small;
        line-height: 1.5;
        color: var(--el-text-color-secondary);
        .rule-item {
            margin-bottom: 4px;
            &:last-child { margin-bottom: 0; }
        }
        .rule-path { font-weight: 500; color: var(--el-text-color-regular); }
        .rule-backend { font-style: italic; }
        .rule-port { font-weight: 500; }
        .rule-host { margin-left: 5px; }
        .el-icon { vertical-align: middle; margin: 0 4px; }
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
