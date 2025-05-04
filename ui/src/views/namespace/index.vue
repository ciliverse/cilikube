<template>
    <div class="namespace-page-container">
      <!-- Breadcrumbs -->
      <el-breadcrumb separator="/" class="page-breadcrumb">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item>集群管理</el-breadcrumb-item>
        <el-breadcrumb-item>命名空间</el-breadcrumb-item>
      </el-breadcrumb>
  
      <!-- Header: Title, Search, Actions -->
      <div class="page-header">
        <h1 class="page-title">命名空间 (Namespaces)</h1>
        <div class="header-actions">
          <el-input
            v-model="searchQuery"
            placeholder="搜索命名空间名称"
            :prefix-icon="SearchIcon"
            clearable
            @input="handleSearchDebounced"
            class="search-input"
          />
          <el-button type="primary" :icon="PlusIcon" @click="showCreateDialog" :loading="loading">
            新建命名空间
          </el-button>
           <el-tooltip content="刷新列表" placement="top">
             <el-button :icon="RefreshIcon" circle @click="fetchNamespaceData" :loading="loading" />
           </el-tooltip>
        </div>
      </div>
  
      <!-- Cluster Knowledge Alert -->
       <el-alert
         title="关于命名空间"
         type="info"
         show-icon
         :closable="true"
         class="info-alert"
       >
          <p>在 Kubernetes 中，名字空间（Namespace） 提供一种机制，将同一集群中的资源划分为相互隔离的组。 同一名字空间内的资源名称要唯一，但跨名字空间时没有这个要求。 名字空间作用域仅针对带有名字空间的对象（例如 Deployment、Service 等），这种作用域对集群范围的对象（例如 StorageClass、Node、PersistentVolume 等）不适用。</p>
       </el-alert>
  
      <!-- Namespace Table -->
      <el-table :data="paginatedData" v-loading="loading" style="width: 100%" border stripe class="namespace-table">
        <el-table-column prop="name" label="名称" min-width="200" sortable fixed>
          <template #default="{ row }">
            <span class="namespace-name">{{ row.name }}</span>
             <el-tooltip v-if="isSystemNamespace(row.name)" content="系统命名空间" placement="top">
               <el-icon class="system-ns-icon"><Setting /></el-icon>
             </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" min-width="120" align="center" sortable>
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small" effect="light">
              <el-icon class="status-icon">
                <CircleCheckFilled v-if="row.status === 'Active'" />
                <WarningFilled v-else />
              </el-icon>
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="creationTimestamp" label="创建时间" min-width="180" sortable />
        <el-table-column label="操作" width="150" align="center" fixed="right">
          <template #default="{ row }">
             <el-tooltip content="查看详情" placement="top">
                <el-button link type="primary" :icon="ViewIcon" @click="viewNamespaceDetails(row)"/>
             </el-tooltip>
             <el-tooltip content="编辑元数据" placement="top">
                <el-button link type="primary" :icon="EditPenIcon" @click="editNamespaceMetadata(row)"/>
             </el-tooltip>
             <el-tooltip content="删除" placement="top">
                <el-button
                  link
                  type="danger"
                  :icon="DeleteIcon"
                  @click="handleDeleteNamespace(row)"
                  :disabled="isSystemNamespace(row.name)"
                />
             </el-tooltip>
          </template>
        </el-table-column>
          <template #empty>
            <el-empty description="未找到命名空间" />
          </template>
      </el-table>
  
      <!-- Pagination -->
      <div class="pagination-container" v-if="totalNamespaces > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="totalNamespaces"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
          :disabled="loading"
        />
      </div>
  
      <!-- Create Namespace Dialog -->
      <el-dialog
        v-model="isDialogVisible"
        title="新建命名空间"
        width="500px"
        :close-on-click-modal="false"
        @closed="resetForm"
      >
        <el-form ref="formRef" :model="form" :rules="formRules" label-width="120px" label-position="right">
          <el-form-item label="名称" prop="name">
            <el-input
              v-model="form.name"
              placeholder="请输入命名空间名称 (小写字母、数字、'-')"
              clearable
            />
             <div class="form-item-help">
               名称必须符合 Kubernetes DNS-1123 标签标准。
             </div>
          </el-form-item>
           <!-- Add fields for labels/annotations if needed -->
           <!--
           <el-form-item label="标签">
             <el-input type="textarea" v-model="form.labels" placeholder="输入键值对, 例如: key1=value1, key2=value2" />
           </el-form-item>
           -->
        </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="isDialogVisible = false" :disabled="createLoading">取 消</el-button>
            <el-button type="primary" @click="createNamespace" :loading="createLoading">
              确 定
            </el-button>
          </div>
        </template>
      </el-dialog>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted, reactive, nextTick } from "vue"
  import { ElMessage, ElMessageBox, ElLoading } from "element-plus" // Keep ElLoading for create/delete potentially
  import type { FormInstance, FormRules } from 'element-plus'
  import { request } from "@/utils/service" // Ensure this path is correct
  import dayjs from "dayjs"
  import { debounce } from 'lodash-es'; // Import debounce
  
  import {
    Search as SearchIcon,
    Plus as PlusIcon,
    Refresh as RefreshIcon,
    View as ViewIcon,
    EditPen as EditPenIcon,
    Delete as DeleteIcon,
    Setting,
    CircleCheckFilled,
    WarningFilled
  } from '@element-plus/icons-vue'
  
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
  const pageSize = ref(10)
  const searchQuery = ref("")
  const loading = ref(false) // Main loading state for table/data fetch
  const isDialogVisible = ref(false)
  const createLoading = ref(false) // Loading state specifically for create button
  
  const formRef = ref<FormInstance>()
  const form = reactive<{ name: string; labels?: string }>({ // Adjust if adding labels/annotations
    name: ""
  });
  
  const systemNamespaces = ['kube-system', 'kube-public', 'kube-node-lease', 'default'] // Add others if needed
  
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
  const VITE_API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://192.168.100:8080";
  // --- API Interaction ---
  const fetchNamespaceData = async () => {
    if (loading.value) return;
    loading.value = true
    try {
      const response = await request<NamespaceApiResponse>({
        url: "/api/v1/namespace", // Corrected endpoint? Verify API doc
        method: "get",
        baseURL: VITE_API_BASE_URL, // Configure in request util if possible
      })
  
      if (response.code === 200 && response.data?.items) {
        allNamespaces.value = response.data.items.map(item => ({
          uid: item.metadata.uid,
          name: item.metadata.name,
          status: item.status.phase,
          creationTimestamp: formatTimestamp(item.metadata.creationTimestamp), // Format here
          labels: item.metadata.labels
        }))
        // Reset pagination if current page is no longer valid after filtering/refresh
         if (currentPage.value > Math.ceil(totalNamespaces.value / pageSize.value)) {
             currentPage.value = 1;
         }
  
      } else {
        ElMessage.error(`获取命名空间数据失败: ${response.message || '未知错误'}`)
        allNamespaces.value = [] // Clear data on error
      }
    } catch (error: any) {
      console.error("获取命名空间数据失败:", error)
      ElMessage.error(`获取命名空间数据出错: ${error.message || '网络请求失败'}`)
      allNamespaces.value = [] // Clear data on error
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
           // Construct the payload according to Kubernetes API spec
           const payload = {
               apiVersion: 'v1',
               kind: 'Namespace',
               metadata: {
                   name: form.name,
                   // Add labels/annotations here if form includes them
                   // labels: parseLabels(form.labels), // Example helper needed
               }
           };
  
          const response = await request<{ code: number; message: string }>({
            url: "/api/v1/namespaces", // POST to the collection endpoint
            method: "post",
            // baseURL: "VITE_API_BASE_URL",
            data: payload // Send the structured K8s object
          })
  
          if (response.code === 201 || response.code === 200) { // Check for 201 Created or 200 OK
            ElMessage.success(`命名空间 "${form.name}" 创建成功`)
            isDialogVisible.value = false
            await fetchNamespaceData() // Refresh the list
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
          ElMessage.warning(`不能删除系统命名空间 "${namespace.name}"`);
          return;
      }
  
      ElMessageBox.confirm(
          `确定要删除命名空间 "${namespace.name}" 吗？此操作将删除该空间下的所有资源且不可恢复！`,
          '危险操作确认',
          {
              confirmButtonText: '确认删除',
              cancelButtonText: '取消',
              type: 'error',
              confirmButtonClass: 'el-button--danger',
          }
      ).then(async () => {
          const loadingInstance = ElLoading.service({ // Use ElLoading for destructive actions
              lock: true, text: `正在删除命名空间 ${namespace.name}...`, background: 'rgba(0, 0, 0, 0.7)'
          });
          try {
              const response = await request<{ code: number; message: string }>({
                  url: `/api/v1/namespaces/${namespace.name}`, // DELETE specific resource URL
                  method: "delete",
                  // baseURL: "VITE_API_BASE_URL",
              });
  
              // Check for successful deletion status codes (200 OK or 202 Accepted)
              if (response.code === 200 || response.code === 202) {
                  ElMessage.success(`命名空间 "${namespace.name}" 已删除`);
                  await fetchNamespaceData(); // Refresh list
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
  
  // --- Lifecycle Hooks ---
  onMounted(() => {
    fetchNamespaceData()
  })
  </script>
  
  <style lang="scss" scoped>
  .namespace-page-container {
    padding: 20px;
    background-color: var(--el-bg-color-page);
  }
  
  .page-breadcrumb {
    margin-bottom: 20px;
  }
  
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 15px;
  }
  
  .page-title {
    font-size: 24px;
    font-weight: 600;
    color: var(--el-text-color-primary);
    margin: 0;
  }
  
  .header-actions {
    display: flex;
    align-items: center;
    gap: 10px;
    flex-wrap: wrap;
  }
  
  .search-input {
    width: 250px; // Adjust width as needed
     :deep(.el-input__wrapper) {
      border-radius: var(--el-border-radius-base); // Ensure consistency
     }
  }
  
  .info-alert {
    margin-bottom: 20px;
    p {
        margin: 5px 0;
        line-height: 1.6;
        font-size: 13px;
    }
  }
  
  
  .namespace-table {
     border-radius: 4px; // Subtle border radius for table
     border: 1px solid var(--el-border-color-lighter); // Add border
  
     :deep(th.el-table__cell) {
         background-color: var(--el-fill-color-lighter);
         color: var(--el-text-color-secondary);
         font-weight: 600;
     }
  
     .namespace-name {
         font-weight: 500;
         color: var(--el-text-color-regular);
     }
  
     .system-ns-icon {
         margin-left: 5px;
         color: var(--el-text-color-secondary);
         font-size: 14px;
         vertical-align: middle;
     }
  
     .status-icon {
         margin-right: 4px;
         vertical-align: middle; // Ensure icon aligns with text
         position: relative;
         top: -1px; // Fine-tune vertical alignment
     }
  }
  
  .el-table .el-button.is-link {
     font-size: 13px; // Ensure action button icons aren't too large
     padding: 4px; // Adjust padding for link buttons
     vertical-align: middle;
  }
  
  .pagination-container {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
  
  .dialog-footer {
    text-align: right;
  }
  
  .form-item-help {
      font-size: 12px;
      color: var(--el-text-color-secondary);
      line-height: 1.4;
      margin-top: 4px;
  }
  </style>