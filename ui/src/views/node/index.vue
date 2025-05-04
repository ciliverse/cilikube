<template>
    <div class="node-page-container">
      <!-- 面包屑导航 (Optional but recommended for context) -->
      <el-breadcrumb separator="/" class="page-breadcrumb">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item>集群管理</el-breadcrumb-item>
        <el-breadcrumb-item>节点列表</el-breadcrumb-item>
      </el-breadcrumb>
  
      <!-- 页面头部：标题和操作按钮 -->
      <div class="page-header">
        <h1 class="page-title">Kubernetes 集群节点</h1>
        <div class="action-buttons">
          <el-button :icon="Plus" type="primary" @click="handleCreateNode" :loading="loading">新建节点</el-button>
          <el-button :icon="View" @click="handleMonitor" :loading="loading">集群监控</el-button>
          <el-dropdown>
            <el-button :icon="Operation" :loading="loading">
              更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleUpgradeCluster">集群升级</el-dropdown-item>
                <el-dropdown-item @click="handleClusterSettings">集群设置</el-dropdown-item>
                <el-dropdown-item divided @click="handleBulkAction">批量操作</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
           <el-tooltip content="刷新数据" placement="top">
             <el-button :icon="Refresh" circle @click="fetchNodeData" :loading="loading" />
           </el-tooltip>
        </div>
      </div>
  
  
      <!-- 节点信息卡片网格 -->
      <el-row :gutter="20" class="node-grid" v-loading="loading">
         <el-col v-if="!loading && !nodeData.length" :span="24">
           <el-empty description="未找到任何节点信息" />
         </el-col>
        <el-col
          v-for="node in currentPageData"
          :key="node.metadata.uid"
          :xs="24"
          :sm="12"
          :lg="8"
          :xl="6"
          class="node-col"
        >
          <el-card class="node-card" shadow="hover">
            <template #header>
              <div class="node-card__header">
                <div class="node-card__title-group">
                   <el-icon class="node-icon" :color="getNodeStatusColor(node.status.conditions)">
                      <Platform v-if="isControlPlaneNode(node)" />
                      <Monitor v-else />
                   </el-icon>
                  <span class="node-card__name">{{ node.metadata.name }}</span>
                  <el-tag
                    :type="getNodeStatusTagType(node.status.conditions)"
                    size="small"
                    effect="light"
                    class="node-card__status-tag"
                  >
                   <el-icon class="status-icon">
                       <CircleCheck v-if="getNodeStatus(node.status.conditions) === 'Ready'" />
                       <CircleClose v-else />
                   </el-icon>
                    {{ getNodeStatus(node.status.conditions) }}
                  </el-tag>
                   <el-tag v-if="isControlPlaneNode(node)" type="info" size="small" effect="plain" class="node-card__role-tag">
                      Control Plane
                   </el-tag>
                </div>
                <el-dropdown trigger="click" @command="handleNodeAction($event, node)">
                  <el-button text :icon="MoreFilled" class="node-card__actions-trigger" />
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="details" :icon="InfoFilled">查看详情</el-dropdown-item>
                      <el-dropdown-item command="monitor" :icon="DataLine">节点监控</el-dropdown-item>
                      <el-dropdown-item command="cordon" :icon="Lock">设为不可调度</el-dropdown-item>
                      <el-dropdown-item command="drain" :icon="DeleteLocation">驱逐Pods</el-dropdown-item>
                      <el-dropdown-item command="edit" :icon="EditPen" divided>编辑标签/注解</el-dropdown-item>
                      <el-dropdown-item command="delete" :icon="Delete" style="color: var(--el-color-danger)">删除节点</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </template>
            <div class="node-card__body">
              <el-descriptions :column="1" size="small" border class="node-descriptions">
                <el-descriptions-item>
                  <template #label><el-icon><Iphone /></el-icon> IP 地址</template>
                   {{ getNodeAddress(node.status.addresses) }}
                </el-descriptions-item>
                 <el-descriptions-item>
                   <template #label><el-icon><Odometer /></el-icon> Kubelet 版本</template>
                   {{ node.status.nodeInfo.kubeletVersion }}
                 </el-descriptions-item>
                <el-descriptions-item>
                  <template #label><el-icon><Cpu /></el-icon> CPU 容量</template>
                   <el-tag type="info" size="small" effect="plain">{{ node.status.capacity.cpu }} 核</el-tag>
                </el-descriptions-item>
                <el-descriptions-item>
                  <template #label><el-icon><SetUp /></el-icon> 内存容量</template>
                  <el-tag type="success" size="small" effect="plain">{{ formatMemory(node.status.capacity.memory) }}</el-tag>
                </el-descriptions-item>
                 <el-descriptions-item>
                   <template #label><el-icon><Box /></el-icon> Pods 容量</template>
                   {{ node.status.capacity.pods }}
                 </el-descriptions-item>
                 <el-descriptions-item>
                   <template #label><el-icon><MessageBox /></el-icon> 容器运行时</template>
                   <div class="runtime-info">{{ node.status.nodeInfo.containerRuntimeVersion }}</div>
                 </el-descriptions-item>
                <el-descriptions-item>
                  <template #label><el-icon><Calendar /></el-icon> 创建时间</template>
                   {{ formatTimestamp(node.metadata.creationTimestamp) }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
            <!-- <div class="node-card__footer">
               Optional Footer Content
            </div> -->
          </el-card>
        </el-col>
      </el-row>
  
      <!-- 分页 -->
      <div class="pagination-container" v-if="totalNodes > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[6, 12, 18, 24]"
          :total="totalNodes"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
          :disabled="loading"
        />
      </div>
  
       <!-- 可以选择性保留或移除 Table -->
       <!-- <div class="table-section">
          <h2>节点详细列表</h2>
          <el-table :data="currentPageData" style="width: 100%" border stripe v-loading="loading">
            <el-table-column prop="metadata.name" label="节点名称" min-width="150" fixed sortable />
            <el-table-column label="状态" min-width="100" sortable :sort-method="(a, b) => getNodeStatusSort(a.status.conditions, b.status.conditions)">
               <template #default="{ row }">
                 <el-tag :type="getNodeStatusTagType(row.status.conditions)" size="small">
                   {{ getNodeStatus(row.status.conditions) }}
                 </el-tag>
               </template>
            </el-table-column>
             <el-table-column label="角色" min-width="120">
               <template #default="{ row }">
                  <el-tag v-if="isControlPlaneNode(row)" type="info" size="small">Control Plane</el-tag>
                   <el-tag v-else type="primary" size="small" effect="plain">Worker</el-tag>
               </template>
             </el-table-column>
            <el-table-column label="IP地址" min-width="150">
               <template #default="{ row }">
                 {{ getNodeAddress(row.status.addresses) }}
               </template>
            </el-table-column>
             <el-table-column prop="status.nodeInfo.kubeletVersion" label="Kubelet版本" min-width="120" sortable />
             <el-table-column prop="status.capacity.cpu" label="CPU" min-width="80" align="center" sortable />
             <el-table-column label="内存" min-width="120" sortable :sort-method="(a, b) => parseMemory(a.status.capacity.memory) - parseMemory(b.status.capacity.memory)">
                 <template #default="{ row }">
                    {{ formatMemory(row.status.capacity.memory) }}
                 </template>
             </el-table-column>
            <el-table-column prop="status.nodeInfo.osImage" label="操作系统" min-width="200" show-overflow-tooltip />
            <el-table-column label="创建时间" min-width="170" sortable>
               <template #default="{ row }">
                  {{ formatTimestamp(row.metadata.creationTimestamp) }}
               </template>
            </el-table-column>
             <el-table-column label="操作" width="100" fixed="right" align="center">
                <template #default="{ row }">
                  <el-tooltip content="详情" placement="top">
                     <el-button link type="primary" :icon="View" @click="handleNodeAction('details', row)" />
                  </el-tooltip>
                  <el-tooltip content="更多" placement="top">
                     <el-dropdown trigger="click" @command="handleNodeAction($event, row)" style="margin-left: 10px;">
                         <el-button link type="primary" :icon="MoreFilled" />
                        <template #dropdown>
                          <el-dropdown-menu>
                              <el-dropdown-item command="monitor" :icon="DataLine">节点监控</el-dropdown-item>
                              <el-dropdown-item command="cordon" :icon="Lock">设为不可调度</el-dropdown-item>
                              <el-dropdown-item command="drain" :icon="DeleteLocation">驱逐Pods</el-dropdown-item>
                              <el-dropdown-item command="edit" :icon="EditPen" divided>编辑标签/注解</el-dropdown-item>
                              <el-dropdown-item command="delete" :icon="Delete" style="color: var(--el-color-danger)">删除节点</el-dropdown-item>
                          </el-dropdown-menu>
                        </template>
                     </el-dropdown>
                  </el-tooltip>
                </template>
             </el-table-column>
          </el-table>
       </div> -->
  
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted } from "vue"
  import { ElMessage, ElMessageBox } from "element-plus"
  import { request } from "@/utils/service" // Ensure this path is correct
  import dayjs from "dayjs"
  import {
    Plus, View, Operation, Refresh, ArrowDown, Platform, Monitor, Cpu, SetUp, Calendar, Iphone,
    CircleCheck, CircleClose, MoreFilled, InfoFilled, DataLine, Lock, DeleteLocation, EditPen, Delete, MessageBox, Box
  } from '@element-plus/icons-vue'
  
  // --- Interfaces (copied from original, assuming they match API) ---
  interface NodeMetadata {
    name: string
    uid: string
    creationTimestamp: string
    labels?: { [key: string]: string } // Added labels for role detection
    annotations?: { [key: string]: string }
    resourceVersion?: string
  }
  
  interface NodeCondition {
    type: string
    status: string
    lastHeartbeatTime?: string
    lastTransitionTime?: string
    reason?: string
    message?: string
  }
  
  interface NodeAddress {
    type: string
    address: string
  }
  
  interface NodeSystemInfo {
    machineID?: string
    systemUUID?: string
    bootID?: string
    kernelVersion: string
    osImage: string
    containerRuntimeVersion: string
    kubeletVersion: string
    kubeProxyVersion: string
    operatingSystem?: string
    architecture?: string
  }
  
  interface NodeStatus {
    capacity: {
      cpu: string
      memory: string
      pods: string
      [key: string]: string // For ephemeral-storage etc.
    }
    allocatable: {
      cpu: string
      memory: string
       pods: string
      [key: string]: string
    }
    conditions: NodeCondition[]
    addresses: NodeAddress[]
    nodeInfo: NodeSystemInfo
    // daemonEndpoints, images etc. can be added if needed
  }
  
  interface NodeSpec {
    podCIDR: string
    podCIDRs: string[]
    providerID?: string
    taints?: any[] // Define taint interface if needed
    unschedulable?: boolean
  }
  
  interface Node {
    metadata: NodeMetadata
    spec: NodeSpec
    status: NodeStatus
  }
  
  interface NodeApiResponse {
    code: number
    data: {
      items: Node[]
      metadata?: { // API might include list metadata
          resourceVersion?: string
      }
    }
    message: string
  }
  
  // --- Reactive State ---
  const nodeData = ref<Node[]>([])
  const currentPage = ref(1)
  const pageSize = ref(12) // Adjust page size for card layout
  const totalNodes = ref(0)
  const loading = ref(false)
  
  // --- Computed Properties ---
  const currentPageData = computed(() => {
    // Client-side pagination (if API doesn't support pagination)
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return nodeData.value.slice(start, end)
  
    // If API supports pagination, `nodeData` would hold only the current page's items
    // return nodeData.value;
  })
  
  const VITE_API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://192.168.1.100:8080";
  // --- API Fetching ---
  const fetchNodeData = async () => {
    if (loading.value) return;
    loading.value = true
    try {
      // Adjust API endpoint and parameters if your API supports pagination/filtering
      const response = await request<NodeApiResponse>({
        url: "/api/v1/nodes", // Make sure this endpoint is correct
        method: "get",
        // Example parameters if API supports pagination:
        // params: {
        //   page: currentPage.value,
        //   limit: pageSize.value,
        // }
        baseURL: VITE_API_BASE_URL // Keep if needed, but prefer configuring in request util
      })
  
      if (response.code === 200 && response.data?.items) {
        nodeData.value = response.data.items
        // If API returns total count:
        // totalNodes.value = response.data.totalCount || response.data.items.length;
        // If no total count, rely on items length (for client-side pagination):
        totalNodes.value = response.data.items.length
         // Reset to page 1 if current page becomes invalid after data refresh
         if (currentPage.value > Math.ceil(totalNodes.value / pageSize.value)) {
             currentPage.value = 1;
         }
  
      } else {
        ElMessage.error(`获取节点数据失败: ${response.message || '未知错误'}`)
        nodeData.value = [] // Clear data on error
        totalNodes.value = 0
      }
    } catch (error: any) {
      console.error("获取节点数据出错:", error)
      ElMessage.error(`获取节点数据出错: ${error.message || '网络请求失败'}`)
      nodeData.value = [] // Clear data on error
      totalNodes.value = 0
    } finally {
      loading.value = false
    }
  }
  
  // --- Helper Functions ---
  
  const getNodeStatus = (conditions: NodeCondition[]): 'Ready' | 'NotReady' | 'Unknown' => {
    const readyCondition = conditions.find(condition => condition.type === "Ready")
    if (!readyCondition) return 'Unknown'
    return readyCondition.status === "True" ? "Ready" : "NotReady"
  }
  
  const getNodeStatusTagType = (conditions: NodeCondition[]): 'success' | 'danger' | 'warning' => {
    const status = getNodeStatus(conditions);
    if (status === 'Ready') return 'success';
    if (status === 'NotReady') return 'danger';
    return 'warning'; // For Unknown
  }
  
  const getNodeStatusColor = (conditions: NodeCondition[]): string => {
      const status = getNodeStatus(conditions);
      if (status === 'Ready') return 'var(--el-color-success)';
      if (status === 'NotReady') return 'var(--el-color-danger)';
      return 'var(--el-color-warning)';
  }
  
  // Check for control plane role label
  const isControlPlaneNode = (node: Node): boolean => {
      const labels = node.metadata.labels || {};
      // Common labels indicating control plane role
      return 'node-role.kubernetes.io/control-plane' in labels ||
             'node-role.kubernetes.io/master' in labels; // Older label
  };
  
  
  const getNodeAddress = (addresses: NodeAddress[]): string => {
    const internalIP = addresses.find(address => address.type === "InternalIP")
    if (internalIP) return internalIP.address
    const externalIP = addresses.find(address => address.type === "ExternalIP")
     if (externalIP) return externalIP.address
    const hostname = addresses.find(address => address.type === "Hostname")
    return hostname ? hostname.address : "N/A"
  }
  
  // Parses memory string (e.g., "16328852Ki") into KiB
  const parseMemory = (memory: string | undefined): number => {
     if (!memory) return 0;
     const value = parseFloat(memory);
     if (isNaN(value)) return 0;
  
     if (memory.includes("Gi")) return value * 1024 * 1024;
     if (memory.includes("Mi")) return value * 1024;
     if (memory.includes("Ki")) return value;
     return value / 1024; // Assume bytes if no unit, convert to Ki (adjust if needed)
  }
  
  
  const formatMemory = (memory: string | undefined): string => {
      if (!memory) return 'N/A';
      const memoryInKi = parseMemory(memory);
      if (memoryInKi === 0 && memory !== '0') return 'N/A'; // Handle parsing errors
  
      const gb = memoryInKi / 1024 / 1024;
      if (gb >= 1) return `${gb.toFixed(1)} GiB`;
      const mb = memoryInKi / 1024;
      if (mb >= 1) return `${mb.toFixed(0)} MiB`;
      return `${memoryInKi.toFixed(0)} KiB`;
  }
  
  const formatTimestamp = (timestamp: string | undefined): string => {
    if (!timestamp) return 'N/A';
    return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss")
  }
  
  // Helper for table sorting
  const getNodeStatusSort = (condA: NodeCondition[], condB: NodeCondition[]) => {
      const statusA = getNodeStatus(condA);
      const statusB = getNodeStatus(condB);
      if (statusA === statusB) return 0;
      if (statusA === 'Ready') return -1;
      if (statusB === 'Ready') return 1;
      if (statusA === 'NotReady') return -1;
      if (statusB === 'NotReady') return 1;
      return 0; // Both Unknown
  }
  
  
  // --- Event Handlers ---
  
  const handlePageChange = (page: number) => {
    currentPage.value = page
    // If using server-side pagination, call fetchNodeData() here
    // fetchNodeData();
  }
  
  const handleSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1 // Reset to page 1 when size changes
    // If using server-side pagination, call fetchNodeData() here
    // fetchNodeData();
  }
  
  
  const handleCreateNode = () => {
    ElMessage.info("模拟: 跳转到新建节点页面或打开表单")
  }
  
  const handleMonitor = () => {
    ElMessage.info("模拟: 跳转到集群监控页面")
  }
  
  const handleUpgradeCluster = () => { ElMessage.info("模拟: 打开集群升级向导") }
  const handleClusterSettings = () => { ElMessage.info("模拟: 跳转到集群设置页面") }
  const handleBulkAction = () => { ElMessage.info("模拟: 启用批量操作模式") }
  
  
  const handleNodeAction = (command: string, node: Node) => {
      console.log(`Action '${command}' triggered for node:`, node.metadata.name);
      switch (command) {
          case 'details':
              ElMessageBox.alert(`显示节点 ${node.metadata.name} 的详细信息弹窗或侧边栏。`, '节点详情 (模拟)');
              break;
          case 'monitor':
              ElMessage.info(`模拟: 跳转到节点 ${node.metadata.name} 的监控页面`);
              break;
          case 'cordon':
              ElMessageBox.confirm(`确定要将节点 ${node.metadata.name} 设置为不可调度吗？`, '确认操作', { type: 'warning' })
                  .then(() => ElMessage.success(`模拟: 节点 ${node.metadata.name} 已设为不可调度`))
                  .catch(() => ElMessage.info('操作已取消'));
              break;
          case 'drain':
               ElMessageBox.confirm(`确定要驱逐节点 ${node.metadata.name} 上的 Pods 吗？此操作可能中断服务。`, '确认驱逐', { type: 'warning' })
                  .then(() => ElMessage.success(`模拟: 开始驱逐节点 ${node.metadata.name} 上的 Pods`))
                  .catch(() => ElMessage.info('操作已取消'));
              break;
          case 'edit':
              ElMessage.info(`模拟: 打开编辑节点 ${node.metadata.name} 标签/注解的表单`);
              break;
           case 'delete':
               ElMessageBox.confirm(`确定要删除节点 ${node.metadata.name} 吗？此操作不可恢复！`, '危险操作确认', { type: 'error' })
                  .then(() => ElMessage.success(`模拟: 正在删除节点 ${node.metadata.name}`))
                  .catch(() => ElMessage.info('操作已取消'));
              break;
          default:
              ElMessage.warning(`未知的操作: ${command}`);
      }
  }
  
  // --- Lifecycle Hooks ---
  onMounted(() => {
    fetchNodeData()
  })
  
  </script>
  
  <style lang="scss" scoped>
  .node-page-container {
    padding: 20px;
    background-color: var(--el-bg-color-page); // Use Element Plus background variable
  }
  
  .page-breadcrumb {
    margin-bottom: 20px;
  }
  
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 25px;
    flex-wrap: wrap; // Allow wrapping on smaller screens
    gap: 15px;
  }
  
  .page-title {
    font-size: 24px;
    font-weight: 600;
    color: var(--el-text-color-primary);
    margin: 0; // Remove default h1 margin
  }
  
  .action-buttons {
    display: flex;
    gap: 10px;
    flex-wrap: wrap; // Allow buttons to wrap
     .el-button {
        // Ensure consistent button height
        // height: 32px; // Or use default size
     }
  }
  
  .node-grid {
    margin-bottom: 25px;
  }
  
  .node-col {
    margin-bottom: 20px; // Ensure vertical space between rows of cards
    display: flex; // Ensure col takes full height for card flex
  }
  
  .node-card {
    width: 100%; // Make card fill the column
    border: 1px solid var(--el-border-color-lighter);
    transition: box-shadow 0.3s ease-in-out, border-color 0.3s ease-in-out;
    display: flex;
    flex-direction: column; // Ensure card content flows vertically
  
    &:hover {
      box-shadow: var(--el-box-shadow-light);
      border-color: var(--el-border-color-darker);
    }
  
    :deep(.el-card__header) {
      padding: 12px 15px; // Reduced header padding
      background-color: var(--el-fill-color-lighter);
      border-bottom: 1px solid var(--el-border-color-extra-light);
    }
  
    :deep(.el-card__body) {
      padding: 15px; // Consistent body padding
      flex-grow: 1; // Allow body to take remaining space
    }
  
    &__header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  
     &__title-group {
         display: flex;
         align-items: center;
         gap: 8px;
         overflow: hidden; // Prevent long names from breaking layout
     }
  
     .node-icon {
         font-size: 18px;
         flex-shrink: 0; // Prevent icon from shrinking
     }
  
    &__name {
      font-weight: 600;
      color: var(--el-text-color-primary);
      font-size: 15px;
       white-space: nowrap;
       overflow: hidden;
       text-overflow: ellipsis;
       margin-right: 5px; // Space before status tag
    }
  
    &__status-tag {
        display: inline-flex; // Align icon and text
        align-items: center;
        gap: 3px;
        font-size: 11px;
        padding: 0 5px;
        height: 20px;
        line-height: 18px;
       .status-icon {
          font-size: 10px;
       }
    }
    &__role-tag {
        font-size: 11px;
        padding: 0 5px;
        height: 20px;
        line-height: 18px;
    }
  
  
    &__actions-trigger {
        padding: 5px; // Make trigger area slightly larger
         color: var(--el-text-color-secondary);
         &:hover {
           color: var(--el-color-primary);
         }
    }
  
    &__body {
      .node-descriptions {
         :deep(.el-descriptions__label) {
            color: var(--el-text-color-secondary);
            font-weight: normal;
             font-size: 12px;
             display: inline-flex;
             align-items: center;
             gap: 4px; // Space between icon and label text
              .el-icon {
                  font-size: 14px;
              }
         }
         :deep(.el-descriptions__content) {
            color: var(--el-text-color-primary);
             font-size: 13px;
             word-break: break-all; // Break long strings like runtime version
         }
          .el-tag { // Style tags within descriptions
              font-size: 12px;
          }
          .runtime-info {
              font-family: monospace; // Use monospace for version strings
              font-size: 12px;
          }
      }
    }
  }
  
  .pagination-container {
    display: flex;
    justify-content: flex-end; // Align pagination to the right
    margin-top: 25px;
  }
  
  .table-section {
      margin-top: 30px;
      h2 {
          font-size: 18px;
          font-weight: 600;
          margin-bottom: 15px;
          color: var(--el-text-color-primary);
      }
       .el-table { // Style table if kept
          :deep(th.el-table__cell) {
              background-color: var(--el-fill-color-lighter);
              color: var(--el-text-color-secondary);
               font-weight: 600;
          }
       }
  }
  
  // Ensure dropdown icons have proper alignment
  .el-icon--right {
    margin-left: 5px;
  }
  </style>