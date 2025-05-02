<template>
    <div class="pod-page-container">
      <!-- Breadcrumbs -->
      <el-breadcrumb separator="/" class="page-breadcrumb">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item>资源管理</el-breadcrumb-item>
        <el-breadcrumb-item>Pods</el-breadcrumb-item>
      </el-breadcrumb>
  
      <!-- Header: Title & Create Button -->
      <div class="page-header">
        <h1 class="page-title">Pods 列表</h1>
         <el-button type="primary" :icon="PlusIcon" @click="handleAddPod" :loading="loading.page">
           创建 Pod
         </el-button>
      </div>
  
       <!-- Cluster Knowledge Alert -->
       <el-alert
         title="关于 Pods"
         type="info"
         show-icon
         :closable="true"
         class="info-alert"
         description="在 Kubernetes 中，Pod 是最小的可部署单元。一个 Pod 表示一个运行中的进程，通常包含一个或多个容器。Pod 提供了容器之间共享的网络和存储资源，并规定了这些容器应如何运行。"
       />
  
      <!-- Filter Bar: Namespace, Search, Refresh -->
      <div class="filter-bar">
          <el-select
              v-model="selectedNamespace"
              placeholder="选择命名空间"
              @change="handleNamespaceChange"
              filterable
              clearable
              :loading="loading.namespaces"
              class="filter-item namespace-select"
          >
              <el-option label="所有命名空间" value="" /> <!-- Option for all namespaces -->
              <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
          </el-select>
  
          <el-input
              v-model="searchQuery"
              placeholder="搜索 Pod 名称 / IP / 节点"
              :prefix-icon="SearchIcon"
              clearable
              @input="handleSearchDebounced"
              class="filter-item search-input"
          />
  
          <el-tooltip content="刷新列表" placement="top">
              <el-button :icon="RefreshIcon" circle @click="fetchPodData" :loading="loading.pods" />
          </el-tooltip>
      </div>
  
      <!-- Pods Table -->
      <el-table
          :data="paginatedData"
          v-loading="loading.pods"
          border
          stripe
          style="width: 100%"
          @sort-change="handleSortChange"
          class="pod-table"
          :default-sort="{ prop: 'createdAt', order: 'descending' }"
      >
          <el-table-column prop="name" label="名称" min-width="250" sortable="custom" fixed show-overflow-tooltip>
               <template #default="{ row }">
                  <el-icon class="pod-icon"><Tickets /></el-icon>
                  <span class="pod-name">{{ row.name }}</span>
              </template>
          </el-table-column>
          <el-table-column prop="namespace" label="命名空间" min-width="150" sortable="custom" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" min-width="130" sortable="custom" align="center">
              <template #default="{ row }">
                  <el-tag :type="getStatusTagType(row.status)" size="small" effect="light" class="status-tag">
                      <el-icon class="status-icon" :class="getSpinClass(row.status)">
                          <component :is="getStatusIcon(row.status)" />
                      </el-icon>
                      {{ row.status }}
                  </el-tag>
              </template>
          </el-table-column>
          <el-table-column prop="ip" label="Pod IP" min-width="150" sortable="custom" show-overflow-tooltip />
          <el-table-column prop="node" label="所在节点" min-width="180" sortable="custom" show-overflow-tooltip />
          <el-table-column prop="createdAt" label="创建时间" min-width="180" sortable="custom" />
          <el-table-column label="操作" width="180" align="center" fixed="right">
              <template #default="{ row }">
                  <el-tooltip content="查看日志" placement="top">
                      <el-button link type="primary" :icon="DocumentIcon" @click="viewPodLogs(row)" />
                  </el-tooltip>
                   <el-tooltip content="进入容器" placement="top">
                      <el-button link type="primary" :icon="MonitorIcon" @click="execIntoPod(row)" />
                  </el-tooltip>
                   <el-tooltip content="编辑 YAML" placement="top">
                      <el-button link type="primary" :icon="EditPenIcon" @click="editPodYaml(row)" />
                  </el-tooltip>
                  <el-tooltip content="删除" placement="top">
                      <el-button link type="danger" :icon="DeleteIcon" @click="handleDeletePod(row)" />
                  </el-tooltip>
              </template>
          </el-table-column>
           <template #empty>
            <el-empty description="未找到 Pods" />
          </template>
      </el-table>
  
      <!-- Pagination -->
      <div class="pagination-container" v-if="totalPods > 0">
          <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="totalPods"
              layout="total, sizes, prev, pager, next, jumper"
              background
              @size-change="handleSizeChange"
              @current-change="handlePageChange"
              :disabled="loading.pods"
          />
      </div>
  
      <!-- Add/Edit Dialog (Simplified - Full YAML Editor Recommended for Real Use) -->
      <el-dialog :title="dialogTitle" v-model="dialogVisible" width="600px" :close-on-click-modal="false">
        <el-alert type="warning" :closable="false" style="margin-bottom: 20px;">
          注意：直接通过表单编辑运行中的 Pod 配置通常不被推荐。建议通过更新部署（Deployment）、状态集（StatefulSet）等控制器来管理 Pod。此对话框仅为模拟。
        </el-alert>
        <el-form :model="currentPodForm" label-width="100px" ref="podFormRef">
          <el-form-item label="名称" prop="name" :rules="[{ required: true, message: 'Pod 名称不能为空', trigger: 'blur' }]">
            <el-input v-model="currentPodForm.name" :disabled="!isAdding" placeholder="请输入 Pod 名称" />
          </el-form-item>
          <el-form-item label="命名空间" prop="namespace" :rules="[{ required: true, message: '请选择命名空间', trigger: 'change' }]">
             <el-select v-model="currentPodForm.namespace" placeholder="选择命名空间" :disabled="!isAdding" filterable style="width: 100%">
                <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
             </el-select>
          </el-form-item>
           <el-form-item label="镜像" prop="image" :rules="[{ required: true, message: '镜像名称不能为空', trigger: 'blur' }]">
            <el-input v-model="currentPodForm.image" placeholder="例如：nginx:latest" />
          </el-form-item>
          <!-- Add more fields as needed (e.g., labels, ports), but keep it simple or use YAML editor -->
        </el-form>
        <template #footer>
          <div class="dialog-footer">
              <el-button @click="dialogVisible = false">取 消</el-button>
              <el-button type="primary" @click="handleSavePod" :loading="loading.dialogSave">保 存</el-button>
          </div>
        </template>
      </el-dialog>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, reactive, computed, onMounted, watch, shallowRef } from "vue"
  import { ElMessage, ElMessageBox } from "element-plus"
  import type { FormInstance } from 'element-plus'
  import { request } from "@/utils/service" // Ensure correct path
  import dayjs from "dayjs"
  import { debounce } from 'lodash-es'
  
  import {
      Plus as PlusIcon,
      Search as SearchIcon,
      Refresh as RefreshIcon,
      Tickets, // Icon for Pod
      CircleCheckFilled,
      WarningFilled,
      CloseBold,
      Loading as LoadingIcon,
      QuestionFilled,
      Document as DocumentIcon,
      Monitor as MonitorIcon,
      EditPen as EditPenIcon,
      Delete as DeleteIcon
  } from '@element-plus/icons-vue'
  
  // --- Interfaces ---
  // Based on the provided API response
  interface PodApiItem {
    name: string
    namespace: string
    labels?: { [key: string]: string }
    annotations?: { [key: string]: string } | null
    status: string // e.g., "Running", "Pending", "Succeeded", "Failed", "Unknown"
    ip?: string   // Pod IP can be missing initially
    node?: string // Node name might be missing if not scheduled
    createdAt: string // ISO Date string
    // Add other fields if your API provides them (like container info, restarts)
  }
  
  // Interface for internal table display and form editing
  interface PodDisplayItem {
    uid?: string // Optional: Add if API provides UID
    name: string
    namespace: string
    status: string
    ip: string
    node: string
    createdAt: string
    labels?: { [key: string]: string }
    // Added simplified image field for the basic dialog
    image?: string
  }
  
  
  interface PodApiResponse {
    code: number
    data: {
      items: PodApiItem[]
      total: number // Assuming API provides total count for pagination
    }
    message: string
  }
  
  // Simple Namespace list response (adjust if API is different)
  interface NamespaceListResponse {
      code: number
      data: string[] // Assuming API returns just an array of names
      message: string
  }
  
  // --- Reactive State ---
  const allPods = ref<PodDisplayItem[]>([])
  const namespaces = ref<string[]>([])
  const selectedNamespace = ref<string>("") // Empty string for "All Namespaces"
  const currentPage = ref(1)
  const pageSize = ref(10)
  const totalPods = ref(0)
  const searchQuery = ref("")
  const sortParams = reactive({ key: 'createdAt', order: 'descending' as ('ascending' | 'descending' | null) }) // Default sort
  
  const loading = reactive({
      page: false, // Overall page/initial load
      namespaces: false,
      pods: false,
      dialogSave: false,
  })
  
  // Dialog state
  const dialogVisible = ref(false)
  const dialogTitle = ref("创建 Pod")
  const isAdding = ref(true)
  const podFormRef = ref<FormInstance>()
  const currentPodForm = reactive<Partial<PodDisplayItem>>({ // Use Partial for form flexibility
      name: "",
      namespace: "default", // Default namespace for creation
      image: ""
  })
  
  // --- Computed Properties ---
  const filteredData = computed(() => {
      const query = searchQuery.value.trim().toLowerCase()
      if (!query) {
          return allPods.value;
      }
      return allPods.value.filter(pod =>
          pod.name.toLowerCase().includes(query) ||
          (pod.ip && pod.ip.toLowerCase().includes(query)) ||
          (pod.node && pod.node.toLowerCase().includes(query))
      );
  });
  
  // Client-side Sorting (If API doesn't support server-side sorting)
  const sortedData = computed(() => {
    const data = [...filteredData.value]; // Create a shallow copy
    const { key, order } = sortParams;
  
    if (key && order) {
      data.sort((a, b) => {
        let valA = a[key as keyof PodDisplayItem];
        let valB = b[key as keyof PodDisplayItem];
  
        // Handle potential undefined values and date strings
        if (key === 'createdAt') {
          valA = valA ? new Date(valA).getTime() : 0;
          valB = valB ? new Date(valB).getTime() : 0;
        } else {
           valA = valA ?? ''; // Default to empty string if undefined
           valB = valB ?? '';
        }
  
  
        let comparison = 0;
        if (valA < valB) {
          comparison = -1;
        } else if (valA > valB) {
          comparison = 1;
        }
  
        return order === 'ascending' ? comparison : -comparison;
      });
    }
    return data;
  });
  
  
  // Final data for the current page
  const paginatedData = computed(() => {
      // If using client-side pagination:
      const start = (currentPage.value - 1) * pageSize.value;
      const end = start + pageSize.value;
      // Apply sorting *before* slicing for pagination
      return sortedData.value.slice(start, end);
  
      // If using server-side pagination, `allPods` would hold only the current page data
      // return allPods.value;
  });
  
  
  // --- Helper Functions ---
  const formatTimestamp = (timestamp: string): string => {
      if (!timestamp) return 'N/A';
      return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss");
  }
  
  const getStatusTagType = (status: string): 'success' | 'warning' | 'danger' | 'info' => {
      const lowerStatus = status?.toLowerCase();
      if (lowerStatus === 'running' || lowerStatus === 'succeeded') return 'success';
      if (lowerStatus === 'pending' || lowerStatus === 'containercreating') return 'warning';
      if (lowerStatus === 'failed' || lowerStatus === 'error') return 'danger';
      return 'info'; // Unknown, Terminating etc.
  }
  
  // Map status to icons
  const getStatusIcon = (status: string) => {
      const lowerStatus = status?.toLowerCase();
      if (lowerStatus === 'running') return CircleCheckFilled;
      if (lowerStatus === 'succeeded') return CircleCheckFilled; // Consider different icon if needed
      if (lowerStatus === 'pending' || lowerStatus === 'containercreating') return LoadingIcon; // Use loading icon for pending
      if (lowerStatus === 'failed' || lowerStatus === 'error') return CloseBold;
      return QuestionFilled; // For Unknown, Terminating etc.
  }
  
  // Add spin class for loading states
  const getSpinClass = (status: string) => {
      const lowerStatus = status?.toLowerCase();
      return (lowerStatus === 'pending' || lowerStatus === 'containercreating') ? 'is-loading' : '';
  }
  
  // --- API Interaction ---
  const fetchNamespaces = async () => {
      loading.namespaces = true;
      try {
          // Assuming API returns just names based on previous example
          const response = await request<NamespaceListResponse>({
              url: "/api/v1/namespaces", // Adjust if needed
              method: "get",
              // baseURL: "VITE_API_BASE_URL",
          });
          if (response.code === 200 && Array.isArray(response.data)) {
              namespaces.value = response.data;
              // Optionally set a default namespace if none is selected and list isn't empty
               if (!selectedNamespace.value && namespaces.value.length > 0) {
                  // selectedNamespace.value = namespaces.value.find(ns => ns === 'default') || namespaces.value[0];
                   selectedNamespace.value = ""; // Default to "All"
               }
          } else {
              ElMessage.error(`获取命名空间失败: ${response.message || '数据格式错误'}`);
          }
      } catch (error: any) {
          console.error("获取命名空间失败:", error);
          ElMessage.error(`获取命名空间失败: ${error.message || '网络请求失败'}`);
      } finally {
          loading.namespaces = false;
      }
  }
  
  const fetchPodData = async () => {
      loading.pods = true;
      try {
          // Construct API parameters for filtering, pagination, sorting (if supported by backend)
          const params: Record<string, any> = {
              // Server-side pagination example:
              // page: currentPage.value,
              // limit: pageSize.value,
  
              // Server-side sorting example:
              // sort_by: sortParams.key,
              // sort_order: sortParams.order,
  
               // Server-side filtering example (if supported):
               // search: searchQuery.value.trim() || undefined
          };
  
          // Determine the URL based on selected namespace
          const url = selectedNamespace.value
              ? `/api/v1/namespaces/${selectedNamespace.value}/pods`
              : "/api/v1/pods"; // Endpoint for all pods across namespaces (verify this endpoint exists)
  
          const response = await request<PodApiResponse>({
              url: url,
              method: "get",
              params: params, // Send params for server-side processing
              baseURL: "VITE_API_BASE_URL",
          });
  
          if (response.code === 200 && response.data?.items) {
              allPods.value = response.data.items.map(item => ({
                  name: item.name,
                  namespace: item.namespace,
                  status: item.status || 'Unknown',
                  ip: item.ip || 'N/A',
                  node: item.node || 'N/A',
                  createdAt: formatTimestamp(item.createdAt),
                  labels: item.labels,
                  // Map other fields if needed
              }));
              // If API provides total count for pagination:
              totalPods.value = response.data.total;
              // If API doesn't provide total, and using client-side pagination:
              // totalPods.value = allPods.value.length; // Calculate total based on fetched data
          } else {
              ElMessage.error(`获取 Pod 数据失败: ${response.message || '未知错误'}`);
              allPods.value = [];
              totalPods.value = 0;
          }
      } catch (error: any) {
          console.error("获取 Pod 数据失败:", error);
          ElMessage.error(`获取 Pod 数据出错: ${error.message || '网络请求失败'}`);
          allPods.value = [];
          totalPods.value = 0;
      } finally {
          loading.pods = false;
      }
  }
  
  // --- Event Handlers ---
  const handleNamespaceChange = () => {
      currentPage.value = 1; // Reset page when namespace changes
      fetchPodData();
  };
  
  const handlePageChange = (page: number) => {
      currentPage.value = page;
      // fetchPodData(); // Fetch data if using server-side pagination
  };
  
  const handleSizeChange = (size: number) => {
      pageSize.value = size;
      currentPage.value = 1; // Reset page when size changes
      // fetchPodData(); // Fetch data if using server-side pagination
  };
  
  // Debounced search handler
  const handleSearchDebounced = debounce(() => {
      currentPage.value = 1; // Reset page when search query changes
      // If using server-side search, call fetchPodData() here
      // fetchPodData();
  }, 300); // 300ms delay
  
  // Handle table sorting change
  const handleSortChange = ({ prop, order }: { prop: string | null; order: 'ascending' | 'descending' | null }) => {
      sortParams.key = prop || 'createdAt'; // Default to createdAt if prop is null
      sortParams.order = order;
      currentPage.value = 1; // Reset page when sorting changes
      // If using server-side sorting, call fetchPodData() here
      // fetchPodData();
  };
  
  
  // --- Dialog and CRUD ---
  const handleAddPod = () => {
      isAdding.value = true;
      dialogTitle.value = "创建 Pod";
      // Reset form fields
      Object.assign(currentPodForm, {
          name: "",
          namespace: selectedNamespace.value || "default", // Default to selected or 'default'
          image: ""
          // Reset other fields if added
      });
      dialogVisible.value = true;
      podFormRef.value?.clearValidate(); // Clear validation on open
  };
  
  // Edit function (using YAML editor is highly recommended for real scenarios)
  const editPodYaml = (pod: PodDisplayItem) => {
       ElMessage.info(`模拟: 打开 YAML 编辑器编辑 Pod "${pod.name}"`);
       // In a real app, you would fetch the full Pod YAML and open a YAML editor component
       // Example:
       // const fullPodYaml = await fetchPodYaml(pod.namespace, pod.name);
       // openYamlEditor(fullPodYaml);
  };
  
  // Simplified Save (Not Recommended for direct Pod creation/update)
  const handleSavePod = async () => {
      if (!podFormRef.value) return;
      await podFormRef.value.validate(async (valid) => {
          if (valid) {
              loading.dialogSave = true;
               ElMessage.warning("模拟保存操作。实际应用中应通过控制器（如 Deployment）管理 Pod。");
               // Simulate save and close
               setTimeout(() => {
                   loading.dialogSave = false;
                   dialogVisible.value = false;
                   fetchPodData(); // Refresh list
               }, 500);
  
              // --- Actual API Call Logic (Use with caution) ---
              /*
              const payload = {
                  apiVersion: 'v1',
                  kind: 'Pod',
                  metadata: {
                      name: currentPodForm.name,
                      namespace: currentPodForm.namespace,
                      // Add labels if needed from form
                  },
                  spec: {
                      containers: [{
                          name: currentPodForm.name || 'container-0', // Basic container name
                          image: currentPodForm.image,
                      }]
                      // Add other necessary spec fields
                  }
              };
  
              try {
                  let response;
                  if (isAdding.value) {
                       ElMessage.info("模拟创建请求...");
                      // response = await request({
                      //     url: `/api/v1/namespaces/${payload.metadata.namespace}/pods`,
                      //     method: "post",
                      //     data: payload,
                      // });
                  } else {
                       ElMessage.info("模拟更新请求...");
                       // Updating Pods directly is complex and usually done via controllers
                       // response = await request({
                       //    url: `/api/v1/namespaces/${payload.metadata.namespace}/pods/${payload.metadata.name}`,
                       //    method: "put", // Or PATCH
                       //    data: payload,
                       // });
                  }
  
                  // Simulate success for demo
                  await new Promise(resolve => setTimeout(resolve, 500));
                  response = { code: 200, message: '操作成功' };
  
  
                  if (response.code === 200 || response.code === 201) {
                      ElMessage.success(`Pod ${isAdding.value ? '创建' : '更新'}成功 (模拟)`);
                      dialogVisible.value = false;
                      fetchPodData();
                  } else {
                      ElMessage.error(`操作失败: ${response.message || '未知错误'}`);
                  }
              } catch (error: any) {
                   console.error("保存 Pod 失败:", error);
                   const errMsg = error.response?.data?.message || error.message || '请求失败';
                   ElMessage.error(`操作失败: ${errMsg}`);
              } finally {
                   loading.dialogSave = false;
              }
              */
          } else {
              console.log('表单验证失败');
              return false;
          }
      });
  };
  
  
  const handleDeletePod = (pod: PodDisplayItem) => {
      ElMessageBox.confirm(
          `确定要删除 Pod "${pod.name}" (命名空间: ${pod.namespace}) 吗？此操作不可恢复。`,
          '确认删除',
          {
              confirmButtonText: '确认删除',
              cancelButtonText: '取消',
              type: 'warning',
          }
      ).then(async () => {
          const loadingInstance = ElLoading.service({ lock: true, text: `正在删除 Pod ${pod.name}...` });
          try {
              const response = await request<{ code: number; message: string }>({
                  url: `/api/v1/namespaces/${pod.namespace}/pods/${pod.name}`,
                  method: "delete",
                  // baseURL: "VITE_API_BASE_URL",
              });
  
              if (response.code === 200 || response.code === 202) {
                  ElMessage.success(`Pod "${pod.name}" 已删除`);
                   // Optimistic update or wait for full refresh
                   allPods.value = allPods.value.filter(p => !(p.name === pod.name && p.namespace === pod.namespace));
                   totalPods.value = allPods.value.length; // Adjust total if using client-side pagination/filtering
                  // Alternatively, call fetchPodData() for full refresh:
                  // await fetchPodData();
              } else {
                  ElMessage.error(`删除 Pod 失败: ${response.message || '未知错误'}`);
              }
          } catch (error: any) {
              console.error("删除 Pod 失败:", error);
              const errMsg = error.response?.data?.message || error.message || '请求失败';
              ElMessage.error(`删除 Pod 失败: ${errMsg}`);
          } finally {
              loadingInstance.close();
          }
      }).catch(() => {
          ElMessage.info('删除操作已取消');
      });
  };
  
  // --- Other Actions ---
  const viewPodLogs = (pod: PodDisplayItem) => {
      ElMessage.info(`模拟: 查看 Pod "${pod.name}" 的日志`);
      // Open log viewer component/dialog, passing pod name and namespace
  }
  
  const execIntoPod = (pod: PodDisplayItem) => {
       ElMessage.info(`模拟: 进入 Pod "${pod.name}" 的容器终端`);
       // Open terminal component/dialog, passing pod name and namespace
  }
  
  
  // --- Lifecycle Hooks ---
  onMounted(async () => {
      loading.page = true;
      await fetchNamespaces();
      // Fetch pods for the default selected namespace (or all if "" is default)
      await fetchPodData();
      loading.page = false;
  });
  
  // Optional: Watch selectedNamespace if you want immediate refetch without waiting for button click
  // watch(selectedNamespace, () => fetchPodData());
  
  </script>
  
  <style lang="scss" scoped>
  .pod-page-container {
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
    margin-bottom: 15px; // Reduced margin
    flex-wrap: wrap;
    gap: 15px;
  }
  
  .page-title {
    font-size: 24px;
    font-weight: 600;
    color: var(--el-text-color-primary);
    margin: 0;
  }
  
  .info-alert {
    margin-bottom: 20px;
    background-color: var(--el-color-info-light-9); // Lighter background
     :deep(.el-alert__description) {
         font-size: 13px;
         color: var(--el-text-color-regular);
         line-height: 1.6;
     }
  }
  
  .filter-bar {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 15px;
    margin-bottom: 20px;
    padding: 15px;
    background-color: var(--el-bg-color);
    border-radius: 4px;
    border: 1px solid var(--el-border-color-lighter);
  }
  
  .filter-item {
    // Base styles for filter elements
  }
  
  .namespace-select {
    width: 240px; // Adjust width
  }
  
  .search-input {
    width: 300px; // Adjust width
  }
  
  .pod-table {
     border-radius: 4px;
     border: 1px solid var(--el-border-color-lighter);
  
      :deep(th.el-table__cell) {
         background-color: var(--el-fill-color-lighter);
         color: var(--el-text-color-secondary);
         font-weight: 600;
     }
  
     .pod-icon {
         margin-right: 6px;
         color: var(--el-text-color-secondary);
         vertical-align: middle;
     }
     .pod-name {
         font-weight: 500;
         vertical-align: middle;
     }
  
     .status-tag {
       display: inline-flex;
       align-items: center;
       gap: 4px;
     }
     .status-icon {
         font-size: 12px; // Adjust icon size in tag
     }
  
     // Loading spin animation
     .is-loading {
         animation: rotating 1.5s linear infinite;
     }
  }
  
  .el-table .el-button.is-link {
     font-size: 14px; // Slightly larger icons for actions
     padding: 4px;
     margin: 0 2px; // Add small horizontal margin
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
  </style>