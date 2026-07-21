<template>
  <div class="namespace-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <h1 class="page-title">{{ $t('namespaceManagement.title') }}</h1>
      <div class="header-actions">
        <el-button :icon="Refresh" @click="fetchNamespaceData" :loading="loading" circle />
        <el-button type="primary" :icon="Plus" @click="showCreateDialog">
          {{ $t('namespaceManagement.actions.create') }}
        </el-button>
      </div>
    </div>

    <!-- 集群选择提示 -->
    <el-alert
      v-if="!selectedClusterName && !s_loadingClusters"
      :title="$t('namespaceManagement.messages.selectCluster')"
      type="info" show-icon :closable="false" style="margin-bottom: 24px;"
    />

    <!-- 工具栏 -->
    <div class="toolbar" v-if="selectedClusterName">
      <div class="search-filters">
        <el-input v-model="searchQuery" :placeholder="$t('namespaceManagement.searchPlaceholder')" :prefix-icon="Search"
          clearable class="search-input" @input="handleSearchDebounced" />
        <el-select v-model="filterStatus" :placeholder="$t('namespaceManagement.filterStatus')" clearable
          class="filter-select">
          <el-option :label="$t('namespaceManagement.allStatuses')" value="" />
          <el-option :label="$t('namespaceManagement.statuses.active')" value="active" />
          <el-option :label="$t('namespaceManagement.statuses.terminating')" value="terminating" />
        </el-select>
      </div>
      <div class="toolbar-right">
        <el-button-group class="view-toggle">
          <el-button :type="viewMode === 'card' ? 'primary' : ''" :icon="Grid" @click="viewMode = 'card'" />
          <el-button :type="viewMode === 'list' ? 'primary' : ''" :icon="List" @click="viewMode = 'list'" />
        </el-button-group>
      </div>
    </div>

    <!-- 命名空间列表 -->
    <div class="namespace-list" v-loading="loading">
      <!-- 卡片视图 -->
      <div v-if="filteredData.length > 0 && viewMode === 'card'" class="namespace-grid">
        <div v-for="namespace in paginatedData" :key="namespace.uid" class="namespace-card"
          @click="handleNamespaceClick(namespace)">
          <div class="card-header">
            <div class="namespace-info">
              <div class="namespace-name-row">
                <div class="status-dot" :class="getStatusClass(namespace.status)"></div>
                <div class="namespace-name">{{ namespace.name }}</div>
                <el-tag v-if="isSystemNamespace(namespace.name)" type="warning" size="small" class="system-badge">
                  {{ $t('namespaceManagement.info.systemNamespace') }}
                </el-tag>
              </div>
              <div class="namespace-status">{{ getStatusText(namespace.status) }}</div>
            </div>
          </div>

          <div class="card-body">
            <div class="namespace-meta">
              <div class="meta-item">
                <span class="meta-label">{{ $t('namespaceManagement.columns.created') }}</span>
                <span class="meta-value">{{ namespace.creationTimestamp }}</span>
              </div>
              <div class="meta-item" v-if="namespace.labels && Object.keys(namespace.labels).length > 0">
                <span class="meta-label">{{ $t('namespaceManagement.columns.labels') }}</span>
                <span class="meta-value">{{ Object.keys(namespace.labels).length }} {{ $t('namespaceManagement.columns.labels') }}</span>
              </div>
            </div>
          </div>

          <div class="card-footer">
            <div class="namespace-actions">
              <el-button type="primary" size="small" @click.stop="viewNamespaceDetails(namespace)">
                {{ $t('namespaceManagement.actions.viewDetails') }}
              </el-button>
              <el-button type="danger" size="small" @click.stop="handleDeleteNamespace(namespace)"
                :disabled="isSystemNamespace(namespace.name)">
                {{ $t('namespaceManagement.actions.delete') }}
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 列表视图 -->
      <div v-else-if="filteredData.length > 0 && viewMode === 'list'" class="namespace-table">
        <el-table :data="paginatedData" @row-click="handleNamespaceClick">
          <el-table-column width="60">
            <template #default="{ row }">
              <div class="status-dot" :class="getStatusClass(row.status)"></div>
            </template>
          </el-table-column>
          <el-table-column prop="name" :label="$t('namespaceManagement.columns.name')" min-width="200">
            <template #default="{ row }">
              <div class="table-namespace-name">
                {{ row.name }}
                <el-tag v-if="isSystemNamespace(row.name)" type="warning" size="small" class="system-badge">
                  {{ $t('namespaceManagement.info.systemNamespace') }}
                </el-tag>
              </div>
            </template>
          </el-table-column>
          <el-table-column :label="$t('namespaceManagement.columns.status')" width="120">
            <template #default="{ row }">
              <span>{{ getStatusText(row.status) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="creationTimestamp" :label="$t('namespaceManagement.columns.created')" width="180" />
          <el-table-column :label="$t('common.actions')" width="200" fixed="right">
            <template #default="{ row }">
              <div class="table-actions">
                <el-button type="primary" size="small" @click.stop="viewNamespaceDetails(row)">
                  {{ $t('namespaceManagement.actions.viewDetails') }}
                </el-button>
                <el-button type="danger" size="small" @click.stop="handleDeleteNamespace(row)"
                  :disabled="isSystemNamespace(row.name)">
                  {{ $t('namespaceManagement.actions.delete') }}
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div v-else-if="!loading && selectedClusterName" class="empty-state">
        <el-empty :description="searchQuery || filterStatus ? 
          $t('namespaceManagement.messages.noMatchingNamespaces') : 
          $t('namespaceManagement.messages.noNamespaces')" />
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination-container" v-if="filteredData.length > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[12, 24, 36, 48]"
        :total="filteredData.length"
        layout="total, sizes, prev, pager, next, jumper"
        background
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>

    <!-- 创建命名空间对话框 -->
    <el-dialog v-model="isDialogVisible" :title="$t('namespaceManagement.actions.create')" width="600px"
      :close-on-click-modal="false" @closed="resetForm" class="namespace-dialog">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="120px" class="namespace-form">
        <div class="form-section">
          <h4 class="section-title">{{ $t('clusterManagement.basicInfo') }}</h4>
          <el-form-item :label="$t('namespaceManagement.form.name')" prop="name">
            <el-input v-model="form.name" :placeholder="$t('namespaceManagement.form.namePlaceholder')" clearable />
            <div class="form-item-help">
              {{ $t('namespaceManagement.form.nameFormat') }}
            </div>
          </el-form-item>
        </div>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="isDialogVisible = false" size="large">
            {{ $t('common.cancel') }}
          </el-button>
          <el-button type="primary" @click="createNamespace" :loading="createLoading" size="large">
            {{ $t('common.confirm') }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted, reactive, nextTick, watch } from "vue"
  import { useI18n } from 'vue-i18n'
  import { ElMessage, ElMessageBox, ElLoading } from "element-plus"
  import type { FormInstance, FormRules } from 'element-plus'
  import { kubernetesRequest, KubernetesApiResponse } from "@/utils/api-config"
  import dayjs from "dayjs"
  import { debounce } from 'lodash-es'
  import { storeToRefs } from "pinia"
  import { useClusterStore } from "@/store/modules/clusterStore"
  
  import {
    Search,
    Plus,
    Refresh,
    Grid,
    List
  } from '@element-plus/icons-vue'

  const { t } = useI18n()
  
  // --- Interfaces based on API response ---
  interface K8sMetadata {
    name: string
    uid: string
    resourceVersion: string
    creationTimestamp: string
    labels?: { [key: string]: string }
    annotations?: { [key: string]: string }
    managedFields?: any[] // Array of managed fields details
  }
  
  interface K8sNamespaceSpec {
    finalizers: string[]
  }
  
  interface K8sNamespaceStatus {
    phase: 'Active' | 'Terminating' | string // API returns string, define known values
  }
  
  interface K8sNamespace {
    metadata: K8sMetadata
    spec: K8sNamespaceSpec
    status: K8sNamespaceStatus
  }
  
  interface NamespaceApiResponse {
    code: number
    data: {
      metadata?: { resourceVersion: string }
      items: K8sNamespace[]
    }
    message: string
  }
  
  // --- Simplified interface for table display ---
  interface NamespaceDisplayItem {
    uid: string
    name: string
    status: string
    creationTimestamp: string
    labels?: { [key: string]: string } // Keep labels if needed for logic/display
  }
  
  // --- Reactive State ---
  const allNamespaces = ref<NamespaceDisplayItem[]>([])
  const currentPage = ref(1)
  const pageSize = ref(12)
  const searchQuery = ref("")
  const filterStatus = ref("")
  const viewMode = ref<'card' | 'list'>('card')
  const loading = ref(false)
  const isDialogVisible = ref(false)
  const createLoading = ref(false)
  
  const formRef = ref<FormInstance>()
  const form = reactive<{ name: string; labels?: string }>({
    name: ""
  });
  
  const systemNamespaces = ['kube-system', 'kube-public', 'kube-node-lease', 'default']
  
  // --- Form Rules ---
  const validateNamespaceName = (rule: any, value: string, callback: any) => {
    if (!value) {
      return callback(new Error('命名空间名称不能为空'))
    }
    // Kubernetes DNS-1123 Label Names:
    // - contain at most 63 characters
    // - contain only lowercase alphanumeric characters or '-'
    // - start with an alphanumeric character
    // - end with an alphanumeric character
    const regex = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/;
    if (value.length > 63) {
        return callback(new Error('名称不能超过 63 个字符'))
    }
    if (!regex.test(value)) {
        return callback(new Error('名称必须由小写字母、数字或 "-" 组成，并以字母或数字开头和结尾'))
    }
    callback()
  }
  
  const formRules = reactive<FormRules>({
    name: [{ required: true, validator: validateNamespaceName, trigger: 'blur' }]
  })
  
  // --- Computed Properties ---
  const filteredData = computed(() => {
    const query = searchQuery.value.trim().toLowerCase()
    if (!query) {
      return allNamespaces.value
    }
    return allNamespaces.value.filter((ns) =>
      ns.name.toLowerCase().includes(query)
    )
  })
  
  const totalNamespaces = computed(() => filteredData.value.length)
  
  const paginatedData = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return filteredData.value.slice(start, end)
  })
  
  // --- Helper Functions ---
  const formatTimestamp = (timestamp: string): string => {
      return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss")
  }
  
  const getStatusTagType = (status: string): 'success' | 'warning' | 'danger' => {
      if (status === 'Active') return 'success';
      if (status === 'Terminating') return 'warning';
      return 'danger'; // For unknown or other states
  }
  
  const isSystemNamespace = (name: string): boolean => {
      return systemNamespaces.includes(name);
  }
  
  // --- Cluster Store Integration ---
  const clusterStore = useClusterStore()
  const { selectedClusterName, availableClusters, loadingClusters: s_loadingClusters } = storeToRefs(clusterStore)
  
  // --- API Interaction ---
  const fetchNamespaceData = async () => {
    if (!selectedClusterName.value) {
      allNamespaces.value = []
      return
    }
    if (loading.value) return;
    loading.value = true
    try {
      const response = await kubernetesRequest<NamespaceApiResponse>({
        url: `/api/v1/namespaces`,
        method: "get"
      })
      if (response.code === 200 && response.data?.items) {
        allNamespaces.value = response.data.items.map(item => ({
          uid: item.metadata.uid,
          name: item.metadata.name,
          status: item.status.phase,
          creationTimestamp: formatTimestamp(item.metadata.creationTimestamp),
          labels: item.metadata.labels
        }))
        if (currentPage.value > Math.ceil(totalNamespaces.value / pageSize.value)) {
          currentPage.value = 1;
        }
      } else {
        ElMessage.error(`获取命名空间数据失败: ${response.message || '未知错误'}`)
        allNamespaces.value = []
      }
    } catch (error: any) {
      console.error("获取命名空间数据失败:", error)
      ElMessage.error(`获取命名空间数据出错: ${error.message || '网络请求失败'}`)
      allNamespaces.value = []
    } finally {
      loading.value = false
    }
  }
  
  const createNamespace = async () => {
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
      if (valid) {
        createLoading.value = true
        try {
          const payload = {
            apiVersion: 'v1',
            kind: 'Namespace',
            metadata: {
              name: form.name,
            }
          };
          const response = await kubernetesRequest<{ code: number; message: string }>({
            url: `/api/v1/namespaces`,
            method: "post",
            data: payload
          })
          if (response.code === 201 || response.code === 200) {
            ElMessage.success(`命名空间 \"${form.name}\" 创建成功`)
            isDialogVisible.value = false
            await fetchNamespaceData()
          } else {
            ElMessage.error(`命名空间创建失败: ${response.message || '未知错误'}`)
          }
        } catch (error: any) {
          console.error("命名空间创建失败:", error)
          const errMsg = error.response?.data?.message || error.message || '请求失败';
          ElMessage.error(`命名空间创建失败: ${errMsg}`)
        } finally {
          createLoading.value = false
        }
      } else {
        console.log('表单验证失败')
        return false
      }
    })
  }
  
  const handleDeleteNamespace = (namespace: NamespaceDisplayItem) => {
    if (isSystemNamespace(namespace.name)) {
      ElMessage.warning(`不能删除系统命名空间 \"${namespace.name}\"`);
      return;
    }
    ElMessageBox.confirm(
      `确定要删除命名空间 \"${namespace.name}\" 吗？此操作将删除该空间下的所有资源且不可恢复！`,
      '危险操作确认',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'error',
        confirmButtonClass: 'el-button--danger',
      }
    ).then(async () => {
      const loadingInstance = ElLoading.service({
        lock: true, text: `正在删除命名空间 ${namespace.name}...`, background: 'rgba(0, 0, 0, 0.7)'
      });
      try {
        const response = await kubernetesRequest<{ code: number; message: string }>({
          url: `/api/v1/namespaces/${namespace.name}`,
          method: "delete"
        });
        if (response.code === 200 || response.code === 202) {
          ElMessage.success(`命名空间 \"${namespace.name}\" 已删除`);
          await fetchNamespaceData();
        } else {
          ElMessage.error(`删除命名空间失败: ${response.message || '未知错误'}`);
        }
      } catch (error: any) {
        console.error("删除命名空间失败:", error);
        const errMsg = error.response?.data?.message || error.message || '请求失败';
        ElMessage.error(`删除命名空间失败: ${errMsg}`);
      } finally {
        loadingInstance.close();
      }
    }).catch(() => {
      ElMessage.info('删除操作已取消');
    });
  }
  
  // --- Watch selectedClusterName for changes and fetch data ---
  watch(selectedClusterName, (newName, oldName) => {
    if (newName) {
      fetchNamespaceData()
    } else {
      allNamespaces.value = []
    }
  }, { immediate: true })
  
  // --- Event Handlers ---
  const handlePageChange = (page: number) => {
    currentPage.value = page
  }
  
  const handleSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1 // Reset to page 1 when size changes
  }
  
  // Debounced search handler
  const handleSearchDebounced = debounce(() => {
      currentPage.value = 1; // Reset page when search query changes
      // No need to manually filter here, computed property `filteredData` handles it
  }, 300); // 300ms debounce delay
  
  const handleSearch = () => {
      // This function is kept in case you need immediate filtering (e.g., on button click)
      // But the debounced input handler is generally better UX
      currentPage.value = 1;
  }
  
  
  const showCreateDialog = () => {
    resetForm(); // Ensure form is clean before showing
    isDialogVisible.value = true
  }
  
  const resetForm = () => {
    form.name = "";
     // Reset other form fields if added (e.g., form.labels = "")
    // Use nextTick to ensure formRef is available after dialog opens/closes
    nextTick(() => {
        formRef.value?.clearValidate(); // Clear validation state
        // formRef.value?.resetFields(); // Also resets model values if needed
    });
  }
  
  const viewNamespaceDetails = (namespace: NamespaceDisplayItem) => {
      ElMessage.info(`模拟: 查看命名空间 "${namespace.name}" 的详情`);
      // router.push(`/namespaces/${namespace.name}`); // Example navigation
  }
  
  const editNamespaceMetadata = (namespace: NamespaceDisplayItem) => {
      ElMessage.info(`模拟: 编辑命名空间 "${namespace.name}" 的标签/注解`);
      // Show a different dialog or navigate to an edit page
  }

  // 新增的辅助函数
  const getStatusClass = (status: string) => {
    if (status === 'Active') return 'status-healthy'
    if (status === 'Terminating') return 'status-unhealthy'
    return 'status-unknown'
  }

  const getStatusText = (status: string) => {
    if (status === 'Active') return t('namespaceManagement.statuses.active')
    if (status === 'Terminating') return t('namespaceManagement.statuses.terminating')
    return t('namespaceManagement.statuses.unknown')
  }

  const handleNamespaceClick = (namespace: NamespaceDisplayItem) => {
    console.log('Namespace clicked:', namespace)
  }
  
  // --- Lifecycle Hooks ---
  onMounted(() => {
    fetchNamespaceData()
  })
  </script>
  
  <style scoped>
.namespace-management {
  padding: 24px;
  background: #f8fafc;
  min-height: 100vh;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 32px;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

/* 工具栏 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 16px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
}

.search-filters {
  display: flex;
  gap: 12px;
  align-items: center;
  flex: 1;
}

.search-input {
  width: 280px;
}

.filter-select {
  width: 160px;
}

.toolbar-right {
  display: flex;
  gap: 16px;
  align-items: center;
}

.view-toggle {
  border-radius: 8px;
}

/* 命名空间列表 */
.namespace-list {
  min-height: 400px;
}

.namespace-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  gap: 20px;
}

.namespace-card {
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  transition: all 0.2s ease;
  cursor: pointer;
  overflow: hidden;
}

.namespace-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: #3b82f6;
}

.card-header {
  padding: 20px 20px 16px;
  border-bottom: 1px solid #f1f5f9;
}

.namespace-info {
  margin-bottom: 12px;
}

.namespace-name {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 4px;
}

.namespace-status {
  font-size: 13px;
  color: #475569;
  font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-weight: 600;
  margin-top: 4px;
}

.card-body {
  padding: 16px 20px;
}

.namespace-meta {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.meta-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.meta-label {
  font-size: 13px;
  color: #64748b;
  font-weight: 500;
}

.meta-value {
  font-size: 13px;
  color: #1e293b;
  font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-weight: 600;
}

.card-footer {
  padding: 16px 20px;
  border-top: 1px solid #f1f5f9;
  background: white;
}

.namespace-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

/* 状态指示器 - 简洁小圆点样式 */
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.status-healthy {
  background: #10b981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.status-unhealthy {
  background: #ef4444;
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.2);
}

.status-unknown {
  background: #64748b;
  box-shadow: 0 0 0 2px rgba(100, 116, 139, 0.2);
}

/* 命名空间名称行 */
.namespace-name-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.system-badge {
  margin-left: auto;
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%) !important;
  border: none !important;
  color: white !important;
  font-weight: 600;
  box-shadow: 0 2px 4px rgba(245, 158, 11, 0.3);
}

/* 列表视图 */
.namespace-table {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}

.table-namespace-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.table-actions {
  display: flex;
  gap: 8px;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 80px 20px;
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

/* 分页 */
.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}

/* 对话框样式 */
.namespace-dialog :deep(.el-dialog__body) {
  padding: 20px 24px;
}

.namespace-form {
  max-height: 70vh;
  overflow-y: auto;
}

.form-section {
  margin-bottom: 32px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 1px solid #e2e8f0;
}

.form-item-help {
  font-size: 12px;
  color: #64748b;
  line-height: 1.4;
  margin-top: 4px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 0 0;
  border-top: 1px solid #e2e8f0;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .namespace-grid {
    grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  }
}

@media (max-width: 768px) {
  .namespace-management {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .page-title {
    font-size: 28px;
  }

  .toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }

  .search-filters {
    flex-direction: column;
    align-items: stretch;
  }

  .search-input,
  .filter-select {
    width: 100%;
  }

  .toolbar-right {
    justify-content: space-between;
    width: 100%;
  }

  .namespace-grid {
    grid-template-columns: 1fr;
  }

  .meta-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .namespace-actions {
    flex-direction: column;
  }
}
  </style>